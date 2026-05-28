# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

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
