package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

type model struct {
	registry  *teautils.KeyRegistry
	showHelp  bool
	width     int
	height    int
	statusMsg string
}

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	categoryStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginTop(1)
	keyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
	descStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	borderStyle   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2)
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
		registry: registry,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			if !m.showHelp {
				return m, tea.Quit
			}
			m.showHelp = false
			return m, nil

		case "?":
			m.showHelp = !m.showHelp
			return m, nil

		case "esc":
			if m.showHelp {
				m.showHelp = false
				return m, nil
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("TeaUtils Key Help Example"))
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

	// Overlay help modal if open
	if m.showHelp {
		helpContent := m.renderHelpModal()
		_, _, row, col := teautils.CenterModal(helpContent, m.width, m.height)
		view = overlayAt(view, helpContent, row, col)
	}

	return view
}

func (m model) renderHelpModal() string {
	var content strings.Builder

	style := teautils.DefaultHelpVisorStyle()
	byCategory := m.registry.ByCategory()
	sorted := teautils.GetSortedCategories(byCategory, style.CategoryOrder)

	content.WriteString(titleStyle.Render("Keyboard Shortcuts"))
	content.WriteString("\n")

	for _, cat := range sorted {
		keys := byCategory[cat]
		content.WriteString("\n")
		content.WriteString(categoryStyle.Render(cat))
		content.WriteString("\n")
		for _, km := range keys {
			display := teautils.FormatKeyDisplay(km)
			paddedKey := fmt.Sprintf("%-16s", display)
			content.WriteString(fmt.Sprintf("  %s %s\n",
				keyStyle.Render(paddedKey),
				descStyle.Render(km.HelpText),
			))
		}
	}

	content.WriteString("\n")
	content.WriteString(helpStyle.Render("Press ? or Esc to close"))

	return borderStyle.Render(content.String())
}

// overlayAt places overlay text on top of base text at the given row/col
func overlayAt(base, overlay string, row, col int) string {
	baseLines := strings.Split(base, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, overlayLine := range overlayLines {
		targetRow := row + i
		if targetRow < 0 || targetRow >= len(baseLines) {
			continue
		}

		baseLine := baseLines[targetRow]
		baseRunes := []rune(baseLine)

		// Pad base line if needed
		for len(baseRunes) < col {
			baseRunes = append(baseRunes, ' ')
		}

		// Build new line: base prefix + overlay + base suffix
		overlayRunes := []rune(overlayLine)
		endCol := col + len(overlayRunes)

		var result []rune
		result = append(result, baseRunes[:col]...)
		result = append(result, overlayRunes...)
		if endCol < len(baseRunes) {
			result = append(result, baseRunes[endCol:]...)
		}

		baseLines[targetRow] = string(result)
	}

	return strings.Join(baseLines, "\n")
}
