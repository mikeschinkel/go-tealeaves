package teautils

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
)

func TestFormatKeyDisplay_SingleKey(t *testing.T) {
	meta := KeyMeta{
		Binding: key.NewBinding(key.WithKeys("enter")),
	}
	result := FormatKeyDisplay(meta)
	if result != "[Enter]" {
		t.Errorf("expected '[Enter]', got %q", result)
	}
}

func TestFormatKeyDisplay_MultipleKeys(t *testing.T) {
	meta := KeyMeta{
		Binding: key.NewBinding(key.WithKeys("up", "k")),
	}
	result := FormatKeyDisplay(meta)
	if result != "[Up]/[K]" {
		t.Errorf("expected '[Up]/[K]', got %q", result)
	}
}

func TestFormatKeyDisplay_Deduplication(t *testing.T) {
	// Both "a" and "A" should title-case to "A" and be deduplicated
	meta := KeyMeta{
		Binding: key.NewBinding(key.WithKeys("a")),
		// Simulate scenario where display would produce duplicates
		DisplayKeys: []string{"A", "A"},
	}
	result := FormatKeyDisplay(meta)
	if result != "[A]" {
		t.Errorf("expected '[A]' after dedup, got %q", result)
	}
}

func TestFormatKeyDisplay_CustomDisplayKeys(t *testing.T) {
	meta := KeyMeta{
		Binding:     key.NewBinding(key.WithKeys(" ")),
		DisplayKeys: []string{"Space"},
	}
	result := FormatKeyDisplay(meta)
	if result != "[Space]" {
		t.Errorf("expected '[Space]', got %q", result)
	}
}

func TestProperCaseShortcut(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"ctrl+c", "Ctrl+C"},
		{"esc", "Esc"},
		{"shift+tab", "Shift+Tab"},
		{"alt+v", "Alt+V"},
		{"a", "A"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ProperCaseShortcut(tt.input)
			if got != tt.want {
				t.Errorf("ProperCaseShortcut(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestProperCaseShortcut_AlreadyProper(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Ctrl+C", "Ctrl+C"},
		{"Esc", "Esc"},
		{"Shift+Tab", "Shift+Tab"},
		{"A", "A"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := ProperCaseShortcut(tt.input)
			if got != tt.want {
				t.Errorf("ProperCaseShortcut(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetSortedCategories_PreferredOrder(t *testing.T) {
	categories := map[string][]KeyMeta{
		"System":     {{ID: "sys.quit"}},
		"Navigation": {{ID: "nav.up"}},
		"Actions":    {{ID: "act.save"}},
		"Zebra":      {{ID: "z.custom"}},
	}
	preferred := []string{"Navigation", "Actions", "System"}
	result := GetSortedCategories(categories, preferred)

	if len(result) != 4 {
		t.Fatalf("expected 4 categories, got %d", len(result))
	}
	if result[0] != "Navigation" {
		t.Errorf("expected first category 'Navigation', got %q", result[0])
	}
	if result[1] != "Actions" {
		t.Errorf("expected second category 'Actions', got %q", result[1])
	}
	if result[2] != "System" {
		t.Errorf("expected third category 'System', got %q", result[2])
	}
	// Zebra not in preferred order, goes last
	if result[3] != "Zebra" {
		t.Errorf("expected fourth category 'Zebra', got %q", result[3])
	}
}

func TestGetSortedCategories_NoPreference(t *testing.T) {
	categories := map[string][]KeyMeta{
		"Zebra":      {{ID: "z.key"}},
		"Alpha":      {{ID: "a.key"}},
		"Middle":     {{ID: "m.key"}},
	}
	result := GetSortedCategories(categories, nil)

	if len(result) != 3 {
		t.Fatalf("expected 3 categories, got %d", len(result))
	}
	if result[0] != "Alpha" {
		t.Errorf("expected first category 'Alpha', got %q", result[0])
	}
	if result[1] != "Middle" {
		t.Errorf("expected second category 'Middle', got %q", result[1])
	}
	if result[2] != "Zebra" {
		t.Errorf("expected third category 'Zebra', got %q", result[2])
	}
}

func TestDefaultHelpVisorStyle(t *testing.T) {
	style := DefaultHelpVisorStyle()
	if style.KeyColumnGap == 0 {
		t.Error("expected non-zero KeyColumnGap")
	}
	if len(style.CategoryOrder) == 0 {
		t.Error("expected non-empty CategoryOrder")
	}
}
