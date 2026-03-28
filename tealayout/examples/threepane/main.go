package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
)

const (
	paneCount  = 3
	resizeStep = 0.05
	minWeight  = 0.05

	baseWeightTree = 0.25
	baseWeightCode = 0.43
	baseWeightDiff = 0.32
)

type paneInfo struct {
	label      string
	fg         string
	border     string
	weight     float64
	pctOfTotal float64
}

// pane is a simple widget that fills its area with a label and border.
type pane struct {
	info   *paneInfo
	style  lipgloss.Style
	width  int
	height int
}

func newPane(info *paneInfo) *pane {
	return &pane{
		info: info,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(info.border)).
			Foreground(lipgloss.Color(info.fg)).
			Padding(0, 1),
	}
}

func (p *pane) SetSize(w, h int)      { p.width = w; p.height = h }
func (p *pane) Style() lipgloss.Style { return p.style }

func (p *pane) View() string {
	label := p.info.label
	label += "\n" + strings.Repeat("─", len(label))
	content := fmt.Sprintf("%s\n- Pane Width: %.0f%%\n- Dimensions: %dx%d",
		label, p.info.pctOfTotal, p.width, p.height)
	return p.style.Width(p.width).Height(p.height).Render(content)
}

func (p *pane) updateStyle(focused bool) {
	borderColor := p.info.border
	if focused {
		borderColor = "#ffffff"
	}
	p.style = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Foreground(lipgloss.Color(p.info.fg)).
		Padding(0, 1)
}

// headerBar is a simple 1-line bar.
type headerBar struct {
	text  string
	width int
	style lipgloss.Style
}

func (h *headerBar) SetSize(w, _ int) { h.width = w }
func (h *headerBar) View() string {
	return h.style.Width(h.width).Render(h.text)
}

// --- Model ---

type model struct {
	layout *tealayout.Layout
	focus  *tealayout.FocusManager

	// Typed element handles for widget access.
	tree   *tealayout.Element[*pane]
	code   *tealayout.Element[*pane]
	diff   *tealayout.Element[*pane]
	header *tealayout.Element[*headerBar]
	footer *tealayout.Element[*headerBar]

	paneInfos [paneCount]*paneInfo
	visCombo  int
	width     int
	height    int
}

// visibility combos: names of visible panes
var visibilityCombos = [][]string{
	{"tree", "code", "diff"},
	{"tree", "code"},
	{"tree", "diff"},
	{"code", "diff"},
	{"tree"},
	{"code"},
	{"diff"},
}

func initialModel() model {
	infos := [paneCount]*paneInfo{
		{label: "Tree", fg: "#a5f3fc", border: "#67e8f9", weight: baseWeightTree},
		{label: "Code", fg: "#d1d5db", border: "#6b7280", weight: baseWeightCode},
		{label: "Diff", fg: "#fca5a5", border: "#f87171", weight: baseWeightDiff},
	}

	tree := tealayout.NewElement(newPane(infos[0]))
	code := tealayout.NewElement(newPane(infos[1]))
	diff := tealayout.NewElement(newPane(infos[2]))

	header := tealayout.NewElement(&headerBar{
		style: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#67e8f9")),
	})

	footer := tealayout.NewElement(&headerBar{
		style: lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2937")).
			Foreground(lipgloss.Color("#9ca3af")),
	})

	root := tealayout.NewColumn(tealayout.Percent100,
		tealayout.NewRow(tealayout.Fixed(1), header).WithName("header").WithAlignment(tealayout.MiddleCenter),
		tealayout.NewRow(tealayout.Flex(1),
			tealayout.NewColumn(tealayout.Flex(baseWeightTree), tree).WithName("tree"),
			tealayout.NewColumn(tealayout.Flex(baseWeightCode), code).WithName("code"),
			tealayout.NewColumn(tealayout.Flex(baseWeightDiff), diff).WithName("diff"),
		).WithName("content"),
		tealayout.NewRow(tealayout.Fixed(1), footer).WithName("footer"),
	)

	layout := tealayout.NewLayout(root)
	focus := tealayout.NewFocusManager(layout)

	return model{
		layout:    layout,
		focus:     focus,
		tree:      tree,
		code:      code,
		diff:      diff,
		header:    header,
		footer:    footer,
		paneInfos: infos,
		visCombo:  0,
	}
}

func (m model) Init() tea.Cmd { return nil }

//goland:noinspection GoAssignmentToReceiver
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.layout.SetSize(m.width, m.height)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.focus.FocusNext()

		case "shift+tab":
			m.focus.FocusPrev()

		case "+", "=":
			m.resizeFocused(resizeStep)

		case "-", "_":
			m.resizeFocused(-resizeStep)

		case "t":
			m.toggleVisibility()
		}
	}
	return m, nil
}

func (m *model) resizeFocused(delta float64) {
	fp := m.focus.FocusedPane()
	if fp == nil {
		return
	}
	// Find the matching paneInfo and update weight
	for _, info := range m.paneInfos {
		if info.label == fp.Name() || strings.EqualFold(info.label, fp.Name()) {
			newWeight := info.weight + delta
			if newWeight < minWeight {
				newWeight = minWeight
			}
			info.weight = newWeight
			fp.SetDimension(tealayout.Flex(newWeight))
			return
		}
	}
}

func (m *model) toggleVisibility() {
	m.visCombo = (m.visCombo + 1) % len(visibilityCombos)
	visible := visibilityCombos[m.visCombo]

	// Build set of visible names
	visSet := make(map[string]bool, len(visible))
	for _, name := range visible {
		visSet[name] = true
	}

	// Show/hide panes by name
	for _, name := range []string{"tree", "code", "diff"} {
		p := m.layout.Pane(name)
		if p != nil {
			p.SetVisible(visSet[name])
		}
	}

	// Ensure focused pane is visible
	m.focus.EnsureFocusedVisible()
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	// Update header/footer text
	m.header.Widget().text = m.headerText()
	m.footer.Widget().text = m.footerText()

	// Update styles based on focus
	m.tree.Widget().updateStyle(m.focus.Focused("tree"))
	m.code.Widget().updateStyle(m.focus.Focused("code"))
	m.diff.Widget().updateStyle(m.focus.Focused("diff"))

	// Compute effective percentages for display
	visible := visibilityCombos[m.visCombo]
	totalWeight := 0.0
	for _, name := range visible {
		for _, info := range m.paneInfos {
			if strings.EqualFold(info.label, name) {
				totalWeight += info.weight
			}
		}
	}
	for _, info := range m.paneInfos {
		info.pctOfTotal = info.weight / totalWeight * 100
	}

	m.layout.MarkDirty()
	output, err := m.layout.Render()
	if err != nil {
		return tea.NewView(fmt.Sprintf("Layout error: %v", err))
	}
	v := tea.NewView(output)
	v.AltScreen = true
	return v
}

func (m model) headerText() string {
	visible := visibilityCombos[m.visCombo]
	var names []string
	for _, name := range visible {
		marker := " "
		if m.focus.Focused(name) {
			marker = "*"
		}
		names = append(names, fmt.Sprintf("[%s%s]", marker, name))
	}
	return fmt.Sprintf(" Three-Pane Demo  %s", strings.Join(names, " "))
}

func (m model) footerText() string {
	return " t:toggle visibility tab:focus  +/-:resize (requires focus)  q:quit"
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		cliutil.Stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}
