package teautils

import (
	"testing"
)

func TestDarkPalette_AllFieldsNonNil(t *testing.T) {
	p := DarkPalette()
	fields := []struct {
		name  string
		value interface{}
	}{
		{"TextPrimary", p.TextPrimary},
		{"TextSecondary", p.TextSecondary},
		{"TextMuted", p.TextMuted},
		{"TextDim", p.TextDim},
		{"Accent", p.Accent},
		{"AccentAlt", p.AccentAlt},
		{"AccentSubtle", p.AccentSubtle},
		{"SelectionBg", p.SelectionBg},
		{"SelectionFg", p.SelectionFg},
		{"FocusBorder", p.FocusBorder},
		{"FocusBg", p.FocusBg},
		{"Border", p.Border},
		{"BorderAccent", p.BorderAccent},
		{"Separator", p.Separator},
		{"StatusSuccess", p.StatusSuccess},
		{"StatusInfo", p.StatusInfo},
		{"StatusWarn", p.StatusWarn},
		{"StatusError", p.StatusError},
		{"ButtonFg", p.ButtonFg},
		{"ButtonFocusBg", p.ButtonFocusBg},
		{"ButtonFocusFg", p.ButtonFocusFg},
		{"EditBg", p.EditBg},
		{"EditFg", p.EditFg},
		{"TintPositive", p.TintPositive},
		{"TintNegative", p.TintNegative},
		{"ScrollTrack", p.ScrollTrack},
		{"ScrollThumb", p.ScrollThumb},
	}
	for _, f := range fields {
		if f.value == nil {
			t.Errorf("DarkPalette().%s is nil", f.name)
		}
	}
}

func TestLightPalette_AllFieldsNonNil(t *testing.T) {
	p := LightPalette()
	fields := []struct {
		name  string
		value interface{}
	}{
		{"TextPrimary", p.TextPrimary},
		{"TextSecondary", p.TextSecondary},
		{"TextMuted", p.TextMuted},
		{"TextDim", p.TextDim},
		{"Accent", p.Accent},
		{"AccentAlt", p.AccentAlt},
		{"AccentSubtle", p.AccentSubtle},
		{"SelectionBg", p.SelectionBg},
		{"SelectionFg", p.SelectionFg},
		{"FocusBorder", p.FocusBorder},
		{"FocusBg", p.FocusBg},
		{"Border", p.Border},
		{"BorderAccent", p.BorderAccent},
		{"Separator", p.Separator},
		{"StatusSuccess", p.StatusSuccess},
		{"StatusInfo", p.StatusInfo},
		{"StatusWarn", p.StatusWarn},
		{"StatusError", p.StatusError},
		{"ButtonFg", p.ButtonFg},
		{"ButtonFocusBg", p.ButtonFocusBg},
		{"ButtonFocusFg", p.ButtonFocusFg},
		{"EditBg", p.EditBg},
		{"EditFg", p.EditFg},
		{"TintPositive", p.TintPositive},
		{"TintNegative", p.TintNegative},
		{"ScrollTrack", p.ScrollTrack},
		{"ScrollThumb", p.ScrollThumb},
	}
	for _, f := range fields {
		if f.value == nil {
			t.Errorf("LightPalette().%s is nil", f.name)
		}
	}
}

func TestDarkPalette_DifferentFromLight(t *testing.T) {
	dark := DarkPalette()
	light := LightPalette()

	// At minimum, primary text should differ (white vs black)
	dr, _, _, _ := dark.TextPrimary.RGBA()
	lr, _, _, _ := light.TextPrimary.RGBA()
	if dr == lr {
		t.Error("DarkPalette and LightPalette have same TextPrimary")
	}
}

func TestPalette_Embedding(t *testing.T) {
	type AppPalette struct {
		Palette
		CustomColor interface{}
	}
	p := DarkPalette()
	ap := AppPalette{
		Palette:     p,
		CustomColor: "test",
	}
	// Embedded field access should work
	if ap.TextPrimary == nil {
		t.Error("embedded Palette.TextPrimary not accessible")
	}
	if ap.Accent == nil {
		t.Error("embedded Palette.Accent not accessible")
	}
}

func TestDefaultPalette_ReturnsValid(t *testing.T) {
	p := DefaultPalette()
	if p.TextPrimary == nil {
		t.Error("DefaultPalette().TextPrimary is nil")
	}
	if p.Accent == nil {
		t.Error("DefaultPalette().Accent is nil")
	}
}
