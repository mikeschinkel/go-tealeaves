package teapane_test

import (
	"fmt"
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
	"github.com/mikeschinkel/go-tealeaves/teapane"
)

func TestStyledPane_InterfaceDetection(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, func(w, h int, focused bool) string {
		return "hello"
	})

	// NewElement should detect all 5 interfaces.
	elem := tealayout.NewElement(sp)
	_ = elem // If this compiles, the types are satisfied.

	// Verify it's detected as ContentProvider
	var _ tealayout.ContentProvider = sp
	var _ tealayout.SetSizer = sp
	var _ tealayout.Styler = sp
	var _ tealayout.SizeHinter = sp
	var _ tealayout.Focusable = sp
}

func TestStyledPane_SetSizeAndView(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, func(w, h int, focused bool) string {
		return fmt.Sprintf("%dx%d", w, h)
	})

	sp.SetSize(20, 10)
	if sp.Width() != 20 {
		t.Errorf("expected width 20, got %d", sp.Width())
	}
	if sp.Height() != 10 {
		t.Errorf("expected height 10, got %d", sp.Height())
	}

	view := sp.Content()
	if !strings.Contains(view, "20x10") {
		t.Errorf("expected view to contain '20x10', got:\n%s", view)
	}
}

func TestStyledPane_FocusBlur(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        lipgloss.Color("#ff0000"),
		FocusedColor: lipgloss.Color("#ffffff"),
		PaddingH:     1,
	}, func(w, h int, focused bool) string {
		if focused {
			return "FOCUSED"
		}
		return "BLURRED"
	})

	if sp.Focused() {
		t.Error("expected not focused initially")
	}

	sp.Focus()
	if !sp.Focused() {
		t.Error("expected focused after Focus()")
	}
	style := sp.Style()
	if style.GetBorderTopForeground() != lipgloss.Color("#ffffff") {
		t.Errorf("expected white border when focused, got %v", style.GetBorderTopForeground())
	}

	sp.Blur()
	if sp.Focused() {
		t.Error("expected not focused after Blur()")
	}
	style = sp.Style()
	if style.GetBorderTopForeground() != lipgloss.Color("#ff0000") {
		t.Errorf("expected red border when blurred, got %v", style.GetBorderTopForeground())
	}
}

func TestStyledPane_SizeHintWithMinWidth(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, nil).WithMinWidth(20)

	hint := sp.SizeHint(100, 50)
	// frame = 4 (2 border + 2 padding), min content = 20
	expected := 24
	if hint.Min.Width != expected {
		t.Errorf("expected SizeHint.Min.Width %d, got %d", expected, hint.Min.Width)
	}
}

func TestStyledPane_SizeHintFunc(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, nil).WithMinWidth(10).WithSizeHintFunc(func(availW, availH int) tealayout.SizeHint {
		return tealayout.SizeHint{
			Desired: tealayout.Size{Width: 42, Height: availH},
			Max:     tealayout.Size{Width: -1, Height: -1},
		}
	})

	hint := sp.SizeHint(100, 50)
	if hint.Desired.Width != 42 {
		t.Errorf("expected SizeHint.Desired.Width 42, got %d", hint.Desired.Width)
	}
	// minWidth should be ignored when SizeHintFunc is set
	if hint.Min.Width == 14 {
		t.Error("SizeHintFunc should override minWidth-based calculation")
	}
}

func TestStyledPane_SizeHintFunc_FallbackToMinWidth(t *testing.T) {
	sp := teapane.NewStyledPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, nil).WithMinWidth(20)

	// No SizeHintFunc set — should fall back to minWidth-based calculation
	hint := sp.SizeHint(100, 50)
	expected := 24 // frame(4) + minWidth(20)
	if hint.Min.Width != expected {
		t.Errorf("expected fallback SizeHint.Min.Width %d, got %d", expected, hint.Min.Width)
	}
}
