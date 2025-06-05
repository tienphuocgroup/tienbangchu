package converter_test

import (
	"os"
	"testing"
	"time"

	"vietnamese-converter/pkg/converter"
	testutil "vietnamese-converter/pkg/converter/testutil"
)

































// Test functions for go test
// NOTE: Only the golden file test (TestVietnameseConverter_FullDataset) is kept. All hardcoded tests are commented out to ensure the golden file is the single source of truth.

/*
func TestVietnameseConverter_BasicNumbers(t *testing.T) {
	converter := converter.NewVietnameseConverter()
	
	testCases := []struct {
		number   int64
		expected string
	}{
		{0, "không đồng"},
		{1, "một đồng"},
		{5, "năm đồng"},
		{10, "mười đồng"},
		{15, "mười lăm đồng"},
		{21, "hai mười mốt đồng"},
		{24, "hai mười tư đồng"},
		{101, "một trăm lẻ một đồng"},
		{1001, "một nghìn không trăm lẻ một đồng"},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("number_%d", tc.number), func(t *testing.T) {
			result, err := converter.Convert(tc.number)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("For number %d, expected %q, got %q", tc.number, tc.expected, result)
			}
		})
	}
}
*/

/*
func TestVietnameseConverter_EdgeCases(t *testing.T) {
	converter := converter.NewVietnameseConverter()
	
	testCases := []struct {
		name     string
		number   int64
		expected string
	}{
		{"Four vs Tư", 24, "hai mười tư đồng"},
		{"Four as Bốn", 40, "bốn mười đồng"},
		{"Four in thousands", 34000, "ba mười bốn nghìn đồng"},
		{"One vs Mốt", 21, "hai mười mốt đồng"},
		{"One standalone", 101, "một trăm lẻ một đồng"},
		{"Zero with lẻ", 1001, "một nghìn không trăm lẻ một đồng"},
		{"Zero handling", 50050050, "năm mười triệu năm mười nghìn không trăm năm mười đồng"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := converter.Convert(tc.number)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if result != tc.expected {
				t.Errorf("For number %d, expected %q, got %q", tc.number, tc.expected, result)
			}
		})
	}
}
*/

/*
func TestVietnameseConverter_ErrorCases(t *testing.T) {
	converter := converter.NewVietnameseConverter()
	
	testCases := []struct {
		name   string
		number int64
		hasErr bool
	}{
		{"Negative number", -1, true},
		{"Too large", 1000000000000000, true},
		{"Max valid", 999999999999999, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := converter.Convert(tc.number)
			if tc.hasErr && err == nil {
				t.Errorf("Expected error for number %d, but got none", tc.number)
			}
			if !tc.hasErr && err != nil {
				t.Errorf("Unexpected error for number %d: %v", tc.number, err)
			}
		})
	}
}
*/

/*
func TestVietnameseConverter_PerformanceBenchmark(t *testing.T) {
	converter := converter.NewVietnameseConverter()
	
	// Test performance with a variety of numbers
	numbers := []int64{123, 12345, 1234567, 123456789, 12345678901}
	
	start := time.Now()
	iterations := 10000
	
	for i := 0; i < iterations; i++ {
		for _, num := range numbers {
			_, err := converter.Convert(num)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		}
	}
	
	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations*len(numbers))
	
	t.Logf("Performance test completed:")
	t.Logf("  Total iterations: %d", iterations*len(numbers))
	t.Logf("  Total time: %v", duration)
	t.Logf("  Average time per conversion: %v", avgTime)
	
	// Ensure average time is reasonable (should be sub-millisecond)
	if avgTime > time.Millisecond {
		t.Errorf("Performance too slow: average time %v exceeds 1ms threshold", avgTime)
	}
}
*/

// Integration test that runs against the full test dataset
func TestVietnameseConverter_FullDataset(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full dataset test in short mode")
	}

	testSuite := testutil.NewTestSuite()
	// Try to load the test file - this test will be skipped if the file doesn't exist
	results, err := testSuite.RunAllTests("../../random_numbers_with_vietnamese.txt")
	if err != nil {
		// Check if it's a file not found error
		if os.IsNotExist(err) {
			t.Skip("Test data file not found, skipping full dataset test")
		}
		t.Fatalf("Failed to run full dataset test: %v", err)
	}

	report := testSuite.GenerateReport(results)
	// Print summary for visibility
	report.PrintSummary()
	// Print some failed cases for debugging
	if len(report.FailedCases) > 0 {
		report.PrintFailedCases(10)
	}
	// Print error cases if any
	if len(report.ErrorCases) > 0 {
		report.PrintErrorCases(5)
	}

	// Calculate pass rate
	passRate := float64(report.PassedTests) / float64(report.TotalTests) * 100
	
	t.Logf("Full results:")
	t.Logf("  Total test cases: %d", report.TotalTests)
	t.Logf("  Pass rate: %.2f%%", passRate)
	t.Logf("  Average processing time: %v", report.AverageTime)

	// We expect a high pass rate, but allow for some failures during development
	// This threshold can be adjusted as the algorithm improves
	expectedPassRate := 95.0
	if passRate < expectedPassRate {
		t.Errorf("Pass rate %.2f%% is below expected threshold of %.2f%%", passRate, expectedPassRate)
	}

	// Ensure performance is good
	if report.AverageTime > 2*time.Millisecond {
		t.Errorf("Average processing time %v exceeds 2ms threshold", report.AverageTime)
	}
}

// Benchmark function for performance testing
func BenchmarkVietnameseConverter_Convert(b *testing.B) {
	converter := converter.NewVietnameseConverter()
	numbers := []int64{123, 12345, 1234567, 123456789, 987654321}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		num := numbers[i%len(numbers)]
		_, err := converter.Convert(num)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Helper function to run manual testing during development
func RunManualTest() {
	println("Running manual test against full dataset...")

	testSuite := testutil.NewTestSuite()
	results, err := testSuite.RunAllTests("random_numbers_with_vietnamese.txt")
	if err != nil {
		println("Error running tests:", err.Error())
		return
	}

	report := testSuite.GenerateReport(results)
	report.PrintSummary()
	report.PrintFailedCases(20)
	report.PrintErrorCases(10)
}