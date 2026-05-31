package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// [test-spec,id=TS-LSP-015,req="REQ-LSP-001",aspice="SWE.4.BP2"]
// TestSplitLines verifies splitLines handles LF, CRLF and empty input correctly.
func TestSplitLines(t *testing.T) {
	cases := []struct {
		input string
		lines int
		first string
	}{
		{"line1\nline2", 2, "line1"},
		{"line1\r\nline2\r\n", 3, "line1"},
		{"", 1, ""},
		{"single", 1, "single"},
	}
	for _, c := range cases {
		got := splitLines(c.input)
		if len(got) != c.lines {
			t.Errorf("splitLines(%q): got %d lines, want %d", c.input, len(got), c.lines)
		}
		if len(got) > 0 && got[0] != c.first {
			t.Errorf("splitLines(%q)[0] = %q, want %q", c.input, got[0], c.first)
		}
	}
}

// [test-spec,id=TS-LSP-016,req="REQ-LSP-001",aspice="SWE.4.BP2"]
// TestLineAt verifies lineAt returns line content or empty string for edge cases.
func TestLineAt(t *testing.T) {
	srv := &Server{docs: map[string][]string{
		"file:///test.adoc": {"line zero", "line one", "line two"},
	}}

	if got := srv.lineAt("file:///test.adoc", 0); got != "line zero" {
		t.Errorf("lineAt(0) = %q, want 'line zero'", got)
	}
	if got := srv.lineAt("file:///test.adoc", 2); got != "line two" {
		t.Errorf("lineAt(2) = %q, want 'line two'", got)
	}
	if got := srv.lineAt("file:///test.adoc", 99); got != "" {
		t.Errorf("lineAt(99) = %q, want empty", got)
	}
	if got := srv.lineAt("file:///unknown.adoc", 0); got != "" {
		t.Errorf("lineAt unknown doc = %q, want empty", got)
	}
}

// [test-spec,id=TS-LSP-017,req="REQ-LSP-001",aspice="SWE.4.BP2"]
// TestNewServer verifies NewServer returns an initialized, non-nil Server.
func TestNewServer(t *testing.T) {
	srv := NewServer(&model.Config{})
	if srv == nil {
		t.Fatal("NewServer returned nil")
	}
	if srv.docs == nil {
		t.Error("docs map not initialized")
	}
}

func makeMinimalGraph() *model.TraceabilityGraph {
	return &model.TraceabilityGraph{
		Requirements:   make(map[string]*model.Requirement),
		ArchElements:   make(map[string]*model.ArchElement),
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestCodes:      make(map[string]*model.TestCode),
		TestResults:    make(map[string]*model.TestResult),
		Links:          []*model.TraceLink{},
	}
}

func makeTestServer(input string) (*Server, *bytes.Buffer) {
	var out bytes.Buffer
	srv := &Server{
		in:    bufio.NewReader(strings.NewReader(input)),
		out:   &out,
		log:   log.New(io.Discard, "", 0),
		docs:  make(map[string][]string),
		graph: makeMinimalGraph(),
	}
	return srv, &out
}

// [test-spec,id=TS-LSP-018,req="REQ-LSP-001",aspice="SWE.5.BP3"]
// TestHandleDidOpen verifies textDocument/didOpen stores document lines.
func TestHandleDidOpen(t *testing.T) {
	params := json.RawMessage(`{"textDocument":{"uri":"file:///test.adoc","languageId":"asciidoc","version":1,"text":"line1\nline2"}}`)
	didOpen := buildMessage(t, message{JSONRPC: "2.0", Method: "textDocument/didOpen", Params: params})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	srv, _ := makeTestServer(didOpen + exit)
	_ = srv.Run()

	lines := srv.docs["file:///test.adoc"]
	if len(lines) != 2 {
		t.Errorf("expected 2 lines after didOpen, got %d", len(lines))
	}
}

// [test-spec,id=TS-LSP-019,req="REQ-LSP-001",aspice="SWE.5.BP3"]
// TestHandleDidChange verifies textDocument/didChange updates document content.
func TestHandleDidChange(t *testing.T) {
	params := json.RawMessage(`{"textDocument":{"uri":"file:///test.adoc","version":2},"contentChanges":[{"text":"new\ncontent\nhere"}]}`)
	didChange := buildMessage(t, message{JSONRPC: "2.0", Method: "textDocument/didChange", Params: params})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	srv, _ := makeTestServer(didChange + exit)
	srv.docs["file:///test.adoc"] = []string{"old content"}
	_ = srv.Run()

	lines := srv.docs["file:///test.adoc"]
	if len(lines) != 3 || lines[0] != "new" {
		t.Errorf("expected 3 updated lines, got %v", lines)
	}
}

// [test-spec,id=TS-LSP-020,req="REQ-LSP-001",aspice="SWE.5.BP3"]
// TestHandleCompletion verifies textDocument/completion returns a valid result.
func TestHandleCompletion(t *testing.T) {
	params := json.RawMessage(`{"textDocument":{"uri":"file:///test.adoc"},"position":{"line":0,"character":5}}`)
	req := buildMessage(t, message{JSONRPC: "2.0", ID: 10, Method: "textDocument/completion", Params: params})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	srv, out := makeTestServer(req + exit)
	srv.docs["file:///test.adoc"] = []string{"[req,"}
	_ = srv.Run()

	resp := parseFirstResponse(t, out)
	if resp.Error != nil {
		t.Fatalf("completion error: %v", resp.Error)
	}
}

// [test-spec,id=TS-LSP-021,req="REQ-LSP-001",aspice="SWE.5.BP3"]
// TestHandleHover verifies textDocument/hover returns result without error.
func TestHandleHover(t *testing.T) {
	params := json.RawMessage(`{"textDocument":{"uri":"file:///test.adoc"},"position":{"line":0,"character":5}}`)
	req := buildMessage(t, message{JSONRPC: "2.0", ID: 11, Method: "textDocument/hover", Params: params})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	srv, out := makeTestServer(req + exit)
	srv.docs["file:///test.adoc"] = []string{"[req,id=REQ-001]"}
	_ = srv.Run()

	resp := parseFirstResponse(t, out)
	if resp.Error != nil {
		t.Fatalf("hover error: %v", resp.Error)
	}
}

// [test-spec,id=TS-LSP-022,req="REQ-LSP-001",aspice="SWE.5.BP3"]
// TestHandleDefinition verifies textDocument/definition returns result without error.
func TestHandleDefinition(t *testing.T) {
	params := json.RawMessage(`{"textDocument":{"uri":"file:///test.adoc"},"position":{"line":0,"character":5}}`)
	req := buildMessage(t, message{JSONRPC: "2.0", ID: 12, Method: "textDocument/definition", Params: params})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	srv, out := makeTestServer(req + exit)
	srv.docs["file:///test.adoc"] = []string{"[req,id=REQ-001]"}
	_ = srv.Run()

	resp := parseFirstResponse(t, out)
	if resp.Error != nil {
		t.Fatalf("definition error: %v", resp.Error)
	}
}
