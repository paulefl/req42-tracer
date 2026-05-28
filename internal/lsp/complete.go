package lsp

import (
	"regexp"
	"sort"
	"strings"

	"github.com/paulefl/req42-tracer/internal/model"
)

// completionContext describes what kind of ID the cursor is completing.
type completionContext int

const (
	ctxNone completionContext = iota
	ctxReq
	ctxArch
	ctxTestSpec
)

// attrPattern matches LSP-relevant attribute contexts immediately before the
// cursor. The cursor sits right after the matched prefix.
//
// Examples matched (| = cursor):
//   req=|            → ctxReq
//   req=REQ-PAR|     → ctxReq  (partial value)
//   arch=com|        → ctxArch
//   test-spec=spec|  → ctxTestSpec
var attrPattern = regexp.MustCompile(
	`(?:^|[,\s])(req|arch|test-spec)=([^,\]"'\s]*)$`,
)

// detectContext returns the completion kind and the already-typed prefix at
// the cursor position (lineUpToCursor is the content of the current line from
// column 0 to the cursor column, inclusive).
func detectContext(lineUpToCursor string) (completionContext, string) {
	m := attrPattern.FindStringSubmatch(lineUpToCursor)
	if m == nil {
		return ctxNone, ""
	}
	prefix := m[2]
	switch m[1] {
	case "req":
		return ctxReq, prefix
	case "arch":
		return ctxArch, prefix
	case "test-spec":
		return ctxTestSpec, prefix
	}
	return ctxNone, ""
}

// CompletionItem is a single LSP completion suggestion.
type CompletionItem struct {
	Label         string `json:"label"`
	Kind          int    `json:"kind"` // 6 = Variable
	Detail        string `json:"detail,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

// CompletionList is the LSP CompletionList response.
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}

// buildCompletions returns a CompletionList for the given line-up-to-cursor
// string using IDs from the provided graph.
func buildCompletions(lineUpToCursor string, g *model.TraceabilityGraph) CompletionList {
	ctx, prefix := detectContext(lineUpToCursor)
	if ctx == ctxNone {
		return CompletionList{}
	}

	var items []CompletionItem
	switch ctx {
	case ctxReq:
		for id, req := range g.Requirements {
			if strings.HasPrefix(id, prefix) {
				items = append(items, CompletionItem{
					Label:         id,
					Kind:          6,
					Detail:        req.Title,
					Documentation: req.Text,
				})
			}
		}
	case ctxArch:
		for id, arch := range g.ArchElements {
			if strings.HasPrefix(id, prefix) {
				items = append(items, CompletionItem{
					Label:  id,
					Kind:   6,
					Detail: arch.Title,
				})
			}
		}
	case ctxTestSpec:
		for id, spec := range g.TestSpecs {
			if strings.HasPrefix(id, prefix) {
				items = append(items, CompletionItem{
					Label:  id,
					Kind:   6,
					Detail: spec.Title,
				})
			}
		}
	}

	sort.Slice(items, func(i, j int) bool { return items[i].Label < items[j].Label })
	return CompletionList{IsIncomplete: false, Items: items}
}
