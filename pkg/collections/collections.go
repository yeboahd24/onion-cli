package collections

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"onioncli/pkg/api"
)

// Collection represents a group of related requests
type Collection struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Requests    []CollectionRequest `json:"requests"`
	Variables   map[string]string   `json:"variables"`
	Auth        *api.AuthConfig     `json:"auth,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

// CollectionRequest represents a request within a collection
type CollectionRequest struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	Auth        *api.AuthConfig   `json:"auth,omitempty"`
	Tests       []string          `json:"tests,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
}

// Environment represents a set of variables for different contexts
type Environment struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Manager handles collections and environments
type Manager struct {
	collections    []Collection
	environments   []Environment
	activeEnv      *Environment
	collectionsDir string
	envFile        string
}

// NewManager creates a new collections manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".onioncli")
	collectionsDir := filepath.Join(configDir, "collections")
	envFile := filepath.Join(configDir, "environments.json")

	// Create directories
	if err := os.MkdirAll(collectionsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create collections directory: %w", err)
	}

	manager := &Manager{
		collections:    make([]Collection, 0),
		environments:   make([]Environment, 0),
		collectionsDir: collectionsDir,
		envFile:        envFile,
	}

	// Load existing data
	if err := manager.LoadCollections(); err != nil {
		return nil, fmt.Errorf("failed to load collections: %w", err)
	}

	if err := manager.LoadEnvironments(); err != nil {
		return nil, fmt.Errorf("failed to load environments: %w", err)
	}

	// Create default environment if none exist
	if len(manager.environments) == 0 {
		defaultEnv := Environment{
			ID:          generateID(),
			Name:        "Default",
			Description: "Default environment",
			Variables: map[string]string{
				"base_url": "http://localhost",
				"api_key":  "",
			},
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		manager.environments = append(manager.environments, defaultEnv)
		manager.activeEnv = &manager.environments[0]
		manager.SaveEnvironments()
	} else {
		// Find active environment
		for i := range manager.environments {
			if manager.environments[i].IsActive {
				manager.activeEnv = &manager.environments[i]
				break
			}
		}
	}

	return manager, nil
}

// CreateCollection creates a new collection
func (m *Manager) CreateCollection(name, description string) *Collection {
	collection := Collection{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Requests:    make([]CollectionRequest, 0),
		Variables:   make(map[string]string),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.collections = append(m.collections, collection)
	m.SaveCollection(&collection)
	return &m.collections[len(m.collections)-1]
}

// AddRequestToCollection adds a request to a collection
func (m *Manager) AddRequestToCollection(collectionID string, req *api.Request, name, description string) error {
	for i := range m.collections {
		if m.collections[i].ID == collectionID {
			collectionReq := CollectionRequest{
				ID:          generateID(),
				Name:        name,
				Description: description,
				Method:      req.Method,
				URL:         req.URL,
				Headers:     make(map[string]string),
				Body:        req.Body,
				CreatedAt:   time.Now(),
			}

			// Copy headers
			for k, v := range req.Headers {
				collectionReq.Headers[k] = v
			}

			m.collections[i].Requests = append(m.collections[i].Requests, collectionReq)
			m.collections[i].UpdatedAt = time.Now()
			return m.SaveCollection(&m.collections[i])
		}
	}
	return fmt.Errorf("collection not found: %s", collectionID)
}

// GetCollections returns all collections
func (m *Manager) GetCollections() []Collection {
	return m.collections
}

// GetCollection returns a specific collection
func (m *Manager) GetCollection(id string) (*Collection, error) {
	for i := range m.collections {
		if m.collections[i].ID == id {
			return &m.collections[i], nil
		}
	}
	return nil, fmt.Errorf("collection not found: %s", id)
}

// DeleteCollection deletes a collection
func (m *Manager) DeleteCollection(id string) error {
	for i, collection := range m.collections {
		if collection.ID == id {
			// Remove from slice
			m.collections = append(m.collections[:i], m.collections[i+1:]...)

			// Delete file
			filename := filepath.Join(m.collectionsDir, fmt.Sprintf("%s.json", id))
			return os.Remove(filename)
		}
	}
	return fmt.Errorf("collection not found: %s", id)
}

// CreateEnvironment creates a new environment
func (m *Manager) CreateEnvironment(name, description string, variables map[string]string) *Environment {
	env := Environment{
		ID:          generateID(),
		Name:        name,
		Description: description,
		Variables:   variables,
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	m.environments = append(m.environments, env)
	m.SaveEnvironments()
	return &m.environments[len(m.environments)-1]
}

// GetEnvironments returns all environments
func (m *Manager) GetEnvironments() []Environment {
	return m.environments
}

// GetActiveEnvironment returns the currently active environment
func (m *Manager) GetActiveEnvironment() *Environment {
	return m.activeEnv
}

// SetActiveEnvironment sets the active environment
func (m *Manager) SetActiveEnvironment(id string) error {
	// Deactivate all environments
	for i := range m.environments {
		m.environments[i].IsActive = false
	}

	// Activate the specified environment
	for i := range m.environments {
		if m.environments[i].ID == id {
			m.environments[i].IsActive = true
			m.activeEnv = &m.environments[i]
			return m.SaveEnvironments()
		}
	}

	return fmt.Errorf("environment not found: %s", id)
}

// SubstituteVariables replaces variables in a string with environment values
func (m *Manager) SubstituteVariables(input string) string {
	if m.activeEnv == nil {
		return input
	}

	result := input
	for key, value := range m.activeEnv.Variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, value)
	}

	return result
}

// ProcessRequest processes a request with variable substitution
func (m *Manager) ProcessRequest(req *api.Request) *api.Request {
	processedReq := &api.Request{
		Method:  req.Method,
		URL:     m.SubstituteVariables(req.URL),
		Headers: make(map[string]string),
		Body:    m.SubstituteVariables(req.Body),
	}

	// Process headers
	for key, value := range req.Headers {
		processedKey := m.SubstituteVariables(key)
		processedValue := m.SubstituteVariables(value)
		processedReq.Headers[processedKey] = processedValue
	}

	return processedReq
}

// LoadCollections loads all collections from disk
func (m *Manager) LoadCollections() error {
	files, err := filepath.Glob(filepath.Join(m.collectionsDir, "*.json"))
	if err != nil {
		return err
	}

	m.collections = make([]Collection, 0)
	for _, file := range files {
		var collection Collection
		data, err := os.ReadFile(file)
		if err != nil {
			continue // Skip corrupted files
		}

		if err := json.Unmarshal(data, &collection); err != nil {
			continue // Skip corrupted files
		}

		m.collections = append(m.collections, collection)
	}

	return nil
}

// SaveCollection saves a collection to disk
func (m *Manager) SaveCollection(collection *Collection) error {
	filename := filepath.Join(m.collectionsDir, fmt.Sprintf("%s.json", collection.ID))
	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadEnvironments loads environments from disk
func (m *Manager) LoadEnvironments() error {
	data, err := os.ReadFile(m.envFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet
		}
		return err
	}

	return json.Unmarshal(data, &m.environments)
}

// SaveEnvironments saves environments to disk
func (m *Manager) SaveEnvironments() error {
	data, err := json.MarshalIndent(m.environments, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.envFile, data, 0644)
}

// ToRequest converts a collection request to an API request
func (cr *CollectionRequest) ToRequest() *api.Request {
	req := api.NewRequest(cr.Method, cr.URL)

	// Set headers
	for k, v := range cr.Headers {
		req.SetHeader(k, v)
	}

	// Set body
	if cr.Body != "" {
		req.SetBody(cr.Body)
	}

	return req
}

// generateID generates a unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
