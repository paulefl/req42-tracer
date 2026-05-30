package model

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleConfig = `
projects:
  software:
    path: .
    docs: docs/requirements

aspice:
  auto-derive: true
  processes:
    - SWE.1
    - SWE.2

reports:
  html:
    output: reports/report.html
  cli:
    format: markdown
`

// [test-spec,id=TS-MODEL-001,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestLoadConfig verifies that a valid YAML config file is loaded correctly.
func TestLoadConfig(t *testing.T) {
	f := writeTempConfig(t, sampleConfig)
	cfg, err := LoadConfig(f)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg == nil {
		t.Fatal("LoadConfig returned nil")
	}
	if _, ok := cfg.Projects["software"]; !ok {
		t.Error("expected 'software' project in config")
	}
	if cfg.Reports.HTML.Output != "reports/report.html" {
		t.Errorf("HTML output = %q", cfg.Reports.HTML.Output)
	}
	if cfg.Reports.CLI.Format != "markdown" {
		t.Errorf("CLI format = %q", cfg.Reports.CLI.Format)
	}
}

// [test-spec,id=TS-MODEL-002,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestLoadConfig_FileNotFound verifies error for missing config file.
func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/.req42.yaml")
	if err == nil {
		t.Error("expected error for missing config file")
	}
}

// [test-spec,id=TS-MODEL-003,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestLoadConfig_Defaults verifies that default values are applied for missing fields.
func TestLoadConfig_Defaults(t *testing.T) {
	f := writeTempConfig(t, "projects: {}")
	cfg, err := LoadConfig(f)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.Reports.HTML.Output == "" {
		t.Error("expected default HTML output path")
	}
	if cfg.Reports.CLI.Format == "" {
		t.Error("expected default CLI format")
	}
	if len(cfg.ASPICE.Processes) == 0 {
		t.Error("expected default ASPICE processes")
	}
}

// [test-spec,id=TS-MODEL-004,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestLoadConfig_InvalidYAML verifies error for malformed YAML.
func TestLoadConfig_InvalidYAML(t *testing.T) {
	f := writeTempConfig(t, "projects: [invalid: yaml: :")
	_, err := LoadConfig(f)
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

// [test-spec,id=TS-MODEL-005,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestConfig_GetRule verifies GetRule returns configured and default rule values.
func TestConfig_GetRule(t *testing.T) {
	cfg := &Config{
		Rules: map[string]string{
			"orphan-req": "error",
		},
	}
	if cfg.GetRule("orphan-req") != "error" {
		t.Errorf("GetRule(orphan-req) = %q, want error", cfg.GetRule("orphan-req"))
	}
	if cfg.GetRule("unknown-rule") != "warning" {
		t.Errorf("GetRule(unknown) = %q, want warning", cfg.GetRule("unknown-rule"))
	}
}

// [test-spec,id=TS-MODEL-006,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestConfig_SetDefault verifies that SetDefault only sets when rule is not already present.
func TestConfig_SetDefault(t *testing.T) {
	cfg := &Config{
		Rules: map[string]string{
			"existing-rule": "error",
		},
	}
	cfg.SetDefault("existing-rule", "warning")
	cfg.SetDefault("new-rule", "off")

	if cfg.Rules["existing-rule"] != "error" {
		t.Errorf("SetDefault should not overwrite existing rule, got %q", cfg.Rules["existing-rule"])
	}
	if cfg.Rules["new-rule"] != "off" {
		t.Errorf("SetDefault should set new rule, got %q", cfg.Rules["new-rule"])
	}
}

// [test-spec,id=TS-MODEL-007,req="REQ-CONFIG-001",aspice="SWE.5.BP3"]
// TestLoadConfig_DefaultMaps verifies that nil maps are initialized on load.
func TestLoadConfig_DefaultMaps(t *testing.T) {
	f := writeTempConfig(t, "{}")
	cfg, err := LoadConfig(f)
	if err != nil {
		t.Fatalf("LoadConfig error: %v", err)
	}
	if cfg.Projects == nil {
		t.Error("Projects should be initialized, not nil")
	}
	if cfg.Rules == nil {
		t.Error("Rules should be initialized, not nil")
	}
	if cfg.ASPICE.ProcessRules == nil {
		t.Error("ASPICE.ProcessRules should be initialized, not nil")
	}
}

// [test-spec,id=TS-MODEL-008,req=REQ-CONFIG-001,aspice=SWE.5-BP3]
// TestConfig_GetDefaultProject verifies priority: explicit field > projects map key > fallback.
func TestConfig_GetDefaultProject(t *testing.T) {
	// Explicit default-project field takes priority
	cfg := &Config{DefaultProject: "firmware", Projects: map[string]*ProjectConfig{"software": {}}}
	if got := cfg.GetDefaultProject(); got != "firmware" {
		t.Errorf("explicit default-project: got %q, want firmware", got)
	}

	// Falls back to first projects key when default-project is empty
	cfg2 := &Config{Projects: map[string]*ProjectConfig{"hardware": {}}}
	if got := cfg2.GetDefaultProject(); got != "hardware" {
		t.Errorf("projects key fallback: got %q, want hardware", got)
	}

	// Falls back to "software" when projects map is empty
	cfg3 := &Config{Projects: map[string]*ProjectConfig{}}
	if got := cfg3.GetDefaultProject(); got != "software" {
		t.Errorf("empty projects fallback: got %q, want software", got)
	}
}

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	f := filepath.Join(dir, ".req42.yaml")
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return f
}
