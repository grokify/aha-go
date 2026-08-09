package aha

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestRateLimitFromHeaders(t *testing.T) {
	tests := []struct {
		name    string
		headers http.Header
		want    *RateLimitInfo
	}{
		{
			name:    "no rate limit headers",
			headers: http.Header{},
			want:    nil,
		},
		{
			name: "all headers present",
			headers: http.Header{
				"X-Ratelimit-Limit":     []string{"300"},
				"X-Ratelimit-Remaining": []string{"42"},
				"X-Ratelimit-Reset":     []string{"1498005537"},
			},
			want: &RateLimitInfo{
				Limit:     300,
				Remaining: 42,
				Reset:     time.Unix(1498005537, 0),
			},
		},
		{
			name: "partial headers still parsed",
			headers: http.Header{
				"X-Ratelimit-Remaining": []string{"0"},
			},
			want: &RateLimitInfo{
				Remaining: 0,
			},
		},
		{
			name: "malformed values ignored",
			headers: http.Header{
				"X-Ratelimit-Limit": []string{"not-a-number"},
			},
			want: &RateLimitInfo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rateLimitFromHeaders(tt.headers)
			if tt.want == nil {
				if got != nil {
					t.Errorf("rateLimitFromHeaders() = %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("rateLimitFromHeaders() = nil, want non-nil")
			}
			if got.Limit != tt.want.Limit {
				t.Errorf("Limit = %d, want %d", got.Limit, tt.want.Limit)
			}
			if got.Remaining != tt.want.Remaining {
				t.Errorf("Remaining = %d, want %d", got.Remaining, tt.want.Remaining)
			}
			if !got.Reset.Equal(tt.want.Reset) {
				t.Errorf("Reset = %v, want %v", got.Reset, tt.want.Reset)
			}
		})
	}
}

func TestRateLimitFromError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want *RateLimitInfo
	}{
		{
			name: "nil error",
			err:  nil,
			want: nil,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: nil,
		},
		{
			name: "API error without rate limit info",
			err:  &APIError{StatusCode: http.StatusNotFound},
			want: nil,
		},
		{
			name: "API error with rate limit info",
			err: &APIError{
				StatusCode: http.StatusTooManyRequests,
				RateLimit:  &RateLimitInfo{Limit: 300, Remaining: 0},
			},
			want: &RateLimitInfo{Limit: 300, Remaining: 0},
		},
		{
			name: "wrapped API error with rate limit info",
			err: fmt.Errorf("wrapped: %w", &APIError{
				StatusCode: http.StatusTooManyRequests,
				RateLimit:  &RateLimitInfo{Limit: 300, Remaining: 5},
			}),
			want: &RateLimitInfo{Limit: 300, Remaining: 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RateLimitFromError(tt.err)
			if tt.want == nil {
				if got != nil {
					t.Errorf("RateLimitFromError() = %+v, want nil", got)
				}
				return
			}
			if got == nil {
				t.Fatal("RateLimitFromError() = nil, want non-nil")
			}
			if got.Limit != tt.want.Limit || got.Remaining != tt.want.Remaining {
				t.Errorf("RateLimitFromError() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
