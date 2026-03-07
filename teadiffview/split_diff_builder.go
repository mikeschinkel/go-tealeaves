package teadiffview

import (
	"fmt"
	"log/slog"

	"github.com/mikeschinkel/go-diffutils"
)

// buildSplitPaneDiff builds split pane diff rows from DiffContent.
// Pre-allocates exact row count and fills array by index.
func buildSplitPaneDiff(content *diffutils.DiffContent, logger *slog.Logger) []SplitPaneRow {
	if content == nil {
		return nil
	}

	// Calculate total rows needed for split pane.
	// We need enough rows to show:
	// - All lines from NewLines
	// - PLUS additional rows for deleted blocks (shown on old side with marker on new side)
	deletedLinesCount := 0
	for _, change := range content.Changes {
		if change.Type == diffutils.LinesDeleted {
			deletedLinesCount += change.OldRange.Count
		}
	}
	maxLines := len(content.NewLines) + deletedLinesCount
	if maxLines == 0 {
		maxLines = len(content.OldLines) // handle all-deleted case
	}
	rows := make([]SplitPaneRow, maxLines)

	if logger != nil {
		logger.Info("Row calculation",
			"oldLines", len(content.OldLines),
			"newLines", len(content.NewLines),
			"deletedCount", deletedLinesCount,
			"totalRows", maxLines,
			"changes", len(content.Changes))
	}

	oldIdx := 0  // 0-based index into OldLines
	newIdx := 0  // 0-based index into NewLines
	changeIdx := 0 // Current Change being processed

	for rowIdx := 0; rowIdx < maxLines; rowIdx++ {
		// Check if we're at the start of a change block
		if changeIdx < len(content.Changes) {
			change := content.Changes[changeIdx]
			oldBlockStart := change.OldRange.Start - 1
			newBlockStart := change.NewRange.Start - 1

			// Handle added lines block
			if change.Type == diffutils.LinesAdded && newIdx == newBlockStart {
				if newIdx >= len(content.NewLines) {
					panic(fmt.Sprintf("LinesAdded block %d: newIdx out of bounds: newIdx=%d, arrayLen=%d",
						changeIdx, newIdx, len(content.NewLines)))
				}
				if newIdx+change.NewRange.Count > len(content.NewLines) {
					panic(fmt.Sprintf("LinesAdded block %d exceeds array: newIdx=%d, count=%d, arrayLen=%d",
						changeIdx, newIdx, change.NewRange.Count, len(content.NewLines)))
				}
				markerLine := oldIdx + 1
				for i := 0; i < change.NewRange.Count; i++ {
					var ref PaneLine
					if i == 0 {
						ref = NewBlockMarker(markerLine, change.NewRange.Count)
					} else {
						ref = NewPlaceholderLine(markerLine, i)
					}
					rows[rowIdx+i] = SplitPaneRow{
						CommitLine: ref,
						ActualLine: NewTextLine(newIdx+1, content.NewLines[newIdx]),
						BlockIndex: changeIdx + 1,
						LineOffset: i,
					}
					newIdx++
				}
				changeIdx++
				rowIdx += change.NewRange.Count - 1
				continue
			}

			// Handle deleted lines block
			if change.Type == diffutils.LinesDeleted && oldIdx == oldBlockStart {
				if oldIdx >= len(content.OldLines) {
					panic(fmt.Sprintf("LinesDeleted block %d: oldIdx out of bounds: oldIdx=%d, arrayLen=%d",
						changeIdx, oldIdx, len(content.OldLines)))
				}
				if oldIdx+change.OldRange.Count > len(content.OldLines) {
					panic(fmt.Sprintf("LinesDeleted block %d exceeds array: oldIdx=%d, count=%d, arrayLen=%d",
						changeIdx, oldIdx, change.OldRange.Count, len(content.OldLines)))
				}
				markerLine := newIdx + 1
				for i := 0; i < change.OldRange.Count; i++ {
					var ref PaneLine
					if i == 0 {
						ref = NewBlockMarker(markerLine, change.OldRange.Count)
					} else {
						ref = NewPlaceholderLine(markerLine, i)
					}
					rows[rowIdx+i] = SplitPaneRow{
						CommitLine: NewTextLine(oldIdx+1, content.OldLines[oldIdx]),
						ActualLine: ref,
						BlockIndex: changeIdx + 1,
						LineOffset: i,
					}
					oldIdx++
				}
				changeIdx++
				rowIdx += change.OldRange.Count - 1
				continue
			}
		}

		// Unchanged line — both sides advance together
		if oldIdx < len(content.OldLines) && newIdx < len(content.NewLines) {
			rows[rowIdx] = SplitPaneRow{
				CommitLine: NewTextLine(oldIdx+1, content.OldLines[oldIdx]),
				ActualLine: NewTextLine(newIdx+1, content.NewLines[newIdx]),
				BlockIndex: 0,
			}
			oldIdx++
			newIdx++
			continue
		}

		// One side exhausted — fill remaining with blanks
		if oldIdx < len(content.OldLines) {
			rows[rowIdx] = SplitPaneRow{
				CommitLine: NewTextLine(oldIdx+1, content.OldLines[oldIdx]),
				ActualLine: NewTextLine(-1, ""),
				BlockIndex: 0,
			}
			oldIdx++
			continue
		}

		if newIdx < len(content.NewLines) {
			rows[rowIdx] = SplitPaneRow{
				CommitLine: NewTextLine(-1, ""),
				ActualLine: NewTextLine(newIdx+1, content.NewLines[newIdx]),
				BlockIndex: 0,
			}
			newIdx++
			continue
		}

		// Both exhausted
		rows[rowIdx] = SplitPaneRow{
			CommitLine: NewTextLine(-1, ""),
			ActualLine: NewTextLine(-1, ""),
			BlockIndex: 0,
		}
	}

	return rows
}
