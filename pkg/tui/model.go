package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"onioncli/pkg/api"
	"onioncli/pkg/collections"
	"onioncli/pkg/history"
)

// AppState represents the current state of the application
type AppState int

const (
	StateRequestBuilder AppState = iota
	StateResponse
	StateHistory
	StateCollections
	StateEnvironments
	StateSettings
)

// FocusedField represents which field is currently focused
type FocusedField int

const (
	FocusURL FocusedField = iota
	FocusMethod
	FocusHeaders
	FocusBody
	FocusSubmit
)

// Model represents the main application model
type Model struct {
	state        AppState
	focusedField FocusedField
	width        int
	height       int

	// Request builder components
	urlInput    textinput.Model
	methodList  list.Model
	headersArea textarea.Model
	bodyArea    textarea.Model

	// API client
	client *api.Client

	// Authentication
	authManager *api.AuthManager
	authDialog  AuthDialog
	authConfig  *api.AuthConfig

	// Collections and environments
	collectionsManager *collections.Manager
	collectionsViewer  CollectionsViewer
	environmentsViewer EnvironmentsViewer

	// History manager
	historyManager *history.Manager
	historyViewer  HistoryViewer
	saveDialog     SaveRequestDialog

	// Current request and response
	currentRequest  *api.Request
	currentResponse *api.Response

	// Response viewer
	responseViewer ResponseViewer

	// Error handling
	errorAnalyzer *api.ErrorAnalyzer
	errorViewer   ErrorViewer
	errorAlert    ErrorAlert

	// Performance and UI enhancements
	loadingSpinner    LoadingSpinner
	statusIndicator   StatusIndicator
	keyboardShortcuts KeyboardShortcuts

	// Status and error messages
	statusMessage string
	errorMessage  string
	loading       bool
}

// HTTPMethod represents an HTTP method for the list
type HTTPMethod struct {
	name string
}

func (m HTTPMethod) FilterValue() string { return m.name }
func (m HTTPMethod) Title() string       { return m.name }
func (m HTTPMethod) Description() string { return "" }

// NewModel creates a new TUI model
func NewModel() (*Model, error) {
	// Initialize API client
	client, err := api.NewClient(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create API client: %w", err)
	}

	// Initialize authentication manager
	authManager := api.NewAuthManager()

	// Initialize error analyzer
	errorAnalyzer := api.NewErrorAnalyzer()

	// Initialize collections manager
	collectionsManager, err := collections.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create collections manager: %w", err)
	}

	// Initialize history manager
	historyManager, err := history.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create history manager: %w", err)
	}

	// Initialize URL input
	urlInput := textinput.New()
	urlInput.Placeholder = "Enter .onion URL (e.g., http://3g2upl4pq6kufc4m.onion)"
	urlInput.Focus()
	urlInput.CharLimit = 500
	urlInput.Width = 80

	// Initialize HTTP method list
	methods := []list.Item{
		HTTPMethod{name: "GET"},
		HTTPMethod{name: "POST"},
		HTTPMethod{name: "PUT"},
		HTTPMethod{name: "DELETE"},
		HTTPMethod{name: "PATCH"},
		HTTPMethod{name: "HEAD"},
		HTTPMethod{name: "OPTIONS"},
	}

	methodList := list.New(methods, list.NewDefaultDelegate(), 20, 8)
	methodList.Title = "HTTP Method"
	methodList.SetShowStatusBar(false)
	methodList.SetFilteringEnabled(false)
	methodList.SetShowHelp(false)

	// Initialize headers textarea
	headersArea := textarea.New()
	headersArea.Placeholder = "Headers (key: value format, one per line)\nUser-Agent: OnionCLI/1.0\nContent-Type: application/json"
	headersArea.SetWidth(80)
	headersArea.SetHeight(6)

	// Initialize body textarea
	bodyArea := textarea.New()
	bodyArea.Placeholder = "Request body (JSON, XML, or plain text)"
	bodyArea.SetWidth(80)
	bodyArea.SetHeight(10)

	model := &Model{
		state:              StateRequestBuilder,
		focusedField:       FocusURL,
		urlInput:           urlInput,
		methodList:         methodList,
		headersArea:        headersArea,
		bodyArea:           bodyArea,
		client:             client,
		authManager:        authManager,
		authDialog:         NewAuthDialog(80, 24),
		collectionsManager: collectionsManager,
		collectionsViewer:  NewCollectionsViewer(collectionsManager, 80, 24),
		environmentsViewer: NewEnvironmentsViewer(collectionsManager, 80, 24),
		historyManager:     historyManager,
		historyViewer:      NewHistoryViewer(historyManager, 80, 24),
		saveDialog:         NewSaveRequestDialog(),
		responseViewer:     NewResponseViewer(80, 24),
		errorAnalyzer:      errorAnalyzer,
		errorViewer:        NewErrorViewer(80, 24),
		errorAlert:         NewErrorAlert(),
		loadingSpinner:     NewLoadingSpinner(),
		statusIndicator:    NewStatusIndicator(),
		keyboardShortcuts:  NewKeyboardShortcuts(),
	}

	return model, nil
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.responseViewer.Resize(msg.Width, msg.Height)
		m.historyViewer.Resize(msg.Width, msg.Height)
		m.collectionsViewer.Resize(msg.Width, msg.Height)
		m.environmentsViewer.Resize(msg.Width, msg.Height)
		m.authDialog.Resize(msg.Width, msg.Height)
		m.errorViewer.Resize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		// Handle auth dialog first
		if m.authDialog.visible {
			m.authDialog, cmd = m.authDialog.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		// Handle save dialog
		if m.saveDialog.visible {
			m.saveDialog, cmd = m.saveDialog.Update(msg)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "h":
			if m.state == StateRequestBuilder {
				m.state = StateHistory
				return m, nil
			}

		case "c":
			if m.state == StateRequestBuilder {
				m.state = StateCollections
				return m, nil
			}

		case "v":
			if m.state == StateRequestBuilder {
				m.state = StateEnvironments
				return m, nil
			}

		case "s":
			if m.state == StateRequestBuilder && m.currentRequest != nil {
				m.saveDialog.Show()
				return m, nil
			}

		case "a":
			if m.state == StateRequestBuilder {
				m.authDialog.Show()
				return m, nil
			}

		case "e":
			if m.errorAlert.IsVisible() {
				// Show detailed error view
				if m.errorAlert.visible {
					// Create a diagnostic error from the alert
					diagnosticError := &api.DiagnosticError{
						Type:        m.errorAlert.errorType,
						Message:     m.errorAlert.message,
						Suggestions: m.errorAlert.suggestions,
					}
					m.errorViewer.Show(diagnosticError)
					return m, nil
				}
			}

		case "r":
			// Retry last request
			if m.currentRequest != nil && !m.loading {
				m.statusIndicator.Show("Retrying request...", StatusLoading)
				return m.sendRequest()
			}

		case "?":
			// Toggle keyboard shortcuts help
			m.keyboardShortcuts.Toggle()
			return m, nil

		case "ctrl+s":
			// Quick save shortcut
			if m.state == StateRequestBuilder && m.currentRequest != nil {
				m.saveDialog.Show()
				return m, nil
			}

		case "tab":
			if m.state == StateRequestBuilder {
				return m.nextField(), nil
			}

		case "shift+tab":
			if m.state == StateRequestBuilder {
				return m.prevField(), nil
			}

		case "enter":
			if m.state == StateRequestBuilder && m.focusedField == FocusSubmit {
				return m.sendRequest()
			} else if m.state == StateHistory {
				if entry := m.historyViewer.GetSelectedEntry(); entry != nil {
					m.loadFromHistory(entry)
					m.state = StateRequestBuilder
					return m, nil
				}
			}

		case "esc":
			if m.state == StateResponse {
				m.state = StateRequestBuilder
				m.focusedField = FocusURL
				m.urlInput.Focus()
				return m, nil
			} else if m.state == StateHistory {
				m.state = StateRequestBuilder
				return m, nil
			} else if m.state == StateCollections {
				m.state = StateRequestBuilder
				return m, nil
			} else if m.state == StateEnvironments {
				m.state = StateRequestBuilder
				return m, nil
			}
			m.errorMessage = ""
			m.statusMessage = ""
			return m, nil
		}

	case SaveRequestMsg:
		if m.currentRequest != nil {
			err := m.historyManager.Save(m.currentRequest, msg.GetName(), msg.GetDescription())
			if err != nil {
				m.errorMessage = fmt.Sprintf("Failed to save request: %v", err)
			} else {
				m.statusMessage = "✅ Request saved to history"
			}
		}
		m.saveDialog.Hide()
		return m, nil

	case AuthConfiguredMsg:
		m.authConfig = msg.config
		m.statusMessage = fmt.Sprintf("✅ Authentication configured: %s", msg.config.Type)
		m.errorMessage = ""
		return m, nil

	case AuthErrorMsg:
		m.errorMessage = fmt.Sprintf("Authentication error: %v", msg.err)
		m.statusMessage = ""
		return m, nil

	case LoadRequestMsg:
		// Load request from collection
		req := msg.request
		m.urlInput.SetValue(req.URL)

		// Set method
		for i, item := range m.methodList.Items() {
			if httpMethod, ok := item.(HTTPMethod); ok && httpMethod.name == req.Method {
				m.methodList.Select(i)
				break
			}
		}

		// Set headers
		var headerLines []string
		for key, value := range req.Headers {
			headerLines = append(headerLines, fmt.Sprintf("%s: %s", key, value))
		}
		m.headersArea.SetValue(strings.Join(headerLines, "\n"))

		// Set body
		m.bodyArea.SetValue(req.Body)

		m.statusMessage = fmt.Sprintf("✅ Loaded request: %s", req.Name)
		m.state = StateRequestBuilder
		return m, nil

	case EnvironmentChangedMsg:
		// Environment changed
		m.statusMessage = fmt.Sprintf("✅ Environment changed to: %s", msg.environment.Name)
		return m, nil

	case RequestSuccessMsg:
		m.currentResponse = msg.response
		m.responseViewer.SetResponse(msg.response)
		m.loading = false
		m.loadingSpinner.Hide()

		// Show success status
		statusMsg := fmt.Sprintf("Request completed successfully (%v)", msg.response.Duration)
		m.statusIndicator.Show(statusMsg, StatusSuccess)
		m.statusMessage = ""
		m.errorMessage = ""
		m.errorAlert.Hide()
		m.state = StateResponse
		return m, nil

	case RequestErrorMsg:
		m.loading = false
		m.loadingSpinner.Hide()

		// Analyze the error for better diagnostics
		diagnosticError := m.errorAnalyzer.AnalyzeError(msg.err, msg.url)
		if diagnosticError != nil {
			m.errorAlert.Show(diagnosticError)
			m.errorMessage = diagnosticError.Message
			m.statusIndicator.Show("Request failed", StatusError)
		} else {
			m.errorMessage = fmt.Sprintf("Request failed: %v", msg.err)
			m.statusIndicator.Show("Request failed", StatusError)
		}

		m.statusMessage = ""
		return m, nil
	}

	// Update error viewer if visible
	if m.errorViewer.IsVisible() {
		m.errorViewer, cmd = m.errorViewer.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	// Update spinner
	m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
	cmds = append(cmds, cmd)

	// Update status indicator
	m.statusIndicator = m.statusIndicator.Update()

	// Update focused component based on current state
	switch m.state {
	case StateResponse:
		m.responseViewer, cmd = m.responseViewer.Update(msg)
		cmds = append(cmds, cmd)
	case StateHistory:
		m.historyViewer, cmd = m.historyViewer.Update(msg)
		cmds = append(cmds, cmd)
	case StateCollections:
		m.collectionsViewer, cmd = m.collectionsViewer.Update(msg)
		cmds = append(cmds, cmd)
	case StateEnvironments:
		m.environmentsViewer, cmd = m.environmentsViewer.Update(msg)
		cmds = append(cmds, cmd)
	default:
		// Update focused component in request builder
		switch m.focusedField {
		case FocusURL:
			m.urlInput, cmd = m.urlInput.Update(msg)
			cmds = append(cmds, cmd)
		case FocusMethod:
			m.methodList, cmd = m.methodList.Update(msg)
			cmds = append(cmds, cmd)
		case FocusHeaders:
			m.headersArea, cmd = m.headersArea.Update(msg)
			cmds = append(cmds, cmd)
		case FocusBody:
			m.bodyArea, cmd = m.bodyArea.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

// loadFromHistory loads a request from history
func (m *Model) loadFromHistory(entry *history.HistoryEntry) {
	req := entry.ToRequest()

	// Set URL
	m.urlInput.SetValue(req.URL)

	// Set method
	for i, item := range m.methodList.Items() {
		if httpMethod, ok := item.(HTTPMethod); ok && httpMethod.name == req.Method {
			m.methodList.Select(i)
			break
		}
	}

	// Set headers
	var headerLines []string
	for key, value := range req.Headers {
		headerLines = append(headerLines, fmt.Sprintf("%s: %s", key, value))
	}
	m.headersArea.SetValue(strings.Join(headerLines, "\n"))

	// Set body
	m.bodyArea.SetValue(req.Body)

	m.statusMessage = fmt.Sprintf("✅ Loaded request: %s", entry.Name)
}

// nextField moves focus to the next field
func (m Model) nextField() Model {
	switch m.focusedField {
	case FocusURL:
		m.focusedField = FocusMethod
		m.urlInput.Blur()
	case FocusMethod:
		m.focusedField = FocusHeaders
		m.headersArea.Focus()
	case FocusHeaders:
		m.focusedField = FocusBody
		m.headersArea.Blur()
		m.bodyArea.Focus()
	case FocusBody:
		m.focusedField = FocusSubmit
		m.bodyArea.Blur()
	case FocusSubmit:
		m.focusedField = FocusURL
		m.urlInput.Focus()
	}
	return m
}

// prevField moves focus to the previous field
func (m Model) prevField() Model {
	switch m.focusedField {
	case FocusURL:
		m.focusedField = FocusSubmit
		m.urlInput.Blur()
	case FocusMethod:
		m.focusedField = FocusURL
		m.urlInput.Focus()
	case FocusHeaders:
		m.focusedField = FocusMethod
		m.headersArea.Blur()
	case FocusBody:
		m.focusedField = FocusHeaders
		m.bodyArea.Blur()
		m.headersArea.Focus()
	case FocusSubmit:
		m.focusedField = FocusBody
		m.bodyArea.Focus()
	}
	return m
}

// sendRequest creates and sends the HTTP request
func (m Model) sendRequest() (Model, tea.Cmd) {
	// Get selected method
	selectedItem := m.methodList.SelectedItem()
	if selectedItem == nil {
		m.errorMessage = "Please select an HTTP method"
		return m, nil
	}
	method := selectedItem.(HTTPMethod).name

	// Get URL
	url := strings.TrimSpace(m.urlInput.Value())
	if url == "" {
		m.errorMessage = "Please enter a URL"
		return m, nil
	}

	// Create request
	req := api.NewRequest(method, url)

	// Parse headers
	headersText := strings.TrimSpace(m.headersArea.Value())
	if headersText != "" {
		headers := m.parseHeaders(headersText)
		for key, value := range headers {
			req.SetHeader(key, value)
		}
	}

	// Set body
	body := strings.TrimSpace(m.bodyArea.Value())
	if body != "" {
		req.SetBody(body)
	}

	// Process request with variable substitution
	req = m.collectionsManager.ProcessRequest(req)

	// Apply authentication if configured
	if m.authConfig != nil {
		if err := m.authManager.ApplyAuth(req, m.authConfig); err != nil {
			m.errorMessage = fmt.Sprintf("Authentication failed: %v", err)
			return m, nil
		}
	}

	// Validate request
	if err := req.Validate(); err != nil {
		m.errorMessage = fmt.Sprintf("Request validation failed: %v", err)
		return m, nil
	}

	m.currentRequest = req
	m.loading = true
	m.errorMessage = ""
	m.statusMessage = ""
	m.errorAlert.Hide()

	// Show loading spinner with appropriate message
	var spinnerMessage string
	if api.IsOnionURL(req.URL) {
		spinnerMessage = "Sending request via Tor (this may take a moment)..."
	} else {
		spinnerMessage = "Sending request..."
	}

	return m, tea.Batch(
		m.loadingSpinner.Show(spinnerMessage),
		m.sendRequestCmd(req),
	)
}

// parseHeaders parses headers from textarea input
func (m Model) parseHeaders(headersText string) map[string]string {
	headers := make(map[string]string)
	lines := strings.Split(headersText, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" && value != "" {
				headers[key] = value
			}
		}
	}

	return headers
}

// sendRequestCmd returns a command to send the HTTP request
func (m Model) sendRequestCmd(req *api.Request) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.Send(req)
		if err != nil {
			return RequestErrorMsg{err: err, url: req.URL}
		}
		return RequestSuccessMsg{response: resp}
	}
}

// RequestSuccessMsg represents a successful request
type RequestSuccessMsg struct {
	response *api.Response
}

// RequestErrorMsg represents a failed request
type RequestErrorMsg struct {
	err error
	url string
}
