# Promotion Code Validator

Validates whether a promotion code exists in both campaign and membership systems.

## Solution Overview

A code is **eligible** only if it exists in **both** data sources.

### Key Features
- ✅ **Memory Efficient**: O(1) space - streams files line-by-line
- ✅ **Handles Large Files**: Processes millions of codes without loading into memory
- ✅ **Clean Architecture**: Repository pattern with dependency injection
- ✅ **Input Validation**: Defensive programming with custom error types
- ✅ **Comprehensive Tests**: 20+ test cases with mocks and benchmarks

---

## Architecture

```
cmd/app/main.go          (Entry point)
    ↓
internal/promotion/
├── service.go           (Service struct)
├── validator.go         (Eligibility logic)
└── validation.go        (Input validation)
    ↓
internal/repository/
├── repository.go        (Interface)
├── file_repository.go   (Streaming implementation)
└── mock_repository.go   (Testing)
```

---

## Design Decisions & Trade-offs

### Current: Streaming Approach

**Why chosen**: Assignment says "files may not fit into memory"

```go
scanner := bufio.NewScanner(file)
for scanner.Scan() {
    if scanner.Text() == code {
        return true, nil  // Early exit when found
    }
}
```

| Metric | Value |
|--------|-------|
| **Time Complexity** | O(n) per lookup |
| **Space Complexity** | O(1) |
| **Memory Usage** | ~4KB buffer |

**Trade-offs**:
- ✅ **Pros**: Minimal memory, handles any file size, simple
- ⚠️ **Cons**: Slower for repeated lookups

### Alternative Considered: Caching

For high-frequency validation (1000+ req/sec), could add hash map caching:
- **Time**: O(1) lookup after initial load
- **Space**: O(n) - ~25MB for 500K codes
- **When**: Production systems with frequent validations

**Not implemented** to avoid over-engineering for this scope.

---

## Input Validation

```go
ValidateCode("solar")   // ✅ Valid
ValidateCode("")        // ❌ Error: code cannot be empty
ValidateCode("PROMO")   // ❌ Error: must be lowercase
ValidateCode("promo1")  // ❌ Error: must contain only a-z
ValidateCode("abcdef")  // ❌ Error: must be at most 5 characters
```

---

## Usage

```bash
# Build
go build -o validator cmd/app/main.go

# Run
./validator <code> <campaign_file> <membership_file>

# Example
./validator promo data/campaign_codes.txt data/membership_codes.txt
```

**Output**: `true` or `false`

---

## Testing

```bash
# Run all tests
go test ./...

# Run with verbose
go test ./... -v

# Run benchmarks
go test -bench=. ./...

# Test coverage
go test ./... -cover
```

### Test Coverage
- ✅ Service tests (9 cases): both systems, single system, errors
- ✅ Validation tests (10 cases): empty, invalid chars, length
- ✅ Repository tests (4 cases): exists, not found, empty file
- ✅ Benchmark tests (3 scenarios)

---

## Project Structure

```
.
├── cmd/app/main.go                    # Entry point
├── internal/
│   ├── promotion/
│   │   ├── service.go                 # Service struct
│   │   ├── validator.go               # IsEligible logic
│   │   ├── validation.go              # Input validation
│   │   ├── validator_test.go          # Service tests
│   │   ├── validation_test.go         # Validation tests
│   │   └── validator_bench_test.go    # Benchmarks
│   └── repository/
│       ├── repository.go              # Interface
│       ├── file_repository.go         # Streaming impl
│       ├── file_repository_test.go    # Repository tests
│       └── mock_repository.go         # Mock for testing
├── data/                              # Test data files
└── scripts/generate_data.go           # Data generator
```

---

## Complexity Analysis

|     Operation     | Time | Space |
|-------------------|------|-------|
| Single validation | O(n) |  O(1) |
| Input validation  | O(m) |  O(1) |

*n = codes in file, m = code length (max 5)*

---

## Assumptions

1. **File format**: One code per line, UTF-8, no trailing whitespace
2. **Code constraints**: 1-5 lowercase letters (a-z)
3. **Files**: Readable, exist at specified path
4. **Codes**: Unique within each file

---

## Future Enhancements

For production with high-frequency validations:
- [ ] Add caching layer (hash map or Redis)
- [ ] Add context.Context for cancellation
- [ ] Add structured logging
- [ ] Add metrics collection
- [ ] Add HTTP/gRPC API

---

**Author**: Le Phuoc Thang
**Last Updated**: 2026-02-05
