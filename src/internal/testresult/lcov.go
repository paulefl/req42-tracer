package testresult

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ParseLCOV parses an LCOV .info coverage file.
// Returns one PackageCoverage per source directory (parent of SF: paths).
func ParseLCOV(filePath string) ([]PackageCoverage, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open lcov file %s: %w", filePath, err)
	}
	defer f.Close()

	type accum struct{ stmts, covered int }
	pkgMap := make(map[string]*accum)
	var order []string

	var currentPkg string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.HasPrefix(line, "SF:"):
			// SF: source file path → use parent directory as package
			src := strings.TrimPrefix(line, "SF:")
			currentPkg = filepath.Base(filepath.Dir(src))
			if currentPkg == "." || currentPkg == "" {
				currentPkg = filepath.Base(src)
			}
			if _, exists := pkgMap[currentPkg]; !exists {
				pkgMap[currentPkg] = &accum{}
				order = append(order, currentPkg)
			}

		case strings.HasPrefix(line, "DA:") && currentPkg != "":
			// DA:line_number,hit_count
			parts := strings.Split(strings.TrimPrefix(line, "DA:"), ",")
			if len(parts) < 2 {
				continue
			}
			hits, err := strconv.Atoi(parts[1])
			if err != nil {
				continue
			}
			pkgMap[currentPkg].stmts++
			if hits > 0 {
				pkgMap[currentPkg].covered++
			}

		case line == "end_of_record":
			currentPkg = ""
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading lcov file: %w", err)
	}

	results := make([]PackageCoverage, 0, len(order))
	for _, pkg := range order {
		a := pkgMap[pkg]
		pct := 0.0
		if a.stmts > 0 {
			pct = float64(a.covered) / float64(a.stmts) * 100
		}
		results = append(results, PackageCoverage{
			Package:    pkg,
			Statements: a.stmts,
			Covered:    a.covered,
			Pct:        pct,
		})
	}
	return results, nil
}
