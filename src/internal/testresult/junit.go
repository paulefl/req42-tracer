package testresult

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"github.com/paulefl/req42-tracer/src/internal/model"
)

// JUnitTestSuite represents a JUnit XML test suite.
type JUnitTestSuite struct {
	XMLName   xml.Name         `xml:"testsuite"`
	Name      string           `xml:"name,attr"`
	Package   string           `xml:"package,attr"`
	Tests     int              `xml:"tests,attr"`
	Skipped   int              `xml:"skipped,attr"`
	Failures  int              `xml:"failures,attr"`
	Errors    int              `xml:"errors,attr"`
	Time      float64          `xml:"time,attr"`
	Timestamp string           `xml:"timestamp,attr"`
	TestCases []JUnitTestCase  `xml:"testcase"`
}

// JUnitTestCase represents a single JUnit test case.
type JUnitTestCase struct {
	Name      string       `xml:"name,attr"`
	Classname string       `xml:"classname,attr"`
	Time      float64      `xml:"time,attr"`
	Skipped   *JUnitStatus `xml:"skipped"`
	Failure   *JUnitStatus `xml:"failure"`
	Error     *JUnitStatus `xml:"error"`
	SystemOut string       `xml:"system-out"`
	SystemErr string       `xml:"system-err"`
}

// JUnitStatus represents a skipped, failure, or error element.
type JUnitStatus struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

// ParseJUnit parses a JUnit XML file and returns TestResult objects.
func ParseJUnit(filePath, project, platform string) ([]*model.TestResult, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JUnit file %s: %w", filePath, err)
	}

	var suite JUnitTestSuite
	if err := xml.Unmarshal(data, &suite); err != nil {
		return nil, fmt.Errorf("failed to parse JUnit XML %s: %w", filePath, err)
	}

	var results []*model.TestResult
	timestamp := time.Now()
	if suite.Timestamp != "" {
		if t, err := time.Parse(time.RFC3339, suite.Timestamp); err == nil {
			timestamp = t
		}
	}

	for i, tc := range suite.TestCases {
		result := &model.TestResult{
			ID:         shortPackage(tc.Classname) + "::" + tc.Name,
			Package:    tc.Classname,
			TestName:   tc.Name,
			FullName:   tc.Classname + "::" + tc.Name,
			Duration:   tc.Time,
			Status:     "passed",
			Timestamp:  timestamp,
			Project:    project,
			Platform:   platform,
			Stdout:     tc.SystemOut,
			Stderr:     tc.SystemErr,
			Attributes: make(map[string]string),
		}

		// Determine status
		if tc.Skipped != nil {
			result.Status = "skipped"
			result.Error = tc.Skipped.Message
		} else if tc.Failure != nil {
			result.Status = "failed"
			result.Error = tc.Failure.Text
			if tc.Failure.Message != "" {
				result.Error = tc.Failure.Message + ": " + result.Error
			}
		} else if tc.Error != nil {
			result.Status = "failed"
			result.Error = tc.Error.Text
			if tc.Error.Message != "" {
				result.Error = tc.Error.Message + ": " + result.Error
			}
		}

		// Store the index for tracing
		result.Attributes["junit-index"] = fmt.Sprintf("%d", i)

		results = append(results, result)
	}

	return results, nil
}
