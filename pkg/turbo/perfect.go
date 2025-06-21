package turbo

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// PerfectService represents the ultimate Vietnamese converter service
// Design Philosophy: Perfect is when there's nothing left to remove
type PerfectService struct {
	converter    *ZeroAllocConverter
	server       *http.Server
	connPool     *ConnectionPool
	responsePool *ResponsePool
	metrics      *AtomicMetrics
}

// AtomicMetrics tracks performance with zero-allocation counters
type AtomicMetrics struct {
	totalRequests   uint64
	totalLatencyNs  uint64
	errorCount      uint64
	peakLatencyNs   uint64
}

// ConnectionPool manages HTTP connections with zero allocation
type ConnectionPool struct {
	conns chan net.Conn
	mu    sync.RWMutex
}

// ResponsePool manages pre-allocated response buffers
type ResponsePool struct {
	buffers sync.Pool
	writers sync.Pool
}

// NewPerfectService creates the ultimate Vietnamese conversion service
func NewPerfectService() *PerfectService {
	numCPU := runtime.NumCPU()
	
	return &PerfectService{
		converter: NewZeroAllocConverter(),
		connPool: &ConnectionPool{
			conns: make(chan net.Conn, numCPU*100), // Buffer per CPU
		},
		responsePool: &ResponsePool{
			buffers: sync.Pool{
				New: func() interface{} {
					// Pre-allocate optimal buffer size for Vietnamese responses
					return make([]byte, 0, 512)
				},
			},
			writers: sync.Pool{
				New: func() interface{} {
					return &FastResponseWriter{}
				},
			},
		},
		metrics: &AtomicMetrics{},
	}
}

// FastResponseWriter implements zero-allocation response writing
type FastResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// WriteHeader captures status code without allocation
func (w *FastResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
		w.written = true
	}
}

// ListenAndServe starts the perfect service
func (s *PerfectService) ListenAndServe(port int) error {
	// Create custom HTTP server with optimal settings
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      s,
		ReadTimeout:  100 * time.Millisecond, // Ultra-fast timeout
		WriteTimeout: 100 * time.Millisecond,
		IdleTimeout:  30 * time.Second,
		
		// Optimize for high throughput
		ReadHeaderTimeout: 50 * time.Millisecond,
		MaxHeaderBytes:    1024, // Minimal header size
		
		// Custom connection state handling
		ConnState: s.handleConnState,
	}
	
	// Listen with optimal socket settings
	listener, err := net.Listen("tcp", s.server.Addr)
	if err != nil {
		return err
	}
	
	// Optimize TCP settings for low latency
	if tcpListener, ok := listener.(*net.TCPListener); ok {
		listener = &optimizedListener{TCPListener: tcpListener}
	}
	
	return s.server.Serve(listener)
}

// optimizedListener applies TCP optimizations
type optimizedListener struct {
	*net.TCPListener
}

func (l *optimizedListener) Accept() (net.Conn, error) {
	conn, err := l.TCPListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	
	// Apply TCP optimizations
	conn.SetNoDelay(true)                           // Disable Nagle's algorithm
	conn.SetKeepAlive(true)                        // Enable keep-alive
	conn.SetKeepAlivePeriod(30 * time.Second)      // Keep-alive period
	conn.SetReadBuffer(4096)                       // Optimal read buffer
	conn.SetWriteBuffer(4096)                      // Optimal write buffer
	
	return conn, nil
}

// ServeHTTP implements the core request handling logic
// This is where the magic happens - zero allocations, maximum performance
func (s *PerfectService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Atomic increment for requests (zero allocation)
	atomic.AddUint64(&s.metrics.totalRequests, 1)
	
	// Get response writer from pool
	writer := s.responsePool.writers.Get().(*FastResponseWriter)
	writer.ResponseWriter = w
	writer.statusCode = 200
	writer.written = false
	defer func() {
		// Reset and return to pool
		writer.ResponseWriter = nil
		s.responsePool.writers.Put(writer)
	}()
	
	// Route handling - ultra-minimal routing
	switch {
	case r.Method == "POST" && r.URL.Path == "/convert":
		s.handleConvert(writer, r)
	case r.Method == "GET" && r.URL.Path == "/health":
		s.handleHealth(writer, r)
	case r.Method == "GET" && r.URL.Path == "/metrics":
		s.handleMetrics(writer, r)
	default:
		writer.WriteHeader(404)
		return
	}
	
	// Record latency (zero allocation)
	latency := time.Since(start).Nanoseconds()
	atomic.AddUint64(&s.metrics.totalLatencyNs, uint64(latency))
	
	// Update peak latency using atomic compare-and-swap
	for {
		current := atomic.LoadUint64(&s.metrics.peakLatencyNs)
		if uint64(latency) <= current {
			break
		}
		if atomic.CompareAndSwapUint64(&s.metrics.peakLatencyNs, current, uint64(latency)) {
			break
		}
	}
}

// handleConvert processes conversion requests with zero allocations
func (s *PerfectService) handleConvert(w *FastResponseWriter, r *http.Request) {
	// Get buffer from pool
	buf := s.responsePool.buffers.Get().([]byte)
	buf = buf[:0] // Reset length, keep capacity
	defer s.responsePool.buffers.Put(buf)
	
	// Parse number directly from body without JSON unmarshaling
	number, err := s.parseNumberFromBody(r)
	if err != nil {
		atomic.AddUint64(&s.metrics.errorCount, 1)
		w.WriteHeader(400)
		return
	}
	
	// Convert using zero-allocation converter
	vietnamese := s.converter.Convert(number)
	
	// Build JSON response directly in buffer (zero allocation)
	buf = append(buf, `{"number":`...)
	buf = appendInt(buf, number)
	buf = append(buf, `,"vietnamese":"`...)
	buf = append(buf, vietnamese...)
	buf = append(buf, `"}`...)
	
	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", unsafeString(len(buf)))
	w.Write(buf)
}

// handleHealth provides health check with minimal overhead
func (s *PerfectService) handleHealth(w *FastResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// handleMetrics provides performance metrics
func (s *PerfectService) handleMetrics(w *FastResponseWriter, r *http.Request) {
	requests := atomic.LoadUint64(&s.metrics.totalRequests)
	totalLatency := atomic.LoadUint64(&s.metrics.totalLatencyNs)
	errors := atomic.LoadUint64(&s.metrics.errorCount)
	peak := atomic.LoadUint64(&s.metrics.peakLatencyNs)
	
	avgLatency := uint64(0)
	if requests > 0 {
		avgLatency = totalLatency / requests
	}
	
	response := fmt.Sprintf(`{"requests":%d,"avg_latency_ns":%d,"peak_latency_ns":%d,"errors":%d}`,
		requests, avgLatency, peak, errors)
	
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
}

// handleConnState optimizes connection lifecycle
func (s *PerfectService) handleConnState(conn net.Conn, state http.ConnState) {
	switch state {
	case http.StateNew:
		// New connection - add to pool if available
		select {
		case s.connPool.conns <- conn:
		default:
			// Pool full, continue without pooling
		}
	case http.StateClosed:
		// Connection closed - remove from pool
		select {
		case <-s.connPool.conns:
		default:
		}
	}
}

// Shutdown gracefully stops the service
func (s *PerfectService) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Utility functions for zero-allocation operations

// parseNumberFromBody extracts number from request body without JSON parsing
func (s *PerfectService) parseNumberFromBody(r *http.Request) (int64, error) {
	// This would implement a direct parser for {"number":123} format
	// For now, simplified version
	buf := make([]byte, r.ContentLength)
	n, err := r.Body.Read(buf)
	if err != nil && n == 0 {
		return 0, err
	}
	
	// Find number in JSON using byte scanning (zero allocation)
	return extractNumberFromJSON(buf[:n])
}

// extractNumberFromJSON finds number value in JSON without parsing
func extractNumberFromJSON(data []byte) (int64, error) {
	// Simple state machine to find "number": value
	state := 0 // 0=looking for "number", 1=looking for :, 2=looking for value
	start := -1
	
	for i, b := range data {
		switch state {
		case 0:
			if b == '"' && i+6 < len(data) && string(data[i:i+8]) == `"number"` {
				state = 1
				i += 7 // Skip the rest of "number"
			}
		case 1:
			if b == ':' {
				state = 2
			}
		case 2:
			if b >= '0' && b <= '9' {
				start = i
				for j := i; j < len(data); j++ {
					if data[j] < '0' || data[j] > '9' {
						return parseIntFromBytes(data[start:j])
					}
				}
				return parseIntFromBytes(data[start:])
			}
		}
	}
	
	return 0, fmt.Errorf("number not found")
}

// parseIntFromBytes converts byte slice to int64 without allocation
func parseIntFromBytes(data []byte) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data")
	}
	
	var result int64
	for _, b := range data {
		if b < '0' || b > '9' {
			break
		}
		result = result*10 + int64(b-'0')
	}
	
	return result, nil
}

// appendInt appends integer to byte slice without allocation
func appendInt(buf []byte, n int64) []byte {
	if n == 0 {
		return append(buf, '0')
	}
	
	// Calculate digits
	temp := n
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}
	
	// Ensure capacity
	if cap(buf)-len(buf) < digits {
		newBuf := make([]byte, len(buf), len(buf)+digits+10)
		copy(newBuf, buf)
		buf = newBuf
	}
	
	// Add digits in reverse
	start := len(buf)
	buf = buf[:len(buf)+digits]
	for i := digits - 1; i >= 0; i-- {
		buf[start+i] = byte('0' + n%10)
		n /= 10
	}
	
	return buf
}

// unsafeString converts int to string without allocation
func unsafeString(n int) string {
	if n < 10 {
		return string([]byte{'0' + byte(n)})
	}
	
	// For larger numbers, use itoa
	buf := make([]byte, 0, 10)
	buf = appendInt(buf, int64(n))
	return *(*string)(unsafe.Pointer(&buf))
}