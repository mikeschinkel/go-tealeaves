package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
)

type model struct {
	statusBar teastatus.Model
	mode      string
	width     int
	height    int
}

func main() {
	helpBinding := key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "Help"))
	tabBinding := key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "Switch"))
	quitBinding := key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit"))

	menuItems := []teastatus.MenuItem{
		teastatus.NewMenuItemFromBinding(helpBinding, "Help"),
		teastatus.NewMenuItemFromBinding(tabBinding, "Switch"),
		teastatus.NewMenuItemFromBinding(quitBinding, "Quit"),
	}

	indicators := []teastatus.StatusIndicator{
		teastatus.NewStatusIndicator("Ready"),
		teastatus.NewStatusIndicator("3 files"),
	}

	sb := teastatus.New().
		SetMenuItems(menuItems).
		SetIndicators(indicators)

	m := model{
		statusBar: sb,
		mode:      "normal",
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
		m.statusBar = m.statusBar.SetSize(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "1":
			// Pipe separator
			styles := teastatus.DefaultStyles()
			styles.SeparatorKind = teastatus.PipeSeparator
			m.statusBar = m.statusBar.WithStyles(styles)
			m.mode = "pipe"
			return m, nil

		case "2":
			// Space separator
			styles := teastatus.DefaultStyles()
			styles.SeparatorKind = teastatus.SpaceSeparator
			m.statusBar = m.statusBar.WithStyles(styles)
			m.mode = "space"
			return m, nil

		case "3":
			// Bracket separator
			styles := teastatus.DefaultStyles()
			styles.SeparatorKind = teastatus.BracketSeparator
			m.statusBar = m.statusBar.WithStyles(styles)
			m.mode = "bracket"
			return m, nil

		case "a":
			// Add indicator
			m.statusBar = m.statusBar.SetIndicators([]teastatus.StatusIndicator{
				teastatus.NewStatusIndicator("Processing").
					WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("226"))),
				teastatus.NewStatusIndicator("5 files"),
				teastatus.NewStatusIndicator("MODIFIED"),
			})
			return m, nil

		case "r":
			// Reset indicators
			m.statusBar = m.statusBar.SetIndicators([]teastatus.StatusIndicator{
				teastatus.NewStatusIndicator("Ready"),
				teastatus.NewStatusIndicator("3 files"),
			})
			m.mode = "normal"
			return m, nil
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content strings.Builder
	content.WriteString("TeaStatus Status Bar Example\n")
	content.WriteString("==============================\n\n")
	content.WriteString("Commands:\n")
	content.WriteString("  1 - Pipe separator style\n")
	content.WriteString("  2 - Space separator style\n")
	content.WriteString("  3 - Bracket separator style\n")
	content.WriteString("  a - Add more indicators\n")
	content.WriteString("  r - Reset to defaults\n")
	content.WriteString("  q - Quit\n\n")
	content.WriteString(fmt.Sprintf("Current mode: %s\n", m.mode))

	// Pad to fill screen (leave room for status bar)
	lines := strings.Split(content.String(), "\n")
	for len(lines) < m.height-1 {
		lines = append(lines, "")
	}

	// Status bar at bottom
	lines = append(lines, m.statusBar.View())

	return strings.Join(lines, "\n")
}
