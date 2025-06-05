package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"vietnamese-converter/pkg/converter"
)

type Config struct {
	TestFile        string `json:"test_file"`
	OutputFile      string `json:"output_file"`
	MaxFailures     int    `json:"max_failures"`
	MaxErrors       int    `json:"max_errors"`
	PerformanceTest bool   `json:"performance_test"`
	Verbose         bool   `json:"verbose"`
}

type DetailedTestReport struct {
	Config          Config                        `json:"config"`
	Summary         TestSummary                   `json:"summary"`
	FailedCases     []converter.TestResult        `json:"failed_cases,omitempty"`
	ErrorCases      []converter.TestResult        `json:"error_cases,omitempty"`
	PerformanceData *PerformanceData              `json:"performance_data,omitempty"`
	Timestamp       time.Time                     `json:"timestamp"`
}

type TestSummary struct {
	TotalTests   int           `json:"total_tests"`
	PassedTests  int           `json:"passed_tests"`
	FailedTests  int           `json:"failed_tests"`
	ErrorTests   int           `json:"error_tests"`
	PassRate     float64       `json:"pass_rate"`
	TotalTime    time.Duration `json:"total_time"`
	AverageTime  time.Duration `json:"average_time"`
}

type PerformanceData struct {
	Iterations      int           `json:"iterations"`
	TotalTime       time.Duration `json:"total_time"`
	AverageTime     time.Duration `json:"average_time"`
	MinTime         time.Duration `json:"min_time"`
	MaxTime         time.Duration `json:"max_time"`
	ConversionsPerSecond float64  `json:"conversions_per_second"`
}

func main() {
	var (
		testFile        = flag.String("file", "random_numbers_with_vietnamese.txt", "Path to test data file")
		outputFile      = flag.String("output", "", "Path to save detailed JSON report (optional)")
		maxFailures     = flag.Int("max-failures", 20, "Maximum number of failed cases to display")
		maxErrors       = flag.Int("max-errors", 10, "Maximum number of error cases to display")
		performanceTest = flag.Bool("perf", false, "Run additional performance tests")
		verbose         = flag.Bool("verbose", false, "Verbose output")
		saveReport      = flag.Bool("save", false, "Save detailed report to JSON file")
		configFile      = flag.String("config", "", "Load configuration from JSON file")
	)
	flag.Parse()

	config := Config{
		TestFile:        *testFile,
		OutputFile:      *outputFile,
		MaxFailures:     *maxFailures,
		MaxErrors:       *maxErrors,
		PerformanceTest: *performanceTest,
		Verbose:         *verbose,
	}

	// Load config from file if specified
	if *configFile != "" {
		if err := loadConfig(*configFile, &config); err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
	}

	// Determine output file name if not specified but save is requested
	if *saveReport && config.OutputFile == "" {
		config.OutputFile = fmt.Sprintf("test_report_%s.json", time.Now().Format("20060102_150405"))
	}

	fmt.Printf("Vietnamese Number Converter Test Suite\n")
	fmt.Printf("=====================================\n")
	fmt.Printf("Test file: %s\n", config.TestFile)
	fmt.Printf("Max failures to show: %d\n", config.MaxFailures)
	fmt.Printf("Max errors to show: %d\n", config.MaxErrors)
	fmt.Printf("Performance test: %v\n", config.PerformanceTest)
	fmt.Printf("Verbose mode: %v\n", config.Verbose)
	if config.OutputFile != "" {
		fmt.Printf("Output file: %s\n", config.OutputFile)
	}
	fmt.Println()

	// Run the main test suite
	fmt.Println("Loading test cases and running conversion tests...")
	testSuite := converter.NewTestSuite()
	
	start := time.Now()
	results, err := testSuite.RunAllTests(config.TestFile)
	if err != nil {
		log.Fatalf("Failed to run tests: %v", err)
	}
	totalTestTime := time.Since(start)

	// Generate report
	report := testSuite.GenerateReport(results)
	
	// Print summary
	printDetailedSummary(report, totalTestTime, config.Verbose)
	
	// Print failed cases
	if len(report.FailedCases) > 0 {
		printFailedCases(report.FailedCases, config.MaxFailures, config.Verbose)
	}
	
	// Print error cases
	if len(report.ErrorCases) > 0 {
		printErrorCases(report.ErrorCases, config.MaxErrors, config.Verbose)
	}

	// Run performance tests if requested
	var perfData *PerformanceData
	if config.PerformanceTest {
		fmt.Println("\nRunning additional performance tests...")
		perfData = runPerformanceTests()
		printPerformanceResults(perfData)
	}

	// Save detailed report if requested
	if config.OutputFile != "" {
		detailedReport := DetailedTestReport{
			Config: config,
			Summary: TestSummary{
				TotalTests:   report.TotalTests,
				PassedTests:  report.PassedTests,
				FailedTests:  report.FailedTests,
				ErrorTests:   report.ErrorTests,
				PassRate:     float64(report.PassedTests) / float64(report.TotalTests) * 100,
				TotalTime:    report.TotalTime,
				AverageTime:  report.AverageTime,
			},
			FailedCases:     report.FailedCases,
			ErrorCases:      report.ErrorCases,
			PerformanceData: perfData,
			Timestamp:       time.Now(),
		}

		if err := saveDetailedReport(config.OutputFile, detailedReport); err != nil {
			log.Printf("Failed to save detailed report: %v", err)
		} else {
			fmt.Printf("\nDetailed report saved to: %s\n", config.OutputFile)
		}
	}

	// Exit with appropriate code
	passRate := float64(report.PassedTests) / float64(report.TotalTests) * 100
	fmt.Printf("\n=== Final Result ===\n")
	fmt.Printf("Pass Rate: %.2f%%\n", passRate)
	
	if passRate < 95.0 {
		fmt.Printf("❌ Test suite FAILED - Pass rate below 95%%\n")
		os.Exit(1)
	} else if len(report.ErrorCases) > 0 {
		fmt.Printf("⚠️  Test suite PASSED with warnings - %d errors encountered\n", len(report.ErrorCases))
		os.Exit(0)
	} else {
		fmt.Printf("✅ Test suite PASSED\n")
		os.Exit(0)
	}
}

func loadConfig(filename string, config *Config) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, config)
}

func printDetailedSummary(report converter.TestReport, totalTime time.Duration, verbose bool) {
	fmt.Printf("=== Detailed Test Results ===\n")
	fmt.Printf("Total Tests: %d\n", report.TotalTests)
	fmt.Printf("Passed: %d (%.2f%%)\n", report.PassedTests, float64(report.PassedTests)/float64(report.TotalTests)*100)
	fmt.Printf("Failed: %d (%.2f%%)\n", report.FailedTests, float64(report.FailedTests)/float64(report.TotalTests)*100)
	fmt.Printf("Errors: %d (%.2f%%)\n", report.ErrorTests, float64(report.ErrorTests)/float64(report.TotalTests)*100)
	fmt.Printf("Total Execution Time: %v\n", totalTime)
	fmt.Printf("Total Conversion Time: %v\n", report.TotalTime)
	fmt.Printf("Average Time per Conversion: %v\n", report.AverageTime)
	
	if verbose {
		fmt.Printf("Fastest Conversion: %v\n", findFastestTime(report))
		fmt.Printf("Slowest Conversion: %v\n", findSlowestTime(report))
		fmt.Printf("Conversions per Second: %.0f\n", float64(report.TotalTests)/report.TotalTime.Seconds())
	}
	fmt.Println()
}

func printFailedCases(failedCases []converter.TestResult, maxToShow int, verbose bool) {
	fmt.Printf("=== Failed Cases (showing first %d of %d) ===\n", min(maxToShow, len(failedCases)), len(failedCases))
	
	count := 0
	for _, result := range failedCases {
		if count >= maxToShow {
			break
		}
		
		fmt.Printf("Line %d: %d\n", result.TestCase.LineNumber, result.TestCase.Number)
		fmt.Printf("  Expected: %s\n", result.Expected)
		fmt.Printf("  Actual:   %s\n", result.ActualResult)
		
		if verbose {
			fmt.Printf("  Time:     %v\n", result.ProcessingTime)
			// Show difference analysis
			fmt.Printf("  Analysis: %s\n", analyzeFailure(result.Expected, result.ActualResult))
		}
		fmt.Println()
		count++
	}
}

func printErrorCases(errorCases []converter.TestResult, maxToShow int, verbose bool) {
	fmt.Printf("=== Error Cases (showing first %d of %d) ===\n", min(maxToShow, len(errorCases)), len(errorCases))
	
	count := 0
	for _, result := range errorCases {
		if count >= maxToShow {
			break
		}
		
		fmt.Printf("Line %d: %d\n", result.TestCase.LineNumber, result.TestCase.Number)
		fmt.Printf("  Error: %v\n", result.Error)
		
		if verbose {
			fmt.Printf("  Expected: %s\n", result.Expected)
		}
		fmt.Println()
		count++
	}
}

func runPerformanceTests() *PerformanceData {
	converter := converter.NewVietnameseConverter()
	
	// Test numbers of varying complexity
	testNumbers := []int64{
		1, 15, 101, 1001, 12345, 123456, 1234567, 12345678, 123456789, 1234567890,
	}
	
	iterations := 10000
	var times []time.Duration
	
	start := time.Now()
	for i := 0; i < iterations; i++ {
		num := testNumbers[i%len(testNumbers)]
		convStart := time.Now()
		_, err := converter.Convert(num)
		convTime := time.Since(convStart)
		
		if err != nil {
			continue // Skip errors for performance test
		}
		
		times = append(times, convTime)
	}
	totalTime := time.Since(start)
	
	// Calculate statistics
	var minTime, maxTime time.Duration
	var totalConvTime time.Duration
	
	if len(times) > 0 {
		minTime = times[0]
		maxTime = times[0]
		
		for _, t := range times {
			totalConvTime += t
			if t < minTime {
				minTime = t
			}
			if t > maxTime {
				maxTime = t
			}
		}
	}
	
	avgTime := totalConvTime / time.Duration(len(times))
	conversionsPerSecond := float64(len(times)) / totalTime.Seconds()
	
	return &PerformanceData{
		Iterations:           len(times),
		TotalTime:           totalTime,
		AverageTime:         avgTime,
		MinTime:             minTime,
		MaxTime:             maxTime,
		ConversionsPerSecond: conversionsPerSecond,
	}
}

func printPerformanceResults(perfData *PerformanceData) {
	fmt.Printf("=== Performance Test Results ===\n")
	fmt.Printf("Iterations: %d\n", perfData.Iterations)
	fmt.Printf("Total Time: %v\n", perfData.TotalTime)
	fmt.Printf("Average Time: %v\n", perfData.AverageTime)
	fmt.Printf("Min Time: %v\n", perfData.MinTime)
	fmt.Printf("Max Time: %v\n", perfData.MaxTime)
	fmt.Printf("Conversions/Second: %.0f\n", perfData.ConversionsPerSecond)
	fmt.Println()
}

func saveDetailedReport(filename string, report DetailedTestReport) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func findFastestTime(report converter.TestReport) time.Duration {
	if len(report.FailedCases) == 0 && len(report.ErrorCases) == 0 {
		return 0 // No individual times available in current report structure
	}
	return 0 // Placeholder - would need to track individual times
}

func findSlowestTime(report converter.TestReport) time.Duration {
	return 0 // Placeholder - would need to track individual times
}

func analyzeFailure(expected, actual string) string {
	if len(actual) == 0 {
		return "Empty result"
	}
	if len(expected) != len(actual) {
		return fmt.Sprintf("Length mismatch (expected: %d, actual: %d)", len(expected), len(actual))
	}
	// Could add more sophisticated analysis here
	return "Content differs"
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
} 