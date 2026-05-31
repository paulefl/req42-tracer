package testresult

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type coberturaXML struct {
	XMLName  xml.Name           `xml:"coverage"`
	Packages []coberturaPackage `xml:"packages>package"`
}

type coberturaPackage struct {
	Name      string  `xml:"name,attr"`
	LineRate  float64 `xml:"line-rate,attr"`
	Classes   []coberturaClass `xml:"classes>class"`
}

type coberturaClass struct {
	Name    string       `xml:"name,attr"`
	Lines   []coberturaLine `xml:"lines>line"`
}

type coberturaLine struct {
	Hits int `xml:"hits,attr"`
}

// ParseCobertura parses a Cobertura XML coverage file.
// Returns one PackageCoverage per package.
func ParseCobertura(filePath string) ([]PackageCoverage, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("open cobertura file %s: %w", filePath, err)
	}

	var cov coberturaXML
	if err := xml.Unmarshal(data, &cov); err != nil {
		return nil, fmt.Errorf("parse cobertura XML %s: %w", filePath, err)
	}

	results := make([]PackageCoverage, 0, len(cov.Packages))
	for _, pkg := range cov.Packages {
		stmts, covered := 0, 0
		for _, cls := range pkg.Classes {
			for _, line := range cls.Lines {
				stmts++
				if line.Hits > 0 {
					covered++
				}
			}
		}
		pct := 0.0
		if stmts > 0 {
			pct = float64(covered) / float64(stmts) * 100
		} else {
			pct = pkg.LineRate * 100
		}
		name := pkg.Name
		if i := strings.LastIndex(name, "."); i >= 0 && !strings.Contains(name, "/") {
			name = name[i+1:]
		} else if i := strings.LastIndex(name, "/"); i >= 0 {
			name = name[i+1:]
		}
		results = append(results, PackageCoverage{
			Package:    name,
			Statements: stmts,
			Covered:    covered,
			Pct:        pct,
		})
	}
	return results, nil
}
