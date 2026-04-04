package examples_test

import (
	"fmt"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teacolor"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
	"github.com/mikeschinkel/go-tealeaves/teapane"
)

// TestCompile_PaneWidgetsQuickExample verifies the quick example from pane-widgets.mdx.
func TestCompile_PaneWidgetsQuickExample(t *testing.T) {
	border := teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        teacolor.SlateGray,
		FocusedColor: teacolor.Teal,
		PaddingH:     1,
	}

	pane := teapane.NewStyledPane(border, func(w, h int, focused bool) string {
		return fmt.Sprintf("content %dx%d focused=%v", w, h, focused)
	})

	layout := tealayout.NewLayout(
		tealayout.NewRow(tealayout.Flex(1),
			tealayout.NewElement(pane),
		),
	)
	layout.SetSize(80, 24)
	_, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
}

// TestCompile_BorderStyle verifies BorderStyle type from pane-widgets.mdx.
func TestCompile_BorderStyle(t *testing.T) {
	bs := teapane.BorderStyle{
		Border:     lipgloss.RoundedBorder(),
		Color:      teacolor.SlateGray,
		PaddingH:   1,
		PaddingV:   0,
	}
	_ = bs.FrameWidth()
	_ = bs.FrameHeight()
}

// TestCompile_StyledPaneModifiers verifies StyledPane modifiers from pane-widgets.mdx.
func TestCompile_StyledPaneModifiers(t *testing.T) {
	border := teapane.BorderStyle{Border: lipgloss.RoundedBorder()}
	pane := teapane.NewStyledPane(border, func(w, h int, focused bool) string {
		return ""
	})
	pane = pane.WithMinWidth(20).WithLabel("sidebar")
	_ = pane.Width()
	_ = pane.Height()
	_ = pane.Focused()
	_ = pane.Label()
}

// TestCompile_ScrollPane verifies ScrollPane from pane-widgets.mdx.
func TestCompile_ScrollPane(t *testing.T) {
	border := teapane.BorderStyle{Border: lipgloss.RoundedBorder()}
	scroll := teapane.NewScrollPane(border, func(w, h, offset int) string {
		return fmt.Sprintf("lines starting at %d", offset)
	})

	scroll.SetTotalLines(100)
	_ = scroll.ScrollOffset()
	scroll.ScrollUp()
	scroll.ScrollDown()
}

// TestCompile_PlainPane verifies PlainPane from pane-widgets.mdx.
func TestCompile_PlainPane(t *testing.T) {
	header := teapane.NewPlainPane(func(w, h int, _ bool) string {
		return lipgloss.NewStyle().Bold(true).Render("My App v1.0")
	})

	header = header.WithStyle(lipgloss.NewStyle().Background(teacolor.SlateGray))
	_ = header
}
