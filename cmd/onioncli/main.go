package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"onioncli/pkg/tui"
)

func main() {
	// Initialize the TUI model
	model, err := tui.NewModel()
	if err != nil {
		log.Fatalf("Failed to initialize TUI: %v", err)
	}

	// Initialize the Bubbletea program
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
