package parser

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// ParseGoTestFiles walks dirPath recursively, reads all *_test.go files and
// extracts [test-spec] comment annotations placed directly before func Test*
// declarations. Each annotation produces a TestCode entry linking the test
// function to the spec ID.
func ParseGoTestFiles(dirPath, project string) (*model.TraceabilityGraph, error) {
	graph := &model.TraceabilityGraph{
		Requirements:   make(map[string]*model.Requirement),
		ArchElements:   make(map[string]*model.ArchElement),
		DesignElements: make(map[string]*model.DesignElement),
		TestSpecs:      make(map[string]*model.TestSpec),
		TestCodes:      make(map[string]*model.TestCode),
		TestResults:    make(map[string]*model.TestResult),
		Links:          []*model.TraceLink{},
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, "_test.go") {
			return err
		}
		codes, parseErr := parseGoTestFile(path, project)
		if parseErr == nil {
			for _, code := range codes {
				graph.TestCodes[code.ID] = code
			}
		}
		return nil
	})

	return graph, err
}

// parseGoTestFile extracts TestCode entries from a single Go test file.
func parseGoTestFile(filePath, project string) ([]*model.TestCode, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var codes []*model.TestCode
	var pendingSpecID string // spec ID from last [test-spec] comment seen

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Detect [test-spec] comment: // [test-spec,id=...,...]
		if strings.HasPrefix(line, "// [test-spec,") {
			inner := line[3:] // strip "// "
			attrs := extractAttributes(inner)
			if id, ok := attrs["id"]; ok && id != "" {
				pendingSpecID = id
			}
			continue
		}

		// Detect test function declaration
		if strings.HasPrefix(line, "func Test") && pendingSpecID != "" {
			funcName := extractFuncName(line)
			if funcName != "" {
				code := &model.TestCode{
					ID:       fmt.Sprintf("%s::%s", filePath, funcName),
					TestSpec: pendingSpecID,
					File:     filePath,
					Function: funcName,
					Language: "go",
					Project:  project,
				}
				codes = append(codes, code)
			}
			pendingSpecID = ""
			continue
		}

		// Reset pending spec if a non-comment, non-empty line appears
		// (guards against spec comments separated from the function by blank lines or other code)
		if line != "" && !strings.HasPrefix(line, "//") && !strings.HasPrefix(line, "func Test") {
			pendingSpecID = ""
		}
	}

	return codes, scanner.Err()
}

// extractFuncName extracts the function name from a "func TestXxx(...)" line.
func extractFuncName(line string) string {
	// line looks like: func TestFoo(t *testing.T) {
	after := strings.TrimPrefix(line, "func ")
	paren := strings.Index(after, "(")
	if paren < 0 {
		return ""
	}
	return strings.TrimSpace(after[:paren])
}
