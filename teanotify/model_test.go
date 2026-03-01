package teanotify

import (
	"errors"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Helpers ---

func extractMsg(cmd tea.Cmd) tea.Msg {
	if cmd == nil {
		return nil
	}
	return cmd()
}

func newTestModel() NotifyModel {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
	})
	m, err := m.Initialize()
	if err != nil {
		panic("newTestModel: " + err.Error())
	}
	return m
}

func newTestModelWithOpts(t *testing.T, opts NotifyOpts) NotifyModel {
	t.Helper()
	m := NewNotifyModel(opts)
	m, err := m.Initialize()
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func activateNotice(t *testing.T, m NotifyModel, key NoticeKey, msg string) (NotifyModel, tea.Cmd) {
	t.Helper()
	nm := notifyMsg{noticeKey: key, msg: msg, dur: m.duration}
	out, cmd := m.Update(nm)
	return out, cmd
}

func makeContent(width, height int) string {
	row := strings.Repeat(".", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return strings.Join(rows, "\n")
}

// --- Constructor tests ---

func TestNewNotifyModel_Defaults(t *testing.T) {
	m := newTestModel()
	if m.position != TopLeftPosition {
		t.Errorf("expected default position TopLeftPosition, got %q", m.position)
	}
	if len(m.noticeTypes) != 4 {
		t.Errorf("expected 4 notice types, got %d", len(m.noticeTypes))
	}
	if m.activeNotice != nil {
		t.Error("expected activeNotice to be nil")
	}
}

func TestNewNotifyModel_ExplicitPosition(t *testing.T) {
	m := newTestModelWithOpts(t, NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
		Position: BottomRightPosition,
	})
	if m.position != BottomRightPosition {
		t.Errorf("expected BottomRightPosition, got %q", m.position)
	}
}

func TestInitialize_InvalidWidth(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    0,
		Duration: 3 * time.Second,
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidWidth) {
		t.Fatalf("expected ErrInvalidWidth, got: %v", err)
	}
}

func TestInitialize_NegativeWidth(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    -5,
		Duration: 3 * time.Second,
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidWidth) {
		t.Fatalf("expected ErrInvalidWidth, got: %v", err)
	}
}

func TestInitialize_InvalidDuration(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: 0,
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidDuration) {
		t.Fatalf("expected ErrInvalidDuration, got: %v", err)
	}
}

func TestInitialize_NegativeDuration(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: -1 * time.Second,
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidDuration) {
		t.Fatalf("expected ErrInvalidDuration, got: %v", err)
	}
}

func TestNewNotifyModel_MinWidthClampedToWidth(t *testing.T) {
	m := newTestModelWithOpts(t, NotifyOpts{
		Width:    30,
		MinWidth: 50,
		Duration: 3 * time.Second,
	})
	if m.minWidth != 30 {
		t.Errorf("expected minWidth clamped to 30, got %d", m.minWidth)
	}
}

func TestNewNotifyModel_NoDefaultNotices(t *testing.T) {
	m := newTestModelWithOpts(t, NotifyOpts{
		Width:            40,
		Duration:         3 * time.Second,
		NoDefaultNotices: true,
	})
	if len(m.noticeTypes) != 0 {
		t.Errorf("expected 0 notice types, got %d", len(m.noticeTypes))
	}
}

func TestNewNotifyModel_CustomNotices(t *testing.T) {
	m := newTestModelWithOpts(t, NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
		CustomNotices: []NoticeDefinition{
			{Key: "Custom", ForeColor: "#AABBCC", Prefix: "[C]"},
		},
	})
	if len(m.noticeTypes) != 5 {
		t.Errorf("expected 5 notice types (4 defaults + 1 custom), got %d", len(m.noticeTypes))
	}
}

func TestInitialize_CustomNoticeInvalidColor(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
		CustomNotices: []NoticeDefinition{
			{Key: "Bad", ForeColor: "bad"},
		},
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidColor) {
		t.Fatalf("expected ErrInvalidColor, got: %v", err)
	}
}

func TestInitialize_CustomNoticeEmptyKey(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
		CustomNotices: []NoticeDefinition{
			{Key: "", ForeColor: "#AABBCC"},
		},
	})
	_, err := m.Initialize()
	if !errors.Is(err, ErrInvalidNoticeKey) {
		t.Fatalf("expected ErrInvalidNoticeKey, got: %v", err)
	}
}

// --- With* method tests ---

func TestWithPosition(t *testing.T) {
	m := newTestModel()
	original := m.position
	out := m.WithPosition(BottomCenterPosition)
	if out.position != BottomCenterPosition {
		t.Errorf("expected BottomCenterPosition, got %q", out.position)
	}
	// Original unchanged (value semantics)
	if m.position != original {
		t.Errorf("expected original position unchanged, got %q", m.position)
	}
}

func TestWithMinWidth(t *testing.T) {
	m := newTestModel()
	out := m.WithMinWidth(20)
	if out.minWidth != 20 {
		t.Errorf("expected minWidth 20, got %d", out.minWidth)
	}
}

func TestWithMinWidth_ClampedToWidth(t *testing.T) {
	m := newTestModel() // width=40
	out := m.WithMinWidth(60)
	if out.minWidth != 40 {
		t.Errorf("expected minWidth clamped to 40, got %d", out.minWidth)
	}
}

func TestWithUnicodePrefix(t *testing.T) {
	m := newTestModel()
	out := m.WithUnicodePrefix()
	// Check that default notice types now have unicode prefixes
	if def, ok := out.noticeTypes[InfoKey]; ok {
		if def.Prefix != InfoUnicodePrefix {
			t.Errorf("expected Unicode prefix for Info, got %q", def.Prefix)
		}
	} else {
		t.Fatal("expected InfoKey in noticeTypes")
	}
	// Original unchanged
	if def, ok := m.noticeTypes[InfoKey]; ok {
		if def.Prefix != InfoASCIIPrefix {
			t.Errorf("expected original to retain ASCII prefix, got %q", def.Prefix)
		}
	}
}

func TestWithAllowEscToClose(t *testing.T) {
	m := newTestModel()
	if m.allowEscToClose {
		t.Fatal("expected allowEscToClose to be false by default")
	}
	out := m.WithAllowEscToClose()
	if !out.allowEscToClose {
		t.Error("expected allowEscToClose to be true after WithAllowEscToClose")
	}
}

// --- Init / View tests ---

func TestInit_ReturnsNilCmd(t *testing.T) {
	m := newTestModel()
	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init to return nil cmd")
	}
}

func TestView_ReturnsEmptyString(t *testing.T) {
	m := newTestModel()
	s := m.View()
	if s != "" {
		t.Errorf("expected empty string from View, got %q", s)
	}
}

// --- Update lifecycle tests ---

func TestUpdate_NotifyMsg_ActivatesNotice(t *testing.T) {
	m := newTestModel()
	out, cmd := activateNotice(t, m, InfoKey, "test message")
	if out.activeNotice == nil {
		t.Fatal("expected activeNotice to be set")
	}
	if cmd == nil {
		t.Fatal("expected cmd to be non-nil")
	}
}

func TestUpdate_NotifyMsg_SetsNoticeFields(t *testing.T) {
	m := newTestModel()
	out, _ := activateNotice(t, m, InfoKey, "test message")
	n := out.activeNotice
	if n == nil {
		t.Fatal("expected activeNotice to be set")
	}
	if n.prefix != InfoASCIIPrefix {
		t.Errorf("expected prefix %q, got %q", InfoASCIIPrefix, n.prefix)
	}
	if n.width != m.width {
		t.Errorf("expected width %d, got %d", m.width, n.width)
	}
	if n.minWidth != m.minWidth {
		t.Errorf("expected minWidth %d, got %d", m.minWidth, n.minWidth)
	}
	if n.position != m.position {
		t.Errorf("expected position %q, got %q", m.position, n.position)
	}
}

func TestUpdate_NotifyMsg_UnregisteredKey(t *testing.T) {
	m := newTestModel()
	out, _ := activateNotice(t, m, "NonExistent", "msg")
	if out.activeNotice != nil {
		t.Error("expected activeNotice to be nil for unregistered key")
	}
}

func TestUpdate_NotifyMsg_EmptyMessage(t *testing.T) {
	m := newTestModel()
	out, _ := activateNotice(t, m, InfoKey, "")
	if out.activeNotice != nil {
		t.Error("expected activeNotice to be nil for empty message")
	}
}

func TestUpdate_TickMsg_AdvancesLerp(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")
	if m.activeNotice == nil {
		t.Fatal("expected activeNotice to be set")
	}
	initialLerp := m.activeNotice.curLerpStep

	// Send a tick before death time
	tick := tickMsg(time.Now())
	out, cmd := m.Update(tick)
	if out.activeNotice == nil {
		t.Fatal("expected activeNotice to still be set")
	}
	expectedLerp := initialLerp + DefaultLerpIncrement
	if expectedLerp > 1 {
		expectedLerp = 1
	}
	if out.activeNotice.curLerpStep != expectedLerp {
		t.Errorf("expected curLerpStep %f, got %f", expectedLerp, out.activeNotice.curLerpStep)
	}
	if cmd == nil {
		t.Error("expected cmd to be non-nil (continued ticking)")
	}
}

func TestUpdate_TickMsg_LerpClampsAtOne(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")
	m.activeNotice.curLerpStep = 0.95

	tick := tickMsg(time.Now())
	out, _ := m.Update(tick)
	if out.activeNotice.curLerpStep != 1.0 {
		t.Errorf("expected curLerpStep clamped at 1.0, got %f", out.activeNotice.curLerpStep)
	}
}

func TestUpdate_TickMsg_ExpiresNotice(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")

	// Tick with time past death
	tick := tickMsg(time.Now().Add(10 * time.Second))
	out, cmd := m.Update(tick)
	if out.activeNotice != nil {
		t.Error("expected activeNotice to be nil after expiry")
	}
	if cmd != nil {
		t.Error("expected cmd to be nil after expiry")
	}
}

func TestUpdate_TickMsg_NoActiveNotice(t *testing.T) {
	m := newTestModel()
	tick := tickMsg(time.Now())
	_, cmd := m.Update(tick)
	if cmd != nil {
		t.Error("expected cmd to be nil with no active notice")
	}
}

func TestUpdate_EscKey_WithAllowEscToClose(t *testing.T) {
	m := newTestModel().WithAllowEscToClose()
	m, _ = activateNotice(t, m, InfoKey, "test")
	if m.activeNotice == nil {
		t.Fatal("expected activeNotice to be set")
	}

	esc := tea.KeyMsg{Type: tea.KeyEscape}
	out, _ := m.Update(esc)
	if out.activeNotice != nil {
		t.Error("expected activeNotice to be cleared by Esc")
	}
}

func TestUpdate_EscKey_WithoutAllowEscToClose(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")
	if m.activeNotice == nil {
		t.Fatal("expected activeNotice to be set")
	}

	esc := tea.KeyMsg{Type: tea.KeyEscape}
	out, _ := m.Update(esc)
	if out.activeNotice == nil {
		t.Error("expected activeNotice to NOT be cleared without allowEscToClose")
	}
}

func TestUpdate_EscKey_NoActiveNotice(t *testing.T) {
	m := newTestModel().WithAllowEscToClose()
	esc := tea.KeyMsg{Type: tea.KeyEscape}
	_, cmd := m.Update(esc)
	if cmd != nil {
		t.Error("expected cmd to be nil with no active notice")
	}
}

func TestUpdate_NonEscKey_NotDismissed(t *testing.T) {
	m := newTestModel().WithAllowEscToClose()
	m, _ = activateNotice(t, m, InfoKey, "test")

	enter := tea.KeyMsg{Type: tea.KeyEnter}
	out, _ := m.Update(enter)
	if out.activeNotice == nil {
		t.Error("expected activeNotice to remain active after non-Esc key")
	}
}

func TestUpdate_UnhandledMsg_TicksIfActive(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")

	wsm := tea.WindowSizeMsg{Width: 80, Height: 24}
	_, cmd := m.Update(wsm)
	if cmd == nil {
		t.Error("expected cmd to be non-nil (keeps ticking while active)")
	}
}

func TestUpdate_UnhandledMsg_NoTickIfInactive(t *testing.T) {
	m := newTestModel()
	wsm := tea.WindowSizeMsg{Width: 80, Height: 24}
	_, cmd := m.Update(wsm)
	if cmd != nil {
		t.Error("expected cmd to be nil with no active notice")
	}
}

// --- Other methods ---

func TestHasActiveNotice_False(t *testing.T) {
	m := newTestModel()
	if m.HasActiveNotice() {
		t.Error("expected HasActiveNotice to be false")
	}
}

func TestHasActiveNotice_True(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "test")
	if !m.HasActiveNotice() {
		t.Error("expected HasActiveNotice to be true after activation")
	}
}

func TestNewNotifyCmd_ProducesNotifyMsg(t *testing.T) {
	m := newTestModel()
	cmd := m.NewNotifyCmd(InfoKey, "hello")
	msg := extractMsg(cmd)
	nm, ok := msg.(notifyMsg)
	if !ok {
		t.Fatalf("expected notifyMsg, got %T", msg)
	}
	if nm.noticeKey != InfoKey {
		t.Errorf("expected key %q, got %q", InfoKey, nm.noticeKey)
	}
	if nm.msg != "hello" {
		t.Errorf("expected msg %q, got %q", "hello", nm.msg)
	}
	if nm.dur != m.duration {
		t.Errorf("expected duration %v, got %v", m.duration, nm.dur)
	}
}

func TestRegisterNoticeType_ImmutableSemantics(t *testing.T) {
	m := newTestModel()
	originalCount := len(m.noticeTypes)

	out, err := m.RegisterNoticeType(NoticeDefinition{
		Key:       "Extra",
		ForeColor: "#112233",
		Prefix:    "[E]",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.noticeTypes) != originalCount+1 {
		t.Errorf("expected %d notice types, got %d", originalCount+1, len(out.noticeTypes))
	}
	if len(m.noticeTypes) != originalCount {
		t.Errorf("expected original to have %d notice types, got %d", originalCount, len(m.noticeTypes))
	}
}

func TestRegisterNoticeType_InvalidColor(t *testing.T) {
	m := newTestModel()
	_, err := m.RegisterNoticeType(NoticeDefinition{
		Key:       "Bad",
		ForeColor: "nope",
	})
	if !errors.Is(err, ErrInvalidColor) {
		t.Fatalf("expected ErrInvalidColor, got: %v", err)
	}
}

func TestNewNotice_ValidKey(t *testing.T) {
	m := newTestModel()
	n := m.newNotice(InfoKey, "hello", 3*time.Second)
	if n == nil {
		t.Fatal("expected non-nil notice")
	}
	if n.message != "hello" {
		t.Errorf("expected message 'hello', got %q", n.message)
	}
}

func TestNewNotice_EmptyKey(t *testing.T) {
	m := newTestModel()
	n := m.newNotice("", "hello", 3*time.Second)
	if n != nil {
		t.Fatal("expected nil notice for empty key")
	}
}

func TestNewNotice_EmptyMsg(t *testing.T) {
	m := newTestModel()
	n := m.newNotice(InfoKey, "", 3*time.Second)
	if n != nil {
		t.Fatal("expected nil notice for empty message")
	}
}

func TestNewNotice_UnregisteredKey(t *testing.T) {
	m := newTestModel()
	n := m.newNotice("DoesNotExist", "hello", 3*time.Second)
	if n != nil {
		t.Fatal("expected nil notice for unregistered key")
	}
}

// --- Render overlay tests (Layer 2) ---

func TestRender_NoActiveNotice(t *testing.T) {
	m := newTestModel()
	content := makeContent(40, 10)
	result := m.Render(content)
	if result != content {
		t.Error("expected content to be unchanged when no active notice")
	}
}

func TestRender_TopLeft(t *testing.T) {
	m := newTestModel().WithPosition(TopLeftPosition)
	m, _ = activateNotice(t, m, InfoKey, "top-left notice")

	content := makeContent(60, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	// Top lines should contain the notice message
	topPortion := strings.Join(lines[:3], "\n")
	if !strings.Contains(topPortion, "top-left notice") {
		t.Error("expected top lines to contain 'top-left notice'")
	}
	// Lower lines should be unchanged
	if lines[9] != strings.Repeat(".", 60) {
		t.Errorf("expected last line unchanged, got %q", lines[9])
	}
}

func TestRender_TopCenter(t *testing.T) {
	m := newTestModel().WithPosition(TopCenterPosition)
	m, _ = activateNotice(t, m, InfoKey, "centered notice")

	content := makeContent(80, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	topPortion := strings.Join(lines[:3], "\n")
	if !strings.Contains(topPortion, "centered notice") {
		t.Error("expected top lines to contain 'centered notice'")
	}
}

func TestRender_TopRight(t *testing.T) {
	m := newTestModel().WithPosition(TopRightPosition)
	m, _ = activateNotice(t, m, InfoKey, "right notice")

	content := makeContent(80, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	topPortion := strings.Join(lines[:3], "\n")
	if !strings.Contains(topPortion, "right notice") {
		t.Error("expected top lines to contain 'right notice'")
	}
}

func TestRender_BottomLeft(t *testing.T) {
	m := newTestModel().WithPosition(BottomLeftPosition)
	m, _ = activateNotice(t, m, InfoKey, "bottom-left notice")

	content := makeContent(60, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	// Top lines should be unchanged
	if lines[0] != strings.Repeat(".", 60) {
		t.Errorf("expected first line unchanged, got %q", lines[0])
	}
	// Bottom lines should contain the notice
	bottomPortion := strings.Join(lines[len(lines)-3:], "\n")
	if !strings.Contains(bottomPortion, "bottom-left notice") {
		t.Error("expected bottom lines to contain 'bottom-left notice'")
	}
}

func TestRender_BottomCenter(t *testing.T) {
	m := newTestModel().WithPosition(BottomCenterPosition)
	m, _ = activateNotice(t, m, InfoKey, "bottom-center")

	content := makeContent(80, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	bottomPortion := strings.Join(lines[len(lines)-3:], "\n")
	if !strings.Contains(bottomPortion, "bottom-center") {
		t.Error("expected bottom lines to contain 'bottom-center'")
	}
}

func TestRender_BottomRight(t *testing.T) {
	m := newTestModel().WithPosition(BottomRightPosition)
	m, _ = activateNotice(t, m, InfoKey, "bottom-right")

	content := makeContent(80, 10)
	result := m.Render(content)
	lines := strings.Split(result, "\n")

	bottomPortion := strings.Join(lines[len(lines)-3:], "\n")
	if !strings.Contains(bottomPortion, "bottom-right") {
		t.Error("expected bottom lines to contain 'bottom-right'")
	}
}

func TestRender_NoticeWiderThanContent(t *testing.T) {
	m := newTestModelWithOpts(t, NotifyOpts{
		Width:    50,
		Duration: 3 * time.Second,
	})
	m, _ = activateNotice(t, m, InfoKey, "wide notice")

	// Narrow content (20 chars wide)
	content := makeContent(20, 10)
	result := m.Render(content)
	if !strings.Contains(result, "wide notice") {
		t.Error("expected result to contain 'wide notice'")
	}
}

func TestRender_NarrowContent(t *testing.T) {
	m := newTestModel()
	m, _ = activateNotice(t, m, InfoKey, "narrow content test")

	content := makeContent(5, 10)
	result := m.Render(content)
	if !strings.Contains(result, "narrow content test") {
		t.Error("expected result to contain 'narrow content test'")
	}
}
