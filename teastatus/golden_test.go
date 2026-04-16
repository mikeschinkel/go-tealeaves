package teastatus

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func TestStatusBarModel_Golden_WithItems(t *testing.T) {
	m := NewStatusBarModel()
	m = m.SetSize(80)
	m = m.SetMenuItems([]MenuItem{
		{KeyText: "?", Label: "Help"},
		{KeyText: "q", Label: "Quit"},
		{KeyText: "Tab", Label: "Focus"},
	})
	m = m.SetIndicators([]StatusIndicator{
		{Text: "3 files"},
		{Text: "VERIFIED"},
	})
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
