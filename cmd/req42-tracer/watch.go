package main

import (
	"context"
	"encoding/json"
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

var watchGeneration atomic.Int64

func runWatchCmd(cmd *cobra.Command, args []string) error {
	configPath, _ := cmd.Flags().GetString("config")
	openBrowser, _ := cmd.Flags().GetBool("open")
	port, _ := cmd.Flags().GetInt("port")
	outputPath, _ := cmd.Flags().GetString("output")

	config, err := model.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Generating initial report...")
	if err := watchGenerateReport(config, outputPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: initial generation failed: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Report generated: %s\n", outputPath)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	for _, dir := range []string{"docs/requirements", "docs/arc42", "docs"} {
		if _, err := os.Stat(dir); err == nil {
			_ = watcher.Add(dir)
		}
	}
	_ = watcher.Add(configPath)

	var (
		mu    sync.Mutex
		timer *time.Timer
	)
	debounce := func() {
		mu.Lock()
		defer mu.Unlock()
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(500*time.Millisecond, func() {
			fmt.Fprintln(os.Stderr, "Change detected, regenerating...")
			if err := watchGenerateReport(config, outputPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			} else {
				watchGeneration.Add(1)
				fmt.Fprintln(os.Stderr, "Report updated")
			}
		})
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) {
					debounce()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
			case <-ctx.Done():
				return
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
		_ = json.NewEncoder(w).Encode(map[string]int64{"generation": watchGeneration.Load()})
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

	server := &http.Server{Addr: fmt.Sprintf(":%d", port), Handler: mux}
	return server.ListenAndServe()
}

func watchGenerateReport(config *model.Config, outputPath string) error {
	builder := graph.NewBuilder()

	if g, err := parser.ParseAllFromDir("docs/requirements", "software"); err == nil {
		_ = builder.MergeGraph(g)
	}
	if g, err := parser.ParseAllFromDir("docs/arc42", "software"); err == nil {
		_ = builder.MergeGraph(g)
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

	summaryPath := filepath.Join(filepath.Dir(outputPath), "summary.html")
	return htmlReporter.GenerateSummaryReport(summaryPath)
}

const liveReloadScript = `<script>(function(){var g=null;function poll(){fetch('/api/generation').then(function(r){return r.json();}).then(function(d){if(g===null){g=d.generation;}else if(d.generation!==g){location.reload();}setTimeout(poll,2000);}).catch(function(){setTimeout(poll,5000);});}poll();})();</script>`

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
