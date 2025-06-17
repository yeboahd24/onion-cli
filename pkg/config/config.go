package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Tor settings
	Tor TorConfig `mapstructure:"tor" json:"tor"`
	
	// HTTP settings
	HTTP HTTPConfig `mapstructure:"http" json:"http"`
	
	// UI settings
	UI UIConfig `mapstructure:"ui" json:"ui"`
	
	// Default headers
	DefaultHeaders map[string]string `mapstructure:"default_headers" json:"default_headers"`
	
	// History settings
	History HistoryConfig `mapstructure:"history" json:"history"`
}

// TorConfig holds Tor-specific configuration
type TorConfig struct {
	Enabled     bool   `mapstructure:"enabled" json:"enabled"`
	ProxyAddr   string `mapstructure:"proxy_addr" json:"proxy_addr"`
	ProxyPort   int    `mapstructure:"proxy_port" json:"proxy_port"`
	Timeout     int    `mapstructure:"timeout" json:"timeout"` // seconds
	AutoDetect  bool   `mapstructure:"auto_detect" json:"auto_detect"`
}

// HTTPConfig holds HTTP-specific configuration
type HTTPConfig struct {
	Timeout         int  `mapstructure:"timeout" json:"timeout"`         // seconds
	FollowRedirects bool `mapstructure:"follow_redirects" json:"follow_redirects"`
	MaxRedirects    int  `mapstructure:"max_redirects" json:"max_redirects"`
	VerifySSL       bool `mapstructure:"verify_ssl" json:"verify_ssl"`
	UserAgent       string `mapstructure:"user_agent" json:"user_agent"`
}

// UIConfig holds UI-specific configuration
type UIConfig struct {
	Theme           string `mapstructure:"theme" json:"theme"`
	ShowLineNumbers bool   `mapstructure:"show_line_numbers" json:"show_line_numbers"`
	AutoSave        bool   `mapstructure:"auto_save" json:"auto_save"`
	ConfirmExit     bool   `mapstructure:"confirm_exit" json:"confirm_exit"`
}

// HistoryConfig holds history-specific configuration
type HistoryConfig struct {
	Enabled    bool `mapstructure:"enabled" json:"enabled"`
	MaxEntries int  `mapstructure:"max_entries" json:"max_entries"`
	AutoSave   bool `mapstructure:"auto_save" json:"auto_save"`
}

// Manager handles configuration loading, saving, and management
type Manager struct {
	config     *Config
	configPath string
	viper      *viper.Viper
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".onioncli")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.yaml")

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	manager := &Manager{
		configPath: configPath,
		viper:      v,
	}

	// Set defaults
	manager.setDefaults()

	// Load existing config or create default
	if err := manager.Load(); err != nil {
		if os.IsNotExist(err) {
			// Create default config
			manager.config = manager.getDefaultConfig()
			if err := manager.Save(); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	return manager, nil
}

// setDefaults sets default values in viper
func (m *Manager) setDefaults() {
	// Tor defaults
	m.viper.SetDefault("tor.enabled", true)
	m.viper.SetDefault("tor.proxy_addr", "127.0.0.1")
	m.viper.SetDefault("tor.proxy_port", 9050)
	m.viper.SetDefault("tor.timeout", 30)
	m.viper.SetDefault("tor.auto_detect", true)

	// HTTP defaults
	m.viper.SetDefault("http.timeout", 30)
	m.viper.SetDefault("http.follow_redirects", true)
	m.viper.SetDefault("http.max_redirects", 10)
	m.viper.SetDefault("http.verify_ssl", true)
	m.viper.SetDefault("http.user_agent", "OnionCLI/1.0")

	// UI defaults
	m.viper.SetDefault("ui.theme", "dark")
	m.viper.SetDefault("ui.show_line_numbers", true)
	m.viper.SetDefault("ui.auto_save", true)
	m.viper.SetDefault("ui.confirm_exit", false)

	// History defaults
	m.viper.SetDefault("history.enabled", true)
	m.viper.SetDefault("history.max_entries", 100)
	m.viper.SetDefault("history.auto_save", true)

	// Default headers
	m.viper.SetDefault("default_headers", map[string]string{
		"User-Agent": "OnionCLI/1.0",
		"Accept":     "application/json, text/plain, */*",
	})
}

// getDefaultConfig returns the default configuration
func (m *Manager) getDefaultConfig() *Config {
	return &Config{
		Tor: TorConfig{
			Enabled:    true,
			ProxyAddr:  "127.0.0.1",
			ProxyPort:  9050,
			Timeout:    30,
			AutoDetect: true,
		},
		HTTP: HTTPConfig{
			Timeout:         30,
			FollowRedirects: true,
			MaxRedirects:    10,
			VerifySSL:       true,
			UserAgent:       "OnionCLI/1.0",
		},
		UI: UIConfig{
			Theme:           "dark",
			ShowLineNumbers: true,
			AutoSave:        true,
			ConfirmExit:     false,
		},
		DefaultHeaders: map[string]string{
			"User-Agent": "OnionCLI/1.0",
			"Accept":     "application/json, text/plain, */*",
		},
		History: HistoryConfig{
			Enabled:    true,
			MaxEntries: 100,
			AutoSave:   true,
		},
	}
}

// Load loads the configuration from file
func (m *Manager) Load() error {
	if err := m.viper.ReadInConfig(); err != nil {
		return err
	}

	config := &Config{}
	if err := m.viper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	m.config = config
	return nil
}

// Save saves the configuration to file
func (m *Manager) Save() error {
	if m.config == nil {
		return fmt.Errorf("no config to save")
	}

	// Update viper with current config values
	m.viper.Set("tor", m.config.Tor)
	m.viper.Set("http", m.config.HTTP)
	m.viper.Set("ui", m.config.UI)
	m.viper.Set("default_headers", m.config.DefaultHeaders)
	m.viper.Set("history", m.config.History)

	return m.viper.WriteConfig()
}

// Get returns the current configuration
func (m *Manager) Get() *Config {
	return m.config
}

// Set updates the configuration
func (m *Manager) Set(config *Config) {
	m.config = config
}

// GetTorProxyAddress returns the full Tor proxy address
func (m *Manager) GetTorProxyAddress() string {
	return fmt.Sprintf("%s:%d", m.config.Tor.ProxyAddr, m.config.Tor.ProxyPort)
}

// GetHTTPTimeout returns the HTTP timeout as a duration
func (m *Manager) GetHTTPTimeout() time.Duration {
	return time.Duration(m.config.HTTP.Timeout) * time.Second
}

// GetTorTimeout returns the Tor timeout as a duration
func (m *Manager) GetTorTimeout() time.Duration {
	return time.Duration(m.config.Tor.Timeout) * time.Second
}

// UpdateTorSettings updates Tor-specific settings
func (m *Manager) UpdateTorSettings(enabled bool, proxyAddr string, proxyPort int, timeout int) {
	m.config.Tor.Enabled = enabled
	m.config.Tor.ProxyAddr = proxyAddr
	m.config.Tor.ProxyPort = proxyPort
	m.config.Tor.Timeout = timeout
}

// UpdateHTTPSettings updates HTTP-specific settings
func (m *Manager) UpdateHTTPSettings(timeout int, followRedirects bool, maxRedirects int, verifySSL bool, userAgent string) {
	m.config.HTTP.Timeout = timeout
	m.config.HTTP.FollowRedirects = followRedirects
	m.config.HTTP.MaxRedirects = maxRedirects
	m.config.HTTP.VerifySSL = verifySSL
	m.config.HTTP.UserAgent = userAgent
}

// UpdateUISettings updates UI-specific settings
func (m *Manager) UpdateUISettings(theme string, showLineNumbers bool, autoSave bool, confirmExit bool) {
	m.config.UI.Theme = theme
	m.config.UI.ShowLineNumbers = showLineNumbers
	m.config.UI.AutoSave = autoSave
	m.config.UI.ConfirmExit = confirmExit
}

// AddDefaultHeader adds or updates a default header
func (m *Manager) AddDefaultHeader(key, value string) {
	if m.config.DefaultHeaders == nil {
		m.config.DefaultHeaders = make(map[string]string)
	}
	m.config.DefaultHeaders[key] = value
}

// RemoveDefaultHeader removes a default header
func (m *Manager) RemoveDefaultHeader(key string) {
	if m.config.DefaultHeaders != nil {
		delete(m.config.DefaultHeaders, key)
	}
}

// GetDefaultHeaders returns a copy of the default headers
func (m *Manager) GetDefaultHeaders() map[string]string {
	headers := make(map[string]string)
	for k, v := range m.config.DefaultHeaders {
		headers[k] = v
	}
	return headers
}

// Validate validates the configuration
func (m *Manager) Validate() error {
	if m.config == nil {
		return fmt.Errorf("config is nil")
	}

	// Validate Tor settings
	if m.config.Tor.ProxyPort < 1 || m.config.Tor.ProxyPort > 65535 {
		return fmt.Errorf("invalid Tor proxy port: %d", m.config.Tor.ProxyPort)
	}

	if m.config.Tor.Timeout < 1 {
		return fmt.Errorf("Tor timeout must be at least 1 second")
	}

	// Validate HTTP settings
	if m.config.HTTP.Timeout < 1 {
		return fmt.Errorf("HTTP timeout must be at least 1 second")
	}

	if m.config.HTTP.MaxRedirects < 0 {
		return fmt.Errorf("max redirects cannot be negative")
	}

	// Validate History settings
	if m.config.History.MaxEntries < 1 {
		return fmt.Errorf("history max entries must be at least 1")
	}

	return nil
}

// Reset resets the configuration to defaults
func (m *Manager) Reset() error {
	m.config = m.getDefaultConfig()
	return m.Save()
}

// GetConfigPath returns the path to the configuration file
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// Export exports the configuration to a file
func (m *Manager) Export(filename string) error {
	tempViper := viper.New()
	tempViper.Set("tor", m.config.Tor)
	tempViper.Set("http", m.config.HTTP)
	tempViper.Set("ui", m.config.UI)
	tempViper.Set("default_headers", m.config.DefaultHeaders)
	tempViper.Set("history", m.config.History)

	return tempViper.WriteConfigAs(filename)
}

// Import imports configuration from a file
func (m *Manager) Import(filename string) error {
	tempViper := viper.New()
	tempViper.SetConfigFile(filename)

	if err := tempViper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := &Config{}
	if err := tempViper.Unmarshal(config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate imported config
	oldConfig := m.config
	m.config = config
	if err := m.Validate(); err != nil {
		m.config = oldConfig
		return fmt.Errorf("imported config is invalid: %w", err)
	}

	return m.Save()
}
