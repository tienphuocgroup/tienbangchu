# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a high-performance Vietnamese number converter service written in Go that provides both REST API endpoints and a web UI for converting numbers to Vietnamese text. The project emphasizes performance optimization with sub-25μs response times and handles Vietnamese linguistic exceptions properly.

## Development Commands

### Building and Running
```bash
# Start development server
make run

# Build production binary
make build

# Format code
make fmt
```

### Testing
```bash
# Run standard unit tests
make test-unit

# Run full test suite (unit + integration)
make test

# Run comprehensive test suite with full dataset
make test-full

# Run performance benchmarks
make test-perf

# Quick test with basic output
make test-quick

# Detailed test reporting
make test-runner
```

### Docker
```bash
# Build Docker image
make docker

# Run container
make docker-run
```

## Architecture

### Core Components

- **cmd/server/**: Main application entry point with graceful shutdown
- **internal/api/**: HTTP layer with Chi router
  - `handlers/`: Request handlers for convert and health endpoints
  - `middleware/`: Rate limiting, logging, recovery, request ID
  - `routes/`: Route definitions and setup
- **internal/config/**: Configuration management with environment variables
- **pkg/converter/**: Core conversion logic
  - `vietnamese.go`: Original implementation
  - `vietnamese_optimized.go`: High-performance version with buffer pooling
  - `testutil/`: Custom testing framework with comprehensive reporting
- **web/static/**: Frontend assets (HTML, CSS, JS)
- **scripts/**: Test runners and performance tools

### Key Design Patterns

- **Interface-based design**: `NumberConverter` interface allows swapping implementations
- **Performance optimization**: Buffer pooling, pre-allocated arrays, minimal allocations
- **Comprehensive testing**: Custom test framework with dataset validation
- **Production readiness**: Rate limiting, structured logging, health checks, graceful shutdown

### Vietnamese Language Handling

The converter implements specific Vietnamese linguistic rules:
- Number 1: "một" vs "mốt" depending on context
- Number 4: "bốn" vs "tư" in tens position
- Number 5: "năm" vs "lăm" in certain contexts
- Zero handling: "lẻ" for gaps in numbers like 101, 1001

### Testing Framework

Uses a custom testing framework (`pkg/converter/testutil/`) that:
- Loads test cases from files (`random_numbers_with_vietnamese.txt`)
- Provides detailed reporting and performance metrics
- Supports both original and optimized implementations
- Tracks conversion accuracy and processing times

### Performance Characteristics

- Target: Sub-25μs response times
- Memory: Minimal allocations through buffer pooling
- Throughput: Thousands of requests per second
- Benchmarking: Dedicated performance test suite

## API Endpoints

- `POST /api/v1/convert`: Convert number to Vietnamese text
- `GET /api/v1/convert`: Convert via URL parameters
- `GET /health`: Service health check
- `GET /ping`: Simple connectivity test
- `GET /`: Web UI interface

## Environment Configuration

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## Development Notes

- Uses Go 1.24.3 with Chi router, UUID generation, and rate limiting
- Two converter implementations: original and optimized for performance comparison
- Comprehensive test data includes edge cases and linguistic exceptions
- Web UI provides real-time conversion for manual testing

## 🚀 TURBO SERVICE - Ultimate Performance Architecture

### Perfect Service Design
The Turbo service represents the pinnacle of Go performance engineering - designed for 1000+ RPS with sub-100μs latency:

```bash
# Build and run the perfect service
make turbo-build
make turbo-run

# Performance testing
make turbo-benchmark
make turbo-load-test

# Production deployment
make turbo-deploy
```

### Architecture Principles
- **Zero Allocations**: All hot paths use memory pools and pre-computed tables
- **Cache-Friendly**: 64-byte aligned structures, NUMA-aware design
- **Minimal Dependencies**: Single binary with no external requirements
- **Perfect is Minimal**: Nothing to add, nothing to remove

### Performance Targets
- **Throughput**: 1000+ requests/second per core
- **Latency**: Sub-100μs response times (P95)
- **Memory**: <256MB footprint for production workload
- **CPU**: <50% utilization at target load

### Key Optimizations
1. **Pre-computed Lookup Tables**: All 0-999 combinations cached
2. **Buffer Pooling**: Zero-allocation request handling
3. **Direct JSON Parsing**: Bypass reflection with state machines
4. **Connection Pooling**: Optimized TCP settings and keep-alive
5. **Vietnamese Rule Engine**: Linguistic exceptions pre-compiled

### Deployment Options
- **Standalone**: `./turbo-service` (8MB binary)
- **Container**: `vietnamese-turbo:latest` (12MB image)
- **Load Balanced**: Nginx + multiple instances
- **Monitoring**: Prometheus metrics integration

### Files Structure
```
cmd/turbo/main.go           # Minimal service entry point
pkg/turbo/perfect.go        # Zero-allocation HTTP service
pkg/turbo/converter.go      # Ultimate converter implementation
pkg/turbo/*_test.go         # Comprehensive benchmarks and load tests
Dockerfile.turbo            # Production container image
docker-compose.turbo.yml    # Complete deployment stack
Makefile.turbo             # Build and deployment automation
```