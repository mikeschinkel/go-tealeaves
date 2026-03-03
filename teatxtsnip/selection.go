package teatxtsnip

// Position represents a cursor position in the textarea
type Position struct {
	Row int
	Col int
}

// Before returns true if p comes before other in document order
func (p Position) Before(other Position) bool {
	if p.Row < other.Row {
		return true
	}
	if p.Row > other.Row {
		return false
	}
	return p.Col < other.Col
}

// After returns true if p comes after other in document order
func (p Position) After(other Position) bool {
	return other.Before(p)
}

// Equal returns true if positions are the same
func (p Position) Equal(other Position) bool {
	return p.Row == other.Row && p.Col == other.Col
}

// Selection represents a text selection range in the textarea
type Selection struct {
	Active bool     // Whether selection is active
	Start  Position // Anchor point (where selection started)
	End    Position // Moving endpoint (follows cursor)
}

// NewSelection creates an empty inactive selection
func NewSelection() Selection {
	return Selection{Active: false}
}

// Clear deactivates the selection
func (s Selection) Clear() Selection {
	s.Active = false
	s.Start = Position{}
	s.End = Position{}
	return s
}

// Begin starts a new selection at the given position
func (s Selection) Begin(pos Position) Selection {
	s.Active = true
	s.Start = pos
	s.End = pos
	return s
}

// Extend extends the selection to the given position
func (s Selection) Extend(pos Position) Selection {
	if !s.Active {
		return s.Begin(pos)
	}
	s.End = pos
	return s
}

// Normalized returns the selection with Start before End
// This is useful for iteration over the selected range
func (s Selection) Normalized() (start, end Position) {
	if s.Start.Before(s.End) || s.Start.Equal(s.End) {
		return s.Start, s.End
	}
	return s.End, s.Start
}

// Contains returns true if the given position is within the selection
func (s Selection) Contains(pos Position) bool {
	if !s.Active {
		return false
	}

	start, end := s.Normalized()

	// Before start
	if pos.Before(start) {
		return false
	}

	// After end
	if pos.After(end) {
		return false
	}

	return true
}

// IsEmpty returns true if the selection has zero width (start == end)
func (s Selection) IsEmpty() bool {
	return !s.Active || s.Start.Equal(s.End)
}

// SelectAll creates a selection spanning the entire content
func SelectAll(lines []string) Selection {
	if len(lines) == 0 {
		return NewSelection()
	}

	lastRow := len(lines) - 1
	lastCol := len([]rune(lines[lastRow]))

	return Selection{
		Active: true,
		Start:  Position{Row: 0, Col: 0},
		End:    Position{Row: lastRow, Col: lastCol},
	}
}
