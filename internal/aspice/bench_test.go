package aspice

import (
	"fmt"
	"testing"

	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
)

func buildLargeASPICEGraph(n int) *graph.Analyzer {
	g := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement, n),
		ArchElements: make(map[string]*model.ArchElement, n+1),
		TestSpecs:    make(map[string]*model.TestSpec, n/2),
		TestCodes:    make(map[string]*model.TestCode),
		TestResults:  make(map[string]*model.TestResult),
		Links:        make([]*model.TraceLink, 0, n*2),
	}
	g.ArchElements["comp.system"] = &model.ArchElement{ID: "comp.system"}
	for i := 0; i < n; i++ {
		reqID := fmt.Sprintf("REQ-%04d", i)
		archID := fmt.Sprintf("comp.m%04d", i)
		specID := fmt.Sprintf("TS-%04d", i)
		g.Requirements[reqID] = &model.Requirement{ID: reqID, ASPICE: "SWE.1"}
		g.ArchElements[archID] = &model.ArchElement{
			ID: archID, Parent: "comp.system", Req: []string{reqID}, Impl: "impl.go",
		}
		if i%2 == 0 {
			g.TestSpecs[specID] = &model.TestSpec{ID: specID, Req: []string{reqID}}
		}
		g.Links = append(g.Links,
			&model.TraceLink{FromID: reqID, ToID: archID, LinkType: "satisfied-by", Status: "active"},
		)
		if i%2 == 0 {
			// arch→spec verified-by: FromType="arch" so testedArchIDs gets populated
			g.Links = append(g.Links,
				&model.TraceLink{FromID: archID, ToID: specID, FromType: "arch", ToType: "test-spec", LinkType: "verified-by", Status: "active"},
			)
		}
	}
	return graph.NewAnalyzer(g)
}

// BenchmarkCheckCompliance measures ASPICE compliance check at scale.
func BenchmarkCheckCompliance(b *testing.B) {
	for _, n := range []int{100, 500, 1000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			analyzer := buildLargeASPICEGraph(n)
			config := &model.Config{
				ASPICE: struct {
					AutoDerive   bool                         `yaml:"auto-derive"`
					Processes    []string                     `yaml:"processes"`
					ProcessRules map[string]map[string]string `yaml:"process-rules"`
				}{
					Processes: []string{"SWE.1", "SWE.2", "SWE.3", "SWE.5"},
				},
			}
			checker := NewChecker(analyzer, config)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				checker.CheckCompliance()
			}
		})
	}
}
