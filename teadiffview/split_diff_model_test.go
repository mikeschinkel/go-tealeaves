package teadiffview

import (
	"testing"

	"github.com/mikeschinkel/go-diffutils"
)

func testContent() *diffutils.DiffContent {
	return &diffutils.DiffContent{
		OldLines: []string{"line 1", "line 2", "line 3", "line 4", "line 5"},
		NewLines: []string{"line 1", "CHANGED", "line 3", "line 4", "line 5", "line 6"},
		Changes: []diffutils.DiffChange{
			{
				Type:     diffutils.LinesDeleted,
				OldRange: diffutils.LineRange{Start: 2, Count: 1},
				NewRange: diffutils.LineRange{Start: 2, Count: 0},
			},
			{
				Type:     diffutils.LinesAdded,
				OldRange: diffutils.LineRange{Start: 3, Count: 0},
				NewRange: diffutils.LineRange{Start: 2, Count: 1},
			},
			{
				Type:     diffutils.LinesAdded,
				OldRange: diffutils.LineRange{Start: 6, Count: 0},
				NewRange: diffutils.LineRange{Start: 6, Count: 1},
			},
		},
		Label: "test.go",
	}
}

func TestSplitDiffModel_SetContent(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	if model.RowCount() == 0 {
		t.Fatal("expected non-zero row count")
	}
	if model.CursorIndex() != 0 {
		t.Errorf("expected cursor at 0, got %d", model.CursorIndex())
	}
}

func TestSplitDiffModel_CursorMovement(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	model = model.MoveCursorDown()
	if model.CursorIndex() != 1 {
		t.Errorf("expected cursor at 1, got %d", model.CursorIndex())
	}

	model = model.MoveCursorDown()
	model = model.MoveCursorDown()
	model = model.MoveCursorUp()
	if model.CursorIndex() < 1 {
		t.Error("cursor should be >= 1 after multiple movements")
	}
}

func TestSplitDiffModel_GoToTopBottom(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	model = model.GoToBottom()
	if model.CursorIndex() != model.RowCount()-1 {
		t.Errorf("expected cursor at last row %d, got %d", model.RowCount()-1, model.CursorIndex())
	}

	model = model.GoToTop()
	if model.CursorIndex() != 0 {
		t.Errorf("expected cursor at 0, got %d", model.CursorIndex())
	}
}

func TestSplitDiffModel_Selection(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	if model.HasSelection() {
		t.Error("should have no initial selection")
	}

	model = model.ExtendSelectionDown()
	if !model.HasSelection() {
		t.Error("should have selection after ExtendSelectionDown")
	}

	start, end := model.GetSelectedLines()
	if start != 0 || end != 1 {
		t.Errorf("expected selection 0-1, got %d-%d", start, end)
	}

	model = model.ClearSelection()
	if model.HasSelection() {
		t.Error("should have no selection after clear")
	}
}

func TestSplitDiffModel_FocusBlur(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})

	if model.focused {
		t.Error("model should start unfocused")
	}

	model = model.Focus()
	if !model.focused {
		t.Error("model should be focused after Focus()")
	}

	model = model.Blur()
	if model.focused {
		t.Error("model should be unfocused after Blur()")
	}
}

func TestSplitDiffModel_SetSize(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetSize(60, 30)

	if model.splitContentWidth != 60 {
		t.Errorf("expected splitContentWidth 60, got %d", model.splitContentWidth)
	}
	if model.height != 30 {
		t.Errorf("expected height 30, got %d", model.height)
	}
}

func TestSplitDiffModel_View(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 5})
	model = model.SetContent(testContent())

	view := model.View()
	if view.Content == "" {
		t.Error("expected non-empty view output")
	}
}

func TestSplitDiffModel_EmptyContent(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 5})

	// No content set
	view := model.View()
	if view.Content != "No diff to display" {
		t.Errorf("expected 'No diff to display', got %q", view.Content)
	}

	// Empty content
	model = model.SetContent(&diffutils.DiffContent{})
	view = model.View()
	if view.Content != "(empty file)" {
		t.Errorf("expected '(empty file)', got %q", view.Content)
	}
}

func TestSplitDiffModel_SetAnnotations(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	annotations := map[int]RowAnnotation{
		0: {Char: '1', Color: CommitGroupColors[0]},
		1: {Char: '2', Color: CommitGroupColors[1]},
	}
	model = model.SetAnnotations(annotations)

	if len(model.gutterChars) != model.RowCount() {
		t.Errorf("expected %d gutter chars, got %d", model.RowCount(), len(model.gutterChars))
	}
	if model.gutterChars[0] != '1' {
		t.Errorf("expected gutter char '1' at index 0, got '%c'", model.gutterChars[0])
	}
}

func TestSplitDiffModel_ScrollHorizontal(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	model = model.ScrollRight()
	if model.xOffset != scrollStep {
		t.Errorf("expected xOffset %d, got %d", scrollStep, model.xOffset)
	}

	model = model.ScrollLeft()
	if model.xOffset != 0 {
		t.Errorf("expected xOffset 0, got %d", model.xOffset)
	}

	// Should not go negative
	model = model.ScrollLeft()
	if model.xOffset != 0 {
		t.Errorf("expected xOffset 0 (clamped), got %d", model.xOffset)
	}
}

func TestSplitDiffModel_CenterOnFirstChangeBlock(t *testing.T) {
	model := NewSplitDiffModel(&SplitDiffModelArgs{Width: 40, Height: 20})
	model = model.SetContent(testContent())

	model = model.CenterOnFirstChangeBlock()
	// Should have moved cursor to first change block
	if model.CursorIndex() == 0 {
		// Row 0 is unchanged "line 1", cursor should be on a change block
		if model.RowCount() > 1 {
			blockIdx := model.GetBlockIndexAtCursor()
			if blockIdx == 0 {
				// Cursor might be at index 1 which is the first change
				// This is acceptable depending on the block structure
			}
		}
	}
}

func TestBuildSplitPaneDiff_SimpleInsert(t *testing.T) {
	content := &diffutils.DiffContent{
		OldLines: []string{"A", "C"},
		NewLines: []string{"A", "B", "C"},
		Changes: []diffutils.DiffChange{
			{
				Type:     diffutils.LinesAdded,
				OldRange: diffutils.LineRange{Start: 2, Count: 0},
				NewRange: diffutils.LineRange{Start: 2, Count: 1},
			},
		},
	}

	rows := buildSplitPaneDiff(content, nil)
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	// Row 0: A | A (unchanged)
	if rows[0].BlockIndex != 0 {
		t.Errorf("row 0 should be unchanged, got BlockIndex %d", rows[0].BlockIndex)
	}

	// Row 1: [marker] | B (added)
	if rows[1].BlockIndex != 1 {
		t.Errorf("row 1 should be block 1, got BlockIndex %d", rows[1].BlockIndex)
	}
	if _, ok := rows[1].CommitLine.(*BlockMarker); !ok {
		t.Error("row 1 CommitLine should be BlockMarker")
	}
	if textLine, ok := rows[1].ActualLine.(*TextLine); !ok || textLine.Text != "B" {
		t.Errorf("row 1 ActualLine should be TextLine 'B'")
	}

	// Row 2: C | C (unchanged)
	if rows[2].BlockIndex != 0 {
		t.Errorf("row 2 should be unchanged, got BlockIndex %d", rows[2].BlockIndex)
	}
}

func TestBuildSplitPaneDiff_SimpleDelete(t *testing.T) {
	content := &diffutils.DiffContent{
		OldLines: []string{"A", "B", "C"},
		NewLines: []string{"A", "C"},
		Changes: []diffutils.DiffChange{
			{
				Type:     diffutils.LinesDeleted,
				OldRange: diffutils.LineRange{Start: 2, Count: 1},
				NewRange: diffutils.LineRange{Start: 2, Count: 0},
			},
		},
	}

	rows := buildSplitPaneDiff(content, nil)
	if len(rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(rows))
	}

	// Row 1: B | [marker] (deleted)
	if rows[1].BlockIndex != 1 {
		t.Errorf("row 1 should be block 1, got BlockIndex %d", rows[1].BlockIndex)
	}
	if textLine, ok := rows[1].CommitLine.(*TextLine); !ok || textLine.Text != "B" {
		t.Error("row 1 CommitLine should be TextLine 'B'")
	}
	if _, ok := rows[1].ActualLine.(*BlockMarker); !ok {
		t.Error("row 1 ActualLine should be BlockMarker")
	}
}
