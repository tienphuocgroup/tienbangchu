# Vietnamese Number Converter - Testing Framework

## Overview

This document describes the comprehensive testing framework set up for the Vietnamese Number Converter algorithm. The framework includes unit tests, integration tests, performance benchmarks, and a full dataset validation system using 1,694 real test cases.

## Test Files Structure

```
├── pkg/converter/
│   ├── vietnamese.go           # Core algorithm
│   └── vietnamese_test.go      # Comprehensive test suite
├── scripts/
│   └── run_tests.go           # Advanced test runner with reporting
├── random_numbers_with_vietnamese.txt  # 1,694 test cases
├── test_setup_validation.go   # Framework validation script
└── Makefile                   # Test automation targets
```

## Test Data

The testing framework uses `random_numbers_with_vietnamese.txt` which contains **1,694 test cases** in the format:
```
2355200847 hai tỷ ba trăm năm mươi lăm triệu hai trăm nghìn tám trăm bốn mươi bảy đồng
9106241390 chín tỷ một trăm lẻ sáu triệu hai trăm bốn mươi mốt nghìn ba trăm chín mươi đồng
...
```

Each line contains:
- A number (up to 13 digits)
- The expected Vietnamese translation with "đồng" currency

## Test Framework Components

### 1. Unit Tests (`vietnamese_test.go`)

**Basic Number Tests:**
- Zero handling
- Single digits (1-9)
- Teens (10-19)
- Tens (20, 30, 40, etc.)
- Hundreds with "lẻ" placement
- Thousands, millions, billions

**Edge Case Tests:**
- `4` vs `tư` vs `bốn` rules
- `1` vs `mốt` rules
- `5` vs `lăm` rules
- Zero handling with "lẻ"
- Scale transitions

**Error Handling Tests:**
- Negative numbers
- Numbers too large (>999 trillion)
- Input validation

**Performance Tests:**
- Sub-millisecond conversion requirement
- Throughput benchmarks
- Memory usage

### 2. Test Data Loader

**Features:**
- Loads test cases from file
- Validates file format
- Provides random access to test cases
- Memory-efficient streaming

**Usage:**
```go
loader := converter.NewTestDataLoader()
err := loader.LoadTestCases("random_numbers_with_vietnamese.txt")
testCases := loader.GetTestCases()
```

### 3. Test Suite Runner

**Capabilities:**
- Runs all 1,694 test cases
- Measures individual conversion times
- Generates detailed reports
- Tracks pass/fail rates
- Identifies specific failure patterns

**Usage:**
```go
suite := converter.NewTestSuite()
results, err := suite.RunAllTests("random_numbers_with_vietnamese.txt")
report := suite.GenerateReport(results)
```

### 4. Advanced Test Runner (`scripts/run_tests.go`)

**Features:**
- Command-line configuration
- Detailed reporting
- JSON output for CI/CD
- Performance benchmarking
- Failure analysis

## Running Tests

### Quick Validation
```bash
# Validate framework setup
go run test_setup_validation.go
```

### Unit Tests Only
```bash
# Standard go test
go test -v ./pkg/...

# Via Makefile
make test-unit
```

### Full Dataset Testing
```bash
# Basic run against all 1,694 cases
go run scripts/run_tests.go

# Verbose output with detailed failures
go run scripts/run_tests.go -verbose

# With performance testing
go run scripts/run_tests.go -perf

# Save detailed JSON report
go run scripts/run_tests.go -save

# Via Makefile
make test-full     # Comprehensive test
make test-quick    # Basic run
make test-runner   # With reporting
```

### Performance Benchmarks
```bash
# Go benchmarks
go test -bench=. -benchmem ./pkg/converter/

# Comprehensive performance test
make test-perf
```

### Custom Testing
```bash
# Custom test file
go run scripts/run_tests.go -file=custom_test_data.txt

# Limited failure output
go run scripts/run_tests.go -max-failures=5 -max-errors=2

# Using Makefile with custom args
make test-custom ARGS="-file=my_tests.txt -verbose"
```

## Test Runner Options

```bash
Usage: go run scripts/run_tests.go [options]

Options:
  -file string          Path to test data file (default: random_numbers_with_vietnamese.txt)
  -output string        Path to save JSON report
  -max-failures int     Max failed cases to display (default: 20)
  -max-errors int       Max error cases to display (default: 10)
  -perf                 Run additional performance tests
  -verbose              Verbose output with detailed analysis
  -save                 Save detailed report to timestamped JSON file
  -config string        Load configuration from JSON file
```

## Expected Performance Targets

### Accuracy Requirements
- **Pass Rate:** ≥95% on the full dataset
- **Zero Errors:** No conversion errors/exceptions
- **Edge Cases:** 100% accuracy on documented Vietnamese rules

### Performance Requirements
- **Average Time:** <2ms per conversion
- **Throughput:** >10,000 conversions/second
- **Memory:** <50MB total footprint

### Test Coverage
- **Unit Tests:** All basic number patterns
- **Integration:** Full 1,694 test case dataset
- **Edge Cases:** All documented Vietnamese language rules
- **Performance:** Benchmarks and stress tests

## Sample Test Report

```
=== Test Summary ===
Total Tests: 1694
Passed: 1598 (94.33%)
Failed: 96 (5.67%)
Errors: 0 (0.00%)
Total Time: 2.1s
Average Time: 1.2ms

=== Failed Cases (showing first 10) ===
Line 45: 34000
  Expected: ba mười bốn nghìn đồng
  Actual:   ba mười tư nghìn đồng

Line 127: 50050050
  Expected: năm mười triệu năm mười nghìn không trăm năm mười đồng
  Actual:   năm mười triệu năm mười nghìn năm mười đồng
```

## Integration with CI/CD

The test framework is designed for automated testing:

```bash
# Exit codes
0 = All tests passed (≥95% pass rate)
1 = Tests failed (<95% pass rate)

# JSON reports for parsing
go run scripts/run_tests.go -save -output=test_results.json
```

## Test-Driven Development Workflow

1. **Baseline Test:** Run full dataset to establish current performance
2. **Identify Issues:** Analyze failed cases and patterns
3. **Implement Fixes:** Modify algorithm based on test results
4. **Validate:** Re-run tests to verify improvements
5. **Performance Check:** Ensure changes don't impact speed
6. **Regression Test:** Verify no new failures introduced

## Key Vietnamese Language Rules Tested

### Number 4 Rules
- `24` → "hai mười **tư**" (not "bốn")
- `40` → "**bốn** mười" (not "tư")
- `34000` → "ba mười **bốn** nghìn" (not "tư")

### Number 1 Rules  
- `21` → "hai mười **mốt**" (not "một")
- `101` → "một trăm lẻ **một**" (not "mốt")

### Zero Handling
- `101` → "một trăm **lẻ** một"
- `1001` → "một nghìn không trăm **lẻ** một"
- `50050050` → proper "lẻ" placement

### Number 5 Rules
- `15` → "mười **lăm**" (not "năm")
- `25` → "hai mười **lăm**" (not "năm")

## Next Steps

1. **Run Baseline Test:** Execute the full test suite to see current algorithm performance
2. **Analyze Results:** Review failed cases to understand improvement areas
3. **Refactor Algorithm:** Make targeted improvements based on test failures
4. **Continuous Validation:** Use the framework during development to ensure progress

The test framework is now ready to support algorithm refactoring with confidence! 