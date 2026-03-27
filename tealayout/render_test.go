package tealayout

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
)

// mockWidget implements SetSizer and Viewer for testing.
type mockWidget struct {
	width  int
	height int
	char   byte // fill character
}

func (m *mockWidget) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m *mockWidget) View() string {
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

	comp := NewRow(Percent100,
		NewColumn(Fixed(20), w1),
		NewColumn(Flex(1), w2),
	)
	comp.SetSize(80, 5)

	output, err := comp.Render()
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 5 {
		t.Errorf("height = %d, want 5", len(lines))
	}

	// w1 should have been set to 20x5
	if w1.width != 20 || w1.height != 5 {
		t.Errorf("w1 size = %dx%d, want 20x5", w1.width, w1.height)
	}

	// w2 should have been set to 60x5
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

	comp := NewRow(Percent100, NewColumn(Fixed(40), w))
	comp.SetSize(80, 10)

	_, err := comp.Render()
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

	comp := NewColumn(Percent100,
		NewRow(Fixed(3), w1),
		NewRow(Flex(1), w2),
	)
	comp.SetSize(80, 24)

	_, err := comp.Render()
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
	comp := NewRow(Percent100, NewColumn(Flex(1), w))
	comp.SetSize(80, 5)

	out1, _ := comp.Render()
	out2, _ := comp.Render()
	if out1 != out2 {
		t.Error("cached render returned different result")
	}
}

func TestRow_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	comp := NewRow(Percent100, NewColumn(Flex(1), w))
	comp.SetSize(80, 5)
	comp.Render()

	comp.MarkDirty()
	comp.ensureInner()
	r := comp.inner.(*row)
	if r.resolved {
		t.Error("MarkDirty should clear resolved flag")
	}
	if r.cachedOutput != "" {
		t.Error("MarkDirty should clear cached output")
	}
}

func TestColumn_MarkDirty(t *testing.T) {
	w := &mockWidget{char: 'X'}
	comp := NewColumn(Percent100, NewRow(Flex(1), w))
	comp.SetSize(80, 24)
	comp.Render()

	comp.MarkDirty()
	comp.ensureInner()
	c := comp.inner.(*column)
	if c.resolved {
		t.Error("MarkDirty should clear resolved flag")
	}
}

func TestRow_Render_Nested(t *testing.T) {
	// Column containing a Row — verify dimensions propagate
	w1 := &mockWidget{char: 'H'} // header
	w2 := &mockWidget{char: 'L'} // left pane
	w3 := &mockWidget{char: 'R'} // right pane
	w4 := &mockWidget{char: 'F'} // footer

	innerRow := NewRow(Percent100,
		NewColumn(Fixed(30), w2),
		NewColumn(Flex(1), w3),
	)

	comp := NewColumn(Percent100,
		NewRow(Fixed(3), w1),
		innerRow,
		NewRow(Fixed(1), w4),
	)
	comp.SetSize(80, 24)

	_, err := comp.Render()
	if err != nil {
		t.Fatal(err)
	}

	// Header: 80x3
	if w1.width != 80 || w1.height != 3 {
		t.Errorf("header = %dx%d, want 80x3", w1.width, w1.height)
	}
	// Footer: 80x1
	if w4.width != 80 || w4.height != 1 {
		t.Errorf("footer = %dx%d, want 80x1", w4.width, w4.height)
	}
	// Inner row gets 80x20, left=30 right=50
	if w2.width != 30 || w2.height != 20 {
		t.Errorf("left = %dx%d, want 30x20", w2.width, w2.height)
	}
	if w3.width != 50 || w3.height != 20 {
		t.Errorf("right = %dx%d, want 50x20", w3.width, w3.height)
	}
}
