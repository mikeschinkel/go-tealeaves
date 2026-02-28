package teautils

import (
	"errors"
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func testBinding(keys ...string) key.Binding {
	return key.NewBinding(key.WithKeys(keys...), key.WithHelp(keys[0], "test help"))
}

func TestKeyRegistry_Register(t *testing.T) {
	r := NewKeyRegistry()
	err := r.Register(KeyMeta{
		ID:      "app.help",
		Binding: testBinding("?"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	meta, err := r.Get("app.help")
	if err != nil {
		t.Fatalf("unexpected error on Get: %v", err)
	}
	if meta.ID != "app.help" {
		t.Errorf("expected ID='app.help', got %q", meta.ID)
	}
}

func TestKeyRegistry_Register_DefaultsHelpText(t *testing.T) {
	r := NewKeyRegistry()
	err := r.Register(KeyMeta{
		ID:      "app.quit",
		Binding: key.NewBinding(key.WithKeys("q"), key.WithHelp("q", "Quit application")),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	meta, err := r.Get("app.quit")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if meta.HelpText != "Quit application" {
		t.Errorf("expected HelpText='Quit application', got %q", meta.HelpText)
	}
}

func TestKeyRegistry_Register_EmptyID(t *testing.T) {
	r := NewKeyRegistry()
	err := r.Register(KeyMeta{
		Binding: testBinding("x"),
	})
	if !errors.Is(err, ErrEmptyKeyID) {
		t.Errorf("expected ErrEmptyKeyID, got %v", err)
	}
}

func TestKeyRegistry_RegisterMany(t *testing.T) {
	r := NewKeyRegistry()
	metas := []KeyMeta{
		{ID: "nav.up", Binding: testBinding("up")},
		{ID: "nav.down", Binding: testBinding("down")},
		{ID: "nav.left", Binding: testBinding("left")},
	}
	err := r.RegisterMany(metas)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, m := range metas {
		got, err := r.Get(m.ID)
		if err != nil {
			t.Errorf("failed to get %q: %v", m.ID, err)
		}
		if got.ID != m.ID {
			t.Errorf("expected ID=%q, got %q", m.ID, got.ID)
		}
	}
}

func TestKeyRegistry_Get_NotFound(t *testing.T) {
	r := NewKeyRegistry()
	_, err := r.Get("nonexistent.key")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound, got %v", err)
	}
}

func TestKeyRegistry_Clear(t *testing.T) {
	r := NewKeyRegistry()
	err := r.Register(KeyMeta{
		ID:      "app.help",
		Binding: testBinding("?"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	r.Clear()

	_, err = r.Get("app.help")
	if !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("expected ErrKeyNotFound after Clear, got %v", err)
	}
}

func TestKeyRegistry_ForStatusBar(t *testing.T) {
	r := NewKeyRegistry()
	err := r.RegisterMany([]KeyMeta{
		{ID: "nav.up", Binding: testBinding("up"), StatusBar: true},
		{ID: "nav.down", Binding: testBinding("down"), StatusBar: false},
		{ID: "app.help", Binding: testBinding("?"), StatusBar: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := r.ForStatusBar()
	if len(result) != 2 {
		t.Fatalf("expected 2 status bar keys, got %d", len(result))
	}
	// Verify registration order preserved
	if result[0].ID != "nav.up" {
		t.Errorf("expected first key 'nav.up', got %q", result[0].ID)
	}
	if result[1].ID != "app.help" {
		t.Errorf("expected second key 'app.help', got %q", result[1].ID)
	}
}

func TestKeyRegistry_ForHelpModal(t *testing.T) {
	r := NewKeyRegistry()
	err := r.RegisterMany([]KeyMeta{
		{ID: "nav.up", Binding: testBinding("up"), HelpModal: true},
		{ID: "nav.down", Binding: testBinding("down"), HelpModal: false},
		{ID: "app.help", Binding: testBinding("?"), HelpModal: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := r.ForHelpModal()
	if len(result) != 2 {
		t.Fatalf("expected 2 help modal keys, got %d", len(result))
	}
	if result[0].ID != "nav.up" {
		t.Errorf("expected first key 'nav.up', got %q", result[0].ID)
	}
	if result[1].ID != "app.help" {
		t.Errorf("expected second key 'app.help', got %q", result[1].ID)
	}
}

func TestKeyRegistry_ByCategory(t *testing.T) {
	r := NewKeyRegistry()
	err := r.RegisterMany([]KeyMeta{
		{ID: "nav.up", Binding: testBinding("up"), HelpModal: true, Category: "Navigation"},
		{ID: "nav.down", Binding: testBinding("down"), HelpModal: true, Category: "Navigation"},
		{ID: "act.save", Binding: testBinding("ctrl+s"), HelpModal: true, Category: "Actions"},
		{ID: "misc.x", Binding: testBinding("x"), HelpModal: true},          // No category → "Other"
		{ID: "hidden.y", Binding: testBinding("y"), HelpModal: false, Category: "Navigation"}, // Not in help modal
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result := r.ByCategory()

	if len(result["Navigation"]) != 2 {
		t.Errorf("expected 2 Navigation keys, got %d", len(result["Navigation"]))
	}
	if len(result["Actions"]) != 1 {
		t.Errorf("expected 1 Actions key, got %d", len(result["Actions"]))
	}
	if len(result["Other"]) != 1 {
		t.Errorf("expected 1 Other key, got %d", len(result["Other"]))
	}

	// Verify registration order within category
	if result["Navigation"][0].ID != "nav.up" {
		t.Errorf("expected first Navigation key 'nav.up', got %q", result["Navigation"][0].ID)
	}
	if result["Navigation"][1].ID != "nav.down" {
		t.Errorf("expected second Navigation key 'nav.down', got %q", result["Navigation"][1].ID)
	}
}
