# Handoff: Implement v1 Test Suite for go-tealeaves

## Your Mission

Write comprehensive v1 tests for all go-tealeaves modules per `test/PLAN.md`. This is the TESTV1 phase of the Charm v2 upgrade — these tests become the regression suite that validates the v2 migration doesn't break anything.

## Project Location

`/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/`

## Current State

- **Baseline verified:** `make test` and `make vet` pass (green)
- **Only existing tests:** `teamodal/choice_model_test.go` (276 lines, Layer 1 only — covers ChoiceModel navigation, selection, hotkeys, cancel, default index, closed modal, overlay)
- **All other modules have zero test files**
- **Test plan written:** `test/PLAN.md` — ~170 test specifications across 8 modules, 3 testing layers

## What You Need to Know

### Repository Structure

7 core modules (each has its own `go.mod`):
- `teautils/` — Pure utility functions (positioning, rendering, key registry, help visor)
- `teadd/` — Popup dropdown selection component
- `teastatus/` — Two-zone status bar (left menu items, right indicators)
- `teamodal/` — Modal dialogs (OK, YesNo, Progress, Choice, List)
- `teadep/` — Dependency path viewer with dropdown (depends on teadd)
- `teatree/` — Generic tree navigation with viewport scrolling
- `teatextsel/` — Text selection wrapping textarea (shift-select, clipboard)

Plus `teanotify/` and `teagrid/` which don't exist yet — skip those. They'll be tested after creation.

### Three Testing Layers

1. **Layer 1 — Direct model tests:** Construct model, call `Update()` with `tea.KeyMsg`/`tea.MouseMsg`, assert on returned model state and commands. Standard `testing` package only.
2. **Layer 2 — View() output assertions:** Call `View()` on models in known states, assert on rendered strings. Check border geometry, alignment, content, ANSI-aware widths.
3. **Layer 3 — teatest program tests:** Use `github.com/charmbracelet/x/exp/teatest` for full program lifecycle. Golden files in `<module>/testdata/<TestName>.golden`.

### Existing Test Patterns to Follow

From `teamodal/choice_model_test.go`:

```go
// Package-internal tests (not _test package)
package teamodal

// Helper to construct and open a test model
func newTestChoiceModel(options []ChoiceOption, defaultIndex int) ChoiceModel {
    m := NewChoiceModel(&ChoiceModelArgs{
        ScreenWidth:  80,
        ScreenHeight: 24,
        Message:      "Test message",
        Options:      options,
        DefaultIndex: defaultIndex,
    })
    m, _ = m.Open()
    return m
}

// Helper to extract message from command
func extractMsg(cmd tea.Cmd) tea.Msg {
    if cmd == nil {
        return nil
    }
    return cmd()
}

// Test constructs key messages directly
result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyTab})
m = result.(ChoiceModel)  // Type assert back to concrete type
```

### House Rules

- **ClearPath style:** `goto end`, named returns, no else chains
- **doterr errors:** No `fmt.Errorf` — but in tests, standard `t.Errorf`/`t.Fatalf` is fine
- **No `Grep` with `output_mode:"content"`** — ever
- **Package-internal tests** (not `_test` suffix on package name) — matches existing pattern and gives access to unexported fields for setup

### Key v1 API Patterns

All modules use Charm v1 types:
- `tea.KeyMsg{Type: tea.KeyTab}` — special keys
- `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}` — character input
- `tea.MouseMsg{Type: tea.MouseLeft, X: 10, Y: 5}` — mouse events (teamodal only)
- `tea.WindowSizeMsg{Width: 80, Height: 24}` — resize events
- `key.Matches(msg, binding)` — key binding matching
- `View() string` — returns rendered view as string

### Module Dependencies

```
teadep → teadd (imports teadd, has replace directive)
teamodal → teautils (imports teautils, has replace directive)
```

All others are independent. You can test modules in any order, but `test/PLAN.md` recommends: teautils → teadd → teastatus → teamodal → teatree → teadep → teatextsel.

### Adding teatest Dependency

For Layer 3 tests, each module that needs teatest must add it to its `go.mod`:

```bash
cd <module> && go get github.com/charmbracelet/x/exp/teatest@latest
```

teatest requires a model that implements `tea.Model` and can quit. For component models (not top-level programs), you'll need a thin wrapper that embeds the component and adds quit-on-specific-key behavior for the test harness.

### Golden File Convention

- Location: `<module>/testdata/<TestName>.golden`
- Terminal size: 80x24 (fixed for reproducibility via `teatest.WithInitialTermSize(80, 24)`)
- Generate: `go test -update` flag
- Committed to git

## What NOT to Do

- Don't create `teanotify/` or `teagrid/` tests — those modules don't exist yet
- Don't modify any production code — tests only
- Don't change Charm v1 imports to v2 — we're testing the v1 baseline
- Don't create a shared test helper module unless patterns genuinely repeat 3+ times
- Don't add comments/docstrings to existing production code

## Verification

After writing tests for each module:
```bash
cd <module> && go test ./... && go vet ./...
```

After all modules:
```bash
make test && make vet
```

## Key Files to Read First

1. `test/PLAN.md` — The detailed test specifications (your primary reference)
2. `teamodal/choice_model_test.go` — Existing test patterns to follow
3. Each module's `model.go` — The Update() and View() methods you're testing
4. `UPGRADE_V2_PLAN.md` — Context on why these tests matter (section TESTV1)
