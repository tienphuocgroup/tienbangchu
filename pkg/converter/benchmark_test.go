package converter_test

import (
	"testing"

	"vietnamese-converter/pkg/converter"
)

// Test numbers from various ranges to get a comprehensive benchmark
var testNumbers = []int64{
	5, 12, 42, 101, 999,
	1000, 12345, 54824722, 123456789,
	1000000000, 2355200847, 9876543210,
}

// BenchmarkOriginalConverter measures the performance of the original implementation
func BenchmarkOriginalConverter(b *testing.B) {
	conv := converter.NewVietnameseConverter()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use modulo to cycle through test numbers
		num := testNumbers[i%len(testNumbers)]
		_, err := conv.Convert(num)
		if err != nil {
			b.Fatalf("Error converting %d: %v", num, err)
		}
	}
}

// BenchmarkOptimizedConverter measures the performance of our optimized implementation
func BenchmarkOptimizedConverter(b *testing.B) {
	conv := converter.NewTurboConverter()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use modulo to cycle through test numbers
		num := testNumbers[i%len(testNumbers)]
		_, err := conv.Convert(num)
		if err != nil {
			b.Fatalf("Error converting %d: %v", num, err)
		}
	}
}

// BenchmarkCompareConverters directly compares both implementations
func BenchmarkCompareConverters(b *testing.B) {
	original := converter.NewVietnameseConverter()
	optimized := converter.NewTurboConverter()
	
	b.Run("Original", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			num := testNumbers[i%len(testNumbers)]
			_, err := original.Convert(num)
			if err != nil {
				b.Fatalf("Error converting %d: %v", num, err)
			}
		}
	})
	
	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			num := testNumbers[i%len(testNumbers)]
			_, err := optimized.Convert(num)
			if err != nil {
				b.Fatalf("Error converting %d: %v", num, err)
			}
		}
	})
}

// BenchmarkParallelOptimized tests performance under concurrent usage
func BenchmarkParallelOptimized(b *testing.B) {
	conv := converter.NewTurboConverter()
	
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			num := testNumbers[i%len(testNumbers)]
			_, err := conv.Convert(num)
			if err != nil {
				b.Fatalf("Error converting %d: %v", num, err)
			}
			i++
		}
	})
}
