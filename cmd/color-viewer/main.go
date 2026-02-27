package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	viewport viewport.Model
	header   string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4 // Header (2 lines) + help text (2 lines)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.header + m.viewport.View() + "\n  ↑/↓ scroll, q quit | Format: bg/fg (e.g., 53/015 = background 53, foreground 15)"
}

// Common foreground colors to test against each background
var fgColors = []int{0, 15, 231, 232, 244, 250, 196, 46, 21, 226}
var fgLabels = []string{"blk", "wht", "brw", "drk", "gry", "ltg", "red", "grn", "blu", "yel"}

const sampleWidth = 11 // Visual width of each sample cell

func buildHeader() string {
	var sb strings.Builder

	// Header row
	sb.WriteString(fmt.Sprintf("%-8s", "BG #"))
	for _, label := range fgLabels {
		// Center the label in sampleWidth
		header := fmt.Sprintf("fg:%s", label)
		pad := sampleWidth - len(header)
		leftPad := pad / 2
		rightPad := pad - leftPad
		sb.WriteString(strings.Repeat(" ", leftPad) + header + strings.Repeat(" ", rightPad))
	}
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("─", 8+sampleWidth*len(fgColors)) + "\n")

	return sb.String()
}

func buildColorContent() string {
	var sb strings.Builder

	// Each row is one background color with multiple foreground samples
	for bg := 0; bg < 256; bg++ {
		bgColor := lipgloss.Color(fmt.Sprintf("%d", bg))

		// Background color number (8 chars wide)
		sb.WriteString(fmt.Sprintf("%-8d", bg))

		// Show this background with each foreground color
		for _, fg := range fgColors {
			fgColor := lipgloss.Color(fmt.Sprintf("%d", fg))

			// Create sample text with fixed visual width (zero-pad fg to 3 digits)
			text := fmt.Sprintf("%3d/%03d", bg, fg)
			sample := lipgloss.NewStyle().
				Foreground(fgColor).
				Background(bgColor).
				Width(sampleWidth).
				Align(lipgloss.Center).
				Render(text)
			sb.WriteString(sample)
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func main() {
	header := buildHeader()
	content := buildColorContent()

	vp := viewport.New(8+sampleWidth*len(fgColors)+5, 20)
	vp.SetContent(content)

	m := model{
		viewport: vp,
		header:   header,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
