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

// CollectionItem represents a collection for the list component
type CollectionItem struct {
	collection collections.Collection
}

func (c CollectionItem) FilterValue() string {
	return c.collection.Name + " " + c.collection.Description
}

func (c CollectionItem) Title() string {
	return c.collection.Name
}

func (c CollectionItem) Description() string {
	requestCount := len(c.collection.Requests)
	return fmt.Sprintf("%s (%d requests)", c.collection.Description, requestCount)
}

// RequestItem represents a request within a collection
type RequestItem struct {
	request collections.CollectionRequest
}

func (r RequestItem) FilterValue() string {
	return r.request.Name + " " + r.request.Method + " " + r.request.URL
}

func (r RequestItem) Title() string {
	return fmt.Sprintf("%s %s", r.request.Method, r.request.Name)
}

func (r RequestItem) Description() string {
	return r.request.URL
}

// CollectionsViewer handles the collections browsing interface
type CollectionsViewer struct {
	manager            *collections.Manager
	collectionsList    list.Model
	requestsList       list.Model
	currentView        CollectionViewState
	selectedCollection *collections.Collection
	width              int
	height             int
	createDialog       CreateCollectionDialog
}

// CollectionViewState represents the current view state
type CollectionViewState int

const (
	ViewCollections CollectionViewState = iota
	ViewRequests
	ViewCreateCollection
)

// NewCollectionsViewer creates a new collections viewer
func NewCollectionsViewer(manager *collections.Manager, width, height int) CollectionsViewer {
	// Create collections list
	collections := manager.GetCollections()
	items := make([]list.Item, len(collections))
	for i, collection := range collections {
		items[i] = CollectionItem{collection: collection}
	}

	collectionsList := list.New(items, list.NewDefaultDelegate(), width-4, height-8)
	collectionsList.Title = "Collections"
	collectionsList.SetShowStatusBar(true)
	collectionsList.SetFilteringEnabled(true)
	collectionsList.SetShowHelp(true)

	// Create requests list (initially empty)
	requestsList := list.New([]list.Item{}, list.NewDefaultDelegate(), width-4, height-8)
	requestsList.Title = "Requests"
	requestsList.SetShowStatusBar(true)
	requestsList.SetFilteringEnabled(true)
	requestsList.SetShowHelp(true)

	return CollectionsViewer{
		manager:         manager,
		collectionsList: collectionsList,
		requestsList:    requestsList,
		currentView:     ViewCollections,
		width:           width,
		height:          height,
		createDialog:    NewCreateCollectionDialog(),
	}
}

// Update handles collections viewer updates
func (cv CollectionsViewer) Update(msg tea.Msg) (CollectionsViewer, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	// Handle create dialog
	if cv.currentView == ViewCreateCollection {
		cv.createDialog, cmd = cv.createDialog.Update(msg)
		cmds = append(cmds, cmd)
		return cv, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			// Create new collection
			cv.currentView = ViewCreateCollection
			cv.createDialog.Show()
			return cv, nil

		case "enter":
			if cv.currentView == ViewCollections {
				// Open selected collection
				if selectedItem := cv.collectionsList.SelectedItem(); selectedItem != nil {
					collectionItem := selectedItem.(CollectionItem)
					cv.selectedCollection = &collectionItem.collection
					cv.loadRequests()
					cv.currentView = ViewRequests
					return cv, nil
				}
			} else if cv.currentView == ViewRequests {
				// Load selected request
				if selectedItem := cv.requestsList.SelectedItem(); selectedItem != nil {
					requestItem := selectedItem.(RequestItem)
					return cv, func() tea.Msg {
						return LoadRequestMsg{request: &requestItem.request}
					}
				}
			}

		case "esc", "backspace":
			if cv.currentView == ViewRequests {
				cv.currentView = ViewCollections
				cv.selectedCollection = nil
				return cv, nil
			}

		case "d":
			// Delete collection or request
			if cv.currentView == ViewCollections {
				if selectedItem := cv.collectionsList.SelectedItem(); selectedItem != nil {
					collectionItem := selectedItem.(CollectionItem)
					cv.manager.DeleteCollection(collectionItem.collection.ID)
					cv.refreshCollections()
					return cv, nil
				}
			}

		case "r":
			// Refresh
			cv.refreshCollections()
			return cv, nil
		}

	case CreateCollectionMsg:
		// Create new collection
		collection := cv.manager.CreateCollection(msg.name, msg.description)
		cv.refreshCollections()
		cv.createDialog.Hide()
		cv.currentView = ViewCollections

		// Select the new collection
		for i, item := range cv.collectionsList.Items() {
			if collectionItem, ok := item.(CollectionItem); ok {
				if collectionItem.collection.ID == collection.ID {
					cv.collectionsList.Select(i)
					break
				}
			}
		}
		return cv, nil
	}

	// Update current list
	switch cv.currentView {
	case ViewCollections:
		cv.collectionsList, cmd = cv.collectionsList.Update(msg)
		cmds = append(cmds, cmd)
	case ViewRequests:
		cv.requestsList, cmd = cv.requestsList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cv, tea.Batch(cmds...)
}

// View renders the collections viewer
func (cv CollectionsViewer) View() string {
	if cv.currentView == ViewCreateCollection {
		return cv.createDialog.View()
	}

	var sections []string

	// Title
	title := titleStyle.Render("Collections & Requests")
	sections = append(sections, title)

	// Current view content
	switch cv.currentView {
	case ViewCollections:
		sections = append(sections, cv.collectionsList.View())
		help := helpStyle.Render("Enter to open, n to create new, d to delete, r to refresh, esc to go back")
		sections = append(sections, help)

	case ViewRequests:
		if cv.selectedCollection != nil {
			collectionTitle := fmt.Sprintf("Collection: %s", cv.selectedCollection.Name)
			sections = append(sections, lipgloss.NewStyle().Bold(true).Render(collectionTitle))
		}
		sections = append(sections, cv.requestsList.View())
		help := helpStyle.Render("Enter to load request, d to delete, esc to go back to collections")
		sections = append(sections, help)
	}

	return strings.Join(sections, "\n\n")
}

// loadRequests loads requests for the selected collection
func (cv *CollectionsViewer) loadRequests() {
	if cv.selectedCollection == nil {
		return
	}

	items := make([]list.Item, len(cv.selectedCollection.Requests))
	for i, request := range cv.selectedCollection.Requests {
		items[i] = RequestItem{request: request}
	}

	cv.requestsList.SetItems(items)
	cv.requestsList.Title = fmt.Sprintf("Requests in %s", cv.selectedCollection.Name)
}

// refreshCollections refreshes the collections list
func (cv *CollectionsViewer) refreshCollections() {
	cv.manager.LoadCollections()
	collections := cv.manager.GetCollections()
	items := make([]list.Item, len(collections))
	for i, collection := range collections {
		items[i] = CollectionItem{collection: collection}
	}
	cv.collectionsList.SetItems(items)
}

// GetSelectedRequest returns the currently selected request
func (cv CollectionsViewer) GetSelectedRequest() *collections.CollectionRequest {
	if cv.currentView == ViewRequests {
		if selectedItem := cv.requestsList.SelectedItem(); selectedItem != nil {
			requestItem := selectedItem.(RequestItem)
			return &requestItem.request
		}
	}
	return nil
}

// Resize updates the viewer size
func (cv *CollectionsViewer) Resize(width, height int) {
	cv.width = width
	cv.height = height
	cv.collectionsList.SetSize(width-4, height-8)
	cv.requestsList.SetSize(width-4, height-8)
}

// CreateCollectionDialog handles creating new collections
type CreateCollectionDialog struct {
	nameInput        textinput.Model
	descriptionInput textinput.Model
	focusedField     int // 0 = name, 1 = description
	visible          bool
}

// NewCreateCollectionDialog creates a new create collection dialog
func NewCreateCollectionDialog() CreateCollectionDialog {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter collection name..."
	nameInput.CharLimit = 100
	nameInput.Width = 50
	nameInput.Focus()

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Enter description (optional)..."
	descriptionInput.CharLimit = 200
	descriptionInput.Width = 50

	return CreateCollectionDialog{
		nameInput:        nameInput,
		descriptionInput: descriptionInput,
		focusedField:     0,
		visible:          false,
	}
}

// Show shows the dialog
func (d *CreateCollectionDialog) Show() {
	d.visible = true
	d.focusedField = 0
	d.nameInput.Focus()
	d.descriptionInput.Blur()
}

// Hide hides the dialog
func (d *CreateCollectionDialog) Hide() {
	d.visible = false
	d.nameInput.SetValue("")
	d.descriptionInput.SetValue("")
	d.nameInput.Blur()
	d.descriptionInput.Blur()
}

// Update handles dialog updates
func (d CreateCollectionDialog) Update(msg tea.Msg) (CreateCollectionDialog, tea.Cmd) {
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
			// Create the collection
			name := strings.TrimSpace(d.nameInput.Value())
			if name == "" {
				return d, nil // Don't create without name
			}
			description := strings.TrimSpace(d.descriptionInput.Value())
			return d, func() tea.Msg {
				return CreateCollectionMsg{
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
func (d CreateCollectionDialog) View() string {
	if !d.visible {
		return ""
	}

	var sections []string

	title := titleStyle.Render("Create New Collection")
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
	help := helpStyle.Render("Tab to switch fields, Enter to create, Esc to cancel")
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

// CreateCollectionMsg represents a create collection message
type CreateCollectionMsg struct {
	name        string
	description string
}

// LoadRequestMsg represents loading a request from collection
type LoadRequestMsg struct {
	request *collections.CollectionRequest
}
