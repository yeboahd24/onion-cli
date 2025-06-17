package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"

	"onioncli/pkg/api"
)

func main() {
	fmt.Println("OnionCLI Error Handling Demo")
	fmt.Println("============================")

	// Create error analyzer
	analyzer := api.NewErrorAnalyzer()

	// Test different types of errors
	testErrors := []struct {
		name        string
		err         error
		url         string
		description string
	}{
		{
			name:        "Tor Connection Error",
			err:         errors.New("dial tcp 127.0.0.1:9050: connect: connection refused"),
			url:         "http://3g2upl4pq6kufc4m.onion",
			description: "Tor proxy not running",
		},
		{
			name:        "SOCKS Proxy Error",
			err:         errors.New("socks connect tcp 3g2upl4pq6kufc4m.onion:80: general socks server failure"),
			url:         "http://3g2upl4pq6kufc4m.onion",
			description: "Onion service unreachable",
		},
		{
			name:        "DNS Error",
			err:         errors.New("lookup nonexistent.example.com: no such host"),
			url:         "https://nonexistent.example.com",
			description: "Domain doesn't exist",
		},
		{
			name:        "Network Timeout",
			err:         &url.Error{Op: "Get", URL: "https://slow.example.com", Err: &net.OpError{Op: "dial", Err: errors.New("i/o timeout")}},
			url:         "https://slow.example.com",
			description: "Request timeout",
		},
		{
			name:        "Connection Refused",
			err:         errors.New("dial tcp 192.168.1.100:8080: connect: connection refused"),
			url:         "http://192.168.1.100:8080",
			description: "Service not running",
		},
		{
			name:        "Authentication Error",
			err:         errors.New("HTTP 401 Unauthorized: invalid credentials"),
			url:         "https://api.example.com/protected",
			description: "Invalid API key",
		},
		{
			name:        "Generic Error",
			err:         errors.New("something went wrong"),
			url:         "https://example.com",
			description: "Unknown error",
		},
	}

	fmt.Println("Analyzing different error types...\n")

	for i, test := range testErrors {
		fmt.Printf("%d. %s\n", i+1, test.name)
		fmt.Printf("   Description: %s\n", test.description)
		fmt.Printf("   Original Error: %v\n", test.err)
		fmt.Printf("   URL: %s\n", test.url)

		// Analyze the error
		diagnosticError := analyzer.AnalyzeError(test.err, test.url)
		if diagnosticError != nil {
			fmt.Printf("   Diagnosed Type: %s\n", diagnosticError.Type)
			fmt.Printf("   Diagnostic Message: %s\n", diagnosticError.Message)
			fmt.Printf("   Retryable: %v\n", diagnosticError.IsRetryable())

			if len(diagnosticError.Suggestions) > 0 {
				fmt.Printf("   Suggestions:\n")
				for j, suggestion := range diagnosticError.Suggestions {
					if j < 3 { // Show only first 3 suggestions for brevity
						fmt.Printf("     • %s\n", suggestion)
					}
				}
				if len(diagnosticError.Suggestions) > 3 {
					fmt.Printf("     ... and %d more suggestions\n", len(diagnosticError.Suggestions)-3)
				}
			}
		} else {
			fmt.Printf("   No diagnostic information available\n")
		}

		fmt.Println()
	}

	// Test diagnostic summary
	fmt.Println("Testing diagnostic summary format...")
	torError := analyzer.AnalyzeError(
		errors.New("dial tcp 127.0.0.1:9050: connect: connection refused"),
		"http://facebookwkhpilnemxj7asaniu7vnjjbiltxjqhye3mhbshg7kx5tfyd.onion",
	)

	if torError != nil {
		fmt.Println("Diagnostic Summary:")
		fmt.Println("==================")
		fmt.Print(torError.GetDiagnosticSummary())
	}

	// Test error type detection
	fmt.Println("\nTesting error type detection...")

	errorTests := []struct {
		err      error
		expected api.ErrorType
	}{
		{errors.New("socks connect failed"), api.ErrorTypeTor},
		{errors.New("connection refused"), api.ErrorTypeNetwork},
		{errors.New("i/o timeout"), api.ErrorTypeTimeout},
		{errors.New("no such host"), api.ErrorTypeDNS},
		{errors.New("401 unauthorized"), api.ErrorTypeAuth},
		{errors.New("random error"), api.ErrorTypeUnknown},
	}

	for _, test := range errorTests {
		diagnostic := analyzer.AnalyzeError(test.err, "https://example.com")
		if diagnostic != nil {
			status := "✅"
			if diagnostic.Type != test.expected {
				status = "❌"
			}
			fmt.Printf("%s Error: '%v' -> Detected: %s (Expected: %s)\n",
				status, test.err, diagnostic.Type, test.expected)
		}
	}

	fmt.Println("\n🎉 Error handling demo completed!")
	fmt.Println("The TUI will now show enhanced error messages with suggestions")
	fmt.Println("Press 'e' when an error occurs to see detailed diagnostic information")
}
