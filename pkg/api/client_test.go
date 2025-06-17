package api

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.TorProxy != "127.0.0.1:9050" {
		t.Errorf("Expected TorProxy to be '127.0.0.1:9050', got '%s'", config.TorProxy)
	}
	
	if !config.TorEnabled {
		t.Error("Expected TorEnabled to be true")
	}
	
	if config.Timeout != 30*time.Second {
		t.Errorf("Expected Timeout to be 30s, got %v", config.Timeout)
	}
}

func TestIsOnionURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://3g2upl4pq6kufc4m.onion", true},
		{"https://facebookwkhpilnemxj7asaniu7vnjjbiltxjqhye3mhbshg7kx5tfyd.onion", true},
		{"http://google.com", false},
		{"https://example.com", false},
		{"invalid-url", false},
		{"http://invalid.onion", false}, // Too short
	}
	
	for _, test := range tests {
		result := IsOnionURL(test.url)
		if result != test.expected {
			t.Errorf("IsOnionURL(%s) = %v, expected %v", test.url, result, test.expected)
		}
	}
}

func TestValidateOnionURL(t *testing.T) {
	tests := []struct {
		url       string
		shouldErr bool
	}{
		{"http://3g2upl4pq6kufc4m.onion", false},
		{"https://facebookwkhpilnemxj7asaniu7vnjjbiltxjqhye3mhbshg7kx5tfyd.onion", false},
		{"ftp://3g2upl4pq6kufc4m.onion", true},  // Invalid scheme
		{"http://google.com", true},             // Not .onion
		{"invalid-url", true},                   // Invalid URL
	}
	
	for _, test := range tests {
		err := ValidateOnionURL(test.url)
		if test.shouldErr && err == nil {
			t.Errorf("ValidateOnionURL(%s) should have returned an error", test.url)
		}
		if !test.shouldErr && err != nil {
			t.Errorf("ValidateOnionURL(%s) should not have returned an error: %v", test.url, err)
		}
	}
}

func TestNewClient(t *testing.T) {
	// Test with default config
	client, err := NewClient(nil)
	if err != nil {
		t.Fatalf("Failed to create client with default config: %v", err)
	}
	
	if !client.IsTorEnabled() {
		t.Error("Expected Tor to be enabled by default")
	}
	
	// Test with custom config
	config := &ClientConfig{
		TorProxy:   "127.0.0.1:9051",
		TorEnabled: false,
		Timeout:    10 * time.Second,
	}
	
	client, err = NewClient(config)
	if err != nil {
		t.Fatalf("Failed to create client with custom config: %v", err)
	}
	
	if client.IsTorEnabled() {
		t.Error("Expected Tor to be disabled")
	}
}

func TestNewRequest(t *testing.T) {
	req := NewRequest("GET", "http://example.onion")
	
	if req.Method != "GET" {
		t.Errorf("Expected method to be 'GET', got '%s'", req.Method)
	}
	
	if req.URL != "http://example.onion" {
		t.Errorf("Expected URL to be 'http://example.onion', got '%s'", req.URL)
	}
	
	if req.Headers == nil {
		t.Error("Expected headers to be initialized")
	}
}

func TestRequestValidation(t *testing.T) {
	// Valid request
	req := NewRequest("GET", "http://example.onion")
	if err := req.Validate(); err != nil {
		t.Errorf("Valid request should not return error: %v", err)
	}
	
	// Missing URL
	req = NewRequest("GET", "")
	if err := req.Validate(); err == nil {
		t.Error("Request with missing URL should return error")
	}
	
	// Missing method
	req = &Request{URL: "http://example.onion"}
	if err := req.Validate(); err == nil {
		t.Error("Request with missing method should return error")
	}
	
	// Invalid JSON body
	req = NewRequest("POST", "http://example.onion")
	req.SetHeader("Content-Type", "application/json")
	req.SetBody("invalid json")
	if err := req.Validate(); err == nil {
		t.Error("Request with invalid JSON should return error")
	}
}
