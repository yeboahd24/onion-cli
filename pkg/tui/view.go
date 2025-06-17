package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Styles for the TUI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			MarginBottom(1)

	focusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	blurredStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(0, 1)

	buttonStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#7D56F4")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 2).
			Margin(1, 0)

	buttonFocusedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#9D7BF4")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 2).
				Margin(1, 0).
				Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true).
			Margin(1, 0)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true).
			Margin(1, 0)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Margin(1, 0)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Margin(1, 0)
)

// View renders the main view
func (m Model) View() string {
	// Handle keyboard shortcuts overlay (highest priority)
	if m.keyboardShortcuts.IsVisible() {
		baseView := m.renderCurrentState()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.keyboardShortcuts.View()) + "\n" + baseView
	}

	// Handle error viewer overlay
	if m.errorViewer.IsVisible() {
		return m.errorViewer.View()
	}

	// Handle auth dialog overlay
	if m.authDialog.visible {
		baseView := m.renderCurrentState()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.authDialog.View()) + "\n" + baseView
	}

	// Handle save dialog overlay
	if m.saveDialog.visible {
		baseView := m.renderCurrentState()
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.saveDialog.View()) + "\n" + baseView
	}

	return m.renderCurrentState()
}

// renderCurrentState renders the current state view
func (m Model) renderCurrentState() string {
	switch m.state {
	case StateRequestBuilder:
		return m.renderRequestBuilder()
	case StateResponse:
		return m.renderResponse()
	case StateHistory:
		return m.renderHistory()
	case StateCollections:
		return m.renderCollections()
	case StateEnvironments:
		return m.renderEnvironments()
	default:
		return m.renderRequestBuilder()
	}
}

// renderHistory renders the history view
func (m Model) renderHistory() string {
	return m.historyViewer.View()
}

// renderCollections renders the collections view
func (m Model) renderCollections() string {
	return m.collectionsViewer.View()
}

// renderEnvironments renders the environments view
func (m Model) renderEnvironments() string {
	return m.environmentsViewer.View()
}

// renderRequestBuilder renders the request builder interface
func (m Model) renderRequestBuilder() string {
	var sections []string

	// Title (more compact)
	title := titleStyle.Render("OnionCLI - .onion API Client")
	sections = append(sections, title)

	// URL input
	urlLabel := "URL:"
	var urlSection string
	if m.focusedField == FocusURL {
		urlSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", urlLabel, m.urlInput.View()))
	} else {
		urlSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", urlLabel, m.urlInput.View()))
	}
	sections = append(sections, urlSection)

	// Method selection
	methodLabel := "HTTP Method:"
	var methodSection string
	if m.focusedField == FocusMethod {
		methodSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", methodLabel, m.methodList.View()))
	} else {
		methodSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", methodLabel, m.methodList.View()))
	}
	sections = append(sections, methodSection)

	// Headers
	headersLabel := "Headers:"
	var headersSection string
	if m.focusedField == FocusHeaders {
		headersSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", headersLabel, m.headersArea.View()))
	} else {
		headersSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", headersLabel, m.headersArea.View()))
	}
	sections = append(sections, headersSection)

	// Body
	bodyLabel := "Request Body:"
	var bodySection string
	if m.focusedField == FocusBody {
		bodySection = focusedStyle.Render(fmt.Sprintf("%s\n%s", bodyLabel, m.bodyArea.View()))
	} else {
		bodySection = blurredStyle.Render(fmt.Sprintf("%s\n%s", bodyLabel, m.bodyArea.View()))
	}
	sections = append(sections, bodySection)

	// Submit button
	var submitButton string
	if m.focusedField == FocusSubmit {
		submitButton = buttonFocusedStyle.Render("Send Request")
	} else {
		submitButton = buttonStyle.Render("Send Request")
	}
	sections = append(sections, submitButton)

	// Error alert (enhanced error display)
	if m.errorAlert.IsVisible() {
		sections = append(sections, m.errorAlert.View())
		// Add hint about detailed error view
		sections = append(sections, helpStyle.Render("Press 'e' for detailed error information"))
	} else if m.errorMessage != "" {
		sections = append(sections, errorStyle.Render("❌ "+m.errorMessage))
	}

	// Loading spinner
	if m.loadingSpinner.IsVisible() {
		sections = append(sections, m.loadingSpinner.View())
	}

	// Status indicator
	if m.statusIndicator.IsVisible() {
		sections = append(sections, m.statusIndicator.View())
	}

	// Status messages (fallback)
	if m.statusMessage != "" && !m.statusIndicator.IsVisible() {
		sections = append(sections, statusStyle.Render("ℹ️  "+m.statusMessage))
	}

	// Help text
	help := helpStyle.Render(m.renderHelp())
	sections = append(sections, help)

	return strings.Join(sections, "\n")
}

// renderResponse renders the response view using the response viewer
func (m Model) renderResponse() string {
	if m.currentResponse == nil {
		return titleStyle.Render("No response to display")
	}

	return m.responseViewer.View()
}

// renderHelp renders the help text
func (m Model) renderHelp() string {
	authStatus := "No auth"
	if m.authConfig != nil {
		authStatus = fmt.Sprintf("Auth: %s", m.authConfig.Type)
	}

	errorHint := ""
	if m.errorAlert.IsVisible() {
		errorHint = ", e for error details"
	}

	retryHint := ""
	if m.currentRequest != nil && !m.loading {
		retryHint = ", r to retry"
	}

	baseHelp := fmt.Sprintf("? for shortcuts%s%s", errorHint, retryHint)

	switch m.focusedField {
	case FocusURL:
		return fmt.Sprintf("Enter a .onion URL. Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, s to save, Enter/Ctrl+Enter to send | %s | %s", authStatus, baseHelp)
	case FocusMethod:
		return fmt.Sprintf("Select HTTP method with ↑/↓ arrows. Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, Ctrl+Enter to send | %s | %s", authStatus, baseHelp)
	case FocusHeaders:
		return fmt.Sprintf("Enter headers in 'key: value' format, one per line. Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, Ctrl+Enter to send | %s | %s", authStatus, baseHelp)
	case FocusBody:
		return fmt.Sprintf("Enter request body (JSON, XML, or plain text). Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, Ctrl+Enter to send | %s | %s", authStatus, baseHelp)
	case FocusSubmit:
		return fmt.Sprintf("Press Enter to send the request. Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, s to save | %s | %s", authStatus, baseHelp)
	default:
		return fmt.Sprintf("Tab/Shift+Tab to navigate, a for auth, c for collections, v for environments, h for history, s to save, Ctrl+Enter to send request, q/Ctrl+C to quit | %s | %s", authStatus, baseHelp)
	}
}
