package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type model struct {
	actionDialog teamodal.ChoiceModel
	statusMsg    string
	width        int
	height       int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var choiceModel tea.Model

	// Let active modal handle message first
	if m.actionDialog.IsOpen() {
		choiceModel, cmd = m.actionDialog.Update(msg)
		if cmd != nil {
			m.actionDialog = choiceModel.(teamodal.ChoiceModel)
			return m, cmd
		}
	}

	// Handle messages that bubble up from modals
	switch msg := msg.(type) {
	case teamodal.ChoiceSelectedMsg:
		m.statusMsg = fmt.Sprintf("Selected: %q (index %d)", msg.OptionID, msg.Index)
		return m, nil

	case teamodal.ChoiceCancelledMsg:
		m.statusMsg = "Cancelled (Esc)"
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "a":
			m.statusMsg = ""
			m.actionDialog, cmd = m.actionDialog.Open()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.actionDialog = m.actionDialog.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() (view string) {
	var baseView strings.Builder
	var lines []string

	baseView.WriteString("TeaModal Vertical Orientation Example\n")
	baseView.WriteString("======================================\n\n")
	baseView.WriteString("Commands:\n")
	baseView.WriteString("  a - Open action picker (vertical buttons)\n")
	baseView.WriteString("  q - Quit\n\n")

	if m.statusMsg != "" {
		baseView.WriteString(fmt.Sprintf("Status: %s\n", m.statusMsg))
	}

	view = baseView.String()

	// Pad view to fill screen
	lines = strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	// Composite modal if open
	if m.actionDialog.IsOpen() {
		view = m.actionDialog.OverlayModal(view)
		goto end
	}

end:
	return view
}

func main() {
	teamodal.EnsureTermGetSize(os.Stdout.Fd())

	actionDialog := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
		Title:   "Choose Action",
		Message: "What would you like to do with this item?",
		Options: []teamodal.ChoiceOption{
			{Label: "Edit", Hotkey: 'e', ID: "edit"},
			{Label: "Duplicate", Hotkey: 'd', ID: "duplicate"},
			{Label: "Archive", Hotkey: 'a', ID: "archive"},
			{Label: "Delete", Hotkey: 'x', ID: "delete"},
		},
		DefaultIndex: 0,
		Orientation:  teamodal.Vertical,
	})

	m := model{
		actionDialog: actionDialog,
		width:        80,
		height:       24,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		cliutil.Stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}
