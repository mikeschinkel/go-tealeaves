package tealayout

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
)

// mockWidget implements SetSizer and ContentProvider for testing.
type mockWidget struct {
	width  int
	height int
	char   byte // fill character
}

func (m *mockWidget) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *mockWidget) Content() string {
	if m.width <= 0 || m.height <= 0 {
		return ""
	}
	line := strings.Repeat(string(m.char), m.width)
	lines := make([]string, m.height)
	for i := range lines {
		lines[i] = line
	}
	return strings.Join(lines, "\n")
}

// mockStyledWidget adds Styler interface.
type mockStyledWidget struct {
	mockWidget
	style lipgloss.Style
}

func (m *mockStyledWidget) Style() lipgloss.Style {
	return m.style
}

func TestRow_Render_Dimensions(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	pane := NewRow(Percent100,
		NewColumn(Fixed(20), NewElement(w1)),
		NewColumn(Flex(1), NewElement(w2)),
	)
	pane.SetSize(80, 5)

	output, err := pane.Render()
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 5 {
		t.Errorf("height = %d, want 5", len(lines))
	}

	if w1.width != 20 || w1.height != 5 {
		t.Errorf("w1 size = %dx%d, want 20x5", w1.width, w1.height)
	}

	if w2.width != 60 || w2.height != 5 {
		t.Errorf("w2 size = %dx%d, want 60x5", w2.width, w2.height)
	}
}

func TestRow_Render_ContentDimensions(t *testing.T) {
	style := lipgloss.NewStyle().Padding(1, 2).Border(lipgloss.NormalBorder())
	w := &mockStyledWidget{
		mockWidget: mockWidget{char: 'X'},
		style:      style,
	}

	pane := NewRow(Percent100, NewColumn(Fixed(40), NewElement(w)))
	pane.SetSize(80, 10)

	_, err := pane.Render()
	if err != nil {
		t.Fatal(err)
	}

	// 40 - 2*2(padding) - 2(border) = 34
	// 10 - 2(padding) - 2(border) = 6
	if w.width != 34 {
		t.Errorf("content width = %d, want 34", w.width)
	}
	if w.height != 6 {
		t.Errorf("content height = %d, want 6", w.height)
	}
}

func TestColumn_Render_Dimensions(t *testing.T) {
	w1 := &mockWidget{char: 'A'}
	w2 := &mockWidget{char: 'B'}

	pane := NewColumn(Percent100,
		NewRow(Fixed(3), NewElement(w1)),
		NewRow(Flex(1), NewElement(w2)),
	)
	pane.SetSize(80, 24)

	_, err := pane.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width != 80 || w1.height != 3 {
		t.Errorf("w1 size = %dx%d, want 80x3", w1.width, w1.height)
	}
	if w2.width != 80 || w2.height != 21 {
		t.Errorf("w2 size = %dx%d, want 80x21", w2.width, w2.height)
	}
}

func TestRow_Render_Cached(t *testing.T) {
	w := &mockWidget{char: 'X'}
	pane := NewRow(Percent100, NewColumn(Flex(1), NewElement(w)))
	pane.SetSize(80, 5)

	out1, _ := pane.Render()
	out2, _ := pane.Render()
	if out1 != out2 {
		t.Error("cached render returned different result")
	}
}

func TestRow_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	pane := NewRow(Percent100, NewColumn(Flex(1), NewElement(w)))
	pane.SetSize(80, 5)
	if _, err := pane.Render(); err != nil {
		t.Fatalf("Render: %v", err)
	}

	pane.MarkDirty()
	if pane.resolved {
		t.Error("MarkDirty should clear resolved flag")
	}
	if pane.cachedOutput != "" {
		t.Error("MarkDirty should clear cached output")
	}
}

func TestColumn_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	pane := NewColumn(Percent100, NewRow(Flex(1), NewElement(w)))
	pane.SetSize(80, 24)
	if _, err := pane.Render(); err != nil {
		t.Fatalf("Render: %v", err)
	}

	pane.MarkDirty()
	if pane.resolved {
		t.Error("MarkDirty should clear resolved flag")
	}
}

func TestRow_Render_Nested(t *testing.T) {
	w1 := &mockWidget{char: 'H'} // header
	w2 := &mockWidget{char: 'L'} // left pane
	w3 := &mockWidget{char: 'R'} // right pane
	w4 := &mockWidget{char: 'F'} // footer

	innerRow := NewRow(Percent100,
		NewColumn(Fixed(30), NewElement(w2)),
		NewColumn(Flex(1), NewElement(w3)),
	)

	pane := NewColumn(Percent100,
		NewRow(Fixed(3), NewElement(w1)),
		innerRow,
		NewRow(Fixed(1), NewElement(w4)),
	)
	pane.SetSize(80, 24)

	_, err := pane.Render()
	if err != nil {
		t.Fatal(err)
	}

	if w1.width != 80 || w1.height != 3 {
		t.Errorf("header = %dx%d, want 80x3", w1.width, w1.height)
	}
	if w4.width != 80 || w4.height != 1 {
		t.Errorf("footer = %dx%d, want 80x1", w4.width, w4.height)
	}
	if w2.width != 30 || w2.height != 20 {
		t.Errorf("left = %dx%d, want 30x20", w2.width, w2.height)
	}
	if w3.width != 50 || w3.height != 20 {
		t.Errorf("right = %dx%d, want 50x20", w3.width, w3.height)
	}
}
