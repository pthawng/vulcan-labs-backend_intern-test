# Promotion Code Validator

Backend intern coding assignment for Vulcan Labs.

## Problem

Determine if a promotion code is eligible for use. A code is eligible if and only if it exists in **both** data sources:
- `campaign_codes.txt`
- `membership_codes.txt`

**Constraints:**
- Files may contain millions of codes
- Files may not fit entirely in memory
- Each file contains one code per line
- Codes are 1-5 lowercase letters (a-z)

---

## High-Level Design

The implementation separates concerns into three layers:

### 1. Application Entry Point
`cmd/app/main.go` - Parses command-line arguments and coordinates the validation flow.

### 2. Promotion Validation Logic
`internal/promotion/` - Contains the core business logic:
- `service.go` - Service struct that depends on repository interfaces
- `validator.go` - Implements `IsEligible()` method with early exit optimization
- `validation.go` - Input validation (length, character constraints)

### 3. Repository Abstraction
`internal/repository/` - Abstracts data access:
- `repository.go` - `CodeRepository` interface with single `Exists(code)` method
- `file_repository.go` - Streaming file implementation
- `mock_repository.go` - In-memory implementation for testing

**Why repository abstraction?**
- Enables unit testing without file I/O
- Separates business logic from data access
- Makes it easier to swap implementations (e.g., database, cache)

---

## Validation Approach

### Algorithm: Streaming File Scan

The `FileCodeRepository` uses `bufio.Scanner` to read files line-by-line:

```go
func (r *FileCodeRepository) Exists(code string) (bool, error) {
    file, err := os.Open(r.filePath)
    if err != nil {
        return false, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        if scanner.Text() == code {
            return true, nil  // Early exit
        }
    }
    return false, scanner.Err()
}
```

### Why This Approach?

**Fits the constraints:**
- Does not load entire file into memory
- Handles arbitrarily large files
- Simple and correct

**Early exit behavior:**
The `IsEligible()` method checks the campaign file first. If the code is not found there, it returns `false` immediately without checking the membership file. This reduces unnecessary I/O for codes that don't exist in the campaign system.

---

## Data Structures & Complexity

### Time Complexity
- **Single validation**: O(n + m) worst case
  - n = number of codes in campaign file
  - m = number of codes in membership file
- **Best case**: O(k) where k is the position of the code in the campaign file (early exit)

### Space Complexity
- **O(1)** - Uses a fixed-size buffer (~4KB) for scanning
- Does not depend on file size

### Realistic Performance
For a file with 1 million codes:
- Worst case: ~2 million line reads (code not in campaign file)
- Average case: ~500K line reads (code found midway in campaign file)
- This is acceptable for a command-line tool with infrequent usage

---

## Trade-offs & Alternatives

### Alternative 1: Load into Hash Map
**Approach:** Read both files into `map[string]bool` at startup.

**Pros:**
- O(1) lookup after initial load
- Fast for repeated validations

**Cons:**
- O(n) memory usage (~25MB for 500K codes)
- Violates "may not fit in memory" constraint
- Overkill for single-use CLI tool

**Why not chosen:** Assignment explicitly states files may not fit in memory.

### Alternative 2: Bloom Filter
**Approach:** Use probabilistic data structure for membership testing.

**Pros:**
- Space-efficient (1-2 bytes per element)
- Fast lookups

**Cons:**
- False positives possible
- Requires external library
- Over-engineering for this scope

**Why not chosen:** Adds complexity without clear benefit for the assignment's requirements.

### Alternative 3: External Sort + Merge
**Approach:** Sort both files, then merge to find intersection.

**Pros:**
- O(n log n) time
- O(1) space if using external sort

**Cons:**
- Requires writing temporary files
- More complex implementation
- Slower for single lookups

**Why not chosen:** Optimizes for batch processing, not single code validation.

---

## Testing & Benchmarking

### Unit Tests
Located in `*_test.go` files:

**Validation tests** (`validation_test.go`):
- Valid codes
- Empty code
- Code too long (>5 chars)
- Invalid characters (uppercase, numbers, special chars)

**Integration tests** (`validator_test.go`):
- Code exists in both systems → `true`
- Code exists in only one system → `false`
- Code exists in neither system → `false`

**Repository tests** (`file_repository_test.go`):
- Code found in file (first line, middle, last line)
- Code not found
- File not found error handling

### Mock Repository
`mock_repository.go` provides an in-memory implementation using a hash map. This allows testing the service layer without file I/O.

### Benchmarks
`validator_bench_test.go` contains a single benchmark:
- `BenchmarkIsEligible_HappyPath` - Measures performance when code exists in both systems

**Note:** Benchmarks use mock repositories with 1000 codes, not real file I/O. They establish a baseline for the service logic, not end-to-end file scanning performance.

---

## Assumptions & Limitations

### File Format Assumptions
- One code per line
- No leading/trailing whitespace
- UTF-8 encoding
- Newline-delimited (LF or CRLF)

### Code Constraints
- Length: 1-5 characters
- Characters: lowercase a-z only
- These constraints are validated before file lookup

### Scope Limitations
- **Single-threaded:** No concurrent validation support
- **No caching:** Each validation scans files from scratch
- **No file watching:** Does not detect file changes
- **No transaction support:** Assumes files are read-only
- **Error handling:** Basic error propagation, no retry logic

### What This Implementation Does NOT Do
- Does not handle concurrent writes to data files
- Does not cache file contents
- Does not provide an HTTP/gRPC API
- Does not include logging or metrics
- Does not handle file rotation or updates

---

## How to Run

### Build
```bash
go build -o validator ./cmd/app
```

### Run
```bash
./validator <code> <campaign_file> <membership_file>
```

**Example:**
```bash
./validator promo data/campaign_codes.txt data/membership_codes.txt
```

**Output:** `true` or `false`

### Run Tests
```bash
# All tests
go test ./...

# With verbose output
go test ./... -v

# With coverage
go test ./... -cover
```

### Run Benchmarks
```bash
go test -bench=. ./internal/promotion/
```

---

## Project Structure

```
.
├── cmd/app/
│   └── main.go                    # Entry point
├── internal/
│   ├── promotion/
│   │   ├── service.go             # Service struct
│   │   ├── validator.go           # IsEligible logic
│   │   ├── validation.go          # Input validation
│   │   ├── validator_test.go      # Integration tests (3 cases)
│   │   ├── validation_test.go     # Validation tests (4 cases)
│   │   └── validator_bench_test.go # Benchmarks (1 case)
│   └── repository/
│       ├── repository.go          # Interface definition
│       ├── file_repository.go     # Streaming implementation
│       ├── file_repository_test.go # Repository tests (3 cases)
│       └── mock_repository.go     # Mock for testing
├── data/
│   ├── campaign_codes.txt         # Campaign data
│   └── membership_codes.txt       # Membership data
└── scripts/
    └── generate_data.go           # Test data generator
```

---

## Implementation Notes

### Why Streaming?
The assignment states "files may not fit entirely into memory." The streaming approach ensures the solution works regardless of file size, even if it means slower lookups.

### Why Early Exit?
Checking the campaign file first and returning early if the code is not found reduces unnecessary I/O. This is a simple optimization that doesn't add complexity.

### Why Repository Pattern?
The repository abstraction allows testing the business logic (`IsEligible`) without file I/O. This makes tests faster and more reliable.

---

**Author:** Le Phuoc Thang  
**Assignment:** Backend Intern Coding Test - Vulcan Labs  
**Date:** February 2026
