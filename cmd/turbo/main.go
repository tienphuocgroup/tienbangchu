package main

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"vietnamese-converter/pkg/turbo"
)

func main() {
	// Set GOMAXPROCS to number of CPU cores for optimal performance
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	// Disable garbage collector for maximum performance in production
	// This is safe since we use zero-allocation pools
	if os.Getenv("DISABLE_GC") == "true" {
		runtime.GC()
		debug.SetGCPercent(-1)
	}
	
	port := 8080
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}
	
	// Create the perfect service
	service := turbo.NewPerfectService()
	
	log.Printf("🚀 Perfect Vietnamese Service starting on port %d", port)
	log.Printf("💡 Target: 1000+ RPS with sub-100μs latency")
	
	if err := service.ListenAndServe(port); err != nil {
		log.Fatal("Service failed:", err)
	}
}