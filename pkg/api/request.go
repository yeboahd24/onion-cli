package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request represents an HTTP request to be sent
type Request struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// Response represents an HTTP response received
type Response struct {
	StatusCode int               `json:"status_code"`
	Status     string            `json:"status"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
	Duration   time.Duration     `json:"duration"`
	Timestamp  time.Time         `json:"timestamp"`
}

// NewRequest creates a new API request
func NewRequest(method, url string) *Request {
	return &Request{
		Method:  strings.ToUpper(method),
		URL:     url,
		Headers: make(map[string]string),
		Body:    "",
	}
}

// SetHeader sets a header for the request
func (r *Request) SetHeader(key, value string) {
	r.Headers[key] = value
}

// SetBody sets the request body
func (r *Request) SetBody(body string) {
	r.Body = body
}

// SetJSONBody sets the request body as JSON and adds appropriate content-type header
func (r *Request) SetJSONBody(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	r.Body = string(jsonData)
	r.SetHeader("Content-Type", "application/json")
	return nil
}

// Validate validates the request
func (r *Request) Validate() error {
	if r.URL == "" {
		return fmt.Errorf("URL is required")
	}

	if r.Method == "" {
		return fmt.Errorf("HTTP method is required")
	}

	// Validate JSON body if Content-Type is application/json
	if contentType, exists := r.Headers["Content-Type"]; exists {
		if strings.Contains(contentType, "application/json") && r.Body != "" {
			var js json.RawMessage
			if err := json.Unmarshal([]byte(r.Body), &js); err != nil {
				return fmt.Errorf("invalid JSON body: %w", err)
			}
		}
	}

	return nil
}

// Send sends the HTTP request using the provided client
func (c *Client) Send(req *Request) (*Response, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("request validation failed: %w", err)
	}

	// Check if this is an .onion URL and Tor is disabled
	if IsOnionURL(req.URL) && !c.torEnabled {
		return nil, fmt.Errorf(".onion URLs require Tor to be enabled")
	}

	// Validate .onion URLs
	if IsOnionURL(req.URL) {
		if err := ValidateOnionURL(req.URL); err != nil {
			return nil, fmt.Errorf("invalid .onion URL: %w", err)
		}
	}

	startTime := time.Now()

	// Create HTTP request
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Send the request
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert response headers to map
	headers := make(map[string]string)
	for key, values := range httpResp.Header {
		if len(values) > 0 {
			headers[key] = values[0] // Take the first value if multiple exist
		}
	}

	duration := time.Since(startTime)

	response := &Response{
		StatusCode: httpResp.StatusCode,
		Status:     httpResp.Status,
		Headers:    headers,
		Body:       string(bodyBytes),
		Duration:   duration,
		Timestamp:  time.Now(),
	}

	return response, nil
}

// PrettyPrintJSON formats JSON response body for better readability
func (r *Response) PrettyPrintJSON() (string, error) {
	if r.Body == "" {
		return "", nil
	}

	// Check if the response is JSON
	contentType, exists := r.Headers["Content-Type"]
	if !exists || !strings.Contains(contentType, "application/json") {
		return r.Body, nil // Return as-is if not JSON
	}

	var jsonData interface{}
	if err := json.Unmarshal([]byte(r.Body), &jsonData); err != nil {
		return r.Body, nil // Return as-is if not valid JSON
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return r.Body, err
	}

	return string(prettyJSON), nil
}

// IsSuccess returns true if the response status code indicates success (2xx)
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if the response status code indicates client error (4xx)
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the response status code indicates server error (5xx)
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}
