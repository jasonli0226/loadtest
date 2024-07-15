package loadgen

import (
	"context"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"loadtest/internal/config"
	"loadtest/internal/monitor"
	"loadtest/internal/results"
	"loadtest/internal/scenarios"
)

type LoadGenerator struct {
	config          *config.Config
	logger          *zap.Logger
	client          *fasthttp.Client
	limiter         *rate.Limiter
	monitor         *monitor.Monitor
	collector       *results.Collector
	scenarioManager *scenarios.ScenarioManager
}

func NewLoadGenerator(cfg *config.Config, logger *zap.Logger, monitor *monitor.Monitor, collector *results.Collector, scenarioManager *scenarios.ScenarioManager) *LoadGenerator {
	client := &fasthttp.Client{
		MaxConnsPerHost: cfg.ConcurrentUsers,
		ReadTimeout:     time.Duration(cfg.Timeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.Timeout) * time.Second,
	}

	return &LoadGenerator{
		config:          cfg,
		logger:          logger,
		client:          client,
		limiter:         rate.NewLimiter(rate.Limit(cfg.RequestRate), 1),
		monitor:         monitor,
		collector:       collector,
		scenarioManager: scenarioManager,
	}
}

func (lg *LoadGenerator) Run(ctx context.Context) error {
	duration, err := time.ParseDuration(lg.config.TestDuration)
	if err != nil {
		return err
	}

	// Warm-up period
	warmupDuration := duration / 10 // 10% of total duration for warm-up
	lg.logger.Debug("Starting warm-up period", zap.Duration("duration", warmupDuration))
	warmupCtx, warmupCancel := context.WithTimeout(ctx, warmupDuration)
	lg.runWorkers(warmupCtx, lg.config.ConcurrentUsers/2) // Use half the users for warm-up
	warmupCancel()

	lg.logger.Debug("Warm-up completed, starting full test")
	lg.monitor.Reset()
	lg.collector.Reset()

	testCtx, testCancel := context.WithTimeout(ctx, duration)
	defer testCancel()

	lg.runWorkers(testCtx, lg.config.ConcurrentUsers)

	return nil
}

func (lg *LoadGenerator) runWorkers(ctx context.Context, numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			lg.runWorker(ctx)
		}()
	}
	wg.Wait()
}

func (lg *LoadGenerator) runWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if err := lg.limiter.Wait(ctx); err != nil {
				lg.logger.Debug("Rate limiter error", zap.Error(err))
				return
			}
			lg.sendRequest(ctx)
		}
	}
}

func (lg *LoadGenerator) sendRequest(_ context.Context) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	// Apply scenario if available
	if lg.scenarioManager != nil {
		scenario := lg.scenarioManager.GetRandomScenario()
		lg.scenarioManager.ApplyScenario(scenario, lg.config)
	}

	req.SetRequestURI(lg.config.TargetURL)
	req.Header.SetMethod(lg.config.HTTPMethod)

	for k, v := range lg.config.CustomHeaders {
		req.Header.Set(k, v)
	}

	if lg.config.RequestPayload != "" {
		req.SetBodyString(lg.config.RequestPayload)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	start := time.Now()
	err := lg.client.DoTimeout(req, resp, time.Duration(lg.config.Timeout)*time.Second)
	duration := time.Since(start)

	if err != nil {
		lg.logger.Error("Request failed", zap.Error(err))
		lg.monitor.RecordRequest(false, duration)
		lg.collector.RecordResponse(0, duration, err)
	} else {
		lg.monitor.RecordRequest(true, duration)
		lg.collector.RecordResponse(resp.StatusCode(), duration, nil)

		// Basic response validation
		if resp.StatusCode() >= 400 {
			lg.logger.Warn("Request returned error status code", zap.Int("status_code", resp.StatusCode()))
		}
	}
}
