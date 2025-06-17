package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"golang.org/x/net/proxy"
)

// Client represents an HTTP client with optional Tor proxy support
type Client struct {
	httpClient *http.Client
	torEnabled bool
	torProxy   string
	timeout    time.Duration
}

// ClientConfig holds configuration for the API client
type ClientConfig struct {
	TorProxy   string        // Tor SOCKS5 proxy address (default: 127.0.0.1:9050)
	TorEnabled bool          // Whether to route requests through Tor
	Timeout    time.Duration // Request timeout (default: 30s)
}

// DefaultConfig returns a default client configuration
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		TorProxy:   "127.0.0.1:9050",
		TorEnabled: true,
		Timeout:    30 * time.Second,
	}
}

// NewClient creates a new API client with the given configuration
func NewClient(config *ClientConfig) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	client := &Client{
		torEnabled: config.TorEnabled,
		torProxy:   config.TorProxy,
		timeout:    config.Timeout,
	}

	if config.TorEnabled {
		httpClient, err := createTorClient(config.TorProxy, config.Timeout)
		if err != nil {
			return nil, fmt.Errorf("failed to create Tor client: %w", err)
		}
		client.httpClient = httpClient
	} else {
		client.httpClient = &http.Client{
			Timeout: config.Timeout,
		}
	}

	return client, nil
}

// createTorClient creates an HTTP client configured to use Tor SOCKS5 proxy
func createTorClient(torProxy string, timeout time.Duration) (*http.Client, error) {
	// Create a SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	// Create a custom transport with the SOCKS5 dialer
	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
		DisableKeepAlives: true, // Recommended for Tor
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}, nil
}

// IsOnionURL checks if a URL is a .onion address
func IsOnionURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	// Regex pattern for .onion domains
	onionPattern := regexp.MustCompile(`^[a-z2-7]{16}\.onion$|^[a-z2-7]{56}\.onion$`)
	return onionPattern.MatchString(u.Host)
}

// ValidateOnionURL validates that a URL is a properly formatted .onion address
func ValidateOnionURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %w", err)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported scheme: %s (use http or https)", u.Scheme)
	}

	if !IsOnionURL(rawURL) {
		return fmt.Errorf("not a valid .onion address: %s", u.Host)
	}

	return nil
}

// TestTorConnection tests if Tor proxy is accessible
func (c *Client) TestTorConnection() error {
	if !c.torEnabled {
		return fmt.Errorf("Tor is not enabled")
	}

	// Try to connect to the Tor proxy
	conn, err := net.DialTimeout("tcp", c.torProxy, 5*time.Second)
	if err != nil {
		return fmt.Errorf("cannot connect to Tor proxy at %s: %w (is Tor running?)", c.torProxy, err)
	}
	conn.Close()

	return nil
}

// GetHTTPClient returns the underlying HTTP client
func (c *Client) GetHTTPClient() *http.Client {
	return c.httpClient
}

// IsTorEnabled returns whether Tor routing is enabled
func (c *Client) IsTorEnabled() bool {
	return c.torEnabled
}

// SetTorEnabled enables or disables Tor routing
func (c *Client) SetTorEnabled(enabled bool) error {
	if c.torEnabled == enabled {
		return nil // No change needed
	}

	c.torEnabled = enabled

	// Recreate the HTTP client with new settings
	config := &ClientConfig{
		TorProxy:   c.torProxy,
		TorEnabled: enabled,
		Timeout:    c.timeout,
	}

	newClient, err := NewClient(config)
	if err != nil {
		return err
	}

	c.httpClient = newClient.httpClient
	return nil
}
