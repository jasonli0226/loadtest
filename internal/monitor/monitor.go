package monitor

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Monitor struct {
	logger *zap.Logger
	mutex  sync.Mutex

	totalRequests   int64
	successRequests int64
	failedRequests  int64
	totalLatency    time.Duration

	startTime time.Time
}

func NewMonitor(logger *zap.Logger) *Monitor {
	return &Monitor{
		logger:    logger,
		startTime: time.Now(),
	}
}

func (m *Monitor) RecordRequest(success bool, latency time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.totalRequests++
	if success {
		m.successRequests++
	} else {
		m.failedRequests++
	}
	m.totalLatency += latency
}

func (m *Monitor) PrintStats() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	elapsed := time.Since(m.startTime)
	rps := float64(m.totalRequests) / elapsed.Seconds()
	avgLatency := m.totalLatency / time.Duration(m.totalRequests)

	fmt.Printf("\nTest Duration: %s\n", elapsed)
	fmt.Printf("Total Requests: %d\n", m.totalRequests)
	fmt.Printf("Successful Requests: %d\n", m.successRequests)
	fmt.Printf("Failed Requests: %d\n", m.failedRequests)
	fmt.Printf("Requests per second: %.2f\n", rps)
	fmt.Printf("Average Latency: %s\n", avgLatency)
}

func (m *Monitor) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.totalRequests = 0
	m.successRequests = 0
	m.failedRequests = 0
	m.totalLatency = 0
	m.startTime = time.Now()
}
