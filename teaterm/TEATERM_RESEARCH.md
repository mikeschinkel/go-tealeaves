# Teaterm: feasibility and architecture for a terminal-in-terminal Bubble Tea v2 component

**A Bubble Tea v2 component that hosts a full terminal emulator is not only feasible — it already exists in production.** TUIOS, a tiling terminal multiplexer built on Bubble Tea v2 and Charm's `x/vt` virtual terminal, demonstrates the complete architecture: PTY management, VT parsing, cell-grid rendering, input forwarding, scrollback, mouse support, and multi-instance sessions with daemon/SSH/web modes. The Charm ecosystem has converged on a stack — ultraviolet primitives, `x/vt` emulator, `x/xpty` cross-platform PTY — that makes teaterm achievable with roughly **300–500 lines of integration code** bridging these existing packages. The primary risk is not feasibility but API instability: every critical dependency (`x/vt`, ultraviolet, Bubble Tea v2 itself) is pre-1.0 and actively evolving.

---

## The Charm stack already provides every layer teaterm needs

The architecture for an embedded terminal component decomposes into five layers, each of which now has a first-party Charm implementation:

**Layer 1 — PTY management** is handled by `charmbracelet/x/xpty`, which wraps `creack/pty` on Unix and Windows ConPTY via `charmbracelet/x/conpty`. The `creack/pty` library (v1.1.24, ~2k stars) is the de facto Go standard for POSIX PTY operations — `Start(cmd)` attaches a subprocess to a PTY and returns the master `*os.File`, while `Setsize()` issues the `TIOCSWINSZ` ioctl for resize. `x/xpty` adds a unified `PTY` interface with a cross-platform `Resize()` method, eliminating the need for platform-specific code. A related library, `aymanbagabas/go-pty` by a core Charm contributor, provides the cleanest API (`pty.New()` → `pty.Command("bash")` → `cmd.Start()`) and could serve as an alternative.

**Layer 2 — VT emulation** is `charmbracelet/x/vt`. The `Emulator` type (renamed from `Terminal` in a recent refactor) implements `io.Writer` — PTY output bytes are fed directly via `Write()` and the emulator parses all escape sequences, maintaining a cell grid internally. The `Emulator` supports **CSI, OSC, DCS, ESC, APC, PM, and SOS** sequence classes, with VT100/VT220/xterm-256color compatibility validated by the vttest CI suite. Key capabilities include alternate screen buffers (`IsAltScreen()`), configurable scrollback (`SetScrollbackSize()`), full cursor management (position, style, visibility, save/restore), **256 indexed colors plus 24-bit true color** via Go's `color.Color` interface, and extensible handler registration (`RegisterCsiHandler`, `RegisterOscHandler`, etc.).

**Layer 3 — Cell primitives** come from `charmbracelet/ultraviolet`, which provides the `Cell`, `Screen`, `Position`, and `LineData` types that both `x/vt` and Bubble Tea v2's Cursed Renderer share. Each `uv.Cell` stores a grapheme, foreground/background colors (as `color.Color`), text attributes (bold, italic, underline, strikethrough, blink, reverse, faint), and hyperlink data. Because `x/vt` and BT v2 share the same cell type, **no conversion layer is needed between emulation and rendering**.

**Layer 4 — Rendering** leverages two paths. `Emulator.Render()` produces a complete ANSI string suitable for returning from `View()`. Alternatively, `Emulator.Draw(scr uv.Screen, area uv.Rectangle)` composites the emulator state directly onto an ultraviolet Screen at a specified rectangle — the same Screen type BT v2's renderer uses internally. Either path feeds into BT v2's **Cursed Renderer**, an ncurses-inspired cell-based diffing algorithm that compares the new frame against the previous one and emits only changed cells, with Mode 2026 (synchronized output) support to eliminate tearing.

**Layer 5 — Input forwarding** is solved by `Emulator.SendKey(KeyPressEvent)`, `SendMouse(Mouse)`, `SendText(string)`, and `Paste(string)`. These methods translate high-level events into the correct byte sequences for the PTY. The `InputPipe() io.Writer` provides a raw byte path for edge cases. `Focus()` and `Blur()` methods support focus event reporting to the child process.

The **damage tracking system** is particularly important for performance. The `Callbacks.Damage` hook fires on every change with typed damage events: `CellDamage` (single cell), `RectDamage` (rectangular region), `ScreenDamage` (full redraw), `ScrollDamage` (scroll event), and `MoveDamage` (region moved). Additionally, `Touched()` returns modified `LineData` entries for line-level dirty checking, and `ClearTouched()` resets the state — enabling a simple cache pattern where `View()` only regenerates when `Touched()` is non-empty.

---

## TUIOS validates the architecture at scale

**TUIOS** (`github.com/Gaurav-Gosain/tuios`) is the strongest feasibility proof. Built on Bubble Tea v2 and LipGloss v2 with `x/vt` as its terminal emulator, it implements a full tiling window manager with **multiple simultaneous terminal sessions**, 9 workspaces, vim-style modal interface, BSP tiling, mouse-driven window management (drag, resize, click), **10,000-line scrollback**, vim copy mode with search, daemon mode with detach/reattach (like tmux), SSH server mode via Wish, and web terminal mode via WebGL. It ships as a Homebrew package, AUR package, and Nix flake. Charm itself promoted it: *"TUIOS is an experimental tiling multiplexer! Built with Bubble Tea and VT, the Charm terminal emulator."*

TUIOS's `getRawKeyBytes` function demonstrates the complete `tea.KeyPressMsg → PTY bytes` conversion, leveraging BT v2's `Key.Code`, `Key.Mod`, and `Key.Text` fields. Its 60Hz polling architecture, style caching, and viewport culling (only rendering visible windows) provide concrete performance benchmarks for teaterm to match or exceed. TUIOS also exposes a library API (`tuios.New()` with functional options) that can be embedded in other BT applications — **exactly the pattern teaterm would follow**.

---

## Go VT emulators compared: x/vt leads decisively

| Capability | `x/vt` | `vito/midterm` | `hinshun/vt10x` | `rcarmo/go-te` | `micro-editor/terminal` |
|---|---|---|---|---|---|
| Color model | 24-bit via `color.Color` | 24-bit via `termenv.Color` | 256 indexed only | Pyte-derived | 256 indexed only |
| Alt screen | ✅ with callback | ✅ (`IsAlt` field) | ✅ (mode flag) | ✅ | ✅ (mode flag) |
| Scrollback | ✅ configurable | ❌ (auto-resize only) | ❌ | ✅ (HistoryScreen) | ❌ |
| Damage tracking | Cell/rect/scroll/move | Per-row change counter | Screen-level flag | Row-level (DiffScreen) | Screen-level flag |
| Input methods | `SendKey`/`SendMouse`/`Paste` | Write bytes only | Write bytes only | `Feed()` strings | Write bytes only |
| Custom handlers | CSI/OSC/DCS/ESC/APC/PM/SOS | OSC forwarding only | ❌ | ❌ | ❌ |
| Thread safety | `SafeEmulator` variant | Unclear | `Lock()`/`Unlock()` | Unclear | `Lock()`/`Unlock()` |
| BT v2 integration | Native (`uv.Cell`, `Draw()`) | Via termenv bridge | None | None | None |
| Maintenance | Active (Charm team, Mar 2026) | Active (Dagger, Oct 2025) | ❌ Unmaintained since 2022 | New, 4 stars | Semi-active (micro-internal) |
| License | MIT | MIT | MIT | MIT | MIT |
| vttest CI | ✅ | ❌ | ❌ | ESCTest2 suite | ❌ |

**`x/vt` is the clear recommendation.** Its native ultraviolet type system eliminates conversion overhead, its `SendKey`/`SendMouse` methods solve input forwarding without manual byte encoding, its damage tracking enables efficient caching, and TUIOS proves it works at scale. The only competitive alternative is `vito/midterm` — production-proven in Dagger with true color via `termenv.Color` — but it lacks native BT v2 integration, scrollback, and custom handler registration. `rcarmo/go-te` has impressive VT220/VT520 standards fidelity through its Python pyte test suite, but is too immature (4 stars, single author) for production use.

---

## Bubbleterm provides a concrete implementation reference

`taigrr/bubbleterm` (v0.2.0, published March 5, 2026, **0BSD license**) already targets Bubble Tea v2 with `View()` returning `tea.View`. Its architecture separates a custom `emulator` sub-package (CSI/OSC/ESC/DCS parsing, cursor/scrollback, 256+true color, line-level damage tracking) from a root `bubbleterm` Model that integrates with BT v2. Three constructors cover different use cases: `New(w, h)` for bare emulator, `NewWithCommand(w, h, cmd)` for PTY-backed shells, and `NewWithPipes(w, h, reader, writer)` for pre-started processes. It implements `Focus()`/`Blur()`/`Focused()`, configurable auto-polling via `SetAutoPoll(bool)`, and manual `UpdateTerminal()` for tick-driven updates.

Key design decisions to learn from or improve upon in teaterm:

- **Pointer receiver model** (`*Model` instead of value type) deviates from BT idiom and requires ugly type assertions in `Update()`. Teaterm should use value-type models or document the trade-off clearly.
- **Raw ANSI rendering** preserves the child process's escape codes faithfully but may conflict with BT v2's automatic color downsampling via `colorprofile`. Using `x/vt`'s `Render()` or `Draw()` would integrate properly with the Cursed Renderer's pipeline.
- **Custom emulator** is the weakest link. The author acknowledges: *"We may decide to use a different emulator library in the future."* Teaterm should build on `x/vt` instead of writing another emulator.
- **`NewWithPipes` is excellent** for embedding already-running processes. Teaterm should offer the same capability — it's the subprocess viewer mode without PTY overhead.

---

## Architectural lessons from BigJk/crt and JediTerm

**BigJk/crt** renders Bubble Tea programs inside Ebitengine graphical windows. Its flat architecture (~10 files) demonstrates the **cell grid as universal interface**: a 2D `[][]Cell` where each cell stores `{rune, fgColor, bgColor, bold, italic, underline}` sits between parsing and rendering. CRT's pipe-based I/O bridge — Bubble Tea writes ANSI to a custom `ReadWriter`, CRT reads and parses it — is the exact inverse of what teaterm needs. Its `Adapter` pattern for bidirectional event translation (Ebitengine keys → BT messages) is directly applicable to teaterm's `tea.KeyPressMsg → PTY bytes` bridge. Most instructively, CRT proves you don't need full VT compliance — it implements only the ~20 sequences Bubble Tea actually emits.

**JediTerm** (`github.com/JetBrains/jediterm`, LGPL/Apache dual-license) demonstrates industrial-grade separation of concerns across five layers: `TtyConnector` (data source) → `JediEmulator` (parser) → `JediTerminal` (state machine, ~60 interface methods) → `TerminalTextBuffer` (buffer with primary/alt screens and scrollback) → `TerminalPanel` (Swing rendering). The critical insight is that **the parser communicates with the state machine via named method calls, not raw escape codes** — making both independently testable. JediTerm's `TextEntry` run-length encoding (grouping consecutive characters with identical styles) is optimal for converting cell grids back to styled strings: walk the grid, group identical-style runs, render each run once.

`ForceTerm` shows minimal JediTerm wiring: create widget, create PTY connector, call `widget.setTtyConnector(connector)` then `widget.start()`. This confirms that with a well-designed emulator library, the integration code is trivially small — exactly what `x/vt`'s API enables.

---

## Input handling, resize, and the escape key problem

BT v2's `tea.KeyPressMsg` carries `Code rune`, `Mod KeyMod` (Ctrl/Alt/Shift/Super flags), and `Text string`. Since BT v2 already runs in raw mode, every keystroke arrives as a message — no OS-level line editing interferes. **`x/vt`'s `SendKey(KeyPressEvent)` accepts events with the same structure**, making conversion straightforward without manual byte encoding. Arrow keys, function keys, Ctrl combinations, and modifier+key sequences are all handled internally by the emulator's input pipeline. Mouse events (`tea.MouseClickMsg`, `tea.MouseReleaseMsg`, `tea.MouseWheelMsg`, `tea.MouseMotionMsg`) require coordinate translation — subtract the teaterm widget's screen offset — before forwarding via `SendMouse()`.

**Resize** flows naturally: `tea.WindowSizeMsg` arrives in `Update()`, teaterm computes the sub-terminal size (component dimensions minus borders/chrome), calls `pty.Setsize()` on the PTY master (which sends `SIGWINCH` to the child process), and `emulator.Resize(w, h)` on the VT emulator. On Windows, `x/xpty`'s `Resize()` calls `ResizePseudoConsole()` transparently.

**The escape key ambiguity** is a genuine design challenge. ESC (`\x1b`) is both a standalone key and the prefix for escape sequences. BT v2's input layer uses a ~50ms timeout disambiguator internally, so `tea.KeyPressMsg` already resolves this before it reaches teaterm. The real question is whether ESC should **exit the teaterm's focus** (like vim's insert mode) or **forward to the PTY** (so programs inside the terminal receive it). TUIOS uses a tmux-style prefix key (`Ctrl+B`) plus a vim-inspired modal system ('i' enters terminal mode, ESC exits to window management mode). This is the recommended pattern: never intercept ESC itself, use a dedicated prefix key or Ctrl combination as the "escape hatch" from terminal mode.

---

## The subprocess viewer should be a separate, lighter component

Not every subprocess needs a PTY. Programs that don't call `isatty()` — simple command output, build logs, test runners — work fine with piped `stdout`/`stderr`. But **most interactive tools change behavior** when they detect a non-TTY: `ls` drops colors and columns, `git` disables its pager, `cargo`/`npm`/`go` suppress progress bars, and anything using readline loses line editing. The threshold is clear: if the user expects the subprocess to behave as if in a terminal, it needs a PTY.

The recommended design offers **two distinct components** rather than modes:

- **`teaterm.Terminal`** — full PTY + `x/vt` emulation for interactive programs (shells, vim, htop, etc.)
- **`teaterm.ProcessViewer`** — piped I/O + ANSI-aware viewport for command output (builds, logs, scripts)

`ProcessViewer` would use `cmd.StdoutPipe()`/`cmd.StderrPipe()`, strip or honor ANSI color codes, and render through Bubbles' existing `viewport` component with automatic scrolling. It needs no VT emulator, no PTY, no resize signaling — dramatically simpler and lighter. `bubbleterm` already validates this split with its `NewWithPipes(w, h, reader, writer)` constructor.

`tea.Exec` is **not relevant to either component**. It pauses the entire BT program and hands the real terminal to the subprocess — the opposite of embedding. It's designed for launching editors or pagers that need the full screen temporarily.

---

## Performance is not a bottleneck at typical terminal sizes

For an **80×24 terminal** (1,920 cells), building the ANSI string representation takes **~50–100μs** — negligible against a 16ms frame budget at 60fps. For **120×40** (4,800 cells), string building takes ~150–300μs, still well within budget. The real optimization comes from caching: check `emulator.Touched()` on each `View()` call, and return the cached string if nothing changed. BT v2's Cursed Renderer then diffs the output at the cell level, emitting only changed cells to the terminal — so even a full 4,800-cell string incurs minimal I/O if only a few cells actually changed.

PTY reads should use **blocking `Read()` in a `tea.Cmd` chain**, not polling. Each Cmd reads one buffer (4–8KB), returns a `ptyOutputMsg`, which triggers the next read in `Update()`. This integrates naturally with BT's message loop and avoids raw goroutines (which BT docs warn against). For rapid output (e.g., `cat large_file`), BT's FPS cap naturally throttles renders — multiple `ptyOutputMsg` messages may arrive between frames, each feeding the emulator, but `View()` only runs once per frame. `bubbleterm`'s alternative approach — a `SetAutoPoll(bool)` toggle with timer-based `UpdateTerminal()` calls — is equally valid for high-throughput scenarios.

The **`SafeEmulator`** variant adds mutex locking for concurrent access, essential if using a separate goroutine for PTY reads rather than pure Cmd chaining. The overhead is minimal for the locking granularity involved.

---

## Proposed architecture with clear layer boundaries

```
┌─────────────────────────────────────────────────┐
│  Bubble Tea v2 Program                          │
│                                                 │
│  ┌───────────────────────────────────────────┐  │
│  │  teaterm.Terminal (Model)                 │  │
│  │                                           │  │
│  │  Init() → startPTY() + readLoop()         │  │
│  │                                           │  │
│  │  Update(tea.KeyPressMsg) ──┐              │  │
│  │  Update(tea.MouseClickMsg) ├─→ SendKey()  │  │
│  │  Update(tea.WindowSizeMsg) │   SendMouse()│  │
│  │         │                  │   Resize()   │  │
│  │         ▼                  ▼              │  │
│  │  ┌──────────┐    ┌─────────────────┐     │  │
│  │  │  x/xpty  │◄──▶│  x/vt.Emulator  │     │  │
│  │  │  master  │    │  (SafeEmulator) │     │  │
│  │  │  fd      │    │                 │     │  │
│  │  └──────────┘    └────────┬────────┘     │  │
│  │       │                   │              │  │
│  │  ptyOutputMsg      Render() or Draw()    │  │
│  │  (tea.Cmd chain)   + Touched() cache     │  │
│  │                           │              │  │
│  │  View() ◄─────────────────┘              │  │
│  │  → tea.View{Content: cachedRender}       │  │
│  └───────────────────────────────────────────┘  │
│                                                 │
│  Cursed Renderer (cell-level diff, Mode 2026)   │
└─────────────────────────────────────────────────┘
```

Public API surface for v1:

- `New(cmd *exec.Cmd, opts ...Option) *Terminal` — create terminal with PTY
- `NewFromPTY(pty io.ReadWriteCloser, opts ...Option) *Terminal` — attach to existing PTY
- `Init() tea.Cmd`, `Update(tea.Msg) (tea.Model, tea.Cmd)`, `View() tea.View`
- `Focus()`, `Blur()`, `Focused() bool`
- `Resize(width, height int)`
- `Close() error`
- Options: `WithSize(w, h)`, `WithScrollback(lines)`, `WithEnv(env []string)`

---

## Minimal viable scope for v1

The v1 should target a **working embedded shell** with these capabilities:

1. PTY creation via `x/xpty` with subprocess spawn (default: user's `$SHELL`)
2. VT emulation via `x/vt.SafeEmulator` with 80×24 default size
3. `View()` returning `emulator.Render()` with `Touched()` caching
4. Input forwarding: `tea.KeyPressMsg` → `emulator.SendKey()`, mouse → `emulator.SendMouse()`
5. Resize: `tea.WindowSizeMsg` → PTY `Setsize()` + `emulator.Resize()`
6. Focus/Blur with a configurable "escape hatch" key (default: Ctrl+\)
7. Lifecycle: clean PTY close, subprocess wait, goroutine cleanup on `Close()`
8. Basic scrollback (1,000 lines default)

**Explicitly out of scope for v1**: clipboard integration, multiple instances in a single model, ProcessViewer mode, custom theming/color remapping, selection/copy, search within scrollback, and serialization.

---

## Blockers, risks, and open questions

**No hard blockers exist.** TUIOS proves every piece works together. However, several risks warrant attention:

- **API instability is the primary risk.** `x/vt` has no tagged release (uses commit-hash versions), recently renamed `Terminal` to `Emulator` and migrated from `cellbuf` to `ultraviolet` types. Pin to a specific commit hash and expect to track upstream changes.
- **Ultraviolet is pre-1.0** and explicitly warns against stability expectations. Since it underpins both `x/vt` and BT v2, breaking changes cascade broadly — but the Charm team's own products (Crush, TUIOS) depend on it, creating strong incentive for stability.
- **Windows support** via ConPTY is architecturally supported by `x/xpty` but less battle-tested than Unix PTY. Empirical testing on Windows is needed before claiming cross-platform support.
- **`tea.View` semantics** need empirical validation. BT v2's `View()` returns a `tea.View` struct with `Content string` plus declarative fields (`AltScreen`, `Cursor`, etc.). How the Cursed Renderer handles pre-formatted ANSI strings in `Content` — especially around color downsampling — needs testing. If the renderer re-parses and downsamples the emulator's ANSI output, colors may shift unexpectedly.
- **Escape hatch design** has no established convention. Tmux uses `Ctrl+B`, screen uses `Ctrl+A`, TUIOS uses a modal system. The choice affects usability significantly and may need to be configurable.
- **Performance under rapid output** (`find /`, `cat huge_file`) needs empirical profiling. The Cmd-chaining pattern introduces one frame of latency per read; under sustained high-throughput output, this may cause visual lag. A channel-based subscription pattern may perform better but is less idiomatic.

**Open questions requiring empirical investigation:**

1. Does `emulator.Render()` output interact correctly with BT v2's color downsampling, or does it need to go through `Draw()` instead?
2. What is the actual frame time for `Render()` at 120×40 with complex styled content (e.g., syntax-highlighted code in vim)?
3. Can `x/vt`'s `SafeEmulator` handle the concurrent access pattern of PTY reads feeding the emulator while `View()` reads the cell grid on BT's render goroutine?
4. How does `x/xpty` behave on Windows with ConPTY for programs expecting Unix PTY semantics (e.g., terminal resizing race conditions)?
5. What is TUIOS's approach to the `View()` rendering path — does it use `Render()`, `Draw()`, or custom cell-grid iteration?

---

## Conclusion

Teaterm should **build on `charmbracelet/x/vt` for emulation and `charmbracelet/x/xpty` for PTY management** — these are first-party Charm packages that share the ultraviolet type system with Bubble Tea v2, eliminating conversion layers and ensuring forward compatibility as the ecosystem evolves. TUIOS's existence as a production terminal multiplexer on this exact stack removes all feasibility doubt. The subprocess viewer should be a **separate component** (`ProcessViewer`) rather than a mode of `Terminal`, because the two share almost no implementation — one needs a full VT emulator and PTY, the other needs only piped I/O and a viewport. Building both under a `teaterm` package with shared option patterns gives users a clean API without overloading a single type. The v1 milestone is achievable with modest integration code: the hard problems (VT parsing, cell rendering, input encoding, PTY management) are already solved by the Charm stack.