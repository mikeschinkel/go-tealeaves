package teautils

import (
	"testing"
)

func TestDarkSystemPalette_AllFieldsNonNil(t *testing.T) {
	p := DarkSystemPalette(nil)
	fields := []struct {
		name  string
		value SemanticColor
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
		if f.value.IsZero() {
			t.Errorf("DarkSystemPalette(nil).%s is zero", f.name)
		}
	}
}

func TestLightSystemPalette_AllFieldsNonNil(t *testing.T) {
	p := LightSystemPalette(nil)
	fields := []struct {
		name  string
		value SemanticColor
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
		if f.value.IsZero() {
			t.Errorf("LightSystemPalette(nil).%s is zero", f.name)
		}
	}
}

func TestDarkSystemPalette_DifferentFromLight(t *testing.T) {
	dark := DarkSystemPalette(nil)
	light := LightSystemPalette(nil)

	// At minimum, primary text should differ (white vs black)
	dr, _, _, _ := dark.TextPrimary.RGBA()
	lr, _, _, _ := light.TextPrimary.RGBA()
	if dr == lr {
		t.Error("DarkSystemPalette and LightSystemPalette have same TextPrimary")
	}
}

func TestPalette_GenericComposition(t *testing.T) {
	type AppColors struct {
		CustomColor interface{}
	}
	sys := DarkSystemPalette(nil)
	p := Palette[AppColors]{
		System: sys,
		App:    AppColors{CustomColor: "test"},
	}
	// System field access should work
	if p.System.TextPrimary.IsZero() {
		t.Error("Palette[T].System.TextPrimary not accessible")
	}
	if p.System.Accent.IsZero() {
		t.Error("Palette[T].System.Accent not accessible")
	}
	// App field access should work
	if p.App.CustomColor != "test" {
		t.Error("Palette[T].App.CustomColor not accessible")
	}
}

func TestDarkSystemPalette_HighlightStyle(t *testing.T) {
	p := DarkSystemPalette(nil)
	if p.HighlightStyle == "" {
		t.Error("DarkSystemPalette(nil).HighlightStyle is empty")
	}
	if p.HighlightStyle != "monokai" {
		t.Errorf("DarkSystemPalette(nil).HighlightStyle = %q, want %q", p.HighlightStyle, "monokai")
	}
}

func TestLightSystemPalette_HighlightStyle(t *testing.T) {
	p := LightSystemPalette(nil)
	if p.HighlightStyle == "" {
		t.Error("LightSystemPalette(nil).HighlightStyle is empty")
	}
	if p.HighlightStyle != "github" {
		t.Errorf("LightSystemPalette(nil).HighlightStyle = %q, want %q", p.HighlightStyle, "github")
	}
}

func TestDefaultSystemPalette_ReturnsValid(t *testing.T) {
	p := DefaultSystemPalette(nil)
	if p.TextPrimary.IsZero() {
		t.Error("DefaultSystemPalette(nil).TextPrimary is zero")
	}
	if p.Accent.IsZero() {
		t.Error("DefaultSystemPalette(nil).Accent is zero")
	}
}

func TestDarkSystemPalette_Adaptive_AllFieldsNonNil(t *testing.T) {
	p := DarkSystemPalette(&PaletteOpts{Adaptive: true})
	fields := []struct {
		name  string
		value SemanticColor
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
		if f.value.IsZero() {
			t.Errorf("DarkSystemPalette(Adaptive).%s is zero", f.name)
		}
	}
}

func TestLightSystemPalette_Adaptive_AllFieldsNonNil(t *testing.T) {
	p := LightSystemPalette(&PaletteOpts{Adaptive: true})
	fields := []struct {
		name  string
		value SemanticColor
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
		if f.value.IsZero() {
			t.Errorf("LightSystemPalette(Adaptive).%s is zero", f.name)
		}
	}
}

func TestDarkSystemPalette_Adaptive_DifferentFromLight(t *testing.T) {
	dark := DarkSystemPalette(&PaletteOpts{Adaptive: true})
	light := LightSystemPalette(&PaletteOpts{Adaptive: true})

	// Primary text should differ (BrightWhite vs Black)
	dr, _, _, _ := dark.TextPrimary.RGBA()
	lr, _, _, _ := light.TextPrimary.RGBA()
	if dr == lr {
		t.Error("adaptive dark and light have same TextPrimary")
	}
}

func TestAdaptiveSystemPalette_Adaptive_ReturnsValid(t *testing.T) {
	p := AdaptiveSystemPalette(&PaletteOpts{Adaptive: true})
	if p.TextPrimary.IsZero() {
		t.Error("AdaptiveSystemPalette(Adaptive).TextPrimary is zero")
	}
	if p.Accent.IsZero() {
		t.Error("AdaptiveSystemPalette(Adaptive).Accent is zero")
	}
}

func TestSemanticColor_CachedStyles(t *testing.T) {
	sc := NewSemanticColor(nil)

	// Even nil-color SemanticColor should return valid styles
	_ = sc.Foreground()
	_ = sc.Background()
	_ = sc.BorderForeground()
	_ = sc.Render("test")

	// Non-nil color
	sc2 := DarkSystemPalette(nil).Accent
	fg := sc2.Foreground()
	if fg.GetForeground() == nil {
		t.Error("Foreground() style has no foreground color set")
	}
	bg := sc2.Background()
	if bg.GetBackground() == nil {
		t.Error("Background() style has no background color set")
	}
}

func TestSemanticColor_NilHandling(t *testing.T) {
	sc := NewSemanticColor(nil)

	if !sc.IsZero() {
		t.Error("NewSemanticColor(nil).IsZero() should be true")
	}
	if sc.Color() != nil {
		t.Error("NewSemanticColor(nil).Color() should be nil")
	}

	r, g, b, a := sc.RGBA()
	if r != 0 || g != 0 || b != 0 || a != 0 {
		t.Errorf("NewSemanticColor(nil).RGBA() = (%d,%d,%d,%d), want (0,0,0,0)", r, g, b, a)
	}
}

func TestSemanticColor_ImplementsColorInterface(t *testing.T) {
	sc := DarkSystemPalette(nil).TextPrimary
	r, _, _, a := sc.RGBA()
	if a == 0 {
		t.Error("non-nil SemanticColor.RGBA() returned zero alpha")
	}
	if r == 0 {
		// Color15 is white, should have non-zero red
		t.Error("TextPrimary (white) has zero red component")
	}
}
