package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"loadtest/pkg/config"
	"loadtest/pkg/histogram"
	"loadtest/pkg/loadgen"
	"loadtest/pkg/monitor"
	"loadtest/pkg/results"
	"loadtest/pkg/scenarios"
)

var (
	logger *zap.Logger
	cfg    *config.Config
)

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	cfg = config.NewConfig()
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "loadtest",
		Short: "A load testing CLI tool",
		Long:  `A flexible and powerful load testing CLI tool implemented in Golang.`,
		Run:   runLoadTest,
	}

	cfg.AddFlags(rootCmd)
	rootCmd.Flags().String("scenario-file", "", "Path to the JSON file containing test scenarios")

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal("Failed to execute root command", zap.Error(err))
	}
}

func runLoadTest(cmd *cobra.Command, args []string) {
	if err := cfg.LoadConfig(); err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigCh
		logger.Info("Received interrupt signal, shutting down...")
		cancel()
	}()

	// Initialize components
	monitor := monitor.NewMonitor(logger)
	collector := results.NewCollector(logger)

	// Load scenarios if a scenario file is provided
	var scenarioManager *scenarios.ScenarioManager
	scenarioFile, _ := cmd.Flags().GetString("scenario-file")
	if scenarioFile != "" {
		var err error
		scenarioManager, err = scenarios.NewScenarioManager(scenarioFile)
		if err != nil {
			logger.Fatal("Failed to load scenarios", zap.Error(err))
		}
		logger.Info("Loaded test scenarios", zap.Int("count", len(scenarioManager.Scenarios)))
	}

	loadGen := loadgen.NewLoadGenerator(cfg, logger, monitor, collector, scenarioManager)

	// Start the load test
	fmt.Println("Starting load test...")
	start := time.Now()
	if err := loadGen.Run(ctx); err != nil {
		logger.Error("Load generation failed", zap.Error(err))
	}
	duration := time.Since(start)

	// Print results
	printResults(monitor, collector, duration)
}

func printResults(monitor *monitor.Monitor, collector *results.Collector, duration time.Duration) {
	fmt.Println("\nLoad Test Results:")
	fmt.Printf("Duration: %s\n", duration)

	// Print monitor stats
	monitor.PrintStats()

	// Print detailed results
	responseTimes, statusCodes, errors := collector.GetResults()

	fmt.Println("\nResponse Time Statistics:")
	printResponseTimeStats(responseTimes)

	// Add histogram output
	hist := histogram.NewHistogram(responseTimes, 20) // 20 bins
	fmt.Println("\n" + hist.String())

	fmt.Println("\nStatus Codes:")
	for code, count := range statusCodes {
		fmt.Printf("%d: %d\n", code, count)
	}

	fmt.Printf("\nTotal Errors: %d\n", len(errors))
	if len(errors) > 0 {
		fmt.Println("Error samples:")
		for i, err := range errors {
			if i >= 5 {
				break
			}
			fmt.Printf("- %s\n", err)
		}
	}
}

func printResponseTimeStats(responseTimes []float64) {
	if len(responseTimes) == 0 {
		fmt.Println("No response times recorded")
		return
	}

	sort.Float64s(responseTimes)
	count := len(responseTimes)

	min := responseTimes[0]
	max := responseTimes[count-1]
	mean := calculateMean(responseTimes)
	median := calculateMedian(responseTimes)
	p50 := calculatePercentile(responseTimes, 50)
	p95 := calculatePercentile(responseTimes, 95)
	p99 := calculatePercentile(responseTimes, 99)

	fmt.Printf("Min: %.2f ms\n", min)
	fmt.Printf("Max: %.2f ms\n", max)
	fmt.Printf("Mean: %.2f ms\n", mean)
	fmt.Printf("Median: %.2f ms\n", median)
	fmt.Printf("50th percentile: %.2f ms\n", p50)
	fmt.Printf("95th percentile: %.2f ms\n", p95)
	fmt.Printf("99th percentile: %.2f ms\n", p99)
}

func calculateMean(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func calculateMedian(data []float64) float64 {
	count := len(data)
	if count%2 == 0 {
		return (data[count/2-1] + data[count/2]) / 2
	}
	return data[count/2]
}

func calculatePercentile(data []float64, percentile float64) float64 {
	index := int(percentile / 100 * float64(len(data)-1))
	return data[index]
}
