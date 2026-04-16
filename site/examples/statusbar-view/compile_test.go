// Source: site/src/content/docs/components/statusbar-view.mdx:22#b811a057,82#3a187afd,108#15050ba0,136#e8c1b86d,154#c071d3b3
package examples_test

import (
	"testing"

	key "charm.land/bubbles/v2/key"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
)

// TestCompile_StatusBarQuickExample verifies the quick example from statusbar-view.mdx.
func TestCompile_StatusBarQuickExample(t *testing.T) {
	sb := teastatus.NewStatusBarModel()

	sb = sb.SetMenuItems([]teastatus.MenuItem{
		{KeyText: "?", Label: "Help"},
		{KeyText: "tab", Label: "Switch pane"},
		{KeyText: "q", Label: "Quit"},
	})

	sb = sb.SetIndicators([]teastatus.StatusIndicator{
		teastatus.NewStatusIndicator("DEPS IN-FLUX"),
		teastatus.NewStatusIndicator("3 batches"),
	})

	sb = sb.SetSize(80)
	_ = sb.View()
}

// TestCompile_StatusBarStyles verifies WithStyles from statusbar-view.mdx.
func TestCompile_StatusBarStyles(t *testing.T) {
	sb := teastatus.NewStatusBarModel()
	sb = sb.
		WithStyles(teastatus.Styles{
			MenuKeyStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true),
			MenuLabelStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
			MenuSeparator:     "  ",
			IndicatorStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
			IndicatorSepStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
			SeparatorKind:     teastatus.PipeSeparator,
			BarStyle:          lipgloss.NewStyle().Background(lipgloss.Color("236")),
		}).
		SetSize(80)
	_ = sb
}

// TestCompile_IndicatorStyle verifies per-indicator styles from statusbar-view.mdx.
func TestCompile_IndicatorStyle(t *testing.T) {
	indicators := []teastatus.StatusIndicator{
		teastatus.NewStatusIndicator("DEPS IN-FLUX").
			WithStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("196"))),
		teastatus.NewStatusIndicator("3 batches"),
	}
	sb := teastatus.NewStatusBarModel().SetIndicators(indicators).SetSize(80)
	_ = sb
}

// TestCompile_NewMenuItemFromBinding verifies NewMenuItem from statusbar-view.mdx.
func TestCompile_NewMenuItemFromBinding(t *testing.T) {
	helpBinding := key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help"))
	quitBinding := key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "quit"))

	sb := teastatus.NewStatusBarModel().SetMenuItems([]teastatus.MenuItem{
		teastatus.NewMenuItem(helpBinding, &teastatus.MenuItemOpts{Label: "Help"}),
		teastatus.NewMenuItem(quitBinding, &teastatus.MenuItemOpts{Label: "Quit"}),
	})
	_ = sb
}

// TestCompile_DynamicUpdates verifies SetMenuItemsMsg from statusbar-view.mdx.
func TestCompile_DynamicUpdates(t *testing.T) {
	msg := teastatus.SetMenuItemsMsg{
		Items: []teastatus.MenuItem{
			{KeyText: "r", Label: "Refresh"},
			{KeyText: "q", Label: "Quit"},
		},
	}
	_ = msg
}
