package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"onioncli/pkg/tui"
)

// DemoModel demonstrates the performance enhancements
type DemoModel struct {
	spinner           tui.LoadingSpinner
	statusIndicator   tui.StatusIndicator
	progressIndicator tui.ProgressIndicator
	keyboardShortcuts tui.KeyboardShortcuts
	currentDemo       int
	maxDemos          int
	width             int
	height            int
}

func NewDemoModel() DemoModel {
	return DemoModel{
		spinner:           tui.NewLoadingSpinner(),
		statusIndicator:   tui.NewStatusIndicator(),
		progressIndicator: tui.NewProgressIndicator(),
		keyboardShortcuts: tui.NewKeyboardShortcuts(),
		currentDemo:       0,
		maxDemos:          5,
		width:             80,
		height:            24,
	}
}

func (m DemoModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

type TickMsg struct{}
type NextDemoMsg struct{}

func (m DemoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "enter", " ":
			return m, func() tea.Msg { return NextDemoMsg{} }
		case "?":
			m.keyboardShortcuts.Toggle()
		}

	case TickMsg:
		// Update status indicator
		m.statusIndicator = m.statusIndicator.Update()

		// Continue ticking
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg {
			return TickMsg{}
		})

	case NextDemoMsg:
		m.currentDemo++
		if m.currentDemo >= m.maxDemos {
			m.currentDemo = 0
		}
		return m.startDemo()
	}

	// Update spinner
	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m DemoModel) startDemo() (DemoModel, tea.Cmd) {
	switch m.currentDemo {
	case 0:
		// Loading spinner demo
		return m, m.spinner.Show("Connecting to Tor network...")
	case 1:
		// Success status demo
		m.spinner.Hide()
		m.statusIndicator.Show("Successfully connected to Tor!", tui.StatusSuccess)
		return m, nil
	case 2:
		// Error status demo
		m.statusIndicator.Show("Failed to connect to .onion service", tui.StatusError)
		return m, nil
	case 3:
		// Warning status demo
		m.statusIndicator.Show("Tor connection is slow, consider checking your network", tui.StatusWarning)
		return m, nil
	case 4:
		// Progress indicator demo
		m.progressIndicator.Show("Processing request history", 10)
		for i := 0; i <= 10; i++ {
			m.progressIndicator.Update(i)
		}
		return m, nil
	default:
		m.currentDemo = 0
		return m.startDemo()
	}
}

func (m DemoModel) View() string {
	if m.keyboardShortcuts.IsVisible() {
		return m.keyboardShortcuts.View()
	}

	var sections []string

	// Title
	title := fmt.Sprintf("OnionCLI Performance & Usability Demo (%d/%d)", m.currentDemo+1, m.maxDemos)
	sections = append(sections, title)
	sections = append(sections, "")

	// Current demo description
	switch m.currentDemo {
	case 0:
		sections = append(sections, "Demo 1: Loading Spinner")
		sections = append(sections, "Shows animated spinner during long operations like Tor connections")
		sections = append(sections, "")
		if m.spinner.IsVisible() {
			sections = append(sections, m.spinner.View())
		}

	case 1:
		sections = append(sections, "Demo 2: Success Status Indicator")
		sections = append(sections, "Shows success messages with auto-hide after 5 seconds")
		sections = append(sections, "")
		if m.statusIndicator.IsVisible() {
			sections = append(sections, m.statusIndicator.View())
		}

	case 2:
		sections = append(sections, "Demo 3: Error Status Indicator")
		sections = append(sections, "Shows error messages with appropriate styling")
		sections = append(sections, "")
		if m.statusIndicator.IsVisible() {
			sections = append(sections, m.statusIndicator.View())
		}

	case 3:
		sections = append(sections, "Demo 4: Warning Status Indicator")
		sections = append(sections, "Shows warning messages for non-critical issues")
		sections = append(sections, "")
		if m.statusIndicator.IsVisible() {
			sections = append(sections, m.statusIndicator.View())
		}

	case 4:
		sections = append(sections, "Demo 5: Progress Indicator")
		sections = append(sections, "Shows progress for operations with known duration")
		sections = append(sections, "")
		if m.progressIndicator.IsVisible() {
			sections = append(sections, m.progressIndicator.View())
		}
	}

	sections = append(sections, "")
	sections = append(sections, "Performance Features:")
	sections = append(sections, "• Animated loading spinners for better user feedback")
	sections = append(sections, "• Smart status indicators with auto-hide")
	sections = append(sections, "• Progress bars for long operations")
	sections = append(sections, "• Keyboard shortcuts for power users")
	sections = append(sections, "• Optimized for Tor's higher latency")
	sections = append(sections, "• Retry functionality for failed requests")
	sections = append(sections, "")

	sections = append(sections, "Keyboard Shortcuts:")
	sections = append(sections, "• Space/Enter/n: Next demo")
	sections = append(sections, "• ?: Toggle keyboard shortcuts help")
	sections = append(sections, "• q/Ctrl+C: Quit")
	sections = append(sections, "")

	sections = append(sections, "In the main application:")
	sections = append(sections, "• Tab/Shift+Tab: Navigate fields")
	sections = append(sections, "• h: View history")
	sections = append(sections, "• a: Configure authentication")
	sections = append(sections, "• s: Save request")
	sections = append(sections, "• r: Retry last request")
	sections = append(sections, "• e: View detailed error information")
	sections = append(sections, "• Ctrl+S: Quick save")

	return fmt.Sprintf("%s\n\nPress Space/Enter/n for next demo, ? for help, q to quit",
		strings.Join(sections, "\n"))
}

func main() {
	fmt.Println("OnionCLI Performance & Usability Demo")
	fmt.Println("=====================================")
	fmt.Println()
	fmt.Println("This demo showcases the performance and usability enhancements:")
	fmt.Println("• Loading spinners and status indicators")
	fmt.Println("• Keyboard shortcuts and quick actions")
	fmt.Println("• Progress indicators and user feedback")
	fmt.Println("• Optimizations for Tor network latency")
	fmt.Println()
	fmt.Println("Starting interactive demo...")
	time.Sleep(2 * time.Second)

	p := tea.NewProgram(NewDemoModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running demo: %v\n", err)
	}

	fmt.Println("\n🎉 Performance demo completed!")
	fmt.Println("The main OnionCLI application now includes all these enhancements")
	fmt.Println("Run './onioncli' to experience the improved interface")
}
