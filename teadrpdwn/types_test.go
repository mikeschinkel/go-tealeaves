package teadrpdwn

import "testing"

func TestToOptions(t *testing.T) {
	strs := []string{"Alpha", "Beta", "Gamma"}
	opts := ToOptions(strs)

	if len(opts) != 3 {
		t.Fatalf("expected 3 options, got %d", len(opts))
	}

	for i, s := range strs {
		if opts[i].Text != s {
			t.Errorf("opts[%d].Text = %q, want %q", i, opts[i].Text, s)
		}
		if opts[i].Value != s {
			t.Errorf("opts[%d].Value = %q, want %q", i, opts[i].Value, s)
		}
	}
}

func TestToOptions_Empty(t *testing.T) {
	opts := ToOptions([]string{})
	if len(opts) != 0 {
		t.Errorf("expected 0 options, got %d", len(opts))
	}
}
