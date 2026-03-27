# jediterm-bug

Minimal reproduction of a rendering bug that affects [Bubble Tea](https://github.com/charmbracelet/bubbletea) v2 applications running inside JetBrains GoLand (and all IntelliJ-based IDEs).

## The Problem

GoLand's built-in terminal uses [JediTerm](https://github.com/JetBrains/jediterm), which does not correctly handle all ANSI escape sequences. Specifically, **CBT (Cursor Backward Tab)** sequences emitted by [ultraviolet](https://github.com/charmbracelet/ultraviolet)'s differential renderer can cause invalid screen rendering when transitioning between layouts of different widths.

This program demonstrates the issue by toggling between two views recorded from a real Bubble Tea application:

- **View A**: A two-pane layout (44-col left + 36-col right = 80 cols total)
- **View B**: A single-pane layout (33-col, no right pane)

When ultraviolet diffs these two fundamentally different layouts, it produces escape sequences that JediTerm mishandles, resulting in broken rendering.

## How to Reproduce

### Requirements

- Go 1.25+
- GoLand (or any JetBrains IDE), **or** the standalone JediTerm app [ForceTerm](https://github.com/sebkur/forceterm)

### Steps

1. Build the program:

   ```bash
   go build -o jediterm-bug .
   ```

2. Run it inside GoLand's **Run/Debug** terminal (not an external terminal):

   ```bash
   ./jediterm-bug
   ```

3. Press **space** to toggle between the two views.

The first view (two-pane layout) renders correctly. After pressing space, the second view (single-pane layout) renders incorrectly — portions of the previous layout remain visible or the new layout is misaligned.

Press **q** to quit.

## Diagnosis

The program writes a `trace.log` file next to the executable that captures all terminal output including escape sequences, which can be analyzed with `xxd` or similar tools.

## Context

This reproduction was created to support PRs that add selective capability control to ultraviolet and Bubble Tea:

- **ultraviolet** [[PR #100](https://github.com/charmbracelet/ultraviolet/pull/100)]: `DisableCaps(...Capability)` method and `UV_DISABLE_CAPS` environment variable 
- **Bubble Tea** [[PR #1641](https://github.com/charmbracelet/bubbletea/pull/1641)]: `WithoutCaps(...uv.Capability)` program option
- **ultraviolet** [[PR #101](https://github.com/charmbracelet/ultraviolet/pull/101)]: Auto-fix CBT use in JediTerm for GoLand/JetBrains' Run/Debug TTY

These allow applications to disable specific terminal optimizations (e.g. CBT) for terminals that don't support them, without degrading the entire rendering pipeline.

## How the Frames Were Captured

The two embedded frames (`frameA` and `frameB`) were recorded from a live [gomion](https://github.com/mikeschinkel/gomion) session using a `RecordingModel` wrapper that captures every `View()` output as JSON. The recording was then whittled down to the two-frame transition that triggers the bug.
