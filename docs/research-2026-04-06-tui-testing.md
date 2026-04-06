# Testing Bubble Tea v2 TUI Apps: A Comprehensive Strategy

**Bubble Tea's Elm Architecture makes TUI apps more testable than most UI frameworks**, but the testing ecosystem remains young and fragmented. The recommended approach is a four-layer pyramid: direct model unit tests (fast, free, deterministic), teatest-based integration tests, PTY + virtual terminal E2E tests for the compiled binary, and optional golden-file visual regression. The tooling exists today to implement all four layers, though some pieces—particularly `charmbracelet/x/vt` integration with teatest—are still experimental. This report provides concrete Go code, library recommendations, and an implementation plan prioritized by value-to-effort ratio.

---

## 1. Bubble Tea's Built-In Testing Support

### teatest Package

The `teatest` package lives at `github.com/charmbracelet/x/exp/teatest/v2` for Bubble Tea v2. It wraps a full `tea.Program` in a test harness, letting you send messages, wait for output conditions, and assert on final model state. It handles the event loop, goroutines, and timing internally.

Core API surface:

| Function | Purpose |
|----------|---------|
| `teatest.NewTestModel(t, model, opts...)` | Create a test program wrapping your model |
| `tm.Send(msg tea.Msg)` | Send any message (key, mouse, window size, custom) |
| `tm.Type("hello")` | Type a string as individual key events |
| `teatest.WaitFor(t, reader, condFn, opts...)` | Block until output matches a condition function |
| `tm.FinalModel(t, opts...)` | Quit program and return final model state |
| `tm.FinalOutput(t)` | Get all accumulated output bytes |
| `teatest.RequireEqualOutput(t, bytes)` | Compare output against golden file |
| `teatest.WithInitialTermSize(w, h)` | Set terminal dimensions for the test |

You can send `tea.KeyPressMsg`, `tea.MouseClickMsg`, `tea.WindowSizeMsg`, and any custom message type programmatically. Assertion on `View()` output supports substring matching via `WaitFor`, golden file comparison via `RequireEqualOutput`, and snapshot testing via the companion `golden` package.

### Key Limitations

teatest has real limitations that surface quickly in complex apps:

**No virtual terminal integration.** Output is a raw byte stream, not rendered frames. Alt-screen apps produce ANSI-encoded output that's hard to assert against directly. Charm has draft PRs (#268, #250) to integrate `charmbracelet/x/vt` into teatest, which would enable proper screen-buffer access — this is the single biggest missing piece.

**No `WaitForString` convenience.** You must write `func(bts []byte) bool` for every condition. A PR for `WaitForMsg` exists but hasn't merged.

**Polling-based timing.** `WaitFor` polls every 50ms by default. On slow CI, tests can flake. Always set explicit `WithDuration` and `WithCheckInterval`.

**Golden file quirks.** The `-update` flag is registered globally (conflicts if not all test files define it). Golden files can break across platforms due to line endings. Add `*.golden -text` to `.gitattributes`.

**Color profile differences.** Golden files generated locally may fail in CI. Use BT v2's `tea.WithColorProfile(colorprofile.Ascii)` to force consistent output, or always strip ANSI before golden comparison.

**Still experimental.** The `x/exp/` import path signals no stability guarantees.

### BT v2 Testing Features Not in v1

BT v2 introduced `tea.WithWindowSize(w, h)` for setting terminal size without a real terminal, `tea.WithColorProfile()` for forcing a color profile in tests, and `Init()` now returns `(tea.Model, tea.Cmd)` instead of just `tea.Cmd`. The key message types were also restructured: `tea.KeyMsg` became `tea.KeyPressMsg`, and `tea.MouseMsg` was split into `tea.MouseClickMsg`, `tea.MouseReleaseMsg`, and others.

---

## 2. Component-Level (Unit) Testing Patterns

The Elm Architecture's greatest testing advantage is that `Update()` is essentially a pure function: message in, model + command out. No terminal, no goroutines, no timing. These tests run in microseconds.

### Testing Update() State Transitions

```go
package gomtui_test

import (
    "testing"
    tea "charm.land/bubbletea/v2"
)

func TestTreeViewCursorDown(t *testing.T) {
    m := NewTreeViewModel(testItems)

    // BT v2 uses tea.KeyPressMsg (not tea.KeyMsg)
    updated, cmd := m.Update(tea.KeyPressMsg{Code: 'j'})
    result := updated.(TreeViewModel)

    if result.Cursor != 1 {
        t.Errorf("cursor: want 1, got %d", result.Cursor)
    }
    if cmd != nil {
        t.Errorf("unexpected command returned")
    }
}
```

### Table-Driven Tests for Key Sequences

```go
func TestKeySequence_NavigateAndSelect(t *testing.T) {
    var m tea.Model = NewTreeViewModel(testItems)

    keys := []tea.KeyPressMsg{
        {Code: 'j'}, {Code: 'j'}, {Code: 'j'},
        {Code: tea.KeyEnter},
        {Code: tea.KeyEscape},
    }
    for _, key := range keys {
        m, _ = m.Update(key)
    }

    result := m.(TreeViewModel)
    if result.Selected != 3 {
        t.Errorf("selected: want 3, got %d", result.Selected)
    }
    if result.ModalOpen {
        t.Error("modal should be closed after Esc")
    }
}
```

### Table-Driven Focus Cycling

```go
func TestSplitPaneFocus(t *testing.T) {
    tests := []struct {
        name    string
        key     rune
        wantFoc int
    }{
        {"Tab cycles forward", tea.KeyTab, 1},
        {"Shift+Tab cycles back", tea.KeyShiftTab, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            m := NewSplitPane(treeView, diffViewer)
            m.Focused = 0
            updated, _ := m.Update(tea.KeyPressMsg{Code: tt.key})
            if updated.(SplitPane).Focused != tt.wantFoc {
                t.Errorf("focus: want %d, got %d", tt.wantFoc, updated.(SplitPane).Focused)
            }
        })
    }
}
```

### Testing View() Output Without Brittleness

Three strategies, from least to most rigid:

```go
import "github.com/charmbracelet/x/ansi"

// Strategy 1: Semantic substring checks (recommended default)
func TestDiffViewerShowsChanges(t *testing.T) {
    m := NewDiffViewer(testDiff)
    m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

    output := ansi.Strip(m.View().String())

    if !strings.Contains(output, "+added line") {
        t.Error("expected added line in diff output")
    }
    if !strings.Contains(output, "-removed line") {
        t.Error("expected removed line in diff output")
    }
}

// Strategy 2: Regex for structural assertions
func TestTreeViewCursorMarker(t *testing.T) {
    m := NewTreeViewModel(testItems)
    output := ansi.Strip(m.View().String())

    re := regexp.MustCompile(`(?m)^▸`)
    if matches := re.FindAllString(output, -1); len(matches) != 1 {
        t.Errorf("expected exactly 1 cursor marker, got %d", len(matches))
    }
}

// Strategy 3: Golden file snapshot (best for complex layouts)
import "github.com/charmbracelet/x/exp/golden"

func TestSplitPaneLayout(t *testing.T) {
    m := NewSplitPane(treeView, diffViewer)
    m, _ = m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
    golden.RequireEqual(t, ansi.Strip(m.View().String()))
}
// First run:  go test -update  (creates testdata/TestSplitPaneLayout.golden)
// Later runs: go test           (validates against golden file)
```

### Testing Commands Returned from Update()

Commands are `func() tea.Msg` — execute them synchronously in tests:

```go
func TestFetchTriggersDataLoad(t *testing.T) {
    m := NewContentModel()
    updated, cmd := m.Update(FetchRequestMsg{ID: "abc"})

    if cmd == nil {
        t.Fatal("expected a command from fetch request")
    }

    // Execute command synchronously — no goroutines, no races
    msg := cmd()

    // Feed result back into Update
    final, _ := updated.Update(msg)
    result := final.(ContentModel)
    if result.Data == nil {
        t.Error("expected data after fetch completes")
    }
}
```

### Testing with Terminal Dimensions

```go
func TestModalCentersInViewport(t *testing.T) {
    m := NewModalOverlay("Confirm delete?")
    m, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

    output := ansi.Strip(m.View().String())
    lines := strings.Split(output, "\n")

    if strings.Contains(lines[0], "Confirm") {
        t.Error("modal should be vertically centered, not at top")
    }
}
```

### Mocking Sub-Models

When testing a parent that composes Tea Leaves components (e.g., `teamodal`, `teadrpdwn`), you can substitute a stub model that records received messages:

```go
type StubModel struct {
    Messages []tea.Msg
    ViewText string
}

func (s StubModel) Init() (tea.Model, tea.Cmd) { return s, nil }
func (s StubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    s.Messages = append(s.Messages, msg)
    return s, nil
}
func (s StubModel) View() string { return s.ViewText }
```

This avoids pulling in the full component dependency tree during unit tests and lets you verify the parent sends the correct messages to children.

**What this layer catches**: Logic bugs in state transitions, cursor math, focus cycling, keyboard handling, command generation, and layout calculations. This is where 80%+ of your test value lives.

**Effort**: Low. No dependencies beyond the standard library and `charmbracelet/x/ansi`. Start here.

---

## 3. Integration Testing — Multi-Component Interactions

### teatest for Full View Compositions

```go
import "github.com/charmbracelet/x/exp/teatest/v2"

func TestFullAppFlow(t *testing.T) {
    m := NewGomionApp(testConfig)
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(120, 40),
    )

    // Wait for initial render
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("gomion"))
    }, teatest.WithDuration(3*time.Second))

    // Navigate tree: down, down, enter
    tm.Send(tea.KeyPressMsg{Code: 'j'})
    tm.Send(tea.KeyPressMsg{Code: 'j'})
    tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})

    // Wait for diff viewer to populate
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("diff"))
    }, teatest.WithDuration(3*time.Second))

    // Open modal with 'd'
    tm.Send(tea.KeyPressMsg{Code: 'd'})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Delete?"))
    })

    // Dismiss modal
    tm.Send(tea.KeyPressMsg{Code: tea.KeyEscape})

    // Quit
    tm.Send(tea.KeyPressMsg{Code: 'q'})

    // Verify final model state
    fm := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second))
    final := fm.(GomionApp)
    if final.ModalOpen {
        t.Error("modal should be closed")
    }
}
```

### Testing Modal Overlay Interactions

The key pattern for testing Tea Leaves overlay components (`teamodal`, `teanotify`, `teadrpdwn`) is verifying that the modal captures keyboard input while open and returns results to the parent on close:

```go
func TestModalCapturesKeys(t *testing.T) {
    m := NewGomionApp(testConfig)
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(120, 40),
    )

    // Trigger modal
    tm.Send(tea.KeyPressMsg{Code: 'd'})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("Confirm"))
    })

    // Press 'j' while modal is open — should NOT move tree cursor
    tm.Send(tea.KeyPressMsg{Code: 'j'})
    time.Sleep(100 * time.Millisecond)

    // Confirm modal
    tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})

    tm.Send(tea.KeyPressMsg{Code: 'q'})
    fm := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second))
    final := fm.(GomionApp)

    // Tree cursor should not have moved while modal was open
    if final.TreeCursor != 0 {
        t.Errorf("tree cursor moved during modal: got %d", final.TreeCursor)
    }
}
```

### Testing Tab Focus Cycling

```go
func TestTabCyclesBetweenPanes(t *testing.T) {
    m := NewGomionApp(testConfig)
    tm := teatest.NewTestModel(t, m,
        teatest.WithInitialTermSize(120, 40),
    )

    // Select a module to enter split-pane view
    tm.Send(tea.KeyPressMsg{Code: tea.KeyEnter})
    teatest.WaitFor(t, tm.Output(), func(bts []byte) bool {
        return bytes.Contains(bts, []byte("diff"))
    }, teatest.WithDuration(3*time.Second))

    // Tab should cycle focus
    tm.Send(tea.KeyPressMsg{Code: tea.KeyTab})
    time.Sleep(100 * time.Millisecond)
    tm.Send(tea.KeyPressMsg{Code: 'q'})

    fm := tm.FinalModel(t, teatest.WithFinalTimeout(2*time.Second))
    final := fm.(GomionApp)
    if final.FocusedPane != 1 {
        t.Errorf("expected focus on pane 1 after Tab, got %d", final.FocusedPane)
    }
}
```

### Testing State Propagation Across Components

```go
func TestTreeSelectionUpdatesDiffViewer(t *testing.T) {
    m := NewSplitPaneView(tree, diffViewer)
    // Select item in tree
    m, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

    result := m.(SplitPaneView)
    if result.DiffViewer.CurrentFile != result.Tree.SelectedFile {
        t.Error("diff viewer should show the file selected in tree")
    }
}
```

**What this layer catches**: Cross-component message passing bugs, race conditions in command execution, full lifecycle issues (Init → Update → View), program startup/shutdown behavior.

**Effort**: Medium. Requires `teatest/v2` dependency. Tests run in 1–3 seconds each due to polling. Best for 10–30 critical-path scenarios, not hundreds.

---

## 4. End-to-End Testing of the Compiled Binary

### Recommended Toolchain

| Tool | Role | Status |
|------|------|--------|
| **`creack/pty`** v1.1.24 | PTY allocation, process spawning | Stable, actively maintained |
| **`charmbracelet/x/vt`** | VT emulator, screen buffer, ANSI parsing | Experimental, actively maintained by Charm team |
| **`charmbracelet/x/xpty`** | Cross-platform PTY (adds Windows via ConPTY) | Experimental, already in your stack for `teaterm` |

Avoid `Netflix/go-expect` and `hinshun/vt10x` — both unmaintained since 2022 and superseded by the Charm tools.

### E2E Test Pattern

```go
package e2e_test

import (
    "os/exec"
    "strings"
    "testing"
    "time"

    "github.com/creack/pty"
    "github.com/charmbracelet/x/vt"
)

func TestGomionBinaryE2E(t *testing.T) {
    // Build the binary
    binPath := t.TempDir() + "/gomion"
    build := exec.Command("go", "build", "-o", binPath, "./cmd/gomion")
    if err := build.Run(); err != nil {
        t.Fatalf("build failed: %v", err)
    }

    // Spawn with PTY at specific size
    cmd := exec.Command(binPath)
    cmd.Env = append(cmd.Environ(), "TERM=xterm-256color")
    ptmx, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: 40, Cols: 120})
    if err != nil {
        t.Fatalf("pty start failed: %v", err)
    }
    defer ptmx.Close()
    defer cmd.Process.Kill()

    // Create VT emulator
    emu := vt.NewEmulator(120, 40)
    defer emu.Close()

    // Pipe PTY output to emulator
    go func() {
        buf := make([]byte, 4096)
        for {
            n, readErr := ptmx.Read(buf)
            if readErr != nil {
                return
            }
            emu.Write(buf[:n])
        }
    }()

    // Wait for initial render
    waitForScreen(t, emu, "gomion", 5*time.Second)

    // Send keystrokes: navigate tree
    ptmx.Write([]byte("j"))  // down
    ptmx.Write([]byte("j"))  // down
    ptmx.Write([]byte("\r")) // enter (select)

    // Wait for diff viewer to show content
    waitForScreen(t, emu, "@@", 3*time.Second)

    // Test Tab focus cycling
    ptmx.Write([]byte("\t")) // Tab to right pane
    time.Sleep(100 * time.Millisecond)

    // Quit
    ptmx.Write([]byte("q"))
    cmd.Wait()
}

func waitForScreen(t *testing.T, emu *vt.Emulator, text string, timeout time.Duration) {
    t.Helper()
    deadline := time.Now().Add(timeout)
    for time.Now().Before(deadline) {
        if strings.Contains(emu.String(), text) {
            return
        }
        time.Sleep(50 * time.Millisecond)
    }
    t.Fatalf("timeout waiting for %q in screen output:\n%s", text, emu.String())
}
```

### Key Escape Sequences for Keystrokes

```go
var keyMap = map[string]string{
    "Enter":     "\r",
    "Escape":    "\x1b",
    "Tab":       "\t",
    "Backspace": "\x7f",
    "Up":        "\x1b[A",
    "Down":      "\x1b[B",
    "Right":     "\x1b[C",
    "Left":      "\x1b[D",
    "Ctrl+C":    "\x03",
    "Ctrl+D":    "\x04",
    "Ctrl+S":    "\x13",
    "Ctrl+L":    "\x0c",
    "Alt+1":     "\x1b1",
    "Alt+A":     "\x1ba",
    "PgUp":      "\x1b[5~",
    "PgDn":      "\x1b[6~",
    "Home":      "\x1b[H",
    "End":       "\x1b[F",
    "ShiftTab":  "\x1b[Z",
}
```

### Event-Driven Synchronization with VT Damage Callbacks

For tighter synchronization than polling, `charmbracelet/x/vt` supports damage callbacks:

```go
func waitForScreenChan(emu *vt.Emulator, text string, timeout time.Duration) error {
    done := make(chan struct{}, 1)
    emu.SetCallbacks(vt.Callbacks{
        Damage: func(d vt.Damage) {
            if strings.Contains(emu.String(), text) {
                select {
                case done <- struct{}{}:
                default:
                }
            }
        },
    })

    select {
    case <-done:
        return nil
    case <-time.After(timeout):
        return fmt.Errorf("timeout waiting for %q", text)
    }
}
```

In practice, the polling approach with `waitForScreen` is simpler and sufficient for most cases.

**What this layer catches**: Binary compilation issues, terminal rendering bugs, alt-screen problems, ANSI escape handling, PTY-specific behavior, real-world keyboard input processing.

**Effort**: High. Requires building the binary, managing PTY lifecycles, and writing synchronization logic. Reserve for 5–10 critical user journeys. Each test runs in 2–5 seconds.

---

## 5. Scripted Terminal Interaction Approaches

### VHS by Charmbracelet

VHS (`github.com/charmbracelet/vhs`) scripts terminal interactions via `.tape` files. Internally it uses ttyd to expose a terminal in a browser, then drives it with a headless Chromium instance, capturing screenshots and GIFs.

```tape
Output testdata/golden/main_view.ascii
Set Width 120
Set Height 40
Set Shell gomion
Sleep 2s
Type "jj"
Enter
Sleep 1s
```

VHS can generate `.ascii` golden files for CI comparison. However, it's heavyweight (requires ttyd + headless Chromium + ffmpeg) and slow (10–30 seconds per test). Use it for documentation screenshots and a handful of critical visual regression tests, not as a primary testing tool.

VHS can be run in CI via `charmbracelet/vhs-action`.

### Expect-Style Testing in Go

`creack/pty` + `charmbracelet/x/vt` gives you the same capabilities as Unix `expect` but with Go's type system and test infrastructure. The E2E pattern in Section 4 is essentially a Go-native expect implementation. There is no need for a separate expect-style library.

### Playwright/Cypress-Style DSL

No established tool provides a high-level DSL for terminal UI testing analogous to Playwright for web UIs. The closest approaches are VHS `.tape` files (declarative but limited) and `knz/catwalk` (data-driven fixture files, described below). Building a thin DSL wrapper around the E2E pattern in Section 4 would be straightforward if you need one.

### Record and Replay

No mature Go tool records terminal sessions for replay-based regression testing. `asciinema` records terminal sessions but doesn't replay them for assertion. The practical alternative is golden-file testing: record expected output once, assert against it on subsequent runs.

---

## 6. Testing Strategies from the Broader TUI Ecosystem

### Textual (Python)

Textual has the most mature TUI testing story. Its `Pilot` class runs the app headlessly, sends input events, and provides `pytest-textual-snapshot` for SVG-based visual regression with HTML diff reports. Key patterns worth borrowing:

**Explicit `pause()`** — waits for all pending messages to process before asserting. The Bubble Tea equivalent is executing all returned commands synchronously before checking model state.

**CSS-selector-based widget targeting** — Bubble Tea lacks this entirely. The equivalent is testing individual components by type assertion on the model, or using `ansi.Strip` + substring matching on `View()` output.

**SVG snapshot testing** with visual HTML diff reports — nothing equivalent exists for Bubble Tea. The closest is golden file text diffs, which show raw text differences rather than visual side-by-side comparisons.

### Ratatui (Rust)

Ratatui provides `TestBackend`, an in-memory terminal buffer with `assert_buffer_lines()` for exact line-by-line comparison. The newer `ratatui-testlib` crate adds PTY-based integration testing. Their architecture maps cleanly to the Bubble Tea approach: PTY → terminal emulation → test harness → snapshot → widget assertions.

The key difference is that Ratatui's `TestBackend` gives you a cell-level buffer (character + style at each position), while Bubble Tea's `View()` gives you a styled string. The `charmbracelet/x/vt` emulator bridges this gap by providing cell-level access for E2E tests.

### Ink (Node.js)

Ink uses `ink-testing-library` with a `render()` function that returns `lastFrame()` (rendered string) and `frames` (all historical frames). The frame-history pattern is useful — Bubble Tea's teatest provides `Output()` as a continuous stream but not discrete frames.

### Universal Patterns

Across all frameworks, three patterns appear consistently: (1) headless mode that bypasses the real terminal, (2) golden file / snapshot comparison for rendered output, and (3) programmatic input event injection. Bubble Tea supports all three, though the headless mode (teatest) is less mature than Textual's or Ratatui's equivalents.

---

## 7. Claude Computer Use for TUI Testing

Claude Computer Use works by controlling a virtual desktop: it takes screenshots, analyzes them visually, and sends keyboard/mouse events. It can launch a terminal, run a TUI app, send keystrokes, and visually inspect the result.

**Practical assessment**: Claude Computer Use is too slow and expensive for automated regression testing. Each interaction cycle (screenshot → analysis → input → wait → screenshot) takes 3–10 seconds and costs approximately $0.03–0.06 per step with Sonnet. For a suite of 100 test scenarios, that's 5–15 minutes and $3–6 per run. Traditional PTY-based tests cover the same ground in seconds for free.

**Where it has value**: Exploratory QA during development and test generation. Claude can interact with the running TUI, discover edge cases, and then write traditional teatest or unit tests that capture those cases. It's a test-authoring aid, not a test-running tool.

### Claude Code + tmux (More Practical Alternative)

The more practical AI-assisted approach is Claude Code driving your TUI through tmux. Claude Code launches the app in a tmux session, sends keystrokes via `tmux send-keys`, and captures rendered output via `tmux capture-pane -p -e`. Claude analyzes the text output directly, which is far more reliable than screenshot analysis.

Uses include test generation (Claude Code interacts with the TUI, then writes tests capturing observed behavior), exploratory QA, and first-pass validation of new features. Not suitable for automated regression suites due to cost and latency.

---

## 8. Alternative: Data-Driven Testing with knz/catwalk

The `knz/catwalk` library (`github.com/knz/catwalk`) offers a data-driven approach using fixture files instead of Go code:

```go
func TestViewport(t *testing.T) {
    m := NewViewportModel(40, 3)
    catwalk.RunModel(t, "testdata/viewport_tests", m)
}
```

Test fixture file (`testdata/viewport_tests/test1`):
```
run
----
-- view:
first line
second line
third line

# Navigate down
run
type j
----
-- view:
second line
third line
fourth line
```

Update expected output: `go test . -args -rewrite`. Supports `type`, `key`, `resize`, and `paste` commands. Particularly useful for testing Tea Leaves components with many input/output scenarios where writing Go code for each case is tedious.

---

## 9. Bubble Tea v2 Breaking Changes That Affect Testing

Every test file must account for these API changes:

| v1 | v2 | Testing Impact |
|----|-----|----------------|
| `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("q")}` | `tea.KeyPressMsg{Code: 'q'}` | All key event construction changes |
| `tea.MouseMsg` (single struct) | `tea.MouseClickMsg`, `tea.MouseReleaseMsg`, etc. | Mouse test messages split into types |
| `View() string` | `View() tea.View` | Call `.String()` on view output |
| `Init() tea.Cmd` | `Init() (tea.Model, tea.Cmd)` | Init returns updated model too |
| N/A | `tea.WithWindowSize(w, h)` | New: set size without real terminal |
| N/A | `tea.WithColorProfile(colorprofile.Ascii)` | New: force color profile in tests |

### LipGloss v2 and ANSI in Tests

LipGloss v2 fundamentally changed rendering: `Style.Render()` always emits full TrueColor ANSI. Downsampling happens at the output layer. For tests:

```go
import "github.com/charmbracelet/x/ansi"

// Strip all ANSI for content assertions
clean := ansi.Strip(styledOutput)

// Or use BT v2's program option to force ASCII profile
p := tea.NewProgram(model{},
    tea.WithColorProfile(colorprofile.Ascii),
    tea.WithOutput(&buf),
)
```

---

## 10. Specific Gotchas and Practical Concerns

### ANSI Escape Sequences in Test Output

Always strip ANSI before content assertions using `ansi.Strip()` from `charmbracelet/x/ansi`. For golden files, decide upfront whether goldens include ANSI (more precise but platform-dependent) or exclude it (more portable). Recommendation: strip ANSI for golden files and test styling separately if needed.

### Alt-Screen Mode

Bubble Tea enters alt-screen mode by default. teatest handles this internally, but raw PTY tests will see alt-screen escape sequences in the output stream. The `charmbracelet/x/vt` emulator handles alt-screen correctly — use `emu.IsAltScreen()` to verify the app entered alt-screen, and `emu.String()` reads from the correct buffer automatically.

### Race Conditions

Bubble Tea uses goroutines for commands. In unit tests, avoid races by executing commands synchronously (call `cmd()` directly and feed the result back to `Update`). In teatest, rely on `WaitFor` with appropriate timeouts instead of `time.Sleep`. Always run tests with `-race` flag.

### CI/CD in Headless Environments

Unit tests (Layer 1) need no terminal at all — they work everywhere. teatest (Layer 2) uses internal PTY mechanisms and works in CI without a real TTY. PTY-based E2E tests (Layer 3) work on Linux/macOS CI runners (PTY allocation doesn't need a display). Docker containers may need `/dev/pts` mounted. Set `TERM=xterm-256color` in CI environment.

GitHub Actions runners have no TTY. `tea.NewProgram` without `WithInput(nil)` will fail with `open /dev/tty: no such device or address`. teatest handles this internally.

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    env:
      TERM: xterm-256color
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.25'
      - name: Unit and integration tests
        run: go test -race -v ./...
      - name: Update golden files (manual trigger)
        if: github.event_name == 'workflow_dispatch'
        run: go test -update ./... && git diff --exit-code testdata/
```

### Test Speed Expectations

| Layer | Per-test time | 100 tests |
|-------|--------------|-----------|
| Direct model unit tests | ~10μs | < 1 second |
| teatest integration | 1–3s | 2–5 minutes |
| PTY E2E | 2–5s | 3–8 minutes |
| VHS visual | 10–30s | Not practical at scale |

### Gitattributes for Golden Files

```gitattributes
*.golden -text
testdata/**/*.golden -text
```

---

## 11. Recommended Implementation Plan

### Phase 1: Direct Model Tests (Start Immediately — Highest ROI)

**Tools**: Standard library + `charmbracelet/x/ansi` + `charmbracelet/x/exp/golden`

**Scope**: Every component's `Update()` and `View()` methods. Cover all keyboard handlers, focus transitions, modal open/close, window resize handling. Prioritize the components where AI-assisted changes most frequently introduce regressions.

**For Gomion specifically**, start with: tree view cursor navigation, split-pane focus cycling (Tab/Shift+Tab), diff viewer content rendering, batch assignment key handlers (1–8, Alt+1–Alt+8, double-tap), modal overlay capture/release in `teamodal`, and the commit target selection flow.

**For Tea Leaves components**, test each component (`teanotify`, `teamodal`, `teadrpdwn`, `tealayout`) in isolation with table-driven key sequences and golden files. These become the foundation that Gomion's integration tests build upon.

**Pattern**: Table-driven tests with key sequences. Strip ANSI for content assertions. Golden files for complex layouts.

**Bugs caught**: ~80% of all TUI bugs (logic errors, off-by-one cursor math, focus cycling mistakes, missing keyboard handlers).

**Effort**: 1–2 days for initial test infrastructure, then 10–30 minutes per component.

### Phase 2: teatest Integration Tests (Week 2)

**Tools**: `github.com/charmbracelet/x/exp/teatest/v2`

**Scope**: 10–20 critical user journeys through the full app. Focus on cross-component interactions: tree selection → diff viewer update, modal overlay → parent state change, Tab cycling between panes, batch assignment → commit group state.

**Effort**: 2–3 days. Each test takes ~30 minutes to write and debug.

### Phase 3: E2E Binary Tests (Week 3–4)

**Tools**: `creack/pty` + `charmbracelet/x/vt`

**Scope**: 5–10 smoke tests for the compiled binary. Verify: app starts, renders main view, keyboard navigation works end-to-end, app exits cleanly. These live in Gomion's `test/` module to avoid circular dependencies.

**Effort**: 3–5 days (significant plumbing for PTY management, timing, and synchronization).

### Phase 4: Golden File Visual Regression (Ongoing)

**Tools**: `charmbracelet/x/exp/golden` for View() snapshots; optionally VHS for full-screen captures.

**Scope**: Every view's rendered output at standard terminal sizes (80×24, 120×40).

**Effort**: Minimal incremental effort once Phase 1 infrastructure exists. Add `golden.RequireEqual` to existing View() tests.

---

## 12. Honest Gaps — What Remains Hard or Impossible

**No structured widget queries.** Unlike Textual's CSS selectors, Bubble Tea provides no way to query "find the focused component" or "get all visible list items" from rendered output. You're limited to string matching on `View()` output or inspecting model state directly.

**No frame-discrete testing.** teatest's `Output()` returns a continuous byte stream, not discrete rendered frames. You can't say "assert on the screen after this specific message was processed." The upcoming `x/vt` integration in teatest would fix this.

**Alt-screen output capture is awkward.** teatest captures raw terminal output bytes (including ANSI escape sequences for alt-screen switching), making golden file testing messy. The workaround is testing via model state (`FinalModel`) rather than output.

**No visual diff tooling.** There's no Bubble Tea equivalent of Textual's SVG snapshot diff report. Golden file diffs show raw text differences, not visual side-by-side comparisons.

**Mouse testing is underdeveloped.** While you can send `tea.MouseClickMsg` in unit tests, E2E mouse testing through a PTY requires computing pixel coordinates and sending raw escape sequences — fragile and rarely worth the effort.

**teatest remains experimental.** The `x/exp/` import path means no backward compatibility guarantees. Pin your dependency version.

**Windows CI.** `creack/pty` is Unix-only. For Windows E2E tests, use `charmbracelet/x/xpty` (which wraps ConPTY), but expect rougher edges.

---

## Summary

The Bubble Tea v2 testing ecosystem is functional but immature. Direct model testing is excellent by design — the Elm Architecture makes `Update()` trivially testable. The gap is in the middle and top of the pyramid: teatest works but lacks virtual terminal integration, and E2E tooling requires manual assembly from `creack/pty` and `charmbracelet/x/vt`.

The most impactful move is starting with comprehensive unit tests using table-driven key sequences and golden files, which catches the vast majority of bugs with minimal infrastructure. For Gomion specifically, prioritize testing keyboard navigation across the tree view, split-pane focus cycling, modal overlay capture/release, batch assignment state management, and diff viewer content rendering — these are the components where AI-assisted code changes are most likely to introduce regressions.
