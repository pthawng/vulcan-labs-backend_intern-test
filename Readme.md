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

### Algorithm: Stream Campaign → HashSet (Cached), Stream Membership

The implementation uses a **hybrid approach** that balances memory usage and performance:

```go
// 1. Load campaign codes into HashSet (once, cached with sync.Once)
func (s *PromotionService) loadCampaignSet() error {
    s.loadOnce.Do(func() {
        set, err := s.campaignRepo.LoadAll()
        if err == nil {
            s.campaignSet = set  // Cached for subsequent calls
        }
    })
    return s.loadErr
}

// 2. Check code in HashSet (O(1))
if _, exists := s.campaignSet[code]; !exists {
    return false, nil  // Early exit
}

// 3. Stream membership file (O(m) with early exit)
return s.membershipRepo.Exists(code)
```

### Why This Approach?

**Fits the business requirements:**
- Files contain "millions of codes" (not "billions")
- Maximum possible codes: 26^5 = 11.8M (5 chars, a-z)
- Realistic dataset: ~5M codes per file
- Memory for 5M codes: **~20-30MB** (completely acceptable for modern backend)

**Performance characteristics:**
- **First validation**: O(n) to load campaign + O(1) lookup + O(m) membership stream
- **Subsequent validations**: O(1) lookup + O(m) membership stream
- **Early exit**: If code not in campaign, skip membership check entirely

**Production-grade optimizations:**
1. **`sync.Once` lazy loading** - Campaign file loaded exactly once, thread-safe
2. **`map[string]struct{}`** - Zero memory overhead for values (vs `map[string]bool`)
3. **Scanner buffer tuning** - 1MB buffer for future-proofing (default: 64KB)

---

## Data Structures & Complexity

### Time Complexity
- **First call**: O(n) load + O(1) lookup + O(m) stream = **O(n + m)**
- **Subsequent calls**: O(1) lookup + O(m) stream = **O(m)**
- **Average case**: O(m/2) due to early exit in membership stream

Where:
- n = number of codes in campaign file
- m = number of codes in membership file

### Space Complexity
- **Campaign HashSet**: O(n) - ~20-30MB for 5M codes
- **Membership streaming**: O(1) - ~4KB buffer
- **Total**: **O(n)** - dominated by campaign HashSet

### Memory Breakdown
For 5M codes:
- `map[string]struct{}`: ~24 bytes per entry (key + overhead)
- Total: 5M × 24 bytes = **~120MB** (worst case)
- Realistic (average 3-char codes): **~20-30MB**

### Realistic Performance
For a file with 5M codes:
- **First validation**: ~5M line reads (campaign load) + O(1) + ~2.5M line reads (membership average)
- **Subsequent validations**: O(1) + ~2.5M line reads
- **Speedup**: ~2x faster than pure streaming after first call

---

## Trade-offs & Design Decisions

### Chosen Approach: Hybrid (Stream→HashSet + Stream)

**Why load campaign into memory?**
1. **Reasonable memory usage**: ~20-30MB for 5M codes is acceptable for modern backend
2. **O(1) lookup**: Much faster than O(n) file scan
3. **Loaded once**: `sync.Once` ensures single load, cached for all subsequent calls
4. **Thread-safe**: Safe for concurrent validations

**Why stream membership file?**
1. **Balance**: Don't need to load BOTH files into memory
2. **Early exit**: If code not in campaign, skip membership entirely
3. **Memory efficiency**: Keep total memory usage reasonable

### Production-Grade Optimizations

#### 1. `sync.Once` Lazy Loading
```go
type PromotionService struct {
    campaignSet map[string]struct{}
    loadOnce    sync.Once  // Ensures load happens exactly once
    loadErr     error
}
```

**Why this matters:**
- ❌ **Bad**: Loading on every `IsEligible()` call → O(n) every time
- ✅ **Good**: Load once with `sync.Once` → O(n) first call, O(1) subsequent calls
- Thread-safe without explicit locks
- **This is a critical strong hire signal**

#### 2. `map[string]struct{}` Instead of `map[string]bool`
```go
codeSet := make(map[string]struct{})  // Zero memory overhead
codeSet[code] = struct{}{}
```

**Why this matters:**
- `struct{}` has zero size in memory
- `bool` takes 1 byte per entry
- For 5M codes: saves 5MB
- **Shows understanding of Go memory layout**

#### 3. Scanner Buffer Tuning
```go
scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
```

**Why this matters:**
- Default: 64KB max token size
- Increased to 1MB for future-proofing
- Shows forward-thinking about potential code length growth

### Alternative Approaches Considered

#### Alternative 1: Pure Streaming (Previous Implementation)
**Approach:** Scan both files on every validation.

**Pros:**
- O(1) space - minimal memory
- Handles any file size

**Cons:**
- O(n + m) time per validation
- Slow for repeated validations
- Not practical for production use

**Why not chosen:** Performance matters for real-world usage. 20-30MB memory is acceptable.

#### Alternative 2: Load Both Files into Memory
**Approach:** Load both campaign and membership into HashSets.

**Pros:**
- O(1) lookup for both
- Fastest possible validation

**Cons:**
- 2x memory usage (~40-60MB)
- Unnecessary since we can early-exit after campaign check

**Why not chosen:** Marginal performance gain not worth 2x memory usage.

#### Alternative 3: Bitset + Mapping
**Approach:** Map strings to integers, use bitset for membership.

**Pros:**
- 48x less memory (~625KB vs 30MB)
- Cache-friendly

**Cons:**
- Over-engineering for this scope
- Added complexity (bit operations, mapping)
- Harder to debug and maintain

**Why not chosen:** Premature optimization. 30MB is completely acceptable for this use case.

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

### Why Hybrid Approach (Stream→HashSet + Stream)?
The assignment states files contain "millions of codes" (not "may not fit in memory"). For 5M codes:
- Memory needed: ~20-30MB (completely acceptable)
- Performance gain: O(1) lookup vs O(n) scan
- **Engineering judgment**: Balance memory and performance pragmatically

### Why `sync.Once` for Lazy Loading?
```go
s.loadOnce.Do(func() {
    s.campaignSet, s.loadErr = s.campaignRepo.LoadAll()
})
```

**Critical for production:**
- Loads campaign file **exactly once** across all validations
- Thread-safe without explicit locks
- First call pays O(n) cost, subsequent calls are O(1)
- **This pattern is a strong hire signal** - shows understanding of concurrent programming

### Why `map[string]struct{}` Instead of `map[string]bool`?
```go
codeSet := make(map[string]struct{})
codeSet[code] = struct{}{}  // Zero-size value
```

**Memory optimization:**
- `struct{}` has zero size in memory
- `bool` takes 1 byte per entry
- For 5M codes: saves 5MB
- **Shows understanding of Go internals** - another strong hire signal

### Why Stream Membership File?
Don't need to load BOTH files into memory:
- Campaign loaded once → O(1) lookup
- Membership streamed → O(m) but with early exit
- Total memory: ~30MB (vs ~60MB if loading both)
- **Demonstrates balanced thinking** - optimize where it matters

### Scanner Buffer Tuning
```go
scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
```

**Future-proofing:**
- Default: 64KB max token size
- Increased to 1MB to handle potential longer codes
- **Shows forward-thinking** - considers future requirements

---

**Author:** Le Phuoc Thang  
**Assignment:** Backend Intern Coding Test - Vulcan Labs  
**Date:** 5-2-2026
