package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LoadingSpinner represents a loading spinner with custom styling
type LoadingSpinner struct {
	spinner spinner.Model
	message string
	visible bool
	style   lipgloss.Style
}

// NewLoadingSpinner creates a new loading spinner
func NewLoadingSpinner() LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4"))

	return LoadingSpinner{
		spinner: s,
		visible: false,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD")).
			Margin(0, 1),
	}
}

// Show displays the spinner with a message
func (ls *LoadingSpinner) Show(message string) tea.Cmd {
	ls.message = message
	ls.visible = true
	return ls.spinner.Tick
}

// Hide hides the spinner
func (ls *LoadingSpinner) Hide() {
	ls.visible = false
	ls.message = ""
}

// Update updates the spinner
func (ls LoadingSpinner) Update(msg tea.Msg) (LoadingSpinner, tea.Cmd) {
	if !ls.visible {
		return ls, nil
	}

	var cmd tea.Cmd
	ls.spinner, cmd = ls.spinner.Update(msg)
	return ls, cmd
}

// View renders the spinner
func (ls LoadingSpinner) View() string {
	if !ls.visible {
		return ""
	}

	return ls.style.Render(ls.spinner.View() + " " + ls.message)
}

// IsVisible returns whether the spinner is currently visible
func (ls LoadingSpinner) IsVisible() bool {
	return ls.visible
}

// ProgressIndicator shows progress for long-running operations
type ProgressIndicator struct {
	current int
	total   int
	message string
	visible bool
	style   lipgloss.Style
}

// NewProgressIndicator creates a new progress indicator
func NewProgressIndicator() ProgressIndicator {
	return ProgressIndicator{
		visible: false,
		style: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Margin(0, 1),
	}
}

// Show displays the progress indicator
func (pi *ProgressIndicator) Show(message string, total int) {
	pi.message = message
	pi.total = total
	pi.current = 0
	pi.visible = true
}

// Update updates the progress
func (pi *ProgressIndicator) Update(current int) {
	if pi.visible {
		pi.current = current
		if pi.current > pi.total {
			pi.current = pi.total
		}
	}
}

// Increment increments the progress by 1
func (pi *ProgressIndicator) Increment() {
	if pi.visible {
		pi.current++
		if pi.current > pi.total {
			pi.current = pi.total
		}
	}
}

// Hide hides the progress indicator
func (pi *ProgressIndicator) Hide() {
	pi.visible = false
	pi.message = ""
	pi.current = 0
	pi.total = 0
}

// View renders the progress indicator
func (pi ProgressIndicator) View() string {
	if !pi.visible {
		return ""
	}

	percentage := 0
	if pi.total > 0 {
		percentage = (pi.current * 100) / pi.total
	}

	// Create a simple progress bar
	barWidth := 20
	filled := (pi.current * barWidth) / pi.total
	if filled > barWidth {
		filled = barWidth
	}

	bar := ""
	for i := 0; i < barWidth; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return pi.style.Render(
		fmt.Sprintf("%s [%s] %d%% (%d/%d)",
			pi.message, bar, percentage, pi.current, pi.total))
}

// IsVisible returns whether the progress indicator is visible
func (pi ProgressIndicator) IsVisible() bool {
	return pi.visible
}

// IsComplete returns whether the progress is complete
func (pi ProgressIndicator) IsComplete() bool {
	return pi.visible && pi.current >= pi.total
}

// StatusIndicator shows various status messages with icons
type StatusIndicator struct {
	message   string
	status    StatusType
	visible   bool
	timestamp time.Time
	timeout   time.Duration
}

// StatusType represents different types of status messages
type StatusType int

const (
	StatusInfo StatusType = iota
	StatusSuccess
	StatusWarning
	StatusError
	StatusLoading
)

// NewStatusIndicator creates a new status indicator
func NewStatusIndicator() StatusIndicator {
	return StatusIndicator{
		visible: false,
		timeout: 5 * time.Second, // Auto-hide after 5 seconds
	}
}

// Show displays a status message
func (si *StatusIndicator) Show(message string, status StatusType) {
	si.message = message
	si.status = status
	si.visible = true
	si.timestamp = time.Now()
}

// ShowWithTimeout displays a status message with custom timeout
func (si *StatusIndicator) ShowWithTimeout(message string, status StatusType, timeout time.Duration) {
	si.Show(message, status)
	si.timeout = timeout
}

// Hide hides the status indicator
func (si *StatusIndicator) Hide() {
	si.visible = false
	si.message = ""
}

// Update updates the status indicator (handles auto-hide)
func (si StatusIndicator) Update() StatusIndicator {
	if si.visible && si.timeout > 0 {
		if time.Since(si.timestamp) > si.timeout {
			si.Hide()
		}
	}
	return si
}

// View renders the status indicator
func (si StatusIndicator) View() string {
	if !si.visible {
		return ""
	}

	var icon string
	var style lipgloss.Style

	switch si.status {
	case StatusInfo:
		icon = "ℹ️"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD"))
	case StatusSuccess:
		icon = "✅"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
	case StatusWarning:
		icon = "⚠️"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F1FA8C"))
	case StatusError:
		icon = "❌"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))
	case StatusLoading:
		icon = "⏳"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#BD93F9"))
	default:
		icon = "•"
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("#F8F8F2"))
	}

	return style.Render(fmt.Sprintf("%s %s", icon, si.message))
}

// IsVisible returns whether the status indicator is visible
func (si StatusIndicator) IsVisible() bool {
	return si.visible
}

// KeyboardShortcuts provides a help display for keyboard shortcuts
type KeyboardShortcuts struct {
	shortcuts map[string]string
	visible   bool
	style     lipgloss.Style
}

// NewKeyboardShortcuts creates a new keyboard shortcuts helper
func NewKeyboardShortcuts() KeyboardShortcuts {
	shortcuts := map[string]string{
		"Tab/Shift+Tab": "Navigate fields",
		"Enter":         "Send request / Select",
		"Esc":           "Go back / Cancel",
		"h":             "View history",
		"a":             "Configure auth",
		"s":             "Save request",
		"e":             "View error details",
		"c":             "Settings",
		"r":             "Retry request",
		"Ctrl+C/q":      "Quit",
		"?":             "Toggle help",
	}

	return KeyboardShortcuts{
		shortcuts: shortcuts,
		visible:   false,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666666")).
			Padding(1).
			Margin(1),
	}
}

// Toggle toggles the visibility of keyboard shortcuts
func (ks *KeyboardShortcuts) Toggle() {
	ks.visible = !ks.visible
}

// Show shows the keyboard shortcuts
func (ks *KeyboardShortcuts) Show() {
	ks.visible = true
}

// Hide hides the keyboard shortcuts
func (ks *KeyboardShortcuts) Hide() {
	ks.visible = false
}

// View renders the keyboard shortcuts
func (ks KeyboardShortcuts) View() string {
	if !ks.visible {
		return ""
	}

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Keyboard Shortcuts:"))
	lines = append(lines, "")

	for key, description := range ks.shortcuts {
		keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#8BE9FD")).Bold(true)
		line := fmt.Sprintf("%s: %s", keyStyle.Render(key), description)
		lines = append(lines, line)
	}

	return ks.style.Render(strings.Join(lines, "\n"))
}

// IsVisible returns whether the shortcuts are visible
func (ks KeyboardShortcuts) IsVisible() bool {
	return ks.visible
}
