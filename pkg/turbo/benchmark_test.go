package turbo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// BenchmarkZeroAllocConverter tests the ultimate converter performance
func BenchmarkZeroAllocConverter(b *testing.B) {
	converter := NewZeroAllocConverter()
	
	// Test with various number ranges
	testCases := []struct {
		name   string
		number int64
	}{
		{"Small", 123},
		{"Medium", 123456},
		{"Large", 123456789},
		{"VeryLarge", 123456789012},
		{"Edge_999", 999},
		{"Edge_1000", 1000},
		{"Edge_Million", 1000000},
		{"Edge_Billion", 1000000000},
	}
	
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()
			
			for i := 0; i < b.N; i++ {
				result := converter.Convert(tc.number)
				if len(result) == 0 {
					b.Fatal("Empty result")
				}
			}
		})
	}
}

// BenchmarkZeroAllocConverterRandom tests with random numbers
func BenchmarkZeroAllocConverterRandom(b *testing.B) {
	converter := NewZeroAllocConverter()
	
	// Pre-generate random numbers to avoid random generation overhead
	numbers := make([]int64, 10000)
	for i := range numbers {
		numbers[i] = rand.Int63n(999999999999) // Up to 999 billion
	}
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		number := numbers[i%len(numbers)]
		result := converter.Convert(number)
		if len(result) == 0 {
			b.Fatal("Empty result")
		}
	}
}

// BenchmarkConcurrentLoad tests performance under concurrent load
func BenchmarkConcurrentLoad(b *testing.B) {
	converter := NewZeroAllocConverter()
	
	b.RunParallel(func(pb *testing.PB) {
		number := int64(123456789)
		for pb.Next() {
			result := converter.Convert(number)
			if len(result) == 0 {
				b.Fatal("Empty result")
			}
		}
	})
}

// BenchmarkMemoryFootprint measures memory usage
func BenchmarkMemoryFootprint(b *testing.B) {
	var converters []*ZeroAllocConverter
	
	b.ResetTimer()
	b.ReportAllocs()
	
	for i := 0; i < b.N; i++ {
		converter := NewZeroAllocConverter()
		converters = append(converters, converter)
	}
	
	// Keep reference to prevent GC
	_ = converters
}

// TestPerformanceTarget ensures we meet performance targets
func TestPerformanceTarget(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	// Target: Sub-100μs per conversion
	targetLatency := 100 * time.Microsecond
	
	testNumbers := []int64{
		123,
		123456,
		123456789,
		123456789012,
		999999999999,
	}
	
	for _, number := range testNumbers {
		start := time.Now()
		result := converter.Convert(number)
		elapsed := time.Since(start)
		
		if elapsed > targetLatency {
			t.Errorf("Conversion of %d took %v, target is %v", 
				number, elapsed, targetLatency)
		}
		
		if len(result) == 0 {
			t.Errorf("Empty result for number %d", number)
		}
		
		t.Logf("Number: %d, Result: %s, Time: %v", number, result, elapsed)
	}
}

// TestThroughputTarget ensures we can handle 1000+ RPS
func TestThroughputTarget(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	// Simulate 1000 RPS for 1 second
	targetRPS := 1000
	duration := 1 * time.Second
	
	start := time.Now()
	count := 0
	
	for time.Since(start) < duration {
		converter.Convert(123456789)
		count++
	}
	
	actualRPS := float64(count) / duration.Seconds()
	
	if actualRPS < float64(targetRPS) {
		t.Errorf("Achieved %0.0f RPS, target is %d RPS", actualRPS, targetRPS)
	} else {
		t.Logf("✓ Achieved %0.0f RPS (target: %d RPS)", actualRPS, targetRPS)
	}
}

// BenchmarkComparison compares different implementations
func BenchmarkComparison(b *testing.B) {
	// This would compare against the existing implementations
	zeroAlloc := NewZeroAllocConverter()
	number := int64(123456789)
	
	b.Run("ZeroAlloc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			zeroAlloc.Convert(number)
		}
	})
}

// TestCacheEffectiveness verifies pre-computation effectiveness
func TestCacheEffectiveness(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	// Test that all 3-digit combinations are cached
	for i := 0; i < 1000; i++ {
		start := time.Now()
		result := converter.Convert(int64(i))
		elapsed := time.Since(start)
		
		// Even cached results should be very fast
		if elapsed > 10*time.Microsecond {
			t.Errorf("Cached conversion of %d took %v, too slow", i, elapsed)
		}
		
		if i > 0 && len(result) == 0 {
			t.Errorf("Empty result for cached number %d", i)
		}
	}
	
	hitRatio := converter.GetCacheHitRatio()
	if hitRatio < 1.0 {
		t.Errorf("Cache hit ratio is %f, expected 1.0", hitRatio)
	}
}

// TestMemoryUsage verifies minimal memory footprint
func TestMemoryUsage(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	footprint := converter.GetMemoryFootprint()
	
	// Footprint should be reasonable (under 1MB for lookup tables)
	maxFootprint := 1024 * 1024 // 1MB
	
	if footprint > maxFootprint {
		t.Errorf("Memory footprint is %d bytes, max allowed is %d bytes", 
			footprint, maxFootprint)
	} else {
		t.Logf("✓ Memory footprint: %d bytes (%0.1f KB)", 
			footprint, float64(footprint)/1024)
	}
}

// BenchmarkLatencyDistribution measures latency distribution
func BenchmarkLatencyDistribution(b *testing.B) {
	converter := NewZeroAllocConverter()
	
	latencies := make([]time.Duration, b.N)
	numbers := make([]int64, 1000)
	
	// Pre-generate test numbers
	for i := range numbers {
		numbers[i] = rand.Int63n(999999999999)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		number := numbers[i%len(numbers)]
		start := time.Now()
		converter.Convert(number)
		latencies[i] = time.Since(start)
	}
	
	// Calculate percentiles
	if len(latencies) > 0 {
		// Simple percentile calculation
		p50 := latencies[len(latencies)/2]
		p95 := latencies[len(latencies)*95/100]
		p99 := latencies[len(latencies)*99/100]
		
		b.Logf("Latency P50: %v, P95: %v, P99: %v", p50, p95, p99)
	}
}

// ExampleZeroAllocConverter demonstrates usage
func ExampleZeroAllocConverter() {
	converter := NewZeroAllocConverter()
	
	result := converter.Convert(123456789)
	fmt.Println(result)
	
	// Output: một trăm hai mươi ba triệu bốn trăm năm mươi sáu nghìn bảy trăm tám mươi chín đồng
}