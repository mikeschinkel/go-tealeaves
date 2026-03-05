package teastatus

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
)

// StatusBarModel is a Bubble Tea model for a two-zone status bar.
// Left zone: key-action menu items (e.g., "[?] Menu  [tab] Switch pane")
// Right zone: text indicators (e.g., "DEPS IN-FLUX | 3 batches")
type StatusBarModel struct {
	Styles     Styles
	menuItems  []MenuItem
	indicators []StatusIndicator
	width      int
}

// NewStatusBarModel creates a new status bar StatusBarModel with default styles.
func NewStatusBarModel() StatusBarModel {
	return StatusBarModel{
		Styles: DefaultStyles(),
	}
}

// WithStyles returns a copy with the given styles override.
func (m StatusBarModel) WithStyles(styles Styles) StatusBarModel {
	m.Styles = styles
	return m
}

// SetSize sets the terminal width for the status bar.
func (m StatusBarModel) SetSize(width int) StatusBarModel {
	m.width = width
	return m
}

// SetMenuItems replaces the current menu items.
func (m StatusBarModel) SetMenuItems(items []MenuItem) StatusBarModel {
	m.menuItems = items
	return m
}

// SetIndicators replaces the current indicators.
func (m StatusBarModel) SetIndicators(indicators []StatusIndicator) StatusBarModel {
	m.indicators = indicators
	return m
}

// Init implements tea.Model. No-op for status bar.
func (m StatusBarModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model. Handles SetMenuItemsMsg and SetIndicatorsMsg.
func (m StatusBarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SetMenuItemsMsg:
		m.menuItems = msg.Items
	case SetIndicatorsMsg:
		m.indicators = msg.Indicators
	}
	return m, nil
}

// View implements tea.Model. Renders the two-zone status bar.
func (m StatusBarModel) View() tea.View {
	var view string
	var left string
	var right string
	var leftWidth int
	var rightWidth int
	var gap int
	var sb strings.Builder

	left = m.renderMenuItems()
	right = m.renderIndicators()

	leftWidth = ansi.StringWidth(left)
	rightWidth = ansi.StringWidth(right)

	// If width not set or too narrow for right side, just show left
	if m.width <= 0 || leftWidth+rightWidth+2 > m.width {
		view = m.Styles.BarStyle.Render(left)
		goto end
	}

	// Fill gap with spaces
	gap = m.width - leftWidth - rightWidth
	if gap < 0 {
		gap = 0
	}

	sb.WriteString(left)
	sb.WriteString(strings.Repeat(" ", gap))
	sb.WriteString(right)

	view = m.Styles.BarStyle.Render(sb.String())

end:
	return tea.NewView(view)
}
