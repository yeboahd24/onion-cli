package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"onioncli/pkg/api"
)

// HistoryEntry represents a saved request with metadata
type HistoryEntry struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Timestamp   time.Time         `json:"timestamp"`
	Description string            `json:"description"`
}

// Manager handles request history persistence
type Manager struct {
	historyFile string
	entries     []HistoryEntry
}

// NewManager creates a new history manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".onioncli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	historyFile := filepath.Join(configDir, "history.json")

	manager := &Manager{
		historyFile: historyFile,
		entries:     make([]HistoryEntry, 0),
	}

	// Load existing history
	if err := manager.Load(); err != nil {
		// If file doesn't exist, that's okay - we'll create it on first save
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load history: %w", err)
		}
	}

	return manager, nil
}

// Save saves a request to history
func (m *Manager) Save(req *api.Request, name, description string) error {
	entry := HistoryEntry{
		ID:          generateID(),
		Name:        name,
		Method:      req.Method,
		URL:         req.URL,
		Headers:     make(map[string]string),
		Body:        req.Body,
		Timestamp:   time.Now(),
		Description: description,
	}

	// Copy headers
	for k, v := range req.Headers {
		entry.Headers[k] = v
	}

	// Add to entries (prepend to show most recent first)
	m.entries = append([]HistoryEntry{entry}, m.entries...)

	// Limit history size to 100 entries
	if len(m.entries) > 100 {
		m.entries = m.entries[:100]
	}

	return m.saveToFile()
}

// Load loads history from file
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.historyFile)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &m.entries)
}

// saveToFile saves history to file
func (m *Manager) saveToFile() error {
	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history: %w", err)
	}

	return os.WriteFile(m.historyFile, data, 0644)
}

// GetEntries returns all history entries
func (m *Manager) GetEntries() []HistoryEntry {
	return m.entries
}

// GetEntry returns a specific history entry by ID
func (m *Manager) GetEntry(id string) (*HistoryEntry, error) {
	for _, entry := range m.entries {
		if entry.ID == id {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("entry with ID %s not found", id)
}

// ToRequest converts a history entry back to an API request
func (entry *HistoryEntry) ToRequest() *api.Request {
	req := api.NewRequest(entry.Method, entry.URL)

	// Set headers
	for k, v := range entry.Headers {
		req.SetHeader(k, v)
	}

	// Set body
	if entry.Body != "" {
		req.SetBody(entry.Body)
	}

	return req
}

// Delete removes an entry from history
func (m *Manager) Delete(id string) error {
	for i, entry := range m.entries {
		if entry.ID == id {
			// Remove entry from slice
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return m.saveToFile()
		}
	}
	return fmt.Errorf("entry with ID %s not found", id)
}

// Clear removes all entries from history
func (m *Manager) Clear() error {
	m.entries = make([]HistoryEntry, 0)
	return m.saveToFile()
}

// Search searches history entries by name, URL, or description
func (m *Manager) Search(query string) []HistoryEntry {
	var results []HistoryEntry

	for _, entry := range m.entries {
		if contains(entry.Name, query) ||
			contains(entry.URL, query) ||
			contains(entry.Description, query) ||
			contains(entry.Method, query) {
			results = append(results, entry)
		}
	}

	return results
}

// GetRecentEntries returns the most recent N entries
func (m *Manager) GetRecentEntries(limit int) []HistoryEntry {
	if limit > len(m.entries) {
		limit = len(m.entries)
	}
	return m.entries[:limit]
}

// generateID generates a unique ID for history entries
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	if substr == "" {
		return true
	}
	// Simple case-insensitive search - convert both to lowercase
	sLower := ""
	substrLower := ""

	// Manual lowercase conversion to avoid importing strings
	for _, r := range s {
		if r >= 'A' && r <= 'Z' {
			sLower += string(r + 32)
		} else {
			sLower += string(r)
		}
	}

	for _, r := range substr {
		if r >= 'A' && r <= 'Z' {
			substrLower += string(r + 32)
		} else {
			substrLower += string(r)
		}
	}

	// Check if substring exists
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}

	return false
}

// Export exports history to a JSON file
func (m *Manager) Export(filename string) error {
	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal history for export: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// Import imports history from a JSON file
func (m *Manager) Import(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read import file: %w", err)
	}

	var importedEntries []HistoryEntry
	if err := json.Unmarshal(data, &importedEntries); err != nil {
		return fmt.Errorf("failed to unmarshal import data: %w", err)
	}

	// Merge with existing entries (imported entries go to the end)
	m.entries = append(m.entries, importedEntries...)

	// Limit total size
	if len(m.entries) > 100 {
		m.entries = m.entries[:100]
	}

	return m.saveToFile()
}

// GetStats returns statistics about the history
func (m *Manager) GetStats() map[string]interface{} {
	stats := make(map[string]interface{})

	stats["total_entries"] = len(m.entries)

	// Count by method
	methodCounts := make(map[string]int)
	for _, entry := range m.entries {
		methodCounts[entry.Method]++
	}
	stats["methods"] = methodCounts

	// Most recent entry
	if len(m.entries) > 0 {
		stats["most_recent"] = m.entries[0].Timestamp
	}

	return stats
}
