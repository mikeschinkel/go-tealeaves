# Bubble Tea Best Practices (Charm v2)

Hard-won lessons from building Bubble Tea UIs, updated for `charm.land/lipgloss/v2` and `charm.land/bubbletea/v2`. Read this BEFORE implementing new models or components to avoid repeating mistakes.

If you are working with v1 (`github.com/charmbracelet/lipgloss` and `github.com/charmbracelet/bubbletea`), see `BEST_PRACTICES_CHARM_V1.md` instead.

## Import Path Changes (v1 to v2)

| v1 | v2 |
|----|-----|
| `github.com/charmbracelet/bubbletea` | `charm.land/bubbletea/v2` |
| `github.com/charmbracelet/lipgloss` | `charm.land/lipgloss/v2` |
| `github.com/charmbracelet/bubbles` | `charm.land/bubbles/v2` |

Note: `github.com/charmbracelet/x/ansi` did NOT move — same import path in v2.

### Key API Renames (v1 to v2)

| v1 | v2 | Notes |
|----|-----|-------|
| `tea.KeyMsg` | `tea.KeyPressMsg` | v2 also adds `tea.KeyReleaseMsg` |
| `View() string` | `View() tea.View` | `tea.View` has `.Content`, `.AltScreen`, `.MouseMode` fields |
| `.BorderStyle(b)` | `.Border(b, sides...)` | `BorderStyle()` still exists but `Border()` is canonical v2 |

In v2, `View()` returns a `tea.View` struct. To get the rendered string, use `.View().Content`:
```go
childContent := m.childModel.View().Content
```

---

## THE GOLDEN RULE

**TRUST LIPGLOSS. DON'T SECOND-GUESS IT.**

If you find yourself with "empirical adjustment" constants to fix widths, **you have a bug in your understanding**, not in lipgloss. Stop, step back, and find the root cause.

This rule has not changed from v1. What HAS changed is what "trusting lipgloss" means for width calculations.

---

## Width/Height Semantics: THE CRITICAL V2 CHANGE

### v1 Behavior (OLD -- do not use with v2)

In lipgloss v1, `Width(n)` set the **content + padding** width. Borders were added OUTSIDE:

```
v1: Width(35) with border
Border(1) | Pad(1) | Content(33) | Pad(0) | Border(1)
|<------ Width(35) ------->|
|<------- Total rendered = 37 -------->|
```

To get a total rendered width of 37, you wrote `Width(37 - 2)` = `Width(35)`.

### v2 Behavior (CURRENT)

In lipgloss v2, `Width(n)` sets the **total rendered width**. Borders, padding, and content ALL fit inside `n`. Margins are applied OUTSIDE.

```
v2: Width(37) with border
Border(1) | Pad(1) | Content(33) | Pad(1) | Border(1)
|<------------- Width(37) ----------------->|
```

To get a total rendered width of 37, you write `Width(37)`. That's it.

### Side-by-Side Comparison

```
Goal: total rendered width of 37, with 1px border on each side, 1px padding on each side

v1:
  borderWidth := 2  // 1 left + 1 right
  style.Width(37 - borderWidth)  // Width(35)
  // Result: 35 (content+padding) + 2 (border) = 37 rendered

v2:
  style.Width(37)
  // Result: 37 rendered, period.
  // Content area = 37 - 2 (border) - 2 (padding) = 33
```

### What About Height?

Same change. `Height(n)` in v2 is total rendered height including top/bottom borders and padding.

```
v2: Height(20) with border
Border(1 top)
Pad(1 top)
Content (16 lines)
Pad(1 bottom)
Border(1 bottom)
|<--- Height(20) --->|
```

### The Migration Trap

**v1 patterns that are BUGS in v2:**

```go
// v1 CORRECT, v2 BUG -- double-subtracting border
const borderWidth = 2
style.Width(totalWidth - borderWidth)  // NO! v2 already accounts for border

// v2 CORRECT
style.Width(totalWidth)  // Border fits inside totalWidth
```

If you are migrating from v1, search your codebase for every `Width(` and `Height(` call that subtracts border sizes. Remove those subtractions.

### How Width Affects Text Wrapping Internally

Lipgloss v2 wraps text based on Width. The actual wrapping limit is computed as:

```
wrapAt = Width - borderSize - padding
       = content area
```

This means **Width must include border + padding + content** for text to wrap correctly. If you only add padding to your content width, the border eats into the wrap budget and text word-wraps prematurely.

**Anti-pattern (causes word-wrapping):**
```go
// BUG: only accounts for padding, not border
style.Width(contentWidth + padding)
// wrapAt = contentWidth + padding - border - padding = contentWidth - border
// Lines at contentWidth will WORD-WRAP because wrapAt < contentWidth!
```

**Correct pattern:**
```go
// v2: Width must include border + padding
style.Width(contentWidth + style.GetHorizontalPadding() + style.GetHorizontalBorderSize())
// wrapAt = contentWidth + padding + border - border - padding = contentWidth ✓
```

This is the converse of the "computing content area" formula. When you HAVE content and need Width:

```go
// Computing Width from known content width:
totalWidth := contentWidth + style.GetHorizontalPadding() + style.GetHorizontalBorderSize()
style.Width(totalWidth)

// Computing content width from known Width:
contentWidth := totalWidth - style.GetHorizontalPadding() - style.GetHorizontalBorderSize()
```

---

## The Border/Padding/Margin Box Model in v2

v2 follows a CSS-like box model:

```
|<------------ Margin (outside Width) ------------>|
|     |<------------ Width(n) ------------->|      |
|     | Border | Padding | Content | Padding | Border |
|     |<------------ n characters ---------->|      |
```

- **Width(n)**: Total rendered width = n. Everything inside (border + padding + content) fits in n characters.
- **Margin**: Applied OUTSIDE Width. `MarginLeft(2)` adds 2 spaces before the box. Total screen usage = margin + Width.
- **Padding**: Inside the border, reduces content area.
- **Border**: Inside Width, reduces content area.
- **Content area** = Width - horizontal border - horizontal padding.

### Computing Content Area

When you need to tell child components how much space they have (e.g., `viewport.SetWidth()`), compute the content area:

```go
const (
    borderH = 2  // 1 left + 1 right
    padH    = 2  // PaddingLeft(1) + PaddingRight(1)
)

totalWidth := 80
style := lipgloss.NewStyle().
    Width(totalWidth).
    PaddingLeft(1).
    PaddingRight(1).
    BorderStyle(lipgloss.RoundedBorder())

contentWidth := totalWidth - borderH - padH  // 80 - 2 - 2 = 76
m.viewport.SetWidth(contentWidth)
```

You can also use lipgloss's built-in helpers:

```go
// When you have NO margins (the common case):
hFrame := style.GetHorizontalFrameSize()  // margins + border + padding
contentWidth := totalWidth - hFrame

// More precise (works regardless of margins):
contentWidth := totalWidth - style.GetHorizontalPadding() - style.GetHorizontalBorderSize()
```

**Caution:** `GetHorizontalFrameSize()` returns `margins + padding + border`. Since Width does NOT include margins, the shorthand `Width - GetHorizontalFrameSize()` is only correct when margins are zero. When margins are present, use `GetHorizontalPadding() + GetHorizontalBorderSize()` explicitly.

---

## The Two-Width Problem (v2 Edition)

You still need TWO width values, but for a clearer reason than in v1:

1. **Total Width** -- For lipgloss `Width()` (the box on screen)
2. **Content Width** -- For child component `SetSize()` (what they can render into)

```go
type MyModel struct {
    paneTotalWidth   int  // For lipgloss Width()
    paneContentWidth int  // For viewport.SetWidth(), tree.SetSize(), etc.
}
```

In v2, deriving one from the other is straightforward and correct:

```go
// v2: This is the RIGHT way to derive content width
func (m MyModel) calculateLayout() MyModel {
    m.paneTotalWidth = allocatedWidth

    hFrame := m.paneStyle.GetHorizontalFrameSize()
    m.paneContentWidth = m.paneTotalWidth - hFrame

    return m
}
```

In v1, `Width()` excluded border, so "total width" and "lipgloss Width() value" were different numbers. In v2, they are the same number. This eliminates an entire class of off-by-one bugs.

### WHO HANDLES BORDERS?

Unchanged from v1. The **parent model** applies borders, not child components.

```go
// Parent model
func (m MyModel) View() tea.View {
    childContent := m.childModel.View().Content  // Child returns tea.View; get string

    // Parent wraps with border
    style := lipgloss.NewStyle().
        Width(m.paneTotalWidth).
        Border(lipgloss.RoundedBorder())
    return tea.NewView(style.Render(childContent))
}
```

Child models (teatree, viewport, etc.) have NO knowledge of borders. Parent calculates total width, derives content width, passes content width to child via `SetSize()`, and applies border when rendering.

---

## TWO APPROACHES TO SIZING

Unchanged from v1.

**Approach 1: Auto-sizing (no Width() set)**
```go
// Lipgloss auto-sizes to content
style := lipgloss.NewStyle().
    PaddingLeft(1).
    PaddingRight(1).
    BorderStyle(lipgloss.RoundedBorder())

// Visual width = content + padding + border (all auto-calculated)
```
Use for: Dynamic content like trees where you don't know the width in advance.

**Approach 2: Constrained sizing (Width() set)**
```go
// v2: Width IS the total rendered width
style := lipgloss.NewStyle().
    Width(totalWidth).  // Total rendered width = totalWidth
    PaddingLeft(1).
    BorderStyle(lipgloss.RoundedBorder())
```
Use for: Fixed-width panes like code viewers.

---

## Lipgloss Border Widths

Unchanged from v1. Borders are 1 character per side, not 2.

**Single borders (RoundedBorder, NormalBorder, etc.):**
- 1 character per side (the box-drawing character)
- 2 characters total width (left + right)

**Double borders (DoubleBorder):**
- Still 1 character per side (different Unicode characters)
- 2 characters total width

**Common mistake:** Thinking borders are 2 chars per side because they "look thick". They are single Unicode characters.

---

## Padding Calculations

Unchanged from v1. Track padding per pane type:

```go
const (
    treePaddingWidth = 2  // Left + Right
    codePaddingWidth = 1  // Left only
)
```

The difference in v2: padding reduces content area within `Width()`, rather than being added on top of `Width()` alongside border being outside. But since `Width()` already includes everything, you only need to think about padding when computing content width for child components.

---

## Style Composition

v2 style composition works the same as v1. Styles are immutable value types:

```go
baseStyle := lipgloss.NewStyle().
    BorderStyle(lipgloss.RoundedBorder()).
    BorderForeground(lipgloss.Color("62"))

// Derive from base -- does NOT mutate baseStyle
activeStyle := baseStyle.
    BorderForeground(lipgloss.Color("205"))
```

Use `Inherit()` to pull properties from another style without overriding what you have set.

---

## Immutable Struct Pattern

Unchanged from v1. Methods should update struct fields and return the model:

```go
// CORRECT - Bubble Tea pattern
func (m MyModel) calculateLayout() MyModel {
    m.contentWidth = /* calculate */
    m.totalWidth = /* calculate */
    return m
}

// WRONG - Leads to verbose code and errors
func (m MyModel) calculateLayout() (contentWidth, totalWidth int) {
    return /* calculate */, /* calculate */
}
```

---

## Equal Width Panes: The Odd Width Problem

Same principle as v1. Use dynamic padding to consume odd columns. Don't lie to child components about their size.

```go
func (m MyModel) calculateLayout() MyModel {
    m.treeRightPadding = 1  // Default

    normalTreeTotal := treeContentWidth + padLeft + 1 + borderH
    potentialRemaining := terminalWidth - normalTreeTotal

    if bothPanesVisible && potentialRemaining%2 == 1 {
        m.treeRightPadding = 2  // Extra column as padding
    }

    return m
}

func (m MyModel) View() string {
    treeStyle := lipgloss.NewStyle().
        Width(m.treeTotalWidth).  // v2: Width IS total rendered width
        PaddingLeft(1).
        PaddingRight(m.treeRightPadding).  // Dynamic!
        BorderStyle(lipgloss.RoundedBorder())

    return treeStyle.Render(m.treePane.View())
}
```

---

## Overlays: Compositing Foreground on Background

Unchanged from v1. Use ANSI-aware string operations from `github.com/charmbracelet/x/ansi`:

```go
import "github.com/charmbracelet/x/ansi"

// CORRECT - Visual width
width := ansi.StringWidth(styledText)

// WRONG - Byte length (includes ANSI codes)
width := len(styledText)
```

### OVERLAY PATTERN

Same proven pattern from teadd and teamodal. See `BEST_PRACTICES_CHARM_V1.md` for full overlay implementation code. The overlay logic is independent of lipgloss width semantics -- it operates on already-rendered strings.

### ANSI-AWARE STRING OPERATIONS

| Operation | Description |
|-----------|-------------|
| `ansi.StringWidth(s)` | Visual width (what you see on screen) |
| `ansi.Truncate(s, width, tail)` | Keep first N visual columns (cut from right) |
| `ansi.TruncateLeft(s, width, tail)` | Skip first N visual columns (cut from left) |

### REFERENCE IMPLEMENTATIONS

- **teadrpdwn/overlay_dropdown.go** - Dropdown pattern
- **teamodal/overlay_modal.go** - Modal pattern

---

## Modal Key Consumption: Preventing ESC from Propagating

Unchanged from v1. Check modal state BEFORE updating, then block key propagation if it was open.

```go
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Track if modal was open BEFORE update
    modalWasOpen := m.modal.IsOpen()

    // Update modal (it may handle the key and close itself)
    m.modal, modalCmd = m.modal.Update(msg)

    // Only pass messages to view stack if no modal was open
    if _, isKey := msg.(tea.KeyPressMsg); isKey && modalWasOpen {
        // Modal consumed the key - don't pass to view stack
    } else {
        viewStackCmd = m.updateViewStack(msg)
    }

    return m, tea.Batch(modalCmd, viewStackCmd)
}
```

**Why "was open" not "is open":** If user presses `esc` to close the modal, `Update()` sets `isOpen = false`. Checking after update would see `false` and let the key propagate. Check before update to capture the pre-close state.

---

## Emojis and Wide Characters: The Unsolvable Alignment Problem

Unchanged from v1. This is a fundamental mismatch between Unicode width tables and terminal rendering.

**For width-sensitive layouts, use ASCII characters only.**

Emojis are fine for status messages, notifications, headers -- anywhere alignment does not matter. They are NOT fine for list items that need to align, table columns, or fixed-width fields.

See `BEST_PRACTICES_CHARM_V1.md` for the full explanation of why this is fundamentally unsolvable.

---

## Debugging Width Issues: The Systematic Approach

Same systematic approach as v1, updated for v2 semantics.

### MEASURE EVERYTHING

```go
m.Logger.Info("WIDTH DEBUG calculateLayout",
    "terminalWidth", m.terminalWidth,
    "paneTotalWidth", m.paneTotalWidth,
    "paneContentWidth", m.paneContentWidth,
    "hFrame", m.paneStyle.GetHorizontalFrameSize(),
)

// In View() - measure what actually renders
rendered := style.Render(content)
actualWidth := 0
if lines := strings.Split(rendered, "\n"); len(lines) > 0 {
    actualWidth = ansi.StringWidth(lines[0])
}

m.Logger.Info("WIDTH DEBUG View",
    "calculated", m.paneTotalWidth,
    "actual", actualWidth,
    "gap", m.paneTotalWidth - actualWidth,
)
```

### v2-SPECIFIC DISCREPANCIES

| Symptom | Likely Cause |
|---------|--------------|
| Actual = Calculated - 2 | Subtracting border from Width() (v1 habit) |
| Actual = Calculated + 2 | NOT subtracting border -- but you set Width to content width, not total |
| Content truncated | Content width not accounting for border+padding inside Width() |
| Too much whitespace | Content width double-subtracting border |
| Long lines word-wrap inside border | Width only includes padding, not border -- `Width(content + padding)` should be `Width(content + padding + border)` |

### THE v2 DEBUGGING CHECKLIST

If actual != calculated:
1. Are you subtracting border from `Width()`? **Stop. In v2, Width IS total.**
2. Are you computing content width as `Width - border - padding`? **Good.**
3. Are you passing content width (not total width) to child `SetSize()`? **Good.**
4. Is `GetHorizontalFrameSize()` returning what you expect? **Log it.**

---

## LOGGING: NEVER USE fmt.Printf/fmt.Fprintf IN TUI CODE

Unchanged from v1. This corrupts the TUI display.

```go
// NEVER DO THIS
fmt.Printf("DEBUG: value=%d\n", value)
fmt.Fprintf(os.Stderr, "DEBUG: %s\n", msg)
log.Printf("DEBUG: %v", thing)

// ALWAYS DO THIS
m.Logger.Debug("debug message", "value", value)
m.Logger.Info("info message", "key", value)
```

Every model that needs logging should have a `Logger *slog.Logger` field, injected from the parent. See `BEST_PRACTICES_CHARM_V1.md` for the full logger injection pattern.

---

## Component Update() Return Types

Unchanged from v1. Some Bubble Tea components return their concrete type from `Update()`, not `tea.Model`.

```go
// viewport.Model.Update() returns viewport.Model directly
m.viewport, cmd = m.viewport.Update(msg)  // Correct
```

Check each component's method signature. The compiler will tell you.

---

## Viewport Horizontal Scrolling: Use Native Support

Unchanged from v1. Use the viewport's built-in horizontal scrolling:

```go
vp := viewport.New(width, height)
vp.SetHorizontalStep(4)  // Must be > 0 to enable
```

Do NOT manually track horizontal offset and apply `TruncateLeft` post-render. See `BEST_PRACTICES_CHARM_V1.md` for the full explanation of why manual scrolling breaks.

---

## Tab Characters in Viewport Content

Unchanged from v1. Replace tabs with spaces before setting viewport content:

```go
diff = strings.ReplaceAll(diff, "\t", "    ")
m.diffViewport.SetContent(diff)
```

---

## Calculation Order and Best Practices

Same principles as v1, simplified by v2 semantics:

1. **Calculate auto-sized pane width first** (e.g., tree from `LayoutWidth()`)
2. **Determine dynamic properties** (e.g., padding based on odd/even)
3. **Calculate auto-sized pane total** (content + dynamic padding + border)
4. **Calculate remaining width** for constrained panes
5. **Distribute remaining width** (guaranteed even if step 2 handled it)

```go
// 1. Auto-sized content width
treeContentWidth := tree.LayoutWidth()

// 2-3. Dynamic padding and total
treeRightPadding := 1
treeTotalWidth := treeContentWidth + padLeft + treeRightPadding + borderH
potentialRemaining := terminalWidth - treeTotalWidth
if potentialRemaining%2 == 1 {
    treeRightPadding = 2
    treeTotalWidth++
}

// 4. Remaining for constrained panes
remainingWidth := terminalWidth - treeTotalWidth

// 5. Distribute
halfWidth := remainingWidth / 2
pane2TotalWidth := halfWidth
pane3TotalWidth := halfWidth

// v2: Content width = total - frame
hFrame := paneStyle.GetHorizontalFrameSize()
pane2ContentWidth := pane2TotalWidth - hFrame
pane3ContentWidth := pane3TotalWidth - hFrame
```

---

## When to Recalculate Layout

Unchanged from v1. Recalculate on:
- `tea.WindowSizeMsg` -- Terminal resized
- Toggling pane visibility
- Changing content that affects natural width (like tree)

```go
case tea.WindowSizeMsg:
    m.terminalWidth = msg.Width
    m.terminalHeight = msg.Height
    m = m.recalculateLayout()
```

---

## Migration Checklist: v1 to v2

When migrating code from lipgloss v1 to v2:

- [ ] Update import paths (`github.com/charmbracelet/*` to `charm.land/*/v2`), including bubbles
- [ ] Rename `tea.KeyMsg` to `tea.KeyPressMsg` in all type switches/assertions
- [ ] Update `View()` return type from `string` to `tea.View`; use `tea.NewView()` to construct
- [ ] Update `.View()` call sites on child models: `.View()` → `.View().Content` to get string
- [ ] Remove ALL border subtraction from `Width()` calls -- `Width(total - border)` becomes `Width(total)`
- [ ] Remove ALL border subtraction from `Height()` calls -- same change
- [ ] When computing Width from content: `Width(content + padding + border)`, NOT `Width(content + padding)`
- [ ] Update content width derivation: `contentWidth = width - GetHorizontalPadding() - GetHorizontalBorderSize()`
- [ ] Verify child `SetSize()` calls use content width, not total width
- [ ] Run debug logging to confirm actual rendered widths match expected
- [ ] Check for any `Width(n - 2)` or `Width(n - borderWidth)` patterns -- these are v1 patterns that are BUGS in v2
- [ ] Check for `Width(content + padding)` patterns -- must also add border in v2
- [ ] Verify margin usage -- margins are still outside Width in v2 (unchanged)

---

## Common Mistakes Checklist (v2)

Before declaring "width calculations are done":

- [ ] **Did you add debug logging?** Measure calculated vs actual widths!
- [ ] **Did you verify with logs?** Run the app and check "WIDTH DEBUG" output
- [ ] **Are you using v2 Width() semantics?** Width = total rendered width (NOT content width)
- [ ] **Are you NOT subtracting border from Width()?** In v2, border is INSIDE Width
- [ ] **Are content widths computed as totalWidth - frame?** Use `GetHorizontalFrameSize()`
- [ ] **Are border widths correct?** 1 char per side (2 total), not 2 per side!
- [ ] **Do you have BOTH content and total width fields** for each pane?
- [ ] **Does `calculateLayout()` return updated model** (not multiple ints)?
- [ ] **Are all empirical adjustments zero?** If not, find the root cause!
- [ ] **Are you using dynamic properties** (like padding) to handle odd widths?
- [ ] **Are you lying to child components?** Don't tell them wrong sizes via `SetSize()`
- [ ] **Do you understand who applies borders?** (Parent model, not child)
- [ ] **Are auto-sized panes truly auto-sizing?** (No `Width()` set)
- [ ] **Are constrained panes properly constrained?** (`Width()` set correctly)

Before implementing overlay compositing:

- [ ] Are you using `ansi.StringWidth()` instead of `len()`?
- [ ] Are you using `ansi.Truncate()` / `ansi.TruncateLeft()` for string cutting?
- [ ] Are both background and foreground fully rendered (with all styling)?
- [ ] Did you handle the case where overlay extends beyond background?
- [ ] Did you test with styled/colored text to ensure ANSI codes don't break?

---

## Quick Reference

**v2 Key Facts:**
- `Width(n)` = total rendered width is `n` characters (border + padding + content all inside)
- `Height(n)` = total rendered height is `n` lines (border + padding + content all inside)
- Content area = `Width - GetHorizontalPadding() - GetHorizontalBorderSize()` (precise)
- Content area = `Width - GetHorizontalFrameSize()` (shorthand, only when margins are 0)
- Margins are OUTSIDE Width (unchanged from v1)
- Borders are 1 character per side (2 total) -- unchanged
- `GetHorizontalFrameSize()` returns total horizontal **margins + border + padding**
- `GetVerticalFrameSize()` returns total vertical **margins + border + padding**
- `GetHorizontalPadding()` + `GetHorizontalBorderSize()` = frame size WITHOUT margins

**Most Common v1-to-v2 Bugs:**
1. Subtracting border from `Width()` -- v2 already includes it
2. Setting `Width(content + padding)` without adding border -- causes premature word-wrapping
3. Passing total width instead of content width to child `SetSize()`
4. Using old import paths (`github.com/charmbracelet/` instead of `charm.land/`)
5. Using `tea.KeyMsg` instead of `tea.KeyPressMsg`
6. Returning `string` from `View()` instead of `tea.View`

---

**When in doubt, refer to this document. When you discover a new pattern, ADD IT HERE.**
