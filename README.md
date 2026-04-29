# Vietnamese Number Converter Service

A high-performance Go backend service that converts numeric values to properly formatted Vietnamese text, handling all linguistic exceptions and edge cases.

## Features

- ⚡ **High Performance**: Sub-25μs response times for most conversions
- 🎯 **Accurate**: Implements Vietnamese number conversion rules including:
  - Special cases for numbers 1 (một/mốt) and 4 (tư/bốn)
  - Proper handling of zero (lẻ) in numbers like 101, 1001, etc.
  - Correct scale transitions (thousands, millions, billions, trillions)
- 🚀 **Efficient**: Optimized implementation with minimal allocations
- 🛡️ **Production Ready**: Comprehensive error handling and logging
- 📦 **Container Ready**: Easy Docker deployment

## Quick Start

1. **Run the service:**
   ```bash
   make run
   ```

2. **Test the conversion:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/convert \
     -H "Content-Type: application/json" \
     -d '{"number": 1433433225}'
   ```

## API Endpoints

### Convert Number to Vietnamese Text

`POST /api/v1/convert`

Convert a number to Vietnamese text representation with currency.

**Request:**
```json
{
  "number": 1433433225
}
```

**Successful Response (200 OK):**
```json
{
  "number": 1433433225,
  "vietnamese": "một tỷ bốn trăm ba mươi ba triệu bốn trăm ba mươi ba nghìn hai trăm hai mươi lăm đồng",
  "processing_time_ms": 0.024084
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "error": "Internal Server Error",
  "code": 500
}
```

### Health Check

`GET /health`

Check if the service is running.

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-06-06T00:47:46+07:00"
}
```

## Usage Examples

```bash
# Convert a number (successful conversion)
curl -X POST http://localhost:8080/api/v1/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 1433433225}'

# Convert a very large number (may hit the limit)
curl -X POST http://localhost:8080/api/v1/convert \
  -H "Content-Type: application/json" \
  -d '{"number": 1433433212125}'

# Health check
curl http://localhost:8080/health
```

## Limitations

- **Number Range**: The service handles numbers up to 999 trillion (999,999,999,999,999) accurately.
- **Large Numbers**: Numbers larger than approximately 10^15 may cause internal server errors.
- **Negative Numbers**: Currently not supported (will result in error).
- **Decimals**: Only whole numbers are supported.

## Performance

- **Typical Response Time**: < 0.05ms for most conversions
- **Memory Usage**: Minimal, with efficient memory pooling
- **Throughput**: Capable of handling thousands of requests per second

## Development

### Prerequisites

- Go 1.21+
- Make (optional, for convenience commands)

### Building and Running

```bash
# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go

# Or use make
make run

# Build binary
make build

# Run tests
make test
```

### Testing

Run the full test suite:
```bash
make test
```

Run performance benchmarks:
```bash
make test-perf
```

## API Documentation

- **[API Reference](docs/api-reference.md)** — Endpoint details, request/response examples, client code samples
- **[OpenAPI 3.0 Spec](docs/api-openapi-spec.yaml)** — Machine-readable spec for auto-generating client SDKs
- **[Deployment Guide](docs/deployment-guide.md)** — Docker, Nginx, Kubernetes, and scaling configurations

## Deployment

### Docker

Build the Docker image:
```bash
docker build -t vietnamese-converter .
```

Run the container:
```bash
docker run -p 8080:8080 vietnamese-converter
```

### Environment Variables

- `PORT`: Port to run the server on (default: 8080)
- `LOG_LEVEL`: Logging level (debug, info, warn, error) (default: info)

## Project Structure

```
.
├── cmd/
│   └── server/           # Main application entry point
│       └── main.go
├── internal/
│   ├── api/             # API handlers and routes
│   ├── config/          # Configuration management
│   └── logger/          # Logging utilities
├── pkg/
│   └── converter/       # Core conversion logic
│       ├── vietnamese.go        # Original implementation
│       ├── vietnamese_test.go   # Tests
│       └── vietnamese_optimized.go  # Optimized implementation
├── scripts/             # Utility scripts
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

[MIT](LICENSE)
