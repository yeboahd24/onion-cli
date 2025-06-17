package main

import (
	"fmt"
	"log"

	"onioncli/pkg/api"
)

func main() {
	fmt.Println("OnionCLI Authentication Demo")
	fmt.Println("===========================")

	// Create authentication manager
	authManager := api.NewAuthManager()

	// Test different authentication types
	fmt.Println("Testing authentication types...")

	// 1. API Key Authentication (Header)
	fmt.Println("\n1. API Key Authentication (Header)")
	apiKeyConfig := &api.AuthConfig{
		Type:     api.AuthAPIKey,
		APIKey:   "sk-1234567890abcdef",
		KeyName:  "X-API-Key",
		Location: "header",
	}

	req1 := api.NewRequest("GET", "https://api.example.com/data")
	err := authManager.ApplyAuth(req1, apiKeyConfig)
	if err != nil {
		log.Printf("Failed to apply API key auth: %v", err)
	} else {
		fmt.Printf("✅ Applied API key auth\n")
		fmt.Printf("   Headers: %v\n", req1.Headers)
	}

	// 2. API Key Authentication (Query)
	fmt.Println("\n2. API Key Authentication (Query)")
	apiKeyQueryConfig := &api.AuthConfig{
		Type:     api.AuthAPIKey,
		APIKey:   "abc123xyz789",
		KeyName:  "api_key",
		Location: "query",
	}

	req2 := api.NewRequest("GET", "https://api.example.com/search")
	err = authManager.ApplyAuth(req2, apiKeyQueryConfig)
	if err != nil {
		log.Printf("Failed to apply API key query auth: %v", err)
	} else {
		fmt.Printf("✅ Applied API key query auth\n")
		fmt.Printf("   URL: %s\n", req2.URL)
	}

	// 3. Bearer Token Authentication
	fmt.Println("\n3. Bearer Token Authentication")
	bearerConfig := &api.AuthConfig{
		Type:  api.AuthBearer,
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
	}

	req3 := api.NewRequest("POST", "https://api.github.com/user/repos")
	err = authManager.ApplyAuth(req3, bearerConfig)
	if err != nil {
		log.Printf("Failed to apply bearer auth: %v", err)
	} else {
		fmt.Printf("✅ Applied bearer token auth\n")
		fmt.Printf("   Authorization header: %s\n", req3.Headers["Authorization"])
	}

	// 4. Basic Authentication
	fmt.Println("\n4. Basic Authentication")
	basicConfig := &api.AuthConfig{
		Type:     api.AuthBasic,
		Username: "user123",
		Password: "secret456",
	}

	req4 := api.NewRequest("GET", "https://api.private.com/data")
	err = authManager.ApplyAuth(req4, basicConfig)
	if err != nil {
		log.Printf("Failed to apply basic auth: %v", err)
	} else {
		fmt.Printf("✅ Applied basic auth\n")
		fmt.Printf("   Authorization header: %s\n", req4.Headers["Authorization"])
	}

	// 5. Custom Headers Authentication
	fmt.Println("\n5. Custom Headers Authentication")
	customConfig := &api.AuthConfig{
		Type: api.AuthCustom,
		Custom: map[string]string{
			"X-Custom-Auth":   "custom-token-123",
			"X-Client-ID":     "client-456",
			"X-Request-ID":    "req-789",
		},
	}

	req5 := api.NewRequest("PUT", "https://api.custom.com/update")
	err = authManager.ApplyAuth(req5, customConfig)
	if err != nil {
		log.Printf("Failed to apply custom auth: %v", err)
	} else {
		fmt.Printf("✅ Applied custom headers auth\n")
		fmt.Printf("   Custom headers: %v\n", req5.Headers)
	}

	// Test validation
	fmt.Println("\n6. Testing Authentication Validation")
	
	// Valid config
	validConfig := &api.AuthConfig{
		Type:   api.AuthBearer,
		Token:  "valid-token",
	}
	if err := authManager.ValidateAuthConfig(validConfig); err != nil {
		fmt.Printf("❌ Valid config failed validation: %v\n", err)
	} else {
		fmt.Printf("✅ Valid config passed validation\n")
	}

	// Invalid config (missing token)
	invalidConfig := &api.AuthConfig{
		Type:  api.AuthBearer,
		Token: "",
	}
	if err := authManager.ValidateAuthConfig(invalidConfig); err != nil {
		fmt.Printf("✅ Invalid config correctly rejected: %v\n", err)
	} else {
		fmt.Printf("❌ Invalid config incorrectly accepted\n")
	}

	// Test masking sensitive data
	fmt.Println("\n7. Testing Sensitive Data Masking")
	sensitiveConfig := &api.AuthConfig{
		Type:     api.AuthAPIKey,
		APIKey:   "sk-1234567890abcdefghijklmnop",
		KeyName:  "Authorization",
		Location: "header",
	}

	maskedConfig := authManager.MaskSensitiveData(sensitiveConfig)
	fmt.Printf("Original API key: %s\n", sensitiveConfig.APIKey)
	fmt.Printf("Masked API key:   %s\n", maskedConfig.APIKey)

	// Test auth type descriptions
	fmt.Println("\n8. Authentication Type Descriptions")
	authTypes := authManager.GetAuthTypes()
	for _, authType := range authTypes {
		description := authManager.GetAuthTypeDescription(authType)
		fmt.Printf("  %s: %s\n", authType, description)
	}

	// Test creating config from input
	fmt.Println("\n9. Creating Config from Input")
	inputs := map[string]string{
		"api_key":   "test-key-123",
		"key_name":  "X-Test-Key",
		"location":  "header",
	}

	configFromInput, err := authManager.CreateAuthConfigFromInput(api.AuthAPIKey, inputs)
	if err != nil {
		fmt.Printf("❌ Failed to create config from input: %v\n", err)
	} else {
		fmt.Printf("✅ Created config from input\n")
		fmt.Printf("   Type: %s\n", configFromInput.Type)
		fmt.Printf("   API Key: %s\n", configFromInput.APIKey)
		fmt.Printf("   Key Name: %s\n", configFromInput.KeyName)
		fmt.Printf("   Location: %s\n", configFromInput.Location)
	}

	fmt.Println("\n🎉 Authentication demo completed!")
	fmt.Println("You can now use the interactive TUI with 'a' key to configure authentication")
}
