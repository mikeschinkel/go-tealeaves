// Source: site/src/content/docs/guides/composition.mdx:15,43,68,85
package examples_test

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	dt "github.com/mikeschinkel/go-dt"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teanotify"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
	"github.com/mikeschinkel/go-tealeaves/teatree"
)

// compModel holds the full set of component fields used across composition examples.
type compModel struct {
	modal        teamodal.ConfirmModel
	dropdown     teafields.DropdownModel
	statusBar    teastatus.StatusBarModel
	notification teanotify.NotifyModel
	tree         teatree.TreeModel[teatree.File]
	width        int
	height       int
	result       string
}

func newCompModel() compModel {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
		Position: teanotify.TopRightPosition,
	})

	nodes := teatree.BuildFileTree([]*teatree.File{
		teatree.NewFile(dt.RelFilepath("main.go"), nil),
	}, teatree.BuildFileTreeArgs{
		RootPath: dt.PathSegment("project"),
	})
	provider := teatree.NewCompactNodeProvider[teatree.File](teatree.TriangleExpanderControls)
	tree := teatree.NewTree(nodes, &teatree.TreeArgs[teatree.File]{
		NodeProvider:     provider,
		ExpanderControls: &teatree.TriangleExpanderControls,
	})
	treeModel := teatree.NewTreeModel(tree, 20)

	return compModel{
		modal: teamodal.NewYesNoModal("Are you sure?", &teamodal.ConfirmModelArgs{
			ScreenWidth:  80,
			ScreenHeight: 24,
		}),
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"Option A", "Option B"}),
			&teafields.DropdownModelArgs{},
		),
		statusBar:    teastatus.NewStatusBarModel().SetSize(80),
		notification: notify,
		tree:         treeModel,
		width:        80,
		height:       24,
	}
}

func (m compModel) Init() tea.Cmd { return nil }

// Update implements the modal-consumption pattern from composition.mdx:15.
func (m compModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	modal, cmd := m.modal.Update(msg)
	if cmd != nil {
		m.modal = modal.(teamodal.ConfirmModel)
		return m, cmd
	}
	dropdown, cmd := m.dropdown.Update(msg)
	if cmd != nil {
		m.dropdown = dropdown.(teafields.DropdownModel)
		return m, cmd
	}
	// Then main content
	return m, nil
}

// View implements overlay composition from composition.mdx:43.
func (m compModel) View() tea.View {
	var b strings.Builder
	b.WriteString("main content")
	view := b.String()

	if m.dropdown.IsOpen {
		view = teafields.OverlayDropdown(view, m.dropdown.View().Content, 0, 0)
	}
	if m.modal.IsOpen() {
		view = m.modal.OverlayModal(view)
	}
	return tea.NewView(view)
}

// TestCompile_CompositionUpdate verifies the modal-consumption Update pattern from composition.mdx:15.
func TestCompile_CompositionUpdate(t *testing.T) {
	m := newCompModel()
	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}

// TestCompile_CompositionView verifies the overlay View pattern from composition.mdx:43.
func TestCompile_CompositionView(t *testing.T) {
	m := newCompModel()
	view := m.View()
	_ = view
}

// TestCompile_CompositionWindowSize verifies WindowSizeMsg handling from composition.mdx:68.
func TestCompile_CompositionWindowSize(t *testing.T) {
	m := newCompModel()
	msg := tea.WindowSizeMsg{Width: 100, Height: 30}

	m.width = msg.Width
	m.height = msg.Height
	m.modal = m.modal.SetSize(msg.Width, msg.Height)
	m.dropdown = m.dropdown.WithScreenSize(msg.Width, msg.Height)
	m.tree = m.tree.SetSize(msg.Width, msg.Height-2)
	m.statusBar = m.statusBar.SetSize(msg.Width)

	_ = m
}

// fullUpdateModel is a dedicated model for the full Update example from composition.mdx:85.
type fullUpdateModel struct {
	modal        teamodal.ConfirmModel
	notification teanotify.NotifyModel
	result       string
	width        int
	height       int
}

func (m fullUpdateModel) Init() tea.Cmd { return nil }

// Update demonstrates the full Update with teastatus + teamodal + teanotify from composition.mdx:85.
// Note: teanotify.NotifyModel.Update returns (NotifyModel, tea.Cmd) directly — it is not a tea.Model.
func (m fullUpdateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Always let notification model process messages (tick, dismiss, etc.)
	var notifyCmd tea.Cmd
	m.notification, notifyCmd = m.notification.Update(msg)

	modal, cmd := m.modal.Update(msg)
	if cmd != nil {
		m.modal = modal.(teamodal.ConfirmModel)
		return m, tea.Batch(notifyCmd, cmd)
	}

	switch msg := msg.(type) {
	case teamodal.AnsweredYesMsg:
		m.result = "Confirmed"
		return m, notifyCmd
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, notifyCmd
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Batch(notifyCmd, tea.Quit)
		}
	}

	return m, notifyCmd
}

func (m fullUpdateModel) View() tea.View {
	return tea.NewView("full update example")
}

// TestCompile_CompositionFullUpdate verifies the full Update pattern from composition.mdx:85.
func TestCompile_CompositionFullUpdate(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
	})
	m := fullUpdateModel{
		modal: teamodal.NewYesNoModal("Continue?", &teamodal.ConfirmModelArgs{
			ScreenWidth:  80,
			ScreenHeight: 24,
		}),
		notification: notify,
	}

	result, cmd := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	_ = result
	_ = cmd
}
