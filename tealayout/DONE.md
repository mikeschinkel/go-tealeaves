# Add ResizeFocused to PaneLayout + MultiPaneLayout Starter

## Completed Steps

- **DIM** — Added `IsFlex()`, `Value()` to Dimension
- **PANE** — Added `minFlexWeight` field, `WithMinFlexWeight()`, `MinFlexWeight()`, `Dimension()` to Pane
- **RESIZE** — Added `ResizeFocused(delta)` to PaneLayout
- **MULTI** — New `MultiPaneLayout` with PaneDef, options, constructor, delegated methods
- **EXAMPLE** — Simplified threepane example to use MultiPaneLayout (removed 30-line resizeFocused, removed weight field from paneInfo)
- **TESTS** — 5 ResizeFocused tests + 8 MultiPaneLayout tests

## Review

- `go test ./...` — PASS
- `go vet ./...` — clean
- `go build ./...` — clean
- threepane example builds successfully
