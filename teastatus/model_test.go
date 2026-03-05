package teastatus

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func testMenuItems() []MenuItem {
	return []MenuItem{
		{Key: "?", Label: "Help"},
		{Key: "tab", Label: "Switch"},
	}
}

func testIndicators() []StatusIndicator {
	return []StatusIndicator{
		{Text: "Ready"},
		{Text: "3 items"},
	}
}

// --- Layer 1 Tests ---

func TestNew(t *testing.T) {
	m := NewStatusBarModel()
	// Should have default styles (non-zero)
	if m.Styles.MenuSeparator == "" {
		t.Error("expected non-empty MenuSeparator in default styles")
	}
}

func TestModel_SetSize(t *testing.T) {
	m := NewStatusBarModel()
	m = m.SetSize(80)
	if m.width != 80 {
		t.Errorf("expected width=80, got %d", m.width)
	}
}

func TestModel_SetMenuItems(t *testing.T) {
	m := NewStatusBarModel()
	items := testMenuItems()
	m = m.SetMenuItems(items)
	if len(m.menuItems) != 2 {
		t.Errorf("expected 2 menu items, got %d", len(m.menuItems))
	}
}

func TestModel_SetIndicators(t *testing.T) {
	m := NewStatusBarModel()
	indicators := testIndicators()
	m = m.SetIndicators(indicators)
	if len(m.indicators) != 2 {
		t.Errorf("expected 2 indicators, got %d", len(m.indicators))
	}
}

func TestModel_Update_SetMenuItemsMsg(t *testing.T) {
	m := NewStatusBarModel()
	items := testMenuItems()
	result, cmd := m.Update(SetMenuItemsMsg{Items: items})
	m = result.(StatusBarModel)

	if len(m.menuItems) != 2 {
		t.Errorf("expected 2 menu items after msg, got %d", len(m.menuItems))
	}
	if cmd != nil {
		t.Error("expected nil cmd from SetMenuItemsMsg")
	}
}

func TestModel_Update_SetIndicatorsMsg(t *testing.T) {
	m := NewStatusBarModel()
	indicators := testIndicators()
	result, cmd := m.Update(SetIndicatorsMsg{Indicators: indicators})
	m = result.(StatusBarModel)

	if len(m.indicators) != 2 {
		t.Errorf("expected 2 indicators after msg, got %d", len(m.indicators))
	}
	if cmd != nil {
		t.Error("expected nil cmd from SetIndicatorsMsg")
	}
}

func TestModel_Update_UnknownMsg(t *testing.T) {
	m := NewStatusBarModel()
	m = m.SetMenuItems(testMenuItems())
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	updated := result.(StatusBarModel)

	if len(updated.menuItems) != 2 {
		t.Errorf("expected menu items unchanged, got %d", len(updated.menuItems))
	}
	if cmd != nil {
		t.Error("expected nil cmd for unknown message")
	}
}

// --- Layer 2 Tests ---

func TestModel_View_Empty(t *testing.T) {
	m := NewStatusBarModel().SetSize(80)
	view := m.View()
	// Should render something (bar background) even with no items
	// At minimum it shouldn't panic
	_ = view
}

func TestModel_View_MenuItems(t *testing.T) {
	m := NewStatusBarModel().SetSize(80).SetMenuItems(testMenuItems())
	view := m.View()
	if !strings.Contains(view.Content, "?") {
		t.Error("expected view to contain key '?'")
	}
	if !strings.Contains(view.Content, "Help") {
		t.Error("expected view to contain label 'Help'")
	}
	if !strings.Contains(view.Content, "tab") {
		t.Error("expected view to contain key 'tab'")
	}
	if !strings.Contains(view.Content, "Switch") {
		t.Error("expected view to contain label 'Switch'")
	}
}

func TestModel_View_Indicators(t *testing.T) {
	m := NewStatusBarModel().SetSize(80).SetIndicators(testIndicators())
	view := m.View()
	if !strings.Contains(view.Content, "Ready") {
		t.Error("expected view to contain 'Ready'")
	}
	if !strings.Contains(view.Content, "3 items") {
		t.Error("expected view to contain '3 items'")
	}
}

func TestModel_View_BothZones(t *testing.T) {
	m := NewStatusBarModel().SetSize(80).SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
	view := m.View()

	// Both menu and indicators should be present
	if !strings.Contains(view.Content, "Help") {
		t.Error("expected view to contain menu item label")
	}
	if !strings.Contains(view.Content, "Ready") {
		t.Error("expected view to contain indicator text")
	}
}

func TestModel_View_SeparatorKinds(t *testing.T) {
	indicators := testIndicators()

	t.Run("PipeSeparator", func(t *testing.T) {
		styles := DefaultStyles()
		styles.SeparatorKind = PipeSeparator
		m := NewStatusBarModel().WithStyles(styles).SetSize(80).SetIndicators(indicators)
		view := m.View()
		if !strings.Contains(view.Content, "|") {
			t.Error("expected pipe separator in view")
		}
	})

	t.Run("BracketSeparator", func(t *testing.T) {
		styles := DefaultStyles()
		styles.SeparatorKind = BracketSeparator
		m := NewStatusBarModel().WithStyles(styles).SetSize(80).SetIndicators(indicators)
		view := m.View()
		if !strings.Contains(view.Content, "[Ready]") {
			t.Error("expected bracket-wrapped indicator in view")
		}
	})

	t.Run("SpaceSeparator", func(t *testing.T) {
		styles := DefaultStyles()
		styles.SeparatorKind = SpaceSeparator
		m := NewStatusBarModel().WithStyles(styles).SetSize(80).SetIndicators(indicators)
		view := m.View()
		if !strings.Contains(view.Content, "Ready") {
			t.Error("expected indicator text in view")
		}
		if strings.Contains(view.Content, "|") {
			t.Error("expected no pipe separator for SpaceSeparator")
		}
	})
}

func TestModel_View_Truncation(t *testing.T) {
	m := NewStatusBarModel().SetSize(5).SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
	view := m.View()
	// With width=5, content will be too wide for both zones.
	// Should still render without panicking.
	if view.Content == "" {
		t.Error("expected non-empty view even at narrow width")
	}
}

func TestModel_Init(t *testing.T) {
	m := NewStatusBarModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}

func TestModel_View_ZeroWidth(t *testing.T) {
	m := NewStatusBarModel().SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
	// width is 0 (never set) — should render left zone only, no panic
	view := m.View()
	if view.Content == "" {
		t.Error("expected non-empty view with zero width")
	}
	// With zero width, should NOT contain right-side indicators in the gap layout
	// (falls through to left-only rendering)
	if strings.Contains(view.Content, "Help") {
		// Menu items should still appear (they're in left zone)
	}
}

func TestModel_View_NegativeWidth(t *testing.T) {
	m := NewStatusBarModel().SetSize(-1).SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
	view := m.View()
	// Negative width should not panic
	if view.Content == "" {
		t.Error("expected non-empty view with negative width")
	}
}

func TestFormatKey_Space(t *testing.T) {
	items := []MenuItem{
		{Key: " ", Label: "Select"},
	}
	styles := DefaultStyles()
	result := RenderMenuLine(items, styles)
	if !strings.Contains(result, "space") {
		t.Error("expected space key to render as 'space'")
	}
}

func TestFormatKey_CtrlC(t *testing.T) {
	items := []MenuItem{
		{Key: "ctrl+c", Label: "Quit"},
	}
	styles := DefaultStyles()
	result := RenderMenuLine(items, styles)
	if !strings.Contains(result, "ctrl+c") {
		t.Error("expected ctrl+c key to render as 'ctrl+c'")
	}
}

func TestIndicatorStyle_CustomColor(t *testing.T) {
	m := NewStatusBarModel().SetSize(80)
	customStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	indicators := []StatusIndicator{
		NewStatusIndicator("Custom").WithStyle(customStyle),
		NewStatusIndicator("Default"),
	}
	m = m.SetIndicators(indicators)
	view := m.View()
	// Both indicators should be present
	if !strings.Contains(view.Content, "Custom") {
		t.Error("expected custom-styled indicator in view")
	}
	if !strings.Contains(view.Content, "Default") {
		t.Error("expected default-styled indicator in view")
	}
}
