package main

import (
	"fmt"
	"log"

	"onioncli/pkg/api"
	"onioncli/pkg/history"
)

func main() {
	fmt.Println("OnionCLI History Demo")
	fmt.Println("====================")

	// Create history manager
	manager, err := history.NewManager()
	if err != nil {
		log.Fatalf("Failed to create history manager: %v", err)
	}

	fmt.Printf("History file location: ~/.onioncli/history.json\n\n")

	// Create some test requests
	testRequests := []*api.Request{
		api.NewRequest("GET", "http://3g2upl4pq6kufc4m.onion"),
		api.NewRequest("POST", "https://httpbin.org/post"),
		api.NewRequest("GET", "https://api.github.com/users/octocat"),
	}

	// Add headers to requests
	testRequests[0].SetHeader("User-Agent", "OnionCLI/1.0")
	testRequests[1].SetHeader("Content-Type", "application/json")
	testRequests[1].SetBody(`{"test": "data"}`)
	testRequests[2].SetHeader("Accept", "application/vnd.github.v3+json")

	// Save requests to history
	names := []string{"DuckDuckGo Search", "HTTPBin POST Test", "GitHub API Test"}
	descriptions := []string{
		"Search via DuckDuckGo .onion service",
		"Test POST request with JSON body",
		"Fetch GitHub user information",
	}

	fmt.Println("Saving test requests to history...")
	for i, req := range testRequests {
		err := manager.Save(req, names[i], descriptions[i])
		if err != nil {
			fmt.Printf("❌ Failed to save request %d: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ Saved: %s\n", names[i])
		}
	}

	// Display history entries
	fmt.Println("\nHistory entries:")
	entries := manager.GetEntries()
	for i, entry := range entries {
		fmt.Printf("%d. %s (%s %s)\n", i+1, entry.Name, entry.Method, entry.URL)
		if entry.Description != "" {
			fmt.Printf("   Description: %s\n", entry.Description)
		}
		fmt.Printf("   Timestamp: %s\n", entry.Timestamp.Format("2006-01-02 15:04:05"))
		if len(entry.Headers) > 0 {
			fmt.Printf("   Headers: %d\n", len(entry.Headers))
		}
		if entry.Body != "" {
			fmt.Printf("   Body: %d bytes\n", len(entry.Body))
		}
		fmt.Println()
	}

	// Test search functionality
	fmt.Println("Testing search functionality...")
	searchResults := manager.Search("github")
	fmt.Printf("Search for 'github': %d results\n", len(searchResults))
	for _, result := range searchResults {
		fmt.Printf("  - %s\n", result.Name)
	}

	searchResults = manager.Search("POST")
	fmt.Printf("Search for 'POST': %d results\n", len(searchResults))
	for _, result := range searchResults {
		fmt.Printf("  - %s (%s)\n", result.Name, result.Method)
	}

	// Test converting back to request
	fmt.Println("\nTesting request reconstruction...")
	if len(entries) > 0 {
		entry := entries[0]
		reconstructed := entry.ToRequest()
		fmt.Printf("Original: %s %s\n", entry.Method, entry.URL)
		fmt.Printf("Reconstructed: %s %s\n", reconstructed.Method, reconstructed.URL)
		fmt.Printf("Headers match: %v\n", len(reconstructed.Headers) == len(entry.Headers))
		fmt.Printf("Body match: %v\n", reconstructed.Body == entry.Body)
	}

	// Display statistics
	fmt.Println("\nHistory statistics:")
	stats := manager.GetStats()
	fmt.Printf("Total entries: %v\n", stats["total_entries"])
	if methods, ok := stats["methods"].(map[string]int); ok {
		fmt.Println("Methods:")
		for method, count := range methods {
			fmt.Printf("  %s: %d\n", method, count)
		}
	}

	fmt.Println("\n🎉 History demo completed!")
	fmt.Println("You can now use the interactive TUI with 'h' key to browse history")
}
