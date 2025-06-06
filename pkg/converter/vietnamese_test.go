package converter_test

import (
	"os"
	"testing"

	"vietnamese-converter/pkg/converter"
	testutil "vietnamese-converter/pkg/converter/testutil"
)

























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
		t.Fatalf("Error running tests: %v", err)
	}

	report := testSuite.GenerateReport(results)
	report.PrintSummary()

	// Print detailed failure information if there are any failures
	if report.FailedTests > 0 || report.ErrorTests > 0 {
		t.Logf("\n=== Detailed Failure Report ===")
		if len(report.FailedCases) > 0 {
			t.Logf("\nFailed Cases (showing first 10):")
			for i, result := range report.FailedCases {
				if i >= 10 {
					t.Logf("  ... and %d more failures", len(report.FailedCases)-10)
					break
				}
				t.Logf("  %d: Expected %q, got %q", 
					result.TestCase.Number, 
					result.TestCase.ExpectedVietnamese,
					result.ActualResult)
			}
		}

		if len(report.ErrorCases) > 0 {
			t.Logf("\nError Cases (showing first 5):")
			for i, result := range report.ErrorCases {
				if i >= 5 {
					t.Logf("  ... and %d more errors", len(report.ErrorCases)-5)
					break
				}
				t.Logf("  %d: %v", result.TestCase.Number, result.Error)
			}
		}

		t.Logf("\nFull results:")
		t.Logf("  Total test cases: %d", report.TotalTests)
		t.Logf("  Pass rate: %.2f%%", float64(report.PassedTests)/float64(report.TotalTests)*100)
		t.Logf("  Average processing time: %v", report.AverageTime)

		t.Fail()
	}
}

// Benchmark function for performance testing
func BenchmarkVietnameseConverter_Convert(b *testing.B) {
	conv := converter.NewConverter()
	
	// Reset the timer to exclude setup time
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Use a fixed number for consistent benchmarking
		_, _ = conv.Convert(1234567890)
	}
}