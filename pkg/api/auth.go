package api

import (
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/zalando/go-keyring"
)

// AuthType represents different authentication methods
type AuthType string

const (
	AuthNone   AuthType = "none"
	AuthAPIKey AuthType = "api_key"
	AuthBearer AuthType = "bearer"
	AuthBasic  AuthType = "basic"
	AuthCustom AuthType = "custom"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Type     AuthType          `json:"type"`
	APIKey   string            `json:"api_key,omitempty"`
	KeyName  string            `json:"key_name,omitempty"`
	Location string            `json:"location,omitempty"` // "header" or "query"
	Token    string            `json:"token,omitempty"`
	Username string            `json:"username,omitempty"`
	Password string            `json:"password,omitempty"`
	Custom   map[string]string `json:"custom,omitempty"`
}

// AuthManager handles authentication for requests
type AuthManager struct {
	serviceName string
}

// NewAuthManager creates a new authentication manager
func NewAuthManager() *AuthManager {
	return &AuthManager{
		serviceName: "onioncli",
	}
}

// ApplyAuth applies authentication to a request based on the auth config
func (am *AuthManager) ApplyAuth(req *Request, config *AuthConfig) error {
	if config == nil || config.Type == AuthNone {
		return nil
	}

	switch config.Type {
	case AuthAPIKey:
		return am.applyAPIKeyAuth(req, config)
	case AuthBearer:
		return am.applyBearerAuth(req, config)
	case AuthBasic:
		return am.applyBasicAuth(req, config)
	case AuthCustom:
		return am.applyCustomAuth(req, config)
	default:
		return fmt.Errorf("unsupported authentication type: %s", config.Type)
	}
}

// applyAPIKeyAuth applies API key authentication
func (am *AuthManager) applyAPIKeyAuth(req *Request, config *AuthConfig) error {
	if config.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	keyName := config.KeyName
	if keyName == "" {
		keyName = "X-API-Key" // Default header name
	}

	switch config.Location {
	case "query":
		// Add to URL query parameters
		u, err := url.Parse(req.URL)
		if err != nil {
			return fmt.Errorf("invalid URL: %w", err)
		}

		query := u.Query()
		query.Set(keyName, config.APIKey)
		u.RawQuery = query.Encode()
		req.URL = u.String()

	case "header", "":
		// Add to headers (default)
		req.SetHeader(keyName, config.APIKey)

	default:
		return fmt.Errorf("invalid API key location: %s (use 'header' or 'query')", config.Location)
	}

	return nil
}

// applyBearerAuth applies Bearer token authentication
func (am *AuthManager) applyBearerAuth(req *Request, config *AuthConfig) error {
	if config.Token == "" {
		return fmt.Errorf("bearer token is required")
	}

	req.SetHeader("Authorization", fmt.Sprintf("Bearer %s", config.Token))
	return nil
}

// applyBasicAuth applies Basic authentication
func (am *AuthManager) applyBasicAuth(req *Request, config *AuthConfig) error {
	if config.Username == "" {
		return fmt.Errorf("username is required for basic auth")
	}

	// Password can be empty for some services
	credentials := fmt.Sprintf("%s:%s", config.Username, config.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.SetHeader("Authorization", fmt.Sprintf("Basic %s", encoded))
	return nil
}

// applyCustomAuth applies custom authentication headers
func (am *AuthManager) applyCustomAuth(req *Request, config *AuthConfig) error {
	if len(config.Custom) == 0 {
		return fmt.Errorf("custom headers are required")
	}

	for key, value := range config.Custom {
		req.SetHeader(key, value)
	}
	return nil
}

// StoreCredentials securely stores credentials using the system keyring
func (am *AuthManager) StoreCredentials(service, username, password string) error {
	return keyring.Set(am.serviceName+"-"+service, username, password)
}

// GetCredentials retrieves stored credentials from the system keyring
func (am *AuthManager) GetCredentials(service, username string) (string, error) {
	return keyring.Get(am.serviceName+"-"+service, username)
}

// DeleteCredentials removes stored credentials from the system keyring
func (am *AuthManager) DeleteCredentials(service, username string) error {
	return keyring.Delete(am.serviceName+"-"+service, username)
}

// ListStoredServices returns a list of services with stored credentials
func (am *AuthManager) ListStoredServices() ([]string, error) {
	// Note: go-keyring doesn't provide a list function, so we'll need to track this separately
	// For now, return an empty list - in a full implementation, we'd store a list of services
	return []string{}, nil
}

// ValidateAuthConfig validates an authentication configuration
func (am *AuthManager) ValidateAuthConfig(config *AuthConfig) error {
	if config == nil {
		return nil
	}

	switch config.Type {
	case AuthNone:
		return nil

	case AuthAPIKey:
		if config.APIKey == "" {
			return fmt.Errorf("API key is required")
		}
		if config.Location != "" && config.Location != "header" && config.Location != "query" {
			return fmt.Errorf("API key location must be 'header' or 'query'")
		}

	case AuthBearer:
		if config.Token == "" {
			return fmt.Errorf("bearer token is required")
		}

	case AuthBasic:
		if config.Username == "" {
			return fmt.Errorf("username is required for basic auth")
		}

	case AuthCustom:
		if len(config.Custom) == 0 {
			return fmt.Errorf("custom headers are required")
		}

	default:
		return fmt.Errorf("unsupported authentication type: %s", config.Type)
	}

	return nil
}

// GetAuthTypes returns all supported authentication types
func (am *AuthManager) GetAuthTypes() []AuthType {
	return []AuthType{
		AuthNone,
		AuthAPIKey,
		AuthBearer,
		AuthBasic,
		AuthCustom,
	}
}

// GetAuthTypeDescription returns a description for an auth type
func (am *AuthManager) GetAuthTypeDescription(authType AuthType) string {
	switch authType {
	case AuthNone:
		return "No authentication"
	case AuthAPIKey:
		return "API Key (header or query parameter)"
	case AuthBearer:
		return "Bearer Token (Authorization header)"
	case AuthBasic:
		return "Basic Authentication (username/password)"
	case AuthCustom:
		return "Custom headers"
	default:
		return "Unknown authentication type"
	}
}

// CreateAuthConfigFromInput creates an auth config from user input
func (am *AuthManager) CreateAuthConfigFromInput(authType AuthType, inputs map[string]string) (*AuthConfig, error) {
	config := &AuthConfig{Type: authType}

	switch authType {
	case AuthNone:
		return config, nil

	case AuthAPIKey:
		config.APIKey = inputs["api_key"]
		config.KeyName = inputs["key_name"]
		config.Location = inputs["location"]
		if config.Location == "" {
			config.Location = "header"
		}

	case AuthBearer:
		config.Token = inputs["token"]

	case AuthBasic:
		config.Username = inputs["username"]
		config.Password = inputs["password"]

	case AuthCustom:
		config.Custom = make(map[string]string)
		// Parse custom headers from input
		if headers := inputs["headers"]; headers != "" {
			lines := strings.Split(headers, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if key != "" && value != "" {
						config.Custom[key] = value
					}
				}
			}
		}

	default:
		return nil, fmt.Errorf("unsupported authentication type: %s", authType)
	}

	return config, am.ValidateAuthConfig(config)
}

// MaskSensitiveData masks sensitive information in auth config for display
func (am *AuthManager) MaskSensitiveData(config *AuthConfig) *AuthConfig {
	if config == nil {
		return nil
	}

	masked := *config // Copy the config

	// Mask sensitive fields
	if masked.APIKey != "" {
		masked.APIKey = am.maskString(masked.APIKey)
	}
	if masked.Token != "" {
		masked.Token = am.maskString(masked.Token)
	}
	if masked.Password != "" {
		masked.Password = "********"
	}

	// Mask custom headers that might contain sensitive data
	if len(masked.Custom) > 0 {
		maskedCustom := make(map[string]string)
		for key, value := range masked.Custom {
			if am.isSensitiveHeader(key) {
				maskedCustom[key] = am.maskString(value)
			} else {
				maskedCustom[key] = value
			}
		}
		masked.Custom = maskedCustom
	}

	return &masked
}

// maskString masks a string showing only first and last few characters
func (am *AuthManager) maskString(s string) string {
	if len(s) <= 8 {
		return "****"
	}
	return s[:3] + "****" + s[len(s)-3:]
}

// isSensitiveHeader checks if a header name typically contains sensitive data
func (am *AuthManager) isSensitiveHeader(headerName string) bool {
	sensitive := []string{
		"authorization", "x-api-key", "x-auth-token", "x-access-token",
		"api-key", "auth-token", "access-token", "secret", "password",
	}

	headerLower := strings.ToLower(headerName)
	for _, s := range sensitive {
		if strings.Contains(headerLower, s) {
			return true
		}
	}
	return false
}
