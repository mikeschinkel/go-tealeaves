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
	hStep       = 4     // horizontal resize step (cells)
	vStep       = 0.1   // vertical resize step (weight delta)
	minVWeight  = vStep // keep positive; layout engine's WithMinSize() handles the visual floor
	goldenSmall = 1.0
	goldenLarge = 1.618

	minWidthMain    = 25
	maxWidthMain    = 60 // caps Main so siblings get more space before collapsing
	minWidthSidebar = 25
	minWidthDetails = 20
)

type panel struct {
	label    string
	minWidth int // 0 = no minimum / not applicable
	maxWidth int // 0 = no maximum / not applicable
	optional bool
	style    lipgloss.Style
	width    int
	height   int
}

func newPanel(label, border string, minWidth, maxWidth int, optional bool) *panel {
	return &panel{
		label:    label,
		minWidth: minWidth,
		maxWidth: maxWidth,
		optional: optional,
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(border)).
			Padding(0, 1),
	}
}

func (p *panel) SetSize(w, h int) { p.width = w; p.height = h }

// minHeight returns the minimum height needed to display this panel's content
// plus border overhead. Update this if View() content changes.
func (p *panel) minHeight() int {
	lines := 2 // label + separator (always present)
	if p.minWidth > 0 {
		lines++
	}
	if p.maxWidth > 0 {
		lines++
	}
	if p.optional {
		lines++
	}
	lines++    // "- Dimensions: WxH" line (always present)
	lines += 2 // top + bottom border (RoundedBorder)
	return lines
}

func (p *panel) View() string {
	// width/height are total (border-box) dimensions — don't also implement
	// Styler, which would cause setChildSize to subtract border before SetSize,
	// then border-box Width/Height would subtract it again.
	label := p.label
	label += "\n" + strings.Repeat("─", len(label))
	lines := []string{label}
	if p.minWidth > 0 {
		lines = append(lines, fmt.Sprintf("- Min Width: %d", p.minWidth))
	}
	if p.maxWidth > 0 {
		lines = append(lines, fmt.Sprintf("- Max Width: %d", p.maxWidth))
	}
	if p.optional {
		lines = append(lines, "- Optional: yes")
	}
	lines = append(lines, fmt.Sprintf("- Dimensions: %dx%d", p.width, p.height))
	content := strings.Join(lines, "\n")
	return p.style.Width(p.width).Height(p.height).Render(content)
}

type footer struct {
	width int
	style lipgloss.Style
}

func (f *footer) SetSize(w, h int) { f.width = w }
func (f *footer) View() string {
	text := " ←/→: shrink/grow total width  ↑/↓: shift Details split  q: quit"
	return f.style.Width(f.width).Render(text)
}

type model struct {
	main         *panel
	sidebar      *panel
	detailsTop   *panel
	detailsBot   *panel
	footer       *footer
	termWidth    int
	termHeight   int
	simWidth     int // simulated total width for the 3 columns
	topWeight    float64
	bottomWeight float64
	layout       *tealayout.Layout
	layoutDirty  bool
}

func initialModel() model {
	return model{
		main:       newPanel("Main", "#67e8f9", minWidthMain, maxWidthMain, false),
		sidebar:    newPanel("Sidebar", "#fbbf24", minWidthSidebar, 0, true),
		detailsTop: newPanel("Details", "#f87171", 0, 0, false),
		detailsBot: newPanel("More Details", "#fb923c", 0, 0, false),
		footer: &footer{
			style: lipgloss.NewStyle().
				Background(lipgloss.Color("#1f2937")).
				Foreground(lipgloss.Color("#9ca3af")),
		},
		topWeight:    goldenLarge,
		bottomWeight: goldenSmall,
	}
}

func (m model) buildLayout() *tealayout.Layout {
	w := m.simWidth
	if w <= 0 {
		w = m.termWidth
	}

	// Details column: top/bottom split by golden ratio weights.
	// Up arrow grows Details (top), down arrow grows More Details (bottom).
	detailsCol := tealayout.NewColumn(tealayout.Percent100,
		tealayout.NewRow(tealayout.Flex(m.topWeight), m.detailsTop).WithMinSize(m.detailsTop.minHeight()),
		tealayout.NewRow(tealayout.Flex(m.bottomWeight), m.detailsBot).WithMinSize(m.detailsBot.minHeight()),
	)

	contentRow := tealayout.NewRow(tealayout.Fixed(w),
		tealayout.NewColumn(tealayout.Percent50, m.main).WithMinSize(minWidthMain).WithMaxSize(maxWidthMain),
		tealayout.NewColumn(tealayout.Percent25, m.sidebar).WithMinSize(minWidthSidebar).WithOptional(true),
		tealayout.NewColumn(tealayout.Percent25, detailsCol).WithMinSize(minWidthDetails).WithOptional(true),
	)

	// Wrap the content row in a fixed-width constraint inside an outer row
	// so the simulated width is respected by the layout engine.
	widthWrapper := tealayout.NewRow(tealayout.Percent100, contentRow)

	root := tealayout.NewColumn(tealayout.Percent100,
		widthWrapper,
		tealayout.NewRow(tealayout.Fixed(1), m.footer),
	)

	layout := tealayout.NewLayout(root)
	layout.SetSize(m.termWidth, m.termHeight)
	return layout
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height
		if m.simWidth == 0 || m.simWidth > msg.Width {
			m.simWidth = msg.Width
		}
		m.layoutDirty = true

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "left":
			if w := m.simWidth - hStep; w >= 20 {
				m.simWidth = w
				m.layoutDirty = true
			}
		case "right":
			if w := m.simWidth + hStep; w <= m.termWidth {
				m.simWidth = w
				m.layoutDirty = true
			}
		case "up":
			if m.topWeight-vStep >= minVWeight {
				m.topWeight -= vStep
				m.bottomWeight += vStep
				m.layoutDirty = true
			}
		case "down":
			if m.bottomWeight-vStep >= minVWeight {
				m.bottomWeight -= vStep
				m.topWeight += vStep
				m.layoutDirty = true
			}
		}
	}
	if m.layoutDirty {
		m.layout = m.buildLayout()
		m.layoutDirty = false
	}
	return m, nil
}

func (m model) View() tea.View {
	if m.layout == nil {
		return tea.NewView("Loading...")
	}
	output, err := m.layout.Render()
	if err != nil {
		return tea.NewView(fmt.Sprintf("Error: %v", err))
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
