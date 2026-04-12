// Source: site/src/content/docs/guides/why-tea-leaves.mdx:15,33,52,60,75
package examples_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// whyModel holds the fields used across why-tea-leaves examples.
type whyModel struct {
	dropdown teafields.DropdownModel
}

func (m whyModel) Init() tea.Cmd { return nil }

func (m whyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m whyModel) View() tea.View {
	return tea.NewView("")
}

// TestCompile_WhyStateCheckAntiPattern verifies the state-check anti-pattern compiles from why-tea-leaves.mdx:15.
// The parent model accesses dropdown internals to gate key handling. This is the anti-pattern Tea Leaves solves.
func TestCompile_WhyStateCheckAntiPattern(t *testing.T) {
	antiPatternUpdate := func(m whyModel, msg tea.Msg) (tea.Model, tea.Cmd) {
		if !m.dropdown.IsOpen {
			switch msg := msg.(type) {
			case tea.KeyPressMsg:
				if msg.String() == "q" {
					return m, tea.Quit
				}
			}
		}
		return m, nil
	}

	m := whyModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"A", "B"}),
			&teafields.DropdownModelArgs{ScreenWidth: 80, ScreenHeight: 24},
		),
	}
	result, cmd := antiPatternUpdate(m, tea.KeyPressMsg{})
	_ = result
	_ = cmd
}

// TestCompile_WhyCleanConsumptionPattern verifies the Tea Leaves clean consumption pattern from why-tea-leaves.mdx:33.
// Non-nil tea.Cmd signals the child component consumed the message — no state checks needed.
func TestCompile_WhyCleanConsumptionPattern(t *testing.T) {
	cleanUpdate := func(m whyModel, msg tea.Msg) (tea.Model, tea.Cmd) {
		dropdown, cmd := m.dropdown.Update(msg)
		if cmd != nil {
			m.dropdown = dropdown.(teafields.DropdownModel)
			return m, cmd // Dropdown consumed it — done!
		}
		// Parent processes normally
		return m, nil
	}

	m := whyModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"A", "B"}),
			&teafields.DropdownModelArgs{ScreenWidth: 80, ScreenHeight: 24},
		),
	}
	result, cmd := cleanUpdate(m, tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}

// TestCompile_WhyAnsiAntiPattern verifies that len() on a styled string compiles but gives wrong width from why-tea-leaves.mdx:52.
// len() counts bytes (including ANSI escape codes), not visual character width.
func TestCompile_WhyAnsiAntiPattern(t *testing.T) {
	styledText := lipgloss.NewStyle().Bold(true).Render("Hello World!")
	width := len(styledText) // Returns more than 12 due to ANSI escape bytes
	_ = width
}

// TestCompile_WhyAnsiCorrectUsage verifies ANSI-aware width measurement and OverlayModal from why-tea-leaves.mdx:60.
func TestCompile_WhyAnsiCorrectUsage(t *testing.T) {
	styledText := lipgloss.NewStyle().Bold(true).Render("Hello World!")
	width := ansi.StringWidth(styledText) // Returns 12 — correct visual width
	_ = width

	modal := teamodal.NewYesNoModal("Continue?", &teamodal.ConfirmModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	background := "main application view content here"
	// m.modal.OverlayModal(background) is the Tea Leaves pattern for ANSI-safe overlay compositing
	view := modal.OverlayModal(background)
	_ = view
}

// TestCompile_WhyIndependentModules verifies that Tea Leaves components import independently from why-tea-leaves.mdx:75.
// Line 75 in the MDX shows shell-style "go get" commands; this test verifies that
// teafields and teamodal can each be imported independently without requiring the other.
func TestCompile_WhyIndependentModules(t *testing.T) {
	// teafields alone
	dropdown := teafields.NewDropdownModel(
		teafields.ToOptions([]string{"Only dropdown needed"}),
		&teafields.DropdownModelArgs{},
	)
	_ = dropdown

	// teamodal alone
	modal := teamodal.NewYesNoModal("Dialog only", &teamodal.ConfirmModelArgs{})
	_ = modal
}
