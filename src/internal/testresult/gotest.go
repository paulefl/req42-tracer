package testresult

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// GoTestEvent represents a single event from go test -json output.
type GoTestEvent struct {
	Time    string `json:"Time"`
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Elapsed float64 `json:"Elapsed"`
	Output  string `json:"Output"`
}

// ParseGoTest parses a go test -json output file and returns TestResult objects.
func ParseGoTest(filePath, project, platform string) ([]*model.TestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open go-test file %s: %w", filePath, err)
	}
	defer file.Close()

	// Map of test-name -> TestResult for aggregation
	tests := make(map[string]*model.TestResult)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		var event GoTestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue // Skip malformed lines
		}

		// Only process test run events
		if event.Test == "" || event.Action == "skip" {
			continue
		}

		key := event.Package + "::" + event.Test

		// Initialize test result on first encounter
		if _, exists := tests[key]; !exists {
			ts := time.Now()
			if event.Time != "" {
				if t, err := time.Parse(time.RFC3339Nano, event.Time); err == nil {
					ts = t
				}
			}

			tests[key] = &model.TestResult{
				ID:         shortPackage(event.Package) + "::" + event.Test,
				Package:    event.Package,
				TestName:   event.Test,
				FullName:   key,
				Status:     "passed", // default
				Timestamp:  ts,
				Project:    project,
				Platform:   platform,
				Attributes: make(map[string]string),
			}
		}

		result := tests[key]

		// Update based on action
		switch event.Action {
		case "run":
			result.Status = "running"
		case "pass":
			result.Status = "passed"
			result.Duration = event.Elapsed
		case "fail":
			result.Status = "failed"
			result.Duration = event.Elapsed
		case "output":
			result.Stdout += event.Output
		}
	}

	// Convert map to slice
	var results []*model.TestResult
	for _, result := range tests {
		results = append(results, result)
	}

	return results, scanner.Err()
}
