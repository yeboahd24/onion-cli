package main

import (
	"fmt"
	"log"
	"os"

	"onioncli/pkg/config"
)

func main() {
	fmt.Println("OnionCLI Configuration Demo")
	fmt.Println("===========================")

	// Create configuration manager
	manager, err := config.NewManager()
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	fmt.Printf("Configuration file location: %s\n\n", manager.GetConfigPath())

	// Display current configuration
	cfg := manager.Get()
	fmt.Println("Current Configuration:")
	fmt.Println("=====================")

	// Tor settings
	fmt.Printf("Tor Settings:\n")
	fmt.Printf("  Enabled: %v\n", cfg.Tor.Enabled)
	fmt.Printf("  Proxy Address: %s:%d\n", cfg.Tor.ProxyAddr, cfg.Tor.ProxyPort)
	fmt.Printf("  Timeout: %d seconds\n", cfg.Tor.Timeout)
	fmt.Printf("  Auto Detect: %v\n", cfg.Tor.AutoDetect)
	fmt.Printf("  Full Proxy Address: %s\n", manager.GetTorProxyAddress())
	fmt.Printf("  Timeout Duration: %v\n", manager.GetTorTimeout())
	fmt.Println()

	// HTTP settings
	fmt.Printf("HTTP Settings:\n")
	fmt.Printf("  Timeout: %d seconds (%v)\n", cfg.HTTP.Timeout, manager.GetHTTPTimeout())
	fmt.Printf("  Follow Redirects: %v\n", cfg.HTTP.FollowRedirects)
	fmt.Printf("  Max Redirects: %d\n", cfg.HTTP.MaxRedirects)
	fmt.Printf("  Verify SSL: %v\n", cfg.HTTP.VerifySSL)
	fmt.Printf("  User Agent: %s\n", cfg.HTTP.UserAgent)
	fmt.Println()

	// UI settings
	fmt.Printf("UI Settings:\n")
	fmt.Printf("  Theme: %s\n", cfg.UI.Theme)
	fmt.Printf("  Show Line Numbers: %v\n", cfg.UI.ShowLineNumbers)
	fmt.Printf("  Auto Save: %v\n", cfg.UI.AutoSave)
	fmt.Printf("  Confirm Exit: %v\n", cfg.UI.ConfirmExit)
	fmt.Println()

	// History settings
	fmt.Printf("History Settings:\n")
	fmt.Printf("  Enabled: %v\n", cfg.History.Enabled)
	fmt.Printf("  Max Entries: %d\n", cfg.History.MaxEntries)
	fmt.Printf("  Auto Save: %v\n", cfg.History.AutoSave)
	fmt.Println()

	// Default headers
	fmt.Printf("Default Headers:\n")
	for key, value := range cfg.DefaultHeaders {
		fmt.Printf("  %s: %s\n", key, value)
	}
	fmt.Println()

	// Test configuration updates
	fmt.Println("Testing configuration updates...")

	// Update Tor settings
	fmt.Println("1. Updating Tor settings...")
	manager.UpdateTorSettings(true, "127.0.0.1", 9051, 45)
	fmt.Printf("   New Tor proxy: %s\n", manager.GetTorProxyAddress())
	fmt.Printf("   New Tor timeout: %v\n", manager.GetTorTimeout())

	// Update HTTP settings
	fmt.Println("2. Updating HTTP settings...")
	manager.UpdateHTTPSettings(60, false, 5, false, "OnionCLI/2.0")
	updatedCfg := manager.Get()
	fmt.Printf("   New HTTP timeout: %v\n", manager.GetHTTPTimeout())
	fmt.Printf("   New User Agent: %s\n", updatedCfg.HTTP.UserAgent)
	fmt.Printf("   Follow Redirects: %v\n", updatedCfg.HTTP.FollowRedirects)

	// Update UI settings
	fmt.Println("3. Updating UI settings...")
	manager.UpdateUISettings("light", false, false, true)
	updatedCfg = manager.Get()
	fmt.Printf("   New Theme: %s\n", updatedCfg.UI.Theme)
	fmt.Printf("   Show Line Numbers: %v\n", updatedCfg.UI.ShowLineNumbers)
	fmt.Printf("   Confirm Exit: %v\n", updatedCfg.UI.ConfirmExit)

	// Test default headers management
	fmt.Println("4. Testing default headers management...")
	manager.AddDefaultHeader("X-Custom-Header", "CustomValue")
	manager.AddDefaultHeader("Authorization", "Bearer token123")

	headers := manager.GetDefaultHeaders()
	fmt.Printf("   Updated headers (%d total):\n", len(headers))
	for key, value := range headers {
		fmt.Printf("     %s: %s\n", key, value)
	}

	// Remove a header
	manager.RemoveDefaultHeader("Accept")
	headers = manager.GetDefaultHeaders()
	fmt.Printf("   After removing 'Accept' header (%d total):\n", len(headers))
	for key, value := range headers {
		fmt.Printf("     %s: %s\n", key, value)
	}

	// Test validation
	fmt.Println("5. Testing configuration validation...")
	if err := manager.Validate(); err != nil {
		fmt.Printf("   ❌ Validation failed: %v\n", err)
	} else {
		fmt.Printf("   ✅ Configuration is valid\n")
	}

	// Test invalid configuration
	fmt.Println("6. Testing invalid configuration...")
	manager.UpdateTorSettings(true, "127.0.0.1", 99999, 45) // Invalid port
	if err := manager.Validate(); err != nil {
		fmt.Printf("   ✅ Correctly detected invalid config: %v\n", err)
	} else {
		fmt.Printf("   ❌ Failed to detect invalid configuration\n")
	}

	// Reset to valid configuration
	manager.UpdateTorSettings(true, "127.0.0.1", 9050, 30)

	// Save configuration
	fmt.Println("7. Saving configuration...")
	if err := manager.Save(); err != nil {
		fmt.Printf("   ❌ Failed to save config: %v\n", err)
	} else {
		fmt.Printf("   ✅ Configuration saved successfully\n")
	}

	// Test export/import
	fmt.Println("8. Testing export/import...")
	exportFile := "/tmp/onioncli_config_export.yaml"

	if err := manager.Export(exportFile); err != nil {
		fmt.Printf("   ❌ Failed to export config: %v\n", err)
	} else {
		fmt.Printf("   ✅ Configuration exported to %s\n", exportFile)

		// Check if file exists
		if _, err := os.Stat(exportFile); err == nil {
			fmt.Printf("   ✅ Export file exists and is readable\n")

			// Clean up
			os.Remove(exportFile)
			fmt.Printf("   🧹 Cleaned up export file\n")
		}
	}

	// Test reset
	fmt.Println("9. Testing configuration reset...")
	if err := manager.Reset(); err != nil {
		fmt.Printf("   ❌ Failed to reset config: %v\n", err)
	} else {
		fmt.Printf("   ✅ Configuration reset to defaults\n")
		resetCfg := manager.Get()
		fmt.Printf("   Tor proxy after reset: %s\n", manager.GetTorProxyAddress())
		fmt.Printf("   HTTP timeout after reset: %v\n", manager.GetHTTPTimeout())
		fmt.Printf("   Theme after reset: %s\n", resetCfg.UI.Theme)
	}

	fmt.Println("\n🎉 Configuration demo completed!")
	fmt.Println("Configuration is now managed and persisted automatically")
	fmt.Println("Settings can be modified through the TUI interface")
}
