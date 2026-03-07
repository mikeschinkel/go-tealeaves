package teatree

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestModel_WithTheme(t *testing.T) {
	root := NewNode("1", "root", "data")
	tree := NewTree([]*Node[string]{root}, nil)
	m := NewTreeModel(tree, 10)

	if m.Theme() != nil {
		t.Error("expected nil theme before WithTheme")
	}

	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	if m.Theme() == nil {
		t.Error("expected non-nil theme after WithTheme")
	}
	if m.Theme().System.TextPrimary.IsZero() {
		t.Error("theme palette TextPrimary is zero")
	}
}
