package aha

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// Sentinel errors for configuration and validation.
var (
	ErrMissingSubdomain = errors.New("aha: subdomain is required")
	ErrMissingAPIKey    = errors.New("aha: api key is required")
)

// APIError represents an error response from the Aha API.
type APIError struct {
	StatusCode int
	Message    string
	RequestID  string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.RequestID != "" {
		return fmt.Sprintf("aha: API error (status %d, request %s): %s",
			e.StatusCode, e.RequestID, e.Message)
	}
	return fmt.Sprintf("aha: API error (status %d): %s", e.StatusCode, e.Message)
}

// IsNotFound returns true if the error indicates a 404 Not Found response.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized returns true if the error indicates a 401 Unauthorized response.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusUnauthorized
	}
	return false
}

// IsForbidden returns true if the error indicates a 403 Forbidden response.
func IsForbidden(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusForbidden
	}
	return false
}

// IsRateLimited returns true if the error indicates a 429 Too Many Requests response.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests
	}
	return false
}

// IsServerError returns true if the error indicates a 5xx server error.
func IsServerError(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode >= 500 && apiErr.StatusCode < 600
	}
	return false
}

// wrapError wraps an error from the API client with operation context.
func wrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	// Extract HTTP status code from error message if present
	// ogen errors often include status code in the message
	errMsg := err.Error()

	// Check for common HTTP error patterns
	statusCode := 0
	switch {
	case strings.Contains(errMsg, "404") || strings.Contains(errMsg, "not found"):
		statusCode = http.StatusNotFound
	case strings.Contains(errMsg, "401") || strings.Contains(errMsg, "unauthorized"):
		statusCode = http.StatusUnauthorized
	case strings.Contains(errMsg, "403") || strings.Contains(errMsg, "forbidden"):
		statusCode = http.StatusForbidden
	case strings.Contains(errMsg, "429") || strings.Contains(errMsg, "rate limit"):
		statusCode = http.StatusTooManyRequests
	case strings.Contains(errMsg, "500") || strings.Contains(errMsg, "internal server"):
		statusCode = http.StatusInternalServerError
	}

	if statusCode > 0 {
		return &APIError{
			StatusCode: statusCode,
			Message:    fmt.Sprintf("%s: %v", operation, err),
		}
	}

	// Return as-is if we can't determine the status code
	return fmt.Errorf("%s: %w", operation, err)
}
