package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
)

// formFieldRow is the row where the dropdown field is rendered.
// This must match the row in your View() output where the label appears.
const formFieldRow = 4

// formFieldCol is the column where the dropdown indicator starts.
// This is the character position after "  Category: "
const formFieldCol = 14

type model struct {
	dropdown     teafields.DropdownModel
	selected     string
	hasSelection bool
	width        int
	height       int
}

var (
	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("246")).
			Bold(true)
	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10")).
			Bold(true)
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205"))
	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2)
)

func main() {
	options := teafields.ToOptions([]string{
		"Electronics",
		"Books",
		"Clothing",
		"Home & Garden",
		"Sports",
	})

	teafields.EnsureTermGetSize(os.Stdout.Fd())

	dropdown := teafields.NewDropdownModel(options, &teafields.DropdownModelArgs{
		FieldRow: formFieldRow,
		FieldCol: formFieldCol,
	})

	p := tea.NewProgram(model{dropdown: dropdown})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// --- Modal consumption pattern ---
	// Let the dropdown handle the message first.
	// If it returns a non-nil cmd, it consumed the message.
	dropdown, cmd := m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teafields.DropdownModel)
		return m, cmd
	}

	// --- Dropdown did not consume; handle ourselves ---
	switch msg := msg.(type) {
	case teafields.OptionSelectedMsg:
		m.selected = msg.Text
		m.hasSelection = true
		m.dropdown, _ = m.dropdown.Close()
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dropdown = m.dropdown.WithScreenSize(msg.Width, msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ", "enter":
			if !m.dropdown.IsOpen {
				m.dropdown, cmd = m.dropdown.Open()
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder

	b.WriteString(titleStyle.Render("New Product Form"))
	b.WriteString("\n\n")
	b.WriteString(labelStyle.Render("  Product:") + "  Widget Pro\n")

	// The dropdown field -- indicator shows open/closed state
	var indicator string
	if m.dropdown.IsOpen {
		if m.dropdown.DisplayAbove {
			indicator = "^"
		} else {
			indicator = "v"
		}
	} else {
		indicator = ">"
	}

	var fieldText string
	if m.hasSelection {
		fieldText = valueStyle.Render(m.selected)
	} else {
		fieldText = "Select..."
	}
	b.WriteString(labelStyle.Render("  Category:") + " " + indicator + " " + fieldText)
	b.WriteString("\n")
	b.WriteString(labelStyle.Render("  Price:") + "    $29.99\n")
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("  Space/Enter: open dropdown | q: quit"))

	// Wrap in a border sized to the terminal
	content := b.String()
	if m.width > 0 && m.height > 0 {
		content = borderStyle.
			Width(m.width - 4).
			Height(m.height - 4).
			Render(content)
	} else {
		content = borderStyle.Render(content)
	}

	// Overlay the dropdown on top of the form when open
	if m.dropdown.IsOpen {
		dropdownView := m.dropdown.View()
		// Adjust for border (1 row) + padding (1 row top, 2 cols left)
		content = teafields.OverlayDropdown(
			content,
			dropdownView.Content,
			m.dropdown.Row+2,
			m.dropdown.Col+3,
		)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	return v
}
