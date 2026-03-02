package teastatus

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
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
	m := New()
	// Should have default styles (non-zero)
	if m.Styles.MenuSeparator == "" {
		t.Error("expected non-empty MenuSeparator in default styles")
	}
}

func TestModel_SetSize(t *testing.T) {
	m := New()
	m = m.SetSize(80)
	if m.width != 80 {
		t.Errorf("expected width=80, got %d", m.width)
	}
}

func TestModel_SetMenuItems(t *testing.T) {
	m := New()
	items := testMenuItems()
	m = m.SetMenuItems(items)
	if len(m.menuItems) != 2 {
		t.Errorf("expected 2 menu items, got %d", len(m.menuItems))
	}
}

func TestModel_SetIndicators(t *testing.T) {
	m := New()
	indicators := testIndicators()
	m = m.SetIndicators(indicators)
	if len(m.indicators) != 2 {
		t.Errorf("expected 2 indicators, got %d", len(m.indicators))
	}
}

func TestModel_Update_SetMenuItemsMsg(t *testing.T) {
	m := New()
	items := testMenuItems()
	result, cmd := m.Update(SetMenuItemsMsg{Items: items})
	m = result.(Model)

	if len(m.menuItems) != 2 {
		t.Errorf("expected 2 menu items after msg, got %d", len(m.menuItems))
	}
	if cmd != nil {
		t.Error("expected nil cmd from SetMenuItemsMsg")
	}
}

func TestModel_Update_SetIndicatorsMsg(t *testing.T) {
	m := New()
	indicators := testIndicators()
	result, cmd := m.Update(SetIndicatorsMsg{Indicators: indicators})
	m = result.(Model)

	if len(m.indicators) != 2 {
		t.Errorf("expected 2 indicators after msg, got %d", len(m.indicators))
	}
	if cmd != nil {
		t.Error("expected nil cmd from SetIndicatorsMsg")
	}
}

func TestModel_Update_UnknownMsg(t *testing.T) {
	m := New()
	m = m.SetMenuItems(testMenuItems())
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	updated := result.(Model)

	if len(updated.menuItems) != 2 {
		t.Errorf("expected menu items unchanged, got %d", len(updated.menuItems))
	}
	if cmd != nil {
		t.Error("expected nil cmd for unknown message")
	}
}

// --- Layer 2 Tests ---

func TestModel_View_Empty(t *testing.T) {
	m := New().SetSize(80)
	view := m.View()
	// Should render something (bar background) even with no items
	// At minimum it shouldn't panic
	_ = view
}

func TestModel_View_MenuItems(t *testing.T) {
	m := New().SetSize(80).SetMenuItems(testMenuItems())
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
	m := New().SetSize(80).SetIndicators(testIndicators())
	view := m.View()
	if !strings.Contains(view.Content, "Ready") {
		t.Error("expected view to contain 'Ready'")
	}
	if !strings.Contains(view.Content, "3 items") {
		t.Error("expected view to contain '3 items'")
	}
}

func TestModel_View_BothZones(t *testing.T) {
	m := New().SetSize(80).SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
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
		m := New().WithStyles(styles).SetSize(80).SetIndicators(indicators)
		view := m.View()
		if !strings.Contains(view.Content, "|") {
			t.Error("expected pipe separator in view")
		}
	})

	t.Run("BracketSeparator", func(t *testing.T) {
		styles := DefaultStyles()
		styles.SeparatorKind = BracketSeparator
		m := New().WithStyles(styles).SetSize(80).SetIndicators(indicators)
		view := m.View()
		if !strings.Contains(view.Content, "[Ready]") {
			t.Error("expected bracket-wrapped indicator in view")
		}
	})

	t.Run("SpaceSeparator", func(t *testing.T) {
		styles := DefaultStyles()
		styles.SeparatorKind = SpaceSeparator
		m := New().WithStyles(styles).SetSize(80).SetIndicators(indicators)
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
	m := New().SetSize(5).SetMenuItems(testMenuItems()).SetIndicators(testIndicators())
	view := m.View()
	// With width=5, content will be too wide for both zones.
	// Should still render without panicking.
	if view.Content == "" {
		t.Error("expected non-empty view even at narrow width")
	}
}
