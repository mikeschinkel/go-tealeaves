# How to trigger CBT sequences in Bubble Tea via ultraviolet

**The short answer: render two rows where both change every frame, with the changing cell on the upper row far to the right and the changing cell on the lower row at a tab-stop-aligned column (8, 16, 24, 32…).** The renderer's left-to-right, top-to-bottom cell scan means backward horizontal movement only happens when transitioning between rows — specifically when the cursor sits far rightward after writing a dirty cell on row N and must jump left to the first dirty cell on row N+1. When that leftward distance crosses multiple 8-column tab stops, CBT costs fewer bytes than every alternative, and ultraviolet emits it.

---

## The render loop always advances the cursor forward within a row

Ultraviolet's `render()` method (in `terminal_renderer.go`) implements an ncurses-inspired cell-based diff. It maintains two screen buffers — "old" (what the terminal currently displays) and "new" (what `View()` just produced) — then iterates in row-major order (top-to-bottom, left-to-right) comparing each cell. When it finds a **dirty cell** (new ≠ old), it calls `relativeCursorMove()` to reposition the cursor, then writes the cell content. The cursor then advances rightward by the cell's width.

Because the scan always proceeds left-to-right within a row, **the cursor only ever moves forward horizontally within a single line**. Skip-ahead jumps across unchanged cells in the middle of a row are always forward moves (CUF or hard tabs). Backward horizontal movement is physically impossible within one row's scan.

The critical moment for backward movement occurs at the **row boundary**: after writing the last dirty cell on row N (cursor now at column X), the renderer moves to the first dirty cell on row N+1 at column Y. When **Y < X**, the horizontal component requires backward movement. This is the only scenario that feeds into `relativeCursorMove()`'s backward path where CBT lives.

## How relativeCursorMove picks CBT over alternatives

When the target column (tx) is less than the current column (fx), `relativeCursorMove()` evaluates several strategies for the horizontal component and picks the one with the lowest byte cost. Following the ncurses `mvcur()` pattern that ultraviolet explicitly implements, the candidates are:

- **CUB (Cursor Backward)**: `\x1b[nD` — costs 3 bytes for n=1, 4 bytes for n=2–9, 5 bytes for n=10–99
- **CBT (Cursor Backward Tab)**: `\x1b[nZ` — jumps backward n tab stops. Costs 3 bytes for n=1, 4 bytes for n=2–9. May require residual CUB if the target isn't exactly on a tab stop
- **CR + CUF**: Carriage return (`\r`, 1 byte) to column 0, then CUF forward to target. Costs 1 + CUF bytes
- **CHA (Cursor Horizontal Absolute)**: `\x1b[colG` (1-indexed column). Costs 4 bytes for columns 1–9, 5 bytes for 10–99
- **CUP (Cursor Position)**: `\x1b[row;colH` — absolute positioning. Most expensive at 6+ bytes

Per the user's confirmed conditions, CBT enters the competition only when `useTabs` is true, `s.tabs` is configured (bubbletea always sets default 8-column stops), `capCBT` is in the terminal's capabilities, and at least one backward tab stop lies between the current and target positions.

**CBT wins the cost comparison in a specific sweet spot.** When the target column aligns exactly with a tab stop and the backward distance crosses 2+ tab stops, CBT's compact encoding beats every alternative. Here's a concrete comparison for cursor at column 34, target at column 16:

| Strategy | Sequence | Bytes |
|---|---|---|
| CUB(18) | `\x1b[18D` | 5 |
| **CBT(3)** | **`\x1b[3Z`** (34→32→24→16) | **4** |
| CR+CUF(16) | `\r\x1b[16C` | 6 |
| CHA(17) | `\x1b[17G` | 5 |

CBT wins by 1 byte. That single byte matters because ultraviolet is optimized for minimal data transfer, especially over SSH.

## The tab-stop alignment requirement is critical

CBT lands on tab stop boundaries, not arbitrary columns. With default 8-column stops at 0, 8, 16, 24, 32, 40, 48…, if the target column isn't exactly a tab stop, CBT overshoots and requires a residual CUF or CUB correction, inflating its total cost and usually losing to simpler alternatives.

**Column 0 is a special case**: CR (`\r`) reaches column 0 in just 1 byte, so CBT can never beat CR for a column-0 target. The productive CBT targets are columns **8, 16, 24, 32, 40, 48** and so on.

For higher tab-stop columns (16+), CBT's advantage over CHA grows because **CHA's byte cost increases with column number** (more digits in the parameter) while **CBT's cost increases with the number of tab stops crossed** (typically 1–4, keeping the parameter small). Column 16 requires `\x1b[17G` (5 bytes) via CHA but only `\x1b[nZ` (4 bytes) via CBT from a few tab stops away.

## What Bubble Tea View() output produces this pattern

The recipe requires exactly two conditions between consecutive frames:

1. **A dirty cell at a high column on row N** — something that changes on the right side of a line
2. **A dirty cell at a tab-stop column on row N+1** — something that changes exactly at column 8, 16, 24, etc. on the next line

Both must change every frame. If only one row changes, the renderer never needs to transition between rows during that render pass, so backward movement never occurs.

Here is a minimal `View()` function that reliably triggers CBT:

```go
func (m model) View() tea.View {
    // Row 0: a counter at column 30 changes every frame.
    // The 30 chars of padding stay constant; only the last 4 chars are dirty.
    row0 := fmt.Sprintf("  Timer value:                %04d", m.count)

    // Row 1: a spinner at column 16 changes every frame.
    // The 16 spaces of padding stay constant; only column 16 is dirty.
    spin := [4]byte{'|', '/', '-', '\\'}
    row1 := fmt.Sprintf("                %c processing...", spin[m.count%4])

    return tea.NewView(row0 + "\n" + row1)
}
```

On each tick, `m.count` increments. The renderer sees:

- **Row 0**: only the digits at columns 30–33 are dirty (e.g., `0001` → `0002`). After writing the last dirty cell, cursor lands at **column 34**.
- **Row 1**: only column 16 is dirty (spinner character rotates). The renderer must move cursor from **(0, 34)** to **(1, 16)**.

The vertical component is CUD(1) = `\x1b[B` (3 bytes). The horizontal component moves backward from column 34 to column 16 — crossing tab stops at 32 and 24. **CBT(3)** goes 34→32→24→16 in exactly `\x1b[3Z` (4 bytes), beating CUB(18) at 5 bytes, CHA(17) at 5 bytes, and CR+CUF(16) at 6 bytes.

## A complete minimal reproduction app

```go
package main

import (
    "fmt"
    "time"

    tea "charm.land/bubbletea/v2"
)

type tickMsg time.Time

type model struct{ count int }

func (m model) Init() tea.Cmd {
    return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg.(type) {
    case tea.KeyMsg:
        return m, tea.Quit
    case tickMsg:
        m.count++
        return m, tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
            return tickMsg(t)
        })
    }
    return m, nil
}

func (m model) View() tea.View {
    // Changing content at column 30 on row 0
    row0 := fmt.Sprintf("  Timer value:                %04d", m.count)
    // Changing content at column 16 (tab stop) on row 1
    spin := [4]byte{'|', '/', '-', '\\'}
    row1 := fmt.Sprintf("                %c processing...", spin[m.count%4])
    return tea.NewView(row0 + "\n" + row1)
}

func main() {
    p := tea.NewProgram(model{}, tea.WithAltScreen())
    if _, err := p.Run(); err != nil {
        fmt.Println(err)
    }
}
```

**Why this works**: Every 500ms, both the counter (column 30+) and spinner (column 16) produce dirty cells. The renderer writes the counter on row 0, then must jump backward across 3 tab stops to reach the spinner on row 1. That backward jump of 18 columns from column 34 to column 16 is cheaper via CBT(3) than any alternative.

## Variations that also trigger CBT

The pattern generalizes to any layout where **right-side content on one row and left-side content on the next row both change simultaneously**:

- **A dashboard** with a timestamp updating on the right of row N and a status indicator on the left of row N+1. The wider the horizontal gap between the two changing regions, the more likely CBT wins.
- **A table** where a value in a rightward column changes on one row while a value in a leftward column changes on the row below. Column-aligned tabular layouts are especially effective because table columns at multiples of 8 naturally land on tab stops.
- **Any two Lip Gloss–styled components** arranged vertically where both re-render each frame. If one component updates content near the right edge and the component below it updates content near its left edge, the row transition triggers backward movement.

The essential conditions are: both rows must produce dirty cells, the upper row's dirty cells must be to the right of the lower row's dirty cells, the gap must cross at least one tab stop, and the lower row's target should align with a tab stop (column 8, 16, 24, 32…) for CBT to beat CHA and CUB on byte cost.

## Conclusion

CBT emission is a **byte-count optimization** in ultraviolet's ncurses-derived diff renderer. It fires exclusively during **inter-row transitions** when the cursor must jump leftward across tab-stop boundaries. The trigger is not exotic content — it's a common layout pattern of two vertically stacked regions that both update simultaneously at different horizontal offsets. The minimal reproduction is two lines: a changing counter on the right of line 1 and a changing indicator at a tab-stop column on line 2. Targeting column 16 or higher ensures CBT reliably beats all competing strategies by at least 1 byte.