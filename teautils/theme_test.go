package teautils

import (
	"testing"
)

func TestNewTheme_FromDarkPalette(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	// Verify palette is stored
	if theme.Palette.TextPrimary == nil {
		t.Error("theme.Palette.TextPrimary is nil")
	}

	// Verify common styles are non-zero (have foreground set)
	if theme.Title.GetForeground() == nil {
		t.Error("theme.Title has no foreground")
	}
	if theme.Message.GetForeground() == nil {
		t.Error("theme.Message has no foreground")
	}
}

func TestNewTheme_StatusBarStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.StatusBar.MenuKeyStyle.GetForeground() == nil {
		t.Error("StatusBar.MenuKeyStyle has no foreground")
	}
	if theme.StatusBar.MenuLabelStyle.GetForeground() == nil {
		t.Error("StatusBar.MenuLabelStyle has no foreground")
	}
	if theme.StatusBar.IndicatorSepStyle.GetForeground() == nil {
		t.Error("StatusBar.IndicatorSepStyle has no foreground")
	}
}

func TestNewTheme_HelpVisorStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.HelpVisor.TitleStyle.GetForeground() == nil {
		t.Error("HelpVisor.TitleStyle has no foreground")
	}
	if theme.HelpVisor.KeyStyle.GetForeground() == nil {
		t.Error("HelpVisor.KeyStyle has no foreground")
	}
	if theme.HelpVisor.DescStyle.GetForeground() == nil {
		t.Error("HelpVisor.DescStyle has no foreground")
	}
}

func TestNewTheme_ModalStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.Modal.TitleStyle.GetForeground() == nil {
		t.Error("Modal.TitleStyle has no foreground")
	}
	if theme.Modal.FocusedButtonStyle.GetBackground() == nil {
		t.Error("Modal.FocusedButtonStyle has no background")
	}
}

func TestNewTheme_DropdownStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.Dropdown.ItemStyle.GetForeground() == nil {
		t.Error("Dropdown.ItemStyle has no foreground")
	}
	if theme.Dropdown.SelectedStyle.GetBackground() == nil {
		t.Error("Dropdown.SelectedStyle has no background")
	}
}

func TestNewTheme_ListStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.List.ItemStyle.GetForeground() == nil {
		t.Error("List.ItemStyle has no foreground")
	}
	if theme.List.ActiveItemStyle.GetForeground() == nil {
		t.Error("List.ActiveItemStyle has no foreground")
	}
	if theme.List.EditItemStyle.GetBackground() == nil {
		t.Error("List.EditItemStyle has no background")
	}
}

func TestNewTheme_GridStyles(t *testing.T) {
	p := DarkPalette()
	theme := NewTheme(p)

	if theme.Grid.HighlightStyle.GetBackground() == nil {
		t.Error("Grid.HighlightStyle has no background")
	}
	if theme.Grid.BorderStyle.GetForeground() == nil {
		t.Error("Grid.BorderStyle has no foreground")
	}
}

func TestDefaultTheme_ReturnsValid(t *testing.T) {
	theme := DefaultTheme()
	if theme.Palette.TextPrimary == nil {
		t.Error("DefaultTheme().Palette.TextPrimary is nil")
	}
	if theme.Title.GetForeground() == nil {
		t.Error("DefaultTheme().Title has no foreground")
	}
}

func TestNewTheme_FromLightPalette(t *testing.T) {
	p := LightPalette()
	theme := NewTheme(p)

	// Should use light palette colors
	if theme.Palette.TextPrimary == nil {
		t.Error("light theme.Palette.TextPrimary is nil")
	}
	if theme.Title.GetForeground() == nil {
		t.Error("light theme.Title has no foreground")
	}
}
