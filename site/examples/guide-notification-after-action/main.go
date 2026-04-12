// Source: site/src/content/docs/cookbook/notification-after-action.mdx:23
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teanotify"
)

// saveCompleteMsg is returned by the async save command.
type saveCompleteMsg struct{ filename string }

// saveFailedMsg is returned when the save fails.
type saveFailedMsg struct{ err error }

type model struct {
	notify teanotify.NotifyModel
	saving bool
	width  int
	height int
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	keyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
)

func main() {
	// Create the notification model with options
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:           40,
		Duration:        3 * time.Second,
		Position:        teanotify.TopRightPosition,
		AllowEscToClose: true,
	})

	// Initialize registers default notice types (Info, Warn, Error, Debug)
	var err error
	notify, err = notify.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	p := tea.NewProgram(model{notify: notify})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// simulateSave returns a tea.Cmd that simulates an async file save.
func simulateSave(filename string) tea.Cmd {
	return func() tea.Msg {
		// Simulate network or disk latency
		time.Sleep(1 * time.Second)
		return saveCompleteMsg{filename: filename}
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Always let teanotify process the message first.
	// This handles tick messages for animation and ESC for dismissal.
	var notifyCmd tea.Cmd
	m.notify, notifyCmd = m.notify.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, notifyCmd

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "s":
			if !m.saving {
				m.saving = true
				return m, tea.Batch(notifyCmd, simulateSave("document.txt"))
			}
		}

	// --- Handle async completion messages ---
	case saveCompleteMsg:
		m.saving = false
		// Trigger an info notification
		cmd := m.notify.NewNotifyCmd(
			teanotify.InfoKey,
			fmt.Sprintf("Saved %s successfully", msg.filename),
		)
		return m, tea.Batch(notifyCmd, cmd)

	case saveFailedMsg:
		m.saving = false
		// Trigger an error notification
		cmd := m.notify.NewNotifyCmd(
			teanotify.ErrorKey,
			fmt.Sprintf("Save failed: %v", msg.err),
		)
		return m, tea.Batch(notifyCmd, cmd)
	}

	return m, notifyCmd
}

func (m model) View() tea.View {
	if m.width == 0 {
		v := tea.NewView("Loading...")
		v.AltScreen = true
		return v
	}

	var b strings.Builder
	b.WriteString(titleStyle.Render("Async Notification Example"))
	b.WriteString("\n\n")
	b.WriteString("  " + keyStyle.Render("s") + "  " + helpStyle.Render("Save file (async)") + "\n")
	b.WriteString("  " + keyStyle.Render("esc") + "  " + helpStyle.Render("Dismiss notification") + "\n")
	b.WriteString("  " + keyStyle.Render("q") + "  " + helpStyle.Render("Quit") + "\n")

	if m.saving {
		b.WriteString("\n  Saving...")
	}

	// Pad content to fill the terminal height
	lines := strings.Split(b.String(), "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	if len(lines) > m.height {
		lines = lines[:m.height]
	}

	fullContent := strings.Join(lines, "\n")

	// Render overlays the notification on top of the content.
	// If no notification is active, it returns the content unchanged.
	v := tea.NewView(m.notify.Render(fullContent))
	v.AltScreen = true
	return v
}
