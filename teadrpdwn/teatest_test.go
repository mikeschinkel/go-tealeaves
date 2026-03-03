//go:build ignore
// Disabled: teatest (charmbracelet/x/exp/teatest) has no v2 equivalent yet.
// Re-enable when charm.land ships a v2-compatible teatest package.

package teadrpdwn

import (
	"io"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// testProgram wraps DropdownModel into a standalone tea.Model for teatest.
// Pre-opens the dropdown and quits when a selection or cancellation occurs.
type testProgram struct {
	dropdown DropdownModel
	done     bool
}

func newTestProgram(opts []Option) testProgram {
	m := NewModel(opts, 2, 5, &ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	// Pre-open since Init() can't mutate state (value receiver)
	m, _ = m.Open()
	return testProgram{dropdown: m}
}

func (tp testProgram) Init() tea.Cmd {
	return nil
}

func (tp testProgram) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case OptionSelectedMsg:
		tp.done = true
		return tp, tea.Quit
	case DropdownCancelledMsg:
		tp.done = true
		return tp, tea.Quit
	}

	result, cmd := tp.dropdown.Update(msg)
	tp.dropdown = result.(DropdownModel)

	return tp, cmd
}

func (tp testProgram) View() tea.View {
	return tp.dropdown.View()
}

func TestDropdownModel_FullLifecycle(t *testing.T) {
	opts := []Option{
		{Text: "First", Value: 1},
		{Text: "Second", Value: 2},
		{Text: "Third", Value: 3},
	}
	p := newTestProgram(opts)

	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	// Wait for initial render
	time.Sleep(200 * time.Millisecond)

	// Navigate down to "Second"
	tm.Send(tea.KeyPressMsg{Code: tea.KeyDown})
	time.Sleep(100 * time.Millisecond)

	// Select with Enter
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})

	// Wait for program to finish
	fm := tm.FinalModel(t, teatest.WithFinalTimeout(3*time.Second))
	tp := fm.(testProgram)

	if !tp.done {
		t.Error("expected program to be done after selection")
	}
}

func TestDropdownModel_RenderGolden(t *testing.T) {
	opts := []Option{
		{Text: "Apple", Value: "apple"},
		{Text: "Banana", Value: "banana"},
		{Text: "Cherry", Value: "cherry"},
	}
	p := newTestProgram(opts)

	tm := teatest.NewTestModel(t, p, teatest.WithInitialTermSize(80, 24))

	// Wait for initial render to stabilize
	time.Sleep(300 * time.Millisecond)

	// Capture rendered output
	out, err := io.ReadAll(tm.Output())
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}

	// Cancel to quit the program
	tm.Send(tea.KeyPressMsg{Code: tea.KeyEsc})
	tm.WaitFinished(t, teatest.WithFinalTimeout(3*time.Second))

	teatest.RequireEqualOutput(t, out)
}
