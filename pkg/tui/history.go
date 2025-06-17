package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"onioncli/pkg/history"
)

// HistoryItem represents a history entry for the list component
type HistoryItem struct {
	entry history.HistoryEntry
}

func (h HistoryItem) FilterValue() string {
	return h.entry.Name + " " + h.entry.URL + " " + h.entry.Method
}

func (h HistoryItem) Title() string {
	if h.entry.Name != "" {
		return h.entry.Name
	}
	return fmt.Sprintf("%s %s", h.entry.Method, h.entry.URL)
}

func (h HistoryItem) Description() string {
	timeStr := h.entry.Timestamp.Format("2006-01-02 15:04")
	if h.entry.Description != "" {
		return fmt.Sprintf("%s - %s", timeStr, h.entry.Description)
	}
	return timeStr
}

// HistoryViewer handles the history browsing interface
type HistoryViewer struct {
	list        list.Model
	searchInput textinput.Model
	manager     *history.Manager
	searching   bool
	width       int
	height      int
	allEntries  []history.HistoryEntry
}

// NewHistoryViewer creates a new history viewer
func NewHistoryViewer(manager *history.Manager, width, height int) HistoryViewer {
	// Create list
	items := make([]list.Item, 0)
	entries := manager.GetEntries()

	for _, entry := range entries {
		items = append(items, HistoryItem{entry: entry})
	}

	l := list.New(items, list.NewDefaultDelegate(), width-4, height-8)
	l.Title = "Request History"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false) // We'll handle our own filtering
	l.SetShowHelp(true)

	// Create search input
	searchInput := textinput.New()
	searchInput.Placeholder = "Search history..."
	searchInput.CharLimit = 100
	searchInput.Width = width - 10

	return HistoryViewer{
		list:        l,
		searchInput: searchInput,
		manager:     manager,
		searching:   false,
		width:       width,
		height:      height,
		allEntries:  entries,
	}
}

// Update handles history viewer updates
func (hv HistoryViewer) Update(msg tea.Msg) (HistoryViewer, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if hv.searching {
			switch msg.String() {
			case "enter":
				// Apply search
				hv.searching = false
				hv.searchInput.Blur()
				hv.applySearch()
				return hv, nil
			case "esc":
				// Cancel search
				hv.searching = false
				hv.searchInput.Blur()
				hv.searchInput.SetValue("")
				hv.resetList()
				return hv, nil
			default:
				hv.searchInput, cmd = hv.searchInput.Update(msg)
				cmds = append(cmds, cmd)
			}
		} else {
			switch msg.String() {
			case "/":
				// Start search
				hv.searching = true
				hv.searchInput.Focus()
				return hv, textinput.Blink
			case "r":
				// Refresh history
				hv.refresh()
				return hv, nil
			case "d":
				// Delete selected entry
				if selectedItem := hv.list.SelectedItem(); selectedItem != nil {
					historyItem := selectedItem.(HistoryItem)
					hv.manager.Delete(historyItem.entry.ID)
					hv.refresh()
					return hv, nil
				}
			case "c":
				// Clear all history
				hv.manager.Clear()
				hv.refresh()
				return hv, nil
			default:
				hv.list, cmd = hv.list.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}

	return hv, tea.Batch(cmds...)
}

// View renders the history viewer
func (hv HistoryViewer) View() string {
	var sections []string

	// Title
	title := titleStyle.Render("Request History")
	sections = append(sections, title)

	// Search input (if searching)
	if hv.searching {
		searchSection := focusedStyle.Render(fmt.Sprintf("Search: %s", hv.searchInput.View()))
		sections = append(sections, searchSection)
	} else if hv.searchInput.Value() != "" {
		searchSection := blurredStyle.Render(fmt.Sprintf("Search: %s (press / to search again)", hv.searchInput.Value()))
		sections = append(sections, searchSection)
	}

	// List
	sections = append(sections, hv.list.View())

	// Help
	if hv.searching {
		help := helpStyle.Render("Enter to search, Esc to cancel")
		sections = append(sections, help)
	} else {
		help := helpStyle.Render("Enter to select, / to search, r to refresh, d to delete, c to clear all, esc to go back")
		sections = append(sections, help)
	}

	return strings.Join(sections, "\n\n")
}

// GetSelectedEntry returns the currently selected history entry
func (hv HistoryViewer) GetSelectedEntry() *history.HistoryEntry {
	if selectedItem := hv.list.SelectedItem(); selectedItem != nil {
		historyItem := selectedItem.(HistoryItem)
		return &historyItem.entry
	}
	return nil
}

// refresh reloads the history from the manager
func (hv *HistoryViewer) refresh() {
	hv.manager.Load() // Reload from file
	hv.allEntries = hv.manager.GetEntries()
	hv.resetList()
}

// resetList resets the list to show all entries
func (hv *HistoryViewer) resetList() {
	items := make([]list.Item, 0)
	for _, entry := range hv.allEntries {
		items = append(items, HistoryItem{entry: entry})
	}
	hv.list.SetItems(items)
}

// applySearch filters the list based on search input
func (hv *HistoryViewer) applySearch() {
	query := hv.searchInput.Value()
	if query == "" {
		hv.resetList()
		return
	}

	// Search through entries
	filteredEntries := hv.manager.Search(query)

	items := make([]list.Item, 0)
	for _, entry := range filteredEntries {
		items = append(items, HistoryItem{entry: entry})
	}

	hv.list.SetItems(items)
}

// Resize updates the viewer size
func (hv *HistoryViewer) Resize(width, height int) {
	hv.width = width
	hv.height = height
	hv.list.SetSize(width-4, height-8)
	hv.searchInput.Width = width - 10
}

// SaveRequestDialog represents a dialog for saving requests
type SaveRequestDialog struct {
	nameInput        textinput.Model
	descriptionInput textinput.Model
	focusedField     int // 0 = name, 1 = description
	visible          bool
}

// NewSaveRequestDialog creates a new save request dialog
func NewSaveRequestDialog() SaveRequestDialog {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter request name..."
	nameInput.CharLimit = 100
	nameInput.Width = 50
	nameInput.Focus()

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Enter description (optional)..."
	descriptionInput.CharLimit = 200
	descriptionInput.Width = 50

	return SaveRequestDialog{
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		focusedField:     0,
		visible:          false,
	}
}

// Show shows the dialog
func (d *SaveRequestDialog) Show() {
	d.visible = true
	d.focusedField = 0
	d.nameInput.Focus()
	d.descriptionInput.Blur()
}

// Hide hides the dialog
func (d *SaveRequestDialog) Hide() {
	d.visible = false
	d.nameInput.SetValue("")
	d.descriptionInput.SetValue("")
	d.nameInput.Blur()
	d.descriptionInput.Blur()
}

// Update handles dialog updates
func (d SaveRequestDialog) Update(msg tea.Msg) (SaveRequestDialog, tea.Cmd) {
	if !d.visible {
		return d, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if d.focusedField == 0 {
				d.focusedField = 1
				d.nameInput.Blur()
				d.descriptionInput.Focus()
			} else {
				d.focusedField = 0
				d.descriptionInput.Blur()
				d.nameInput.Focus()
			}
			return d, nil
		case "enter":
			// Save the request
			name := d.nameInput.Value()
			description := d.descriptionInput.Value()
			return d, func() tea.Msg {
				return SaveRequestMsg{
					name:        name,
					description: description,
				}
			}
		case "esc":
			d.Hide()
			return d, nil
		}
	}

	// Update focused input
	if d.focusedField == 0 {
		d.nameInput, cmd = d.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		d.descriptionInput, cmd = d.descriptionInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return d, tea.Batch(cmds...)
}

// View renders the dialog
func (d SaveRequestDialog) View() string {
	if !d.visible {
		return ""
	}

	var sections []string

	title := titleStyle.Render("Save Request")
	sections = append(sections, title)

	// Name input
	nameLabel := "Name:"
	var nameSection string
	if d.focusedField == 0 {
		nameSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", nameLabel, d.nameInput.View()))
	} else {
		nameSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", nameLabel, d.nameInput.View()))
	}
	sections = append(sections, nameSection)

	// Description input
	descLabel := "Description:"
	var descSection string
	if d.focusedField == 1 {
		descSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", descLabel, d.descriptionInput.View()))
	} else {
		descSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", descLabel, d.descriptionInput.View()))
	}
	sections = append(sections, descSection)

	// Help
	help := helpStyle.Render("Tab to switch fields, Enter to save, Esc to cancel")
	sections = append(sections, help)

	// Center the dialog
	content := strings.Join(sections, "\n\n")
	return lipgloss.Place(80, 20, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Render(content))
}

// SaveRequestMsg represents a save request message
type SaveRequestMsg struct {
	name        string
	description string
}

// GetName returns the request name
func (msg SaveRequestMsg) GetName() string {
	return msg.name
}

// GetDescription returns the request description
func (msg SaveRequestMsg) GetDescription() string {
	return msg.description
}
