package main

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// repoRoot locates the workspace root relative to this source file.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	// src/cmd/req42-tracer/cli_test.go → Dir → src/cmd/req42-tracer → ../../.. → repo root
	abs, err := filepath.Abs(filepath.Join(filepath.Dir(thisFile), "..", "..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// execCmd runs a command from the repo root, returns cobra error only.
// Commands write directly to os.Stdout/Stderr so buffer only captures cobra usage errors.
func execCmd(t *testing.T, args ...string) error {
	t.Helper()
	root := repoRoot(t)
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir %s: %v", root, err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs(args)
	return cmd.Execute()
}

// [test-spec,id=TS-CMD-020,req=REQ-REPORT-001,aspice=SWE.5-BP3]
// TestRunRootHelp verifies the root command exits without error.
func TestRunRootHelp(t *testing.T) {
	if err := execCmd(t, "--help"); err != nil {
		t.Errorf("root --help returned error: %v", err)
	}
}

// [test-spec,id=TS-CMD-021,req=REQ-REPORT-001,aspice=SWE.5-BP3]
// TestRunTraceCmd verifies the trace command exits without error on the project.
func TestRunTraceCmd(t *testing.T) {
	if err := execCmd(t, "trace", "--config", "project/req42-tracer/.req42.yaml"); err != nil {
		t.Errorf("trace command error: %v", err)
	}
}

// [test-spec,id=TS-CMD-022,req=REQ-REPORT-002,aspice=SWE.5-BP3]
// TestRunGapsCmd verifies the gaps command exits without error on the project.
func TestRunGapsCmd(t *testing.T) {
	if err := execCmd(t, "gaps", "--config", "project/req42-tracer/.req42.yaml"); err != nil {
		t.Errorf("gaps command error: %v", err)
	}
}

// [test-spec,id=TS-CMD-023,req=REQ-ASPICE-001,aspice=SWE.5-BP3]
// TestRunAspiceCmd verifies the aspice command exits without error on the project.
func TestRunAspiceCmd(t *testing.T) {
	if err := execCmd(t, "aspice", "--config", "project/req42-tracer/.req42.yaml"); err != nil {
		t.Errorf("aspice command error: %v", err)
	}
}

// [test-spec,id=TS-CMD-024,req=REQ-VALIDATE-001,aspice=SWE.5-BP3]
// TestRunValidateCmd verifies the validate command runs without panic.
// The project may have validation warnings — non-zero exit is acceptable.
func TestRunValidateCmd(t *testing.T) {
	// validate may return an error for warnings/violations — that is expected behaviour
	_ = execCmd(t, "validate", "--config", "project/req42-tracer/.req42.yaml")
}

// [test-spec,id=TS-CMD-025,req=REQ-INIT-001,aspice=SWE.5-BP3]
// TestRunInitCmd verifies init creates expected files in a temp directory.
func TestRunInitCmd(t *testing.T) {
	dir := t.TempDir()
	if err := execCmd(t, "init",
		"--dir", dir,
		"--name", "TestProject",
		"--module", "github.com/test/testproject",
		"--interactive=false"); err != nil {
		t.Fatalf("init command error: %v", err)
	}
	for _, rel := range []string{
		".req42.yaml",
		"docs/requirements/req42.adoc",
		"docs/arc42/arc42.adoc",
	} {
		if _, statErr := os.Stat(filepath.Join(dir, rel)); statErr != nil {
			t.Errorf("expected %s to be created: %v", rel, statErr)
		}
	}
}

// [test-spec,id=TS-CMD-026,req=REQ-REPORT-001,aspice=SWE.5-BP3]
// TestRunTraceCmd_HTMLOutput verifies trace --output writes a non-empty HTML report.
func TestRunTraceCmd_HTMLOutput(t *testing.T) {
	root := repoRoot(t)
	out := filepath.Join(t.TempDir(), "report.html")

	orig, _ := os.Getwd()
	t.Cleanup(func() { _ = os.Chdir(orig) })
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cmd := NewRootCmd()
	cmd.SetOut(new(bytes.Buffer))
	cmd.SetErr(new(bytes.Buffer))
	cmd.SetArgs([]string{"trace", "--config", "project/req42-tracer/.req42.yaml", "--output", out})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("trace --output error: %v", err)
	}

	info, err := os.Stat(out)
	if err != nil {
		t.Fatalf("HTML report not created at %s: %v", out, err)
	}
	if info.Size() == 0 {
		t.Error("HTML report is empty")
	}
}

// [test-spec,id=TS-CMD-027,req=REQ-REPORT-001,aspice=SWE.5-BP3]
// TestRunTraceCmd_MarkdownFormat verifies trace --format markdown runs without error.
func TestRunTraceCmd_MarkdownFormat(t *testing.T) {
	if err := execCmd(t, "trace",
		"--config", "project/req42-tracer/.req42.yaml",
		"--format", "markdown"); err != nil {
		t.Errorf("trace --format markdown error: %v", err)
	}
}

// [test-spec,id=TS-CMD-028,req=REQ-REPORT-002,aspice=SWE.5-BP3]
// TestRunGapsCmd_MarkdownFormat verifies gaps --format markdown runs without error.
func TestRunGapsCmd_MarkdownFormat(t *testing.T) {
	if err := execCmd(t, "gaps",
		"--config", "project/req42-tracer/.req42.yaml",
		"--format", "markdown"); err != nil {
		t.Errorf("gaps --format markdown error: %v", err)
	}
}

// [test-spec,id=TS-CMD-029,req=REQ-VALIDATE-001,aspice=SWE.5-BP3]
// TestRunValidateCmd_Verbose verifies validate --verbose runs without panic.
func TestRunValidateCmd_Verbose(t *testing.T) {
	_ = execCmd(t, "validate",
		"--config", "project/req42-tracer/.req42.yaml",
		"--verbose")
}

// [test-spec,id=TS-CMD-030,req=REQ-REPORT-001,aspice=SWE.5-BP3]
// TestUnknownCommand verifies an unknown subcommand returns a usage error.
func TestUnknownCommand(t *testing.T) {
	err := execCmd(t, "nonexistent-subcommand")
	if err == nil {
		t.Error("expected error for unknown command, got nil")
	}
	if !strings.Contains(err.Error(), "unknown command") {
		t.Errorf("unexpected error: %v", err)
	}
}
