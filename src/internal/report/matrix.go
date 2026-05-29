package report

import (
	"fmt"
	"sort"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// MatrixCell represents a single cell in the traceability matrix.
type MatrixCell struct {
	Status   string `json:"status"`   // "covered", "missing", "stale"
	Evidence string `json:"evidence"` // Details about the coverage
}

// MatrixRow represents a single requirement row in the matrix.
type MatrixRow struct {
	RequirementID string
	Title         string
	Priority      string
	Status        string
	Cells         map[string]MatrixCell // column_id -> MatrixCell
}

// MatrixData represents the complete traceability matrix.
type MatrixData struct {
	Rows        []MatrixRow `json:"rows"`
	Columns     []MatrixColumn `json:"columns"`
	Matrix      map[string]map[string]MatrixCell `json:"matrix"` // req_id -> col_id -> cell
	Statistics  MatrixStats `json:"statistics"`
}

// MatrixColumn represents a column in the matrix.
type MatrixColumn struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"` // "arch", "test-spec"
}

// MatrixStats contains coverage statistics.
type MatrixStats struct {
	TotalRequirements     int     `json:"total_requirements"`
	CoveredRequirements   int     `json:"covered_requirements"`
	MissingRequirements   int     `json:"missing_requirements"`
	StaleLinks            int     `json:"stale_links"`
	CoveragePercentage    float64 `json:"coverage_percentage"`
	AverageCoveragePerReq float64 `json:"average_coverage_per_req"`
}

// BuildMatrixData constructs a traceability matrix from the graph.
func BuildMatrixData(g *model.TraceabilityGraph) *MatrixData {
	data := &MatrixData{
		Rows:    make([]MatrixRow, 0),
		Columns: make([]MatrixColumn, 0),
		Matrix:  make(map[string]map[string]MatrixCell),
	}

	// Build columns: Architecture elements first, then test specs
	columnOrder := make(map[string]int) // col_id -> order

	// Add architecture columns
	for id, arch := range g.ArchElements {
		data.Columns = append(data.Columns, MatrixColumn{
			ID:    id,
			Title: fmt.Sprintf("%s: %s", id, arch.Title),
			Type:  "arch",
		})
		columnOrder[id] = len(data.Columns) - 1
	}

	// Add test spec columns
	for id, spec := range g.TestSpecs {
		data.Columns = append(data.Columns, MatrixColumn{
			ID:    id,
			Title: fmt.Sprintf("%s: %s", id, spec.Title),
			Type:  "test-spec",
		})
		columnOrder[id] = len(data.Columns) - 1
	}

	// Build rows: Requirements
	for _, req := range g.Requirements {
		row := MatrixRow{
			RequirementID: req.ID,
			Title:         req.Title,
			Priority:      req.Priority,
			Status:        req.Status,
			Cells:         make(map[string]MatrixCell),
		}

		// Initialize all cells as missing
		for _, col := range data.Columns {
			row.Cells[col.ID] = MatrixCell{
				Status:   "missing",
				Evidence: "No trace link found",
			}
		}

		// Fill cells based on trace links
		for _, link := range g.Links {
			if link.FromID == req.ID && link.FromType == "requirement" {
				status := "covered"
				if link.Status == "stale" {
					status = "stale"
				}

				if cell, exists := row.Cells[link.ToID]; exists {
					cell.Status = status
					cell.Evidence = fmt.Sprintf("%s (%s)", link.LinkType, link.Reason)
					row.Cells[link.ToID] = cell
				}
			}
		}

		data.Rows = append(data.Rows, row)
		data.Matrix[req.ID] = row.Cells
	}

	// Sort rows by requirement ID
	sort.Slice(data.Rows, func(i, j int) bool {
		return data.Rows[i].RequirementID < data.Rows[j].RequirementID
	})

	// Calculate statistics
	data.calculateStatistics()

	return data
}

// FilterMatrix returns a filtered matrix by priority and status.
func FilterMatrix(data *MatrixData, priorities []string, statuses []string) *MatrixData {
	filtered := &MatrixData{
		Rows:       make([]MatrixRow, 0),
		Columns:    data.Columns,
		Matrix:     make(map[string]map[string]MatrixCell),
	}

	priorityMap := make(map[string]bool)
	statusMap := make(map[string]bool)

	for _, p := range priorities {
		priorityMap[p] = true
	}
	for _, s := range statuses {
		statusMap[s] = true
	}

	for _, row := range data.Rows {
		if priorityMap[row.Priority] && statusMap[row.Status] {
			filtered.Rows = append(filtered.Rows, row)
			filtered.Matrix[row.RequirementID] = row.Cells
		}
	}

	filtered.calculateStatistics()

	return filtered
}

// calculateStatistics computes coverage metrics for the matrix.
func (m *MatrixData) calculateStatistics() {
	m.Statistics.TotalRequirements = len(m.Rows)

	var coveredCount int
	var totalCells int
	var cellsCovered int

	for _, row := range m.Rows {
		rowCovered := false
		for _, cell := range row.Cells {
			totalCells++
			if cell.Status == "covered" {
				cellsCovered++
				rowCovered = true
			} else if cell.Status == "stale" {
				m.Statistics.StaleLinks++
			}
		}

		if rowCovered {
			coveredCount++
		}
	}

	m.Statistics.CoveredRequirements = coveredCount
	m.Statistics.MissingRequirements = m.Statistics.TotalRequirements - coveredCount

	if m.Statistics.TotalRequirements > 0 {
		m.Statistics.CoveragePercentage = (float64(coveredCount) / float64(m.Statistics.TotalRequirements)) * 100
		m.Statistics.AverageCoveragePerReq = float64(cellsCovered) / float64(totalCells) * 100
	}
}

// ExportMatrixToCSV generates CSV content from the matrix.
func ExportMatrixToCSV(data *MatrixData) string {
	csv := "Requirement,Priority,Status"

	// Add column headers
	for _, col := range data.Columns {
		csv += "," + col.ID
	}
	csv += "\n"

	// Add data rows
	for _, row := range data.Rows {
		csv += fmt.Sprintf("%s,%s,%s", row.RequirementID, row.Priority, row.Status)

		for _, col := range data.Columns {
			status := "✗"
			if cell, exists := row.Cells[col.ID]; exists {
				if cell.Status == "covered" {
					status = "✓"
				} else if cell.Status == "stale" {
					status = "⚠"
				}
			}
			csv += "," + status
		}
		csv += "\n"
	}

	return csv
}
