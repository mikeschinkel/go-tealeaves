// Source: site/src/content/docs/components/text-selection.mdx:21
package examples_test

import (
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teatext"
)

// TestCompile_TextSelectionQuickExample verifies the quick example from text-selection.mdx.
func TestCompile_TextSelectionQuickExample(t *testing.T) {
	editor := teatext.NewTextSnipModel(nil)
	editor.SetWidth(60)
	editor.SetHeight(20)
	editor.SetValue("Hello, world!\nSelect me with Shift+Arrow keys.")
	editor.Focus()
	_ = editor
}

// TestCompile_TextSnipModelArgs verifies TextSnipModelArgs from text-selection.mdx.
func TestCompile_TextSnipModelArgs(t *testing.T) {
	// Single-line mode
	singleLine := teatext.NewTextSnipModel(&teatext.TextSnipModelArgs{
		SingleLine: true,
	})
	_ = singleLine.IsSingleLine()
}

// TestCompile_SelectionType verifies Selection type and methods from text-selection.mdx.
func TestCompile_SelectionType(t *testing.T) {
	sel := teatext.NewSelection()
	_ = sel.IsEmpty()
	_ = sel.Active

	pos := teatext.Position{Row: 0, Col: 5}
	sel = sel.Begin(pos)
	sel = sel.Extend(teatext.Position{Row: 0, Col: 10})
	start, end := sel.Normalized()
	_, _ = start, end

	_ = sel.Contains(pos)
	sel = sel.Clear()
	_ = sel
}

// TestCompile_PositionType verifies Position type and methods from text-selection.mdx.
func TestCompile_PositionType(t *testing.T) {
	p1 := teatext.Position{Row: 0, Col: 5}
	p2 := teatext.Position{Row: 1, Col: 0}

	_ = p1.Before(p2)
	_ = p1.After(p2)
	_ = p1.Equal(p2)
}

// TestCompile_SelectAll verifies SelectAll from text-selection.mdx.
func TestCompile_SelectAll(t *testing.T) {
	lines := []string{"hello", "world"}
	sel := teatext.SelectAll(lines)
	_ = sel
}

// TestCompile_DefaultSelectionKeyMap verifies DefaultSelectionKeyMap from text-selection.mdx.
func TestCompile_DefaultSelectionKeyMap(t *testing.T) {
	km := teatext.DefaultSelectionKeyMap()
	_ = km
}

// TestCompile_SelectionStyle verifies SelectionStyle package variable from text-selection.mdx.
func TestCompile_SelectionStyle(t *testing.T) {
	style := teatext.SelectionStyle
	_ = style

	// Override the default selection style
	teatext.SelectionStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("25")).
		Foreground(lipgloss.Color("255"))
}

// TestCompile_ModelMethods verifies TextSnipModel methods from text-selection.mdx.
func TestCompile_ModelMethods(t *testing.T) {
	editor := teatext.NewTextSnipModel(nil)
	editor.SetValue("Hello, world!")

	_ = editor.Selection()
	_ = editor.HasSelection()
	_ = editor.SelectedText()

	sel := teatext.NewSelection()
	editor = editor.SetSelection(sel)
	editor = editor.ClearSelection()
	editor = editor.Copy()
	editor = editor.Cut()
}
