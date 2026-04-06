# Research Prompt: Testing Strategies for Go TUI Apps Built with Bubble Tea

## Context

I am building a complex Go TUI application called **Gomion** using **Bubble Tea v2** (by Charmbracelet) with a companion component library called **Tea Leaves**. The app has grown complex enough that AI-assisted code changes frequently introduce regressions. I have no tests currently because TUI testing is harder than standard unit testing — but the complexity now demands it.

The app is structured as:
- A Go workspace with multiple modules (`cmd/gomion`, `gommod/`, `test/`)
- Heavy use of Bubble Tea's `Model`/`Update`/`View` pattern
- Multiple TUI views with keyboard navigation, modal overlays, split panes, tree views, and diff viewers
- Components from my `go-tealeaves` library (dropdowns, modals, notifications, layout engine, text selection, grid)
- The app processes keystrokes, renders styled terminal output via LipGloss v2, and manages complex state transitions across views

## What I Need Researched

### 1. Bubble Tea's Built-In Testing Support

- What testing utilities does Bubble Tea v2 provide? (I know there's a `teatest` package — what exactly does it offer?)
- How does `teatest.NewModel()` or equivalent work for sending messages and inspecting state?
- Can you send sequences of `tea.KeyMsg`, `tea.MouseMsg`, and `tea.WindowSizeMsg` programmatically?
- How do you assert on `View()` output — substring matching, golden files, snapshot testing?
- What are the limitations? What can't `teatest` test?
- Are there any newer testing features in Bubble Tea v2 that weren't in v1?
- Search the Bubble Tea v2 GitHub repo (github.com/charmbracelet/bubbletea) for test examples, especially in `teatest/` and any `_test.go` files

### 2. Component-Level (Unit) Testing Patterns

- Best practices for testing individual Bubble Tea `Model` components in isolation
- How to test `Update()` state transitions: send a message, check resulting model state
- How to test `View()` output: strategies for asserting rendered content without brittle string matching
- How to mock or stub sub-models when testing a parent that composes multiple child models
- How to test components that depend on terminal size (`tea.WindowSizeMsg`)
- How to test async operations (commands returned from `Update()` that produce messages)
- Testing keyboard navigation sequences (e.g., "press j three times, then Enter, then Esc")
- Testing focus management across composed components

### 3. Integration Testing — Multi-Component Interactions

- Strategies for testing full view compositions (e.g., a split-pane view with a tree on the left and diff viewer on the right)
- How to test modal overlay interactions (modal appears, captures keys, returns result to parent)
- Testing Tab-based focus cycling between panes
- Testing state propagation: action in component A triggers update in component B
- Testing the full `Init()` → multiple `Update()` cycles → `View()` pipeline

### 4. End-to-End Testing of Compiled TUI Apps

- **How to test a compiled Go TUI binary end-to-end** — spawning the process, sending keystrokes to its PTY, and asserting on terminal output
- Tools and libraries for PTY-based testing in Go:
  - `github.com/Netflix/go-expect` — is this still maintained? How does it work?
  - `github.com/hinshun/vt10x` — terminal emulator for test assertions
  - `github.com/creack/pty` — raw PTY allocation
  - `charmbracelet/x/xpty` — Charm's own PTY wrapper (I already plan to use this for my `teaterm` component)
  - `charmbracelet/x/vt` — Charm's VT emulator (also already in my stack)
  - Any other PTY/terminal testing libraries in the Go ecosystem
- How to handle timing: waiting for the TUI to render before sending next keystroke
- How to parse ANSI-styled terminal output for assertions (stripping escape sequences vs. parsing them)
- How to handle terminal size configuration in tests
- Golden file / snapshot testing for full terminal screens
- How does `charmbracelet/vhs` (the GIF recorder) work internally? Could its approach to scripting terminal interactions be adapted for testing?

### 5. Scripted Terminal Interaction Approaches

- **VHS by Charmbracelet** (github.com/charmbracelet/vhs): It scripts terminal interactions via `.tape` files. Could this be used or adapted for testing? What's its architecture?
- **Expect-style testing**: Are there Go equivalents of the Unix `expect` tool for scripting interactive terminal sessions?
- **Playwright/Cypress-style approaches**: Are there any tools that provide a higher-level DSL for terminal UI testing, analogous to how Playwright tests web UIs?
- **Record-and-replay**: Any tools that can record a terminal session (keystrokes + output) and replay it for regression testing?

### 6. Testing Strategies from the Broader TUI Ecosystem

- How do other TUI frameworks handle testing?
  - **Textual** (Python) — has a testing framework; what patterns does it use?
  - **tui-rs / Ratatui** (Rust) — any testing patterns from the Rust TUI ecosystem?
  - **blessed / ink** (Node.js) — testing approaches?
- Are there universal patterns that apply regardless of framework?
- How do terminal emulator projects themselves test rendering correctness? (e.g., how does `xterm.js` or `alacritty` test?)

### 7. Claude Computer Use for TUI Testing

- **Claude's computer use capability**: Claude can now control a computer (keyboard, mouse, screenshots). Could this be used to interact with a running TUI app for testing or validation?
- What are the mechanics? Does Claude spawn a virtual desktop, take screenshots, and send input events?
- Could Claude Computer Use be scripted to: launch a TUI app in a terminal, send keystrokes, take a screenshot, and assert that the visual output looks correct?
- What are the latency and reliability characteristics? Is this practical for regression testing, or only for exploratory/manual QA?
- Are there APIs or programmatic ways to use Claude Computer Use, or is it only available interactively?
- How does this compare to traditional PTY-based testing in terms of reliability and speed?

### 8. Practical Implementation Plan

Given all the above research, recommend a **layered testing strategy** for my specific situation:

1. **Layer 1 — Component unit tests**: What tools and patterns for testing individual Tea Leaves components and Gomion view models in isolation
2. **Layer 2 — Integration tests**: What approach for testing composed views with multiple interacting components
3. **Layer 3 — E2E tests**: What approach for testing the full compiled `gomion` binary
4. **Layer 4 — Visual regression**: Whether and how to implement screenshot/golden-file testing for rendered TUI output

For each layer, specify:
- Recommended Go libraries/tools
- Example test structure (pseudocode or real Go code)
- What types of bugs each layer catches
- Rough effort to implement
- Any gotchas specific to Bubble Tea v2

### 9. Specific Gotchas and Practical Concerns

- **ANSI escape sequences in test output**: How to handle colored/styled output in assertions
- **Terminal state**: Bubble Tea enters alt-screen mode — how does this affect testing?
- **Race conditions**: Bubble Tea uses goroutines for commands — how to avoid flaky tests?
- **CI/CD**: Can these tests run in headless CI environments (GitHub Actions, etc.) without a real terminal?
- **Test speed**: What's realistic for test execution time across hundreds of TUI interaction tests?
- **Bubble Tea v2 specifics**: Any testing differences between BT v1 and v2 that I should know about?

## Output Format

Please provide:
1. A comprehensive survey of available tools and approaches, with links to repos and docs
2. Concrete code examples where possible (Go code preferred)
3. A recommended implementation plan prioritized by value-to-effort ratio
4. An honest assessment of what's mature/reliable vs. experimental/fragile in this space
5. Any notable gaps — things that are hard or impossible to test in TUI apps today
