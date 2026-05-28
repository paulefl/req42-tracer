package lsp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/paulefl/req42-tracer/internal/model"
	"github.com/paulefl/req42-tracer/internal/parser"
)

// Server is a minimal LSP server over stdio using JSON-RPC 2.0.
type Server struct {
	in    *bufio.Reader
	out   io.Writer
	log   *log.Logger
	docs  map[string][]string    // uri → lines
	graph *model.TraceabilityGraph
}

// NewServer creates a server that reads from stdin and writes to stdout.
// Diagnostic messages go to stderr.
func NewServer() *Server {
	return &Server{
		in:   bufio.NewReader(os.Stdin),
		out:  os.Stdout,
		log:  log.New(os.Stderr, "[lsp] ", log.LstdFlags),
		docs: make(map[string][]string),
	}
}

// Run processes incoming JSON-RPC messages until shutdown/exit or EOF.
func (s *Server) Run() error {
	s.log.Println("LSP server started on stdio")
	s.reloadGraph()
	shutdownReceived := false
	for {
		msg, err := s.readMessage()
		if err == io.EOF {
			s.log.Println("stdin closed, shutting down")
			return nil
		}
		if err != nil {
			// Non-fatal parse errors: log and skip to next message.
			s.log.Printf("read error (skipping): %v", err)
			continue
		}
		done, writeErr := s.dispatch(msg, &shutdownReceived)
		if writeErr != nil {
			return fmt.Errorf("write: %w", writeErr)
		}
		if done {
			return nil
		}
	}
}

// reloadGraph rebuilds the traceability graph from the project's doc dirs.
func (s *Server) reloadGraph() {
	builder := newGraphBuilder()
	for _, dir := range []string{"docs/requirements", "docs/arc42"} {
		if g, err := parser.ParseAllFromDir(dir, "software"); err == nil {
			builder.merge(g)
		}
	}
	s.graph = builder.graph
	s.log.Printf("graph loaded: %d reqs, %d arch, %d specs",
		len(s.graph.Requirements), len(s.graph.ArchElements), len(s.graph.TestSpecs))
}

// --- JSON-RPC types ---

type message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// --- Wire protocol ---

func (s *Server) readMessage() (*message, error) {
	contentLength := -1
	for {
		line, err := s.in.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		// Case-insensitive header match per RFC 7230.
		lower := strings.ToLower(line)
		if strings.HasPrefix(lower, "content-length:") {
			val := strings.TrimSpace(line[len("content-length:"):])
			n, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("bad Content-Length: %w", err)
			}
			contentLength = n
		}
	}
	if contentLength < 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}
	buf := make([]byte, contentLength)
	if _, err := io.ReadFull(s.in, buf); err != nil {
		return nil, err
	}
	var msg message
	if err := json.Unmarshal(buf, &msg); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}
	return &msg, nil
}

func (s *Server) send(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		s.log.Printf("marshal error: %v", err)
		return err
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := io.WriteString(s.out, header); err != nil {
		return err
	}
	if _, err := s.out.Write(data); err != nil {
		return err
	}
	return nil
}

func (s *Server) reply(id interface{}, result interface{}) error {
	return s.send(message{JSONRPC: "2.0", ID: id, Result: result})
}

func (s *Server) replyError(id interface{}, code int, msg string) error {
	return s.send(message{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}})
}

// --- Dispatch ---

// dispatch handles one message. It returns (exit, writeError).
func (s *Server) dispatch(msg *message, shutdownReceived *bool) (exit bool, writeErr error) {
	s.log.Printf("← %s (id=%v)", msg.Method, msg.ID)
	switch msg.Method {
	case "initialize":
		writeErr = s.handleInitialize(msg)
	case "initialized":
		// notification, no response required
	case "shutdown":
		*shutdownReceived = true
		writeErr = s.reply(msg.ID, nil)
	case "exit":
		return true, nil
	case "textDocument/didOpen":
		s.handleDidOpen(msg)
	case "textDocument/didChange":
		s.handleDidChange(msg)
	case "textDocument/completion":
		writeErr = s.handleCompletion(msg)
	default:
		if msg.ID != nil {
			writeErr = s.replyError(msg.ID, -32601, "method not found: "+msg.Method)
		}
	}
	return false, writeErr
}

// --- initialize ---

type initializeResult struct {
	Capabilities serverCapabilities `json:"capabilities"`
	ServerInfo   serverInfo         `json:"serverInfo"`
}

type serverCapabilities struct {
	TextDocumentSync   int                  `json:"textDocumentSync"`
	CompletionProvider *completionOptions   `json:"completionProvider,omitempty"`
}

type completionOptions struct {
	TriggerCharacters []string `json:"triggerCharacters"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (s *Server) handleInitialize(msg *message) error {
	s.log.Println("initialize: sending capabilities")
	return s.reply(msg.ID, initializeResult{
		Capabilities: serverCapabilities{
			TextDocumentSync: 1, // full sync
			CompletionProvider: &completionOptions{
				TriggerCharacters: []string{"=", ","},
			},
		},
		ServerInfo: serverInfo{Name: "req42-tracer", Version: "0.1.0"},
	})
}

// --- Document sync ---

type didOpenParams struct {
	TextDocument struct {
		URI  string `json:"uri"`
		Text string `json:"text"`
	} `json:"textDocument"`
}

type didChangeParams struct {
	TextDocument struct {
		URI string `json:"uri"`
	} `json:"textDocument"`
	ContentChanges []struct {
		Text string `json:"text"`
	} `json:"contentChanges"`
}

func (s *Server) handleDidOpen(msg *message) {
	var p didOpenParams
	if err := json.Unmarshal(msg.Params, &p); err != nil {
		s.log.Printf("didOpen parse error: %v", err)
		return
	}
	s.docs[p.TextDocument.URI] = strings.Split(p.TextDocument.Text, "\n")
	s.reloadGraph()
}

func (s *Server) handleDidChange(msg *message) {
	var p didChangeParams
	if err := json.Unmarshal(msg.Params, &p); err != nil {
		s.log.Printf("didChange parse error: %v", err)
		return
	}
	if len(p.ContentChanges) > 0 {
		s.docs[p.TextDocument.URI] = strings.Split(p.ContentChanges[0].Text, "\n")
	}
	s.reloadGraph()
}

// --- Completion ---

type completionParams struct {
	TextDocument struct {
		URI string `json:"uri"`
	} `json:"textDocument"`
	Position struct {
		Line      int `json:"line"`
		Character int `json:"character"`
	} `json:"position"`
}

func (s *Server) handleCompletion(msg *message) error {
	var p completionParams
	if err := json.Unmarshal(msg.Params, &p); err != nil {
		return s.replyError(msg.ID, -32602, "invalid params: "+err.Error())
	}

	lineUpToCursor := ""
	if lines, ok := s.docs[p.TextDocument.URI]; ok {
		line := p.Position.Line
		col := p.Position.Character
		if line >= 0 && line < len(lines) {
			row := lines[line]
			if col > len(row) {
				col = len(row)
			}
			lineUpToCursor = row[:col]
		}
	}

	list := buildCompletions(lineUpToCursor, s.graph)
	return s.reply(msg.ID, list)
}

// --- graph builder helper ---

type graphBuilder struct {
	graph *model.TraceabilityGraph
}

func newGraphBuilder() *graphBuilder {
	return &graphBuilder{graph: &model.TraceabilityGraph{
		Requirements: make(map[string]*model.Requirement),
		ArchElements: make(map[string]*model.ArchElement),
		TestSpecs:    make(map[string]*model.TestSpec),
		TestCodes:    make(map[string]*model.TestCode),
		TestResults:  make(map[string]*model.TestResult),
		Links:        nil,
	}}
}

func (b *graphBuilder) merge(src *model.TraceabilityGraph) {
	for k, v := range src.Requirements {
		b.graph.Requirements[k] = v
	}
	for k, v := range src.ArchElements {
		b.graph.ArchElements[k] = v
	}
	for k, v := range src.TestSpecs {
		b.graph.TestSpecs[k] = v
	}
}
