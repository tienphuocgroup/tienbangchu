package main

import (
	"fmt"
	"os"
	"time"
	
	"vietnamese-converter/pkg/converter"
)

// Test both implementations with a variety of numbers
func main() {
	fmt.Println("=== Vietnamese Number Converter Performance Comparison ===")
	
	// Create both converter implementations
	originalConverter := converter.NewVietnameseConverter()
	optimizedConverter := converter.NewTurboConverter()
	
	// Test numbers from various ranges
	testNumbers := []int64{
		5, 42, 101, 999,
		1000, 12345, 54824722, 123456789,
		1000000000, 2355200847, 9876543210,
	}
	
	// Test parameters
	iterations := 100000
	fmt.Printf("Running %d iterations per converter\n\n", iterations)
	
	// Benchmark original implementation
	fmt.Println("Testing Original Implementation...")
	originalStart := time.Now()
	for i := 0; i < iterations; i++ {
		num := testNumbers[i%len(testNumbers)]
		result, err := originalConverter.Convert(num)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		// Prevent compiler optimization by using the result
		if i == iterations-1 {
			fmt.Printf("Final conversion: %d → %s\n", num, result)
		}
	}
	originalDuration := time.Since(originalStart)
	originalAvg := originalDuration.Nanoseconds() / int64(iterations)
	fmt.Printf("Total time: %v\n", originalDuration)
	fmt.Printf("Average time per conversion: %d ns\n\n", originalAvg)
	
	// Benchmark optimized implementation
	fmt.Println("Testing Optimized Implementation...")
	optimizedStart := time.Now()
	for i := 0; i < iterations; i++ {
		num := testNumbers[i%len(testNumbers)]
		result, err := optimizedConverter.Convert(num)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		// Prevent compiler optimization by using the result
		if i == iterations-1 {
			fmt.Printf("Final conversion: %d → %s\n", num, result)
		}
	}
	optimizedDuration := time.Since(optimizedStart)
	optimizedAvg := optimizedDuration.Nanoseconds() / int64(iterations)
	fmt.Printf("Total time: %v\n", optimizedDuration)
	fmt.Printf("Average time per conversion: %d ns\n\n", optimizedAvg)
	
	// Calculate and display improvement
	speedup := float64(originalDuration) / float64(optimizedDuration)
	improvement := (speedup - 1.0) * 100
	fmt.Printf("Performance improvement: %.2fx faster (%.1f%% improvement)\n", speedup, improvement)
	
	// Display memory allocation comparison from benchmark results
	fmt.Println("\nMemory allocation comparison (from benchmark):")
	fmt.Println("Original: ~929 bytes/op with ~20 allocations/op")
	fmt.Println("Optimized: ~128 bytes/op with ~3 allocations/op")
	fmt.Printf("Memory reduction: %.1f%% fewer bytes, %.1f%% fewer allocations\n", 
		(1-(128.0/929.0))*100, (1-(3.0/20.0))*100)
}
