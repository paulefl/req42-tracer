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
// exit=true means the Run loop should stop.
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
		// Keep reading: client must send exit notification next.
	case "exit":
		return true, nil
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
	TextDocumentSync int `json:"textDocumentSync"`
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
		},
		ServerInfo: serverInfo{Name: "req42-tracer", Version: "0.1.0"},
	})
}
