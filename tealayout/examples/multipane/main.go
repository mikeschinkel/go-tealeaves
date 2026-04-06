package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/mikeschinkel/go-tealeaves/teacolor"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
	"github.com/mikeschinkel/go-tealeaves/teapane"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

const (
	minWeight = 5.0
)

// panelDisplay holds demo-specific display metadata for each panel.
type panelDisplay struct {
	minSize  int
	maxSize  int
	optional bool
}

func panelContentFunc(sp *teapane.StyledPane, pd *panelDisplay) teapane.ContentFunc {
	return func(w, h int, focused bool) string {
		label := sp.Label()
		label += "\n" + strings.Repeat("─", len(label))
		lines := []string{label}
		if pd.minSize > 0 {
			lines = append(lines, fmt.Sprintf("- Min Size: %d", pd.minSize))
		}
		if pd.maxSize > 0 {
			lines = append(lines, fmt.Sprintf("- Max Size: %d", pd.maxSize))
		}
		if pd.optional {
			lines = append(lines, "- Optional: yes")
		}
		pct := sp.FlexPercent()
		if pct > 0 {
			lines = append(lines, fmt.Sprintf("- Flex: %.0f%%", pct))
		}
		lines = append(lines, fmt.Sprintf("- Size: %dx%d", w, h))
		return strings.Join(lines, "\n")
	}
}

func newPane(label string, border teapane.BorderStyle, pd *panelDisplay) *teapane.StyledPane {
	minContentWidth := len("- Size: 999x99")
	sp := teapane.NewStyledPane(border, nil).WithLabel(label).WithMinWidth(minContentWidth)
	sp.SetContentFunc(panelContentFunc(sp, pd))
	return sp
}

// --- Model ---

type model struct {
	mpl *tealayout.MultiPaneLayout

	panes     [5]*teapane.StyledPane
	paneNames [5]string

	headerStr *string
	footerStr *string
	rotator   *tealayout.VisibilityRotator
	width     int
	height    int
}

func initialModel() model {
	paneNames := [5]string{"tree", "main", "sidebar", "det-top", "det-bot"}

	treePD := &panelDisplay{}
	mainPD := &panelDisplay{maxSize: 80}
	sidebarPD := &panelDisplay{optional: true, minSize: 15}
	detTopPD := &panelDisplay{}
	detBotPD := &panelDisplay{}

	treeSP := newPane("Tree", teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.ElectricCyan,
		FocusedColor: teacolor.TrueWhite,
		Foreground:   teacolor.IceCyan,
		PaddingH:     1,
	}, treePD)
	mainSP := newPane("Main", teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.SlateGray,
		FocusedColor: teacolor.TrueWhite,
		Foreground:   teacolor.LightGray,
		PaddingH:     1,
	}, mainPD)
	sidebarSP := newPane("Sidebar", teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.Amber,
		FocusedColor: teacolor.TrueWhite,
		Foreground:   teacolor.Gold,
		PaddingH:     1,
	}, sidebarPD)
	detTopSP := newPane("Details Top", teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.Coral,
		FocusedColor: teacolor.TrueWhite,
		Foreground:   teacolor.Salmon,
		PaddingH:     1,
	}, detTopPD)
	detBotSP := newPane("Details Bot", teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.Tangerine,
		FocusedColor: teacolor.TrueWhite,
		Foreground:   teacolor.Peach,
		PaddingH:     1,
	}, detBotPD)

	tree := tealayout.NewElement(treeSP)
	main := tealayout.NewElement(mainSP)
	sidebar := tealayout.NewElement(sidebarSP)
	detTop := tealayout.NewElement(detTopSP)
	detBot := tealayout.NewElement(detBotSP)

	headerStr := new(string)
	footerStr := new(string)

	header := tealayout.NewElement(teapane.NewPlainPane(func(w, h int, focused bool) string {
		return *headerStr
	}).WithStyle(
		lipgloss.NewStyle().
			Bold(true).
			Background(teacolor.CharcoalGray).
			Foreground(teacolor.ElectricCyan),
	))

	footer := tealayout.NewElement(teapane.NewPlainPane(func(w, h int, focused bool) string {
		return *footerStr
	}).WithStyle(
		lipgloss.NewStyle().
			Background(teacolor.Gunmetal).
			Foreground(teacolor.Gray),
	))

	mpl := tealayout.NewMultiPaneLayout(
		[]tealayout.PaneDef{
			{Name: "tree", Element: tree, Dim: tealayout.Percent(20), MinFlexWeight: minWeight, MinSizeFit: true},
			{Name: "main", Element: main, Dim: tealayout.Percent(30), MinFlexWeight: minWeight, MinSizeFit: true, MaxSize: 80},
			{Name: "sidebar", Element: sidebar, Dim: tealayout.Percent(25), MinFlexWeight: minWeight, Optional: true, MinSize: 15},
			{Name: "details", Dim: tealayout.Percent(25), MinFlexWeight: minWeight, Optional: true, MinSize: 15, Children: []tealayout.PaneDef{
				{Name: "det-top", Element: detTop, Dim: tealayout.Percent(62), MinFlexWeight: minWeight, MinSizeFit: true},
				{Name: "det-bot", Element: detBot, Dim: tealayout.Percent(38), MinFlexWeight: minWeight, MinSizeFit: true},
			}},
		},
		tealayout.WithHeader(header),
		tealayout.WithFooter(footer),
	)

	return model{
		mpl:       mpl,
		panes:     [5]*teapane.StyledPane{treeSP, mainSP, sidebarSP, detTopSP, detBotSP},
		paneNames: paneNames,
		headerStr: headerStr,
		footerStr: footerStr,
		rotator: tealayout.NewVisibilityRotator(mpl, [][]string{
			{"tree", "main", "sidebar", "details"},
			{"tree", "main", "sidebar"},
			{"tree", "main", "details"},
			{"tree", "main"},
			{"main", "sidebar", "details"},
			{"main", "sidebar"},
			{"main", "details"},
			{"tree"},
			{"main"},
		}),
	}
}

func (m model) Init() tea.Cmd { return nil }

//goland:noinspection GoAssignmentToReceiver
func (m model) Update(msg tea.Msg) (_ tea.Model, cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.mpl.SetSize(m.width, m.height)

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			cmd = tea.Quit
			goto end

		case "tab":
			m.mpl.FocusNext()

		case "shift+tab":
			m.mpl.FocusPrev()

		case "right":
			m.mpl.ResizeFocusedColumnByCells(1)

		case "left":
			m.mpl.ResizeFocusedColumnByCells(-1)

		case "up":
			m.mpl.PaneLayout().ResizePaneByCells("det-top", -1)

		case "down":
			m.mpl.PaneLayout().ResizePaneByCells("det-top", 1)

		case "t":
			m.rotator.Next()
		}
	}
end:
	return m, cmd
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		return tea.NewView("Loading...")
	}

	*m.headerStr = m.headerText()
	*m.footerStr = m.footerText()

	// Push flex percentages to pane widgets for display.
	pcts := m.mpl.VisibleFlexPercents()
	for i, name := range m.paneNames {
		m.panes[i].SetFlexPercent(pcts[name])
	}

	m.mpl.MarkDirty()
	output, err := m.mpl.Render()
	if err != nil {
		return tea.NewView(fmt.Sprintf("Layout error: %v", err))
	}
	v := tea.NewView(output)
	v.AltScreen = true
	return v
}

func (m model) headerText() string {
	visible := m.rotator.Current()
	var names []string
	for _, name := range visible {
		marker := " "
		if m.mpl.Focused(name) {
			marker = "*"
		}
		names = append(names, fmt.Sprintf("[%s%s]", marker, name))
	}

	// Also show focused child pane in header if focus is inside a group.
	fp := m.mpl.FocusedPane()
	if fp != nil {
		fpName := fp.Name()
		// Check if focused pane is a child (not in the rotator's visible list).
		found := false
		for _, name := range visible {
			if name == fpName {
				found = true
				break
			}
		}
		if !found {
			names = append(names, fmt.Sprintf("[*%s]", fpName))
		}
	}

	return fmt.Sprintf(" Layout Demo  %s", strings.Join(names, " "))
}

func (m model) footerText() string {
	return " t:toggle  tab/shift+tab:focus  \u2190/\u2192:resize cols  \u2191/\u2193:resize rows  q:quit"
}

// overrideTERM replaces the TERM entry in an environ slice.
func overrideTERM(environ []string, term string) []string {
	result := make([]string, 0, len(environ))
	for _, e := range environ {
		if !strings.HasPrefix(e, "TERM=") {
			result = append(result, e)
		}
	}
	return append(result, "TERM="+term)
}

func main() {
	var opts []tea.ProgramOption

	if teautils.IsJediTerm() {
		opts = append(opts,
			tea.WithEnvironment(overrideTERM(os.Environ(), "ansi")),
			tea.WithColorProfile(colorprofile.TrueColor),
		)
	}

	p := tea.NewProgram(initialModel(), opts...)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
