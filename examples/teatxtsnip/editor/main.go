package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teatxtsnip"
)

type model struct {
	editor   teatxtsnip.Model
	width    int
	height   int
	quitting bool
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

func main() {
	editor := teatxtsnip.NewTextSnipModel(nil)
	editor.SetWidth(76)
	editor.SetHeight(15)
	editor.ShowLineNumbers = true
	editor.SetValue("Welcome to the TeaTextSel editor example!\n\nThis demonstrates text selection and clipboard operations.\nTry these keyboard shortcuts:\n\n  Shift+Arrow keys  - Select text\n  Ctrl+A            - Select all\n  Ctrl+C            - Copy selection\n  Ctrl+X            - Cut selection\n  Ctrl+V            - Paste\n  Esc               - Clear selection\n\nEdit this text freely. Use arrow keys to navigate.")
	editor.Focus()

	m := model{
		editor: editor,
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
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		editorWidth := msg.Width - 4
		if editorWidth > 120 {
			editorWidth = 120
		}
		editorHeight := msg.Height - 8
		if editorHeight < 5 {
			editorHeight = 5
		}
		m.editor.SetWidth(editorWidth)
		m.editor.SetHeight(editorHeight)
		return m, nil

	case tea.KeyPressMsg:
		// Ctrl+Q to quit (since Ctrl+C is used for copy)
		if msg.String() == "ctrl+q" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.editor, cmd = m.editor.Update(msg)
	return m, cmd
}

func (m model) View() tea.View {
	if m.quitting {
		return tea.NewView("")
	}

	var b strings.Builder

	b.WriteString("TeaTextSel Editor Example\n")
	b.WriteString("==========================\n\n")

	b.WriteString(m.editor.View().Content)
	b.WriteString("\n\n")

	// Status line
	sel := m.editor.Selection()
	if m.editor.HasSelection() {
		selectedText := m.editor.SelectedText()
		charCount := len([]rune(selectedText))
		b.WriteString(helpStyle.Render(fmt.Sprintf(
			"Selection: (%d,%d)-(%d,%d) | %d chars selected | Ctrl+Q quit",
			sel.Start.Row, sel.Start.Col,
			sel.End.Row, sel.End.Col,
			charCount,
		)))
	} else {
		b.WriteString(helpStyle.Render("Shift+Arrows: select | Ctrl+A: select all | Ctrl+C/X/V: copy/cut/paste | Ctrl+Q: quit"))
	}

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
