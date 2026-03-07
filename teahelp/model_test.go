package teahelp

import (
	"strings"
	"testing"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func testKeys() map[string][]teautils.KeyMeta {
	return map[string][]teautils.KeyMeta{
		"Navigation": {
			{
				Binding: key.NewBinding(
					key.WithKeys("up", "k"),
					key.WithHelp("up/k", "move up"),
				),
				HelpText: "Move cursor up",
			},
			{
				Binding: key.NewBinding(
					key.WithKeys("down", "j"),
					key.WithHelp("down/j", "move down"),
				),
				HelpText: "Move cursor down",
			},
		},
		"Actions": {
			{
				Binding: key.NewBinding(
					key.WithKeys("enter"),
					key.WithHelp("enter", "confirm"),
				),
				HelpText: "Confirm selection",
			},
		},
	}
}

func TestNewHelpVisorModel(t *testing.T) {
	m := NewHelpVisorModel()
	if m.IsOpen() {
		t.Error("expected new model to be closed")
	}
	if m.Page() != 0 {
		t.Errorf("expected page 0, got %d", m.Page())
	}
	if m.TotalPages() != 0 {
		t.Errorf("expected 0 total pages, got %d", m.TotalPages())
	}
}

func TestOpenClose(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 40)

	// Open
	m = m.Open(testKeys())
	if !m.IsOpen() {
		t.Error("expected model to be open after Open()")
	}
	if m.TotalPages() == 0 {
		t.Error("expected at least 1 page after Open()")
	}

	// Close
	m = m.Close()
	if m.IsOpen() {
		t.Error("expected model to be closed after Close()")
	}
	if m.Page() != 0 {
		t.Errorf("expected page reset to 0, got %d", m.Page())
	}
}

func TestSetSizeRebuildsPages(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 40)
	m = m.Open(testKeys())

	pagesAt40 := m.TotalPages()

	// Shrink height to force more pages
	m = m.SetSize(80, 10)
	pagesAt10 := m.TotalPages()

	if pagesAt10 <= pagesAt40 {
		t.Errorf("expected more pages at height 10 (%d) than height 40 (%d)",
			pagesAt10, pagesAt40)
	}
}

func TestViewWhenClosed(t *testing.T) {
	m := NewHelpVisorModel()
	view := m.View()
	if view.Content != "" {
		t.Errorf("expected empty view when closed, got %q", view.Content)
	}
}

func TestViewWhenOpen(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 40)
	m = m.Open(testKeys())

	view := m.View()
	if view.Content == "" {
		t.Error("expected non-empty view when open")
	}
	if !strings.Contains(view.Content, "Keyboard Shortcuts") {
		t.Error("expected view to contain 'Keyboard Shortcuts' title")
	}
}

func TestUpdateCloseEmitsClosedMsg(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 40)
	m = m.Open(testKeys())

	// Press Esc to close
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if m.IsOpen() {
		t.Error("expected model to be closed after Esc")
	}
	if cmd == nil {
		t.Fatal("expected cmd to emit ClosedMsg")
	}

	msg := cmd()
	if _, ok := msg.(ClosedMsg); !ok {
		t.Errorf("expected ClosedMsg, got %T", msg)
	}
}

func TestUpdateCloseWithQuestionMark(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 40)
	m = m.Open(testKeys())

	// Press ? to close
	m, cmd := m.Update(tea.KeyPressMsg{Code: '?', Text: "?"})
	if m.IsOpen() {
		t.Error("expected model to be closed after ?")
	}
	if cmd == nil {
		t.Fatal("expected cmd to emit ClosedMsg")
	}
}

func TestUpdatePagination(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 10) // small height to force multiple pages
	m = m.Open(testKeys())

	if m.TotalPages() < 2 {
		t.Skip("not enough content to test pagination at this height")
	}

	// Navigate right
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyRight})
	if m.Page() != 1 {
		t.Errorf("expected page 1 after right, got %d", m.Page())
	}

	// Navigate left
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.Page() != 0 {
		t.Errorf("expected page 0 after left, got %d", m.Page())
	}

	// Left at page 0 stays at 0
	m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyLeft})
	if m.Page() != 0 {
		t.Errorf("expected page 0 when already at first page, got %d", m.Page())
	}
}

func TestUpdateDigitPageNavigation(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 10) // small height to force multiple pages
	m = m.Open(testKeys())

	if m.TotalPages() < 2 {
		t.Skip("not enough content to test digit navigation")
	}

	// Press "2" to go to page 2
	m, _ = m.Update(tea.KeyPressMsg{Code: '2', Text: "2"})
	if m.Page() != 1 {
		t.Errorf("expected page 1 (0-based) after pressing '2', got %d", m.Page())
	}

	// Press "1" to go back to page 1
	m, _ = m.Update(tea.KeyPressMsg{Code: '1', Text: "1"})
	if m.Page() != 0 {
		t.Errorf("expected page 0 after pressing '1', got %d", m.Page())
	}
}

func TestUpdateIgnoredWhenClosed(t *testing.T) {
	m := NewHelpVisorModel()
	m, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEsc})
	if cmd != nil {
		t.Error("expected nil cmd when model is closed")
	}
	if m.IsOpen() {
		t.Error("expected model to remain closed")
	}
}

func TestWithContentStyle(t *testing.T) {
	m := NewHelpVisorModel()
	custom := teautils.DefaultHelpVisorStyle()
	custom.KeyColumnGap = 8
	m2 := m.WithContentStyle(custom)
	if m2.contentStyle.KeyColumnGap != 8 {
		t.Errorf("expected KeyColumnGap 8, got %d", m2.contentStyle.KeyColumnGap)
	}
}

func TestWithKeys(t *testing.T) {
	m := NewHelpVisorModel()
	customKeys := DefaultHelpVisorKeyMap()
	customKeys.Close = key.NewBinding(key.WithKeys("q"))
	m = m.WithKeys(customKeys)
	m = m.SetSize(80, 40)
	m = m.Open(testKeys())

	// Press q to close (custom binding)
	m, cmd := m.Update(tea.KeyPressMsg{Code: 'q', Text: "q"})
	if m.IsOpen() {
		t.Error("expected model to be closed after custom close key 'q'")
	}
	if cmd == nil {
		t.Fatal("expected cmd to emit ClosedMsg")
	}
}

func TestFooterAppearsOnMultiPage(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 10) // small height to force multiple pages
	m = m.Open(testKeys())

	if m.TotalPages() < 2 {
		t.Skip("not enough content for multi-page footer test")
	}

	view := m.View()
	if !strings.Contains(view.Content, "Page 1/") {
		t.Error("expected footer with page indicator on multi-page visor")
	}
}

func TestFooterAbsentOnSinglePage(t *testing.T) {
	m := NewHelpVisorModel()
	m = m.SetSize(80, 80) // large height for single page
	m = m.Open(testKeys())

	if m.TotalPages() != 1 {
		t.Skip("expected single page with large height")
	}

	view := m.View()
	if strings.Contains(view.Content, "Page ") {
		t.Error("expected no footer on single-page visor")
	}
}

func TestInitReturnsNil(t *testing.T) {
	m := NewHelpVisorModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init() to return nil")
	}
}
