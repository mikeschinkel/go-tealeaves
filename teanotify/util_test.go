package teanotify

import (
	"strings"
	"testing"
)

// --- getLines tests ---

func TestGetLines_MultiLine(t *testing.T) {
	lines, widest := getLines("abc\nde\nfghij")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if widest != 5 {
		t.Fatalf("expected widest 5, got %d", widest)
	}
}

func TestGetLines_EmptyString(t *testing.T) {
	lines, widest := getLines("")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
	if widest != 0 {
		t.Fatalf("expected widest 0, got %d", widest)
	}
}

func TestGetLines_VaryingWidths(t *testing.T) {
	lines, widest := getLines("ab\nabcdef\nabc")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if widest != 6 {
		t.Fatalf("expected widest 6, got %d", widest)
	}
}

// --- isANSITerminator tests ---

func TestIsANSITerminator_ValidTerminators(t *testing.T) {
	valid := []rune{'m', 'H', 'A', 0x40, 0x5A, 0x61, 0x7A}
	for _, c := range valid {
		t.Run(string(c), func(t *testing.T) {
			if !isANSITerminator(c) {
				t.Errorf("expected %q (0x%02X) to be a terminator", c, c)
			}
		})
	}
}

func TestIsANSITerminator_NonTerminators(t *testing.T) {
	invalid := []rune{'0', ' ', '!', 0x3F, 0x5B, 0x7B}
	for _, c := range invalid {
		t.Run(string(c), func(t *testing.T) {
			if isANSITerminator(c) {
				t.Errorf("expected %q (0x%02X) to NOT be a terminator", c, c)
			}
		})
	}
}

// --- cutLeft tests ---

func TestCutLeft_PlainText_CutZero(t *testing.T) {
	result := cutLeft("hello world", 0)
	if result != "hello world" {
		t.Fatalf("expected %q, got %q", "hello world", result)
	}
}

func TestCutLeft_PlainText_CutPartial(t *testing.T) {
	result := cutLeft("hello world", 6)
	if result != "world" {
		t.Fatalf("expected %q, got %q", "world", result)
	}
}

func TestCutLeft_PlainText_CutExceedsLength(t *testing.T) {
	result := cutLeft("hello", 20)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestCutLeft_WithANSI(t *testing.T) {
	// Red text: \x1b[31mHello\x1b[0m World
	input := "\x1b[31mHello\x1b[0m World"
	result := cutLeft(input, 6)
	// After cutting 6 cells ("Hello "), we should get "World"
	// The reset sequence ([0m) should be stripped; non-reset ANSI preserved
	if !strings.Contains(result, "World") {
		t.Fatalf("expected result to contain 'World', got %q", result)
	}
	// Reset sequence should not be preserved since it was stripped
	if strings.Contains(result, "\x1b[0m") {
		t.Fatalf("expected reset sequence to be stripped, got %q", result)
	}
}

// --- cutRight tests ---

func TestCutRight_PlainText_KeepZero(t *testing.T) {
	result := cutRight("hello world", 0)
	if result != "" {
		t.Fatalf("expected empty string, got %q", result)
	}
}

func TestCutRight_PlainText_KeepPartial(t *testing.T) {
	result := cutRight("hello world", 5)
	// Should keep first 5 cells: "hello"
	// cutRight appends a reset suffix
	if !strings.HasPrefix(result, "hello") {
		t.Fatalf("expected result to start with 'hello', got %q", result)
	}
}

func TestCutRight_PlainText_KeepExceedsLength(t *testing.T) {
	result := cutRight("hello", 20)
	if !strings.Contains(result, "hello") {
		t.Fatalf("expected result to contain 'hello', got %q", result)
	}
}

func TestCutRight_WithANSI(t *testing.T) {
	// Red text: \x1b[31mHello World\x1b[0m
	input := "\x1b[31mHello World\x1b[0m"
	result := cutRight(input, 5)
	// Should keep 5 visible cells: "Hello"
	if !strings.Contains(result, "Hello") {
		t.Fatalf("expected result to contain 'Hello', got %q", result)
	}
	// Should not contain the content beyond the keep width
	if strings.Contains(result, "World") {
		t.Fatalf("expected result to not contain 'World', got %q", result)
	}
}

func TestCutRight_AppendsReset(t *testing.T) {
	result := cutRight("hello world", 5)
	if !strings.HasSuffix(result, "\x1b[0m") {
		t.Fatalf("expected result to end with reset sequence, got %q", result)
	}
}

// --- hangingWrap tests ---

func TestHangingWrap_ShortMessage(t *testing.T) {
	result := hangingWrap("INFO", "ok", 30)
	if result != "INFO ok" {
		t.Fatalf("expected %q, got %q", "INFO ok", result)
	}
}

func TestHangingWrap_LongMessage(t *testing.T) {
	result := hangingWrap("(!)", "this is a long message that should wrap onto multiple lines", 20)
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected multiple lines, got %d", len(lines))
	}
	// First line starts with prefix
	if !strings.HasPrefix(lines[0], "(!) ") {
		t.Fatalf("expected first line to start with prefix, got %q", lines[0])
	}
	// Continuation lines should be indented with spaces matching prefix width
	prefixWidth := len("(!) ")
	indent := strings.Repeat(" ", prefixWidth)
	for i := 1; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], indent) {
			t.Fatalf("expected line %d to start with %d-space indent, got %q", i, prefixWidth, lines[i])
		}
	}
}

func TestHangingWrap_NarrowWidth(t *testing.T) {
	// When avail < 1, fallback to prefix + msg (no wrapping)
	result := hangingWrap("(!) ", "msg", 1)
	expected := "(!)  msg"
	if result != expected {
		t.Fatalf("expected %q, got %q", expected, result)
	}
}
