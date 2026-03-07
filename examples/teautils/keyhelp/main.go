package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teahelp"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

type model struct {
	registry  *teautils.KeyRegistry
	helpVisor teahelp.HelpVisorModel
	width     int
	height    int
}

var (
	pageTitleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func main() {
	registry := teautils.NewKeyRegistry()

	// Register navigation keys
	registry.MustRegisterMany([]teautils.KeyMeta{
		{
			ID:        teautils.MustParseKeyIdentifier("nav.up"),
			Binding:   key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("up/k", "Move up")),
			HelpModal: true,
			Category:  "Navigation",
			HelpText:  "Move cursor up one line",
		},
		{
			ID:        teautils.MustParseKeyIdentifier("nav.down"),
			Binding:   key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("down/j", "Move down")),
			HelpModal: true,
			Category:  "Navigation",
			HelpText:  "Move cursor down one line",
		},
		{
			ID:        teautils.MustParseKeyIdentifier("nav.left"),
			Binding:   key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("left/h", "Collapse")),
			HelpModal: true,
			Category:  "Navigation",
			HelpText:  "Collapse node or move to parent",
		},
		{
			ID:        teautils.MustParseKeyIdentifier("nav.right"),
			Binding:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("right/l", "Expand")),
			HelpModal: true,
			Category:  "Navigation",
			HelpText:  "Expand node or move to first child",
		},
	})

	// Register action keys
	registry.MustRegisterMany([]teautils.KeyMeta{
		{
			ID:             teautils.MustParseKeyIdentifier("action.open"),
			Binding:        key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Open")),
			StatusBar:      true,
			StatusBarLabel: "Open",
			HelpModal:      true,
			Category:       "Actions",
			HelpText:       "Open the selected file",
		},
		{
			ID:             teautils.MustParseKeyIdentifier("action.delete"),
			Binding:        key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "Delete")),
			StatusBar:      true,
			StatusBarLabel: "Delete",
			HelpModal:      true,
			Category:       "Actions",
			HelpText:       "Delete the selected item",
		},
		{
			ID:        teautils.MustParseKeyIdentifier("action.rename"),
			Binding:   key.NewBinding(key.WithKeys("r"), key.WithHelp("r", "Rename")),
			HelpModal: true,
			Category:  "Actions",
			HelpText:  "Rename the selected item",
		},
	})

	// Register system keys
	registry.MustRegisterMany([]teautils.KeyMeta{
		{
			ID:             teautils.MustParseKeyIdentifier("sys.help"),
			Binding:        key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "Help")),
			StatusBar:      true,
			StatusBarLabel: "Help",
			HelpModal:      true,
			Category:       "System",
			HelpText:       "Toggle help modal",
		},
		{
			ID:             teautils.MustParseKeyIdentifier("sys.quit"),
			Binding:        key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit")),
			StatusBar:      true,
			StatusBarLabel: "Quit",
			HelpModal:      true,
			Category:       "System",
			HelpText:       "Quit the application",
		},
	})

	m := model{
		registry:  registry,
		helpVisor: teahelp.NewHelpVisorModel(),
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.helpVisor = m.helpVisor.SetSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		// When visor is open, delegate all keys to it
		if m.helpVisor.IsOpen() {
			var cmd tea.Cmd
			m.helpVisor, cmd = m.helpVisor.Update(msg)
			return m, cmd
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.helpVisor = m.helpVisor.Open(m.registry.ByCategory())
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	var b strings.Builder
	b.WriteString(pageTitleStyle.Render("TeaUtils Key Help Example"))
	b.WriteString("\n\n")
	b.WriteString("This example demonstrates the KeyRegistry system.\n")
	b.WriteString("Press ? to toggle the help modal.\n\n")

	// Show status bar keys
	b.WriteString("Status bar keys:\n")
	for _, km := range m.registry.ForStatusBar() {
		display := teautils.FormatKeyDisplay(km)
		label := km.StatusBarLabel
		if label == "" {
			label = km.Binding.Help().Desc
		}
		b.WriteString(fmt.Sprintf("  [%s] %s\n", display, label))
	}

	// Show categories
	b.WriteString("\nRegistered categories:\n")
	byCategory := m.registry.ByCategory()
	sorted := teautils.GetSortedCategories(byCategory, []string{"Navigation", "Actions", "System"})
	for _, cat := range sorted {
		keys := byCategory[cat]
		b.WriteString(fmt.Sprintf("  %s: %d keys\n", cat, len(keys)))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press ? for help | q to quit"))

	view := b.String()

	// Pad to fill screen
	lines := strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	// Overlay help visor if open
	if m.helpVisor.IsOpen() {
		helpView := m.helpVisor.View().Content
		_, _, row, col := teautils.CenterModal(helpView, m.width, m.height)
		view = teamodal.OverlayModal(view, helpView, row, col)
	}

	v := tea.NewView(view)
	v.AltScreen = true
	return v
}
