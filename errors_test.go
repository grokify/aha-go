package aha

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestAPIErrorError(t *testing.T) {
	tests := []struct {
		name    string
		err     *APIError
		wantMsg string
	}{
		{
			name: "without request ID",
			err: &APIError{
				StatusCode: 404,
				Message:    "feature not found",
			},
			wantMsg: "aha: API error (status 404): feature not found",
		},
		{
			name: "with request ID",
			err: &APIError{
				StatusCode: 500,
				Message:    "internal server error",
				RequestID:  "req-123",
			},
			wantMsg: "aha: API error (status 500, request req-123): internal server error",
		},
		{
			name: "empty message",
			err: &APIError{
				StatusCode: 401,
				Message:    "",
			},
			wantMsg: "aha: API error (status 401): ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantMsg {
				t.Errorf("Error() = %q, want %q", got, tt.wantMsg)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "404 API error",
			err:  &APIError{StatusCode: http.StatusNotFound},
			want: true,
		},
		{
			name: "wrapped 404 API error",
			err:  fmt.Errorf("wrapped: %w", &APIError{StatusCode: http.StatusNotFound}),
			want: true,
		},
		{
			name: "different status code",
			err:  &APIError{StatusCode: http.StatusUnauthorized},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "401 API error",
			err:  &APIError{StatusCode: http.StatusUnauthorized},
			want: true,
		},
		{
			name: "wrapped 401 API error",
			err:  fmt.Errorf("wrapped: %w", &APIError{StatusCode: http.StatusUnauthorized}),
			want: true,
		},
		{
			name: "different status code",
			err:  &APIError{StatusCode: http.StatusForbidden},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnauthorized(tt.err); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsForbidden(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "403 API error",
			err:  &APIError{StatusCode: http.StatusForbidden},
			want: true,
		},
		{
			name: "wrapped 403 API error",
			err:  fmt.Errorf("wrapped: %w", &APIError{StatusCode: http.StatusForbidden}),
			want: true,
		},
		{
			name: "different status code",
			err:  &APIError{StatusCode: http.StatusUnauthorized},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsForbidden(tt.err); got != tt.want {
				t.Errorf("IsForbidden() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "429 API error",
			err:  &APIError{StatusCode: http.StatusTooManyRequests},
			want: true,
		},
		{
			name: "wrapped 429 API error",
			err:  fmt.Errorf("wrapped: %w", &APIError{StatusCode: http.StatusTooManyRequests}),
			want: true,
		},
		{
			name: "different status code",
			err:  &APIError{StatusCode: http.StatusBadRequest},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRateLimited(tt.err); got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsServerError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "500 Internal Server Error",
			err:  &APIError{StatusCode: http.StatusInternalServerError},
			want: true,
		},
		{
			name: "502 Bad Gateway",
			err:  &APIError{StatusCode: http.StatusBadGateway},
			want: true,
		},
		{
			name: "503 Service Unavailable",
			err:  &APIError{StatusCode: http.StatusServiceUnavailable},
			want: true,
		},
		{
			name: "wrapped 500 error",
			err:  fmt.Errorf("wrapped: %w", &APIError{StatusCode: 500}),
			want: true,
		},
		{
			name: "499 is not server error",
			err:  &APIError{StatusCode: 499},
			want: false,
		},
		{
			name: "600 is not server error",
			err:  &APIError{StatusCode: 600},
			want: false,
		},
		{
			name: "400 is not server error",
			err:  &APIError{StatusCode: http.StatusBadRequest},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsServerError(tt.err); got != tt.want {
				t.Errorf("IsServerError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	tests := []struct {
		name          string
		operation     string
		err           error
		wantNil       bool
		wantAPIError  bool
		wantStatus    int
		wantContains  string
	}{
		{
			name:      "nil error returns nil",
			operation: "get feature",
			err:       nil,
			wantNil:   true,
		},
		{
			name:         "404 in message",
			operation:    "get feature",
			err:          errors.New("response returned 404"),
			wantAPIError: true,
			wantStatus:   http.StatusNotFound,
			wantContains: "get feature",
		},
		{
			name:         "not found in message",
			operation:    "get feature",
			err:          errors.New("feature not found"),
			wantAPIError: true,
			wantStatus:   http.StatusNotFound,
			wantContains: "get feature",
		},
		{
			name:         "401 in message",
			operation:    "list products",
			err:          errors.New("status 401"),
			wantAPIError: true,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "unauthorized in message",
			operation:    "list products",
			err:          errors.New("unauthorized access"),
			wantAPIError: true,
			wantStatus:   http.StatusUnauthorized,
		},
		{
			name:         "403 in message",
			operation:    "delete feature",
			err:          errors.New("403 forbidden"),
			wantAPIError: true,
			wantStatus:   http.StatusForbidden,
		},
		{
			name:         "429 in message",
			operation:    "batch update",
			err:          errors.New("429 too many requests"),
			wantAPIError: true,
			wantStatus:   http.StatusTooManyRequests,
		},
		{
			name:         "rate limit in message",
			operation:    "batch update",
			err:          errors.New("rate limit exceeded"),
			wantAPIError: true,
			wantStatus:   http.StatusTooManyRequests,
		},
		{
			name:         "500 in message",
			operation:    "create feature",
			err:          errors.New("500 error"),
			wantAPIError: true,
			wantStatus:   http.StatusInternalServerError,
		},
		{
			name:         "internal server in message",
			operation:    "create feature",
			err:          errors.New("internal server error"),
			wantAPIError: true,
			wantStatus:   http.StatusInternalServerError,
		},
		{
			name:         "unknown error passes through",
			operation:    "connect",
			err:          errors.New("connection refused"),
			wantAPIError: false,
			wantContains: "connect: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wrapError(tt.operation, tt.err)

			if tt.wantNil {
				if got != nil {
					t.Errorf("wrapError() = %v, want nil", got)
				}
				return
			}

			if got == nil {
				t.Fatal("wrapError() = nil, want non-nil")
			}

			if tt.wantAPIError {
				var apiErr *APIError
				if !errors.As(got, &apiErr) {
					t.Errorf("wrapError() did not return *APIError, got %T", got)
					return
				}
				if apiErr.StatusCode != tt.wantStatus {
					t.Errorf("APIError.StatusCode = %d, want %d", apiErr.StatusCode, tt.wantStatus)
				}
			}

			if tt.wantContains != "" {
				if gotMsg := got.Error(); !contains(gotMsg, tt.wantContains) {
					t.Errorf("error message %q does not contain %q", gotMsg, tt.wantContains)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && searchSubstring(s, substr)))
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestSentinelErrors(t *testing.T) {
	// Verify sentinel errors have expected messages
	if ErrMissingSubdomain.Error() != "aha: subdomain is required" {
		t.Errorf("ErrMissingSubdomain = %q", ErrMissingSubdomain.Error())
	}
	if ErrMissingAPIKey.Error() != "aha: api key is required" {
		t.Errorf("ErrMissingAPIKey = %q", ErrMissingAPIKey.Error())
	}
}
