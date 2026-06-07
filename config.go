package aha

import (
	"net/http"
	"os"
	"time"
)

// Config holds the configuration for the Aha client.
type Config struct {
	// Subdomain is your Aha account subdomain (e.g., "mycompany" for mycompany.aha.io).
	Subdomain string

	// APIKey is your Aha API key for authentication.
	APIKey string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, http.DefaultClient is used.
	HTTPClient *http.Client

	// Timeout is the request timeout.
	// Default is 60 seconds.
	Timeout time.Duration

	// BaseURL overrides the default API URL.
	// If empty, https://{subdomain}.aha.io/api/v1 is used.
	BaseURL string
}

// loadDefaults sets default values for unset fields.
func (c *Config) loadDefaults() {
	if c.HTTPClient == nil {
		c.HTTPClient = http.DefaultClient
	}
	if c.Timeout == 0 {
		c.Timeout = 60 * time.Second
	}
}

// loadEnv loads configuration from environment variables.
// Environment variables are only used if the corresponding field is empty.
func (c *Config) loadEnv() {
	if c.Subdomain == "" {
		if v := os.Getenv("AHA_SUBDOMAIN"); v != "" {
			c.Subdomain = v
		} else if v := os.Getenv("AHA_DOMAIN"); v != "" {
			// Also support AHA_DOMAIN for backward compatibility
			c.Subdomain = v
		}
	}
	if c.APIKey == "" {
		if v := os.Getenv("AHA_API_KEY"); v != "" {
			c.APIKey = v
		} else if v := os.Getenv("AHA_API_TOKEN"); v != "" {
			// Also support AHA_API_TOKEN for backward compatibility
			c.APIKey = v
		}
	}
}

// validate checks that required fields are set.
func (c *Config) validate() error {
	if c.Subdomain == "" {
		return ErrMissingSubdomain
	}
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	return nil
}

// buildBaseURL returns the API base URL.
func (c *Config) buildBaseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return "https://" + c.Subdomain + ".aha.io/api/v1"
}
