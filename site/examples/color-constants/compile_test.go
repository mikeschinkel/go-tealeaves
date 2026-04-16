// Source: site/src/content/docs/components/color-constants.mdx:26#62d558e3,45#9a7d58cb,57#d981446c,69#388b98e3,80#f10b4507,112#0a3ea19d
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teacolor"
	lipgloss "charm.land/lipgloss/v2"
)

// TestCompile_ColorQuickExample verifies the quick example from color-constants.mdx.
func TestCompile_ColorQuickExample(t *testing.T) {
	// Use a named ANSI color
	style := lipgloss.NewStyle().Foreground(teacolor.Coral)
	_ = style

	// Create a SemanticColor for zero-allocation rendering
	accent := teacolor.NewSemanticColor(teacolor.Teal)
	rendered := accent.Render("important text")
	_ = rendered
}

// TestCompile_ANSI256 verifies ANSI 256 indexed colors.
func TestCompile_ANSI256(t *testing.T) {
	_ = teacolor.Color0
	_ = teacolor.Color1
	_ = teacolor.Color46
	_ = teacolor.Color255
}

// TestCompile_StandardNames verifies standard ANSI names (colors 0-15).
func TestCompile_StandardNames(t *testing.T) {
	_ = teacolor.Black
	_ = teacolor.Red
	_ = teacolor.Green
	_ = teacolor.Yellow
	_ = teacolor.Blue
	_ = teacolor.Magenta
	_ = teacolor.Cyan
	_ = teacolor.White
	_ = teacolor.BrightBlack
	_ = teacolor.BrightRed
	_ = teacolor.BrightGreen
	_ = teacolor.BrightYellow
	_ = teacolor.BrightBlue
	_ = teacolor.BrightMagenta
	_ = teacolor.BrightCyan
	_ = teacolor.BrightWhite
}

// TestCompile_SemanticAliases verifies curated semantic aliases.
func TestCompile_SemanticAliases(t *testing.T) {
	_ = teacolor.Coral
	_ = teacolor.SlateGray
	_ = teacolor.Teal
	_ = teacolor.Gold
	_ = teacolor.SteelBlue
	_ = teacolor.Lavender
}

// TestCompile_SemanticColor verifies SemanticColor methods.
func TestCompile_SemanticColor(t *testing.T) {
	accent := teacolor.NewSemanticColor(teacolor.Teal)

	text := accent.Render("highlighted text")
	_ = text

	bg := accent.Background().Render("background text")
	_ = bg

	border := accent.BorderForeground()
	_ = border

	style := lipgloss.NewStyle().Foreground(accent)
	_ = style
}

// TestCompile_ColorHelper verifies the Color() helper.
func TestCompile_ColorHelper(t *testing.T) {
	c := teacolor.Color("46")
	_ = c
}
