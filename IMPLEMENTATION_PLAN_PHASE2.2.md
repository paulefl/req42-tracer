# Implementierungsplan: Phase 2.2 — Matrix View Tab for HTML Reports

## 1. Überblick
Ergänze den HTML-Report um einen zweiten Tab mit einer interaktiven Traceability-Matrix. Die Matrix zeigt Requirements als Zeilen und Architecture-Elemente sowie Test-Spezifikationen als Spalten. Jede Zelle ist farbcodiert nach Coverage-Status.

## 2. Acceptance Criteria
- [x] Matrix layout: rows=requirements, columns=architecture/tests
- [x] Cells color-coded: green=covered, red=missing, orange=stale
- [x] Sortable columns, filterable by priority/status
- [x] Export to CSV

## 3. Implementierungsschritte

### Phase 2.2.1: Matrix Data Structure
- [ ] Funktion `BuildMatrixData()` in `report/matrix.go`
- [ ] Datenstruktur für Matrix mit Requirements × Architecture
- [ ] Coverage-Status pro Cell: covered, missing, stale
- [ ] JSON-Export für JavaScript

### Phase 2.2.2: Matrix HTML Template
- [ ] Zweiter Tab in HTML-Report ("Matrix" neben "Graph")
- [ ] HTML-Tabelle mit Bootstrap-Styling
- [ ] Responsives Design (horizontal scroll auf Mobile)
- [ ] Farb-Legende für Coverage-Status

### Phase 2.2.3: Matrix Interaktivität (JavaScript)
- [ ] Spalten-Sortierung (A-Z, Coverage-%)
- [ ] Zeilen-Filter nach Priority (high/medium/low)
- [ ] Status-Filter (approved/draft/deprecated)
- [ ] Such-Feld für Requirements

### Phase 2.2.4: CSV Export
- [ ] Button "Export to CSV"
- [ ] CSV-Format: Header + Requirements + Coverage-Daten
- [ ] Download im Browser (JavaScript Blob)

### Phase 2.2.5: CLI Integration
- [ ] Keine neuen Flags nötig (Matrix wird automatisch generiert)
- [ ] Matrix-Daten in HTML einbetten
- [ ] Optional: CSV-Export via CLI (`--format=csv`)

### Phase 2.2.6: Testing & Validation
- [ ] Unit Tests für Matrix-Daten
- [ ] Browser-Test (Sortierung, Filter, Export)
- [ ] Edge Cases: Leere Spalten, große Matrizen

## 4. Technologie-Stack
- **Frontend:** Vanilla JavaScript, HTML5, CSS3
- **Styling:** Bootstrap oder Custom CSS
- **Export:** JavaScript Blob API für CSV
- **Backend:** Go JSON Export aus TraceabilityGraph

## 5. Datenstruktur (MatrixData)

```go
type MatrixCell struct {
    Status   string // "covered", "missing", "stale"
    Evidence string // Details warum dieser Status
}

type MatrixData struct {
    Requirements []Requirement          // Zeilen
    ArchElements []ArchElement          // Spalten
    TestSpecs    []TestSpec             // Zusätzliche Spalten
    Matrix       map[string]map[string]MatrixCell // req_id -> arch_id -> status
}
```

## 6. Farb-Schema
- 🟢 **Green (#7ed321):** Covered (Requirement → Architecture → Test)
- 🔴 **Red (#e74c3c):** Missing (Requirement ohne Architecture)
- 🟠 **Orange (#f5a623):** Stale (Veraltete Versionen)
- ⚪ **Gray (#ccc):** N/A (nicht anwendbar)

## 7. CSV-Format
```
Requirement,Priority,Status,arch-001,arch-002,spec-001
REQ-001,high,approved,✓,✗,✓
REQ-002,medium,draft,✗,✓,✓
...
```

## 8. Erfolgs-Kriterien
- ✅ Matrix wird in HTML-Report angezeigt
- ✅ Sortierung funktioniert (5+ Spalten)
- ✅ Filter nach Priority/Status funktioniert
- ✅ CSV-Export generiert valid CSV
- ✅ Große Matrizen (100+ rows) sind performant
- ✅ Responsive auf Tablets/Mobiles

## 9. Nächste Schritte
- Phase 2.3: ASPICE Dashboard Tab
- Phase 2.4: Watch Mode with Live-Reload
