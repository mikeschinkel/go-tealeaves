package main

import (
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type model struct {
	confirmDialog   teamodal.ConfirmModel
	alertDialog     teamodal.ConfirmModel
	multilineDialog teamodal.ConfirmModel
	styledDialog    teamodal.ConfirmModel
	currentModal    string // "confirm", "alert", "multiline", or "styled"
	confirmed       bool
	cancelled       bool
	width           int
	height          int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var modal tea.Model

	// Let active modal handle message first
	if m.currentModal == "confirm" {
		modal, cmd = m.confirmDialog.Update(msg)
		if cmd != nil {
			m.confirmDialog = modal.(teamodal.ConfirmModel)
			return m, cmd
		}
	} else if m.currentModal == "alert" {
		modal, cmd = m.alertDialog.Update(msg)
		if cmd != nil {
			m.alertDialog = modal.(teamodal.ConfirmModel)
			return m, cmd
		}
	} else if m.currentModal == "multiline" {
		modal, cmd = m.multilineDialog.Update(msg)
		if cmd != nil {
			m.multilineDialog = modal.(teamodal.ConfirmModel)
			return m, cmd
		}
	} else if m.currentModal == "styled" {
		modal, cmd = m.styledDialog.Update(msg)
		if cmd != nil {
			m.styledDialog = modal.(teamodal.ConfirmModel)
			return m, cmd
		}
	}

	// Modal didn't handle - parent processes
	switch msg := msg.(type) {
	case teamodal.AnsweredYesMsg:
		m.confirmed = true
		m.currentModal = ""
		return m, nil

	case teamodal.AnsweredNoMsg:
		m.cancelled = true
		m.currentModal = ""
		return m, nil

	case teamodal.ClosedMsg:
		m.currentModal = ""
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "y":
			// Open confirmation dialog
			m.currentModal = "confirm"
			m.confirmed = false
			m.cancelled = false
			m.confirmDialog, cmd = m.confirmDialog.Open()
			return m, cmd

		case "o":
			// Open alert dialog
			m.currentModal = "alert"
			m.alertDialog, cmd = m.alertDialog.Open()
			return m, cmd

		case "m":
			// Open multiline dialog
			m.currentModal = "multiline"
			m.multilineDialog, cmd = m.multilineDialog.Open()
			return m, cmd

		case "s":
			// Open styled dialog
			m.currentModal = "styled"
			m.styledDialog, cmd = m.styledDialog.Open()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.confirmDialog = m.confirmDialog.SetSize(msg.Width, msg.Height)
		m.alertDialog = m.alertDialog.SetSize(msg.Width, msg.Height)
		m.multilineDialog = m.multilineDialog.SetSize(msg.Width, msg.Height)
		m.styledDialog = m.styledDialog.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	var view string
	var baseView strings.Builder
	var lines []string

	// Build base view
	baseView.WriteString("TeaModal Example\n")
	baseView.WriteString("===================\n\n")
	baseView.WriteString("Commands:\n")
	baseView.WriteString("  y - Open confirmation dialog (Yes/No)\n")
	baseView.WriteString("  o - Open alert dialog (OK)\n")
	baseView.WriteString("  m - Open multiline dialog\n")
	baseView.WriteString("  s - Open styled dialog (custom colors)\n")
	baseView.WriteString("  q - Quit\n\n")

	if m.confirmed {
		baseView.WriteString("Status: User confirmed! ✓\n")
	} else if m.cancelled {
		baseView.WriteString("Status: User cancelled.\n")
	}

	view = baseView.String()

	// Pad view to fill screen
	lines = strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	// Composite modal if open
	if m.confirmDialog.IsOpen() {
		view = m.confirmDialog.OverlayModal(view)
		goto end
	}

	if m.alertDialog.IsOpen() {
		view = m.alertDialog.OverlayModal(view)
		goto end
	}

	if m.multilineDialog.IsOpen() {
		view = m.multilineDialog.OverlayModal(view)
		goto end
	}

	if m.styledDialog.IsOpen() {
		view = m.styledDialog.OverlayModal(view)
		goto end
	}

end:
	v := tea.NewView(view)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	// Ensure terminal size is available
	teamodal.EnsureTermGetSize(os.Stdout.Fd())

	confirmModal := teamodal.NewYesNoModal("Do you want to proceed with this operation?", &teamodal.ConfirmModelArgs{
		Title:      "Confirmation Required",
		DefaultYes: true,
	})

	alertModal := teamodal.NewOKModal("Operation completed successfully!", &teamodal.ConfirmModelArgs{
		Title: "Success",
	})

	multilineModal := teamodal.NewOKModal("This is a multi-line message.\n\nIt demonstrates that you can use \\n\nto create line breaks in your modal text.", &teamodal.ConfirmModelArgs{
		Title: "Multi-line Example",
	})

	styledModal := teamodal.NewYesNoModal("This modal demonstrates custom styling.\n\nNotice the custom colors for title,\nmessage, and buttons!", &teamodal.ConfirmModelArgs{
		Title:      "Custom Styling",
		DefaultYes: true,
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")). // Hot pink
			Background(lipgloss.Color("235")),
		MessageStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Cyan
			Italic(true),
		ButtonStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Background(lipgloss.Color("252")),
		FocusedButtonStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("230")). // Yellow
			Background(lipgloss.Color("63")),  // Purple
	})

	m := model{
		confirmDialog:   confirmModal,
		alertDialog:     alertModal,
		multilineDialog: multilineModal,
		styledDialog:    styledModal,
		width:           80,
		height:          24,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		cliutil.Stderr("Error: %v\n", err)
		os.Exit(1)
	}
}
