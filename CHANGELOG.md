# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Bug Fixes

- **fix(ci): Benchmark-Job — gh-pages Branch initialisiert** (#55)
  - Leerer Orphan-Branch `gh-pages` remote angelegt (via `git commit-tree`)
  - `github-action-benchmark@v1`: `gh-pages-branch: gh-pages` und `benchmark-data-dir-path: dev/bench` explizit gesetzt
  - Behebt `fatal: couldn't find remote ref gh-pages` bei jedem Master-Push

### CI/CD

- **ci: GoReleaser-Integration für alle GOOS/GOARCH-Kombinationen** (#53)
  - `.goreleaser.yaml` v2: 7 Targets (linux/darwin/windows × amd64/arm64/386, darwin-386 + windows-arm64 excluded)
  - `.tar.gz` für Unix, `.zip` für Windows, `SHA256SUMS` via GoReleaser automatisch generiert
  - SBOM-Generierung via `syft` (in GHA via `anchore/sbom-action/download-syft@v0`)
  - Changelog gruppiert nach `feat`, `fix`, `refactor`
  - `release`-Job in `ci.yml`: manuelles Shell-Script + `softprops/action-gh-release` → `goreleaser/goreleaser-action@v6`
  - `fetch-depth: 0` für vollständigen Git-History-Zugriff, `id-token: write` für SBOM-Attestation

### Features

- **feat(lsp): Phase 3.3 — Hover Tooltips** (#7, PR #32)
  - `textDocument/hover` → Markdown-Tooltip für `req=`, `arch=`, `test-spec=` Werte unter dem Cursor
  - `detectHoverValue()`: findet Attribut-Wert an Cursor-Position (byte-genaue Bereichsprüfung)
  - `buildHoverContent()`: gibt Titel, Text und Metadaten (Priority/Status/ASPICE, impl/parent) zurück
  - `hoverProvider: true` in Capabilities advertised
  - Traceability: `comp.lsp.hover`, `TS-LSP-010/011/012`, `bausteinsicht validate` ✅

### Features

- **feat(lsp): Phase 3.1 — LSP Server Skeleton** (#5, PR #30)
  - `internal/lsp/server.go`: JSON-RPC 2.0 over stdio mit Content-Length-Framing
  - `initialize` → `InitializeResult` (textDocumentSync:1, serverInfo)
  - `shutdown` + `exit` Lifecycle korrekt implementiert (shutdown wartet auf exit)
  - Write-Fehler propagiert; non-fatal read errors werden geloggt und übersprungen
  - Case-insensitive `Content-Length`-Header-Matching
  - `cmd/req42-tracer/lsp.go`: Cobra-Command aktiviert in root.go
  - 3 Tests: initialize handshake, shutdown/exit sequence, unknown method
  - `REVIEW.md`: Traceability-Check Sektion (7-Punkte-Checkliste, Test-Spec-Format)

### Documentation

- **docs(arch): sync architecture.jsonc with arc42.adoc** (#25)
  - Added missing Level-2 subcomponents to `architecture.jsonc`:
    `backend_parser_adoc`, `backend_parser_arch`, `backend_parser_jsonc`,
    `backend_graph_builder`, `backend_graph_analyzer`, `cli_init`
  - Added `comp.trace` and `comp.gaps` arch blocks to `docs/arc42/arc42.adoc`
  - Added `impl=` fields to all new and existing JSONC elements
  - Added new view `Component — Backend Level 2` to `architecture.jsonc`
  - `req42-tracer validate` now passes without warnings; `req42-tracer gaps` shows no gaps
- **docs(arch): fix impl paths and add missing arch elements** (follow-up to #25)
  - Fixed `comp.init` impl path (`internal/cmd` → `cmd/req42-tracer/init.go`)
  - Fixed `comp.lsp`: removed non-existent `impl=internal/lsp`; gap report now honestly flags LSP as unimplemented
  - Added `comp.validate` → `cmd/req42-tracer/validate.go`
  - Added `comp.model` → `internal/model` (shared types + config loading)
  - Added `comp.testresult` → `internal/testresult` (JUnit XML / go-test loader)
  - Mirrored all changes in `architecture.jsonc`; all 15 `impl=` paths verified on disk

## [0.3.0] — 2025-05-18

### Features

- **feat(watch):** Watch mode with HTTP server and live-reload (Phase 2.4) (#24)
- **feat(report):** ASPICE Dashboard Tab in HTML reports (Phase 2.3) (#23)

### Features (Phase 3.4 + 3.5)

- **feat(lsp): Phase 3.4 — Diagnostics** (#8, PR #34)
  - `textDocument/publishDiagnostics` nach jedem `didOpen`/`didChange`
  - Unbekannte `req=`/`arch=`/`test-spec=`-Werte → rote Unterstreichung (error severity)
  - `source: "req42-tracer"`, byte-exakter Range

- **feat(lsp): Phase 3.5 — Go-to-Definition F12** (#9, PR #35)
  - `textDocument/definition` → `Location {file:// URI, 0-based line}`
  - Cross-File-Navigation zu `.adoc` und `_test.go` Definitionen
  - RFC-8089-konformes URI-Format (Unix + Windows)
  - `lineAt()` Helper extrahiert, `TS-LSP-016..019`
