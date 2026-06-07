// Package browser provides browser automation for Aha operations not available via API.
package browser

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

// Client wraps rod browser for Aha automation.
type Client struct {
	browser   *rod.Browser
	page      *rod.Page
	subdomain string
	email     string
	password  string
	headless  bool
	timeout   time.Duration
}

// Option configures the browser client.
type Option func(*Client)

// WithSubdomain sets the Aha subdomain.
func WithSubdomain(subdomain string) Option {
	return func(c *Client) {
		c.subdomain = subdomain
	}
}

// WithCredentials sets email and password for authentication.
func WithCredentials(email, password string) Option {
	return func(c *Client) {
		c.email = email
		c.password = password
	}
}

// WithHeadless sets whether to run browser in headless mode.
func WithHeadless(headless bool) Option {
	return func(c *Client) {
		c.headless = headless
	}
}

// WithTimeout sets the default timeout for operations.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
	}
}

// NewClient creates a new browser automation client.
// It reads configuration from environment variables if not provided via options:
//   - AHA_SUBDOMAIN: Aha account subdomain
//   - AHA_EMAIL: Aha login email
//   - AHA_PASSWORD: Aha login password
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		subdomain: os.Getenv("AHA_SUBDOMAIN"),
		email:     os.Getenv("AHA_EMAIL"),
		password:  os.Getenv("AHA_PASSWORD"),
		headless:  true,
		timeout:   30 * time.Second,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.subdomain == "" {
		return nil, fmt.Errorf("subdomain is required (set AHA_SUBDOMAIN or use WithSubdomain)")
	}
	if c.email == "" {
		return nil, fmt.Errorf("email is required (set AHA_EMAIL or use WithCredentials)")
	}
	if c.password == "" {
		return nil, fmt.Errorf("password is required (set AHA_PASSWORD or use WithCredentials)")
	}

	return c, nil
}

// BaseURL returns the Aha base URL for the configured subdomain.
func (c *Client) BaseURL() string {
	return fmt.Sprintf("https://%s.aha.io", c.subdomain)
}

// Connect launches the browser and establishes a connection.
func (c *Client) Connect(ctx context.Context) error {
	l := launcher.New().Headless(c.headless)
	url, err := l.Launch()
	if err != nil {
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	c.browser = rod.New().ControlURL(url)
	if err := c.browser.Connect(); err != nil {
		return fmt.Errorf("failed to connect to browser: %w", err)
	}

	return nil
}

// Close closes the browser connection.
func (c *Client) Close() error {
	if c.page != nil {
		_ = c.page.Close()
	}
	if c.browser != nil {
		return c.browser.Close()
	}
	return nil
}

// Login authenticates with Aha using email and password.
func (c *Client) Login(ctx context.Context) error {
	if c.browser == nil {
		return fmt.Errorf("browser not connected, call Connect first")
	}

	page := c.browser.MustPage()
	c.page = page

	// Navigate to login page
	loginURL := fmt.Sprintf("%s/users/sign_in", c.BaseURL())
	if err := page.Timeout(c.timeout).Navigate(loginURL); err != nil {
		return fmt.Errorf("failed to navigate to login page: %w", err)
	}

	// Wait for page to load
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page load: %w", err)
	}

	// Fill in email
	emailInput, err := page.Element("input[name='user[email]']")
	if err != nil {
		return fmt.Errorf("failed to find email input: %w", err)
	}
	if err := emailInput.Input(c.email); err != nil {
		return fmt.Errorf("failed to input email: %w", err)
	}

	// Fill in password
	passwordInput, err := page.Element("input[name='user[password]']")
	if err != nil {
		return fmt.Errorf("failed to find password input: %w", err)
	}
	if err := passwordInput.Input(c.password); err != nil {
		return fmt.Errorf("failed to input password: %w", err)
	}

	// Click submit button
	submitBtn, err := page.Element("input[type='submit'], button[type='submit']")
	if err != nil {
		return fmt.Errorf("failed to find submit button: %w", err)
	}
	if err := submitBtn.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click submit: %w", err)
	}

	// Wait for navigation to complete
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for login completion: %w", err)
	}

	// Verify login success by checking URL or element
	currentURL := page.MustInfo().URL
	if currentURL == loginURL {
		return fmt.Errorf("login failed, still on login page")
	}

	return nil
}

// Page returns the current page for advanced operations.
func (c *Client) Page() *rod.Page {
	return c.page
}

// Browser returns the underlying rod browser for advanced operations.
func (c *Client) Browser() *rod.Browser {
	return c.browser
}
