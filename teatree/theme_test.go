package teatree

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestModel_WithTheme(t *testing.T) {
	root := NewNode("1", "root", "data")
	tree := NewTree([]*Node[string]{root}, nil)
	m := NewModel(tree, 10)

	if m.Theme() != nil {
		t.Error("expected nil theme before WithTheme")
	}

	theme := teautils.NewTheme(teautils.DarkPalette())
	m = m.WithTheme(theme)

	if m.Theme() == nil {
		t.Error("expected non-nil theme after WithTheme")
	}
	if m.Theme().Palette.TextPrimary == nil {
		t.Error("theme palette TextPrimary is nil")
	}
}
