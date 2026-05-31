# Requirements Review Protocol

ASPICE SWE.1 BP4 — Requirements shall be reviewed for consistency and testability.

## Review Rounds

### Round 1 — Initial Review (2026-05-07)

| Field | Value |
|---|---|
| Date | 2026-05-07 |
| Reviewer | paulefl |
| Document | `docs/requirements/req42.adoc` |
| Version | all requirements at version=1 |
| Scope | All 16 REQ-* requirements |

**Outcome:** Approved — all requirements are testable, non-ambiguous, and consistently structured.

**Checklist:**
- [x] All requirements have unique IDs (REQ-*)
- [x] All requirements have priority and status attributes
- [x] All requirements have at least one Acceptance Criterion
- [x] All requirements have a corresponding test-spec block
- [x] All requirements reference a relevant ASPICE process area
- [x] No contradictory requirements detected
- [x] REQ-LSP-001 marked `status=draft` (LSP is MVP, not fully specified)

**Open Items:**
- REQ-LSP-001: status=draft — finalize when LSP feature set is confirmed
- All REQ-PERF-001 acceptance criteria are measurable but depend on target hardware

---

### Round 2 — SYS-Level Addition (2026-05-31)

| Field | Value |
|---|---|
| Date | 2026-05-31 |
| Reviewer | paulefl |
| Document | `docs/requirements/req42.adoc`, `docs/requirements/sys-requirements.adoc` |
| Scope | SYS-level requirements added; `derives=SYS-*` links added to all REQ-* |

**Outcome:** Approved — SYS→SW derivation chain established, ASPICE SWE.1 BP1 addressed.

**Checklist:**
- [x] All 16 SW requirements have `derives=SYS-*` link
- [x] All SYS requirements have SW derivatives listed
- [x] SYS requirements cover all functional areas (Traceability, ASPICE, Tooling, Quality, Config, Validation, Testing)
- [x] No SW requirement is orphaned (all derive from at least one SYS requirement)
- [x] SYS requirements are reviewed and approved by paulefl

**Changes made:**
- Added `sys-requirements.adoc` with 7 SYS-level requirements
- Added `derives=SYS-*` to all 16 REQ-* requirements in `req42.adoc`
