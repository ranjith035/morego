package reporter

import (
	"encoding/xml"
	"time"
)

// JUnitProperty represents variables logged in JUnit XML.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitProperties holds configurations list.
type JUnitProperties struct {
	Property []JUnitProperty `xml:"property"`
}

// JUnitFailure structures failed cases messages inside JUnit.
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

// JUnitTestCase maps individual runs metrics in JUnit.
type JUnitTestCase struct {
	Classname string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      float64       `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

// JUnitTestSuite structures execution suites metadata tags.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Time       float64         `xml:"time,attr"`
	Timestamp  string          `xml:"timestamp,attr"`
	Properties JUnitProperties `xml:"properties"`
	TestCases  []JUnitTestCase `xml:"testcase"`
}

// JUnitTestSuites aggregates suites collections.
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	Time       float64          `xml:"time,attr"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// ExportToJUnit converts SuiteResult logs into standard JUnit XML formats.
func ExportToJUnit(result *SuiteResult) ([]byte, error) {
	suite := JUnitTestSuite{
		Name:      result.Name,
		Tests:     result.TotalCount,
		Failures:  result.FailCount,
		Errors:    0,
		Time:      result.Duration.Seconds(),
		Timestamp: result.StartTime.Format(time.RFC3339),
	}

	for _, t := range result.Tests {
		tc := JUnitTestCase{
			Classname: result.Name,
			Name:      t.Name,
			Time:      t.Duration.Seconds(),
		}

		if t.Status == StatusFailed {
			tc.Failure = &JUnitFailure{
				Message: t.ErrorMessage,
				Type:    "AssertionError",
				Content: t.StackTrace,
			}
		}

		suite.TestCases = append(suite.TestCases, tc)
	}

	suites := JUnitTestSuites{
		Time:       result.Duration.Seconds(),
		TestSuites: []JUnitTestSuite{suite},
	}

	output, err := xml.MarshalIndent(suites, "", "  ")
	if err != nil {
		return nil, err
	}

	return []byte(xml.Header + string(output)), nil
}
