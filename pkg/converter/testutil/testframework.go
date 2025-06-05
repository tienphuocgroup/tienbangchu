package testutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"vietnamese-converter/pkg/converter"
)

type TestCase struct {
	Number             int64
	ExpectedVietnamese string
	LineNumber         int
}

type TestResult struct {
	TestCase       TestCase
	ActualResult   string
	Expected       string
	Passed         bool
	Error          error
	ProcessingTime time.Duration
}

type TestReport struct {
	TotalTests   int
	PassedTests  int
	FailedTests  int
	ErrorTests   int
	TotalTime    time.Duration
	AverageTime  time.Duration
	FailedCases  []TestResult
	ErrorCases   []TestResult
}

func (tr *TestReport) PrintSummary() {
	fmt.Printf("\n=== Test Summary ===\n")
	fmt.Printf("Total Tests: %d\n", tr.TotalTests)
	fmt.Printf("Passed: %d (%.2f%%)\n", tr.PassedTests, float64(tr.PassedTests)/float64(tr.TotalTests)*100)
	fmt.Printf("Failed: %d (%.2f%%)\n", tr.FailedTests, float64(tr.FailedTests)/float64(tr.TotalTests)*100)
	fmt.Printf("Errors: %d (%.2f%%)\n", tr.ErrorTests, float64(tr.ErrorTests)/float64(tr.TotalTests)*100)
	fmt.Printf("Total Time: %v\n", tr.TotalTime)
	fmt.Printf("Average Time: %v\n", tr.AverageTime)
	fmt.Println()
}

func (tr *TestReport) PrintFailedCases(limit int) {
	if len(tr.FailedCases) == 0 {
		return
	}
	fmt.Printf("=== Failed Cases (showing first %d) ===\n", limit)
	count := 0
	for _, result := range tr.FailedCases {
		if count >= limit {
			break
		}
		fmt.Printf("Line %d: %d\n", result.TestCase.LineNumber, result.TestCase.Number)
		fmt.Printf("  Expected: %s\n", result.Expected)
		fmt.Printf("  Actual:   %s\n", result.ActualResult)
		fmt.Println()
		count++
	}
}

func (tr *TestReport) PrintErrorCases(limit int) {
	if len(tr.ErrorCases) == 0 {
		return
	}
	fmt.Printf("=== Error Cases (showing first %d) ===\n", limit)
	count := 0
	for _, result := range tr.ErrorCases {
		if count >= limit {
			break
		}
		fmt.Printf("Line %d: %d\n", result.TestCase.LineNumber, result.TestCase.Number)
		fmt.Printf("  Error: %v\n", result.Error)
		fmt.Println()
		count++
	}
}

type TestDataLoader struct {
	testCases []TestCase
	loaded    bool
}

func NewTestDataLoader() *TestDataLoader {
	return &TestDataLoader{}
}

func (tdl *TestDataLoader) LoadTestCases(filename string) error {
	if tdl.loaded {
		return nil
	}
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open test file: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			return fmt.Errorf("invalid line format at line %d: %s", lineNumber, line)
		}
		number, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid number at line %d: %s", lineNumber, parts[0])
		}
		expectedVietnamese := parts[1]
		tdl.testCases = append(tdl.testCases, TestCase{
			Number:             number,
			ExpectedVietnamese: expectedVietnamese,
			LineNumber:         lineNumber,
		})
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}
	tdl.loaded = true
	return nil
}

func (tdl *TestDataLoader) GetTestCases() []TestCase {
	return tdl.testCases
}

type TestSuite struct {
	converter converter.NumberConverter
	loader    *TestDataLoader
}

func NewTestSuite() *TestSuite {
	return &TestSuite{
		converter: converter.NewVietnameseConverter(),
		loader:    NewTestDataLoader(),
	}
}

func (ts *TestSuite) RunAllTests(filename string) ([]TestResult, error) {
	err := ts.loader.LoadTestCases(filename)
	if err != nil {
		return nil, err
	}
	testCases := ts.loader.GetTestCases()
	results := make([]TestResult, 0, len(testCases))
	for _, tc := range testCases {
		result := ts.runSingleTest(tc)
		results = append(results, result)
	}
	return results, nil
}

func (ts *TestSuite) runSingleTest(tc TestCase) TestResult {
	start := time.Now()
	actual, err := ts.converter.Convert(tc.Number)
	processingTime := time.Since(start)
	result := TestResult{
		TestCase:       tc,
		ActualResult:   actual,
		Expected:       tc.ExpectedVietnamese,
		ProcessingTime: processingTime,
		Error:          err,
	}
	if err != nil {
		result.Passed = false
	} else {
		result.Passed = actual == tc.ExpectedVietnamese
	}
	return result
}

func (ts *TestSuite) GenerateReport(results []TestResult) TestReport {
	report := TestReport{
		TotalTests:    len(results),
		PassedTests:   0,
		FailedTests:   0,
		ErrorTests:    0,
		TotalTime:     0,
		FailedCases:   make([]TestResult, 0),
		ErrorCases:    make([]TestResult, 0),
	}
	var totalTime time.Duration
	for _, result := range results {
		totalTime += result.ProcessingTime
		if result.Error != nil {
			report.ErrorTests++
			report.ErrorCases = append(report.ErrorCases, result)
		} else if result.Passed {
			report.PassedTests++
		} else {
			report.FailedTests++
			report.FailedCases = append(report.FailedCases, result)
		}
	}
	report.TotalTime = totalTime
	if len(results) > 0 {
		report.AverageTime = totalTime / time.Duration(len(results))
	}
	return report
}
