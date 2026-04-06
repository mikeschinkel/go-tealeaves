package teacolor

import (
	"testing"
)

func TestColor_ReturnsNonNil(t *testing.T) {
	c := Color("42")
	if c == nil {
		t.Fatal("Color(\"42\") returned nil")
	}
}

func TestColor_RGBA(t *testing.T) {
	c := Color("1")
	r, g, b, a := c.RGBA()
	// ANSI color 1 (red) should have non-zero red channel
	if r == 0 && g == 0 && b == 0 && a == 0 {
		t.Error("expected non-zero RGBA for ANSI color 1")
	}
}

func TestANSI256Constants_NonNil(t *testing.T) {
	constants := []struct {
		name  string
		value interface{}
	}{
		{"Color0", Color0},
		{"Color127", Color127},
		{"Color255", Color255},
	}
	for _, tc := range constants {
		if tc.value == nil {
			t.Errorf("%s is nil", tc.name)
		}
	}
}

func TestANSINames_NonNil(t *testing.T) {
	names := []struct {
		name  string
		value interface{}
	}{
		{"Black", Black},
		{"Red", Red},
		{"Green", Green},
		{"Yellow", Yellow},
		{"Blue", Blue},
		{"Magenta", Magenta},
		{"Cyan", Cyan},
		{"White", White},
		{"BrightBlack", BrightBlack},
		{"BrightRed", BrightRed},
		{"BrightGreen", BrightGreen},
		{"BrightYellow", BrightYellow},
		{"BrightBlue", BrightBlue},
		{"BrightMagenta", BrightMagenta},
		{"BrightCyan", BrightCyan},
		{"BrightWhite", BrightWhite},
	}
	for _, tc := range names {
		if tc.value == nil {
			t.Errorf("%s is nil", tc.name)
		}
	}
}

func TestNamedColors_NonNil(t *testing.T) {
	named := []struct {
		name  string
		value interface{}
	}{
		{"DarkGray", DarkGray},
		{"LightGray", LightGray},
		{"Coral", Coral},
		{"SkyBlue", SkyBlue},
		{"Gold", Gold},
		{"Crimson", Crimson},
		{"DodgerBlue", DodgerBlue},
		{"Teal", Teal},
		{"Salmon", Salmon},
		{"Olive", Olive},
		{"Plum", Plum},
		{"SlateGray", SlateGray},
		{"Indigo", Indigo},
	}
	for _, tc := range named {
		if tc.value == nil {
			t.Errorf("%s is nil", tc.name)
		}
	}
}

func TestHexColor_RGBA(t *testing.T) {
	c := Color("#FF0000")
	r, g, b, _ := c.RGBA()
	if r == 0 {
		t.Error("expected non-zero red channel for #FF0000")
	}
	if g != 0 || b != 0 {
		t.Error("expected zero green/blue channels for #FF0000")
	}
}
