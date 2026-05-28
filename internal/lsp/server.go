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
)

// Server is a minimal LSP server over stdio using JSON-RPC 2.0.
type Server struct {
	in  *bufio.Reader
	out io.Writer
	log *log.Logger
}

// NewServer creates a server that reads from stdin and writes to stdout.
// Diagnostic messages go to stderr.
func NewServer() *Server {
	return &Server{
		in:  bufio.NewReader(os.Stdin),
		out: os.Stdout,
		log: log.New(os.Stderr, "[lsp] ", log.LstdFlags),
	}
}

// Run processes incoming JSON-RPC messages until shutdown/exit or EOF.
func (s *Server) Run() error {
	s.log.Println("LSP server started on stdio")
	for {
		msg, err := s.readMessage()
		if err == io.EOF {
			s.log.Println("stdin closed, shutting down")
			return nil
		}
		if err != nil {
			return fmt.Errorf("read: %w", err)
		}
		if done := s.dispatch(msg); done {
			return nil
		}
	}
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
		if strings.HasPrefix(line, "Content-Length: ") {
			n, err := strconv.Atoi(strings.TrimPrefix(line, "Content-Length: "))
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

func (s *Server) send(v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		s.log.Printf("marshal error: %v", err)
		return
	}
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	_, _ = io.WriteString(s.out, header)
	_, _ = s.out.Write(data)
}

func (s *Server) reply(id interface{}, result interface{}) {
	s.send(message{JSONRPC: "2.0", ID: id, Result: result})
}

func (s *Server) replyError(id interface{}, code int, msg string) {
	s.send(message{JSONRPC: "2.0", ID: id, Error: &rpcError{Code: code, Message: msg}})
}

// --- Dispatch ---

func (s *Server) dispatch(msg *message) (shutdown bool) {
	s.log.Printf("← %s (id=%v)", msg.Method, msg.ID)
	switch msg.Method {
	case "initialize":
		s.handleInitialize(msg)
	case "initialized":
		// notification, no response required
	case "shutdown":
		s.reply(msg.ID, nil)
		return true
	case "exit":
		return true
	default:
		if msg.ID != nil {
			s.replyError(msg.ID, -32601, "method not found: "+msg.Method)
		}
	}
	return false
}

// --- initialize ---

type initializeResult struct {
	Capabilities serverCapabilities `json:"capabilities"`
	ServerInfo   serverInfo         `json:"serverInfo"`
}

type serverCapabilities struct {
	TextDocumentSync   int  `json:"textDocumentSync"`
	CompletionProvider bool `json:"completionProvider,omitempty"`
	HoverProvider      bool `json:"hoverProvider,omitempty"`
	DefinitionProvider bool `json:"definitionProvider,omitempty"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func (s *Server) handleInitialize(msg *message) {
	s.log.Println("initialize: sending capabilities")
	s.reply(msg.ID, initializeResult{
		Capabilities: serverCapabilities{
			TextDocumentSync: 1, // full sync
		},
		ServerInfo: serverInfo{Name: "req42-tracer", Version: "0.1.0"},
	})
}
