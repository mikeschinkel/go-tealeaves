package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func runColors(args []string) error {
	fs := flag.NewFlagSet("colors", flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `tlcli colors — Interactive 256-color palette viewer

Displays all 256 ANSI colors in a grid with multiple foreground
colors against each background color.

Usage:
  tlcli colors

Keys:
  ↑/↓   Scroll
  q     Quit
`)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	header := buildHeader()
	content := buildColorContent()

	vp := viewport.New(viewport.WithWidth(8+sampleWidth*len(fgColors)+5), viewport.WithHeight(20))
	vp.SetContent(content)

	m := colorModel{
		viewport: vp,
		header:   header,
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}

type colorModel struct {
	viewport viewport.Model
	header   string
}

func (m colorModel) Init() tea.Cmd {
	return nil
}

func (m colorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.viewport.SetWidth(msg.Width)
		m.viewport.SetHeight(msg.Height - 4)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m colorModel) View() tea.View {
	v := tea.NewView(m.header + m.viewport.View() + "\n  ↑/↓ scroll, q quit | Format: bg/fg (e.g., 53/015 = background 53, foreground 15)")
	v.AltScreen = true
	return v
}

// Common foreground colors to test against each background
var fgColors = []int{0, 15, 231, 232, 244, 250, 196, 46, 21, 226}
var fgLabels = []string{"blk", "wht", "brw", "drk", "gry", "ltg", "red", "grn", "blu", "yel"}

const sampleWidth = 11

func buildHeader() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("%-8s", "BG #"))
	for _, label := range fgLabels {
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

	for bg := 0; bg < 256; bg++ {
		bgColor := lipgloss.Color(fmt.Sprintf("%d", bg))

		sb.WriteString(fmt.Sprintf("%-8d", bg))

		for _, fg := range fgColors {
			fgColor := lipgloss.Color(fmt.Sprintf("%d", fg))

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
