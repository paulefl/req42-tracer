package graph

import (
	"fmt"

	"github.com/paulefl/req42-tracer/internal/model"
)

// Analyzer performs gap analysis and coverage reporting on a traceability graph.
type Analyzer struct {
	graph *model.TraceabilityGraph
}

// NewAnalyzer creates a new graph analyzer.
func NewAnalyzer(graph *model.TraceabilityGraph) *Analyzer {
	return &Analyzer{graph: graph}
}

// AnalyzeGaps performs a comprehensive gap analysis.
func (a *Analyzer) AnalyzeGaps() *model.GapAnalysisResult {
	gap := &model.GapAnalysisResult{
		OrphanRequirements:    []*model.Requirement{},
		OrphanArchElements:    []*model.ArchElement{},
		OrphanTestSpecs:       []*model.TestSpec{},
		UntracedTestResults:   []*model.TestResult{},
		MissingImplementation: []*model.ArchElement{},
		StaleTraces:           []*model.TraceLink{},
	}

	// Find orphan requirements (not referenced by architecture)
	for _, req := range a.graph.Requirements {
		if !a.isRequirementCovered(req.ID) {
			gap.OrphanRequirements = append(gap.OrphanRequirements, req)
		}
	}

	// Find orphan architecture elements (not referencing requirements)
	for _, arch := range a.graph.ArchElements {
		if len(arch.Req) == 0 && arch.Parent != "" {
			// Only flag child elements without requirements
			gap.OrphanArchElements = append(gap.OrphanArchElements, arch)
		}
	}

	// Find orphan test specs (not linked to requirements or arch)
	for _, spec := range a.graph.TestSpecs {
		if len(spec.Req) == 0 && len(spec.Arch) == 0 {
			gap.OrphanTestSpecs = append(gap.OrphanTestSpecs, spec)
		}
	}

	// Find untraced test results (not linked to specs or codes)
	for _, result := range a.graph.TestResults {
		if result.LinkedSpec == "" && result.LinkedCode == "" {
			gap.UntracedTestResults = append(gap.UntracedTestResults, result)
		}
	}

	// Find architecture without implementation
	for _, arch := range a.graph.ArchElements {
		if arch.Impl == "" && arch.Parent != "" {
			// Only flag detailed elements without impl
			gap.MissingImplementation = append(gap.MissingImplementation, arch)
		}
	}

	// Find stale traces (references to outdated versions)
	for _, link := range a.graph.Links {
		if link.Status == "stale" {
			gap.StaleTraces = append(gap.StaleTraces, link)
		}
	}

	return gap
}

// CalculateCoverage computes traceability coverage metrics.
func (a *Analyzer) CalculateCoverage() *model.CoverageReport {
	report := &model.CoverageReport{
		TotalRequirements: len(a.graph.Requirements),
		TotalTestResults:  len(a.graph.TestResults),
	}

	// Count covered requirements
	for _, req := range a.graph.Requirements {
		if a.isRequirementCovered(req.ID) {
			report.CoveredByArch++
			if a.isRequirementTestedByArch(req.ID) {
				report.CoveredByTests++
			}
		}
	}

	// Calculate percentages
	if report.TotalRequirements > 0 {
		report.RequirementCoverage = float64(report.CoveredByArch) / float64(report.TotalRequirements) * 100
		report.TestCoverage = float64(report.CoveredByTests) / float64(report.TotalRequirements) * 100
	}

	// Count test results by status
	for _, result := range a.graph.TestResults {
		switch result.Status {
		case "passed":
			report.PassedTests++
		case "failed":
			report.FailedTests++
		case "skipped":
			report.SkippedTests++
		}
	}

	// Calculate pass rate
	if report.TotalTestResults > 0 {
		report.PassRate = float64(report.PassedTests) / float64(report.TotalTestResults-report.SkippedTests) * 100
	}

	// Architecture coverage (all elements referenced by at least one requirement)
	totalArch := len(a.graph.ArchElements)
	coveredArch := 0
	for _, arch := range a.graph.ArchElements {
		if len(arch.Req) > 0 {
			coveredArch++
		}
	}
	if totalArch > 0 {
		report.ArchCoverage = float64(coveredArch) / float64(totalArch) * 100
	}

	return report
}

// isRequirementCovered checks if a requirement is referenced by at least one architecture element.
func (a *Analyzer) isRequirementCovered(reqID string) bool {
	for _, arch := range a.graph.ArchElements {
		for _, ref := range arch.Req {
			if ref == reqID {
				return true
			}
		}
	}
	return false
}

// isRequirementTestedByArch checks if a requirement is transitively covered by tests.
func (a *Analyzer) isRequirementTestedByArch(reqID string) bool {
	// Find architecture elements covering this requirement
	for _, arch := range a.graph.ArchElements {
		for _, ref := range arch.Req {
			if ref == reqID {
				// Check if this arch element is tested
				for _, link := range a.graph.Links {
					if link.FromID == arch.ID && link.LinkType == "verified-by" {
						return true
					}
				}
			}
		}
	}
	return false
}

// GetOrphanRequirementsByProject returns orphan requirements for a specific project.
func (a *Analyzer) GetOrphanRequirementsByProject(project string) []*model.Requirement {
	var orphans []*model.Requirement
	for _, req := range a.graph.Requirements {
		if req.Project == project && !a.isRequirementCovered(req.ID) {
			orphans = append(orphans, req)
		}
	}
	return orphans
}

// GetCoverageByProject returns coverage metrics for a specific project.
func (a *Analyzer) GetCoverageByProject(project string) *model.CoverageReport {
	report := &model.CoverageReport{}

	for _, req := range a.graph.Requirements {
		if req.Project == project {
			report.TotalRequirements++
			if a.isRequirementCovered(req.ID) {
				report.CoveredByArch++
				if a.isRequirementTestedByArch(req.ID) {
					report.CoveredByTests++
				}
			}
		}
	}

	// Calculate percentages
	if report.TotalRequirements > 0 {
		report.RequirementCoverage = float64(report.CoveredByArch) / float64(report.TotalRequirements) * 100
		report.TestCoverage = float64(report.CoveredByTests) / float64(report.TotalRequirements) * 100
	}

	// Count test results by project
	for _, result := range a.graph.TestResults {
		if result.Project == project {
			report.TotalTestResults++
			switch result.Status {
			case "passed":
				report.PassedTests++
			case "failed":
				report.FailedTests++
			case "skipped":
				report.SkippedTests++
			}
		}
	}

	// Calculate pass rate
	activeTests := report.TotalTestResults - report.SkippedTests
	if activeTests > 0 {
		report.PassRate = float64(report.PassedTests) / float64(activeTests) * 100
	}

	return report
}

// GetLinksFor returns all trace links from or to a given artifact.
func (a *Analyzer) GetLinksFor(artifactID string) []*model.TraceLink {
	var links []*model.TraceLink
	for _, link := range a.graph.Links {
		if link.FromID == artifactID || link.ToID == artifactID {
			links = append(links, link)
		}
	}
	return links
}

// ValidateReferences checks for broken links (references to non-existent artifacts).
func (a *Analyzer) ValidateReferences() []string {
	var errors []string

	// Check requirement references in architecture
	for archID, arch := range a.graph.ArchElements {
		for _, reqID := range arch.Req {
			if _, exists := a.graph.Requirements[reqID]; !exists {
				errors = append(errors, fmt.Sprintf("Architecture %s references unknown requirement %s", archID, reqID))
			}
		}
	}

	// Check architecture references in test specs
	for specID, spec := range a.graph.TestSpecs {
		for _, archID := range spec.Arch {
			if _, exists := a.graph.ArchElements[archID]; !exists {
				errors = append(errors, fmt.Sprintf("TestSpec %s references unknown architecture %s", specID, archID))
			}
		}
	}

	// Check requirement references in test specs
	for specID, spec := range a.graph.TestSpecs {
		for _, reqID := range spec.Req {
			if _, exists := a.graph.Requirements[reqID]; !exists {
				errors = append(errors, fmt.Sprintf("TestSpec %s references unknown requirement %s", specID, reqID))
			}
		}
	}

	// Check test spec references in test codes
	for codeID, code := range a.graph.TestCodes {
		if code.TestSpec != "" {
			if _, exists := a.graph.TestSpecs[code.TestSpec]; !exists {
				errors = append(errors, fmt.Sprintf("TestCode %s references unknown TestSpec %s", codeID, code.TestSpec))
			}
		}
	}

	return errors
}
