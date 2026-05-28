package lsp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"testing"
)

// [test-spec,id=spec.lsp.initialize,req=REQ-LSP-001,arch=comp.lsp,aspice=SWE.4]
// TestInitializeHandshake verifies the LSP initialize/initialized handshake.
func TestInitializeHandshake(t *testing.T) {
	req := buildMessage(t, message{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params:  json.RawMessage(`{"processId":null,"capabilities":{}}`),
	})

	var out bytes.Buffer
	srv := &Server{
		in:  bufio.NewReader(strings.NewReader(req)),
		out: &out,
		log: log.New(io.Discard, "", 0),
	}
	_ = srv.Run()

	resp := parseFirstResponse(t, &out)
	if resp.Error != nil {
		t.Fatalf("initialize error: %v", resp.Error)
	}
	if resp.Result == nil {
		t.Fatal("expected non-nil result")
	}
	var result initializeResult
	data, _ := json.Marshal(resp.Result)
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal result: %v", err)
	}
	if result.ServerInfo.Name != "req42-tracer" {
		t.Errorf("serverInfo.name = %q, want req42-tracer", result.ServerInfo.Name)
	}
	if result.Capabilities.TextDocumentSync != 1 {
		t.Errorf("textDocumentSync = %d, want 1", result.Capabilities.TextDocumentSync)
	}
}

// [test-spec,id=spec.lsp.shutdown,req=REQ-LSP-001,arch=comp.lsp,aspice=SWE.4]
// TestShutdown verifies the server exits cleanly on shutdown → exit sequence.
func TestShutdown(t *testing.T) {
	shutdown := buildMessage(t, message{JSONRPC: "2.0", ID: 2, Method: "shutdown"})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	var out bytes.Buffer
	srv := &Server{
		in:  bufio.NewReader(strings.NewReader(shutdown + exit)),
		out: &out,
		log: log.New(io.Discard, "", 0),
	}
	if err := srv.Run(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resp := parseFirstResponse(t, &out)
	if resp.Error != nil {
		t.Fatalf("shutdown error: %v", resp.Error)
	}
}

// [test-spec,id=spec.lsp.unknown-method,req=REQ-LSP-001,arch=comp.lsp,aspice=SWE.4]
// TestUnknownMethod verifies that unknown request methods return a -32601 error.
func TestUnknownMethod(t *testing.T) {
	req := buildMessage(t, message{JSONRPC: "2.0", ID: 3, Method: "textDocument/foobar"})
	exit := buildMessage(t, message{JSONRPC: "2.0", Method: "exit"})

	var out bytes.Buffer
	srv := &Server{
		in:  bufio.NewReader(strings.NewReader(req + exit)),
		out: &out,
		log: log.New(io.Discard, "", 0),
	}
	_ = srv.Run()

	resp := parseFirstResponse(t, &out)
	if resp.Error == nil {
		t.Fatal("expected error response for unknown method")
	}
	if resp.Error.Code != -32601 {
		t.Errorf("error code = %d, want -32601", resp.Error.Code)
	}
}

// helpers

func buildMessage(t *testing.T, msg message) string {
	t.Helper()
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(data), data)
}

func parseFirstResponse(t *testing.T, buf *bytes.Buffer) message {
	t.Helper()
	r := bufio.NewReader(buf)
	contentLength := -1
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			t.Fatalf("read header: %v", err)
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		if strings.HasPrefix(line, "Content-Length: ") {
			fmt.Sscanf(strings.TrimPrefix(line, "Content-Length: "), "%d", &contentLength)
		}
	}
	if contentLength < 0 {
		t.Fatal("no Content-Length in response")
	}
	body := make([]byte, contentLength)
	if _, err := io.ReadFull(r, body); err != nil {
		t.Fatalf("read body: %v", err)
	}
	var msg message
	if err := json.Unmarshal(body, &msg); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	return msg
}
