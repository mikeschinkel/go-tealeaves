package teagrid

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestWithTheme_AppliesGridStyles(t *testing.T) {
	cols := []Column{NewColumn("id", "ID", 5)}
	m := NewGridModel(cols)

	theme := teautils.NewTheme(teautils.DarkPalette())
	m = m.WithTheme(theme)

	if m.highlightStyle.GetBackground() == nil {
		t.Error("themed highlightStyle has no background")
	}
	if m.headerStyle.GetForeground() == nil {
		t.Error("themed headerStyle has no foreground")
	}
}

func TestWithTheme_ThenWithHighlightStyle_Overrides(t *testing.T) {
	cols := []Column{NewColumn("id", "ID", 5)}
	theme := teautils.NewTheme(teautils.DarkPalette())
	m := NewGridModel(cols).WithTheme(theme).WithHighlightStyle(defaultHighlightStyle)

	// WithHighlightStyle should override theme
	if m.highlightStyle.GetBackground() != defaultHighlightStyle.GetBackground() {
		t.Error("WithHighlightStyle did not override WithTheme")
	}
}
