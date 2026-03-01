package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikeschinkel/go-tealeaves/teanotify"
)

var positions = []teanotify.Position{
	teanotify.TopLeftPosition,
	teanotify.TopCenterPosition,
	teanotify.TopRightPosition,
	teanotify.BottomLeftPosition,
	teanotify.BottomCenterPosition,
	teanotify.BottomRightPosition,
}

type model struct {
	notify   teanotify.NotifyModel
	posIndex int
	width    int
	height   int
}

func main() {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:           40,
		Duration:        3 * time.Second,
		AllowEscToClose: true,
	})
	notify, err := notify.Initialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	m := model{
		notify: notify,
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

		case "i":
			cmd := m.notify.NewNotifyCmd(teanotify.InfoKey, "File saved successfully")
			return m, tea.Batch(notifyCmd, cmd)

		case "w":
			cmd := m.notify.NewNotifyCmd(teanotify.WarnKey, "Disk space low")
			return m, tea.Batch(notifyCmd, cmd)

		case "e":
			cmd := m.notify.NewNotifyCmd(teanotify.ErrorKey, "Connection failed")
			return m, tea.Batch(notifyCmd, cmd)

		case "d":
			cmd := m.notify.NewNotifyCmd(teanotify.DebugKey, "Request took 42ms")
			return m, tea.Batch(notifyCmd, cmd)

		case "p":
			m.posIndex = (m.posIndex + 1) % len(positions)
			m.notify = m.notify.WithPosition(positions[m.posIndex])
			return m, notifyCmd
		}
	}

	return m, notifyCmd
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	keyStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	descStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
)

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var content strings.Builder
	content.WriteString(titleStyle.Render("TeaNotify Example"))
	content.WriteString("\n\n")

	bindings := []struct {
		key  string
		desc string
	}{
		{"i", "Info notice"},
		{"w", "Warn notice"},
		{"e", "Error notice"},
		{"d", "Debug notice"},
		{"p", "Cycle position"},
		{"esc", "Dismiss notice"},
		{"q", "Quit"},
	}

	for _, b := range bindings {
		content.WriteString("  ")
		content.WriteString(keyStyle.Render(b.key))
		content.WriteString("  ")
		content.WriteString(descStyle.Render(b.desc))
		content.WriteString("\n")
	}

	content.WriteString("\n")
	pos := positions[m.posIndex]
	content.WriteString(descStyle.Render(fmt.Sprintf("Position: %s", pos.Label())))
	content.WriteString("\n")

	// Pad content to fill the screen
	lines := strings.Split(content.String(), "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	// Trim to exact screen height
	if len(lines) > m.height {
		lines = lines[:m.height]
	}

	fullContent := strings.Join(lines, "\n")
	return m.notify.Render(fullContent)
}
