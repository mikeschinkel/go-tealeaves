// Source: site/src/content/docs/components/help-visor.mdx:22#ab966188,125#df64d61b,136#ff8aefd8
package examples_test

import (
	"testing"

	key "charm.land/bubbles/v2/key"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teacolor"
	"github.com/mikeschinkel/go-tealeaves/teahelp"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_HelpVisorQuickExample verifies the quick example from help-visor.mdx.
func TestCompile_HelpVisorQuickExample(t *testing.T) {
	help := teahelp.NewHelpVisorModel()
	help = help.SetSize(80, 24)

	upKey := key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "Move up"))
	downKey := key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "Move down"))
	enterKey := key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "Select"))
	escKey := key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "Cancel"))

	keysByCategory := map[string][]teautils.KeyMeta{
		"Navigation": {
			{Binding: upKey, HelpText: "Move up"},
			{Binding: downKey, HelpText: "Move down"},
		},
		"Actions": {
			{Binding: enterKey, HelpText: "Select"},
			{Binding: escKey, HelpText: "Cancel"},
		},
	}
	help = help.Open(keysByCategory)
	_ = help.IsOpen()
}

// TestCompile_HelpVisorKeyMap verifies WithKeys from help-visor.mdx.
func TestCompile_HelpVisorKeyMap(t *testing.T) {
	km := teahelp.DefaultHelpVisorKeyMap()
	km.Close = key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "close help"))
	help := teahelp.NewHelpVisorModel().WithKeys(km)
	_ = help
}

// TestCompile_HelpVisorContentStyle verifies WithContentStyle from help-visor.mdx.
func TestCompile_HelpVisorContentStyle(t *testing.T) {
	help := teahelp.NewHelpVisorModel().WithContentStyle(teautils.HelpVisorStyle{
		TitleStyle:    lipgloss.NewStyle().Bold(true),
		CategoryStyle: lipgloss.NewStyle().Foreground(teacolor.Teal),
		KeyStyle:      lipgloss.NewStyle().Foreground(teacolor.BrightWhite),
		DescStyle:     lipgloss.NewStyle().Foreground(teacolor.SlateGray),
		KeyColumnGap:  4,
		CategoryOrder: []string{"Navigation", "Actions", "System"},
	})
	_ = help
}
