package main

import (
	"fmt"
	"log"

	"onioncli/pkg/api"
	"onioncli/pkg/collections"
)

func main() {
	fmt.Println("OnionCLI Collections & Environments Demo")
	fmt.Println("=======================================")

	// Create collections manager
	manager, err := collections.NewManager()
	if err != nil {
		log.Fatalf("Failed to create collections manager: %v", err)
	}

	fmt.Printf("Collections directory: ~/.onioncli/collections/\n")
	fmt.Printf("Environments file: ~/.onioncli/environments.json\n\n")

	// Display current environments
	fmt.Println("Current Environments:")
	fmt.Println("====================")
	environments := manager.GetEnvironments()
	for _, env := range environments {
		status := ""
		if env.IsActive {
			status = " (Active)"
		}
		fmt.Printf("• %s%s - %s\n", env.Name, status, env.Description)
		for key, value := range env.Variables {
			fmt.Printf("  %s = %s\n", key, value)
		}
		fmt.Println()
	}

	// Create a new environment for testing
	fmt.Println("Creating test environments...")
	
	// Development environment
	devVars := map[string]string{
		"base_url":    "http://dev-api.example.onion:8080",
		"api_key":     "dev-key-123",
		"timeout":     "30",
		"debug_mode": "true",
	}
	devEnv := manager.CreateEnvironment("Development", "Development environment for .onion APIs", devVars)
	fmt.Printf("✅ Created environment: %s\n", devEnv.Name)

	// Production environment
	prodVars := map[string]string{
		"base_url":    "http://prod-api.example.onion",
		"api_key":     "prod-key-456",
		"timeout":     "60",
		"debug_mode": "false",
	}
	prodEnv := manager.CreateEnvironment("Production", "Production environment for .onion APIs", prodVars)
	fmt.Printf("✅ Created environment: %s\n", prodEnv.Name)

	// Test environment
	testVars := map[string]string{
		"base_url":    "http://test-api.example.onion:3000",
		"api_key":     "test-key-789",
		"timeout":     "15",
		"debug_mode": "true",
	}
	testEnv := manager.CreateEnvironment("Testing", "Testing environment for .onion APIs", testVars)
	fmt.Printf("✅ Created environment: %s\n", testEnv.Name)

	// Switch to development environment
	fmt.Println("\nSwitching to Development environment...")
	manager.SetActiveEnvironment(devEnv.ID)
	activeEnv := manager.GetActiveEnvironment()
	fmt.Printf("✅ Active environment: %s\n", activeEnv.Name)

	// Test variable substitution
	fmt.Println("\nTesting variable substitution...")
	testURL := "{{base_url}}/api/v1/users"
	substitutedURL := manager.SubstituteVariables(testURL)
	fmt.Printf("Original URL: %s\n", testURL)
	fmt.Printf("Substituted URL: %s\n", substitutedURL)

	testHeader := "Authorization: Bearer {{api_key}}"
	substitutedHeader := manager.SubstituteVariables(testHeader)
	fmt.Printf("Original Header: %s\n", testHeader)
	fmt.Printf("Substituted Header: %s\n", substitutedHeader)

	// Create a test collection
	fmt.Println("\nCreating test collection...")
	collection := manager.CreateCollection("Onion API Tests", "Collection of requests for testing .onion APIs")
	fmt.Printf("✅ Created collection: %s (ID: %s)\n", collection.Name, collection.ID)

	// Add some test requests to the collection
	fmt.Println("\nAdding test requests to collection...")

	// Request 1: Get users
	getUsersReq := api.NewRequest("GET", "{{base_url}}/api/v1/users")
	getUsersReq.SetHeader("Authorization", "Bearer {{api_key}}")
	getUsersReq.SetHeader("Accept", "application/json")
	err = manager.AddRequestToCollection(collection.ID, getUsersReq, "Get Users", "Retrieve list of all users")
	if err != nil {
		fmt.Printf("❌ Failed to add request: %v\n", err)
	} else {
		fmt.Printf("✅ Added request: Get Users\n")
	}

	// Request 2: Create user
	createUserReq := api.NewRequest("POST", "{{base_url}}/api/v1/users")
	createUserReq.SetHeader("Authorization", "Bearer {{api_key}}")
	createUserReq.SetHeader("Content-Type", "application/json")
	createUserReq.SetBody(`{
  "name": "Test User",
  "email": "test@example.com",
  "role": "user"
}`)
	err = manager.AddRequestToCollection(collection.ID, createUserReq, "Create User", "Create a new user account")
	if err != nil {
		fmt.Printf("❌ Failed to add request: %v\n", err)
	} else {
		fmt.Printf("✅ Added request: Create User\n")
	}

	// Request 3: Health check
	healthReq := api.NewRequest("GET", "{{base_url}}/health")
	healthReq.SetHeader("User-Agent", "OnionCLI/1.0")
	err = manager.AddRequestToCollection(collection.ID, healthReq, "Health Check", "Check API health status")
	if err != nil {
		fmt.Printf("❌ Failed to add request: %v\n", err)
	} else {
		fmt.Printf("✅ Added request: Health Check\n")
	}

	// Display collection summary
	fmt.Println("\nCollection Summary:")
	fmt.Println("==================")
	collections := manager.GetCollections()
	for _, col := range collections {
		fmt.Printf("Collection: %s\n", col.Name)
		fmt.Printf("Description: %s\n", col.Description)
		fmt.Printf("Requests: %d\n", len(col.Requests))
		fmt.Printf("Created: %s\n", col.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		for i, req := range col.Requests {
			fmt.Printf("  %d. %s %s - %s\n", i+1, req.Method, req.Name, req.Description)
			fmt.Printf("     URL: %s\n", req.URL)
			if len(req.Headers) > 0 {
				fmt.Printf("     Headers: %d\n", len(req.Headers))
			}
			if req.Body != "" {
				fmt.Printf("     Body: %d bytes\n", len(req.Body))
			}
			fmt.Println()
		}
	}

	// Test request processing with variable substitution
	fmt.Println("Testing request processing with variable substitution...")
	if len(collections) > 0 && len(collections[0].Requests) > 0 {
		originalReq := &collections[0].Requests[0]
		
		// Convert to API request
		apiReq := originalReq.ToRequest()
		fmt.Printf("Original request URL: %s\n", apiReq.URL)
		
		// Process with variable substitution
		processedReq := manager.ProcessRequest(apiReq)
		fmt.Printf("Processed request URL: %s\n", processedReq.URL)
		
		fmt.Printf("Original headers:\n")
		for key, value := range apiReq.Headers {
			fmt.Printf("  %s: %s\n", key, value)
		}
		
		fmt.Printf("Processed headers:\n")
		for key, value := range processedReq.Headers {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Test environment switching
	fmt.Println("\nTesting environment switching...")
	fmt.Printf("Current environment: %s\n", manager.GetActiveEnvironment().Name)

	// Switch to production
	manager.SetActiveEnvironment(prodEnv.ID)
	fmt.Printf("Switched to: %s\n", manager.GetActiveEnvironment().Name)

	// Test variable substitution with new environment
	newSubstitutedURL := manager.SubstituteVariables("{{base_url}}/api/v1/users")
	fmt.Printf("URL with production environment: %s\n", newSubstitutedURL)

	fmt.Println("\n🎉 Collections & Environments demo completed!")
	fmt.Println("\nFeatures demonstrated:")
	fmt.Println("• Environment management with variables")
	fmt.Println("• Variable substitution in URLs and headers")
	fmt.Println("• Request collections and organization")
	fmt.Println("• Persistent storage of collections and environments")
	fmt.Println("• Environment switching")
	fmt.Println("\nIn the TUI:")
	fmt.Println("• Press 'c' to browse collections")
	fmt.Println("• Press 'v' to manage environments")
	fmt.Println("• Use {{variable}} syntax in URLs and headers")
	fmt.Println("• Save requests to collections for reuse")
}
