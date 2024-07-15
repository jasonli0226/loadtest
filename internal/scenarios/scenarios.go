package scenarios

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"loadtest/internal/config"
)

type Scenario struct {
	Name     string            `json:"name"`
	Endpoint string            `json:"endpoint"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
	Payload  string            `json:"payload"`
}

type ScenarioManager struct {
	Scenarios []Scenario
}

func NewScenarioManager(scenarioFile string) (*ScenarioManager, error) {
	data, err := os.ReadFile(scenarioFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read scenario file: %w", err)
	}

	var scenarios []Scenario
	if err := json.Unmarshal(data, &scenarios); err != nil {
		return nil, fmt.Errorf("failed to parse scenario file: %w", err)
	}

	return &ScenarioManager{Scenarios: scenarios}, nil
}

func (sm *ScenarioManager) GetRandomScenario() *Scenario {
	if len(sm.Scenarios) == 0 {
		return nil
	}
	return &sm.Scenarios[rand.Intn(len(sm.Scenarios))]
}

func (sm *ScenarioManager) ApplyScenario(scenario *Scenario, cfg *config.Config) {
	if scenario == nil {
		return
	}

	cfg.TargetURL = scenario.Endpoint
	cfg.HTTPMethod = scenario.Method
	cfg.CustomHeaders = scenario.Headers
	cfg.RequestPayload = scenario.Payload
}
