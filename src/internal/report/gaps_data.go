package report

import "github.com/paulefl/req42-tracer/src/internal/model"

// GapItem is a single entry in a gap category.
type GapItem struct {
	ID   string `json:"id"`
	Title string `json:"title"`
	Info string `json:"info"`
}

// GapsData holds all gap categories for the HTML Gaps-Tab.
type GapsData struct {
	OrphanRequirements     []GapItem `json:"orphan_requirements"`
	OrphanArchElements     []GapItem `json:"orphan_arch_elements"`
	OrphanDesignElements   []GapItem `json:"orphan_design_elements"`
	UntestedArchElements   []GapItem `json:"untested_arch_elements"`
	UntestedDesignElements []GapItem `json:"untested_design_elements"`
	OrphanTestSpecs        []GapItem `json:"orphan_test_specs"`
	UntracedTestResults    []GapItem `json:"untraced_test_results"`
	HasGaps                bool      `json:"has_gaps"`
}

// BuildGapsData converts a GapAnalysisResult into the JSON-serializable GapsData.
func BuildGapsData(gaps *model.GapAnalysisResult) *GapsData {
	d := &GapsData{
		OrphanRequirements:     make([]GapItem, 0, len(gaps.OrphanRequirements)),
		OrphanArchElements:     make([]GapItem, 0, len(gaps.OrphanArchElements)),
		OrphanDesignElements:   make([]GapItem, 0, len(gaps.OrphanDesignElements)),
		UntestedArchElements:   make([]GapItem, 0, len(gaps.UntestedArchElements)),
		UntestedDesignElements: make([]GapItem, 0, len(gaps.UntestedDesignElements)),
		OrphanTestSpecs:        make([]GapItem, 0, len(gaps.OrphanTestSpecs)),
		UntracedTestResults:    make([]GapItem, 0, len(gaps.UntracedTestResults)),
	}

	for _, r := range gaps.OrphanRequirements {
		d.OrphanRequirements = append(d.OrphanRequirements, GapItem{ID: r.ID, Title: r.Title, Info: r.Priority})
	}
	for _, a := range gaps.OrphanArchElements {
		d.OrphanArchElements = append(d.OrphanArchElements, GapItem{ID: a.ID, Title: a.Title, Info: a.ASPICE})
	}
	for _, dsn := range gaps.OrphanDesignElements {
		d.OrphanDesignElements = append(d.OrphanDesignElements, GapItem{ID: dsn.ID, Title: dsn.Title, Info: dsn.Arch})
	}
	for _, a := range gaps.UntestedArchElements {
		d.UntestedArchElements = append(d.UntestedArchElements, GapItem{ID: a.ID, Title: a.Title, Info: a.ASPICE})
	}
	for _, dsn := range gaps.UntestedDesignElements {
		d.UntestedDesignElements = append(d.UntestedDesignElements, GapItem{ID: dsn.ID, Title: dsn.Title, Info: dsn.Arch})
	}
	for _, s := range gaps.OrphanTestSpecs {
		d.OrphanTestSpecs = append(d.OrphanTestSpecs, GapItem{ID: s.ID, Title: s.Title})
	}
	for _, r := range gaps.UntracedTestResults {
		d.UntracedTestResults = append(d.UntracedTestResults, GapItem{ID: r.ID, Title: r.TestName, Info: r.Status})
	}

	d.HasGaps = len(d.OrphanRequirements) > 0 ||
		len(d.OrphanArchElements) > 0 ||
		len(d.OrphanDesignElements) > 0 ||
		len(d.UntestedArchElements) > 0 ||
		len(d.UntestedDesignElements) > 0 ||
		len(d.OrphanTestSpecs) > 0 ||
		len(d.UntracedTestResults) > 0

	return d
}
