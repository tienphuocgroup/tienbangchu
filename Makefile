.PHONY: build test test-unit test-full test-perf test-runner run docker clean fmt

build:
	CGO_ENABLED=0 GOOS=linux go build -o bin/server cmd/server/main.go

# Standard unit tests only
test-unit:
	go test -v ./pkg/...

# Full test suite including unit tests and integration tests
test:
	go test -v ./...

# Run the comprehensive test suite against the full dataset
test-full:
	go run scripts/run_tests.go -file=random_numbers_with_vietnamese.txt -verbose

# Run performance benchmarks
test-perf:
	go test -bench=. -benchmem ./pkg/converter/
	go run scripts/run_tests.go -file=random_numbers_with_vietnamese.txt -perf

# Run test suite with detailed reporting
test-runner:
	go run scripts/run_tests.go -file=random_numbers_with_vietnamese.txt -verbose -save

# Run test suite with custom parameters
test-custom:
	@echo "Usage: make test-custom ARGS='-file=your_test_file.txt -verbose'"
	go run scripts/run_tests.go $(ARGS)

# Quick test with basic output
test-quick:
	go run scripts/run_tests.go -file=random_numbers_with_vietnamese.txt

run:
	go run cmd/server/main.go

docker:
	docker build -t vietnamese-converter:latest .

docker-run:
	docker run -p 8080:8080 vietnamese-converter:latest

clean:
	rm -rf bin/ test_report_*.json

fmt:
	go fmt ./...
