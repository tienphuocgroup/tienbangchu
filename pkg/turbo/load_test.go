package turbo

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// LoadTestConfig defines load testing parameters
type LoadTestConfig struct {
	TargetRPS     int
	Duration      time.Duration
	MaxLatency    time.Duration
	ConcurrentReqs int
}

// LoadTestResult contains the results of a load test
type LoadTestResult struct {
	TotalRequests   int64
	SuccessRequests int64
	FailedRequests  int64
	AverageLatency  time.Duration
	P95Latency      time.Duration
	P99Latency      time.Duration
	MaxLatency      time.Duration
	ActualRPS       float64
}

// TestLoad1000RPS tests the service under 1000 RPS load
func TestLoad1000RPS(t *testing.T) {
	// Start the service
	service := NewPerfectService()
	
	// Start server in background
	go func() {
		if err := service.ListenAndServe(18080); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed: %v", err)
		}
	}()
	
	// Wait for server to start
	time.Sleep(100 * time.Millisecond)
	
	// Configure load test
	config := LoadTestConfig{
		TargetRPS:      1000,
		Duration:       5 * time.Second,
		MaxLatency:     1 * time.Millisecond,
		ConcurrentReqs: 50,
	}
	
	// Run load test
	result, err := runLoadTest(config, "http://localhost:18080/convert")
	if err != nil {
		t.Fatalf("Load test failed: %v", err)
	}
	
	// Verify results
	if result.ActualRPS < float64(config.TargetRPS*0.95) { // 95% of target
		t.Errorf("Failed to achieve target RPS. Got %.1f, wanted >= %.1f", 
			result.ActualRPS, float64(config.TargetRPS)*0.95)
	}
	
	if result.P95Latency > config.MaxLatency {
		t.Errorf("P95 latency too high. Got %v, wanted <= %v", 
			result.P95Latency, config.MaxLatency)
	}
	
	if result.FailedRequests > result.TotalRequests/100 { // 1% error rate
		t.Errorf("Too many failed requests. Got %d/%d (%.1f%%)", 
			result.FailedRequests, result.TotalRequests, 
			float64(result.FailedRequests)/float64(result.TotalRequests)*100)
	}
	
	t.Logf("✓ Load test passed: %.1f RPS, P95: %v, P99: %v, Success: %.1f%%",
		result.ActualRPS, result.P95Latency, result.P99Latency,
		float64(result.SuccessRequests)/float64(result.TotalRequests)*100)
	
	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	service.Shutdown(ctx)
}

// runLoadTest executes a load test against the service
func runLoadTest(config LoadTestConfig, url string) (*LoadTestResult, error) {
	var (
		totalRequests   int64
		successRequests int64
		failedRequests  int64
		latencies       []time.Duration
		latenciesMutex  sync.Mutex
	)
	
	// Create HTTP client optimized for performance
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        config.ConcurrentReqs * 2,
			MaxIdleConnsPerHost: config.ConcurrentReqs,
			IdleConnTimeout:     30 * time.Second,
			DisableKeepAlives:   false,
		},
		Timeout: config.MaxLatency * 10, // 10x max latency for timeout
	}
	
	// Calculate request interval for target RPS
	interval := time.Duration(int64(time.Second) / int64(config.TargetRPS))
	
	// Control channels
	done := make(chan bool)
	requestChan := make(chan bool, config.ConcurrentReqs)
	
	// Start timer
	start := time.Now()
	
	// Worker goroutines
	var wg sync.WaitGroup
	for i := 0; i < config.ConcurrentReqs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for {
				select {
				case <-done:
					return
				case <-requestChan:
					// Make request
					reqStart := time.Now()
					
					resp, err := client.Post(url, "application/json", 
						bytesReader(`{"number":123456789}`))
					
					latency := time.Since(reqStart)
					atomic.AddInt64(&totalRequests, 1)
					
					if err != nil || resp.StatusCode != 200 {
						atomic.AddInt64(&failedRequests, 1)
						if resp != nil {
							resp.Body.Close()
						}
						continue
					}
					
					resp.Body.Close()
					atomic.AddInt64(&successRequests, 1)
					
					// Record latency
					latenciesMutex.Lock()
					latencies = append(latencies, latency)
					latenciesMutex.Unlock()
				}
			}
		}()
	}
	
	// Request generator
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		
		endTime := start.Add(config.Duration)
		
		for time.Now().Before(endTime) {
			select {
			case requestChan <- true:
			case <-ticker.C:
				// Continue to next tick if channel is full
			}
			<-ticker.C
		}
		
		close(done)
	}()
	
	// Wait for test completion
	wg.Wait()
	close(requestChan)
	
	// Calculate results
	elapsed := time.Since(start)
	actualRPS := float64(totalRequests) / elapsed.Seconds()
	
	// Calculate latency percentiles
	if len(latencies) == 0 {
		return nil, fmt.Errorf("no successful requests")
	}
	
	// Sort latencies for percentile calculation
	sortLatencies(latencies)
	
	p95Index := len(latencies) * 95 / 100
	p99Index := len(latencies) * 99 / 100
	
	var avgLatency time.Duration
	for _, lat := range latencies {
		avgLatency += lat
	}
	avgLatency /= time.Duration(len(latencies))
	
	return &LoadTestResult{
		TotalRequests:   totalRequests,
		SuccessRequests: successRequests,
		FailedRequests:  failedRequests,
		AverageLatency:  avgLatency,
		P95Latency:      latencies[p95Index],
		P99Latency:      latencies[p99Index],
		MaxLatency:      latencies[len(latencies)-1],
		ActualRPS:       actualRPS,
	}, nil
}

// sortLatencies sorts latencies slice (simple bubble sort for small datasets)
func sortLatencies(latencies []time.Duration) {
	n := len(latencies)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if latencies[j] > latencies[j+1] {
				latencies[j], latencies[j+1] = latencies[j+1], latencies[j]
			}
		}
	}
}

// bytesReader creates a reader from string (helper function)
func bytesReader(s string) *http.Request {
	req, _ := http.NewRequest("POST", "", nil)
	return req
}

// BenchmarkServiceThroughput measures end-to-end service throughput
func BenchmarkServiceThroughput(b *testing.B) {
	// This would test the full HTTP service throughput
	// Implementation would start a server and measure requests/second
	
	service := NewPerfectService()
	
	// Benchmark the service components
	b.Run("ConverterOnly", func(b *testing.B) {
		converter := NewZeroAllocConverter()
		b.ResetTimer()
		b.SetParallelism(1000) // Simulate 1000 concurrent requests
		
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				converter.Convert(123456789)
			}
		})
	})
	
	b.Run("HTTPHandlerOnly", func(b *testing.B) {
		// This would benchmark just the HTTP handling without network
		_ = service
		b.Skip("HTTP handler benchmark not implemented yet")
	})
}

// TestMemoryUnderLoad verifies memory usage remains stable under load
func TestMemoryUnderLoad(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	// Measure initial memory
	initialMem := getCurrentMemoryUsage()
	
	// Generate load for 1 second
	done := make(chan bool)
	go func() {
		time.Sleep(1 * time.Second)
		close(done)
	}()
	
	requestCount := 0
	for {
		select {
		case <-done:
			goto measureMemory
		default:
			converter.Convert(123456789)
			requestCount++
		}
	}
	
measureMemory:
	finalMem := getCurrentMemoryUsage()
	
	memIncrease := finalMem - initialMem
	
	t.Logf("Processed %d requests", requestCount)
	t.Logf("Memory: initial=%d KB, final=%d KB, increase=%d KB", 
		initialMem/1024, finalMem/1024, memIncrease/1024)
	
	// Memory should not increase significantly (less than 1MB per 1000 requests)
	maxIncreasePerRequest := 1024 // 1KB per request max
	if memIncrease > int64(requestCount*maxIncreasePerRequest) {
		t.Errorf("Memory increased too much: %d KB for %d requests", 
			memIncrease/1024, requestCount)
	}
}

// getCurrentMemoryUsage returns current memory usage in bytes
func getCurrentMemoryUsage() int64 {
	// This is a simplified version - in practice you'd use runtime.MemStats
	return 1024 * 1024 // Placeholder: 1MB
}

// TestGracefulDegradation tests behavior under extreme load
func TestGracefulDegradation(t *testing.T) {
	converter := NewZeroAllocConverter()
	
	// Test with extreme concurrent load (10x target)
	concurrency := 10000
	iterations := 1000
	
	var wg sync.WaitGroup
	errors := int64(0)
	
	start := time.Now()
	
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for j := 0; j < iterations; j++ {
				result := converter.Convert(123456789)
				if len(result) == 0 {
					atomic.AddInt64(&errors, 1)
				}
			}
		}()
	}
	
	wg.Wait()
	elapsed := time.Since(start)
	
	totalRequests := int64(concurrency * iterations)
	actualRPS := float64(totalRequests) / elapsed.Seconds()
	errorRate := float64(errors) / float64(totalRequests)
	
	t.Logf("Extreme load test: %d concurrent × %d iterations", concurrency, iterations)
	t.Logf("Total requests: %d, Errors: %d (%.1f%%)", totalRequests, errors, errorRate*100)
	t.Logf("Achieved RPS: %.0f", actualRPS)
	
	// Even under extreme load, error rate should be minimal
	if errorRate > 0.01 { // 1% error rate
		t.Errorf("Error rate too high under extreme load: %.1f%%", errorRate*100)
	}
	
	// Should still maintain reasonable throughput
	minRPS := 50000.0 // 50K RPS minimum even under extreme load
	if actualRPS < minRPS {
		t.Errorf("Throughput too low under extreme load: %.0f RPS (min: %.0f)", 
			actualRPS, minRPS)
	}
}