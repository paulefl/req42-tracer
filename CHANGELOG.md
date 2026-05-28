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

## [0.3.0] — 2025-05-18

### Features

- **feat(watch):** Watch mode with HTTP server and live-reload (Phase 2.4) (#24)
- **feat(report):** ASPICE Dashboard Tab in HTML reports (Phase 2.3) (#23)
