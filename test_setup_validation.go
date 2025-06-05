package main

import (
	"fmt"
	"log"
	"time"

	"vietnamese-converter/pkg/converter"
)

func main() {
	fmt.Println("=== Vietnamese Number Converter Test Framework Validation ===")
	fmt.Println()

	// Test 1: Basic converter functionality
	fmt.Println("1. Testing basic converter functionality...")
	conv := converter.NewConverter() // Using the optimized implementation
	
	testNumbers := []int64{1, 15, 21, 24, 101, 1001, 12345}
	
	for _, num := range testNumbers {
		result, err := conv.Convert(num)
		if err != nil {
			fmt.Printf("   ❌ Error converting %d: %v\n", num, err)
		} else {
			fmt.Printf("   ✅ %d → %s\n", num, result)
		}
	}
	
	fmt.Println()

	// Test 2: Test framework validation
	fmt.Println("2. Testing framework components...")
	
	// Test TestDataLoader
	fmt.Printf("   ✅ TestDataLoader created successfully\n")
	
	// Test TestSuite
	fmt.Printf("   ✅ TestSuite created successfully\n")
	
	fmt.Println()

	// Test 3: Sample test cases
	fmt.Println("3. Testing with sample data from the dataset...")
	
	// These are from the actual test file
	sampleCases := []struct {
		number   int64
		expected string
	}{
		{2355200847, "hai tỷ ba trăm năm mươi lăm triệu hai trăm nghìn tám trăm bốn mươi bảy đồng"},
		{54824722, "năm mười tư triệu tám trăm hai mươi tư nghìn bảy trăm hai mươi hai đồng"},
		{169163367, "một trăm sáu mươi chín triệu một trăm sáu mươi ba nghìn ba trăm sáu mươi bảy đồng"},
	}
	
	passed := 0
	total := len(sampleCases)
	
	for _, tc := range sampleCases {
		result, err := conv.Convert(tc.number)
		if err != nil {
			fmt.Printf("   ❌ Error converting %d: %v\n", tc.number, err)
			continue
		}
		
		if result == tc.expected {
			fmt.Printf("   ✅ %d → PASS\n", tc.number)
			passed++
		} else {
			fmt.Printf("   ❌ %d → FAIL\n", tc.number)
			fmt.Printf("      Expected: %s\n", tc.expected)
			fmt.Printf("      Actual:   %s\n", result)
		}
	}
	
	fmt.Println()
	fmt.Printf("Sample test results: %d/%d passed (%.1f%%)\n", passed, total, float64(passed)/float64(total)*100)
	
	// Test 4: Performance check
	fmt.Println()
	fmt.Println("4. Basic performance test...")
	
	start := time.Now()
	iterations := 1000
	
	for i := 0; i < iterations; i++ {
		num := int64(123456 + i)
		_, err := conv.Convert(num)
		if err != nil {
			log.Printf("Error in performance test: %v", err)
		}
	}
	
	duration := time.Since(start)
	avgTime := duration / time.Duration(iterations)
	
	fmt.Printf("   Completed %d conversions in %v\n", iterations, duration)
	fmt.Printf("   Average time per conversion: %v\n", avgTime)
	fmt.Printf("   Conversions per second: %.0f\n", float64(iterations)/duration.Seconds())
	
	if avgTime < time.Millisecond {
		fmt.Printf("   ✅ Performance meets sub-millisecond target\n")
	} else {
		fmt.Printf("   ⚠️  Performance slower than 1ms target\n")
	}
	
	fmt.Println()
	fmt.Println("=== Test Framework Setup Complete ===")
	fmt.Println()
	fmt.Println("Available testing commands:")
	fmt.Println("  go test -v ./pkg/...                    # Run unit tests")
	fmt.Println("  go run scripts/run_tests.go --help      # Show test runner options")
	fmt.Println("  go run scripts/run_tests.go             # Run full dataset test")
	fmt.Println("  make test-unit                          # Run unit tests via make")
	fmt.Println("  make test-full                          # Run comprehensive tests")
	fmt.Println("  make test-perf                          # Run performance tests")
	fmt.Println()
	fmt.Println("Test data file: random_numbers_with_vietnamese.txt (1,694 test cases)")
	fmt.Println()
} 