// Source: site/src/content/docs/guides/modal-consumption.mdx:17,47,76,128
package examples_test

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// ----- Anti-pattern model (composition.mdx:17) -----

// antiPatternModel demonstrates the infected state-check anti-pattern.
type antiPatternModel struct {
	dropdown teafields.DropdownModel
}

func (m antiPatternModel) Init() tea.Cmd { return nil }

// Update shows the anti-pattern: checking modal state before handling input.
func (m antiPatternModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.dropdown.IsOpen { // State check "infection"
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m antiPatternModel) View() tea.View {
	return tea.NewView("anti-pattern example")
}

// TestCompile_AntiPatternStateCheck verifies the anti-pattern code from modal-consumption.mdx:17.
func TestCompile_AntiPatternStateCheck(t *testing.T) {
	m := antiPatternModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"A", "B"}),
			&teafields.DropdownModelArgs{},
		),
	}
	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}

// ----- Clean consumption pattern (modal-consumption.mdx:47) -----

// cleanConsumptionModel demonstrates the Tea Leaves modal consumption pattern.
type cleanConsumptionModel struct {
	dropdown teafields.DropdownModel
}

func (m cleanConsumptionModel) Init() tea.Cmd { return nil }

// Update implements the Tea Leaves consumption pattern: let the modal go first.
func (m cleanConsumptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	dropdown, cmd := m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teafields.DropdownModel)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "space":
			m.dropdown, cmd = m.dropdown.Open()
			return m, cmd
		}
	}
	return m, nil
}

func (m cleanConsumptionModel) View() tea.View {
	return tea.NewView("clean consumption example")
}

// TestCompile_ConsumptionPatternClean verifies the clean consumption pattern from modal-consumption.mdx:47.
func TestCompile_ConsumptionPatternClean(t *testing.T) {
	m := cleanConsumptionModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"X", "Y"}),
			&teafields.DropdownModelArgs{},
		),
	}
	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}

// ----- Dropdown internals / goto pattern (modal-consumption.mdx:76) -----

// illustrativeDropdown illustrates the internal goto pattern used inside a
// component's Update — NOT the real teafields.DropdownModel API.

type optionSelectedMsg struct{ Index int }
type dropdownCancelledMsg struct{}

type illustrativeDropdown struct {
	IsOpen   bool
	Selected int
}

func (m illustrativeDropdown) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if !m.IsOpen {
		goto end
	}
	{
		keyMsg, ok := msg.(tea.KeyPressMsg)
		if !ok {
			goto end
		}
		switch keyMsg.String() {
		case "up", "k":
			m.Selected--
			cmd = func() tea.Msg { return nil }
			goto end
		case "enter":
			m.IsOpen = false
			selected := m.Selected
			cmd = func() tea.Msg { return optionSelectedMsg{Index: selected} }
			goto end
		case "esc":
			m.IsOpen = false
			cmd = func() tea.Msg { return dropdownCancelledMsg{} }
			goto end
		}
	}
end:
	return m, cmd
}

func (m illustrativeDropdown) Init() tea.Cmd { return nil }

func (m illustrativeDropdown) View() tea.View {
	return tea.NewView("")
}

// TestCompile_DropdownGotoPattern verifies the illustrative goto pattern from modal-consumption.mdx:76.
func TestCompile_DropdownGotoPattern(t *testing.T) {
	m := illustrativeDropdown{IsOpen: true, Selected: 1}
	result, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	_ = result
	_ = cmd
}

// ----- Multiple modals pattern (modal-consumption.mdx:128) -----

// modalConsumptionModel demonstrates composing multiple modal components.
type modalConsumptionModel struct {
	dropdown    teafields.DropdownModel
	confirmModal teamodal.ConfirmModel
}

func (m modalConsumptionModel) Init() tea.Cmd { return nil }

// Update applies the consumption pattern for two separate modals in sequence.
func (m modalConsumptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	dropdown, cmd := m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teafields.DropdownModel)
		return m, cmd
	}
	modal, cmd := m.confirmModal.Update(msg)
	if cmd != nil {
		m.confirmModal = modal.(teamodal.ConfirmModel)
		return m, cmd
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		_ = msg
	}
	return m, nil
}

func (m modalConsumptionModel) View() tea.View {
	return tea.NewView("multiple modals example")
}

// TestCompile_MultipleModalsPattern verifies multiple modals composition from modal-consumption.mdx:128.
func TestCompile_MultipleModalsPattern(t *testing.T) {
	m := modalConsumptionModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"Choice 1", "Choice 2"}),
			&teafields.DropdownModelArgs{},
		),
		confirmModal: teamodal.NewYesNoModal("Delete item?", &teamodal.ConfirmModelArgs{
			ScreenWidth:  80,
			ScreenHeight: 24,
		}),
	}
	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}
