package teamodal

// calcScrollbar computes scrollbar thumb position and height for a scrollable list.
// Maps offset ∈ [0, totalItems-maxVisible] → pos ∈ [0, maxVisible-height].
func calcScrollbar(offset, maxVisible, totalItems int) (pos, height int) {
	height = maxVisible * maxVisible / totalItems
	if height < 1 {
		height = 1
	}
	maxOffset := totalItems - maxVisible
	maxPos := maxVisible - height
	if maxOffset > 0 && maxPos > 0 {
		pos = offset * maxPos / maxOffset
	}
	return pos, height
}
