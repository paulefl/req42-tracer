package testresult

import (
	"os"
	"path/filepath"
	"testing"
)

const sampleGoCoverage = `github.com/foo/pkg/parser/adoc.go:ParseAll	85.7%
github.com/foo/pkg/parser/adoc.go:parseBlock	100.0%
github.com/foo/pkg/graph/build.go:Build	90.0%
github.com/foo/pkg/graph/build.go:merge	60.0%
total:	(statements)	84.5%
`

// Raw coverage.out profile: mode line + block lines (numStmts count)
const sampleGoCoverageRaw = `mode: set
github.com/foo/pkg/parser/adoc.go:10.56,12.3 5 1
github.com/foo/pkg/parser/adoc.go:14.10,16.2 3 0
github.com/foo/pkg/graph/build.go:20.5,22.1 4 1
github.com/foo/pkg/graph/build.go:24.1,26.2 2 1
`

const sampleCobertura = `<?xml version="1.0"?>
<coverage line-rate="0.85" branch-rate="0.7" version="1" timestamp="1234567890">
  <packages>
    <package name="com.example.parser" line-rate="0.9" branch-rate="0.8" complexity="0">
      <classes>
        <class name="Parser.java" filename="com/example/Parser.java">
          <lines>
            <line number="1" hits="1"/>
            <line number="2" hits="1"/>
            <line number="3" hits="0"/>
            <line number="4" hits="1"/>
          </lines>
        </class>
      </classes>
    </package>
  </packages>
</coverage>`

const sampleLCOV = `SF:src/parser/adoc.c
FN:10,parse_block
DA:10,1
DA:11,1
DA:12,0
DA:13,1
end_of_record
SF:src/graph/build.c
DA:5,1
DA:6,1
DA:7,1
end_of_record
`

// [test-spec,id=TS-COV-001,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseGoCoverage verifies go tool cover -func output is parsed into packages.
func TestParseGoCoverage(t *testing.T) {
	f := writeTempFile(t, sampleGoCoverage, "coverage.out")
	pkgs, err := ParseGoCoverage(f)
	if err != nil {
		t.Fatalf("ParseGoCoverage error: %v", err)
	}
	if len(pkgs) == 0 {
		t.Fatal("expected at least one package")
	}
	// Should have parser and graph packages
	names := make(map[string]bool)
	for _, p := range pkgs {
		names[p.Package] = true
	}
	if !names["parser"] {
		t.Error("expected 'parser' package")
	}
	if !names["graph"] {
		t.Error("expected 'graph' package")
	}
}

// [test-spec,id=TS-COV-002,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseGoCoverage_FileNotFound verifies error for missing file.
func TestParseGoCoverage_FileNotFound(t *testing.T) {
	_, err := ParseGoCoverage("/nonexistent/coverage.out")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// [test-spec,id=TS-COV-003,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseCobertura verifies Cobertura XML is parsed into package coverage entries.
func TestParseCobertura(t *testing.T) {
	f := writeTempFile(t, sampleCobertura, "coverage.xml")
	pkgs, err := ParseCobertura(f)
	if err != nil {
		t.Fatalf("ParseCobertura error: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Statements != 4 {
		t.Errorf("statements = %d, want 4", pkgs[0].Statements)
	}
	if pkgs[0].Covered != 3 {
		t.Errorf("covered = %d, want 3", pkgs[0].Covered)
	}
	if pkgs[0].Pct < 74 || pkgs[0].Pct > 76 {
		t.Errorf("pct = %.1f, want ~75.0", pkgs[0].Pct)
	}
}

// [test-spec,id=TS-COV-004,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseCobertura_FileNotFound verifies error for missing Cobertura file.
func TestParseCobertura_FileNotFound(t *testing.T) {
	_, err := ParseCobertura("/nonexistent/coverage.xml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// [test-spec,id=TS-COV-005,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseCobertura_InvalidXML verifies error for malformed Cobertura XML.
func TestParseCobertura_InvalidXML(t *testing.T) {
	f := writeTempFile(t, "<not valid xml", "bad.xml")
	_, err := ParseCobertura(f)
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

// [test-spec,id=TS-COV-006,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseLCOV verifies LCOV .info format is parsed into package coverage entries.
func TestParseLCOV(t *testing.T) {
	f := writeTempFile(t, sampleLCOV, "coverage.info")
	pkgs, err := ParseLCOV(f)
	if err != nil {
		t.Fatalf("ParseLCOV error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}
	// parser: 4 lines, 3 covered = 75%
	if pkgs[0].Statements != 4 {
		t.Errorf("parser statements = %d, want 4", pkgs[0].Statements)
	}
	if pkgs[0].Covered != 3 {
		t.Errorf("parser covered = %d, want 3", pkgs[0].Covered)
	}
	// graph: 3 lines, 3 covered = 100%
	if pkgs[1].Covered != 3 {
		t.Errorf("graph covered = %d, want 3", pkgs[1].Covered)
	}
}

// [test-spec,id=TS-COV-007,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseLCOV_FileNotFound verifies error for missing LCOV file.
func TestParseLCOV_FileNotFound(t *testing.T) {
	_, err := ParseLCOV("/nonexistent/coverage.info")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

// [test-spec,id=TS-COV-008,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseGoCoverageRaw verifies raw coverage.out profile format is parsed correctly.
func TestParseGoCoverageRaw(t *testing.T) {
	f := writeTempFile(t, sampleGoCoverageRaw, "coverage.out")
	pkgs, err := ParseGoCoverage(f)
	if err != nil {
		t.Fatalf("ParseGoCoverage (raw) error: %v", err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("expected 2 packages, got %d", len(pkgs))
	}
	names := make(map[string]bool)
	for _, p := range pkgs {
		names[p.Package] = true
	}
	if !names["parser"] {
		t.Error("expected 'parser' package")
	}
	if !names["graph"] {
		t.Error("expected 'graph' package")
	}
}

// [test-spec,id=TS-COV-009,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseGoCoverageRaw_Counts verifies statement and covered counts from raw profile.
func TestParseGoCoverageRaw_Counts(t *testing.T) {
	f := writeTempFile(t, sampleGoCoverageRaw, "coverage.out")
	pkgs, err := ParseGoCoverage(f)
	if err != nil {
		t.Fatalf("ParseGoCoverage (raw) error: %v", err)
	}
	// parser: 5+3=8 stmts, 5 covered (first block hit, second not)
	var parser PackageCoverage
	for _, p := range pkgs {
		if p.Package == "parser" {
			parser = p
		}
	}
	if parser.Statements != 8 {
		t.Errorf("parser statements = %d, want 8", parser.Statements)
	}
	if parser.Covered != 5 {
		t.Errorf("parser covered = %d, want 5", parser.Covered)
	}
	wantPct := 5.0 / 8.0 * 100
	if parser.Pct < wantPct-0.1 || parser.Pct > wantPct+0.1 {
		t.Errorf("parser pct = %.2f, want %.2f", parser.Pct, wantPct)
	}
}

// [test-spec,id=TS-COV-010,req="REQ-TESTING-001",aspice="SWE.4.BP2"]
// TestParseGoCoverageRaw_AllCovered verifies 100% coverage is computed correctly.
func TestParseGoCoverageRaw_AllCovered(t *testing.T) {
	content := "mode: set\ngithub.com/foo/pkg/graph/build.go:1.1,5.1 4 2\n"
	f := writeTempFile(t, content, "coverage.out")
	pkgs, err := ParseGoCoverage(f)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(pkgs) != 1 {
		t.Fatalf("expected 1 package, got %d", len(pkgs))
	}
	if pkgs[0].Statements != 4 {
		t.Errorf("statements = %d, want 4", pkgs[0].Statements)
	}
	if pkgs[0].Covered != 4 {
		t.Errorf("covered = %d, want 4", pkgs[0].Covered)
	}
	if pkgs[0].Pct != 100.0 {
		t.Errorf("pct = %.1f, want 100.0", pkgs[0].Pct)
	}
}

func writeTempFile(t *testing.T, content, name string) string {
	t.Helper()
	f := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return f
}
