package teastatus

import (
	"testing"

	"charm.land/bubbles/v2/key"
	"charm.land/lipgloss/v2"
)

func TestNewMenuItemFromBinding(t *testing.T) {
	binding := key.NewBinding(key.WithKeys("?", "h"), key.WithHelp("?", "Show help"))
	item := NewMenuItemFromBinding(binding, "Help")

	if item.Key != "?" {
		t.Errorf("expected Key='?' (first binding key), got %q", item.Key)
	}
	if item.Label != "Help" {
		t.Errorf("expected Label='Help', got %q", item.Label)
	}
}

func TestNewStatusIndicator(t *testing.T) {
	si := NewStatusIndicator("Ready")
	if si.Text != "Ready" {
		t.Errorf("expected Text='Ready', got %q", si.Text)
	}
}

func TestStatusIndicator_WithStyle(t *testing.T) {
	si := NewStatusIndicator("Active")
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	si2 := si.WithStyle(style)

	// Original should be unchanged
	_, isNoColor := si.Style.GetForeground().(lipgloss.NoColor)
	if !isNoColor {
		t.Error("expected original indicator to have no custom foreground")
	}

	// Copy should have the style
	_, isNoColor = si2.Style.GetForeground().(lipgloss.NoColor)
	if isNoColor {
		t.Error("expected WithStyle copy to have custom foreground")
	}
}
