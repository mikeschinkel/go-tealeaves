# Charm V1 Table (bubble-table) Problems

Issues discovered during teagrid development that should be fixed in the
Charm V2 rewrite. Each item violates the principle of least surprise,
makes examples harder to write, or has an outright bug.

---

## 1. Default text alignment is right-justified

**File:** `model.go:136`

```go
baseStyle: lipgloss.NewStyle().Align(lipgloss.Right),
```

The default `baseStyle` aligns all cell content to the right. Most table
UIs default to left-aligned text (right-align is typically reserved for
numeric columns). Every user who wants a normal-looking table must add
`WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left))`.

**Fix:** Default to `Align(lipgloss.Left)`. Offer per-column alignment
via `Column.WithStyle()` for numeric columns that need right-alignment.

---

## 2. No cell padding support

**Files:** `row.go:57-58`, `header.go:16-22`

Cell content is rendered flush against the border characters with zero
padding. There is no `WithCellPadding(left, right)` option.

`WithBaseStyle` accepts lipgloss `PaddingLeft`/`PaddingRight`, but these
do not work correctly because:

1. Cell content is pre-truncated to `column.width` by `limitStr()` before
   the style is applied (`row.go:114`).
2. Padding then makes each cell wider than `column.width`, breaking the
   table's width calculations in `dimensions.go`.
3. The footer inherits `baseStyle` including padding, which reduces its
   usable content area and can cause the footer text to wrap to a second
   line.

The only workaround is to manually prefix data values and column titles
with spaces (e.g. `" Title"`, `" " + value`), which is fragile and
pollutes the data layer.

**Fix:** Add a dedicated `WithCellPadding(left, right int)` option that
is accounted for in `limitStr`, `recalculateWidth`, column width
calculations, and is excluded from the footer rendering.

---

## 3. Footer inherits baseStyle alignment, breaking filter/pagination layout

**File:** `footer.go:21` (before our fix), `border.go:110-111`

The footer style is built as:

```go
styleFooter := m.baseStyle.Inherit(m.border.styleFooter)
```

The border's `styleFooter` sets `Align(lipgloss.Right)`, but
`baseStyle.Inherit()` only applies properties not already set on the
receiver. So if `baseStyle` has `Align(lipgloss.Left)`, the footer
becomes left-aligned — including the pagination indicator (`1/1`), which
should always be right-aligned.

Conversely, if `baseStyle` keeps the default `Align(lipgloss.Right)`,
the filter input (`/searchterm`) is right-aligned, which no text-mode
application does (vim, less, tmux all left-align search input).

**The real problem:** filter input and pagination are joined into a single
string and rendered with one alignment. They need independent alignment:
filter input left, pagination right.

**Fix:** Render filter and pagination as separate layout regions within
the footer. We applied a workaround that manually space-fills between
them, but a proper implementation should use a two-zone footer layout.

---

## 4. Row highlight is barely visible

**File:** `model.go:15`

```go
defaultHighlightStyle = lipgloss.NewStyle().Background(lipgloss.Color("#334"))
```

The default highlight color `#334` is extremely close to typical terminal
background colors (especially dark themes). On most terminals, the
highlighted row is nearly indistinguishable from non-highlighted rows.

**Fix:** Use a more visible default, e.g. `#445` or `#336`, or use
reverse video. At minimum, the default should be clearly visible on both
common dark and light terminal themes.

---

## 5. No filter match highlighting in cells

**Files:** `filter.go`, `row.go:57-121`

When filtering is active, matched rows are shown but the matching
substring within each cell is not highlighted. The user types a filter
term and sees rows appear/disappear, but has no visual indication of
*where* in the row the match occurred. This makes it unclear why a
particular row matched, especially when the match is in a truncated
portion of the text.

**Fix:** After filtering, highlight the matched substring within each
filterable column's rendered text (e.g., bold or contrasting background).

---

## 6. SelectableRows adds a column with no explanation

**File:** `options.go:109-127`

`SelectableRows(true)` silently prepends a `[x]`/`[ ]` column to the
table. If the example doesn't actually use selection (calling
`SelectedRows()`), this column is confusing — it appears interactive but
does nothing visible.

The column header shows `[x]` (the *selected* text), which implies
something is selected when nothing is.

**Fix:** Consider making the select column opt-in via a separate
`WithSelectColumn(true)` rather than coupling it to `SelectableRows()`.
At minimum, the header should show a neutral label (e.g. blank or a
checkbox icon) rather than the "selected" text.

---

## 7. Column width is content-only with no margin abstraction

**Files:** `column.go:22-29`, `strlimit.go:10-26`, `dimensions.go:27-77`

Column `width` is the raw content width in characters. There is no
concept of margin or padding at the column level. The `limitStr` function
truncates content to exactly `column.width`, and `recalculateWidth` sums
column widths plus border characters.

This means:
- Users must manually account for visual padding when setting column
  widths (e.g., adding +2 to the longest data value).
- `NewFlexColumn` distributes remaining width as raw content width, so
  flex columns also have no padding.
- There is no way to set different left/right margins per column.

**Fix:** Add `Column.WithPadding(left, right int)` that is respected by
`limitStr` (reduce content area) and `recalculateWidth` (include padding
in total width). The column's rendered width would be
`left + content + right`.

---

## 8. WithFormatString only affects data cells, not headers

**Files:** `column.go:63-73`, `header.go:19`, `row.go:82-84`

`WithFormatString` sets a format string used by `fmt.Sprintf` when
rendering data cells, but headers always render via
`limitStr(column.title, column.width)`. This means you cannot use
format strings to consistently style both headers and data (e.g., adding
left padding via `" %v"` only pads data, not the header).

This is a minor issue but contributes to the padding problem (#2).

**Fix:** Either apply format strings to headers too, or (better) solve
padding at the cell/column level so format strings aren't abused for
layout.

---

## 9. baseStyle leaks into footer rendering

**File:** `footer.go:21`

The footer is rendered with `m.baseStyle.Inherit(m.border.styleFooter)`.
Any property set on `baseStyle` (alignment, padding, colors, font weight)
affects the footer, even when the intent was to style data cells only.

This means setting `PaddingLeft(1)` on baseStyle for cell content also
shrinks the footer's usable width. Setting a background color for cells
also colors the footer.

**Fix:** The footer should have its own independent style, or at minimum
should selectively inherit only border/color properties from baseStyle,
not layout properties like alignment and padding.

---

## 10. Cell cursor mode has no per-cell highlight

**Files:** `row.go:123-142`, `model.go:51-53`

The `cellCursorMode` field and `cellCursorColumnIndex` field exist on the
model, and `WithCellCursorMode(true)` is part of the API, but the
rendering in `renderRow` always applies the highlight style to the
**entire row**:

```go
} else if m.focused && highlighted {
    rowStyle = rowStyle.Inherit(m.highlightStyle)
}
return m.renderRowData(row, rowStyle, last)
```

There is no code path that applies the highlight to only the cell at
`cellCursorColumnIndex`. When using horizontal scrolling with cell cursor
mode, the user expects to see a single highlighted cell they can navigate
left/right through columns. Instead, the entire row is highlighted
identically regardless of which cell the cursor is on.

This was one of the primary motivations for forking bubble-table into
teagrid.

**Fix:** In `renderRow` (or `renderRowData`), when `cellCursorMode` is
enabled, apply `highlightStyle` only to the cell at
`cellCursorColumnIndex`, not to the entire `rowStyle`. The non-cursor
cells in the highlighted row should either use the default style or a
subdued variant.

---

## 11. Overflow indicator is heavy-handed and non-standard

**Files:** `overflow.go`, `row.go:69-71`, `row.go:210-234`,
`header.go:72-79`

When horizontal scrolling is active and columns extend beyond
`maxTotalWidth`, bubble-table renders a hard-coded `>` character in a
dedicated overflow column on **every row** (including the header). This
column consumes `overflowColWidth` (2 chars) of horizontal space.

```go
case columnKeyOverflowRight:
    cellStyle = cellStyle.Align(lipgloss.Right)
    str = ">"
```

Problems:

- The `>` on every row is visually noisy and unconventional. No standard
  TUI framework (htop, lazygit, tig, midnight commander) uses per-row
  `>` indicators. The standard convention is a scrollbar, a subtle
  gradient/fade, or simply truncation with an ellipsis.
- The indicator is not configurable. There is no way to change the
  character, hide it, or replace it with a different visual treatment.
- It steals 2 columns of width from actual content on every render.
- The same `<` indicator appears on the left when scrolled right, further
  reducing usable width.

**Fix:** Make the overflow indicator configurable:
`WithOverflowIndicator(right, left string)` or
`WithOverflowStyle(style)`. Allow disabling it entirely. Consider
replacing the default with a subtler treatment — e.g., a single-width
Unicode fade character, or no indicator at all (let the border itself
communicate the boundary).

---

## 12. No cell or row editing support

bubble-table is strictly a read-only data display component. There is no
built-in support for editing cell values or row data, and the
architecture does not facilitate adding it:

- Row data (`RowData map[string]any`) is copied into the model on
  creation and can only be replaced wholesale via `WithRows()`.
- There is no concept of an "active cell" that can receive text input
  (cell cursor mode tracks position but has no edit state).
- There is no event or hook for entering/exiting an edit mode.
- There is no integration point for embedding a text input or form
  within or alongside the table.

### Desired editing capabilities for teagrid V2

**A. Cell-level editing (spreadsheet style)**

The cell cursor (issue #10) should support an edit mode where:
- The user navigates to a cell and presses Enter (or a configurable key)
  to enter edit mode.
- The cell background becomes an inline text input for editing the
  single cell value.
- Escape cancels, Enter confirms, and the model emits a
  `CellEditedMsg` with the old/new value, row index, and column key.
- The developer can hook `WithCellValidator(func)` to accept/reject
  edits before they commit.

**B. Row-level editing (form style)**

With or without cell cursor mode, the user should be able to enter a
row edit mode that presents all editable columns as form fields. Three
layout options:

1. **Inline single-row** — the highlighted row transforms into an
   inline editable form (each cell becomes a text input within the
   same row space).
2. **Modal popup** — a multi-field form appears as a modal overlay
   (using `teamodal`) with one field per editable column, stacked
   vertically.
3. **Split pane** — the form appears in a pane split horizontally or
   vertically from the table (the table remains visible and the form
   occupies the other half).

All three modes emit a `RowEditedMsg` on confirm or `RowEditCancelledMsg`
on cancel. The developer configures which columns are editable and
can supply per-column validators.

### Dependency: teaform

Row-level editing (option B) requires a general-purpose form component
(`teaform`) that does not yet exist in go-tealeaves. This component
would provide:

- A model that takes a list of field definitions (label, key, type,
  validator, default value).
- Keyboard navigation between fields (tab/shift-tab or up/down).
- Per-field validation with inline error display.
- Submit/cancel events.
- Embeddable as a child model within teamodal, a split pane, or
  inline within teagrid.

`teaform` should be planned and built as a standalone module before
integrating it into teagrid's row editing modes.

---

## 13. No sort key / display value separation

**Files:** `sort.go:104-121`, `sort.go:123-131`

Sorting always uses the raw `row.Data[column]` value — the same value
that is rendered in the cell. There is no way to provide a separate sort
key (e.g., a numeric timestamp or integer) while displaying a
human-friendly string.

```go
func (s *sortableTable) extractString(i int, column string) string {
    iData, exists := s.rows[i].Data[column]
    // ...
    return fmt.Sprintf("%v", iData)
}
```

This means:

- A column showing `"47d 12h"` (uptime) cannot sort numerically by
  total seconds without a workaround.
- A column showing `"Mar 1, 2026"` cannot sort by an underlying
  `time.Time` or Unix timestamp.
- The only workaround is to store a numeric value in `RowData` and use
  `WithFormatString` or `StyledCell` to control display, but this is
  fragile and not well-documented.

**Fix:** Add a `SortValue` field to `RowData` entries (or to
`StyledCell`) that, when present, is used by `extractString` and
`extractNumber` instead of the display value. Example API:

```go
teagrid.NewRow(teagrid.RowData{
    "uptime": teagrid.StyledCell{
        Data:      "47d 12h",
        SortValue: 4104000, // seconds
    },
})
```

---

## 14. No rich / mixed coloration within a cell

**Files:** `cell.go:10-21`, `row.go:91-107`

`StyledCell` applies a single `lipgloss.Style` to the **entire** cell
content. There is no way to render mixed styles within one cell — for
example, a status column showing "running" in green and "degraded" in
yellow **within the same table**, or a cell containing `key: value`
where `key` is bold and `value` is normal.

The rendering pipeline converts cell data to a plain string via
`fmt.Sprintf("%v", entry.Data)` and then applies one uniform style:

```go
case StyledCell:
    str = fmt.Sprintf(fmtString, entry.Data)
    cellStyle = entry.Style.Inherit(cellStyle)
```

Per-cell styling via `StyledCell` or `StyledCellFunc` can vary style
**between** cells (e.g., color one cell red, another green), but cannot
produce mixed formatting **within** a single cell.

Workarounds are limited:

- Pre-rendering ANSI escape sequences into the data string is fragile
  and breaks `limitStr` truncation, `lipgloss.Width` measurements, and
  style inheritance.
- Using `StyledCellFunc` to conditionally color entire cells based on
  value works for simple cases (status badges) but still cannot mix
  styles within a cell.

### Desired capabilities for teagrid V2

1. **Rich cell content** — Allow cell `Data` to be a pre-styled
   `string` containing lipgloss-rendered segments that the table
   respects during width measurement and truncation.
2. **Inline markup** — Provide a lightweight API for marking up
   substrings within a cell, e.g.:
   ```go
   teagrid.RichText(
       teagrid.Span("running", lipgloss.NewStyle().Foreground(lipgloss.Color("10"))),
   )
   ```
3. **ANSI-aware truncation** — `limitStr` (or its V2 equivalent) must
   measure visible width excluding ANSI sequences and truncate without
   splitting escape codes.
4. **ANSI-aware width measurement** — Column auto-sizing and flex
   calculations must use visible width, not byte/rune length.

---

## 15. No per-region border customization

**Files:** `border.go:6-51`, `border.go:56-94`, `border.go:102-115`

The `Border` struct is a flat set of 14 character fields applied
uniformly to the entire table. There is no way to use different border
styles for different regions — for example:

- Double-line outer border (`╔═╗║╚╝`) with single-line column dividers
  (`│`).
- Heavy header separator (`━`) with thin row dividers (`─`).
- Prominent footer border with subtle inner borders.
- Different junction characters where the header separator meets the
  outer border vs where inner dividers meet each other.

The two built-in presets (`borderDefault` with heavy box-drawing and
`borderRounded` with thin rounded) demonstrate the limitation: each
is internally consistent but there is no way to combine elements from
both, or to introduce a third weight for specific regions.

The `Border(border Border)` method (`border.go:336-342`) accepts a
custom struct, but since `generateStyles()` mechanically distributes
the same characters across all regions, fine-grained control is
impossible without replacing the entire style generation pipeline.

**Fix for teagrid V2:** Restructure borders into logical regions:

```go
type BorderConfig struct {
    Outer    BorderChars  // top, bottom, left, right, corners
    Header   BorderChars  // separator below header row
    Inner    BorderChars  // column dividers and row separators
    Footer   BorderChars  // separator above footer
}
```

Each region should independently specify its characters and whether it
is visible. This enables mixed-weight borders, borderless inner
regions with a visible outer frame, and per-region styling (color,
bold) on border characters.

---

## 16. Header cannot be fully removed

**Files:** `view.go:30-36`, `options.go:405-413`

`WithHeaderVisibility(false)` hides the header **text** but not the top
border. When the header is hidden and rows exist, `View()` still
renders the header, extracts the first line (the top border), and
includes it:

```go
} else if numRows > 0 || padding > 0 {
    split := strings.SplitN(headers, "\n", 2)
    rowStrs = append(rowStrs, split[0])
}
```

This means there is no way to render a table without its top border
frame. A developer who wants a minimal, borderless list of rows
cannot achieve it.

**Fix for teagrid V2:** Decouple header visibility from top border
rendering. Provide independent controls:

- `WithHeaderVisibility(bool)` — show/hide header text row.
- `WithBorderVisibility(bool)` or per-region border visibility — show/
  hide the outer frame, header separator, inner dividers, and footer
  separator independently.

A fully borderless table (data rows only with no chrome) should be a
supported configuration.

---

## 17. No borderless or minimal-chrome mode

**Files:** `border.go`, `view.go`

There is no `NoBorder()`, `Borderless()`, or equivalent option. Every
table renders with full box-drawing borders. The only way to
approximate a borderless table is to supply a custom `Border` struct
with all characters set to spaces, but this still consumes width for
the invisible border characters and leaves ghost spacing.

Common use cases for borderless or minimal-chrome tables:

- Embedding a table within a larger TUI layout where outer borders
  conflict with the parent container's borders.
- Lightweight data lists where borders add visual noise (e.g., a
  simple key-value display, a log viewer).
- Tables that use only horizontal rules (no vertical dividers) for a
  cleaner look.

**Fix for teagrid V2:** Support a `Borderless()` option that:

1. Removes all border characters and their associated width from
   layout calculations.
2. Optionally preserves the header separator as a horizontal rule
   (`WithHeaderSeparator(bool)`).
3. Optionally preserves column dividers without outer borders.

These should compose with the per-region border config from issue #15.
