package results

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type Collector struct {
	logger *zap.Logger
	mutex  sync.Mutex

	responseTimesMs []float64
	statusCodes     map[int]int
	errors          []string
}

func NewCollector(logger *zap.Logger) *Collector {
	return &Collector{
		logger:      logger,
		statusCodes: make(map[int]int),
	}
}

func (c *Collector) RecordResponse(statusCode int, responseTime time.Duration, err error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.responseTimesMs = append(c.responseTimesMs, float64(responseTime.Milliseconds()))
	c.statusCodes[statusCode]++

	if err != nil {
		c.errors = append(c.errors, err.Error())
	}
}

func (c *Collector) GetResults() ([]float64, map[int]int, []string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	responseTimes := make([]float64, len(c.responseTimesMs))
	copy(responseTimes, c.responseTimesMs)

	statusCodes := make(map[int]int)
	for k, v := range c.statusCodes {
		statusCodes[k] = v
	}

	errors := make([]string, len(c.errors))
	copy(errors, c.errors)

	return responseTimes, statusCodes, errors
}

func (c *Collector) Reset() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.responseTimesMs = nil
	c.statusCodes = make(map[int]int)
	c.errors = nil
}
