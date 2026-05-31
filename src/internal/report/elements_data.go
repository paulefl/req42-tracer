package report

import "github.com/paulefl/req42-tracer/src/internal/model"

// ElementItem represents one traceability element for the Elements-Tab.
type ElementItem struct {
	ID        string   `json:"id"`
	Type      string   `json:"type"`   // req, arch, dsn, test-spec, test-result
	Title     string   `json:"title"`
	ASPICE    string   `json:"aspice,omitempty"`
	Impl      string   `json:"impl,omitempty"`
	Status    string   `json:"status,omitempty"`
	TraceUp   []string `json:"trace_up"`   // IDs of elements that link TO this element
	TraceDown []string `json:"trace_down"` // IDs of elements this element links TO
}

// ElementsData is the payload injected into the Elements-Tab.
type ElementsData struct {
	Items []ElementItem `json:"items"`
}

// BuildElementsData collects all graph elements with their up/down traces.
func BuildElementsData(g *model.TraceabilityGraph) *ElementsData {
	data := &ElementsData{Items: make([]ElementItem, 0)}

	// Build reverse-link index: toID → []fromID
	upIndex := make(map[string][]string)
	downIndex := make(map[string][]string)
	for _, link := range g.Links {
		downIndex[link.FromID] = append(downIndex[link.FromID], link.ToID)
		upIndex[link.ToID] = append(upIndex[link.ToID], link.FromID)
	}

	for _, r := range g.Requirements {
		data.Items = append(data.Items, ElementItem{
			ID:        r.ID,
			Type:      "req",
			Title:     r.Title,
			ASPICE:    r.ASPICE,
			Status:    r.Status,
			TraceUp:   orEmpty(upIndex[r.ID]),
			TraceDown: orEmpty(downIndex[r.ID]),
		})
	}
	for _, a := range g.ArchElements {
		data.Items = append(data.Items, ElementItem{
			ID:        a.ID,
			Type:      "arch",
			Title:     a.Title,
			ASPICE:    a.ASPICE,
			Impl:      a.Impl,
			TraceUp:   orEmpty(upIndex[a.ID]),
			TraceDown: orEmpty(downIndex[a.ID]),
		})
	}
	for _, d := range g.DesignElements {
		data.Items = append(data.Items, ElementItem{
			ID:        d.ID,
			Type:      "dsn",
			Title:     d.Title,
			ASPICE:    d.ASPICE,
			Impl:      d.Impl,
			TraceUp:   orEmpty(upIndex[d.ID]),
			TraceDown: orEmpty(downIndex[d.ID]),
		})
	}
	for _, s := range g.TestSpecs {
		data.Items = append(data.Items, ElementItem{
			ID:        s.ID,
			Type:      "test-spec",
			Title:     s.Title,
			TraceUp:   orEmpty(upIndex[s.ID]),
			TraceDown: orEmpty(downIndex[s.ID]),
		})
	}
	for _, r := range g.TestResults {
		data.Items = append(data.Items, ElementItem{
			ID:        r.ID,
			Type:      "test-result",
			Title:     r.TestName,
			Status:    r.Status,
			TraceUp:   orEmpty(upIndex[r.ID]),
			TraceDown: orEmpty(downIndex[r.ID]),
		})
	}

	return data
}

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
