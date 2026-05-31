package testresult

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PackageCoverage holds statement-level coverage data for one package or file.
type PackageCoverage struct {
	Package    string  // Go package path or Cobertura/LCOV source path
	Statements int     // total statements
	Covered    int     // covered statements
	Pct        float64 // coverage percentage (0–100)
}

// ParseGoCoverage parses the output of `go tool cover -func=coverage.out`.
// Each line has the form:
//
//	github.com/foo/bar/pkg/file.go:FuncName	75.0%
//
// or the total line:
//
//	total:	(statements)	75.0%
//
// Returns one PackageCoverage per package (aggregated across functions).
func ParseGoCoverage(filePath string) ([]PackageCoverage, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open go coverage file %s: %w", filePath, err)
	}
	defer f.Close()

	// Aggregate per package: covered + total
	type accum struct{ stmts, covered int }
	pkgMap := make(map[string]*accum)
	var order []string // preserve insertion order

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "total:") {
			continue
		}
		// Format: path/file.go:FuncName\tpct%
		// Split on last tab
		tabIdx := strings.LastIndex(line, "\t")
		if tabIdx < 0 {
			continue
		}
		fileFunc := line[:tabIdx]
		pctStr := strings.TrimSuffix(strings.TrimSpace(line[tabIdx+1:]), "%")

		pct, err := strconv.ParseFloat(pctStr, 64)
		if err != nil {
			continue
		}

		// Extract package from file path (everything before the last '/')
		pkgPath := fileFunc
		if colonIdx := strings.LastIndex(fileFunc, ":"); colonIdx > 0 {
			pkgPath = fileFunc[:colonIdx]
		}
		if slashIdx := strings.LastIndex(pkgPath, "/"); slashIdx > 0 {
			pkgPath = pkgPath[:slashIdx]
		}

		if _, exists := pkgMap[pkgPath]; !exists {
			pkgMap[pkgPath] = &accum{}
			order = append(order, pkgPath)
		}
		// Approximate: treat each function as ~10 statements, weighted by pct
		// go tool cover -func doesn't give raw statement counts per function,
		// so we use 10 as unit weight per function entry.
		const weight = 10
		pkgMap[pkgPath].stmts += weight
		pkgMap[pkgPath].covered += int(float64(weight) * pct / 100)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading go coverage file: %w", err)
	}

	results := make([]PackageCoverage, 0, len(order))
	for _, pkg := range order {
		a := pkgMap[pkg]
		pct := 0.0
		if a.stmts > 0 {
			pct = float64(a.covered) / float64(a.stmts) * 100
		}
		// Use only the short package name (last segment)
		short := pkg
		if i := strings.LastIndex(pkg, "/"); i >= 0 {
			short = pkg[i+1:]
		}
		results = append(results, PackageCoverage{
			Package:    short,
			Statements: a.stmts,
			Covered:    a.covered,
			Pct:        pct,
		})
	}
	return results, nil
}
