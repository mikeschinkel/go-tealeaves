package teapane_test

import (
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teapane"
)

func TestBorderStyle_Build(t *testing.T) {
	bs := teapane.BorderStyle{
		Border:       lipgloss.RoundedBorder(),
		Color:        lipgloss.Color("#ff0000"),
		FocusedColor: lipgloss.Color("#ffffff"),
		Foreground:   lipgloss.Color("#aabbcc"),
		PaddingH:     1,
		PaddingV:     0,
	}

	normal := bs.Build(false)
	if normal.GetBorderTopForeground() != lipgloss.Color("#ff0000") {
		t.Errorf("expected normal border color #ff0000, got %v", normal.GetBorderTopForeground())
	}

	focused := bs.Build(true)
	if focused.GetBorderTopForeground() != lipgloss.Color("#ffffff") {
		t.Errorf("expected focused border color #ffffff, got %v", focused.GetBorderTopForeground())
	}
}

func TestBorderStyle_FrameDimensions(t *testing.T) {
	bs := teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
		PaddingV: 0,
	}

	// Rounded border = 1 cell on each side; padding 1 on each side horizontally
	fw := bs.FrameWidth()
	if fw != 4 { // 2 border + 2 padding
		t.Errorf("expected FrameWidth 4, got %d", fw)
	}

	fh := bs.FrameHeight()
	if fh != 2 { // 2 border + 0 padding
		t.Errorf("expected FrameHeight 2, got %d", fh)
	}
}
