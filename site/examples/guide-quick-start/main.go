// Source: site/src/content/docs/guides/quick-start.mdx:39#972a53f5
package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type model struct {
	modal  teamodal.ConfirmModel
	result string
	width  int
	height int
}

func main() {
	modal := teamodal.NewYesNoModal(
		"Would you like to continue?",
		&teamodal.ConfirmModelArgs{Title: "Quick Start"},
	)
	m := model{modal: modal}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var modal tea.Model

	modal, cmd = m.modal.Update(msg)
	if cmd != nil {
		m.modal = modal.(teamodal.ConfirmModel)
		return m, cmd
	}

	switch msg := msg.(type) {
	case teamodal.AnsweredYesMsg:
		m.result = "You said Yes!"
		return m, nil
	case teamodal.AnsweredNoMsg:
		m.result = "You said No."
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.modal = m.modal.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "space", "enter":
			if !m.modal.IsOpen() {
				m.modal, cmd = m.modal.Open()
				return m, cmd
			}
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	var b strings.Builder
	b.WriteString("Quick Start Example\n\n")
	b.WriteString("Press Space to open the modal, q to quit\n\n")
	if m.result != "" {
		b.WriteString(m.result + "\n")
	}

	view := b.String()
	lines := strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	if m.modal.IsOpen() {
		view = m.modal.OverlayModal(view)
	}

	v := tea.NewView(view)
	v.AltScreen = true
	return v
}
