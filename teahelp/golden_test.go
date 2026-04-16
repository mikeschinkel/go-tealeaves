package teahelp

import (
	"testing"

	"charm.land/bubbles/v2/key"
	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestHelpVisorModel_Golden_Open(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 24)
	keys := map[string][]teautils.KeyMeta{
		"Navigation": {
			{ID: "nav.up", Binding: key.NewBinding(key.WithKeys("up")), HelpText: "Move up"},
			{ID: "nav.down", Binding: key.NewBinding(key.WithKeys("down")), HelpText: "Move down"},
			{ID: "nav.select", Binding: key.NewBinding(key.WithKeys("enter")), HelpText: "Select item"},
		},
		"Actions": {
			{ID: "act.quit", Binding: key.NewBinding(key.WithKeys("q")), HelpText: "Quit"},
			{ID: "act.help", Binding: key.NewBinding(key.WithKeys("?")), HelpText: "Toggle help"},
		},
	}
	m = m.Open(keys)
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
