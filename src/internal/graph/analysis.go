package graph

import (
	"fmt"
	"sync"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// Analyzer performs gap analysis and coverage reporting on a traceability graph.
// The graph must not be modified after the first call to any analysis method —
// the internal index is built lazily via sync.Once and is not invalidated on mutation.
type Analyzer struct {
	graph *model.TraceabilityGraph

	once          sync.Once
	reqToArch     map[string]map[string]struct{} // requirement ID → arch IDs referencing it
	testedArchIDs map[string]struct{}             // arch IDs with at least one verified-by link
	testedDsnIDs  map[string]struct{}             // design element IDs with at least one verified-by link
}

// NewAnalyzer creates a new graph analyzer.
func NewAnalyzer(graph *model.TraceabilityGraph) *Analyzer {
	return &Analyzer{graph: graph}
}

// GetGraph returns the underlying traceability graph.
func (a *Analyzer) GetGraph() *model.TraceabilityGraph {
	return a.graph
}

// buildIndex computes reqToArch and testedArchIDs exactly once (thread-safe).
// O(|ArchElements| × avg|Req|) + O(|Links|)
func (a *Analyzer) buildIndex() {
	a.once.Do(func() {
		a.reqToArch = make(map[string]map[string]struct{}, len(a.graph.Requirements))
		a.testedArchIDs = make(map[string]struct{})
		a.testedDsnIDs = make(map[string]struct{})

		for archID, arch := range a.graph.ArchElements {
			for _, reqID := range arch.Req {
				if a.reqToArch[reqID] == nil {
					a.reqToArch[reqID] = make(map[string]struct{})
				}
				a.reqToArch[reqID][archID] = struct{}{}
			}
		}

		for _, link := range a.graph.Links {
			if link.LinkType == "verified-by" && link.FromType == "arch" {
				a.testedArchIDs[link.FromID] = struct{}{}
			}
			if link.LinkType == "verified-by" && link.FromType == "design" {
				a.testedDsnIDs[link.FromID] = struct{}{}
			}
		}
	})
}

// AnalyzeGaps performs a comprehensive gap analysis.
func (a *Analyzer) AnalyzeGaps() *model.GapAnalysisResult {
	a.buildIndex()

	gap := &model.GapAnalysisResult{
		OrphanRequirements:     []*model.Requirement{},
		OrphanArchElements:     []*model.ArchElement{},
		OrphanTestSpecs:        []*model.TestSpec{},
		UntracedTestResults:    []*model.TestResult{},
		MissingImplementation:  []*model.ArchElement{},
		StaleTraces:            []*model.TraceLink{},
		UntestedArchElements:   []*model.ArchElement{},
		OrphanDesignElements:   []*model.DesignElement{},
		UntestedDesignElements: []*model.DesignElement{},
	}

	// O(|Requirements|) — O(1) lookup via index
	for _, req := range a.graph.Requirements {
		if len(a.reqToArch[req.ID]) == 0 {
			gap.OrphanRequirements = append(gap.OrphanRequirements, req)
		}
	}

	// O(|ArchElements|)
	for _, arch := range a.graph.ArchElements {
		if len(arch.Req) == 0 && arch.Parent != "" {
			gap.OrphanArchElements = append(gap.OrphanArchElements, arch)
		}
		if arch.Impl == "" && arch.Parent != "" {
			gap.MissingImplementation = append(gap.MissingImplementation, arch)
		}
		// SWE.5: top-level arch elements (SWE.2) need at least one integration test (arch= on test-spec)
		if arch.Parent == "" {
			if _, tested := a.testedArchIDs[arch.ID]; !tested {
				gap.UntestedArchElements = append(gap.UntestedArchElements, arch)
			}
		}
	}

	// O(|DesignElements|) — SWE.3 gap rules
	for _, dsn := range a.graph.DesignElements {
		if dsn.Arch == "" {
			gap.OrphanDesignElements = append(gap.OrphanDesignElements, dsn)
		}
		if _, tested := a.testedDsnIDs[dsn.ID]; !tested {
			gap.UntestedDesignElements = append(gap.UntestedDesignElements, dsn)
		}
	}

	// O(|TestSpecs|)
	for _, spec := range a.graph.TestSpecs {
		if len(spec.Req) == 0 && len(spec.Arch) == 0 && len(spec.Dsn) == 0 {
			gap.OrphanTestSpecs = append(gap.OrphanTestSpecs, spec)
		}
	}

	// O(|TestResults|)
	for _, result := range a.graph.TestResults {
		if result.LinkedSpec == "" && result.LinkedCode == "" {
			gap.UntracedTestResults = append(gap.UntracedTestResults, result)
		}
	}

	// O(|Links|)
	for _, link := range a.graph.Links {
		if link.Status == "stale" {
			gap.StaleTraces = append(gap.StaleTraces, link)
		}
	}

	return gap
}

// CalculateCoverage computes traceability coverage metrics.
func (a *Analyzer) CalculateCoverage() *model.CoverageReport {
	a.buildIndex()

	report := &model.CoverageReport{
		TotalRequirements: len(a.graph.Requirements),
		TotalTestResults:  len(a.graph.TestResults),
	}

	// O(|Requirements|) — O(1) index lookups
	for _, req := range a.graph.Requirements {
		archIDs := a.reqToArch[req.ID]
		if len(archIDs) > 0 {
			report.CoveredByArch++
			for archID := range archIDs {
				if _, tested := a.testedArchIDs[archID]; tested {
					report.CoveredByTests++
					break
				}
			}
		}
	}

	if report.TotalRequirements > 0 {
		report.RequirementCoverage = float64(report.CoveredByArch) / float64(report.TotalRequirements) * 100
		report.TestCoverage = float64(report.CoveredByTests) / float64(report.TotalRequirements) * 100
	}

	// O(|TestResults|)
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

	if report.TotalTestResults > 0 {
		report.PassRate = float64(report.PassedTests) / float64(report.TotalTestResults-report.SkippedTests) * 100
	}

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

// GetOrphanRequirementsByProject returns orphan requirements for a specific project.
func (a *Analyzer) GetOrphanRequirementsByProject(project string) []*model.Requirement {
	a.buildIndex()
	var orphans []*model.Requirement
	for _, req := range a.graph.Requirements {
		if req.Project == project && len(a.reqToArch[req.ID]) == 0 {
			orphans = append(orphans, req)
		}
	}
	return orphans
}

// GetCoverageByProject returns coverage metrics for a specific project.
func (a *Analyzer) GetCoverageByProject(project string) *model.CoverageReport {
	a.buildIndex()
	report := &model.CoverageReport{}

	for _, req := range a.graph.Requirements {
		if req.Project == project {
			report.TotalRequirements++
			archIDs := a.reqToArch[req.ID]
			if len(archIDs) > 0 {
				report.CoveredByArch++
				for archID := range archIDs {
					if _, ok := a.testedArchIDs[archID]; ok {
						report.CoveredByTests++
						break
					}
				}
			}
		}
	}

	if report.TotalRequirements > 0 {
		report.RequirementCoverage = float64(report.CoveredByArch) / float64(report.TotalRequirements) * 100
		report.TestCoverage = float64(report.CoveredByTests) / float64(report.TotalRequirements) * 100
	}

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

	for archID, arch := range a.graph.ArchElements {
		for _, reqID := range arch.Req {
			if _, exists := a.graph.Requirements[reqID]; !exists {
				errors = append(errors, fmt.Sprintf("Architecture %s references unknown requirement %s", archID, reqID))
			}
		}
	}

	for dsnID, dsn := range a.graph.DesignElements {
		if dsn.Arch != "" {
			if _, exists := a.graph.ArchElements[dsn.Arch]; !exists {
				errors = append(errors, fmt.Sprintf("DesignElement %s references unknown architecture %s", dsnID, dsn.Arch))
			}
		}
	}

	for specID, spec := range a.graph.TestSpecs {
		for _, archID := range spec.Arch {
			if _, exists := a.graph.ArchElements[archID]; !exists {
				errors = append(errors, fmt.Sprintf("TestSpec %s references unknown architecture %s", specID, archID))
			}
		}
		for _, dsnID := range spec.Dsn {
			if _, exists := a.graph.DesignElements[dsnID]; !exists {
				errors = append(errors, fmt.Sprintf("TestSpec %s references unknown design element %s", specID, dsnID))
			}
		}
		for _, reqID := range spec.Req {
			if _, exists := a.graph.Requirements[reqID]; !exists {
				errors = append(errors, fmt.Sprintf("TestSpec %s references unknown requirement %s", specID, reqID))
			}
		}
	}

	for codeID, code := range a.graph.TestCodes {
		if code.TestSpec != "" {
			if _, exists := a.graph.TestSpecs[code.TestSpec]; !exists {
				errors = append(errors, fmt.Sprintf("TestCode %s references unknown TestSpec %s", codeID, code.TestSpec))
			}
		}
	}

	return errors
}
