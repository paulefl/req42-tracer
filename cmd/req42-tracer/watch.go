package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/paulefl/req42-tracer/internal/graph"
	"github.com/paulefl/req42-tracer/internal/model"
	"github.com/paulefl/req42-tracer/internal/parser"
	"github.com/paulefl/req42-tracer/internal/report"
	"github.com/spf13/cobra"
)

func newWatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Watch files and serve live-reloading HTML report",
		Long: `Watch documentation files for changes and automatically regenerate
the HTML report. Starts an HTTP server with live-reload so the browser
refreshes automatically when sources change.`,
		RunE: runWatchCmd,
	}

	cmd.Flags().Bool("open", false, "Open report in browser on start")
	cmd.Flags().Int("port", 8042, "HTTP server port")
	cmd.Flags().String("output", "reports/traceability-report.html", "Output path for HTML report")

	return cmd
}

func runWatchCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	openBrowser, _ := cmd.Flags().GetBool("open")
	port, _ := cmd.Flags().GetInt("port")
	outputPath, _ := cmd.Flags().GetString("output")

	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Respect config output path when --output flag was not explicitly set
	if !cmd.Flags().Changed("output") && config.Reports.HTML.Output != "" {
		outputPath = config.Reports.HTML.Output
	}

	// Derive watch/parse directories from config instead of hardcoding
	docsPath := "docs"
	if p, ok := config.Projects["software"]; ok && p.Docs != "" {
		docsPath = p.Docs
	}
	reqDir := filepath.Join(docsPath, "requirements")
	arcDir := filepath.Join(docsPath, "arc42")

	fmt.Fprintln(os.Stderr, "Generating initial report...")
	if err := watchGenerateReport(config, outputPath, reqDir, arcDir); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: initial generation failed: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Report generated: %s\n", outputPath)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	for _, dir := range []string{reqDir, arcDir, docsPath} {
		if _, err := os.Stat(dir); err == nil {
			_ = watcher.Add(dir)
		}
	}
	_ = watcher.Add(configPath)

	// Local counter — avoids package-level state bleed between invocations/tests
	var generation atomic.Int64

	var (
		mu    sync.Mutex
		genMu sync.Mutex // serializes report generation; prevents concurrent file writes
		timer *time.Timer
	)
	debounce := func() {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(500*time.Millisecond, func() {
			genMu.Lock()
			defer genMu.Unlock()
			fmt.Fprintln(os.Stderr, "Change detected, regenerating...")
			// Reload config on every regeneration so config-file edits take effect
			cfg, cfgErr := model.LoadConfig(configPath)
			if cfgErr != nil {
				fmt.Fprintf(os.Stderr, "Warning: config reload failed: %v\n", cfgErr)
				cfg = config
			}
			if err := watchGenerateReport(cfg, outputPath, reqDir, arcDir); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			} else {
				generation.Add(1)
				fmt.Fprintln(os.Stderr, "Report updated")
			}
		})
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Include Rename: many editors (vim, emacs, JetBrains) write via atomic rename
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) ||
					event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
					debounce()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
			}
		}
	}()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprint(w, injectLiveReload(string(data)))
	})

	mux.HandleFunc("/api/generation", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int64{"generation": generation.Load()})
	})

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Printf("Watching for changes. Report available at %s\n", url)
	fmt.Println("Press Ctrl+C to stop.")

	if openBrowser {
		go func() {
			time.Sleep(500 * time.Millisecond)
			watchOpenBrowser(url)
		}()
	}

	// Bind to loopback only — report contains confidential project documentation
	server := &http.Server{
		Addr:         fmt.Sprintf("127.0.0.1:%d", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func watchGenerateReport(config *model.Config, outputPath, reqDir, arcDir string) error {
	builder := graph.NewBuilder()

	if g, err := parser.ParseAllFromDir(reqDir, "software"); err == nil {
		if err := builder.MergeGraph(g); err != nil {
			return fmt.Errorf("requirements merge: %w", err)
		}
	}
	if g, err := parser.ParseAllFromDir(arcDir, "software"); err == nil {
		if err := builder.MergeGraph(g); err != nil {
			return fmt.Errorf("architecture merge: %w", err)
		}
	}

	builder.DeriveASPICELevels()
	if err := builder.BuildLinks(); err != nil {
		return err
	}

	analyzer := graph.NewAnalyzer(builder.GetGraph())
	htmlReporter := report.NewHTMLReporter(analyzer, config, outputPath)

	if err := htmlReporter.GenerateReport(); err != nil {
		return err
	}

	// Summary failure is non-fatal — main report is already written and browser reload still fires
	summaryPath := filepath.Join(filepath.Dir(outputPath), "summary.html")
	if err := htmlReporter.GenerateSummaryReport(summaryPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: summary report failed: %v\n", err)
	}
	return nil
}

const liveReloadScript = `<style>#req42-live-badge{position:fixed;bottom:12px;right:12px;padding:4px 10px;border-radius:12px;font-size:12px;font-family:monospace;z-index:9999;background:#27ae60;color:#fff;opacity:.85;transition:background .3s;}#req42-live-badge.rebuilding{background:#e67e22;}</style><div id="req42-live-badge">● Live</div><script>(function(){var g=null;var badge=document.getElementById('req42-live-badge');function poll(){fetch('/api/generation').then(function(r){return r.json();}).then(function(d){if(g===null){g=d.generation;}else if(d.generation!==g){if(badge){badge.textContent='⟳ Rebuilding…';badge.className='rebuilding';}setTimeout(function(){location.reload();},200);}setTimeout(poll,2000);}).catch(function(){setTimeout(poll,5000);});}poll();})();</script>`

func injectLiveReload(content string) string {
	return strings.Replace(content, "</body>", liveReloadScript+"</body>", 1)
}

func watchOpenBrowser(url string) {
	var cmdName string
	var cmdArgs []string
	switch runtime.GOOS {
	case "darwin":
		cmdName, cmdArgs = "open", []string{url}
	case "windows":
		cmdName, cmdArgs = "rundll32", []string{"url.dll,FileProtocolHandler", url}
	default:
		cmdName, cmdArgs = "xdg-open", []string{url}
	}
	_ = exec.Command(cmdName, cmdArgs...).Start()
}
