package teapane_test

import (
	"fmt"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teapane"
)

func TestScrollPane_OffsetClamping(t *testing.T) {
	sp := teapane.NewScrollPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, func(w, h, offset int) string {
		return fmt.Sprintf("offset=%d", offset)
	})

	sp.SetSize(20, 10)
	sp.SetTotalLines(15)

	// Initial offset should be 0
	if sp.ScrollOffset() != 0 {
		t.Errorf("expected initial offset 0, got %d", sp.ScrollOffset())
	}

	// Scroll down 3 times
	sp.ScrollDown()
	sp.ScrollDown()
	sp.ScrollDown()
	if sp.ScrollOffset() != 3 {
		t.Errorf("expected offset 3, got %d", sp.ScrollOffset())
	}

	// Scroll up once
	sp.ScrollUp()
	if sp.ScrollOffset() != 2 {
		t.Errorf("expected offset 2, got %d", sp.ScrollOffset())
	}

	// Scroll past max (totalLines=15, height=10, max offset=5)
	for range 20 {
		sp.ScrollDown()
	}
	if sp.ScrollOffset() != 5 {
		t.Errorf("expected offset clamped to 5, got %d", sp.ScrollOffset())
	}

	// ScrollUp at 0 stays at 0
	for range 10 {
		sp.ScrollUp()
	}
	if sp.ScrollOffset() != 0 {
		t.Errorf("expected offset clamped to 0, got %d", sp.ScrollOffset())
	}
}

func TestScrollPane_SetTotalLinesClampsOffset(t *testing.T) {
	sp := teapane.NewScrollPane(teapane.BorderStyle{
		Border:   lipgloss.RoundedBorder(),
		PaddingH: 1,
	}, func(w, h, offset int) string {
		return ""
	})

	sp.SetSize(20, 10)
	sp.SetTotalLines(20)

	// Scroll to offset 10
	for range 10 {
		sp.ScrollDown()
	}
	if sp.ScrollOffset() != 10 {
		t.Errorf("expected offset 10, got %d", sp.ScrollOffset())
	}

	// Reduce total lines — offset should clamp
	sp.SetTotalLines(12)
	if sp.ScrollOffset() != 2 { // max = 12-10 = 2
		t.Errorf("expected offset clamped to 2, got %d", sp.ScrollOffset())
	}
}
