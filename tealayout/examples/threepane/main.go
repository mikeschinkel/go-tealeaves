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
	paneTree   = 0
	paneCode   = 1
	paneDiff   = 2
	resizeStep = 0.05
	minWeight  = 0.05

	// Base flex weights — proportions that define relative pane widths.
	// When panes are hidden, visible panes grow proportionally because
	// Flex distributes by weight ratio: w_i / sum(w_visible).
	baseWeightTree = 0.25 // 25% of total
	baseWeightCode = 0.43 // 43% of total
	baseWeightDiff = 0.32 // 32% of total
)

// visibility combos: indices of visible panes
var visibilityCombos = [][]int{
	{paneTree, paneCode, paneDiff}, // L+M+R
	{paneTree, paneCode},           // L+M
	{paneTree, paneDiff},           // L+R
	{paneCode, paneDiff},           // M+R
	{paneTree},                     // L
	{paneCode},                     // M
	{paneDiff},                     // R
}

type paneInfo struct {
	label      string
	fg         string
	border     string
	weight     float64
	focused    bool
	pctOfTotal float64 // effective percentage of visible space, set before render
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
	if p.info.focused {
		label = "> " + label
	}
	label += "\n" + strings.Repeat("─", len(label))
	content := fmt.Sprintf("%s\n- Pane Width: %.0f%%\n- Dimensions: %dx%d",
		label, p.info.pctOfTotal, p.width, p.height)
	return p.style.Width(p.width).Height(p.height).Render(content)
}

func (p *pane) updateStyle() {
	borderColor := p.info.border
	if p.info.focused {
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

type model struct {
	panes       [paneCount]*paneInfo
	paneWidgets [paneCount]*pane
	header      *headerBar
	footer      *headerBar
	focused     int // index into panes
	visCombo    int // index into visibilityCombos
	width       int
	height      int
}

func initialModel() model {
	panes := [paneCount]*paneInfo{
		{label: "Tree", fg: "#a5f3fc", border: "#67e8f9", weight: baseWeightTree},
		{label: "Code", fg: "#d1d5db", border: "#6b7280", weight: baseWeightCode},
		{label: "Diff", fg: "#fca5a5", border: "#f87171", weight: baseWeightDiff},
	}
	panes[0].focused = true

	widgets := [paneCount]*pane{
		newPane(panes[0]),
		newPane(panes[1]),
		newPane(panes[2]),
	}

	header := &headerBar{
		style: lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#67e8f9")),
	}

	footer := &headerBar{
		style: lipgloss.NewStyle().
			Background(lipgloss.Color("#1f2937")).
			Foreground(lipgloss.Color("#9ca3af")),
	}

	return model{
		panes:       panes,
		paneWidgets: widgets,
		header:      header,
		footer:      footer,
		focused:     paneTree,
		visCombo:    0,
	}
}

func (m model) visiblePanes() []int {
	return visibilityCombos[m.visCombo]
}

func (m model) isFocusedVisible() bool {
	for _, idx := range m.visiblePanes() {
		if idx == m.focused {
			return true
		}
	}
	return false
}

// buildLayout constructs a fresh layout from the current state.
func (m model) buildLayout() *tealayout.Layout {
	visible := m.visiblePanes()

	// Compute effective percentage for each visible pane
	totalWeight := 0.0
	for _, idx := range visible {
		totalWeight += m.panes[idx].weight
	}
	for _, idx := range visible {
		m.panes[idx].pctOfTotal = m.panes[idx].weight / totalWeight * 100
	}

	elements := make([]tealayout.Element, 0, len(visible))
	for _, idx := range visible {
		p := m.paneWidgets[idx]
		p.updateStyle()
		elements = append(elements, tealayout.NewColumn(tealayout.Flex(m.panes[idx].weight), p))
	}

	m.header.text = m.headerText()
	m.footer.text = m.footerText()

	contentRow := tealayout.NewRow(tealayout.Flex(1), elements...)

	root := tealayout.NewColumn(tealayout.Percent100,
		tealayout.NewRow(tealayout.Fixed(1), m.header),
		contentRow,
		tealayout.NewRow(tealayout.Fixed(1), m.footer),
	)

	layout := tealayout.NewLayout(root)
	layout.SetSize(m.width, m.height)
	return layout
}

func (m model) headerText() string {
	visible := m.visiblePanes()
	var names []string
	for _, idx := range visible {
		marker := " "
		if idx == m.focused {
			marker = "*"
		}
		names = append(names, fmt.Sprintf("[%s%s]", marker, m.panes[idx].label))
	}
	return fmt.Sprintf(" Three-Pane Demo  %s", strings.Join(names, " "))
}

func (m model) footerText() string {
	return " t:toggle visibility tab:focus  +/-:resize (requires focus)  q:quit"
}

func (m model) Init() tea.Cmd { return nil }

//goland:noinspection GoAssignmentToReceiver
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m = m.focusNext()

		case "shift+tab":
			m = m.focusPrev()

		case "+", "=":
			m = m.resizeFocused(resizeStep)

		case "-", "_":
			m = m.resizeFocused(-resizeStep)

		case "t":
			m = m.toggleVisibility()
		}
	}
	return m, nil
}

func (m model) focusNext() model {
	visible := m.visiblePanes()
	// Find current position in visible list
	cur := 0
	for i, idx := range visible {
		if idx == m.focused {
			cur = i
			break
		}
	}
	m.panes[m.focused].focused = false
	next := visible[(cur+1)%len(visible)]
	m.focused = next
	m.panes[m.focused].focused = true
	return m
}

func (m model) focusPrev() model {
	visible := m.visiblePanes()
	cur := 0
	for i, idx := range visible {
		if idx == m.focused {
			cur = i
			break
		}
	}
	m.panes[m.focused].focused = false
	prev := visible[(cur-1+len(visible))%len(visible)]
	m.focused = prev
	m.panes[m.focused].focused = true
	return m
}

func (m model) resizeFocused(delta float64) model {
	newWeight := m.panes[m.focused].weight + delta
	if newWeight < minWeight {
		newWeight = minWeight
	}
	m.panes[m.focused].weight = newWeight
	return m
}

func (m model) toggleVisibility() model {
	m.visCombo = (m.visCombo + 1) % len(visibilityCombos)
	// If focused pane is no longer visible, move focus to first visible
	if !m.isFocusedVisible() {
		m.panes[m.focused].focused = false
		m.focused = m.visiblePanes()[0]
		m.panes[m.focused].focused = true
	}
	return m
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	layout := m.buildLayout()
	output, err := layout.Render()
	if err != nil {
		return tea.NewView(fmt.Sprintf("Layout error: %v", err))
	}
	v := tea.NewView(output)
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
