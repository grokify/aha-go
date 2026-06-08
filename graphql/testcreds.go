//go:build integration

package graphql

import (
	"fmt"
	"os"

	"github.com/grokify/goauth"
)

// TestCredentials holds credentials for integration tests.
type TestCredentials struct {
	Subdomain string
	APIKey    string
}

// LoadTestCredentials loads credentials for integration tests.
// It supports two patterns:
//
//  1. Direct env vars: AHA_SUBDOMAIN + AHA_API_KEY
//  2. Goauth file: GOAUTH_CREDENTIALS_FILE + GOAUTH_ACCOUNT
//
// Direct env vars take precedence if both are set.
func LoadTestCredentials() (*TestCredentials, error) {
	// Pattern 1: Direct env vars (preferred)
	subdomain := os.Getenv("AHA_SUBDOMAIN")
	apiKey := os.Getenv("AHA_API_KEY")
	if subdomain != "" && apiKey != "" {
		return &TestCredentials{
			Subdomain: subdomain,
			APIKey:    apiKey,
		}, nil
	}

	// Pattern 2: Goauth credentials file
	credsFile := os.Getenv("GOAUTH_CREDENTIALS_FILE")
	account := os.Getenv("GOAUTH_ACCOUNT")
	if credsFile != "" && account != "" {
		credSet, err := goauth.ReadFileCredentialsSet(credsFile, false)
		if err != nil {
			return nil, fmt.Errorf("loading goauth credentials file: %w", err)
		}
		creds, ok := credSet.Credentials[account]
		if !ok {
			return nil, fmt.Errorf("account %q not found in credentials file (available: %v)", account, credSet.Keys())
		}
		// Extract subdomain and API key from goauth credentials
		tc := &TestCredentials{
			Subdomain: creds.Subdomain,
		}
		// Check various places where access token might be stored
		if creds.Token != nil && creds.Token.AccessToken != "" {
			tc.APIKey = creds.Token.AccessToken
		} else if creds.OAuth2 != nil && creds.OAuth2.Token != nil && creds.OAuth2.Token.AccessToken != "" {
			tc.APIKey = creds.OAuth2.Token.AccessToken
		}
		if tc.Subdomain == "" || tc.APIKey == "" {
			return nil, fmt.Errorf("goauth credentials missing subdomain or access token for account %q", account)
		}
		return tc, nil
	}

	// No credentials found
	return nil, fmt.Errorf("no credentials found: set AHA_SUBDOMAIN+AHA_API_KEY or GOAUTH_CREDENTIALS_FILE+GOAUTH_ACCOUNT")
}

// SkipReason returns a message explaining why credentials are not available.
func SkipReason() string {
	return "credentials not set: need AHA_SUBDOMAIN+AHA_API_KEY or GOAUTH_CREDENTIALS_FILE+GOAUTH_ACCOUNT"
}
