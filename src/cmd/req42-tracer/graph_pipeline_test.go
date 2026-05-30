package main

// [test-spec,id=TS-CMD-010,req=REQ-GRAPH-001,aspice=SWE.5-BP3]
// Test: buildGraph returns non-nil graph with empty dirs (no crash, empty result)
// [end]

// [test-spec,id=TS-CMD-011,req=REQ-GRAPH-001,aspice=SWE.5-BP3]
// Test: buildGraph loads Bausteinsicht elements when config has model path
// [end]

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

func TestBuildGraph_EmptyDirs(t *testing.T) {
	cfg := &model.Config{
		Projects: map[string]*model.ProjectConfig{"software": {}},
		Rules:    map[string]string{},
	}
	empty := t.TempDir()

	g, err := buildGraph(cfg, empty, empty, "software", false)
	if err != nil {
		t.Fatalf("buildGraph with empty dirs: %v", err)
	}
	if g == nil {
		t.Fatal("expected non-nil graph")
	}
}

func TestBuildGraph_LoadsBausteinsicht(t *testing.T) {
	dir := t.TempDir()
	jsonc := filepath.Join(dir, "architecture.jsonc")
	content := `{"model":{"proj":{"description":"test","elements":{"comp":{"description":"c"}}}}}`
	if err := os.WriteFile(jsonc, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := &model.Config{
		Projects: map[string]*model.ProjectConfig{"software": {}},
		Rules:    map[string]string{},
	}
	cfg.Bausteinsicht.Model = jsonc

	g, err := buildGraph(cfg, t.TempDir(), t.TempDir(), "software", false)
	if err != nil {
		t.Fatalf("buildGraph: %v", err)
	}
	if _, ok := g.ArchElements["comp"]; !ok {
		t.Error("expected Bausteinsicht element 'comp' in graph")
	}
}
