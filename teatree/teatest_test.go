//go:build ignore
// Disabled: teatest (charmbracelet/x/exp/teatest) has no v2 equivalent yet.
// Re-enable when charm.land ships a v2-compatible teatest package.

package teatree

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// treeProgram wraps TreeModel[string] into a standalone tea.Model for teatest.
type treeProgram struct {
	model TreeModel[string]
}

func newTreeProgram() treeProgram {
	tree, _ := buildTestTree()
	m := NewModel(tree, 10)
	m = m.SetSize(80, 24)
	return treeProgram{model: m}
}

func (p treeProgram) Init() tea.Cmd { return nil }

func (p treeProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.Code == 'c' && msg.Mod.Contains(tea.ModCtrl) {
			return p, tea.Quit
		}
	}

	var cmd tea.Cmd
	p.model, cmd = p.model.Update(msg)
	return p, cmd
}

func (p treeProgram) View() tea.View {
	return p.model.View()
}

// --- Layer 3 Golden Tests ---

func TestTreeModel_NavigationGolden(t *testing.T) {
	p := newTreeProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	// Expand root1
	tm.Send(tea.KeyPressMsg{Code: tea.KeyRight})
	time.Sleep(100 * time.Millisecond)

	// Navigate down to child1
	tm.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	time.Sleep(100 * time.Millisecond)

	// Expand child1
	tm.Send(tea.KeyPressMsg{Code: tea.KeyRight})
	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Quit the program
	tm.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}

func TestTreeModel_ExpandCollapseGolden(t *testing.T) {
	p := newTreeProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	// Expand root1
	tm.Send(tea.KeyPressMsg{Code: tea.KeyRight})
	time.Sleep(100 * time.Millisecond)

	// Collapse root1
	tm.Send(tea.KeyPressMsg{Code: tea.KeyLeft})
	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Quit the program
	tm.Send(tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}
