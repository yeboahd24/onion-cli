package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"onioncli/pkg/collections"
)

// EnvironmentItem represents an environment for the list component
type EnvironmentItem struct {
	environment collections.Environment
}

func (e EnvironmentItem) FilterValue() string {
	return e.environment.Name + " " + e.environment.Description
}

func (e EnvironmentItem) Title() string {
	title := e.environment.Name
	if e.environment.IsActive {
		title += " (Active)"
	}
	return title
}

func (e EnvironmentItem) Description() string {
	varCount := len(e.environment.Variables)
	return fmt.Sprintf("%s (%d variables)", e.environment.Description, varCount)
}

// EnvironmentsViewer handles the environments management interface
type EnvironmentsViewer struct {
	manager      *collections.Manager
	envList      list.Model
	currentView  EnvViewState
	width        int
	height       int
	createDialog CreateEnvironmentDialog
	editDialog   EditEnvironmentDialog
}

// EnvViewState represents the current view state
type EnvViewState int

const (
	ViewEnvironments EnvViewState = iota
	ViewCreateEnvironment
	ViewEditEnvironment
)

// NewEnvironmentsViewer creates a new environments viewer
func NewEnvironmentsViewer(manager *collections.Manager, width, height int) EnvironmentsViewer {
	// Create environments list
	environments := manager.GetEnvironments()
	items := make([]list.Item, len(environments))
	for i, env := range environments {
		items[i] = EnvironmentItem{environment: env}
	}

	envList := list.New(items, list.NewDefaultDelegate(), width-4, height-8)
	envList.Title = "Environments"
	envList.SetShowStatusBar(true)
	envList.SetFilteringEnabled(true)
	envList.SetShowHelp(true)

	return EnvironmentsViewer{
		manager:      manager,
		envList:      envList,
		currentView:  ViewEnvironments,
		width:        width,
		height:       height,
		createDialog: NewCreateEnvironmentDialog(),
		editDialog:   NewEditEnvironmentDialog(),
	}
}

// Update handles environments viewer updates
func (ev EnvironmentsViewer) Update(msg tea.Msg) (EnvironmentsViewer, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle dialogs
	switch ev.currentView {
	case ViewCreateEnvironment:
		ev.createDialog, cmd = ev.createDialog.Update(msg)
		cmds = append(cmds, cmd)
		return ev, tea.Batch(cmds...)
	case ViewEditEnvironment:
		ev.editDialog, cmd = ev.editDialog.Update(msg)
		cmds = append(cmds, cmd)
		return ev, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Create new environment
			ev.currentView = ViewCreateEnvironment
			ev.createDialog.Show()
			return ev, nil

		case "enter", " ":
			// Activate selected environment
			if selectedItem := ev.envList.SelectedItem(); selectedItem != nil {
				envItem := selectedItem.(EnvironmentItem)
				ev.manager.SetActiveEnvironment(envItem.environment.ID)
				ev.refreshEnvironments()
				return ev, func() tea.Msg {
					return EnvironmentChangedMsg{environment: &envItem.environment}
				}
			}

		case "e":
			// Edit selected environment
			if selectedItem := ev.envList.SelectedItem(); selectedItem != nil {
				envItem := selectedItem.(EnvironmentItem)
				ev.editDialog.Show(&envItem.environment)
				ev.currentView = ViewEditEnvironment
				return ev, nil
			}

		case "d":
			// Delete environment (except if it's the only one or active)
			if selectedItem := ev.envList.SelectedItem(); selectedItem != nil {
				envItem := selectedItem.(EnvironmentItem)
				if !envItem.environment.IsActive && len(ev.manager.GetEnvironments()) > 1 {
					// TODO: Implement delete environment
					ev.refreshEnvironments()
				}
				return ev, nil
			}

		case "r":
			// Refresh
			ev.refreshEnvironments()
			return ev, nil
		}

	case CreateEnvironmentMsg:
		// Create new environment
		ev.manager.CreateEnvironment(msg.name, msg.description, msg.variables)
		ev.refreshEnvironments()
		ev.createDialog.Hide()
		ev.currentView = ViewEnvironments
		return ev, nil

	case EditEnvironmentMsg:
		// Update environment
		// TODO: Implement environment update
		ev.refreshEnvironments()
		ev.editDialog.Hide()
		ev.currentView = ViewEnvironments
		return ev, nil
	}

	// Update environment list
	ev.envList, cmd = ev.envList.Update(msg)
	cmds = append(cmds, cmd)

	return ev, tea.Batch(cmds...)
}

// View renders the environments viewer
func (ev EnvironmentsViewer) View() string {
	switch ev.currentView {
	case ViewCreateEnvironment:
		return ev.createDialog.View()
	case ViewEditEnvironment:
		return ev.editDialog.View()
	}

	var sections []string

	// Title
	title := titleStyle.Render("Environment Management")
	sections = append(sections, title)

	// Active environment info
	if activeEnv := ev.manager.GetActiveEnvironment(); activeEnv != nil {
		activeInfo := fmt.Sprintf("Active Environment: %s", activeEnv.Name)
		sections = append(sections, successStyle.Render(activeInfo))
	}

	// Environment list
	sections = append(sections, ev.envList.View())

	// Help
	help := helpStyle.Render("Enter/Space to activate, n to create new, e to edit, d to delete, r to refresh, esc to go back")
	sections = append(sections, help)

	return strings.Join(sections, "\n\n")
}

// refreshEnvironments refreshes the environments list
func (ev *EnvironmentsViewer) refreshEnvironments() {
	environments := ev.manager.GetEnvironments()
	items := make([]list.Item, len(environments))
	for i, env := range environments {
		items[i] = EnvironmentItem{environment: env}
	}
	ev.envList.SetItems(items)
}

// Resize updates the viewer size
func (ev *EnvironmentsViewer) Resize(width, height int) {
	ev.width = width
	ev.height = height
	ev.envList.SetSize(width-4, height-8)
}

// CreateEnvironmentDialog handles creating new environments
type CreateEnvironmentDialog struct {
	nameInput        textinput.Model
	descriptionInput textinput.Model
	variablesInput   textinput.Model
	focusedField     int // 0 = name, 1 = description, 2 = variables
	visible          bool
}

// NewCreateEnvironmentDialog creates a new create environment dialog
func NewCreateEnvironmentDialog() CreateEnvironmentDialog {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter environment name..."
	nameInput.CharLimit = 100
	nameInput.Width = 50
	nameInput.Focus()

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Enter description (optional)..."
	descriptionInput.CharLimit = 200
	descriptionInput.Width = 50

	variablesInput := textinput.New()
	variablesInput.Placeholder = "Variables (key=value, separated by commas)..."
	variablesInput.CharLimit = 500
	variablesInput.Width = 50

	return CreateEnvironmentDialog{
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		variablesInput:   variablesInput,
		focusedField:     0,
		visible:          false,
	}
}

// Show shows the dialog
func (d *CreateEnvironmentDialog) Show() {
	d.visible = true
	d.focusedField = 0
	d.nameInput.Focus()
	d.descriptionInput.Blur()
	d.variablesInput.Blur()
}

// Hide hides the dialog
func (d *CreateEnvironmentDialog) Hide() {
	d.visible = false
	d.nameInput.SetValue("")
	d.descriptionInput.SetValue("")
	d.variablesInput.SetValue("")
	d.nameInput.Blur()
	d.descriptionInput.Blur()
	d.variablesInput.Blur()
}

// Update handles dialog updates
func (d CreateEnvironmentDialog) Update(msg tea.Msg) (CreateEnvironmentDialog, tea.Cmd) {
	if !d.visible {
		return d, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			d.focusedField = (d.focusedField + 1) % 3
			d.updateFocus()
			return d, nil
		case "enter":
			// Create the environment
			name := strings.TrimSpace(d.nameInput.Value())
			if name == "" {
				return d, nil // Don't create without name
			}
			description := strings.TrimSpace(d.descriptionInput.Value())
			variables := d.parseVariables(d.variablesInput.Value())
			return d, func() tea.Msg {
				return CreateEnvironmentMsg{
					name:        name,
					description: description,
					variables:   variables,
				}
			}
		case "esc":
			d.Hide()
			return d, nil
		}
	}

	// Update focused input
	switch d.focusedField {
	case 0:
		d.nameInput, cmd = d.nameInput.Update(msg)
	case 1:
		d.descriptionInput, cmd = d.descriptionInput.Update(msg)
	case 2:
		d.variablesInput, cmd = d.variablesInput.Update(msg)
	}
	cmds = append(cmds, cmd)

	return d, tea.Batch(cmds...)
}

// updateFocus updates the focus state of inputs
func (d *CreateEnvironmentDialog) updateFocus() {
	d.nameInput.Blur()
	d.descriptionInput.Blur()
	d.variablesInput.Blur()

	switch d.focusedField {
	case 0:
		d.nameInput.Focus()
	case 1:
		d.descriptionInput.Focus()
	case 2:
		d.variablesInput.Focus()
	}
}

// parseVariables parses variables from input string
func (d *CreateEnvironmentDialog) parseVariables(input string) map[string]string {
	variables := make(map[string]string)
	if input == "" {
		return variables
	}

	pairs := strings.Split(input, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			if key != "" {
				variables[key] = value
			}
		}
	}

	return variables
}

// View renders the dialog
func (d CreateEnvironmentDialog) View() string {
	if !d.visible {
		return ""
	}

	var sections []string

	title := titleStyle.Render("Create New Environment")
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

	// Variables input
	varLabel := "Variables (key=value, comma separated):"
	var varSection string
	if d.focusedField == 2 {
		varSection = focusedStyle.Render(fmt.Sprintf("%s\n%s", varLabel, d.variablesInput.View()))
	} else {
		varSection = blurredStyle.Render(fmt.Sprintf("%s\n%s", varLabel, d.variablesInput.View()))
	}
	sections = append(sections, varSection)

	// Help
	help := helpStyle.Render("Tab to switch fields, Enter to create, Esc to cancel")
	sections = append(sections, help)

	// Center the dialog
	content := strings.Join(sections, "\n\n")
	return lipgloss.Place(80, 25, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Render(content))
}

// EditEnvironmentDialog handles editing environments (placeholder)
type EditEnvironmentDialog struct {
	visible bool
}

// NewEditEnvironmentDialog creates a new edit environment dialog
func NewEditEnvironmentDialog() EditEnvironmentDialog {
	return EditEnvironmentDialog{visible: false}
}

// Show shows the dialog
func (d *EditEnvironmentDialog) Show(env *collections.Environment) {
	d.visible = true
	// TODO: Implement environment editing
}

// Hide hides the dialog
func (d *EditEnvironmentDialog) Hide() {
	d.visible = false
}

// Update handles dialog updates
func (d EditEnvironmentDialog) Update(msg tea.Msg) (EditEnvironmentDialog, tea.Cmd) {
	// TODO: Implement environment editing
	return d, nil
}

// View renders the dialog
func (d EditEnvironmentDialog) View() string {
	if !d.visible {
		return ""
	}
	return "Environment editing coming soon..."
}

// Message types
type CreateEnvironmentMsg struct {
	name        string
	description string
	variables   map[string]string
}

type EditEnvironmentMsg struct {
	id          string
	name        string
	description string
	variables   map[string]string
}

type EnvironmentChangedMsg struct {
	environment *collections.Environment
}
