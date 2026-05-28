package lsp

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/paulefl/req42-tracer/internal/model"
)

// hoverPattern finds req=/arch=/test-spec= attribute values anywhere in a line.
// Group 1 = attribute name, group 2 = start index (unused), group 3 = value.
var hoverPattern = regexp.MustCompile(`(?:^|[,\s])(req|arch|test-spec)=([^,\]"'\s]+)`)

// detectHoverValue returns the attribute name and ID value under the given
// column in the line, or ("", "", false) if the cursor is not over a value.
func detectHoverValue(line string, col int) (attr, value string, ok bool) {
	for _, m := range hoverPattern.FindAllStringSubmatchIndex(line, -1) {
		// m[4],m[5] = start/end of attr name, m[6],m[7] = start/end of value
		attrStart, attrEnd := m[2], m[3]
		valStart, valEnd := m[4], m[5]
		_ = attrStart
		_ = attrEnd
		if col >= valStart && col <= valEnd {
			return line[m[2]:m[3]], line[valStart:valEnd], true
		}
	}
	return "", "", false
}

// HoverResult is the LSP Hover response.
type HoverResult struct {
	Contents markupContent `json:"contents"`
}

type markupContent struct {
	Kind  string `json:"kind"`  // "markdown"
	Value string `json:"value"`
}

// buildHoverContent looks up the ID in the graph and returns a HoverResult,
// or nil if the ID is not found.
func buildHoverContent(attr, id string, g *model.TraceabilityGraph) *HoverResult {
	var md string
	switch attr {
	case "req":
		req, ok := g.Requirements[id]
		if !ok {
			return nil
		}
		md = fmt.Sprintf("**%s** *(requirement)*\n\n%s", req.ID, formatReq(req))
	case "arch":
		arch, ok := g.ArchElements[id]
		if !ok {
			return nil
		}
		md = fmt.Sprintf("**%s** *(arch element)*\n\n%s", arch.ID, formatArch(arch))
	case "test-spec":
		spec, ok := g.TestSpecs[id]
		if !ok {
			return nil
		}
		md = fmt.Sprintf("**%s** *(test spec)*\n\n%s", spec.ID, formatSpec(spec))
	default:
		return nil
	}
	return &HoverResult{Contents: markupContent{Kind: "markdown", Value: md}}
}

func formatReq(r *model.Requirement) string {
	var b strings.Builder
	if r.Title != "" {
		b.WriteString(r.Title + "\n\n")
	}
	if r.Text != "" {
		b.WriteString(r.Text + "\n\n")
	}
	if r.Priority != "" {
		fmt.Fprintf(&b, "- Priority: %s\n", r.Priority)
	}
	if r.Status != "" {
		fmt.Fprintf(&b, "- Status: %s\n", r.Status)
	}
	if r.ASPICE != "" {
		fmt.Fprintf(&b, "- ASPICE: %s\n", r.ASPICE)
	}
	return b.String()
}

func formatArch(a *model.ArchElement) string {
	var b strings.Builder
	if a.Title != "" {
		b.WriteString(a.Title + "\n\n")
	}
	if a.Text != "" {
		b.WriteString(a.Text + "\n\n")
	}
	if a.Impl != "" {
		fmt.Fprintf(&b, "- impl: `%s`\n", a.Impl)
	}
	if a.Parent != "" {
		fmt.Fprintf(&b, "- parent: %s\n", a.Parent)
	}
	return b.String()
}

func formatSpec(s *model.TestSpec) string {
	var b strings.Builder
	if s.Title != "" {
		b.WriteString(s.Title + "\n\n")
	}
	if s.Text != "" {
		b.WriteString(s.Text + "\n\n")
	}
	return b.String()
}
