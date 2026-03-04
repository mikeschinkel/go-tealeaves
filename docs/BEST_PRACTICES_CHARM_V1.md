# Bubble Tea Best Practices

Hard-won lessons from building Bubble Tea UIs. Read this BEFORE implementing new models or components to avoid repeating mistakes.

## THE GOLDEN RULE

**TRUST LIPGLOSS. DON'T SECOND-GUESS IT.**

If you find yourself with "empirical adjustment" constants to fix widths, **you have a bug in your understanding**, not in lipgloss. Stop, step back, and find the root cause.

## Quick Reference

**Most Common Mistakes:**
1. Using fmt.Printf/fmt.Fprintf for debugging -- Use `*slog.Logger` instead!
2. Making fixes without measuring actual values -- Add debug logging first!
3. Thinking borders are 2 chars per side -- They're 1 char per side (2 total)
4. Not understanding Width() semantics -- Width = totalWidth - borderWidth
5. Manually subtracting padding/border from widths -- Trust lipgloss to handle it
6. Lying to child components about their size via `SetSize()` -- Use natural widths
7. Modifying content width to handle odd columns -- Use dynamic padding instead
8. Thinking child components handle borders -- Parent model applies borders
9. Using empirical adjustments as band-aids -- Find and fix the root cause
10. Checking modal.IsOpen() AFTER update to consume keys -- Check BEFORE update!
11. Using tea.Model type assertion on viewport/textinput/textarea -- They return concrete types, not tea.Model!


**Key Facts:**
- **DEBUG FIRST!** Measure calculated vs actual widths before making assumptions
- `Width()` in lipgloss includes padding but excludes border
- Borders are 1 character per side (2 total), not 2 per side (4 total)
- To get total width of 37: use `Width(37 - 2)` = `Width(35)`
- Parent models apply borders, child components don't know about them
- Auto-sized panes (no `Width()` set) determine their own size
- Constrained panes (with `Width()` set) use allocated space
- Dynamic properties (like padding) are better than modifying content width

---

## Width Calculations: Trust Lipgloss

### THE CRITICAL RULE

**TRUST LIPGLOSS - DON'T SECOND-GUESS IT**

You will be tempted to manually calculate widths by subtracting padding and border. **DON'T.** Lipgloss handles spacing internally. Your job is to:
1. Allocate space
2. Let lipgloss fit content into that space

### HOW LIPGLOSS WIDTH() WORKS

**Key fact:** `Width()` in lipgloss **includes padding but excludes border**.

```go
style := lipgloss.NewStyle().
    Width(100).              // Total pane width (includes padding, excludes border)
    PaddingLeft(1).          // 1 char inside the width
    BorderStyle(RoundedBorder())  // 4 chars OUTSIDE the width (2 per side)
```

Visual breakdown:
```
Border(1) | Pad(1) | Content | Pad(1) | Border(1)
|<------- Width(100) ------->|
|<------- Total Visual Width = 102 -------->|
```

**CRITICAL: Width() includes padding but excludes border!**

If you want total rendered width of 37:
```go
const borderWidth = 2 // 1 left + 1 right
totalWidth := 37
widthForLipgloss := totalWidth - borderWidth  // 35

style := lipgloss.NewStyle().
    Width(widthForLipgloss).  // 35 (includes padding)
    PaddingLeft(1).
    PaddingRight(2).
    BorderStyle(lipgloss.RoundedBorder())

// Result: 35 + 2 (border) = 37 total rendered width
```

### TWO APPROACHES TO SIZING

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
// You specify the width, lipgloss fits content inside
style := lipgloss.NewStyle().
    Width(totalWidth).  // You control this
    PaddingLeft(1).
    BorderStyle(lipgloss.RoundedBorder())

// Lipgloss ensures final visual width = totalWidth + border
```
Use for: Fixed-width panes like code viewers.

### THE TWO-WIDTH PROBLEM

**You still need TWO width values, but for a different reason than you think:**

1. **Total Width** - For lipgloss `Width()` (what the pane occupies)
2. **Content Width** - For child components `SetSize()` (what they can render)

**Track both in your model:**

```go
type MyModel struct {
    paneContentWidth int  // For viewport.SetSize()
    paneTotalWidth   int  // For lipgloss Width()
}
```

**But DON'T try to derive one from the other by subtracting constants!**

### THE CORRECT CALCULATION PATTERN

```go
func (m MyModel) calculateLayout() MyModel {
    // Calculate total width based on available space
    halfTotal := remainingWidth / 2
    m.paneTotalWidth = halfTotal  // This is what we allocate

    // For content width: TRUST LIPGLOSS
    // Just use the same base value - lipgloss handles the rest
    m.paneContentWidth = halfTotal  // Same value!

    return m
}
```

**Why are they the same?** Because lipgloss handles padding and border internally when you call `Width()`. You don't need to account for it manually.

### WHO HANDLES BORDERS?

**Critical understanding:** The **parent model** applies borders, not child components.

```go
// Parent model (e.g., CommitReviewModel)
func (m MyModel) View() string {
    childContent := m.childModel.View()  // Child returns raw content

    // Parent wraps with border
    style := lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
    return style.Render(childContent)
}
```

This means:
- Child models (teatree, viewport, etc.) have NO knowledge of borders
- Parent calculates total width including border
- Child receives content width via `SetSize()`
- Parent applies border when rendering

---

## Lipgloss Border Widths

**CRITICAL: Border width is 1 character per side, not 2!**

**Single borders (RoundedBorder, NormalBorder, etc.):**
- 1 character per side (the box-drawing character)
- 2 characters total width (left + right)
- Use `const borderWidth = 2` for single borders

**Double borders (DoubleBorder):**
- Still 1 character per side (just different Unicode characters)
- 2 characters total width
- Use `const borderWidth = 2`

**Common mistake:** Thinking borders are 2 chars per side because they "look thick". They're single Unicode characters!

---

## Padding Calculations

**Common mistake:** Forgetting to account for ALL padding.

```go
// Tree pane with padding on both sides
treeStyle := lipgloss.NewStyle().
    PaddingLeft(1).   // 1 char
    PaddingRight(1)   // 1 char
// Total padding width = 2

// Code pane with padding on one side only
codeStyle := lipgloss.NewStyle().
    PaddingLeft(1)    // 1 char
// Total padding width = 1
```

**Track padding separately per pane type:**
```go
const (
    treePaddingWidth = 2  // Left + Right
    codePaddingWidth = 1  // Left only
)
```

---

## Immutable Struct Pattern

### THE PATTERN

**Methods should update struct fields and return the model:**

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

### WHY

1. **Follows Bubble Tea conventions** - All model updates return updated model
2. **Caches values** - Calculated once, used multiple times
3. **Reduces errors** - No manual assignment of multiple return values
4. **Cleaner call sites:**

```go
// CORRECT
m = m.calculateLayout()
// Values available as m.contentWidth, m.totalWidth

// WRONG
contentWidth, totalWidth := m.calculateLayout()
// Now you have to track these variables separately
```

---

## Equal Width Panes: The Odd Width Problem

**Problem:** When splitting remaining width between two panes, odd numbers create unequal panes.

**WRONG SOLUTION:** Modify content width artificially

```go
// DON'T DO THIS!
if remainingWidth % 2 == 1 {
    treeContentWidth++   // Lying to the child model about its size
    treeTotalWidth++
    remainingWidth--
}
```

**Why this is wrong:**
- You're telling the child model it has more width than it naturally needs
- The child doesn't use the extra space
- Creates confusion between calculated and actual widths
- You're fighting lipgloss instead of working with it

**CORRECT SOLUTION:** Use dynamic padding to consume odd columns

```go
// In your model struct
type MyModel struct {
    treeRightPadding int  // 1 or 2, depending on odd/even
}

// In calculateLayout()
func (m MyModel) calculateLayout() MyModel {
    m.treeRightPadding = 1  // Default

    // Calculate what remaining width WOULD be
    normalTreeTotal := treeContentWidth + 1 + 1 + 4  // left pad + right pad + border
    potentialRemaining := terminalWidth - normalTreeTotal

    // If both panes visible and remaining would be odd, add extra padding
    if bothPanesVisible && potentialRemaining % 2 == 1 {
        m.treeRightPadding = 2  // Extra column as padding
    }

    // Now calculate actual tree total with dynamic padding
    m.treeTotalWidth = treeContentWidth + 1 + m.treeRightPadding + 4

    remainingWidth = terminalWidth - m.treeTotalWidth
    // remainingWidth is now guaranteed even!

    return m
}

// In View()
func (m MyModel) View() string {
    treeStyle := lipgloss.NewStyle().
        PaddingLeft(1).
        PaddingRight(m.treeRightPadding).  // Dynamic!
        BorderStyle(lipgloss.RoundedBorder())

    return treeStyle.Render(m.treePane.View())
}
```

**Why this is correct:**
- Tree content width stays natural (from `LayoutWidth()`)
- Extra column becomes whitespace inside the border
- Lipgloss handles it all - you're working WITH it
- No lies to child components about their size
- Code panes guaranteed equal

---

## Empirical Verification vs. Empirical Adjustments

### EMPIRICAL VERIFICATION (Good!)

**When debugging width issues, verify empirically:**

1. Do NOT use `fmt.Printf()` for debugging as it will throw off TUI display. You MUST use `*slog.Logger`.
2. **Use test pattern in file:**
   ```
   01234567890123456789012345678901234567890...
   ```
3. **Take screenshot showing exact cutoff point**
4. **Count characters to measure actual vs expected**
5. **Trust the compiler over gopls** - `go build` is source of truth

**This helps you find the ROOT CAUSE.**

### EMPIRICAL ADJUSTMENTS (Code Smell!)

**RED FLAG: Constants like this:**

```go
const (
    paneTotalAdjustment   = 3  // "Empirically determined"
    paneContentAdjustment = 1  // "Empirically determined"
)

// Used like:
m.paneTotalWidth = halfTotal + paneTotalAdjustment
m.paneContentWidth = halfTotal + paneContentAdjustment
```

**Why this is a problem:**
- You're compensating for a calculation error, not fixing it
- Adjustments are magic numbers - no one knows WHY they're needed
- Brittle - breaks when layout changes
- You're fighting lipgloss instead of working with it

**What to do instead:**
1. **Find the root cause** - Why do you need the adjustment?
2. **Common root causes:**
   - Not trusting lipgloss (manually subtracting padding/border)
   - Lying to child components about their size
   - Not accounting for dynamic properties (like padding)
   - Mixing up Width() vs content width
3. **Fix the calculation** - Eliminate the need for adjustments
4. **If you truly need adjustments** - Document WHY with code comments explaining the reason

**The goal:** All width adjustments should be 0. If they're not, you have a bug in your understanding or implementation.

---

## Calculation Order and Best Practices

**Correct order when mixing auto-sized and constrained panes:**

1. **Calculate auto-sized pane width first** (e.g., tree from `LayoutWidth()`)
2. **Determine dynamic properties** (e.g., padding based on odd/even)
3. **Calculate auto-sized pane total** (content + dynamic padding + border)
4. **Calculate remaining width** for constrained panes
5. **Distribute remaining width** (guaranteed even if step 2 handled it)

```go
// 1. Auto-sized content width
treeContentWidth = tree.LayoutWidth()

// 2. Dynamic padding to handle odd widths
treeRightPadding = 1  // Default
normalTreeTotal := treeContentWidth + 1 + 1 + 4
potentialRemaining := terminalWidth - normalTreeTotal
if potentialRemaining % 2 == 1 {
    treeRightPadding = 2  // Consume odd column
}

// 3. Auto-sized total with dynamic padding
treeTotalWidth = treeContentWidth + 1 + treeRightPadding + 4

// 4. Remaining for constrained panes
remainingWidth = terminalWidth - treeTotalWidth

// 5. Distribute (now guaranteed even)
halfWidth := remainingWidth / 2
pane2TotalWidth = halfWidth
pane3TotalWidth = halfWidth
pane2ContentWidth = halfWidth  // Trust lipgloss
pane3ContentWidth = halfWidth  // Trust lipgloss
```

**Key principles:**
- Auto-sized panes (no `Width()` set) determine their own size
- Constrained panes (with `Width()` set) use allocated space
- Handle odd columns with dynamic padding, NOT by lying about content width
- Trust lipgloss to handle padding and border - don't manually subtract

---

## When to Recalculate Layout

**Recalculate on:**
- `tea.WindowSizeMsg` - Terminal resized
- Toggling pane visibility
- Changing content that affects natural width (like tree)

**Store results in model fields:**
```go
case tea.WindowSizeMsg:
    m.terminalWidth = msg.Width
    m.terminalHeight = msg.Height
    m = m.recalculateLayout()  // Recalc and update viewports
```

---

## Overlays: Compositing Foreground on Background

### THE PROBLEM

You need to overlay one Bubble Tea component on top of another (modal dialogs, dropdowns, tooltips). Standard string concatenation doesn't work because:

1. **ANSI escape codes** - Syntax highlighting and lipgloss styling embed ANSI codes
2. **String width vs byte length** - `len(str)` counts ANSI codes, breaking positioning
3. **Positioning must be visual** - Overlays need to appear at specific screen columns

### THE SOLUTION

**Use ANSI-aware string operations from `github.com/charmbracelet/x/ansi`:**

```go
import "github.com/charmbracelet/x/ansi"

// CORRECT - Visual width
width := ansi.StringWidth(styledText)

// WRONG - Byte length (includes ANSI codes)
width := len(styledText)
```

### OVERLAY PATTERN

**Proven pattern from teadd and teamodal:**

```go
// OverlayDropdown overlays foreground view on background view at specified position.
// Uses ANSI-aware string operations to correctly handle styled text.
//
// Parameters:
//   - background: The base view (fully rendered string with ANSI codes)
//   - foreground: The overlay view (fully rendered string with ANSI codes)
//   - row: Line number in background where foreground row 0 should appear (0-indexed)
//   - col: Display column in background where foreground col 0 should appear (0-indexed)
//
// Returns:
//   - Composited view with foreground overlaid on background
func OverlayDropdown(background, foreground string, row, col int) string {
    var result strings.Builder

    bgLines := strings.Split(background, "\n")
    fgLines := strings.Split(foreground, "\n")

    for i, bgLine := range bgLines {
        fgRow := i - row

        // This line has no foreground overlay
        if fgRow < 0 || fgRow >= len(fgLines) {
            result.WriteString(bgLine)
            result.WriteString("\n")
            continue
        }

        // Overlay foreground line onto background line
        fgLine := fgLines[fgRow]
        composited := overlayLine(bgLine, fgLine, col)
        result.WriteString(composited)
        result.WriteString("\n")
    }

    // Remove trailing newline
    output := result.String()
    if len(output) > 0 && output[len(output)-1] == '\n' {
        output = output[:len(output)-1]
    }

    return output
}
```

**Per-line overlay (the critical piece):**

```go
// overlayLine overlays foreground onto background at column position (ANSI-aware).
// Pattern: left part of background + foreground + right part of background
//
// The key insight: Standard Go string operations (len, slicing) count ANSI escape
// codes as characters, which breaks positioning. We use ansi.StringWidth() for
// visual width and ansi.Truncate/TruncateLeft for ANSI-safe string cutting.
func overlayLine(background, foreground string, col int) string {
    if col < 0 {
        col = 0
    }

    bgWidth := ansi.StringWidth(background)
    fgWidth := ansi.StringWidth(foreground)

    var result strings.Builder

    // Left part: truncate background to col width
    if col > 0 {
        if col <= bgWidth {
            left := ansi.Truncate(background, col, "")
            result.WriteString(left)
        } else {
            // Need padding beyond background width
            result.WriteString(background)
            result.WriteString(strings.Repeat(" ", col-bgWidth))
        }
    }

    // Middle part: foreground content (the overlay)
    result.WriteString(foreground)

    // Right part: remainder of background after foreground
    endCol := col + fgWidth
    if endCol < bgWidth {
        // TruncateLeft(s, n) skips the first n display columns
        remaining := ansi.TruncateLeft(background, endCol, "")
        result.WriteString(remaining)
    }

    return result.String()
}
```

### ANSI-AWARE STRING OPERATIONS

**From `github.com/charmbracelet/x/ansi`:**

| Operation | Description |
|-----------|-------------|
| `ansi.StringWidth(s)` | Visual width (what you see on screen) |
| `ansi.Truncate(s, width, tail)` | Keep first N visual columns (cut from right) |
| `ansi.TruncateLeft(s, width, tail)` | Skip first N visual columns (cut from left) |

**Example:**

```go
styled := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Hello")
// styled contains ANSI escape codes

len(styled)                    // WRONG: 18 (includes escape codes)
ansi.StringWidth(styled)       // CORRECT: 5  (visual width)

// Keep first 3 visual characters
truncated := ansi.Truncate(styled, 3, "")  // "Hel" (still styled!)

// Skip first 2 visual characters
remaining := ansi.TruncateLeft(styled, 2, "")  // "llo" (still styled!)
```

### USAGE IN YOUR MODEL

**In your View() method:**

```go
func (m MyModel) View() string {
    // Render main content
    mainView := m.mainContent.View()

    // If modal/dropdown is active, overlay it
    if m.showModal {
        modalView := m.modal.View()

        // Calculate position (typically centered)
        row := (m.height - m.modal.Height()) / 2
        col := (m.width - m.modal.Width()) / 2

        // Overlay modal on main view
        return OverlayModal(mainView, modalView, row, col)
    }

    return mainView
}
```

### POSITIONING HELPERS

**Center overlay on screen:**

```go
row := (screenHeight - overlayHeight) / 2
col := (screenWidth - overlayWidth) / 2
```

**Position below trigger element:**

```go
row := triggerRow + 1
col := triggerCol
```

**Right-align overlay:**

```go
row := /* desired row */
col := screenWidth - overlayWidth
```

### REFERENCE IMPLEMENTATIONS

- **teadd/overlay_dropdown.go** - Dropdown pattern
- **teamodal/overlay_modal.go** - Modal pattern

Both use identical overlay logic. The difference is in usage context, not implementation.

---

## Modal Key Consumption: Preventing ESC from Propagating

### THE PROBLEM

You have a modal dialog that closes when the user presses `esc`. The modal handles the key, closes itself, but the `esc` key **continues to propagate** to the parent model, causing unintended side effects like:

- Closing the current view (drilling up)
- Quitting the application
- Closing another component that also listens for `esc`

**Why this happens:** In Bubble Tea, messages flow through the entire model tree. When a modal handles a key in its `Update()`, it returns a command but the **same key message** is also passed to other components unless explicitly consumed.

### THE PATTERN: Check Modal State Before Propagating

**Key insight:** Check if the modal was open *before* updating it, then block key propagation if it was.

```go
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var modalCmd tea.Cmd
    var viewStackCmd tea.Cmd

    // Track if modal was open BEFORE update (to consume key events)
    modalWasOpen := m.modal.IsOpen()

    // Update modal (it may handle the key and close itself)
    m.modal, modalCmd = m.modal.Update(msg)

    // Only pass messages to view stack if no modal was open
    // This prevents key events (like ESC) from propagating through
    if _, isKey := msg.(tea.KeyMsg); isKey && modalWasOpen {
        // Modal consumed the key - don't pass to view stack
    } else {
        viewStackCmd = m.updateViewStack(msg)
    }

    return m, tea.Batch(modalCmd, viewStackCmd)
}
```

### WHY "WAS OPEN" NOT "IS OPEN"

**WRONG:** Checking `m.modal.IsOpen()` *after* update

```go
// DON'T DO THIS
m.modal, modalCmd = m.modal.Update(msg)
if m.modal.IsOpen() {  // Modal already closed! Returns false!
    // This block never executes
}
```

**Why it fails:** If user presses `esc` to close the modal:
1. `Update()` handles `esc` and sets `isOpen = false`
2. After update, `IsOpen()` returns `false`
3. Your guard condition doesn't trigger
4. Key propagates to view stack

**CORRECT:** Check state *before* update

```go
// Check BEFORE update
modalWasOpen := m.modal.IsOpen()  // True! Modal is still open

m.modal, modalCmd = m.modal.Update(msg)
// Now modal is closed, but modalWasOpen is still true

if _, isKey := msg.(tea.KeyMsg); isKey && modalWasOpen {
    // This block executes - key is consumed!
}
```

### MULTIPLE MODALS

If you have multiple modals, check all of them:

```go
// Any modal open = block key propagation
anyModalOpen := m.verifyModal.IsOpen() ||
                m.warningModal.IsOpen() ||
                m.confirmModal.IsOpen()

m.verifyModal, verifyCmd = m.verifyModal.Update(msg)
m.warningModal, warningCmd = m.warningModal.Update(msg)
m.confirmModal, confirmCmd = m.confirmModal.Update(msg)

if _, isKey := msg.(tea.KeyMsg); isKey && anyModalOpen {
    // Don't pass to view stack
} else {
    viewStackCmd = m.updateViewStack(msg)
}
```

### EXISTING ALERT PATTERN

This pattern already exists for alerts in the codebase:

```go
// Check if alert is active before updating (ESC will close it)
hadActiveAlert := m.Alert.HasActiveAlert()

// Delegate all messages to Alert first (for ticks and ESC)
m, alertCmd = m.alertCmd(msg)

// If alert was active and this is ESC, consume it (don't pass to views)
if hadActiveAlert {
    if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
        return m, alertCmd
    }
}
```

The modal pattern is the same concept, applied at the end of `Update()` rather than the beginning.

---

## Emojis and Wide Characters: The Unsolvable Alignment Problem

**Emojis and many UTF-8 characters CANNOT be reliably aligned in terminal UIs.**

This is not a bug in your code. This is not a bug in lipgloss. This is not a bug in `ansi.StringWidth()`. This is a **fundamental mismatch** between Unicode width tables and terminal rendering that cannot be fixed programmatically.

### THE PROBLEM

When you use emojis or wide Unicode characters in width-sensitive layouts (like list items that need to align), you will encounter off-by-one (or more) alignment errors:

```
👤 Commit Batch #1
✨ Recommendation: 6 commits [ACTIVE]
```

Even though your code calculates both lines as the same width, they will render at different visual widths. The reverse-video highlight on the selected item won't align with other items.

### WHY THIS HAPPENS

**`ansi.StringWidth()` and the terminal disagree on character widths.**

Unicode Standard Annex #11 defines "East Asian Width" properties for characters:

| Property | Meaning | Cell Width |
|----------|---------|------------|
| N (Neutral) | Not East Asian | 1 cell |
| Na (Narrow) | Narrow | 1 cell |
| W (Wide) | Wide | 2 cells |
| F (Fullwidth) | Fullwidth | 2 cells |
| H (Halfwidth) | Halfwidth | 1 cell |
| A (Ambiguous) | Context-dependent | 1 or 2 cells |

But:

1. **Different libraries use different tables** - Go's unicode package, charmbracelet/x/ansi, glibc wcwidth, musl wcwidth all have different implementations
2. **Terminals make their own decisions** - Based on font, Unicode version, OS rendering engine
3. **Emojis are especially problematic** - Many are classified as "Neutral" or "Ambiguous" width
4. **There's no way to query the terminal** - You cannot ask "how wide will this character render?"

### THE MISMATCH IN DETAIL

**Example emojis:**

| Emoji | Codepoint | Unicode Name | East Asian Width |
|-------|-----------|--------------|------------------|
| 👤 | U+1F464 | Bust in Silhouette | Neutral (N) |
| ✨ | U+2728 | Sparkles | Neutral (N) |

Both are "Neutral" in the Unicode standard. However `ansi.StringWidth()` uses wcwidth-style lookup tables that may classify these as 2-cell wide based on emoji presentation rules, while the terminal makes its OWN decision based on font, Unicode version, and OS rendering engine.

```
What ansi.StringWidth() reports:
  "👤 Commit Batch #1"           -> 18 columns (assuming 👤 = 2 cols)
  "✨ Recommendation: 6 commits" -> 28 columns (assuming ✨ = 2 cols)

What the terminal MIGHT render:
  "👤 Commit Batch #1"           -> 17 columns (if 👤 = 1 col)
  "✨ Recommendation: 6 commits" -> 28 columns (if ✨ = 2 cols)
```

If we pad item 1 with 10 spaces (28 - 18 = 10) based on `ansi.StringWidth()`:
- Our calculation: 18 + 10 = 28 columns (equal to item 2)
- Terminal reality: 17 + 10 = 27 columns (1 less than item 2)

**Result: Visual misalignment of 1 character.**

### WHY EACH FIX ATTEMPT FAILS

**Attempt 1: Use `ansi.StringWidth()` instead of rune count**

`ansi.StringWidth()` returns a value based on Unicode tables, but those tables don't match what the specific terminal actually renders.

**Attempt 2: Manual padding with `padLabel()`**

Same underlying issue. We calculate padding based on `ansi.StringWidth()`, but if that doesn't match terminal rendering, the padding amount is wrong.

**Attempt 3: Two-phase rendering (measure then pad)**

We're still using `ansi.StringWidth()` to measure. If it reports both lines as equal width (e.g., both 46 columns), we add zero padding. But if the terminal renders them at different widths, they still misalign.

Debug output confirmed both lines measure as 46 columns, yet visually they're different widths.

**Attempt 4: Avoid lipgloss Width() since it uses rune count**

The replacement (`ansi.StringWidth()`) has its own terminal mismatch problem.

### WHY THIS IS FUNDAMENTALLY UNSOLVABLE

**There is no universal width oracle.** There is no programmatic way to ask a terminal "how many columns wide will this character render?"

Terminals vary wildly. The width depends on terminal version, font, Unicode version support, emoji presentation selectors, and OS rendering engine.

Different systems have different wcwidth implementations (glibc, musl, macOS, Go's unicode tables, charmbracelet/x/ansi). None are guaranteed to match what your specific terminal renders.

### WHAT DOESN'T WORK

| Approach | Why It Fails |
|----------|--------------|
| Use `ansi.StringWidth()` | Returns library's table value, not terminal's rendering |
| Use `lipgloss.Width()` | Uses rune count, not display width |
| Measure after styling | ANSI codes are ignored, but emoji width mismatch remains |
| Two-phase rendering | Still uses `ansi.StringWidth()` to measure |
| Empirical adjustments | Different terminals need different adjustments |
| Terminal detection | Maintenance nightmare, incomplete coverage |

### THE ONLY RELIABLE SOLUTIONS

**Option 1: Don't use emojis (Recommended)**

```go
// Unreliable
"👤 Commit Batch #1"
"✨ Recommendation: 6 commits"

// Reliable
"[U] Commit Batch #1"
"[R] Recommendation: 6 commits"
```

ASCII characters have consistent 1-column width across all terminals.

**Option 2: Use symbols that are reliably single-cell**

```go
// More likely to be consistent (but still not guaranteed)
"● Commit Batch #1"     // U+25CF Black Circle
"★ Recommendation: ..." // U+2605 Black Star
```

**Option 3: Accept imperfect alignment**

If emojis are essential to your design, accept that alignment will be imperfect on some terminals.

**Option 4: Right-align variable content**

Use a layout that doesn't depend on character-level alignment. Avoids the padding problem entirely but may not fit the desired UI design.

### SYMPTOMS OF THIS PROBLEM

You know you've hit this issue when:

1. Debug logging shows both lines have identical `ansi.StringWidth()` values
2. Your padding math is correct
3. You're using `ansi.StringWidth()` not `len()` or rune count
4. But the visual alignment is still off by 1+ characters

### THE BOTTOM LINE

**For width-sensitive layouts, use ASCII characters only.**

Emojis are fine for:
- Status messages
- Notifications
- Headers and titles
- Anywhere alignment doesn't matter

Emojis are NOT fine for:
- List items that need to align
- Table columns
- Fixed-width fields
- Anywhere padding must be precise

### References

- [Unicode Standard Annex #11: East Asian Width](https://www.unicode.org/reports/tr11/)
- [wcwidth and its problems](https://github.com/jquast/wcwidth)
- [Terminal emoji rendering issues](https://github.com/charmbracelet/lipgloss/issues)
- [The state of emoji on the command line](https://tonsky.me/blog/unicode/)

---

## Debugging Width Issues: The Systematic Approach

**CRITICAL LESSON: Debug first, fix second. Never make "confident" claims without measuring actual values!**

### THE PROBLEM PATTERN

You calculate widths mathematically, it looks perfect on paper, but the UI has a gap. You try fixes based on assumptions. Nothing works. User gets frustrated.

**Why:** You're debugging your mental model, not the actual code.

### THE SOLUTION: MEASURE EVERYTHING

**Step 1: Add debug logging IMMEDIATELY**

Don't guess. Don't assume. Measure actual runtime values:

```go
// In calculateLayout()
m.Logger.Info("WIDTH DEBUG calculateLayout",
    "terminalWidth", m.terminalWidth,
    "treeContentWidth", m.treeContentWidth,
    "treeTotalWidth", m.treeTotalWidth,
    "splitTotalWidth", m.splitTotalWidth,
    "splitContentWidth", m.splitContentWidth,
)

// In View() - measure what actually renders
treeRendered := treeStyle.Render(m.treePane.View())
treeActualWidth := 0
if lines := strings.Split(treeRendered, "\n"); len(lines) > 0 {
    treeActualWidth = ansi.StringWidth(lines[0])
}

m.Logger.Info("WIDTH DEBUG View",
    "treeCalculated", m.treeTotalWidth,
    "treeActual", treeActualWidth,
    "gap", m.treeTotalWidth - treeActualWidth,
)
```

**Step 2: Run the app, capture logs**

```bash
cd gommod && ../cmd/gomion/gomion commit 2>&1 | tee /tmp/debug.log
# Then: grep "WIDTH DEBUG" /tmp/debug.log
```

**Step 3: Compare calculated vs actual**

Look for discrepancies:
- `treeCalculated: 37` but `treeActual: 34` -- Tree is 3 chars too narrow!
- `splitCalculated: 283` and `splitActual: 283` -- Split pane is correct!

**Step 4: Find the root cause**

In the example above:
- Tree is 3 chars short
- 3 = border width (2) + 1
- Hypothesis: Width() might include padding but exclude border
- Test: Change `Width(contentWidth)` to `Width(totalWidth - borderWidth)`
- Verify: Check logs again

## LOGGING: NEVER USE fmt.Printf/fmt.Fprintf IN TUI CODE

**This is CRITICAL. fmt.Printf() and fmt.Fprintf(os.Stderr, ...) will COMPLETELY DESTROY the TUI display.**

### THE RULE

```go
// NEVER DO THIS - corrupts TUI display
fmt.Printf("DEBUG: value=%d\n", value)
fmt.Fprintf(os.Stderr, "DEBUG: %s\n", msg)
log.Printf("DEBUG: %v", thing)

// ALWAYS DO THIS - goes to log file, not TUI
m.Logger.Debug("debug message", "value", value)
m.Logger.Info("info message", "key", value)
```

### WHY

`fmt.Printf` and `fmt.Fprintf(os.Stderr, ...)` write directly to stdout/stderr. In a TUI application:
- stdout IS the terminal display
- stderr also goes to the terminal
- Your "debug" output gets interleaved with the rendered TUI
- The display becomes corrupted, unreadable garbage
- You can't even see if your fix worked because the logging broke the display!

### THE PATTERN

**Every model that needs logging should have a Logger field:**

```go
type MyModel struct {
    Logger *slog.Logger
    // ... other fields
}

type MyModelArgs struct {
    Logger *slog.Logger
    // ... other args
}

func NewMyModel(args *MyModelArgs) MyModel {
    return MyModel{
        Logger: args.Logger,
        // ...
    }
}
```

**The logger is created ONCE at app startup and injected into all models:**

```go
// In main.go or run_cli.go - logger is created EARLY
logger := slog.New(slog.NewTextHandler(logFile, nil))

// Pass to all models
appModel := NewAppModel(args.Logger)
// AppModel passes to child models...
```

**DO NOT create your own logger. Use the injected one.**

**If a model doesn't have a Logger yet, ADD IT:**

```go
// 1. Add to model struct
type MyModel struct {
    Logger *slog.Logger  // ADD THIS
    // ... existing fields
}

// 2. Add to args struct
type MyModelArgs struct {
    Logger *slog.Logger  // ADD THIS
    // ... existing args
}

// 3. Set in constructor
func NewMyModel(args *MyModelArgs) MyModel {
    return MyModel{
        Logger: args.Logger,  // ADD THIS
        // ...
    }
}

// 4. Pass from parent when creating
childModel := NewMyModel(&MyModelArgs{
    Logger: m.Logger,  // Pass parent's logger
    // ...
})
```

This ensures the logger flows down from app startup through all models.

### DEBUG LOGGING EXAMPLE

```go
func (m MyModel) calculateLayout() MyModel {
    // ... calculations ...

    if m.Logger != nil {
        m.Logger.Debug("layout calculated",
            "screenW", m.screenWidth,
            "screenH", m.screenHeight,
            "vpWidth", vpWidth,
            "vpHeight", vpHeight)
    }

    return m
}
```

### VIEWING LOGS

Logs go to the configured log file (typically in ~/.config/gomion/ or similar).
Use `tail -f /path/to/log` in another terminal to watch in real-time.

---

### COMMON DISCREPANCIES AND CAUSES

| Symptom | Likely Cause |
|---------|--------------|
| Actual = Calculated - 2 | Forgot to account for borders (1 per side) |
| Actual = Calculated - 3 | Width() excludes border, you didn't subtract it |
| Actual < Calculated | Content is narrower than you think, not being padded |
| Actual > Calculated | You're adding padding/border twice |

### MEASURE CHILD CONTENT TOO

Don't just measure the final rendered output. Measure intermediate values:

```go
// Measure tree content BEFORE styling
treeContent := m.treePane.View()
treeContentActualWidth := 0
if lines := strings.Split(treeContent, "\n"); len(lines) > 0 {
    treeContentActualWidth = ansi.StringWidth(lines[0])
}

m.Logger.Info("WIDTH DEBUG tree details",
    "treeContentWidth", m.treeContentWidth,      // What you calculated
    "treeContentActualWidth", treeContentActualWidth,  // What it really is
    "treeTotalWidth", m.treeTotalWidth,          // What you expect total
    "treeRenderedWidth", treeActualWidth,        // What actually rendered
)
```

This reveals:
- Is the child producing the width you expect?
- Is lipgloss styling it correctly?
- Where exactly is the discrepancy?

### THE GOLDEN RULE

**Never make a "fix" without first understanding the discrepancy.**

Bad: "I'll add 3 to the width, that should fix it"
Bad: "Let me try subtracting the border here"
Bad: "Maybe if I adjust the padding..."

Good: "Logs show actual is 34 but calculated is 37. That's a 3-char gap. Let me investigate why."
Good: "Tree content is 17 chars but I'm setting Width(32). Is lipgloss padding it correctly?"
Good: "Width() might include padding. Let me verify with documentation."

### AFTER THE FIX

Leave the debug logging in place! Comment it out if needed, but don't delete it. Next time layout breaks, you'll be glad you can uncomment and immediately see the values.

```go
// DEBUG: Uncomment to debug width calculations
// m.Logger.Info("WIDTH DEBUG calculateLayout", ...)
```

---

## Component Update() Return Types: Not All Return tea.Model

### THE GOTCHA

Some Bubble Tea components return their **concrete type** from `Update()`, not `tea.Model`. This means you cannot use `tea.Model` type assertion with them.

**`viewport.Model`** is the most common example:

```go
// WRONG - viewport.Model.Update() does NOT return tea.Model
var updated tea.Model
updated, cmd = m.contextViewport.Update(msg)          // Compile error!
m.contextViewport = updated.(viewport.Model)           // Compile error!

// CORRECT - viewport.Model.Update() returns viewport.Model directly
m.contextViewport, cmd = m.contextViewport.Update(msg) // Works!
```

**The error message is clear but surprising:**

```
cannot use m.contextViewport.Update(msg) as tea.Model value:
  viewport.Model does not implement tea.Model
  (wrong type for method Update)
    have Update(tea.Msg) (viewport.Model, tea.Cmd)
    want Update(tea.Msg) (tea.Model, tea.Cmd)
```

### WHY THIS HAPPENS

Bubble Tea's `tea.Model` interface requires:

```go
type Model interface {
    Init() tea.Cmd
    Update(tea.Msg) (Model, tea.Cmd)  // Returns tea.Model
    View() string
}
```

But `viewport.Model.Update()` returns `viewport.Model`, not `tea.Model`. This is intentional -- viewport is designed as an **embedded component**, not a standalone model. The same applies to other charmbracelet components like `textinput.Model` and `textarea.Model`.

### WHICH COMPONENTS ARE AFFECTED

| Component | Update() Returns | Use tea.Model assertion? |
|-----------|-----------------|--------------------------|
| `viewport.Model` | `viewport.Model` | No -- assign directly |
| `textinput.Model` | `textinput.Model` | No -- assign directly |
| `textarea.Model` | `textarea.Model` | No -- assign directly |
| `teamodal.ChoiceModel` | `tea.Model` | Yes -- type assert |
| `teamodal.ModalModel` | `tea.Model` | Yes -- type assert |
| `teadep.PathViewerModel` | `tea.Model` | Yes -- type assert |

### THE RULE

**If a component's `Update()` returns `tea.Model`, use type assertion. If it returns its concrete type, assign directly.** When unsure, check the component's method signature -- the compiler will tell you.

## Viewport Horizontal Scrolling: Use Native Support

### THE WRONG WAY

**Do NOT manually track horizontal offset and apply it post-render:**

```go
// WRONG — manual horizontal scrolling
type MyModel struct {
    diffViewport viewport.Model
    hOffset      int            // Manual offset tracking
}

func (m MyModel) handleKey(msg tea.KeyMsg) MyModel {
    switch {
    case key.Matches(msg, m.Keys.ScrollLeft):
        if m.hOffset > 0 {
            m.hOffset--
        }
    case key.Matches(msg, m.Keys.ScrollRight):
        m.hOffset++
    }
    return m
}

func (m MyModel) View() string {
    content := m.diffViewport.View()
    // Apply offset after viewport already truncated lines to its width
    if m.hOffset > 0 {
        lines := strings.Split(content, "\n")
        for i, line := range lines {
            lines[i] = ansi.TruncateLeft(line, m.hOffset, "")
        }
        content = strings.Join(lines, "\n")
    }
    return content
}
```

**Why this fails:** `viewport.View()` internally truncates lines to the viewport width. Then `TruncateLeft` removes more characters from the left. Content is lost on BOTH sides — left (by TruncateLeft) and right (by viewport truncation). The visible window shrinks as you scroll right.

### THE RIGHT WAY

**Use the viewport's native horizontal scrolling — it's built in, just disabled by default:**

```go
// CORRECT — native horizontal scrolling
type MyModel struct {
    diffViewport viewport.Model
    // No manual hOffset field needed!
}

func NewMyModel() MyModel {
    vp := viewport.New(0, 0)
    vp.SetHorizontalStep(4) // Enable with 4-char scroll steps
    return MyModel{diffViewport: vp}
}

func (m MyModel) handleKey(msg tea.KeyMsg) (MyModel, tea.Cmd) {
    switch {
    case key.Matches(msg, m.Keys.Home):
        m.diffViewport.GotoTop()
        m.diffViewport.SetXOffset(0) // Reset horizontal position
    default:
        // Viewport handles left/right/up/down/pgup/pgdn natively
        var cmd tea.Cmd
        m.diffViewport, cmd = m.diffViewport.Update(msg)
        return m, cmd
    }
    return m, nil
}

func (m MyModel) View() string {
    // Just use View() directly — viewport handles horizontal windowing
    return m.diffViewport.View()
}
```

**Why this works:** The viewport's internal `visibleLines()` uses `ansi.Cut(line, xOffset, xOffset+width)` to extract the correct visible window from each line. It also:
- Automatically tracks the longest line width on `SetContent()`
- Clamps scroll offset so you can't scroll past content
- Supports mouse wheel (Shift+wheel for horizontal)
- Reports `HorizontalScrollPercent()` for scroll indicators

### KEY FACTS

| Feature | Manual Approach | Native viewport |
|---------|----------------|-----------------|
| Offset tracking | Manual field | Internal `xOffset` |
| Line windowing | `TruncateLeft` (breaks) | `ansi.Cut` (correct) |
| Max offset clamping | None | Automatic |
| Longest line tracking | None | Automatic on `SetContent()` |
| Mouse support | None | Shift+wheel |
| Scroll step | 1 char (slow) | Configurable via `SetHorizontalStep()` |

### ENABLING NATIVE HORIZONTAL SCROLLING

```go
vp := viewport.New(width, height)
vp.SetHorizontalStep(4)  // Must be > 0 to enable; 0 = disabled (default)
```

The viewport's default keymap already binds `left`/`h` and `right`/`l` for horizontal scrolling — they just do nothing until `SetHorizontalStep()` is called with a positive value.

### RESETTING ON CONTENT CHANGE

```go
m.diffViewport.SetContent(newDiff)
m.diffViewport.GotoTop()
m.diffViewport.SetXOffset(0)  // Reset horizontal scroll too
```

---

## Tab Characters in Viewport Content

### THE PROBLEM

Tab characters (`\t`) have variable visual width depending on column position. When the viewport applies horizontal scrolling via `ansi.Cut()`, it counts characters, not visual columns. This causes content to jump and misalign as you scroll horizontally — tabs that appeared aligned at offset 0 become ragged at other offsets.

### THE FIX

Replace tabs with spaces before setting viewport content:

```go
diff = strings.ReplaceAll(diff, "\t", "    ")
m.diffViewport.SetContent(diff)
```

This is safe for Go source code (which uses tabs for indentation) because each tab is consistently one indentation level. Fixed-width spaces preserve alignment at every scroll position.

**Do this at the content-setting boundary**, not in rendering. The viewport needs consistent character widths to scroll correctly.

---


---

## Common Mistakes Checklist

Before declaring "width calculations are done":

- [ ] **Did you add debug logging?** Measure calculated vs actual widths!
- [ ] **Did you verify with logs?** Run the app and check "WIDTH DEBUG" output
- [ ] **Are calculated and actual widths equal?** If not, find root cause first!
- [ ] **Are you trusting lipgloss?** No manual subtraction of padding/border?
- [ ] **Do you understand Width() semantics?** Width = totalWidth - borderWidth
- [ ] **Are border widths correct?** 1 char per side (2 total), not 2 per side!
- [ ] **Do you have BOTH content and total width fields** for each pane?
- [ ] **Does `calculateLayout()` return updated model** (not multiple ints)?
- [ ] **Are all empirical adjustments zero?** If not, find the root cause!
- [ ] **Are you using dynamic properties** (like padding) to handle odd widths?
- [ ] **Are you lying to child components?** Don't tell them wrong sizes via `SetSize()`
- [ ] **Do you understand who applies borders?** (Parent model, not child)
- [ ] **Did you verify empirically** with test pattern?
- [ ] **Are auto-sized panes truly auto-sizing?** (No `Width()` set)
- [ ] **Are constrained panes properly constrained?** (`Width()` set correctly)

Before implementing overlay compositing:

- [ ] Are you using `ansi.StringWidth()` instead of `len()`?
- [ ] Are you using `ansi.Truncate()` / `ansi.TruncateLeft()` for string cutting?
- [ ] Are both background and foreground fully rendered (with all styling)?
- [ ] Did you handle the case where overlay extends beyond background?
- [ ] Did you test with styled/colored text to ensure ANSI codes don't break?

---

## Reference Implementations

### Width Calculations
See `gommod/gomtui/commit_review_model.go` for complete reference:
- Width cache fields in struct (lines ~47-52)
- `calculateLayout()` populates cache and returns model (lines ~480-530)
- `View()` uses total widths directly (lines ~221-270)
- `recalculateLayout()` uses content widths for viewports (lines ~558-573)

### Overlay Compositing
See:
- `teadd/overlay_dropdown.go` - Complete overlay implementation
- `teamodal/overlay_modal.go` - Same pattern, different context

---

**When in doubt, refer to this document. When you discover a new pattern, ADD IT HERE.**
