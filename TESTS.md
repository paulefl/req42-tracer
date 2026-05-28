# Test Conventions

Dieses Dokument definiert die verbindlichen Regeln für Tests in diesem Repository.

## Pflicht: Test-Spec Block + ASPICE-Linkage

Jeder Test **muss** mit einem `[test-spec]`-Block verknüpft sein. Dieser Block:
- referenziert die Anforderung (`req=`) die der Test verifiziert
- benennt den ASPICE Best-Practice (`aspice=`) dem er zugeordnet ist
- steht als Kommentar direkt oberhalb der Testfunktion

### Format

```go
// [test-spec,id=TS-PKG-NNN,req="REQ-PKG-NNN",aspice="SWE.X.BPY"]
// Test: Kurzbeschreibung was getestet wird
// [end]
func TestFunctionName(t *testing.T) {
    ...
}
```

### Felder

| Feld | Pflicht | Beschreibung |
|---|---|---|
| `id` | ✅ | Eindeutige Test-Spec-ID, Schema: `TS-<PAKET>-<NNN>` |
| `req` | ✅ | Anforderungs-ID aus `docs/requirements/req42.adoc` |
| `aspice` | ✅ | ASPICE BP z.B. `SWE.5.BP3` (Establish Traceability between tests and requirements) |

### ASPICE-Zuordnung für Tests

| ASPICE BP | Bedeutung | Wann verwenden |
|---|---|---|
| `SWE.5.BP3` | Traceability zwischen Tests und Anforderungen | Alle Unit-Tests |
| `SWE.5.BP2` | Test-Implementierung | Integrations- und End-to-End-Tests |
| `SWE.5.BP4` | Test-Ausführung | Tests die Ausführungsverhalten prüfen |

## Coverage-Ziele (Phase 4)

| Paket | Ziel |
|---|---|
| `internal/parser` | ≥ 80% |
| `internal/graph` | ≥ 80% |
| `internal/aspice` | ≥ 75% |
| `internal/report` | ≥ 70% |
| `internal/model` | ≥ 60% |
| `internal/testresult` | ≥ 70% |

Coverage prüfen:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

## Dateinamen-Konvention

```
adoc.go       →  adoc_test.go
aspice_data.go →  aspice_data_test.go
matrix.go     →  matrix_test.go
```

## Beispiel (vollständig)

```go
// [test-spec,id=TS-RPT-003,req="REQ-RPT-011",aspice="SWE.5.BP3"]
// Test: getCoverageLevel gibt "good" bei ≥80%, "warning" bei ≥50%, sonst "danger"
// [end]
func TestGetCoverageLevel(t *testing.T) {
    cases := []struct {
        pct      float64
        expected string
    }{
        {100.0, "good"},
        {80.0, "good"},
        {79.9, "warning"},
        {50.0, "warning"},
        {49.9, "danger"},
        {0.0, "danger"},
    }
    for _, tc := range cases {
        got := getCoverageLevel(tc.pct)
        if got != tc.expected {
            t.Errorf("getCoverageLevel(%.1f) = %q, want %q", tc.pct, got, tc.expected)
        }
    }
}
```

## Warum diese Konvention?

`req42-tracer` ist ein ASPICE-Tracing-Tool — es muss seine eigenen Anforderungen nachvollziehbar testen (Dogfooding). Die `[test-spec]`-Blöcke werden vom Tool selbst geparst und in den Traceability-Graph eingebunden, sodass `req42-tracer trace` auch die Test-Coverage des Tools selbst visualisiert.
