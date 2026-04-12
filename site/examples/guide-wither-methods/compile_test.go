// Source: site/src/content/docs/guides/wither-methods.mdx:15,30,58,72,82,94
package examples_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
	"github.com/mikeschinkel/go-tealeaves/teanotify"
)

// witherMyModel is a local stand-in for the MyModel type referenced in wither-methods.mdx.
type witherMyModel struct {
	contentWidth int
	totalWidth   int
}

// calculateLayout demonstrates the wither method pattern from wither-methods.mdx:15.
func (m witherMyModel) calculateLayout() witherMyModel {
	m.contentWidth = 80 - 4  // example calculation
	m.totalWidth = 80
	return m
}

// TestCompile_WitherMethodPattern verifies the correct wither pattern from wither-methods.mdx:15.
func TestCompile_WitherMethodPattern(t *testing.T) {
	m := witherMyModel{}
	m = m.calculateLayout()
	_ = m.contentWidth
	_ = m.totalWidth
}

// calculateLayoutMultiReturn demonstrates the multiple-return comparison from wither-methods.mdx:30.
func (m witherMyModel) calculateLayoutMultiReturn() (contentWidth, totalWidth int) {
	return 80 - 4, 80 // example calculation
}

// TestCompile_MultipleReturnComparison verifies the verbose multi-return alternative from wither-methods.mdx:30.
func TestCompile_MultipleReturnComparison(t *testing.T) {
	m := witherMyModel{}
	contentWidth, totalWidth := m.calculateLayoutMultiReturn()
	_ = contentWidth
	_ = totalWidth
}

// TestCompile_TeamodalWitherPattern verifies the teamodal SetSize wither and NewYesNoModal from wither-methods.mdx:58.
func TestCompile_TeamodalWitherPattern(t *testing.T) {
	type witherParentModel struct {
		confirmDialog teamodal.ConfirmModel
	}
	pm := witherParentModel{
		confirmDialog: teamodal.NewYesNoModal("Proceed?", &teamodal.ConfirmModelArgs{
			ScreenWidth:  80,
			ScreenHeight: 24,
		}),
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	pm.confirmDialog = pm.confirmDialog.SetSize(msg.Width, msg.Height)

	modal := teamodal.NewYesNoModal("Proceed?", &teamodal.ConfirmModelArgs{
		Title:      "Confirmation",
		DefaultYes: true,
	})
	_ = modal
}

// TestCompile_TeafieldsDropdownWither verifies the teafields dropdown wither chain from wither-methods.mdx:72.
func TestCompile_TeafieldsDropdownWither(t *testing.T) {
	type dropdownParentModel struct {
		dropdown teafields.DropdownModel
	}
	pm := dropdownParentModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"A", "B", "C"}),
			&teafields.DropdownModelArgs{},
		),
	}
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}
	row, col := 5, 10
	pm.dropdown = pm.dropdown.
		WithScreenSize(msg.Width, msg.Height).
		WithPosition(row, col)
	_ = pm
}

// TestCompile_TeanotifyPositionWither verifies the teanotify WithPosition wither from wither-methods.mdx:82.
func TestCompile_TeanotifyPositionWither(t *testing.T) {
	type notifyParentModel struct {
		notify teanotify.NotifyModel
	}
	pm := notifyParentModel{
		notify: teanotify.NewNotifyModel(teanotify.NotifyOpts{
			Width:    40,
			Duration: 3e9, // 3 seconds
		}),
	}
	pm.notify = pm.notify.WithPosition(teanotify.TopRightPosition)
	_ = pm
}

// witherValueReceiverModel demonstrates the correct value-receiver pattern from wither-methods.mdx:94.
// The MDX shows a pointer-receiver anti-pattern that conflicts with Bubble Tea; this test
// shows the correct value-receiver equivalent that compiles and works correctly.
type witherValueReceiverModel struct {
	contentWidth int
}

func (m witherValueReceiverModel) calculateLayoutCorrect() witherValueReceiverModel {
	m.contentWidth = 76 // value receiver: works correctly — returns updated copy
	return m
}

func (m witherValueReceiverModel) Init() tea.Cmd { return nil }

func (m witherValueReceiverModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m = m.calculateLayoutCorrect() // value receiver: compiles and works
	return m, nil
}

func (m witherValueReceiverModel) View() tea.View {
	return tea.NewView("")
}

// TestCompile_ValueReceiverPattern verifies the correct value-receiver pattern from wither-methods.mdx:94.
// (The MDX illustrates a pointer-receiver anti-pattern; this test shows the compilable correct version.)
func TestCompile_ValueReceiverPattern(t *testing.T) {
	m := witherValueReceiverModel{}
	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}
