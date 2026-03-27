package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-cliutil"

	"github.com/mikeschinkel/go-tealeaves/teaguide"
)

// workflowState represents the current deployment workflow state.
type workflowState int

const (
	stateNotTested workflowState = iota
	stateTested
	stateDeployed
)

func (s workflowState) String() string {
	switch s {
	case stateNotTested:
		return "Not Tested"
	case stateTested:
		return "Tested"
	case stateDeployed:
		return "Deployed"
	default:
		return "Unknown"
	}
}

type model struct {
	guide  teaguide.GuideModel
	state  workflowState
	width  int
	height int
	status string
}

func newModel() model {
	return model{
		guide: teaguide.NewGuideModel(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Track if guide was open before update (modal key consumption pattern)
	guideWasOpen := m.guide.IsOpen()

	// Let guide handle messages first
	if guideWasOpen {
		var guideModel tea.Model
		guideModel, cmd = m.guide.Update(msg)
		m.guide = guideModel.(teaguide.GuideModel)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.guide = m.guide.SetSize(msg.Width, msg.Height)
		goto end

	case teaguide.ActionSelectedMsg:
		m = m.handleAction(msg.ActionKey)
		goto end

	case teaguide.GuideDismissedMsg:
		m.status = "Guide dismissed"
		goto end

	case tea.KeyPressMsg:
		if guideWasOpen {
			// Guide consumed the key
			goto end
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "n", "N":
			m.guide, cmd = m.guide.Open(m.buildGuideData())
			goto end
		case "t":
			m = m.handleAction("t")
			goto end
		case "d":
			m = m.handleAction("d")
			goto end
		case "l":
			m = m.handleAction("l")
			goto end
		case "r":
			m = m.handleAction("r")
			goto end
		}
	}

end:
	return m, cmd
}

func (m model) handleAction(actionKey string) model {
	switch actionKey {
	case "t":
		if m.state == stateNotTested || m.state == stateTested {
			m.state = stateTested
			m.status = "Tests passed!"
		}
	case "d":
		m = m.handleDeploy()
	case "l":
		m = m.handleRelease()
	case "r":
		m.status = "Data refreshed"
	}
	return m
}

func (m model) handleDeploy() model {
	if m.state == stateTested || m.state == stateDeployed {
		m.state = stateDeployed
		m.status = "Deployed successfully!"
		return m
	}
	m.status = "Cannot deploy: run tests first"
	return m
}

func (m model) handleRelease() model {
	if m.state == stateDeployed {
		m.status = "Released! Workflow complete."
		return m
	}
	m.status = "Cannot release: deploy first"
	return m
}

func (m model) buildGuideData() teaguide.GuideData {
	data := teaguide.GuideData{
		Title: "What's Next?",
	}

	switch m.state {
	case stateNotTested:
		data.Sections = []teaguide.GuideSection{
			{
				Priority: teaguide.PriorityRecommended,
				Heading:  "Recommended Next",
				Items: []teaguide.GuideItem{
					{
						ActionKey:  "t",
						KeyDisplay: "[T]",
						Label:      "Run Tests",
						Prose:      "Tests have not been run yet. Run them to verify the project builds correctly.",
					},
				},
			},
			{
				Priority: teaguide.PriorityAvailable,
				Heading:  "Also Available",
				Items: []teaguide.GuideItem{
					{ActionKey: "r", KeyDisplay: "[R]", Label: "Refresh Data"},
					{ActionKey: "q", KeyDisplay: "[Q]", Label: "Quit"},
				},
			},
			{
				Priority: teaguide.PriorityBlocked,
				Heading:  "Not Yet Available",
				Items: []teaguide.GuideItem{
					{Label: "Deploy", BlockReason: "run tests first"},
					{Label: "Release", BlockReason: "deploy first"},
				},
			},
		}

	case stateTested:
		data.Sections = []teaguide.GuideSection{
			{
				Priority: teaguide.PriorityRecommended,
				Heading:  "Recommended Next",
				Items: []teaguide.GuideItem{
					{
						ActionKey:  "d",
						KeyDisplay: "[D]",
						Label:      "Deploy",
						Prose:      "Tests passed. Deploy to staging to verify in a real environment.",
					},
				},
			},
			{
				Priority: teaguide.PriorityAvailable,
				Heading:  "Also Available",
				Items: []teaguide.GuideItem{
					{ActionKey: "t", KeyDisplay: "[T]", Label: "Re-run Tests"},
					{ActionKey: "r", KeyDisplay: "[R]", Label: "Refresh Data"},
					{ActionKey: "q", KeyDisplay: "[Q]", Label: "Quit"},
				},
			},
			{
				Priority: teaguide.PriorityBlocked,
				Heading:  "Not Yet Available",
				Items: []teaguide.GuideItem{
					{Label: "Release", BlockReason: "deploy first"},
				},
			},
		}

	case stateDeployed:
		data.Sections = []teaguide.GuideSection{
			{
				Priority: teaguide.PriorityRecommended,
				Heading:  "Recommended Next",
				Items: []teaguide.GuideItem{
					{
						ActionKey:  "l",
						KeyDisplay: "[L]",
						Label:      "Release",
						Prose:      "Deployment verified. Create a release to publish the new version.",
					},
				},
			},
			{
				Priority: teaguide.PriorityAvailable,

				Heading: "Also Available",
				Items: []teaguide.GuideItem{
					{ActionKey: "d", KeyDisplay: "[D]", Label: "Re-deploy"},
					{ActionKey: "t", KeyDisplay: "[T]", Label: "Re-run Tests"},
					{ActionKey: "q", KeyDisplay: "[Q]", Label: "Quit"},
				},
			},
		}
	}

	return data
}

func (m model) View() tea.View {
	var bg strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("219"))
	stateStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("114"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	statusStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))

	bg.WriteString("\n")
	bg.WriteString(titleStyle.Render("  teaguide Example — Project Deployment Workflow"))
	bg.WriteString("\n\n")
	bg.WriteString(fmt.Sprintf("  Current State: %s", stateStyle.Render(m.state.String())))
	bg.WriteString("\n\n")

	if m.status != "" {
		bg.WriteString(fmt.Sprintf("  %s", statusStyle.Render(m.status)))
		bg.WriteString("\n\n")
	}

	bg.WriteString(hintStyle.Render("  [N] Next"))
	bg.WriteString("\n")

	// Pad to fill screen
	lines := strings.Count(bg.String(), "\n") + 1
	for i := lines; i < m.height; i++ {
		bg.WriteString("\n")
	}

	view := bg.String()
	view = m.guide.OverlayModal(view)

	v := tea.NewView(view)
	v.AltScreen = true
	return v
}

func main() {
	p := tea.NewProgram(newModel())
	_, err := p.Run()
	if err != nil {
		cliutil.Stderrf("Error: %v\n", err)
		os.Exit(1)
	}
}
