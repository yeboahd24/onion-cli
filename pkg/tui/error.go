package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"onioncli/pkg/api"
)

// ErrorViewer displays detailed error information and suggestions
type ErrorViewer struct {
	viewport    viewport.Model
	error       *api.DiagnosticError
	visible     bool
	width       int
	height      int
}

// NewErrorViewer creates a new error viewer
func NewErrorViewer(width, height int) ErrorViewer {
	vp := viewport.New(width-4, height-8)
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF5555")).
		Padding(1)

	return ErrorViewer{
		viewport: vp,
		visible:  false,
		width:    width,
		height:   height,
	}
}

// Show displays the error viewer with the given error
func (ev *ErrorViewer) Show(err *api.DiagnosticError) {
	ev.error = err
	ev.visible = true
	content := ev.formatError(err)
	ev.viewport.SetContent(content)
}

// Hide hides the error viewer
func (ev *ErrorViewer) Hide() {
	ev.visible = false
	ev.error = nil
}

// Update handles error viewer updates
func (ev ErrorViewer) Update(msg tea.Msg) (ErrorViewer, tea.Cmd) {
	if !ev.visible {
		return ev, nil
	}

	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			ev.Hide()
			return ev, nil
		}
	}

	ev.viewport, cmd = ev.viewport.Update(msg)
	return ev, cmd
}

// View renders the error viewer
func (ev ErrorViewer) View() string {
	if !ev.visible || ev.error == nil {
		return ""
	}

	// Header with error type and icon
	header := ev.renderErrorHeader()
	
	// Viewport with error details
	content := ev.viewport.View()
	
	// Footer with navigation help
	footer := ev.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderErrorHeader renders the error header with type and icon
func (ev ErrorViewer) renderErrorHeader() string {
	if ev.error == nil {
		return ""
	}

	// Choose icon and color based on error type
	var icon string
	var color lipgloss.Color
	
	switch ev.error.Type {
	case api.ErrorTypeTor:
		icon = "🧅"
		color = lipgloss.Color("#FF6B6B")
	case api.ErrorTypeNetwork:
		icon = "🌐"
		color = lipgloss.Color("#FF8E53")
	case api.ErrorTypeTimeout:
		icon = "⏱️"
		color = lipgloss.Color("#FFD93D")
	case api.ErrorTypeDNS:
		icon = "🔍"
		color = lipgloss.Color("#6BCF7F")
	case api.ErrorTypeAuth:
		icon = "🔐"
		color = lipgloss.Color("#4D96FF")
	default:
		icon = "❌"
		color = lipgloss.Color("#FF5555")
	}

	title := fmt.Sprintf("%s %s Error", icon, strings.Title(string(ev.error.Type)))
	
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Render(title)
}

// renderFooter renders navigation help
func (ev ErrorViewer) renderFooter() string {
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("↑/↓ scroll • esc/q close error details")
	
	return help
}

// formatError formats the error for display
func (ev ErrorViewer) formatError(err *api.DiagnosticError) string {
	var sections []string

	// Error message
	sections = append(sections, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5555")).
		Bold(true).
		Render("Error Message:"))
	sections = append(sections, err.Message)
	sections = append(sections, "")

	// URL if available
	if err.URL != "" {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Bold(true).
			Render("URL:"))
		sections = append(sections, err.URL)
		sections = append(sections, "")
	}

	// Status code if available
	if err.StatusCode != 0 {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C")).
			Bold(true).
			Render("Status Code:"))
		sections = append(sections, fmt.Sprintf("%d", err.StatusCode))
		sections = append(sections, "")
	}

	// Error type specific information
	sections = append(sections, ev.renderTypeSpecificInfo(err))

	// Suggestions
	if len(err.Suggestions) > 0 {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true).
			Render("💡 Suggestions:"))
		
		for i, suggestion := range err.Suggestions {
			suggestionText := fmt.Sprintf("  %d. %s", i+1, suggestion)
			sections = append(sections, suggestionText)
		}
		sections = append(sections, "")
	}

	// Retry information
	if err.IsRetryable() {
		retryInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F1FA8C")).
			Render("🔄 This error might be resolved by retrying the request.")
		sections = append(sections, retryInfo)
	} else {
		retryInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF79C6")).
			Render("⚠️  This error is unlikely to be resolved by retrying.")
		sections = append(sections, retryInfo)
	}

	return strings.Join(sections, "\n")
}

// renderTypeSpecificInfo renders additional information based on error type
func (ev ErrorViewer) renderTypeSpecificInfo(err *api.DiagnosticError) string {
	var info []string

	switch err.Type {
	case api.ErrorTypeTor:
		info = append(info, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true).
			Render("🧅 Tor-Specific Information:"))
		info = append(info, "• Tor requests require the Tor service to be running")
		info = append(info, "• Default Tor proxy: 127.0.0.1:9050")
		info = append(info, "• .onion services can be slow or temporarily unavailable")
		info = append(info, "")

	case api.ErrorTypeNetwork:
		info = append(info, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true).
			Render("🌐 Network Information:"))
		info = append(info, "• Check your internet connection")
		info = append(info, "• Verify firewall settings")
		info = append(info, "• Some networks block certain ports or protocols")
		info = append(info, "")

	case api.ErrorTypeTimeout:
		info = append(info, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true).
			Render("⏱️ Timeout Information:"))
		info = append(info, "• Tor requests typically take 10-30 seconds")
		info = append(info, "• .onion services can be slower than regular websites")
		info = append(info, "• Consider increasing timeout in settings")
		info = append(info, "")

	case api.ErrorTypeAuth:
		info = append(info, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9")).
			Bold(true).
			Render("🔐 Authentication Information:"))
		info = append(info, "• Verify your credentials are correct")
		info = append(info, "• Check if your API key/token has expired")
		info = append(info, "• Ensure you have the required permissions")
		info = append(info, "")
	}

	return strings.Join(info, "\n")
}

// Resize updates the error viewer size
func (ev *ErrorViewer) Resize(width, height int) {
	ev.width = width
	ev.height = height
	ev.viewport.Width = width - 4
	ev.viewport.Height = height - 8
}

// IsVisible returns whether the error viewer is currently visible
func (ev ErrorViewer) IsVisible() bool {
	return ev.visible
}

// ErrorAlert represents a simple error alert for inline display
type ErrorAlert struct {
	message     string
	errorType   api.ErrorType
	suggestions []string
	visible     bool
}

// NewErrorAlert creates a new error alert
func NewErrorAlert() ErrorAlert {
	return ErrorAlert{
		visible: false,
	}
}

// Show displays the error alert
func (ea *ErrorAlert) Show(err *api.DiagnosticError) {
	if err == nil {
		ea.Hide()
		return
	}
	
	ea.message = err.Message
	ea.errorType = err.Type
	ea.suggestions = err.Suggestions
	ea.visible = true
}

// Hide hides the error alert
func (ea *ErrorAlert) Hide() {
	ea.visible = false
	ea.message = ""
	ea.suggestions = nil
}

// View renders the error alert
func (ea ErrorAlert) View() string {
	if !ea.visible {
		return ""
	}

	// Choose color based on error type
	var color lipgloss.Color
	switch ea.errorType {
	case api.ErrorTypeTor:
		color = lipgloss.Color("#FF6B6B")
	case api.ErrorTypeNetwork:
		color = lipgloss.Color("#FF8E53")
	case api.ErrorTypeTimeout:
		color = lipgloss.Color("#FFD93D")
	case api.ErrorTypeAuth:
		color = lipgloss.Color("#4D96FF")
	default:
		color = lipgloss.Color("#FF5555")
	}

	style := lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color)

	content := fmt.Sprintf("❌ %s", ea.message)
	
	// Add first suggestion if available
	if len(ea.suggestions) > 0 {
		content += fmt.Sprintf("\n💡 %s", ea.suggestions[0])
	}

	return style.Render(content)
}

// IsVisible returns whether the error alert is visible
func (ea ErrorAlert) IsVisible() bool {
	return ea.visible
}
