package aha

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/time/rate"
)

// RateLimitInfo captures Aha's rate limit headers from an API response.
//
// Aha allows up to 300 requests per minute and 20 requests per second per
// source IP. On every response it sends X-Ratelimit-Limit,
// X-Ratelimit-Remaining, and X-Ratelimit-Reset (the UTC unix time when the
// limit resets and it's safe to retry). See https://www.aha.io/api.
type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

// rateLimitFromHeaders parses Aha's rate limit headers from an HTTP response.
// It returns nil if none of the rate limit headers are present.
func rateLimitFromHeaders(h http.Header) *RateLimitInfo {
	limitStr := h.Get("X-Ratelimit-Limit")
	remainingStr := h.Get("X-Ratelimit-Remaining")
	resetStr := h.Get("X-Ratelimit-Reset")

	if limitStr == "" && remainingStr == "" && resetStr == "" {
		return nil
	}

	info := &RateLimitInfo{}
	if v, err := strconv.Atoi(limitStr); err == nil {
		info.Limit = v
	}
	if v, err := strconv.Atoi(remainingStr); err == nil {
		info.Remaining = v
	}
	if v, err := strconv.ParseInt(resetStr, 10, 64); err == nil {
		info.Reset = time.Unix(v, 0)
	}
	return info
}

// rateLimitedTransport is an http.RoundTripper that throttles outgoing
// requests to a fixed rate using a token bucket, waiting for a token before
// delegating to the next transport.
type rateLimitedTransport struct {
	next    http.RoundTripper
	limiter *rate.Limiter
}

// RoundTrip implements http.RoundTripper.
func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.next.RoundTrip(req)
}

// RateLimitFromError extracts rate limit information from an error returned
// by the client, if present. It returns nil if err is not an *APIError or
// carries no rate limit information (e.g. the response predates this SDK
// version capturing headers, or the request never reached Aha).
func RateLimitFromError(err error) *RateLimitInfo {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.RateLimit
	}
	return nil
}
