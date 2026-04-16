package teaguide

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func TestGuideModel_Golden_Open(t *testing.T) {
	m := NewGuideModel()
	m = m.SetSize(80, 24)
	m, _ = m.Open(GuideData{
		Title: "Getting Started",
		Sections: []GuideSection{
			{
				Heading: "Recommended Next Steps",
				Items: []GuideItem{
					{KeyDisplay: "Enter", Label: "Select module to commit"},
					{KeyDisplay: "Tab", Label: "Switch between panes"},
					{KeyDisplay: "?", Label: "Show help"},
				},
			},
			{
				Heading: "Other Actions",
				Items: []GuideItem{
					{KeyDisplay: "q", Label: "Quit"},
					{KeyDisplay: "v", Label: "Run verification"},
				},
			},
		},
	})
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
