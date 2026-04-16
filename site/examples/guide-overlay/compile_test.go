// Source: site/src/content/docs/components/guide-overlay.mdx:26#6ce86b4b,78#95332a15,169#f59ea972,178#0d3dacea,217#9dd2de53,228#7aa6a9f9,243#5d9563a8,257#42407fc9
package examples_test

import (
	"testing"

	key "charm.land/bubbles/v2/key"
	"github.com/mikeschinkel/go-tealeaves/teaguide"
)

// TestCompile_GuideQuickExample verifies the quick example from guide-overlay.mdx.
func TestCompile_GuideQuickExample(t *testing.T) {
	guide := teaguide.NewGuideModel()
	guide = guide.SetSize(80, 24)

	var cmd interface{}
	guide, cmd = guide.Open(teaguide.GuideData{
		Title: "What's Next?",
		Sections: []teaguide.GuideSection{
			{
				Priority: teaguide.PriorityRecommended,
				Heading:  "Recommended",
				Items: []teaguide.GuideItem{
					{
						ActionKey:  "enter",
						KeyDisplay: "[Enter]",
						Label:      "Select module",
						Prose:      "This module has uncommitted changes.",
					},
				},
			},
			{
				Priority: teaguide.PriorityAvailable,
				Heading:  "Also available",
				Items: []teaguide.GuideItem{
					{ActionKey: "r", KeyDisplay: "[r]", Label: "Refresh"},
				},
			},
		},
	})
	_ = cmd
	_ = guide.IsOpen()
}

// TestCompile_GuideDataStruct verifies GuideData and related types from guide-overlay.mdx.
func TestCompile_GuideDataStruct(t *testing.T) {
	data := teaguide.GuideData{
		Title: "Guide",
		Sections: []teaguide.GuideSection{
			{
				Priority: teaguide.PriorityBlocked,
				Heading:  "Blocked",
				Items: []teaguide.GuideItem{
					{
						Label:       "Publish package",
						BlockReason: "No remote configured.",
					},
				},
			},
		},
	}
	_ = data
}

// TestCompile_GuideKeyBindings verifies key binding replacement from guide-overlay.mdx.
func TestCompile_GuideKeyBindings(t *testing.T) {
	guide := teaguide.NewGuideModel()
	guide.Keys = teaguide.GuideKeyMap{
		ScrollUp:   key.NewBinding(key.WithKeys("w")),
		ScrollDown: key.NewBinding(key.WithKeys("s")),
	}
	guide.Styles = teaguide.DefaultGuideStyles()
	_ = guide
}

// TestCompile_ActionSelectedMsg verifies ActionSelectedMsg from guide-overlay.mdx.
func TestCompile_ActionSelectedMsg(t *testing.T) {
	msg := teaguide.ActionSelectedMsg{ActionKey: "enter"}
	_ = msg.ActionKey
}

// TestCompile_GuideDismissedMsg verifies GuideDismissedMsg type from guide-overlay.mdx.
func TestCompile_GuideDismissedMsg(t *testing.T) {
	var msg teaguide.GuideDismissedMsg
	_ = msg
}
