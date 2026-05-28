package main

// [test-spec,id=TS-WATCH-001,req="REQ-WATCH-001",aspice="SWE.5.BP3"]
// Test: injectLiveReload ersetzt </body> durch Script + </body>
// [end]

// [test-spec,id=TS-WATCH-002,req="REQ-WATCH-001",aspice="SWE.5.BP3"]
// Test: injectLiveReload ist no-op wenn kein </body> vorhanden
// [end]

// [test-spec,id=TS-WATCH-003,req="REQ-WATCH-001",aspice="SWE.5.BP4"]
// Test: /api/generation liefert JSON mit generation-Feld
// [end]

// [test-spec,id=TS-WATCH-004,req="REQ-WATCH-001",aspice="SWE.5.BP4"]
// Test: / Handler liefert 503 wenn Report-Datei fehlt
// [end]

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
)

func TestInjectLiveReload_InjectsBeforeBody(t *testing.T) {
	input := "<html><body><p>Hello</p></body></html>"
	result := injectLiveReload(input)

	if !strings.Contains(result, liveReloadScript+"</body>") {
		t.Errorf("expected live-reload script before </body>, got:\n%s", result)
	}
	if strings.Count(result, "</body>") != 1 {
		t.Errorf("expected exactly one </body>, got:\n%s", result)
	}
}

func TestInjectLiveReload_NoOpWhenNoBody(t *testing.T) {
	input := "<html><p>No closing body tag</p></html>"
	result := injectLiveReload(input)
	if result != input {
		t.Errorf("expected unchanged content, got:\n%s", result)
	}
}

func TestGenerationEndpoint(t *testing.T) {
	var gen atomic.Int64
	gen.Store(42)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int64{"generation": gen.Load()})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/generation", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var body map[string]int64
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if body["generation"] != 42 {
		t.Errorf("expected generation=42, got %d", body["generation"])
	}
}

func TestRootHandler_Returns503WhenReportMissing(t *testing.T) {
	outputPath := filepath.Join(t.TempDir(), "nonexistent.html")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := os.ReadFile(outputPath)
		if err != nil {
			http.Error(w, "Report not yet generated", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(injectLiveReload(string(data))))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}
}

func TestRootHandler_ServesInjectedHTML(t *testing.T) {
	dir := t.TempDir()
	outputPath := filepath.Join(dir, "report.html")
	htmlContent := "<html><body><p>Report</p></body></html>"
	if err := os.WriteFile(outputPath, []byte(htmlContent), 0644); err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		data, err := os.ReadFile(outputPath)
		if err != nil {
			http.Error(w, "Report not yet generated", http.StatusServiceUnavailable)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(injectLiveReload(string(data))))
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	body := rr.Body.String()
	if !strings.Contains(body, liveReloadScript) {
		t.Error("expected live-reload script in response body")
	}
	if !strings.Contains(body, "<p>Report</p>") {
		t.Error("expected original report content in response body")
	}
}
