# ADR-2025-01-24: String Compositing for Dropdown Overlay

## Status

**Accepted**

## Date

2025-01-24

## Context

The teadd dropdown component needs to overlay on top of a parent view. There are two primary approaches for rendering overlays in Bubble Tea applications:

1. **ANSI Escape Codes**: Use cursor positioning and save/restore escape sequences to draw the overlay
2. **String Compositing**: Render both views as strings, then merge them line-by-line

Key constraints:

* Dropdown must be a self-contained component (not wrap parent models)
* Parent retains control over final rendering
* Must work with lipgloss-styled content (ANSI codes present)
* Should be simple to use and test
* Must support accurate positioning with styled text

Existing patterns:

* **bubbletea-overlay** uses string compositing, based on Superfile's implementation
* **bubbleup** demonstrates ANSI-aware string operations for correct overlay positioning

## Decision

**Use string-compositing approach with ANSI-aware operations**, inspired by bubbletea-overlay and bubbleup.

### Implementation

The dropdown:
1. Renders its own box with borders and items via `View()` (returns styled string)
2. Parent is responsible for compositing the dropdown onto its base view
3. `OverlayDropdown(background, foreground, row, col)` helper function provided

Example usage:
```go
func (m model) View() string {
    baseView := renderParentContent()

    if m.dropdown.IsOpen {
        dropdownView := m.dropdown.View()
        return teadd.OverlayDropdown(baseView, dropdownView, m.dropdown.Row, m.dropdown.Col)
    }

    return baseView
}
```

### ANSI-Aware Operations

Use `github.com/charmbracelet/x/ansi` for string operations:
* `ansi.StringWidth()` - Get visual display width (ignores ANSI codes)
* `ansi.Truncate()` - Truncate to visual width, preserving ANSI codes
* `ansi.TruncateLeft()` - Skip first N visual columns

This ensures correct positioning when styled text contains ANSI escape sequences.

## Rationale

### Why String Compositing?

1. **Simpler than ANSI escape codes**
   - No cursor positioning state management
   - Easier to understand and debug
   - Testable with plain string comparison

2. **Parent control**
   - Parent decides when/where to composite
   - No hidden rendering side effects
   - Clean separation of concerns

3. **Self-contained component**
   - Dropdown doesn't wrap parent models
   - No two-model constraint
   - Better for standalone package

4. **Proven pattern**
   - bubbletea-overlay demonstrates viability
   - bubbleup shows ANSI-aware operations work
   - Used successfully in Superfile

### Why Not ANSI Escape Codes?

* More complex (cursor positioning, save/restore state)
* Harder to test (escape sequences in test strings)
* Less composable (parent loses rendering control)
* Requires careful escape sequence management

## Consequences

### Positive

* Simple API - parent just calls `OverlayDropdown()`
* Works correctly with lipgloss-styled content
* Easy to test (string comparison)
* Clean component boundaries
* Parent retains full control over rendering

### Trade-offs

* Parent must call `OverlayDropdown()` when dropdown is open
* String operations have O(n) overhead per overlay
* Requires ANSI-aware string operations for correct positioning

## Alternatives Considered

### ANSI Escape Codes

**Rejected** because:
* Complexity doesn't match use case
* Harder to test and debug
* Parent loses control over rendering order

### Wrapping Parent Model

**Rejected** because:
* Violates single responsibility (dropdown shouldn't manage parent)
* Two-model constraint limits composability
* Harder to integrate into existing apps

## Summary

String compositing with ANSI-aware operations provides the simplest, most composable approach for dropdown overlay rendering. The pattern is proven, testable, and gives parents full control while keeping the dropdown component self-contained.

---

*End of ADR*
