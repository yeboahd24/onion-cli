package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"onioncli/pkg/api"
)

// AuthDialog handles authentication configuration
type AuthDialog struct {
	visible      bool
	authManager  *api.AuthManager
	authTypeList list.Model
	inputs       map[string]textinput.Model
	currentStep  int // 0 = select type, 1+ = input fields
	authConfig   *api.AuthConfig
	width        int
	height       int
}

// AuthTypeItem represents an auth type for the list
type AuthTypeItem struct {
	authType api.AuthType
	manager  *api.AuthManager
}

func (a AuthTypeItem) FilterValue() string {
	return string(a.authType)
}

func (a AuthTypeItem) Title() string {
	return string(a.authType)
}

func (a AuthTypeItem) Description() string {
	return a.manager.GetAuthTypeDescription(a.authType)
}

// NewAuthDialog creates a new authentication dialog
func NewAuthDialog(width, height int) AuthDialog {
	authManager := api.NewAuthManager()

	// Create auth type list
	authTypes := authManager.GetAuthTypes()
	items := make([]list.Item, len(authTypes))
	for i, authType := range authTypes {
		items[i] = AuthTypeItem{authType: authType, manager: authManager}
	}

	authTypeList := list.New(items, list.NewDefaultDelegate(), width-10, 8)
	authTypeList.Title = "Select Authentication Type"
	authTypeList.SetShowStatusBar(false)
	authTypeList.SetFilteringEnabled(false)
	authTypeList.SetShowHelp(false)

	// Create input fields
	inputs := make(map[string]textinput.Model)

	// API Key inputs
	apiKeyInput := textinput.New()
	apiKeyInput.Placeholder = "Enter API key..."
	apiKeyInput.Width = width - 20
	inputs["api_key"] = apiKeyInput

	keyNameInput := textinput.New()
	keyNameInput.Placeholder = "Header/parameter name (default: X-API-Key)"
	keyNameInput.Width = width - 20
	inputs["key_name"] = keyNameInput

	locationInput := textinput.New()
	locationInput.Placeholder = "Location: header or query (default: header)"
	locationInput.Width = width - 20
	inputs["location"] = locationInput

	// Bearer token input
	tokenInput := textinput.New()
	tokenInput.Placeholder = "Enter bearer token..."
	tokenInput.Width = width - 20
	inputs["token"] = tokenInput

	// Basic auth inputs
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Enter username..."
	usernameInput.Width = width - 20
	inputs["username"] = usernameInput

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password..."
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Width = width - 20
	inputs["password"] = passwordInput

	// Custom headers input
	headersInput := textinput.New()
	headersInput.Placeholder = "Enter custom headers (key: value, one per line)..."
	headersInput.Width = width - 20
	inputs["headers"] = headersInput

	return AuthDialog{
		visible:      false,
		authManager:  authManager,
		authTypeList: authTypeList,
		inputs:       inputs,
		currentStep:  0,
		width:        width,
		height:       height,
	}
}

// Show displays the auth dialog
func (ad *AuthDialog) Show() {
	ad.visible = true
	ad.currentStep = 0
	ad.authConfig = nil

	// Reset all inputs
	for _, input := range ad.inputs {
		input.SetValue("")
		input.Blur()
	}
}

// Hide hides the auth dialog
func (ad *AuthDialog) Hide() {
	ad.visible = false
	ad.currentStep = 0
	ad.authConfig = nil
}

// Update handles auth dialog updates
func (ad AuthDialog) Update(msg tea.Msg) (AuthDialog, tea.Cmd) {
	if !ad.visible {
		return ad, nil
	}

	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			ad.Hide()
			return ad, nil

		case "enter":
			if ad.currentStep == 0 {
				// Auth type selected, move to input fields
				if selectedItem := ad.authTypeList.SelectedItem(); selectedItem != nil {
					authTypeItem := selectedItem.(AuthTypeItem)
					ad.currentStep = 1

					// Focus first input for selected auth type
					ad.focusFirstInput(authTypeItem.authType)
				}
				return ad, nil
			} else {
				// Complete authentication setup
				return ad.completeAuth()
			}

		case "tab":
			if ad.currentStep > 0 {
				ad.focusNextInput()
				return ad, nil
			}
		}
	}

	// Update current component
	if ad.currentStep == 0 {
		ad.authTypeList, cmd = ad.authTypeList.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Update focused input
		for name, input := range ad.inputs {
			if input.Focused() {
				ad.inputs[name], cmd = input.Update(msg)
				cmds = append(cmds, cmd)
				break
			}
		}
	}

	return ad, tea.Batch(cmds...)
}

// View renders the auth dialog
func (ad AuthDialog) View() string {
	if !ad.visible {
		return ""
	}

	var sections []string

	title := titleStyle.Render("Authentication Setup")
	sections = append(sections, title)

	if ad.currentStep == 0 {
		// Show auth type selection
		sections = append(sections, ad.authTypeList.View())
		help := helpStyle.Render("↑/↓ to select, Enter to confirm, Esc to cancel")
		sections = append(sections, help)
	} else {
		// Show input fields based on selected auth type
		if selectedItem := ad.authTypeList.SelectedItem(); selectedItem != nil {
			authTypeItem := selectedItem.(AuthTypeItem)
			sections = append(sections, ad.renderInputFields(authTypeItem.authType))
		}
	}

	// Center the dialog
	content := strings.Join(sections, "\n\n")
	return lipgloss.Place(ad.width, ad.height, lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Render(content))
}

// renderInputFields renders input fields for the selected auth type
func (ad AuthDialog) renderInputFields(authType api.AuthType) string {
	var sections []string

	typeTitle := fmt.Sprintf("Configure %s Authentication", authType)
	sections = append(sections, lipgloss.NewStyle().Bold(true).Render(typeTitle))

	switch authType {
	case api.AuthNone:
		sections = append(sections, "No authentication required.")

	case api.AuthAPIKey:
		sections = append(sections, ad.renderInput("api_key", "API Key:"))
		sections = append(sections, ad.renderInput("key_name", "Key Name:"))
		sections = append(sections, ad.renderInput("location", "Location:"))

	case api.AuthBearer:
		sections = append(sections, ad.renderInput("token", "Bearer Token:"))

	case api.AuthBasic:
		sections = append(sections, ad.renderInput("username", "Username:"))
		sections = append(sections, ad.renderInput("password", "Password:"))

	case api.AuthCustom:
		sections = append(sections, ad.renderInput("headers", "Custom Headers:"))
	}

	help := helpStyle.Render("Tab to switch fields, Enter to save, Esc to cancel")
	sections = append(sections, help)

	return strings.Join(sections, "\n\n")
}

// renderInput renders a single input field
func (ad AuthDialog) renderInput(name, label string) string {
	input, exists := ad.inputs[name]
	if !exists {
		return ""
	}

	var style lipgloss.Style
	if input.Focused() {
		style = focusedStyle
	} else {
		style = blurredStyle
	}

	return style.Render(fmt.Sprintf("%s\n%s", label, input.View()))
}

// focusFirstInput focuses the first input for the given auth type
func (ad *AuthDialog) focusFirstInput(authType api.AuthType) {
	// Blur all inputs first
	for name, input := range ad.inputs {
		input.Blur()
		ad.inputs[name] = input
	}

	// Focus the first relevant input
	switch authType {
	case api.AuthAPIKey:
		input := ad.inputs["api_key"]
		input.Focus()
		ad.inputs["api_key"] = input
	case api.AuthBearer:
		input := ad.inputs["token"]
		input.Focus()
		ad.inputs["token"] = input
	case api.AuthBasic:
		input := ad.inputs["username"]
		input.Focus()
		ad.inputs["username"] = input
	case api.AuthCustom:
		input := ad.inputs["headers"]
		input.Focus()
		ad.inputs["headers"] = input
	}
}

// focusNextInput moves focus to the next input field
func (ad *AuthDialog) focusNextInput() {
	if selectedItem := ad.authTypeList.SelectedItem(); selectedItem != nil {
		authTypeItem := selectedItem.(AuthTypeItem)

		var inputOrder []string
		switch authTypeItem.authType {
		case api.AuthAPIKey:
			inputOrder = []string{"api_key", "key_name", "location"}
		case api.AuthBasic:
			inputOrder = []string{"username", "password"}
		case api.AuthCustom:
			inputOrder = []string{"headers"}
		default:
			return
		}

		// Find currently focused input and move to next
		for i, name := range inputOrder {
			if ad.inputs[name].Focused() {
				input := ad.inputs[name]
				input.Blur()
				ad.inputs[name] = input

				nextIndex := (i + 1) % len(inputOrder)
				nextInput := ad.inputs[inputOrder[nextIndex]]
				nextInput.Focus()
				ad.inputs[inputOrder[nextIndex]] = nextInput
				break
			}
		}
	}
}

// completeAuth completes the authentication setup
func (ad AuthDialog) completeAuth() (AuthDialog, tea.Cmd) {
	if selectedItem := ad.authTypeList.SelectedItem(); selectedItem != nil {
		authTypeItem := selectedItem.(AuthTypeItem)

		// Collect input values
		inputs := make(map[string]string)
		for name, input := range ad.inputs {
			inputs[name] = input.Value()
		}

		// Create auth config
		config, err := ad.authManager.CreateAuthConfigFromInput(authTypeItem.authType, inputs)
		if err != nil {
			return ad, func() tea.Msg {
				return AuthErrorMsg{err: err}
			}
		}

		ad.authConfig = config
		ad.Hide()

		return ad, func() tea.Msg {
			return AuthConfiguredMsg{config: config}
		}
	}

	return ad, nil
}

// GetAuthConfig returns the configured auth config
func (ad AuthDialog) GetAuthConfig() *api.AuthConfig {
	return ad.authConfig
}

// Resize updates the dialog size
func (ad *AuthDialog) Resize(width, height int) {
	ad.width = width
	ad.height = height
	ad.authTypeList.SetSize(width-10, 8)

	for name, input := range ad.inputs {
		input.Width = width - 20
		ad.inputs[name] = input
	}
}

// AuthConfiguredMsg represents a successful auth configuration
type AuthConfiguredMsg struct {
	config *api.AuthConfig
}

// AuthErrorMsg represents an auth configuration error
type AuthErrorMsg struct {
	err error
}
