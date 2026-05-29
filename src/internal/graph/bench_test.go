package graph

import (
	"fmt"
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// buildLargeGraph builds a graph with n requirements, n arch elements, n/2 test specs and links.
func buildLargeGraph(n int) *model.TraceabilityGraph {
	g := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement, n),
		ArchElements: make(map[string]*model.ArchElement, n+1),
		TestSpecs:    make(map[string]*model.TestSpec, n/2),
		TestCodes:    make(map[string]*model.TestCode),
		TestResults:  make(map[string]*model.TestResult, n/2),
		Links:        make([]*model.TraceLink, 0, n*2),
	}

	g.ArchElements["comp.system"] = &model.ArchElement{
		ID: "comp.system", Title: "System", Project: "bench",
	}

	for i := 0; i < n; i++ {
		reqID := fmt.Sprintf("REQ-%04d", i)
		archID := fmt.Sprintf("comp.module%04d", i)
		specID := fmt.Sprintf("TS-%04d", i)

		g.Requirements[reqID] = &model.Requirement{
			ID: reqID, Title: fmt.Sprintf("Requirement %d", i),
			Priority: "high", Status: "approved", ASPICE: "SWE.1", Project: "bench",
		}

		g.ArchElements[archID] = &model.ArchElement{
			ID: archID, Parent: "comp.system", Req: []string{reqID},
			Impl: fmt.Sprintf("internal/module%d.go", i), Project: "bench",
		}

		if i%2 == 0 {
			g.TestSpecs[specID] = &model.TestSpec{
				ID: specID, Req: []string{reqID}, Arch: []string{archID}, Project: "bench",
			}
		}

		g.Links = append(g.Links, &model.TraceLink{
			FromID: reqID, ToID: archID,
			FromType: "requirement", ToType: "arch",
			LinkType: "satisfied-by", Status: "active",
		})
		if i%2 == 0 {
			// arch→spec verified-by link: FromType must be "arch" so testedArchIDs is populated
			g.Links = append(g.Links, &model.TraceLink{
				FromID: archID, ToID: specID,
				FromType: "arch", ToType: "test-spec",
				LinkType: "verified-by", Status: "active",
			})
		}
	}

	return g
}

// BenchmarkBuildLinks measures graph link construction at scale.
func BenchmarkBuildLinks(b *testing.B) {
	for _, n := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			src := buildLargeGraph(n)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				builder := NewBuilder()
				builder.MergeGraph(src)
				builder.BuildLinks()
			}
		})
	}
}

// BenchmarkAnalyzeGaps measures gap analysis at scale.
func BenchmarkAnalyzeGaps(b *testing.B) {
	for _, n := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			g := buildLargeGraph(n)
			analyzer := NewAnalyzer(g)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				analyzer.AnalyzeGaps()
			}
		})
	}
}

// BenchmarkCalculateCoverage measures coverage calculation at scale.
func BenchmarkCalculateCoverage(b *testing.B) {
	for _, n := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			g := buildLargeGraph(n)
			analyzer := NewAnalyzer(g)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				analyzer.CalculateCoverage()
			}
		})
	}
}

// BenchmarkValidateReferences measures reference validation at scale.
func BenchmarkValidateReferences(b *testing.B) {
	for _, n := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			g := buildLargeGraph(n)
			analyzer := NewAnalyzer(g)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				analyzer.ValidateReferences()
			}
		})
	}
}
