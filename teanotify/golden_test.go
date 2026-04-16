package teanotify

import (
	"testing"
	"time"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func TestNotifyModel_Golden_Idle(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    50,
		Duration: 3 * time.Second,
		Position: TopRightPosition,
	})
	// No active notification — renders background unchanged
	output := ansi.Strip(m.Render("Background content here"))
	golden.RequireEqual(t, []byte(output))
}
