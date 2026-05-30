package main

// [test-spec,id=TS-CMD-001,req=REQ-PARSE-002,aspice=SWE.5-BP3]
// Test: loadBausteinsicht warnt bei fehlendem JSONC-File (non-verbose)
// [end]

// [test-spec,id=TS-CMD-002,req=REQ-PARSE-002,aspice=SWE.5-BP3]
// Test: loadBausteinsicht mergt JSONC-Arch-Elemente in den Builder
// [end]

// [test-spec,id=TS-CMD-003,req=REQ-PARSE-002,aspice=SWE.5-BP3]
// Test: loadBausteinsicht überspringt doppelte IDs statt Fatal-Error
// [end]

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/graph"
	"github.com/paulefl/req42-tracer/src/internal/model"
)

const sampleBausteinsichtCmd = `{
	"model": {
		"proj": {
			"description": "test project",
			"elements": {
				"comp": {
					"description": "a component",
					"type": "component"
				}
			}
		}
	}
}`

func writeTempJSONC(t *testing.T, content string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), "architecture.jsonc")
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return f
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stderr
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = old
	buf := make([]byte, 4096)
	n, _ := r.Read(buf)
	return string(buf[:n])
}

func TestLoadBausteinsicht_WarnOnMissingFile(t *testing.T) {
	builder := graph.NewBuilder()
	stderr := captureStderr(t, func() {
		loadBausteinsicht(builder, "/nonexistent/architecture.jsonc", "software", false)
	})
	if !strings.Contains(stderr, "Warning") {
		t.Errorf("expected warning on stderr for missing file, got: %q", stderr)
	}
	if len(builder.GetGraph().ArchElements) != 0 {
		t.Error("expected empty graph after failed load")
	}
}

func TestLoadBausteinsicht_MergesElements(t *testing.T) {
	f := writeTempJSONC(t, sampleBausteinsichtCmd)
	builder := graph.NewBuilder()
	loadBausteinsicht(builder, f, "software", false)

	g := builder.GetGraph()
	if len(g.ArchElements) == 0 {
		t.Error("expected arch elements after loading JSONC")
	}
	if _, ok := g.ArchElements["comp"]; !ok {
		t.Error("expected 'comp' element in graph")
	}
}

func TestLoadBausteinsicht_SkipsDuplicateIDs(t *testing.T) {
	f := writeTempJSONC(t, sampleBausteinsichtCmd)
	builder := graph.NewBuilder()

	// Pre-populate 'comp' to create a conflict
	existing := &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: map[string]*model.ArchElement{
			"comp": {ID: "comp", Title: "existing", Project: "software"},
		},
		TestSpecs:   make(map[string]*model.TestSpec),
		TestCodes:   make(map[string]*model.TestCode),
		TestResults: make(map[string]*model.TestResult),
	}
	if err := builder.MergeGraph(existing); err != nil {
		t.Fatalf("setup MergeGraph: %v", err)
	}

	stderr := captureStderr(t, func() {
		loadBausteinsicht(builder, f, "software", false)
	})

	// Should warn about skipped duplicate, not fatal
	if !strings.Contains(stderr, "duplicate") && !strings.Contains(stderr, "skipped") {
		t.Errorf("expected duplicate-skip warning, got: %q", stderr)
	}
	// Original element must still be present (AsciiDoc takes precedence)
	elem := builder.GetGraph().ArchElements["comp"]
	if elem == nil || elem.Title != "existing" {
		t.Error("expected original 'comp' element to be preserved")
	}
}
