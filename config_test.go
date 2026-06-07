package aha

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestConfigLoadDefaults(t *testing.T) {
	tests := []struct {
		name           string
		config         Config
		wantHTTPClient bool
		wantTimeout    time.Duration
	}{
		{
			name:           "empty config gets defaults",
			config:         Config{},
			wantHTTPClient: true,
			wantTimeout:    60 * time.Second,
		},
		{
			name: "existing values preserved",
			config: Config{
				HTTPClient: &http.Client{Timeout: 10 * time.Second},
				Timeout:    30 * time.Second,
			},
			wantHTTPClient: true,
			wantTimeout:    30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.loadDefaults()

			if tt.wantHTTPClient && tt.config.HTTPClient == nil {
				t.Error("expected HTTPClient to be set")
			}
			if tt.config.Timeout != tt.wantTimeout {
				t.Errorf("Timeout = %v, want %v", tt.config.Timeout, tt.wantTimeout)
			}
		})
	}
}

func TestConfigLoadEnv(t *testing.T) {
	// Save original env vars
	origSubdomain := os.Getenv("AHA_SUBDOMAIN")
	origDomain := os.Getenv("AHA_DOMAIN")
	origAPIKey := os.Getenv("AHA_API_KEY")
	origAPIToken := os.Getenv("AHA_API_TOKEN")

	// Restore after test
	defer func() {
		os.Setenv("AHA_SUBDOMAIN", origSubdomain)
		os.Setenv("AHA_DOMAIN", origDomain)
		os.Setenv("AHA_API_KEY", origAPIKey)
		os.Setenv("AHA_API_TOKEN", origAPIToken)
	}()

	tests := []struct {
		name          string
		config        Config
		envVars       map[string]string
		wantSubdomain string
		wantAPIKey    string
	}{
		{
			name:   "loads AHA_SUBDOMAIN and AHA_API_KEY",
			config: Config{},
			envVars: map[string]string{
				"AHA_SUBDOMAIN": "mycompany",
				"AHA_API_KEY":   "secret123",
			},
			wantSubdomain: "mycompany",
			wantAPIKey:    "secret123",
		},
		{
			name:   "falls back to AHA_DOMAIN",
			config: Config{},
			envVars: map[string]string{
				"AHA_DOMAIN": "fallback",
			},
			wantSubdomain: "fallback",
		},
		{
			name:   "falls back to AHA_API_TOKEN",
			config: Config{},
			envVars: map[string]string{
				"AHA_API_TOKEN": "token456",
			},
			wantAPIKey: "token456",
		},
		{
			name: "existing values not overwritten",
			config: Config{
				Subdomain: "preset",
				APIKey:    "presetkey",
			},
			envVars: map[string]string{
				"AHA_SUBDOMAIN": "ignored",
				"AHA_API_KEY":   "ignored",
			},
			wantSubdomain: "preset",
			wantAPIKey:    "presetkey",
		},
		{
			name:   "AHA_SUBDOMAIN takes precedence over AHA_DOMAIN",
			config: Config{},
			envVars: map[string]string{
				"AHA_SUBDOMAIN": "primary",
				"AHA_DOMAIN":    "fallback",
			},
			wantSubdomain: "primary",
		},
		{
			name:   "AHA_API_KEY takes precedence over AHA_API_TOKEN",
			config: Config{},
			envVars: map[string]string{
				"AHA_API_KEY":   "primary",
				"AHA_API_TOKEN": "fallback",
			},
			wantAPIKey: "primary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env vars
			os.Unsetenv("AHA_SUBDOMAIN")
			os.Unsetenv("AHA_DOMAIN")
			os.Unsetenv("AHA_API_KEY")
			os.Unsetenv("AHA_API_TOKEN")

			// Set test env vars
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			tt.config.loadEnv()

			if tt.config.Subdomain != tt.wantSubdomain {
				t.Errorf("Subdomain = %q, want %q", tt.config.Subdomain, tt.wantSubdomain)
			}
			if tt.config.APIKey != tt.wantAPIKey {
				t.Errorf("APIKey = %q, want %q", tt.config.APIKey, tt.wantAPIKey)
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr error
	}{
		{
			name:    "empty config",
			config:  Config{},
			wantErr: ErrMissingSubdomain,
		},
		{
			name: "missing subdomain",
			config: Config{
				APIKey: "key123",
			},
			wantErr: ErrMissingSubdomain,
		},
		{
			name: "missing API key",
			config: Config{
				Subdomain: "mycompany",
			},
			wantErr: ErrMissingAPIKey,
		},
		{
			name: "valid config",
			config: Config{
				Subdomain: "mycompany",
				APIKey:    "key123",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if err != tt.wantErr {
				t.Errorf("validate() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfigBuildBaseURL(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantURL string
	}{
		{
			name: "default URL from subdomain",
			config: Config{
				Subdomain: "mycompany",
			},
			wantURL: "https://mycompany.aha.io/api/v1",
		},
		{
			name: "custom BaseURL overrides",
			config: Config{
				Subdomain: "mycompany",
				BaseURL:   "https://custom.example.com/api",
			},
			wantURL: "https://custom.example.com/api",
		},
		{
			name: "empty subdomain still builds URL",
			config: Config{
				Subdomain: "",
			},
			wantURL: "https://.aha.io/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.buildBaseURL()
			if got != tt.wantURL {
				t.Errorf("buildBaseURL() = %q, want %q", got, tt.wantURL)
			}
		})
	}
}
