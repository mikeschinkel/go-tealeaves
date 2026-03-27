package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
)

type panel struct {
	label  string
	style  lipgloss.Style
	width  int
	height int
}

func newPanel(label, border string) *panel {
	return &panel{
		label: label,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(border)).
			Padding(0, 1),
	}
}

func (p *panel) SetSize(w, h int)      { p.width = w; p.height = h }
func (p *panel) Style() lipgloss.Style { return p.style }
func (p *panel) View() string {
	content := fmt.Sprintf("%s\n%dx%d", p.label, p.width, p.height)
	return p.style.Width(p.width).Height(p.height).Render(content)
}

type model struct {
	layout *tealayout.Layout
}

func initialModel() model {
	left := newPanel("Narrow (1.0)", "#67e8f9")
	right := newPanel("Wide (1.618)", "#fbbf24")

	root := tealayout.NewRow(tealayout.Percent100,
		tealayout.NewColumn(tealayout.Flex(1.0), left),
		tealayout.NewColumn(tealayout.Flex(1.618), right),
	)

	return model{layout: tealayout.NewLayout(root)}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.layout.SetSize(msg.Width, msg.Height)
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	s, err := m.layout.Render()
	if err != nil {
		return tea.NewView(fmt.Sprintf("Error: %v", err))
	}
	v := tea.NewView(s)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		cliutil.Stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}
