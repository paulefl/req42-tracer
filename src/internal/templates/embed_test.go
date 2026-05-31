package templates

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// [test-spec,id=TS-TPL-001,req="REQ-INIT-001",aspice="SWE.4.BP2"]
// TestFS_AllFilesEmbedded verifies all expected template files are embedded and non-empty.
func TestFS_AllFilesEmbedded(t *testing.T) {
	files := []string{
		"req42.adoc",
		"arc42.adoc",
		"architecture.jsonc",
		".req42.yaml",
		".gitignore",
	}
	for _, f := range files {
		data, err := FS.ReadFile(f)
		if err != nil {
			t.Errorf("template %s not embedded: %v", f, err)
			continue
		}
		if len(data) == 0 {
			t.Errorf("template %s is empty", f)
		}
	}
}

// [test-spec,id=TS-TPL-002,req="REQ-INIT-001",aspice="SWE.4.BP2"]
// TestFS_ReqAdocHasReqBlock verifies req42.adoc contains at least one [req,id=] block.
func TestFS_ReqAdocHasReqBlock(t *testing.T) {
	data, err := FS.ReadFile("req42.adoc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "[req,id=") {
		t.Error("req42.adoc template missing [req,id=] block")
	}
}

// [test-spec,id=TS-TPL-003,req="REQ-INIT-001",aspice="SWE.4.BP2"]
// TestFS_Arc42AdocHasArchBlock verifies arc42.adoc contains at least one [arch,id=] block.
func TestFS_Arc42AdocHasArchBlock(t *testing.T) {
	data, err := FS.ReadFile("arc42.adoc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "[arch,id=") {
		t.Error("arc42.adoc template missing [arch,id=] block")
	}
}

// [test-spec,id=TS-TPL-004,req="REQ-INIT-001",aspice="SWE.4.BP2"]
// TestFS_ConfigIsValidYAML verifies .req42.yaml is parseable YAML with required fields.
func TestFS_ConfigIsValidYAML(t *testing.T) {
	data, err := FS.ReadFile(".req42.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]interface{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf(".req42.yaml is not valid YAML: %v", err)
	}
	if _, ok := cfg["projects"]; !ok {
		t.Error(".req42.yaml missing 'projects' key")
	}
}

// [test-spec,id=TS-TPL-005,req="REQ-INIT-001",aspice="SWE.4.BP2"]
// TestFS_ArchitectureJSONCHasModel verifies architecture.jsonc contains a model key.
func TestFS_ArchitectureJSONCHasModel(t *testing.T) {
	data, err := FS.ReadFile("architecture.jsonc")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"model"`) {
		t.Error("architecture.jsonc template missing \"model\" key")
	}
}
