# Implementierungsplan: Phase 2 — HTML Report Generator with D3.js Graph

## 1. Überblick
Erstelle einen interaktiven HTML-Report mit D3.js-Graph-Visualisierung für die Traceability-Matrix. Der Report zeigt Requirements, Architecture-Elemente, Test-Spezifikationen und Test-Ergebnisse als Knoten mit Kanten (satisfies, implements, verifies, derives).

## 2. Acceptance Criteria
- [x] Graph view: D3.js Visualization der Traceability-Knoten
- [x] Node colors: blue=req, green=arch, orange=test-spec, gray=test-result
- [x] Edge labels: satisfies, implements, verifies, derives
- [x] Interactive: hover für Details, click zum Filtern

## 3. Implementierungsschritte

### Phase 2.1: Graph Data Export (internal/report/graph.go)
- [ ] Struktur `GraphData` mit Nodes und Edges definieren
- [ ] Funktion `ExportGraphData(graph *graph.Graph) *GraphData` implementieren
- [ ] JSON-Serialisierung für D3.js vorbereiten
- **Abhängigkeit:** Funktioniert mit bestehendem `graph.Graph` aus `internal/graph`

### Phase 2.2: HTML Template mit D3.js (internal/templates/report.html)
- [ ] HTML-Template mit D3.js Script erstellen
- [ ] Graph-Rendering mit force-directed layout
- [ ] Node-Styling basierend auf Typ
- [ ] Edge-Labels und Pfeile
- **Dependencies:** D3.js v7.0+ (CDN), CSS für Styling

### Phase 2.3: Report Generator (internal/report/html.go)
- [ ] Funktion `GenerateHTMLReport(graph *graph.Graph, outputPath string)` implementieren
- [ ] Graph-Daten einbetten in Template
- [ ] CSS und JavaScript minifizieren (optional)
- **Dependencies:** `html/template`, `graph.Graph`

### Phase 2.4: Interaktivität (internal/templates/report-interactions.js)
- [ ] Hover-Tooltips für Node-Details
- [ ] Click-Filter nach Node-Typ
- [ ] Zoom und Pan Support
- [ ] Legend für Farben und Edge-Types

### Phase 2.5: Integration in CLI (cmd/req42-tracer/report.go)
- [ ] Report-Command aktualisieren für HTML-Output
- [ ] Flag `--format=html` unterstützen
- [ ] Output-Pfad konfigurierbar

### Phase 2.6: Tests und Demo
- [ ] Unit Tests für Graph-Export
- [ ] Integration Test mit Demo-Projekt
- [ ] Manuelle Validierung im Browser

## 4. Technologie-Stack
- **Frontend:** D3.js v7+ (Force-Directed Layout)
- **Backend:** Go `html/template`, JSON
- **Styling:** Embedded CSS
- **Interaktivität:** Vanilla JavaScript

## 5. Datenstruktur (GraphData)

```go
type Node struct {
    ID       string `json:"id"`
    Label    string `json:"label"`
    Type     string `json:"type"` // "req", "arch", "test-spec", "test-result"
    Metadata map[string]interface{} `json:"metadata"`
}

type Edge struct {
    Source string `json:"source"`
    Target string `json:"target"`
    Label  string `json:"label"` // "satisfies", "implements", "verifies", "derives"
}

type GraphData struct {
    Nodes []Node `json:"nodes"`
    Edges []Edge `json:"edges"`
}
```

## 6. Risiken und Mitigation
- **Große Graphen:** D3.js kann mit 1000+ Nodes langsam werden
  - Mitigation: Quadtree oder Canvas-Rendering verwenden
- **Browser-Kompatibilität:** Ältere Browser unterstützen D3 nicht
  - Mitigation: Moderne Browser-Versionen voraussetzen (ES6+)

## 7. Erfolgs-Kriterien
- ✅ HTML-Report wird generiert ohne Fehler
- ✅ Graph mit allen Knoten und Kanten sichtbar
- ✅ Farben korrekt entsprechend Node-Typ
- ✅ Hover-Tooltips funktionieren
- ✅ Click-Filter funktioniert
- ✅ Report bei Demo-Projekt getestet
