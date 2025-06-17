package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"onioncli/pkg/api"
)

// ResponseViewer handles the display of HTTP responses
type ResponseViewer struct {
	viewport viewport.Model
	response *api.Response
	width    int
	height   int
}

// NewResponseViewer creates a new response viewer
func NewResponseViewer(width, height int) ResponseViewer {
	vp := viewport.New(width-4, height-10) // Leave space for borders and controls
	vp.Style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(1)

	return ResponseViewer{
		viewport: vp,
		width:    width,
		height:   height,
	}
}

// SetResponse sets the response to display
func (rv *ResponseViewer) SetResponse(response *api.Response) {
	rv.response = response
	content := rv.formatResponse(response)
	rv.viewport.SetContent(content)
}

// Update handles viewport updates
func (rv ResponseViewer) Update(msg tea.Msg) (ResponseViewer, tea.Cmd) {
	var cmd tea.Cmd
	rv.viewport, cmd = rv.viewport.Update(msg)
	return rv, cmd
}

// View renders the response viewer
func (rv ResponseViewer) View() string {
	if rv.response == nil {
		return rv.viewport.View()
	}

	// Header with response summary
	header := rv.renderResponseHeader()
	
	// Viewport with response details
	content := rv.viewport.View()
	
	// Footer with navigation help
	footer := rv.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

// renderResponseHeader renders the response status and timing information
func (rv ResponseViewer) renderResponseHeader() string {
	if rv.response == nil {
		return ""
	}

	// Status styling based on response code
	var statusStyle lipgloss.Style
	if rv.response.IsSuccess() {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Bold(true)
	} else if rv.response.IsClientError() {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")).Bold(true)
	} else if rv.response.IsServerError() {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Bold(true)
	} else {
		statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Bold(true)
	}

	status := statusStyle.Render(fmt.Sprintf("Status: %s", rv.response.Status))
	duration := lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C")).Render(
		fmt.Sprintf("Duration: %v", rv.response.Duration))
	timestamp := lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9")).Render(
		fmt.Sprintf("Time: %s", rv.response.Timestamp.Format("15:04:05")))

	return lipgloss.JoinHorizontal(lipgloss.Left, status, "  ", duration, "  ", timestamp)
}

// renderFooter renders navigation help
func (rv ResponseViewer) renderFooter() string {
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666")).
		Render("↑/↓ scroll • esc back to request builder • q quit")
	
	return help
}

// formatResponse formats the response for display
func (rv ResponseViewer) formatResponse(response *api.Response) string {
	var sections []string

	// Request summary (if available)
	sections = append(sections, lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Bold(true).
		Render("Response Details"))
	sections = append(sections, "")

	// Headers section
	if len(response.Headers) > 0 {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true).
			Render("Headers:"))
		
		for key, value := range response.Headers {
			headerLine := fmt.Sprintf("  %s: %s", 
				lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Render(key),
				value)
			sections = append(sections, headerLine)
		}
		sections = append(sections, "")
	}

	// Body section
	if response.Body != "" {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true).
			Render("Response Body:"))
		
		// Try to pretty-print JSON
		prettyBody, err := response.PrettyPrintJSON()
		if err != nil {
			prettyBody = response.Body
		}

		// Syntax highlighting for JSON (basic)
		if strings.Contains(response.Headers["Content-Type"], "application/json") {
			prettyBody = rv.highlightJSON(prettyBody)
		}

		sections = append(sections, prettyBody)
	} else {
		sections = append(sections, lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Italic(true).
			Render("(No response body)"))
	}

	return strings.Join(sections, "\n")
}

// highlightJSON provides basic JSON syntax highlighting
func (rv ResponseViewer) highlightJSON(jsonStr string) string {
	// Basic JSON highlighting - this is a simple implementation
	// In a production app, you might want to use a proper JSON syntax highlighter
	
	lines := strings.Split(jsonStr, "\n")
	var highlighted []string
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Highlight different JSON elements
		if strings.Contains(line, ":") && (strings.Contains(line, "\"") || strings.Contains(line, "'")) {
			// Key-value pairs
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Render(parts[0])
				value := strings.TrimSpace(parts[1])
				
				// Highlight string values
				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					value = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C")).Render(value)
				} else if value == "true" || value == "false" {
					// Boolean values
					value = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render(value)
				} else if value == "null" {
					// Null values
					value = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272A4")).Render(value)
				}
				
				line = key + ":" + value
			}
		}
		
		highlighted = append(highlighted, line)
	}
	
	return strings.Join(highlighted, "\n")
}

// Resize updates the viewport size
func (rv *ResponseViewer) Resize(width, height int) {
	rv.width = width
	rv.height = height
	rv.viewport.Width = width - 4
	rv.viewport.Height = height - 10
}
