package aha

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// TestNewClient_RetryOn429 verifies that requests are automatically retried
// on 429 by default, and that retry applies to DoRaw() (not just the typed
// ogen client) since both share the same wrapped transport.
func TestNewClient_RetryOn429(t *testing.T) {
	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.AddInt32(&reqCount, 1) == 1 {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		WithSubdomain("test"),
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithMaxRetries(1),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.DoRaw(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoRaw() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if got := atomic.LoadInt32(&reqCount); got != 2 {
		t.Errorf("reqCount = %d, want 2 (initial attempt + 1 retry)", got)
	}
}

// TestNewClient_RetryDisabled verifies WithRetryDisabled() prevents retry.
func TestNewClient_RetryDisabled(t *testing.T) {
	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client, err := NewClient(
		WithSubdomain("test"),
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithRetryDisabled(),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	resp, err := client.DoRaw(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("DoRaw() error = %v", err)
	}
	defer resp.Body.Close()

	if got := atomic.LoadInt32(&reqCount); got != 1 {
		t.Errorf("reqCount = %d, want 1 (retry disabled)", got)
	}
}

// TestNewClient_RequestsPerSecond verifies WithRequestsPerSecond throttles
// requests once the initial burst is exhausted, and that throttling applies
// to DoRaw() as well as the typed ogen client.
func TestNewClient_RequestsPerSecond(t *testing.T) {
	var reqCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&reqCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		WithSubdomain("test"),
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
		WithRequestsPerSecond(5),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	const numRequests = 8 // 5 immediate (burst) + 3 throttled at 5/s
	start := time.Now()
	for i := 0; i < numRequests; i++ {
		resp, err := client.DoRaw(context.Background(), http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("DoRaw() error = %v", err)
		}
		resp.Body.Close()
	}
	elapsed := time.Since(start)

	if got := atomic.LoadInt32(&reqCount); got != numRequests {
		t.Errorf("reqCount = %d, want %d", got, numRequests)
	}
	// 3 requests beyond the burst of 5, at 5 req/s, take ~600ms; assert a
	// conservative lower bound to avoid flakiness while still proving
	// throttling occurred (unthrottled, this loop completes in <10ms).
	if elapsed < 400*time.Millisecond {
		t.Errorf("elapsed = %v, want at least 400ms (requests were not throttled)", elapsed)
	}
}

// TestNewClient_NoThrottleByDefault verifies that without
// WithRequestsPerSecond, requests are not throttled.
func TestNewClient_NoThrottleByDefault(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, err := NewClient(
		WithSubdomain("test"),
		WithAPIKey("test-key"),
		WithBaseURL(server.URL),
	)
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	start := time.Now()
	for i := 0; i < 20; i++ {
		resp, err := client.DoRaw(context.Background(), http.MethodGet, "/test", nil)
		if err != nil {
			t.Fatalf("DoRaw() error = %v", err)
		}
		resp.Body.Close()
	}
	elapsed := time.Since(start)

	if elapsed > 400*time.Millisecond {
		t.Errorf("elapsed = %v, want under 400ms (requests should not be throttled by default)", elapsed)
	}
}
