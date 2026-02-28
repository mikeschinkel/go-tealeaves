package teanotify

import (
	"errors"
	"strings"
	"testing"
)

// --- defaultNotices tests ---

func TestDefaultNotices_ASCIIPrefixes(t *testing.T) {
	defs := defaultNotices(false, false)
	if len(defs) != 4 {
		t.Fatalf("expected 4 definitions, got %d", len(defs))
	}
	if defs[0].Prefix != InfoASCIIPrefix {
		t.Errorf("Info prefix: expected %q, got %q", InfoASCIIPrefix, defs[0].Prefix)
	}
	if defs[1].Prefix != WarningASCIIPrefix {
		t.Errorf("Warn prefix: expected %q, got %q", WarningASCIIPrefix, defs[1].Prefix)
	}
	if defs[2].Prefix != ErrorASCIIPrefix {
		t.Errorf("Error prefix: expected %q, got %q", ErrorASCIIPrefix, defs[2].Prefix)
	}
	if defs[3].Prefix != DebugASCIIPrefix {
		t.Errorf("Debug prefix: expected %q, got %q", DebugASCIIPrefix, defs[3].Prefix)
	}
	if defs[0].ForeColor != InfoColor {
		t.Errorf("Info color: expected %q, got %q", InfoColor, defs[0].ForeColor)
	}
	if defs[1].ForeColor != WarnColor {
		t.Errorf("Warn color: expected %q, got %q", WarnColor, defs[1].ForeColor)
	}
	if defs[2].ForeColor != ErrorColor {
		t.Errorf("Error color: expected %q, got %q", ErrorColor, defs[2].ForeColor)
	}
	if defs[3].ForeColor != DebugColor {
		t.Errorf("Debug color: expected %q, got %q", DebugColor, defs[3].ForeColor)
	}
}

func TestDefaultNotices_NerdFontPrefixes(t *testing.T) {
	defs := defaultNotices(true, false)
	if defs[0].Prefix != InfoNerdSymbol {
		t.Errorf("Info prefix: expected NerdFont symbol, got %q", defs[0].Prefix)
	}
	if defs[1].Prefix != WarnNerdSymbol {
		t.Errorf("Warn prefix: expected NerdFont symbol, got %q", defs[1].Prefix)
	}
	if defs[2].Prefix != ErrorNerdSymbol {
		t.Errorf("Error prefix: expected NerdFont symbol, got %q", defs[2].Prefix)
	}
	if defs[3].Prefix != DebugNerdSymbol {
		t.Errorf("Debug prefix: expected NerdFont symbol, got %q", defs[3].Prefix)
	}
}

func TestDefaultNotices_UnicodePrefixes(t *testing.T) {
	defs := defaultNotices(false, true)
	if defs[0].Prefix != InfoUnicodePrefix {
		t.Errorf("Info prefix: expected %q, got %q", InfoUnicodePrefix, defs[0].Prefix)
	}
	if defs[1].Prefix != WarningUnicodePrefix {
		t.Errorf("Warn prefix: expected %q, got %q", WarningUnicodePrefix, defs[1].Prefix)
	}
	if defs[2].Prefix != ErrorUnicodePrefix {
		t.Errorf("Error prefix: expected %q, got %q", ErrorUnicodePrefix, defs[2].Prefix)
	}
	if defs[3].Prefix != DebugUnicodePrefix {
		t.Errorf("Debug prefix: expected %q, got %q", DebugUnicodePrefix, defs[3].Prefix)
	}
}

func TestDefaultNotices_NerdFontTakesPrecedence(t *testing.T) {
	defs := defaultNotices(true, true)
	// NerdFont wins when both are true (switch precedence)
	if defs[0].Prefix != InfoNerdSymbol {
		t.Errorf("expected NerdFont to take precedence, Info prefix: got %q", defs[0].Prefix)
	}
}

// --- registerNotice tests ---

func TestRegisterNotice_ValidDefinition(t *testing.T) {
	m := NotifyModel{
		noticeTypes: make(map[NoticeKey]NoticeDefinition),
	}
	def := NoticeDefinition{
		Key:       "Custom",
		ForeColor: "#AABBCC",
		Prefix:    "[C]",
	}
	out, err := registerNotice(m, def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out.noticeTypes["Custom"]; !ok {
		t.Fatal("expected 'Custom' to be present in noticeTypes")
	}
}

func TestRegisterNotice_EmptyKey(t *testing.T) {
	m := NotifyModel{
		noticeTypes: make(map[NoticeKey]NoticeDefinition),
	}
	def := NoticeDefinition{
		Key:       "",
		ForeColor: "#AABBCC",
	}
	_, err := registerNotice(m, def)
	if !errors.Is(err, ErrInvalidNoticeKey) {
		t.Fatalf("expected ErrInvalidNoticeKey, got: %v", err)
	}
}

func TestRegisterNotice_InvalidColor(t *testing.T) {
	m := NotifyModel{
		noticeTypes: make(map[NoticeKey]NoticeDefinition),
	}
	def := NoticeDefinition{
		Key:       "Bad",
		ForeColor: "not-a-color",
	}
	_, err := registerNotice(m, def)
	if !errors.Is(err, ErrInvalidColor) {
		t.Fatalf("expected ErrInvalidColor, got: %v", err)
	}
}

// --- notice render tests (Layer 2) ---

func TestNotice_Render_ContainsMessage(t *testing.T) {
	n := &notice{
		message:     "File saved",
		prefix:      "(i)",
		foreColor:   infoColor,
		style:       baseStyle,
		width:       40,
		curLerpStep: 0.5,
	}
	result := n.render()
	if !strings.Contains(result, "File saved") {
		t.Fatalf("expected render output to contain 'File saved', got %q", result)
	}
	if !strings.Contains(result, "(i)") {
		t.Fatalf("expected render output to contain prefix '(i)', got %q", result)
	}
}

func TestNotice_Render_DynamicWidth(t *testing.T) {
	// Short message with minWidth: width clamps to minWidth
	n := &notice{
		message:     "ok",
		prefix:      "(i)",
		foreColor:   infoColor,
		style:       baseStyle,
		width:       60,
		minWidth:    20,
		curLerpStep: 0.5,
	}
	result := n.render()
	if !strings.Contains(result, "ok") {
		t.Fatalf("expected render output to contain 'ok', got %q", result)
	}

	// Long message: width clamps at max width
	n2 := &notice{
		message:     strings.Repeat("x", 100),
		prefix:      "(i)",
		foreColor:   infoColor,
		style:       baseStyle,
		width:       30,
		minWidth:    10,
		curLerpStep: 0.5,
	}
	result2 := n2.render()
	if !strings.Contains(result2, "xxx") {
		t.Fatalf("expected render output to contain message content, got %q", result2)
	}
}
