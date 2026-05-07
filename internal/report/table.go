package report

import (
	"fmt"
	"strings"

	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
)

// TableReporter generates text and markdown reports.
type TableReporter struct {
	analyzer *graph.Analyzer
	format   string // "text", "markdown", "json"
}

// NewTableReporter creates a new table reporter.
func NewTableReporter(analyzer *graph.Analyzer, format string) *TableReporter {
	return &TableReporter{analyzer: analyzer, format: format}
}

// TraceabilityMatrix generates a traceability matrix report.
func (tr *TableReporter) TraceabilityMatrix() string {
	g := tr.analyzer.GetGraph()

	switch tr.format {
	case "markdown":
		return tr.traceabilityMatrixMarkdown(g)
	case "json":
		return tr.traceabilityMatrixJSON(g)
	default:
		return tr.traceabilityMatrixText(g)
	}
}

// traceabilityMatrixText generates text format.
func (tr *TableReporter) traceabilityMatrixText(g *model.TraceabilityGraph) string {
	var buf strings.Builder

	buf.WriteString(strings.Repeat("=", 80) + "\n")
	buf.WriteString("TRACEABILITY MATRIX\n")
	buf.WriteString(strings.Repeat("=", 80) + "\n\n")

	// Group by requirement
	for _, req := range g.Requirements {
		buf.WriteString(fmt.Sprintf("Req: %s (%s)\n", req.ID, req.Priority))
		buf.WriteString(fmt.Sprintf("  %s\n", req.Title))

		// Find architecture covering this requirement
		covered := false
		for _, arch := range g.ArchElements {
			for _, ref := range arch.Req {
				if ref == req.ID {
					covered = true
					buf.WriteString(fmt.Sprintf("  -> Arch: %s\n", arch.ID))

					// Check if architecture is tested
					for _, link := range g.Links {
						if link.FromID == arch.ID && link.LinkType == "verified-by" {
							buf.WriteString(fmt.Sprintf("     -> Test: %s\n", link.ToID))
						}
					}
					break
				}
			}
		}

		if !covered {
			buf.WriteString("  ⚠ NOT COVERED BY ARCHITECTURE\n")
		}

		buf.WriteString("\n")
	}

	return buf.String()
}

// traceabilityMatrixMarkdown generates markdown format.
func (tr *TableReporter) traceabilityMatrixMarkdown(g *model.TraceabilityGraph) string {
	var buf strings.Builder

	buf.WriteString("# Traceability Matrix\n\n")
	buf.WriteString("| Requirement | Priority | Arch | Test | Coverage |\n")
	buf.WriteString("|-------------|----------|------|------|----------|\n")

	for _, req := range g.Requirements {
		archCovered := false
		testCovered := false
		archID := ""

		// Find coverage
		for _, arch := range g.ArchElements {
			for _, ref := range arch.Req {
				if ref == req.ID {
					archCovered = true
					archID = arch.ID

					for _, link := range g.Links {
						if link.FromID == arch.ID && link.LinkType == "verified-by" {
							testCovered = true
							break
						}
					}
					break
				}
			}
			if archCovered {
				break
			}
		}

		archStr := "❌"
		testStr := "❌"
		coverage := "0%"

		if archCovered {
			archStr = "✅"
			coverage = "50%"
		}
		if testCovered {
			testStr = "✅"
			coverage = "100%"
		}

		buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			req.ID, req.Priority, archStr, testStr, coverage))
	}

	return buf.String()
}

// traceabilityMatrixJSON generates JSON format.
func (tr *TableReporter) traceabilityMatrixJSON(g *model.TraceabilityGraph) string {
	buf := strings.Builder{}
	buf.WriteString("{\n")
	buf.WriteString("  \"matrix\": [\n")

	requirements := g.Requirements
	first := true
	for _, req := range requirements {
		if !first {
			buf.WriteString(",\n")
		}
		first = false

		archCovered := false
		testCovered := false

		for _, arch := range g.ArchElements {
			for _, ref := range arch.Req {
				if ref == req.ID {
					archCovered = true
					for _, link := range g.Links {
						if link.FromID == arch.ID && link.LinkType == "verified-by" {
							testCovered = true
							break
						}
					}
					break
				}
			}
			if archCovered {
				break
			}
		}

		buf.WriteString(fmt.Sprintf("    {\n"))
		buf.WriteString(fmt.Sprintf("      \"id\": \"%s\",\n", req.ID))
		buf.WriteString(fmt.Sprintf("      \"priority\": \"%s\",\n", req.Priority))
		buf.WriteString(fmt.Sprintf("      \"arch_covered\": %v,\n", archCovered))
		buf.WriteString(fmt.Sprintf("      \"test_covered\": %v\n", testCovered))
		buf.WriteString(fmt.Sprintf("    }"))
	}

	buf.WriteString("\n  ]\n")
	buf.WriteString("}\n")

	return buf.String()
}

// GapReport generates a gap analysis report.
func (tr *TableReporter) GapReport() string {
	gaps := tr.analyzer.AnalyzeGaps()

	switch tr.format {
	case "markdown":
		return tr.gapReportMarkdown(gaps)
	case "json":
		return tr.gapReportJSON(gaps)
	default:
		return tr.gapReportText(gaps)
	}
}

// gapReportText generates text format gap report.
func (tr *TableReporter) gapReportText(gaps *model.GapAnalysisResult) string {
	var buf strings.Builder

	buf.WriteString(strings.Repeat("=", 80) + "\n")
	buf.WriteString("GAP ANALYSIS\n")
	buf.WriteString(strings.Repeat("=", 80) + "\n\n")

	if len(gaps.OrphanRequirements) > 0 {
		buf.WriteString(fmt.Sprintf("ORPHAN REQUIREMENTS (%d):\n", len(gaps.OrphanRequirements)))
		for _, req := range gaps.OrphanRequirements {
			buf.WriteString(fmt.Sprintf("  ❌ %s: %s\n", req.ID, req.Title))
		}
		buf.WriteString("\n")
	}

	if len(gaps.OrphanArchElements) > 0 {
		buf.WriteString(fmt.Sprintf("ORPHAN ARCHITECTURE ELEMENTS (%d):\n", len(gaps.OrphanArchElements)))
		for _, arch := range gaps.OrphanArchElements {
			buf.WriteString(fmt.Sprintf("  ❌ %s: %s\n", arch.ID, arch.Title))
		}
		buf.WriteString("\n")
	}

	if len(gaps.OrphanTestSpecs) > 0 {
		buf.WriteString(fmt.Sprintf("ORPHAN TEST SPECS (%d):\n", len(gaps.OrphanTestSpecs)))
		for _, spec := range gaps.OrphanTestSpecs {
			buf.WriteString(fmt.Sprintf("  ❌ %s: %s\n", spec.ID, spec.Title))
		}
		buf.WriteString("\n")
	}

	if len(gaps.MissingImplementation) > 0 {
		buf.WriteString(fmt.Sprintf("MISSING IMPLEMENTATION (%d):\n", len(gaps.MissingImplementation)))
		for _, arch := range gaps.MissingImplementation {
			buf.WriteString(fmt.Sprintf("  ⚠ %s: %s\n", arch.ID, arch.Title))
		}
		buf.WriteString("\n")
	}

	if len(gaps.StaleTraces) > 0 {
		buf.WriteString(fmt.Sprintf("STALE TRACES (%d):\n", len(gaps.StaleTraces)))
		for _, link := range gaps.StaleTraces {
			buf.WriteString(fmt.Sprintf("  ⚠ %s -> %s (version mismatch)\n", link.FromID, link.ToID))
		}
		buf.WriteString("\n")
	}

	if len(gaps.OrphanRequirements) == 0 && len(gaps.OrphanArchElements) == 0 &&
		len(gaps.OrphanTestSpecs) == 0 && len(gaps.MissingImplementation) == 0 {
		buf.WriteString("✅ No gaps detected!\n")
	}

	return buf.String()
}

// gapReportMarkdown generates markdown format gap report.
func (tr *TableReporter) gapReportMarkdown(gaps *model.GapAnalysisResult) string {
	var buf strings.Builder

	buf.WriteString("# Gap Analysis Report\n\n")

	if len(gaps.OrphanRequirements) > 0 {
		buf.WriteString(fmt.Sprintf("## Orphan Requirements (%d)\n\n", len(gaps.OrphanRequirements)))
		for _, req := range gaps.OrphanRequirements {
			buf.WriteString(fmt.Sprintf("- **%s**: %s\n", req.ID, req.Title))
		}
		buf.WriteString("\n")
	}

	if len(gaps.OrphanArchElements) > 0 {
		buf.WriteString(fmt.Sprintf("## Orphan Architecture Elements (%d)\n\n", len(gaps.OrphanArchElements)))
		for _, arch := range gaps.OrphanArchElements {
			buf.WriteString(fmt.Sprintf("- **%s**: %s\n", arch.ID, arch.Title))
		}
		buf.WriteString("\n")
	}

	return buf.String()
}

// gapReportJSON generates JSON format gap report.
func (tr *TableReporter) gapReportJSON(gaps *model.GapAnalysisResult) string {
	buf := strings.Builder{}
	buf.WriteString("{\n")
	buf.WriteString(fmt.Sprintf("  \"orphan_requirements\": %d,\n", len(gaps.OrphanRequirements)))
	buf.WriteString(fmt.Sprintf("  \"orphan_arch_elements\": %d,\n", len(gaps.OrphanArchElements)))
	buf.WriteString(fmt.Sprintf("  \"orphan_test_specs\": %d,\n", len(gaps.OrphanTestSpecs)))
	buf.WriteString(fmt.Sprintf("  \"missing_implementation\": %d,\n", len(gaps.MissingImplementation)))
	buf.WriteString(fmt.Sprintf("  \"stale_traces\": %d\n", len(gaps.StaleTraces)))
	buf.WriteString("}\n")

	return buf.String()
}

// CoverageReport generates a coverage report.
func (tr *TableReporter) CoverageReport() string {
	coverage := tr.analyzer.CalculateCoverage()

	switch tr.format {
	case "markdown":
		return tr.coverageReportMarkdown(coverage)
	case "json":
		return tr.coverageReportJSON(coverage)
	default:
		return tr.coverageReportText(coverage)
	}
}

// coverageReportText generates text format coverage report.
func (tr *TableReporter) coverageReportText(coverage *model.CoverageReport) string {
	var buf strings.Builder

	buf.WriteString(strings.Repeat("=", 80) + "\n")
	buf.WriteString("COVERAGE REPORT\n")
	buf.WriteString(strings.Repeat("=", 80) + "\n\n")

	buf.WriteString("Requirements:\n")
	buf.WriteString(fmt.Sprintf("  Total: %d\n", coverage.TotalRequirements))
	buf.WriteString(fmt.Sprintf("  Covered by Architecture: %d (%.1f%%)\n", coverage.CoveredByArch, coverage.RequirementCoverage))
	buf.WriteString(fmt.Sprintf("  Covered by Tests: %d (%.1f%%)\n\n", coverage.CoveredByTests, coverage.TestCoverage))

	buf.WriteString("Tests:\n")
	buf.WriteString(fmt.Sprintf("  Total: %d\n", coverage.TotalTestResults))
	buf.WriteString(fmt.Sprintf("  Passed: %d\n", coverage.PassedTests))
	buf.WriteString(fmt.Sprintf("  Failed: %d\n", coverage.FailedTests))
	buf.WriteString(fmt.Sprintf("  Skipped: %d\n", coverage.SkippedTests))
	if coverage.TotalTestResults > 0 {
		buf.WriteString(fmt.Sprintf("  Pass Rate: %.1f%%\n", coverage.PassRate))
	}

	return buf.String()
}

// coverageReportMarkdown generates markdown format coverage report.
func (tr *TableReporter) coverageReportMarkdown(coverage *model.CoverageReport) string {
	var buf strings.Builder

	buf.WriteString("# Coverage Report\n\n")
	buf.WriteString("## Requirements\n\n")
	buf.WriteString(fmt.Sprintf("| Metric | Count | Coverage |\n"))
	buf.WriteString(fmt.Sprintf("|--------|-------|----------|\n"))
	buf.WriteString(fmt.Sprintf("| Total | %d | — |\n", coverage.TotalRequirements))
	buf.WriteString(fmt.Sprintf("| Covered by Architecture | %d | %.1f%% |\n", coverage.CoveredByArch, coverage.RequirementCoverage))
	buf.WriteString(fmt.Sprintf("| Covered by Tests | %d | %.1f%% |\n\n", coverage.CoveredByTests, coverage.TestCoverage))

	buf.WriteString("## Tests\n\n")
	buf.WriteString(fmt.Sprintf("| Status | Count |\n"))
	buf.WriteString(fmt.Sprintf("|--------|-------|\n"))
	buf.WriteString(fmt.Sprintf("| Passed | %d |\n", coverage.PassedTests))
	buf.WriteString(fmt.Sprintf("| Failed | %d |\n", coverage.FailedTests))
	buf.WriteString(fmt.Sprintf("| Skipped | %d |\n", coverage.SkippedTests))

	return buf.String()
}

// coverageReportJSON generates JSON format coverage report.
func (tr *TableReporter) coverageReportJSON(coverage *model.CoverageReport) string {
	buf := strings.Builder{}
	buf.WriteString("{\n")
	buf.WriteString("  \"requirements\": {\n")
	buf.WriteString(fmt.Sprintf("    \"total\": %d,\n", coverage.TotalRequirements))
	buf.WriteString(fmt.Sprintf("    \"covered_by_arch\": %d,\n", coverage.CoveredByArch))
	buf.WriteString(fmt.Sprintf("    \"coverage_pct\": %.1f\n", coverage.RequirementCoverage))
	buf.WriteString("  },\n")
	buf.WriteString("  \"tests\": {\n")
	buf.WriteString(fmt.Sprintf("    \"total\": %d,\n", coverage.TotalTestResults))
	buf.WriteString(fmt.Sprintf("    \"passed\": %d,\n", coverage.PassedTests))
	buf.WriteString(fmt.Sprintf("    \"failed\": %d,\n", coverage.FailedTests))
	buf.WriteString(fmt.Sprintf("    \"skipped\": %d,\n", coverage.SkippedTests))
	buf.WriteString(fmt.Sprintf("    \"pass_rate\": %.1f\n", coverage.PassRate))
	buf.WriteString("  }\n")
	buf.WriteString("}\n")

	return buf.String()
}
