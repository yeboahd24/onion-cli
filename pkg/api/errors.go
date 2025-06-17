package api

import (
	"fmt"
	"net"
	"strings"
)

// ErrorType represents different categories of errors
type ErrorType string

const (
	ErrorTypeNetwork    ErrorType = "network"
	ErrorTypeTor        ErrorType = "tor"
	ErrorTypeAuth       ErrorType = "auth"
	ErrorTypeValidation ErrorType = "validation"
	ErrorTypeTimeout    ErrorType = "timeout"
	ErrorTypeDNS        ErrorType = "dns"
	ErrorTypeHTTP       ErrorType = "http"
	ErrorTypeUnknown    ErrorType = "unknown"
)

// DiagnosticError provides enhanced error information with suggestions
type DiagnosticError struct {
	Type        ErrorType `json:"type"`
	Message     string    `json:"message"`
	Cause       error     `json:"cause,omitempty"`
	Suggestions []string  `json:"suggestions"`
	URL         string    `json:"url,omitempty"`
	StatusCode  int       `json:"status_code,omitempty"`
}

// Error implements the error interface
func (de *DiagnosticError) Error() string {
	return de.Message
}

// Unwrap returns the underlying error
func (de *DiagnosticError) Unwrap() error {
	return de.Cause
}

// ErrorAnalyzer analyzes errors and provides diagnostic information
type ErrorAnalyzer struct{}

// NewErrorAnalyzer creates a new error analyzer
func NewErrorAnalyzer() *ErrorAnalyzer {
	return &ErrorAnalyzer{}
}

// AnalyzeError analyzes an error and returns a diagnostic error with suggestions
func (ea *ErrorAnalyzer) AnalyzeError(err error, requestURL string) *DiagnosticError {
	if err == nil {
		return nil
	}

	// Parse URL for context
	isOnion := IsOnionURL(requestURL)

	// Analyze different error types
	switch {
	case ea.isTorError(err):
		return ea.analyzeTorError(err, requestURL, isOnion)
	case ea.isNetworkError(err):
		return ea.analyzeNetworkError(err, requestURL, isOnion)
	case ea.isTimeoutError(err):
		return ea.analyzeTimeoutError(err, requestURL, isOnion)
	case ea.isDNSError(err):
		return ea.analyzeDNSError(err, requestURL, isOnion)
	case ea.isAuthError(err):
		return ea.analyzeAuthError(err, requestURL)
	default:
		return ea.analyzeGenericError(err, requestURL)
	}
}

// isTorError checks if the error is related to Tor
func (ea *ErrorAnalyzer) isTorError(err error) bool {
	errStr := strings.ToLower(err.Error())
	torKeywords := []string{
		"socks",
		"proxy",
		"tor",
		"general socks server failure",
		"connection refused",
		"127.0.0.1:9050",
	}

	for _, keyword := range torKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

// isNetworkError checks if the error is a network-related error
func (ea *ErrorAnalyzer) isNetworkError(err error) bool {
	if netErr, ok := err.(net.Error); ok {
		return netErr.Temporary() || netErr.Timeout()
	}

	errStr := strings.ToLower(err.Error())
	networkKeywords := []string{
		"connection refused",
		"connection reset",
		"network unreachable",
		"host unreachable",
		"no route to host",
	}

	for _, keyword := range networkKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

// isTimeoutError checks if the error is a timeout
func (ea *ErrorAnalyzer) isTimeoutError(err error) bool {
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}

	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "timeout") || strings.Contains(errStr, "deadline exceeded")
}

// isDNSError checks if the error is DNS-related
func (ea *ErrorAnalyzer) isDNSError(err error) bool {
	errStr := strings.ToLower(err.Error())
	dnsKeywords := []string{
		"no such host",
		"dns",
		"name resolution",
		"lookup",
	}

	for _, keyword := range dnsKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

// isAuthError checks if the error is authentication-related
func (ea *ErrorAnalyzer) isAuthError(err error) bool {
	errStr := strings.ToLower(err.Error())
	authKeywords := []string{
		"unauthorized",
		"authentication",
		"401",
		"403",
		"forbidden",
		"invalid credentials",
	}

	for _, keyword := range authKeywords {
		if strings.Contains(errStr, keyword) {
			return true
		}
	}
	return false
}

// analyzeTorError analyzes Tor-specific errors
func (ea *ErrorAnalyzer) analyzeTorError(err error, requestURL string, isOnion bool) *DiagnosticError {
	suggestions := []string{
		"Check if Tor is installed and running",
		"Verify Tor is listening on port 9050: netstat -tlnp | grep 9050",
		"Start Tor service: sudo systemctl start tor (Linux) or brew services start tor (macOS)",
		"Check Tor configuration in /etc/tor/torrc",
	}

	if strings.Contains(err.Error(), "connection refused") {
		suggestions = append(suggestions, "Tor proxy is not running or not accessible on 127.0.0.1:9050")
	}

	if strings.Contains(err.Error(), "general socks server failure") {
		suggestions = append(suggestions,
			"The .onion service might be down or unreachable",
			"Try a different .onion URL to test Tor connectivity",
		)
	}

	return &DiagnosticError{
		Type:        ErrorTypeTor,
		Message:     fmt.Sprintf("Tor connection failed: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// analyzeNetworkError analyzes network-related errors
func (ea *ErrorAnalyzer) analyzeNetworkError(err error, requestURL string, isOnion bool) *DiagnosticError {
	suggestions := []string{
		"Check your internet connection",
		"Verify the URL is correct and accessible",
	}

	if isOnion {
		suggestions = append(suggestions,
			"Ensure Tor is running and properly configured",
			"Try accessing a regular website to test connectivity",
		)
	} else {
		suggestions = append(suggestions,
			"Try accessing the URL in a web browser",
			"Check if the service is currently available",
		)
	}

	if strings.Contains(err.Error(), "connection refused") {
		suggestions = append(suggestions, "The server is not accepting connections on the specified port")
	}

	return &DiagnosticError{
		Type:        ErrorTypeNetwork,
		Message:     fmt.Sprintf("Network error: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// analyzeTimeoutError analyzes timeout errors
func (ea *ErrorAnalyzer) analyzeTimeoutError(err error, requestURL string, isOnion bool) *DiagnosticError {
	suggestions := []string{
		"Increase the request timeout in settings",
		"Check your internet connection speed",
	}

	if isOnion {
		suggestions = append(suggestions,
			"Tor requests typically take longer - consider increasing timeout to 60+ seconds",
			"The .onion service might be slow or overloaded",
			"Try the request again as Tor circuits can be slow",
		)
	} else {
		suggestions = append(suggestions,
			"The server might be overloaded or slow to respond",
			"Try the request again later",
		)
	}

	return &DiagnosticError{
		Type:        ErrorTypeTimeout,
		Message:     fmt.Sprintf("Request timeout: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// analyzeDNSError analyzes DNS-related errors
func (ea *ErrorAnalyzer) analyzeDNSError(err error, requestURL string, isOnion bool) *DiagnosticError {
	suggestions := []string{}

	if isOnion {
		suggestions = append(suggestions,
			"DNS errors for .onion URLs indicate a Tor configuration issue",
			"Ensure requests are routed through Tor proxy",
			"Check that Tor is running and properly configured",
		)
	} else {
		suggestions = append(suggestions,
			"Check if the domain name is spelled correctly",
			"Try using a different DNS server (8.8.8.8, 1.1.1.1)",
			"Check your network's DNS configuration",
		)
	}

	return &DiagnosticError{
		Type:        ErrorTypeDNS,
		Message:     fmt.Sprintf("DNS resolution failed: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// analyzeAuthError analyzes authentication-related errors
func (ea *ErrorAnalyzer) analyzeAuthError(err error, requestURL string) *DiagnosticError {
	suggestions := []string{
		"Check your authentication credentials",
		"Verify the authentication method is correct",
		"Ensure API keys or tokens are valid and not expired",
		"Check if the authentication headers are properly formatted",
	}

	return &DiagnosticError{
		Type:        ErrorTypeAuth,
		Message:     fmt.Sprintf("Authentication failed: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// analyzeGenericError analyzes generic errors
func (ea *ErrorAnalyzer) analyzeGenericError(err error, requestURL string) *DiagnosticError {
	suggestions := []string{
		"Check the error message for specific details",
		"Verify the request URL and parameters",
		"Try the request again",
	}

	return &DiagnosticError{
		Type:        ErrorTypeUnknown,
		Message:     fmt.Sprintf("Request failed: %v", err),
		Cause:       err,
		Suggestions: suggestions,
		URL:         requestURL,
	}
}

// GetDiagnosticSummary returns a formatted summary of the diagnostic error
func (de *DiagnosticError) GetDiagnosticSummary() string {
	var summary strings.Builder

	summary.WriteString(fmt.Sprintf("Error Type: %s\n", de.Type))
	summary.WriteString(fmt.Sprintf("Message: %s\n", de.Message))

	if de.URL != "" {
		summary.WriteString(fmt.Sprintf("URL: %s\n", de.URL))
	}

	if de.StatusCode != 0 {
		summary.WriteString(fmt.Sprintf("Status Code: %d\n", de.StatusCode))
	}

	if len(de.Suggestions) > 0 {
		summary.WriteString("\nSuggestions:\n")
		for i, suggestion := range de.Suggestions {
			summary.WriteString(fmt.Sprintf("  %d. %s\n", i+1, suggestion))
		}
	}

	return summary.String()
}

// IsRetryable returns true if the error might be resolved by retrying
func (de *DiagnosticError) IsRetryable() bool {
	switch de.Type {
	case ErrorTypeTimeout, ErrorTypeNetwork:
		return true
	case ErrorTypeTor:
		// Some Tor errors are retryable (circuit issues), others are not (Tor not running)
		return strings.Contains(strings.ToLower(de.Message), "circuit") ||
			strings.Contains(strings.ToLower(de.Message), "temporary")
	default:
		return false
	}
}
