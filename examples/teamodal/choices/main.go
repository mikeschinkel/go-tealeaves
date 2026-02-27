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
	exitDialog   teamodal.ChoiceModel
	deleteDialog teamodal.ChoiceModel
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
	if m.exitDialog.IsOpen() {
		choiceModel, cmd = m.exitDialog.Update(msg)
		if cmd != nil {
			m.exitDialog = choiceModel.(teamodal.ChoiceModel)
			return m, cmd
		}
	}

	if m.deleteDialog.IsOpen() {
		choiceModel, cmd = m.deleteDialog.Update(msg)
		if cmd != nil {
			m.deleteDialog = choiceModel.(teamodal.ChoiceModel)
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

		case "e":
			m.statusMsg = ""
			m.exitDialog, cmd = m.exitDialog.Open()
			return m, cmd

		case "d":
			m.statusMsg = ""
			m.deleteDialog, cmd = m.deleteDialog.Open()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.exitDialog = m.exitDialog.SetSize(msg.Width, msg.Height)
		m.deleteDialog = m.deleteDialog.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() (view string) {
	var baseView strings.Builder
	var lines []string

	baseView.WriteString("TeaModal Choice Example\n")
	baseView.WriteString("==========================\n\n")
	baseView.WriteString("Commands:\n")
	baseView.WriteString("  e - Exit with pending changes (3 options + hotkeys)\n")
	baseView.WriteString("  d - Delete confirmation (2 options)\n")
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
	if m.exitDialog.IsOpen() {
		view = m.exitDialog.OverlayModal(view)
		goto end
	}

	if m.deleteDialog.IsOpen() {
		view = m.deleteDialog.OverlayModal(view)
		goto end
	}

end:
	return view
}

func main() {
	teamodal.EnsureTermGetSize(os.Stdout.Fd())

	exitDialog := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
		Title:   "Exit with Pending Changes",
		Message: "Some files have been reassigned:",
		Options: []teamodal.ChoiceOption{
			{Label: "Reorg & exit", Hotkey: 'r', ID: "reorganize"},
			{Label: "Save & exit", Hotkey: 's', ID: "save"},
			{Label: "Cancel", Hotkey: 'c', ID: "cancel"},
		},
		DefaultIndex: 0,
	})

	deleteDialog := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
		Title:   "Confirm Delete",
		Message: "Are you sure you want to delete this item?\nThis action cannot be undone.",
		Options: []teamodal.ChoiceOption{
			{Label: "Delete", Hotkey: 'd', ID: "delete"},
			{Label: "Keep", Hotkey: 'k', ID: "keep"},
		},
		DefaultIndex: 1,
	})

	m := model{
		exitDialog:   exitDialog,
		deleteDialog: deleteDialog,
		width:        80,
		height:       24,
	}

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())

	if _, err := p.Run(); err != nil {
		cliutil.Stderr("Error: %v\n", err)
		os.Exit(1)
	}
}
