# teamodal Design Fix Instructions

## Problem

The current design requires consumers to manually calculate center position even though the modal has all the data needed to center itself.

**Current issues:**
- `OverlayModal(background, foreground, row, col)` is a package function that requires explicit positioning
- Consumers must manually calculate center position even though the modal has all the data needed (`ScreenWidth`, `ScreenHeight`, `width`, `height`)
- `calculateCenter()` is private, forcing consumers to duplicate centering logic

**Example of current awkward usage:**
```go
// Consumer must manually calculate center
func (m FileIntentModel) overlayModal(baseView, modalView string) string {
    modalHeight := strings.Count(modalView, "\n") + 1
    row := (m.terminalHeight - modalHeight) / 2
    col := (m.terminalWidth - 60) / 2 // Assume ~60 width modal
    if row < 0 { row = 0 }
    if col < 0 { col = 0 }
    return teamodal.OverlayModal(baseView, modalView, row, col)
}
```

## Solution

Change `OverlayModal()` to overlays itself centered on the background.

**Add this method to `model.go`:**

```go
// OverlayModal renders the modal centered over the background view.
// This is a convenience method that handles positioning automatically.
// For custom positioning, use the package-level OverlayModal() function.
func (m ModalModel) OverlayModal(background string) (view string) {
	var row, col int
	var err error

	if !m.IsOpen {
		view = background
		goto end
	}

	// Calculate center position
	row, col, err = m.calculateCenter()
	if err != nil {
		// Fallback to top-left if calculation fails
		row, col = 0, 0
	}

	// Render modal view and overlay it
	modalView := m.View()
	// TODO: Put the former logic for OverlayModal() here 
	view = ...

end:
	return view
}
```

## Benefits

- **Self-contained:** Modal knows how to position itself
- **Simpler API:** Consumers just call `m.OverlayModal(baseView)`
- **Single responsibility:** Modal handles its own rendering and positioning
- **No code duplication:** Consumers don't duplicate centering logic

## Usage Example

**Before (manual centering):**
```go
func (m model) View() string {
    view := baseView

    if m.confirmDialog.IsOpen {
        modalView := m.confirmDialog.View()
        modalHeight := strings.Count(modalView, "\n") + 1
        row := (m.height - modalHeight) / 2
        col := (m.width - 60) / 2
        if row < 0 { row = 0 }
        if col < 0 { col = 0 }
        view = teamodal.OverlayModal(view, modalView, row, col)
    }

    return view
}
```

**After (automatic centering):**
```go
func (m model) View() string {
    view := baseView

    if m.confirmDialog.IsOpen {
        view = m.confirmDialog.OverlayModal(view)
    }

    return view
}
```

## Update Example

Update `example/main.go` to demonstrate the simpler API:

```go
func (m model) View() (view string) {
    // ... build baseView ...

    // Old way (still works for custom positioning):
    // view = teamodal.OverlayModal(baseView, modalView, row, col)

    // New way (recommended for centered modals):
    if m.confirmDialog.IsOpen {
        view = m.confirmDialog.OverlayModal(view)
    }

    if m.alertDialog.IsOpen {
        view = m.alertDialog.OverlayModal(view)
    }

    return view
}
```

## Notes

- `calculateCenter()` stays private (no need to export it)
