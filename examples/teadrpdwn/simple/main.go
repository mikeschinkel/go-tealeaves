package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teadrpdwn"
)

type model struct {
	dropdown     teadrpdwn.DropdownModel
	selected     string
	quitting     bool
	hasSelection bool
	screenWidth  int
	screenHeight int
}

var (
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true) // Bright green
	borderStyle   = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)
)

func main() {
	items := teadrpdwn.ToOptions([]string{"Apple", "Banana", "Cherry", "Date", "Elderberry"})

	// Ensure that term.GetSize() is initialized before continuing.
	// This is needed in GoLand terminal for debugging, but is not harmful if not needed.
	teadrpdwn.EnsureTermGetSize(os.Stdout.Fd())

	// Create dropdown at row 3, col 18 (after "     Fruit Selected: ")
	// Screen dimensions will be set automatically from tea.WindowSizeMsg
	dropdown := teadrpdwn.NewDropdownModel(items, &teadrpdwn.DropdownModelArgs{
		FieldRow: 3,
		FieldCol: 18,
	})

	initialModel := model{
		dropdown:     dropdown,
		selected:     "",
		hasSelection: false,
	}

	p := tea.NewProgram(initialModel)

	finalModel, err := p.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Show selected value when exiting
	m := finalModel.(model)
	if m.hasSelection {
		fmt.Printf("You selected fruit %s\n", selectedStyle.Render(m.selected))
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Let dropdown handle message first
	dropdown, cmd := m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teadrpdwn.DropdownModel)
		return m, cmd
	}

	// Dropdown didn't handle - process message
	switch msg := msg.(type) {
	case teadrpdwn.OptionSelectedMsg:
		m.selected = msg.Text
		m.hasSelection = true
		m.dropdown, _ = m.dropdown.Close()
		return m, nil

	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
		m.screenHeight = msg.Height
		m.dropdown = m.dropdown.WithScreenSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "space", "enter":
			if !m.dropdown.IsOpen {
				m.dropdown, cmd = m.dropdown.Open()
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	// Build interior content
	var content strings.Builder

	content.WriteString("TeaDD Dropdown Example - Full Screen TUI\n")
	content.WriteString("Press Space/Enter to open dropdown, q to quit\n")
	content.WriteString("\n")

	// Field with triangle indicator
	var fieldSymbol string
	if m.dropdown.IsOpen {
		if m.dropdown.DisplayAbove {
			fieldSymbol = "▲"
		} else {
			fieldSymbol = "▼"
		}
	} else {
		fieldSymbol = "▶"
	}

	// Show "Select Fruit" or selected value in bright green
	var fieldText string
	if m.hasSelection {
		fieldText = selectedStyle.Render(m.selected)
	} else {
		fieldText = "Select Fruit"
	}
	content.WriteString(fmt.Sprintf("Fruit Selected: %s %s", fieldSymbol, fieldText))

	// Apply border sized to fill screen
	baseView := content.String()
	if m.screenWidth > 0 && m.screenHeight > 0 {
		// Size the border to fill the screen
		baseView = borderStyle.
			Width(m.screenWidth - 4).   // -4 for border and padding
			Height(m.screenHeight - 4). // -4 for border and padding
			Render(content.String())
	} else {
		baseView = borderStyle.Render(content.String())
	}

	// Overlay dropdown if open
	if m.dropdown.IsOpen {
		dropdownView := m.dropdown.View().Content
		// Adjust for border (1 row) + padding (1 row top, 2 cols left)
		baseView = teadrpdwn.OverlayDropdown(baseView, dropdownView, m.dropdown.Row+2, m.dropdown.Col+3)
	}

	v := tea.NewView(baseView)
	v.AltScreen = true
	return v
}
