package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Uncomment the function you want to test:
	// TestForHandleConcurencyInHandlerRequestWeb()
	// Example1ForContext()
	// TestForHttpClient()
	// ExampleConnectionPool()
	// MainForBackpressure()
	// MainForIOWork()
}

// ============================================================================
// SECTION 1: HTTP SERVER CONCURRENCY IN GO
// ============================================================================
//
// IMPORTANT: Does Go handle each request concurrently?
// YES! When you use http.ListenAndServe(), for each new request:
//   - A new goroutine is created
//   - Your handler runs inside that goroutine
//   - All requests are completely concurrent
//
// Example: If 10 people call the same endpoint simultaneously:
//   - 10 goroutines will run concurrently
//
// ⚠️ CRITICAL WARNING: If you have shared state:
//   - Modifying it without synchronization creates DATA RACES
//   - Always use: sync.Mutex, sync/atomic, channels, or worker pools
// ============================================================================

var (
	counter int        // Shared resource - requires synchronization
	mu      sync.Mutex // Protects counter from concurrent access
)

// TestForHandleConcurencyInHandlerRequestWeb demonstrates safe concurrent counter access
// Each request increments and displays the counter safely using mutex locking
func TestForHandleConcurencyInHandlerRequestWeb() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()   // Acquire lock to prevent race conditions
		counter++   // Critical section - only one goroutine at a time
		mu.Unlock() // Release lock for others
		fmt.Fprintf(w, "Counter is : %d", counter)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// ============================================================================
// SECTION 2: CONTEXT IN HTTP HANDLERS
// ============================================================================
//
// context.Context is ONE OF THE MOST IMPORTANT concepts in Go.
// In HTTP, each request automatically has a context: r.Context()
//
// What can Context do?
//   1. Cancellation (when client disconnects)
//   2. Timeout handling
//   3. Deadline management
//   4. Metadata transfer (e.g., correlation IDs)
//
// WHY Context is CRITICAL?
// Without Context:
//   - Request is cancelled but goroutine keeps working
//   - Database queries continue unnecessarily
//   - Memory leaks occur
//   - Resources are wasted
//
// With Context:
//   - Everything cascades and cancels properly
// ============================================================================

// Example1ForContext demonstrates basic context cancellation detection
func Example1ForContext(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	select {
	case <-ctx.Done(): // Client disconnected/cancelled
		fmt.Fprintf(w, "Context is done")
	case <-time.After(1 * time.Second): // Operation completes
		fmt.Fprintf(w, "Context One Second TimeOut")
	default:
		return
	}
}

// Example2ForContext demonstrates timeout handling with context
// This is useful for operations that might take too long
func Example2ForContext(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel() // IMPORTANT: Always call cancel to release resources

	result := make(chan string)

	go func() {
		time.Sleep(3 * time.Second) // Simulates slow operation
		result <- "Finished"
	}()

	select {
	case res := <-result:
		fmt.Fprintln(w, res)
	case <-ctx.Done(): // Timeout occurred (2 seconds passed)
		http.Error(w, "Request Timeout", http.StatusRequestTimeout)
	}
}

type contextKey string

const CorrelationKey contextKey = "cid"

// correlationMiddleware adds correlation ID tracking to each request
// This is essential for distributed tracing and debugging
func correlationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get correlation ID from request header
		cid := r.Header.Get("X-Correlation-ID")
		if cid == "" {
			// Generate new ID if none exists
			cid = uuid.New().String()
		}

		// Add correlation ID to context for downstream functions
		ctx := context.WithValue(r.Context(), CorrelationKey, cid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// PRODUCTION-LEVEL BEST PRACTICES FOR CONTEXT:
// ✅ ALWAYS pass context to lower layers (database, API calls, etc.)
// ✅ ALWAYS call cancel() when using context.WithTimeout/WithCancel
// ❌ NEVER store context inside structs
// ❌ NEVER use context for optional parameters
// ✅ ONLY use context.WithValue for request-scoped metadata

// ============================================================================
// SECTION 3: PROPER HTTP CLIENT CONFIGURATION
// ============================================================================
//
// In Go, DON'T use http.Get directly. Instead create an http.Client with:
//   - Timeout settings (prevents hanging requests)
//   - Proper Transport (enables connection reuse)
//   - Resource limits (controls memory/CPU usage)
// ============================================================================

func TestForHttpClient() {
	client := &http.Client{
		Timeout: 5 * time.Second, // Total request timeout
		Transport: &http.Transport{
			MaxIdleConns:        10,               // Total idle connections
			MaxIdleConnsPerHost: 5,                // Per-host connection limit
			IdleConnTimeout:     30 * time.Second, // Close idle connections after 30s
		},
	}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close() // IMPORTANT: Always close response body

	fmt.Println("Status:", resp.Status)
}

// ============================================================================
// SECTION 4: CONNECTION POOLING
// ============================================================================
//
// Connection Pool = Reusing HTTP connections instead of creating new ones
// Benefits:
//   - Lower latency (no TCP handshake per request)
//   - Reduced CPU usage
//   - Better resource utilization
//   - Higher throughput
// ============================================================================

func ExampleConnectionPool() {
	tr := &http.Transport{
		MaxIdleConns:        100,              // Maximum idle connections pool-wide
		MaxIdleConnsPerHost: 10,               // Per-host limit
		IdleConnTimeout:     60 * time.Second, // Keep connections alive
	}

	client := &http.Client{Transport: tr}

	// These requests will reuse connections, not create new ones each time
	for i := 0; i < 5; i++ {
		resp, err := client.Get("https://httpbin.org/get")
		if err != nil {
			panic(err)
		}
		fmt.Println("Request", i, resp.Status)
		resp.Body.Close() // Body must be closed to reuse connection
	}
}

// ============================================================================
// SECTION 5: BACKPRESSURE IN I/O
// ============================================================================
//
// Backpressure = When consumer is slow, producer must wait to prevent overload
// In Go, this is typically achieved using buffered channels:
//   - Buffer acts as a queue
//   - When buffer fills, producer blocks
//   - Prevents system overload
// ============================================================================

func worker(jobs chan int) {
	for j := range jobs {
		fmt.Println("processing", j)
		time.Sleep(time.Second) // Simulates slow processing
	}
}

func MainForBackpressure() {
	// Buffered channel with capacity 3 creates backpressure
	// When 3 jobs are queued, producer will block
	jobs := make(chan int, 3)

	go worker(jobs)

	// Producer sends 10 jobs but will block when buffer is full
	for i := 0; i < 10; i++ {
		fmt.Println("sending", i)
		jobs <- i // This blocks when channel buffer is full
	}

	close(jobs) // Signal no more jobs
	time.Sleep(5 * time.Second)
}

// ============================================================================
// SECTION 6: FILE I/O OPERATIONS
// ============================================================================
//
// Go uses 'os' and 'io' packages for file operations
// For large files, use streaming or buffered I/O:
//   - bufio for buffered reading/writing
//   - io.Reader/io.Writer for streaming
//   - os.File for direct file access
// ============================================================================

func MainForIOWork() {
	// CREATE AND WRITE to file
	file, err := os.Create("test.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close() // Always close files to avoid leaks

	file.WriteString("Hello Go File\n")
	// No need to call file.Close() here if using defer above

	// READ entire file (good for small files)
	data, err := os.ReadFile("test.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(data))
}

// ============================================================================
// KEY TAKEAWAYS & SUMMARY
// ============================================================================
//
// 1. CONCURRENCY:
//    - Each HTTP request runs in its own goroutine
//    - Always protect shared state with mutex, atomic, or channels
//
// 2. CONTEXT:
//    - Always propagate context to lower layers
//    - Use timeouts to prevent hanging operations
//    - Never store context in structs
//
// 3. HTTP CLIENT:
//    - Don't use http.Get directly
//    - Configure timeout and transport always
//    - Close response bodies
//
// 4. CONNECTION POOLING:
//    - Reuse connections for better performance
//    - Configure MaxIdleConns and IdleConnTimeout
//
// 5. BACKPRESSURE:
//    - Use buffered channels to control flow
//    - Prevents system overload
//
// 6. FILE I/O:
//    - Always close files (use defer)
//    - Use streaming for large files
//    - Handle errors properly
// ============================================================================
