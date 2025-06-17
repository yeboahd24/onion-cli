package main

import (
	"fmt"
	"log"

	"onioncli/pkg/api"
)

func main() {
	fmt.Println("OnionCLI Tor Integration Test")
	fmt.Println("=============================")

	// Create a client with default Tor configuration
	client, err := api.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	fmt.Printf("Tor enabled: %v\n", client.IsTorEnabled())

	// Test Tor connection
	fmt.Println("\nTesting Tor connection...")
	if err := client.TestTorConnection(); err != nil {
		fmt.Printf("❌ Tor connection failed: %v\n", err)
		fmt.Println("\nTo fix this:")
		fmt.Println("1. Install Tor: sudo apt install tor (Ubuntu/Debian) or brew install tor (macOS)")
		fmt.Println("2. Start Tor service: sudo systemctl start tor")
		fmt.Println("3. Verify Tor is running on port 9050: netstat -tlnp | grep 9050")
	} else {
		fmt.Println("✅ Tor connection successful!")
	}

	// Test .onion URL validation
	fmt.Println("\nTesting .onion URL validation...")
	testURLs := []string{
		"http://3g2upl4pq6kufc4m.onion",                                           // DuckDuckGo (v2)
		"https://facebookwkhpilnemxj7asaniu7vnjjbiltxjqhye3mhbshg7kx5tfyd.onion", // Facebook (v3)
		"http://google.com",                                                      // Regular URL
		"invalid-url",                                                            // Invalid URL
	}

	for _, url := range testURLs {
		isOnion := api.IsOnionURL(url)
		fmt.Printf("  %s -> .onion: %v", url, isOnion)
		
		if isOnion {
			if err := api.ValidateOnionURL(url); err != nil {
				fmt.Printf(" (validation failed: %v)", err)
			} else {
				fmt.Printf(" ✅")
			}
		}
		fmt.Println()
	}

	// Test creating requests
	fmt.Println("\nTesting request creation...")
	req := api.NewRequest("GET", "http://3g2upl4pq6kufc4m.onion")
	req.SetHeader("User-Agent", "OnionCLI/1.0")
	
	if err := req.Validate(); err != nil {
		fmt.Printf("❌ Request validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Request validation successful!")
		fmt.Printf("  Method: %s\n", req.Method)
		fmt.Printf("  URL: %s\n", req.URL)
		fmt.Printf("  Headers: %v\n", req.Headers)
	}

	fmt.Println("\n🎉 Tor integration setup complete!")
	fmt.Println("Next: Run the main CLI to start making requests to .onion services")
}
