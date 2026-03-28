package tealayout

import "strings"

// Alignment is a bitmap controlling how content is positioned within a pane
// when the rendered content is smaller than the allocated space.
type Alignment int8

// Horizontal alignment bits.
const (
	Left   Alignment = 1 << iota // 0b0000_0001
	Center                       // 0b0000_0010
	Right                        // 0b0000_0100
)

// Vertical alignment bits.
const (
	Top    Alignment = 1 << (iota + 3) // 0b0000_1000
	Middle                              // 0b0001_0000
	Bottom                              // 0b0010_0000
)

// Composed alignment constants.
const (
	TopLeft      = Top | Left
	TopCenter    = Top | Center
	TopRight     = Top | Right
	MiddleLeft   = Middle | Left
	MiddleCenter = Middle | Center
	MiddleRight  = Middle | Right
	BottomLeft   = Bottom | Left
	BottomCenter = Bottom | Center
	BottomRight  = Bottom | Right
)

const (
	hMask Alignment = Left | Center | Right
	vMask Alignment = Top | Middle | Bottom
)

// mergeAlignment composably merges new alignment bits into old.
// If new has horizontal bits, they replace old horizontal bits (and vice versa).
// Bits on the other axis are preserved.
func mergeAlignment(old, new Alignment) Alignment {
	if new&hMask != 0 {
		old = old &^ hMask // clear old horizontal
	}
	if new&vMask != 0 {
		old = old &^ vMask // clear old vertical
	}
	return old | new
}

// alignContent positions content within allocW x allocH space according to align.
func alignContent(content string, allocW, allocH int, align Alignment) string {
	if content == "" || allocW <= 0 || allocH <= 0 {
		return content
	}

	lines := strings.Split(content, "\n")

	// Horizontal alignment: pad each line
	if align&hMask != 0 {
		for i, line := range lines {
			lineW := visibleWidth(line)
			pad := allocW - lineW
			if pad <= 0 {
				continue
			}
			switch {
			case align&Right != 0:
				lines[i] = strings.Repeat(" ", pad) + line
			case align&Center != 0:
				leftPad := pad / 2
				lines[i] = strings.Repeat(" ", leftPad) + line
			// Left: no padding needed (default)
			}
		}
	}

	// Vertical alignment: prepend/append blank lines
	contentH := len(lines)
	vPad := allocH - contentH
	if vPad > 0 {
		blank := strings.Repeat(" ", allocW)
		switch {
		case align&Bottom != 0:
			padLines := make([]string, vPad)
			for i := range padLines {
				padLines[i] = blank
			}
			lines = append(padLines, lines...)
		case align&Middle != 0:
			topPad := vPad / 2
			topLines := make([]string, topPad)
			for i := range topLines {
				topLines[i] = blank
			}
			lines = append(topLines, lines...)
		// Top: no vertical padding needed (default)
		}
	}

	return strings.Join(lines, "\n")
}

// visibleWidth returns the visible width of a string, accounting for ANSI
// sequences. For now uses a simple byte-length approach; can be upgraded
// to use lipgloss.Width if needed.
func visibleWidth(s string) int {
	// Use simple len for now — lipgloss.Width handles ANSI but adds import.
	// Since content is typically already rendered by lipgloss, the width
	// should match allocW. We use len as a fallback.
	return len(s)
}
