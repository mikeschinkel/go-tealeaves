//go:build ignore
// Disabled: teatest (charmbracelet/x/exp/teatest) has no v2 equivalent yet.
// Re-enable when charm.land ships a v2-compatible teatest package.

package teamodal

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// --- OK Modal wrapper ---

type okModalProgram struct {
	modal ModalModel
	done  bool
}

func newOKModalProgram() okModalProgram {
	m := NewOKModal("Test alert message", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m, _ = m.Open()
	return okModalProgram{modal: m}
}

func (p okModalProgram) Init() tea.Cmd { return nil }

func (p okModalProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ClosedMsg:
		p.done = true
		return p, tea.Quit
	}

	result, cmd := p.modal.Update(msg)
	p.modal = result.(ModalModel)
	return p, cmd
}

func (p okModalProgram) View() tea.View {
	return p.modal.View()
}

// --- YesNo Modal wrapper ---

type yesNoModalProgram struct {
	modal ModalModel
	done  bool
}

func newYesNoModalProgram() yesNoModalProgram {
	m := NewYesNoModal("Are you sure?", &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		DefaultYes:   true,
	})
	m, _ = m.Open()
	return yesNoModalProgram{modal: m}
}

func (p yesNoModalProgram) Init() tea.Cmd { return nil }

func (p yesNoModalProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case AnsweredYesMsg, AnsweredNoMsg:
		p.done = true
		return p, tea.Quit
	}

	result, cmd := p.modal.Update(msg)
	p.modal = result.(ModalModel)
	return p, cmd
}

func (p yesNoModalProgram) View() tea.View {
	return p.modal.View()
}

// --- Choice Modal wrapper ---

type choiceModalProgram struct {
	modal ChoiceModel
	done  bool
}

func newChoiceModalProgram() choiceModalProgram {
	m := NewChoiceModel(&ChoiceModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Action",
		Message:      "What would you like to do?",
		Options: []ChoiceOption{
			{Label: "Save", ID: "save", Hotkey: 's'},
			{Label: "Discard", ID: "discard", Hotkey: 'd'},
			{Label: "Cancel", ID: "cancel", Hotkey: 'c'},
		},
	})
	m, _ = m.Open()
	return choiceModalProgram{modal: m}
}

func (p choiceModalProgram) Init() tea.Cmd { return nil }

func (p choiceModalProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ChoiceSelectedMsg, ChoiceCancelledMsg:
		p.done = true
		return p, tea.Quit
	}

	result, cmd := p.modal.Update(msg)
	p.modal = result.(ChoiceModel)
	return p, cmd
}

func (p choiceModalProgram) View() tea.View {
	return p.modal.View()
}

// --- List Modal wrapper ---

type listModalProgram struct {
	list ListModel[testItem]
	done bool
}

func newListModalProgram() listModalProgram {
	items := []testItem{
		{id: "1", label: "Alpha", active: false},
		{id: "2", label: "Beta", active: true},
		{id: "3", label: "Gamma", active: false},
	}
	m := NewListModel(items, &ListModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Select Item",
		MaxVisible:   8,
	})
	m = m.Open()
	return listModalProgram{list: m}
}

func (p listModalProgram) Init() tea.Cmd { return nil }

func (p listModalProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case ListAcceptedMsg[testItem], ListCancelledMsg:
		p.done = true
		return p, tea.Quit
	}

	var cmd tea.Cmd
	p.list, cmd = p.list.Update(msg)
	return p, cmd
}

func (p listModalProgram) View() tea.View {
	return p.list.View()
}

// --- Layer 3 Golden Tests ---

func TestOKModal_GoldenRender(t *testing.T) {
	p := newOKModalProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Close with Enter to quit
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}

func TestYesNoModal_GoldenRender(t *testing.T) {
	p := newYesNoModalProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Cancel with Esc to quit
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}

func TestChoiceModel_GoldenRender(t *testing.T) {
	p := newChoiceModalProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Cancel with Esc to quit
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}

func TestListModel_GoldenRender(t *testing.T) {
	p := newListModalProgram()
	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	time.Sleep(300 * time.Millisecond)

	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Cancel with Esc to quit
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}
