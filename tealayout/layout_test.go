package tealayout

import (
	"errors"
	"testing"
)

func TestLayout_RenderWithoutSetSize_AutoDetectDisabled(t *testing.T) {
	layout := NewLayout(NewRow(Percent100, NewColumn(Flex(1))))
	_, err := layout.Render()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrZeroDimensions) {
		t.Errorf("expected ErrZeroDimensions, got %v", err)
	}
}

func TestLayout_RenderWithSetSize(t *testing.T) {
	w := &mockWidget{char: 'X'}
	layout := NewLayout(NewRow(Percent100, NewColumn(Flex(1), w)))
	layout.SetSize(80, 24)

	output, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
	if w.width != 80 || w.height != 24 {
		t.Errorf("widget size = %dx%d, want 80x24", w.width, w.height)
	}
}

func TestLayout_AutoDetectWithSizeSource(t *testing.T) {
	w := &mockWidget{char: 'X'}
	layout := NewLayout(
		NewRow(Percent100, NewColumn(Flex(1), w)),
		WithAutoDetectSize(true),
		WithSizeSource(func() (int, int, error) {
			return 120, 40, nil
		}),
	)

	output, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
	if output == "" {
		t.Error("expected non-empty output")
	}
	if w.width != 120 || w.height != 40 {
		t.Errorf("widget size = %dx%d, want 120x40", w.width, w.height)
	}
}

func TestLayout_AutoDetectFailure(t *testing.T) {
	layout := NewLayout(
		NewRow(Percent100, NewColumn(Flex(1))),
		WithAutoDetectSize(true),
		WithSizeSource(func() (int, int, error) {
			return 0, 0, errors.New("no terminal")
		}),
	)

	_, err := layout.Render()
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrZeroDimensions) {
		t.Errorf("expected ErrZeroDimensions, got %v", err)
	}
}

func TestLayout_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	layout := NewLayout(NewRow(Percent100, NewColumn(Flex(1), w)))
	layout.SetSize(80, 24)
	layout.Render()

	layout.MarkDirty()
	// Should re-render successfully
	_, err := layout.Render()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLayout_Accessors(t *testing.T) {
	root := NewRow(Percent100, NewColumn(Flex(1)))
	layout := NewLayout(root)
	layout.SetSize(80, 24)

	if layout.Width() != 80 {
		t.Errorf("Width() = %d, want 80", layout.Width())
	}
	if layout.Height() != 24 {
		t.Errorf("Height() = %d, want 24", layout.Height())
	}
	if layout.Root() != root {
		t.Error("Root() != original root")
	}
}
