# Plan: Add Vertical Orientation to ChoiceModel

## Context

The gomion project needs a vertical companion selector modal. The current `ChoiceModel` only renders buttons horizontally (joined with `"  "`). This plan adds an `Orientation` field so callers can request vertical layout.

## Changes to `choice_model.go`

### 1. Add Orientation Type

```go
// Orientation controls how choice buttons are laid out
type Orientation int

const (
    Horizontal Orientation = iota // Default: buttons in a row
    Vertical                      // Buttons stacked vertically
)
```

### 2. Add to ChoiceModelArgs

```go
type ChoiceModelArgs struct {
    // ... existing fields ...
    Orientation Orientation // Horizontal (default) or Vertical
}
```

### 3. Store in ChoiceModel

```go
type ChoiceModel struct {
    // ... existing fields ...
    orientation Orientation
}
```

### 4. Propagate in NewChoiceModel

```go
m = ChoiceModel{
    // ... existing fields ...
    orientation: args.Orientation,
}
```

### 5. Modify renderButtons()

Current (line ~353):
```go
func (m ChoiceModel) renderButtons() (line string) {
    parts := make([]string, 0, len(m.options))
    for i, opt := range m.options {
        parts = append(parts, m.renderButton(opt, i == m.focusButton))
    }
    line = strings.Join(parts, "  ")
    return line
}
```

Change to:
```go
func (m ChoiceModel) renderButtons() (line string) {
    parts := make([]string, 0, len(m.options))
    for i, opt := range m.options {
        parts = append(parts, m.renderButton(opt, i == m.focusButton))
    }
    sep := "  "
    if m.orientation == Vertical {
        sep = "\n"
    }
    line = strings.Join(parts, sep)
    return line
}
```

### 6. Modify renderModal() Height Calculation

In `renderModal()` (line ~310), the height calculation currently assumes 1 line for all buttons:
```go
m.height = messageLines + 7
```

When vertical, buttons take `len(m.options)` lines instead of 1:
```go
buttonLines := 1
if m.orientation == Vertical {
    buttonLines = len(m.options)
}
m.height = messageLines + 6 + buttonLines
```

(The original `7` = message spacing (2) + buttons (1) + padding (2) + borders (2). Replace the `1` with `buttonLines`.)

### 7. Modify Update() Navigation Keys

When vertical, `NextButton`/`PrevButton` should use up/down semantics. The current Tab/Shift-Tab key bindings work for both orientations since they cycle through options. No changes needed — Tab/Shift-Tab already cycle regardless of layout direction.

## No Breaking Changes

- `Orientation` defaults to `Horizontal` (zero value of `iota`)
- All existing callers pass no `Orientation` field → get horizontal (existing behavior)
- No API removals or renames

## Verification

```bash
cd ~/Projects/go-pkgs/go-tealeaves && go build -o /dev/null ./teamodal/...
```

If teamodal has tests:
```bash
go test ./teamodal/...
```

---

# Plan: Add AllowCancel + [esc] Cancel hint to ChoiceModel

## Status: Code complete — needs manual visual verification

All code changes have been implemented and build/tests pass.

### What was done

- **CHG-INIT**: `NewChoiceModel` now initializes `cancelKeyStyle`/`cancelTextStyle` defaults, resolves `AllowCancel` (nil → true) and `ShowCancelHint` (nil → follows `AllowCancel`), and applies custom cancel style overrides
- **CHG-UPDATE**: `Update()` Cancel key handler guarded with `if !m.allowCancel`
- **CHG-RENDER**: `renderModal()` adds `[esc] Cancel` hint (centered, styled) below options when `showCancelHint` is true, with height adjusted by +2
- **CHG-EXAMPLE**: Removed manual "Cancel" option from `examples/teamodal/vertical/main.go`
- **Vet fix**: Changed `cliutil.Stderr` → `cliutil.Stderrf` in vertical example

### Remaining: Manual visual verification

Run the vertical example and confirm:

```bash
cd ~/Projects/go-pkgs/go-tealeaves/examples/teamodal/vertical && go run .
```

1. `[esc]` appears in green (color 46), `Cancel` in dim gray (color 244)
2. One blank line separates the last option from the hint
3. The hint is centered
4. Pressing Esc sends `ChoiceCancelledMsg`
5. Setting `AllowCancel: ptr(false)` hides the hint AND disables Esc
