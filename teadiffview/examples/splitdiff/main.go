package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-diffutils"
	"github.com/mikeschinkel/go-tealeaves/teadiffview"
)

// Sample "before" and "after" texts for demonstration.
var oldText = `package greeting

import "fmt"

func Hello(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func Goodbye(name string) string {
	return fmt.Sprintf("Goodbye, %s.", name)
}
`

var newText = `package greeting

import (
	"fmt"
	"strings"
)

// Hello greets the given name.
func Hello(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "World"
	}
	return fmt.Sprintf("Hello, %s!", name)
}

func Farewell(name string) string {
	return fmt.Sprintf("Farewell, %s!", name)
}
`

type model struct {
	split  teadiffview.SplitDiffModel
	width  int
	height int
}

func main() {
	oldLines := strings.Split(strings.TrimRight(oldText, "\n"), "\n")
	newLines := strings.Split(strings.TrimRight(newText, "\n"), "\n")

	content, err := diffutils.DiffLines(oldLines, newLines, "greeting.go")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error computing diff: %v\n", err)
		os.Exit(1)
	}

	split := teadiffview.NewSplitDiffModel(&teadiffview.SplitDiffModelArgs{
		Width:  80,
		Height: 24,
	})
	split = split.SetContent(content)

	p := tea.NewProgram(model{split: split})
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.split = m.split.SetSize(msg.Width, msg.Height-2)
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "g":
			m.split = m.split.GoToTop()
			return m, nil
		case "G":
			m.split = m.split.GoToBottom()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.split, cmd = m.split.Update(msg)
	return m, cmd
}

var headerStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("205")).
	PaddingBottom(1)

func (m model) View() tea.View {
	if m.width == 0 {
		return tea.NewView("Loading...")
	}

	var b strings.Builder
	b.WriteString(headerStyle.Render("SplitDiffModel Example"))
	b.WriteString("\n")
	b.WriteString(m.split.View().Content)

	v := tea.NewView(b.String())
	v.AltScreen = true
	return v
}
