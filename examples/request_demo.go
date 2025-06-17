package main

import (
	"fmt"
	"log"

	"onioncli/pkg/api"
)

func main() {
	fmt.Println("OnionCLI Request Demo")
	fmt.Println("====================")

	// Create a client with default Tor configuration
	client, err := api.NewClient(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Test Tor connection first
	fmt.Println("Testing Tor connection...")
	if err := client.TestTorConnection(); err != nil {
		fmt.Printf("❌ Tor connection failed: %v\n", err)
		fmt.Println("Please ensure Tor is running on port 9050")
		return
	}
	fmt.Println("✅ Tor connection successful!")

	// Test with a regular HTTP service first (Tor disabled)
	fmt.Println("\nTesting with regular HTTP service (Tor disabled)...")
	client.SetTorEnabled(false)

	req := api.NewRequest("GET", "https://httpbin.org/json")
	req.SetHeader("User-Agent", "OnionCLI/1.0")

	// Send the request
	fmt.Println("⏳ Sending request...")
	resp, err := client.Send(req)
	if err != nil {
		fmt.Printf("❌ Request failed: %v\n", err)
		return
	}

	// Display response summary
	fmt.Printf("✅ Request completed successfully!\n")
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Printf("Response size: %d bytes\n", len(resp.Body))

	// Display some headers
	fmt.Println("\nResponse Headers:")
	for key, value := range resp.Headers {
		if key == "Content-Type" || key == "Server" || key == "Date" {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Display first 200 characters of response body
	if len(resp.Body) > 0 {
		fmt.Println("\nResponse Body (first 200 chars):")
		body := resp.Body
		if len(body) > 200 {
			body = body[:200] + "..."
		}
		fmt.Printf("%s\n", body)
	}

	// Test JSON pretty printing
	if prettyBody, err := resp.PrettyPrintJSON(); err == nil && prettyBody != resp.Body {
		fmt.Println("\nJSON Pretty Print Test:")
		fmt.Printf("Original: %s\n", resp.Body[:100])
		fmt.Printf("Pretty:   %s\n", prettyBody[:100])
	}

	fmt.Println("\n🎉 Request demo completed successfully!")
	fmt.Println("You can now use the interactive TUI by running: ./onioncli")
	fmt.Println("\nNote: For .onion services, ensure Tor is running and try the TUI interface.")
}
