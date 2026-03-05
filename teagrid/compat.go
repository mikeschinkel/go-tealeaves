package teagrid

import "charm.land/lipgloss/v2"

// --- v0.1.0 backward compatibility aliases ---

// StyledCell is an alias for CellValue for backward compatibility.
type StyledCell = CellValue

// StyledCellFunc is an alias for CellStyleFunc for backward compatibility.
type StyledCellFunc = CellStyleFunc

// StyledCellFuncInput is an alias for CellStyleInput for backward compatibility.
type StyledCellFuncInput = CellStyleInput

// NewStyledCell creates a StyledCell (CellValue) for backward compatibility.
func NewStyledCell(data any, style lipgloss.Style) StyledCell {
	return NewCellValue(data, style)
}

// NewStyledCellWithStyleFunc creates a styled cell with a dynamic style function.
func NewStyledCellWithStyleFunc(data any, fn StyledCellFunc) StyledCell {
	return NewCellValueWithStyleFunc(data, fn)
}

// BorderDefault applies the default heavy border to the model.
// Deprecated: Use WithBorder(BorderDefault()) instead.
func (m GridModel) BorderDefault() GridModel {
	return m.WithBorder(teagridBorderDefault())
}

// BorderRounded applies the rounded border to the model.
// Deprecated: Use WithBorder(BorderRounded()) instead.
func (m GridModel) BorderRounded() GridModel {
	return m.WithBorder(teagridBorderRounded())
}

// SetBorder applies a custom border to the model.
// Deprecated: Use WithBorder() instead.
func (m GridModel) SetBorder(border BorderConfig) GridModel {
	return m.WithBorder(border)
}

// teagridBorderDefault returns the default border config.
// Named differently to avoid collision with the standalone BorderDefault function.
func teagridBorderDefault() BorderConfig {
	return BorderConfig{
		Outer:  BorderRegion{Visible: true},
		Header: BorderRegion{Visible: true},
		Inner:  BorderRegion{Visible: true},
		Footer: BorderRegion{Visible: true},
		Chars:  charsDefault,
	}
}

// teagridBorderRounded returns the rounded border config.
func teagridBorderRounded() BorderConfig {
	return BorderRounded()
}

// HeaderStyle sets the header text style.
// Deprecated: Use WithHeaderStyle() instead.
func (m GridModel) HeaderStyle(style lipgloss.Style) GridModel {
	return m.WithHeaderStyle(style)
}

// HighlightStyle sets the highlighted row style.
// Deprecated: Use WithHighlightStyle() instead.
func (m GridModel) HighlightStyle(style lipgloss.Style) GridModel {
	return m.WithHighlightStyle(style)
}

// SelectableRows enables/disables row selection.
// Deprecated: Use WithSelectableRows() instead.
func (m GridModel) SelectableRows(selectable bool) GridModel {
	return m.WithSelectableRows(selectable)
}

// WithTargetWidth sets the viewport width for flex columns.
// Deprecated: Use SetSize() instead.
func (m GridModel) WithTargetWidth(width int) GridModel {
	m.viewportWidth = width
	m.recalculateWidth()
	return m
}

// WithMaxTotalWidth sets the maximum width for overflow/scrolling.
// Deprecated: Use SetSize() instead.
func (m GridModel) WithMaxTotalWidth(width int) GridModel {
	m.viewportWidth = width
	m.recalculateWidth()
	return m
}

// WithGlobalMetadata sets grid-level metadata.
// Deprecated: Use WithMetadata() instead.
func (m GridModel) WithGlobalMetadata(metadata map[string]any) GridModel {
	return m.WithMetadata(metadata)
}

// WithFilterFunc sets a custom filter function.
// Backward-compatible alias.
func (m GridModel) WithCustomFilterFunc(fn FilterFunc) GridModel {
	return m.WithFilterFunc(fn)
}

// WithMissingDataIndicatorStyled sets a styled missing data indicator.
// Deprecated: Use WithMissingDataIndicator with a CellValue.
func (m GridModel) WithMissingDataIndicatorStyled(styled StyledCell) GridModel {
	m.missingDataIndicator = styled
	return m
}
