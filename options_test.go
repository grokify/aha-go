package aha

import (
	"net/http"
	"testing"
	"time"
)

func TestWithSubdomain(t *testing.T) {
	cfg := &Config{}
	opt := WithSubdomain("mycompany")
	opt(cfg)

	if cfg.Subdomain != "mycompany" {
		t.Errorf("Subdomain = %q, want %q", cfg.Subdomain, "mycompany")
	}
}

func TestWithAPIKey(t *testing.T) {
	cfg := &Config{}
	opt := WithAPIKey("secret123")
	opt(cfg)

	if cfg.APIKey != "secret123" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "secret123")
	}
}

func TestWithHTTPClient(t *testing.T) {
	customClient := &http.Client{Timeout: 5 * time.Second}
	cfg := &Config{}
	opt := WithHTTPClient(customClient)
	opt(cfg)

	if cfg.HTTPClient != customClient {
		t.Error("HTTPClient was not set correctly")
	}
}

func TestWithTimeout(t *testing.T) {
	cfg := &Config{}
	opt := WithTimeout(30 * time.Second)
	opt(cfg)

	if cfg.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 30*time.Second)
	}
}

func TestWithBaseURL(t *testing.T) {
	cfg := &Config{}
	opt := WithBaseURL("https://custom.example.com/api")
	opt(cfg)

	if cfg.BaseURL != "https://custom.example.com/api" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.example.com/api")
	}
}

func TestOptionChaining(t *testing.T) {
	cfg := &Config{}

	opts := []Option{
		WithSubdomain("mycompany"),
		WithAPIKey("key123"),
		WithTimeout(45 * time.Second),
		WithBaseURL("https://custom.aha.io/api"),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.Subdomain != "mycompany" {
		t.Errorf("Subdomain = %q, want %q", cfg.Subdomain, "mycompany")
	}
	if cfg.APIKey != "key123" {
		t.Errorf("APIKey = %q, want %q", cfg.APIKey, "key123")
	}
	if cfg.Timeout != 45*time.Second {
		t.Errorf("Timeout = %v, want %v", cfg.Timeout, 45*time.Second)
	}
	if cfg.BaseURL != "https://custom.aha.io/api" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.aha.io/api")
	}
}

func TestWithPage(t *testing.T) {
	opts := &ListOptions{}
	opt := WithPage(5)
	opt(opts)

	if opts.Page != 5 {
		t.Errorf("Page = %d, want %d", opts.Page, 5)
	}
}

func TestWithPerPage(t *testing.T) {
	opts := &ListOptions{}
	opt := WithPerPage(50)
	opt(opts)

	if opts.PerPage != 50 {
		t.Errorf("PerPage = %d, want %d", opts.PerPage, 50)
	}
}

func TestApplyListOptions(t *testing.T) {
	tests := []struct {
		name        string
		opts        []ListOption
		wantPage    int
		wantPerPage int
	}{
		{
			name:        "no options",
			opts:        nil,
			wantPage:    0,
			wantPerPage: 0,
		},
		{
			name:        "empty options",
			opts:        []ListOption{},
			wantPage:    0,
			wantPerPage: 0,
		},
		{
			name:        "page only",
			opts:        []ListOption{WithPage(3)},
			wantPage:    3,
			wantPerPage: 0,
		},
		{
			name:        "perPage only",
			opts:        []ListOption{WithPerPage(25)},
			wantPage:    0,
			wantPerPage: 25,
		},
		{
			name:        "both options",
			opts:        []ListOption{WithPage(2), WithPerPage(100)},
			wantPage:    2,
			wantPerPage: 100,
		},
		{
			name:        "later option overrides",
			opts:        []ListOption{WithPage(1), WithPage(10)},
			wantPage:    10,
			wantPerPage: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyListOptions(tt.opts...)

			if got.Page != tt.wantPage {
				t.Errorf("Page = %d, want %d", got.Page, tt.wantPage)
			}
			if got.PerPage != tt.wantPerPage {
				t.Errorf("PerPage = %d, want %d", got.PerPage, tt.wantPerPage)
			}
		})
	}
}

func TestListOptionsZeroValues(t *testing.T) {
	opts := applyListOptions()

	// Verify zero values
	if opts.Page != 0 {
		t.Errorf("Page = %d, want 0", opts.Page)
	}
	if opts.PerPage != 0 {
		t.Errorf("PerPage = %d, want 0", opts.PerPage)
	}
}
