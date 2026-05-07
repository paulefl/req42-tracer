package aspice

import (
	"fmt"

	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
)

// Checker validates ASPICE compliance based on a traceability graph.
type Checker struct {
	registry *ProcessRegistry
	analyzer *graph.Analyzer
	config   *model.Config
}

// NewChecker creates a new ASPICE compliance checker.
func NewChecker(analyzer *graph.Analyzer, config *model.Config) *Checker {
	return &Checker{
		registry: NewProcessRegistry(),
		analyzer: analyzer,
		config:   config,
	}
}

// CheckCompliance validates ASPICE compliance against the configured processes.
func (c *Checker) CheckCompliance() *model.ASPICEReport {
	report := &model.ASPICEReport{
		Processes: make(map[string][]*model.ASPICECheckResult),
	}

	// Get configured processes or default to all
	processesToCheck := c.config.ASPICE.Processes
	if len(processesToCheck) == 0 {
		// Default to key processes
		processesToCheck = []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"}
	}

	totalCoverage := 0.0
	processCount := 0

	for _, processID := range processesToCheck {
		process := c.registry.GetProcess(processID)
		if process == nil {
			continue
		}

		bestPractices := c.registry.ListBestPracticesForProcess(processID)
		results := []*model.ASPICECheckResult{}

		for _, bp := range bestPractices {
			result := c.checkBestPractice(bp)
			results = append(results, result)
			totalCoverage += result.Coverage
		}

		report.Processes[processID] = results
		processCount += len(bestPractices)
	}

	if processCount > 0 {
		report.Overall = totalCoverage / float64(processCount)
	}

	return report
}

// checkBestPractice evaluates a single best practice.
func (c *Checker) checkBestPractice(bp *model.ASPICEBestPractice) *model.ASPICECheckResult {
	result := &model.ASPICECheckResult{
		BP:       bp,
		Evidence: []*model.TraceLink{},
		Gaps:     []string{},
	}

	switch bp.ProcessID {
	case "SWE.1":
		result = c.checkSWE1BP(bp, result)
	case "SWE.2":
		result = c.checkSWE2BP(bp, result)
	case "SWE.3":
		result = c.checkSWE3BP(bp, result)
	case "SWE.5":
		result = c.checkSWE5BP(bp, result)
	default:
		result.Status = "not-applicable"
		result.Coverage = 0
	}

	return result
}

// checkSWE1BP checks SWE.1 (Software Requirement Analysis) best practices.
func (c *Checker) checkSWE1BP(bp *model.ASPICEBestPractice, result *model.ASPICECheckResult) *model.ASPICECheckResult {
	g := c.analyzer.GetGraph()

	switch bp.ID {
	case "SWE.1.BP2":
		// Bidirectional traceability to system requirements
		// Check that requirements have architecture coverage
		totalReqs := 0
		coveredReqs := 0
		for _, req := range g.Requirements {
			totalReqs++
			for _, link := range g.Links {
				if link.FromID == req.ID && link.LinkType == "satisfied-by" {
					coveredReqs++
					result.Evidence = append(result.Evidence, link)
					break
				}
			}
		}
		if totalReqs > 0 {
			result.Coverage = float64(coveredReqs) / float64(totalReqs) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d requirements traced", coveredReqs, totalReqs))
			}
		}

	case "SWE.1.BP6":
		// Ensure testability: requirements should have test coverage
		totalReqs := 0
		testedReqs := 0
		for _, req := range g.Requirements {
			totalReqs++
			for _, link := range g.Links {
				if link.FromID == req.ID && link.LinkType == "verified-by" {
					testedReqs++
					result.Evidence = append(result.Evidence, link)
					break
				}
			}
		}
		if totalReqs > 0 {
			result.Coverage = float64(testedReqs) / float64(totalReqs) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d requirements have test coverage", testedReqs, totalReqs))
			}
		}

	case "SWE.1.BP8":
		// Bidirectional traceability between requirements and design
		totalArch := 0
		tracedArch := 0
		for _, arch := range g.ArchElements {
			if arch.Parent != "" {
				totalArch++
				if len(arch.Req) > 0 {
					tracedArch++
				}
			}
		}
		if totalArch > 0 {
			result.Coverage = float64(tracedArch) / float64(totalArch) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d architecture elements traced to requirements", tracedArch, totalArch))
			}
		}

	default:
		result.Status = "not-evaluated"
		result.Coverage = 0
	}

	return result
}

// checkSWE2BP checks SWE.2 (Software Design) best practices.
func (c *Checker) checkSWE2BP(bp *model.ASPICEBestPractice, result *model.ASPICECheckResult) *model.ASPICECheckResult {
	g := c.analyzer.GetGraph()

	switch bp.ID {
	case "SWE.2.BP4":
		// Ensure traceability to software requirements
		totalArch := len(g.ArchElements)
		tracedArch := 0
		for _, arch := range g.ArchElements {
			if len(arch.Req) > 0 {
				tracedArch++
				for _, reqID := range arch.Req {
					for _, link := range g.Links {
						if link.FromID == arch.ID && link.ToID == reqID {
							result.Evidence = append(result.Evidence, link)
							break
						}
					}
				}
			}
		}
		if totalArch > 0 {
			result.Coverage = float64(tracedArch) / float64(totalArch) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d architecture elements traced to requirements", tracedArch, totalArch))
			}
		}

	default:
		result.Status = "not-evaluated"
		result.Coverage = 0
	}

	return result
}

// checkSWE3BP checks SWE.3 (Software Unit Implementation) best practices.
func (c *Checker) checkSWE3BP(bp *model.ASPICEBestPractice, result *model.ASPICECheckResult) *model.ASPICECheckResult {
	g := c.analyzer.GetGraph()

	switch bp.ID {
	case "SWE.3.BP3":
		// Establish traceability between implementation and design
		totalArch := 0
		implArch := 0
		for _, arch := range g.ArchElements {
			if arch.Parent != "" {
				totalArch++
				if arch.Impl != "" {
					implArch++
				}
			}
		}
		if totalArch > 0 {
			result.Coverage = float64(implArch) / float64(totalArch) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d architecture elements have implementation references", implArch, totalArch))
			}
		}

	default:
		result.Status = "not-evaluated"
		result.Coverage = 0
	}

	return result
}

// checkSWE5BP checks SWE.5 (Software Testing) best practices.
func (c *Checker) checkSWE5BP(bp *model.ASPICEBestPractice, result *model.ASPICECheckResult) *model.ASPICECheckResult {
	g := c.analyzer.GetGraph()

	switch bp.ID {
	case "SWE.5.BP3":
		// Establish traceability between tests and requirements
		totalTests := len(g.TestSpecs)
		tracedTests := 0
		for _, spec := range g.TestSpecs {
			if len(spec.Req) > 0 {
				tracedTests++
				for _, reqID := range spec.Req {
					for _, link := range g.Links {
						if link.FromID == spec.ID && link.ToID == reqID {
							result.Evidence = append(result.Evidence, link)
							break
						}
					}
				}
			}
		}
		if totalTests > 0 {
			result.Coverage = float64(tracedTests) / float64(totalTests) * 100
			result.Status = "satisfied"
			if result.Coverage < 100 {
				result.Status = "partial"
				result.Gaps = append(result.Gaps, fmt.Sprintf("%d/%d test specifications traced to requirements", tracedTests, totalTests))
			}
		}

	default:
		result.Status = "not-evaluated"
		result.Coverage = 0
	}

	return result
}

// GetProcessCoverage returns the overall coverage for a specific process.
func (c *Checker) GetProcessCoverage(processID string) (float64, error) {
	bps := c.registry.ListBestPracticesForProcess(processID)
	if len(bps) == 0 {
		return 0, fmt.Errorf("unknown process: %s", processID)
	}

	totalCoverage := 0.0
	for _, bp := range bps {
		result := c.checkBestPractice(bp)
		totalCoverage += result.Coverage
	}

	return totalCoverage / float64(len(bps)), nil
}
