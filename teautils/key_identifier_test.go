package teautils

import (
	"errors"
	"testing"
)

func TestParseKeyIdentifier_Valid(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"app.help"},
		{"tree.nav.up"},
		{"file-intent.commit"},
		{"my_app.key1"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			id, err := ParseKeyIdentifier(tt.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tt.input, err)
			}
			if string(id) != tt.input {
				t.Errorf("expected %q, got %q", tt.input, string(id))
			}
		})
	}
}

func TestParseKeyIdentifier_Empty(t *testing.T) {
	_, err := ParseKeyIdentifier("")
	if !errors.Is(err, ErrEmptyKeyIdentifier) {
		t.Errorf("expected ErrEmptyKeyIdentifier, got %v", err)
	}
}

func TestParseKeyIdentifier_NoDot(t *testing.T) {
	_, err := ParseKeyIdentifier("nodotshere")
	if !errors.Is(err, ErrKeyIdentifierMissingDot) {
		t.Errorf("expected ErrKeyIdentifierMissingDot, got %v", err)
	}
}

func TestParseKeyIdentifier_EmptyPart(t *testing.T) {
	_, err := ParseKeyIdentifier("app..help")
	if !errors.Is(err, ErrKeyIdentifierEmptyPart) {
		t.Errorf("expected ErrKeyIdentifierEmptyPart, got %v", err)
	}
}

func TestParseKeyIdentifier_InvalidPart(t *testing.T) {
	_, err := ParseKeyIdentifier("app.he lp")
	if !errors.Is(err, ErrKeyIdentifierInvalidPart) {
		t.Errorf("expected ErrKeyIdentifierInvalidPart, got %v", err)
	}
}

func TestMustParseKeyIdentifier_Panics(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expected panic for invalid input")
		}
	}()
	MustParseKeyIdentifier("invalid no dot")
}

func TestKeyIdentifier_String(t *testing.T) {
	id := MustParseKeyIdentifier("app.help")
	if id.String() != "app.help" {
		t.Errorf("expected 'app.help', got %q", id.String())
	}
}
