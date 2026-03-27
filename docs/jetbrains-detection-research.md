# Detecting JetBrains IDE Processes on Windows and Linux

**No reliable, automatically-injected environment variable exists for processes launched via JetBrains Run/Debug configurations.** The widely-used `TERMINAL_EMULATOR=JetBrains-JediTerm` is set only in the Terminal tool window, not in Run/Debug configs. However, the practical impact is narrower than it appears: JediTerm's CBT bug only matters in two scenarios (the Terminal tool window and the "Emulate terminal" Run/Debug option), and the first scenario already has a clean detection path. The second scenario—and any future edge cases—require either `TOOLBOX_VERSION` detection or parent process inspection as a fallback. Here is a complete breakdown of every detection method, its scope, and a recommended strategy.

## The three execution contexts and why they matter

JetBrains IDEs launch user processes through three distinct code paths, each with different terminal emulation behavior. Understanding these is essential because **CBT (Cursor Backward Tabulation, CSI Z) only matters when JediTerm is actively interpreting the output stream**.

**Standard Run/Debug console** (the default) pipes stdout/stderr directly—no PTY is allocated, `isatty()` returns `false`, and no terminal emulation occurs. A well-behaved terminal library will skip escape sequences entirely in this mode, making CBT detection irrelevant. As JetBrains developer Dmitry Jemerov confirmed in relation to PY-4853: the IDE starts the program as an external process, redirecting its standard input and output and displaying the standard output in the IDE's UI. The same applies to GoLand's Go run configurations.

**Terminal tool window** spawns a full interactive shell via Pty4J with a real PTY. The environment passes through `LocalTerminalDirectRunner` → `LocalOptionsConfigurer`, which injects `TERMINAL_EMULATOR=JetBrains-JediTerm`, `TERM_SESSION_ID`, shell integration variables, and sets `TERM=xterm-256color`. Detection here is trivial.

**"Emulate terminal in output console"** (an opt-in checkbox in Run/Debug configs, available in CLion, PyCharm, GoLand, and some IntelliJ run types) wraps the process output in a JediTerm widget with a real PTY—so `isatty()` returns `true`—but takes an entirely different code path via `TerminalExecutionConsole`. This path does **not** go through `LocalOptionsConfigurer` and **does not set `TERMINAL_EMULATOR`**. This is the problematic gap.

## Confirmed environment variables and their exact scope

The research examined every `IDEA_*`, `JETBRAINS_*`, and `INTELLIJ_*` prefixed variable across JetBrains documentation, the `intellij-community` and `jediterm` GitHub repos, YouTrack, community sources, and hands-on testing of GoLand on Windows.

### Variables present in BOTH Terminal tool window AND Run/Debug configs

**`TOOLBOX_VERSION`** (e.g., `3.4.1.78303`) is injected by JetBrains Toolbox into the IDE process environment and inherited by all child processes—both Terminal tool window and Run/Debug configurations. **Critically, it is NOT a system-level or user-level environment variable.** Testing on Windows confirmed that `echo %TOOLBOX_VERSION%` in a plain `cmd.exe` outside the IDE prints the literal string (unset), while it is present in both the Run/Debug env dump and the Terminal tool window env dump. This makes it a valid detection signal for "running inside a JetBrains IDE launched via Toolbox." Since Toolbox is JetBrains' default recommended installation method, this covers the majority of users.

### Variables present ONLY in Terminal tool window

**`TERMINAL_EMULATOR=JetBrains-JediTerm`** is the gold-standard detection signal for the Terminal tool window. It is set cross-platform (Windows, Linux, macOS) but exclusively in the Terminal tool window. JetBrains implemented this after GitHub issue JetBrains/jediterm#253 (filed October 2022 by @mikehearn, now closed). The variable is injected in the `LocalOptionsConfigurer.kt` code within the Terminal plugin.

**`TERM_SESSION_ID`** (a UUID, e.g., `0e83f686-dca2-4910-a169-26bd4dd3442e`) is set alongside `TERMINAL_EMULATOR` in the Terminal tool window. Not present in Run/Debug configs.

**`INTELLIJ_TERMINAL_COMMAND_BLOCKS=1`** is set in the Terminal tool window (likely related to the "new terminal" block-based UI). Not present in Run/Debug configs.

### Variables NOT set by JetBrains

**`IDEA_INITIAL_DIRECTORY`** and **`JETBRAINS_IDE`** do not exist. Exhaustive searches across GitHub code search, YouTrack, JetBrains documentation, and community forums returned zero results for either variable.

**`INTELLIJ_ENVIRONMENT_READER`** is set transiently during IDE startup when the IDE spawns an interactive shell to capture user environment variables from `~/.bashrc` or `~/.zshrc`. It exists only for that brief subprocess, primarily on macOS and Linux. It is not present in any user-facing session.

**`FIG_TERM=1`**, **`PROCESS_LAUNCHED_BY_CW=1`**, and **`PROCESS_LAUNCHED_BY_Q=1`** are injected by the Amazon Q / CodeWhisperer JetBrains plugin, not by JetBrains itself. They appear in Terminal tool window sessions when the plugin is installed.

### Summary table of tested Windows GoLand environment

| Variable | Run/Debug | Terminal | Source |
|---|---|---|---|
| `TOOLBOX_VERSION` | ✅ Present | ✅ Present | JetBrains Toolbox |
| `TERMINAL_EMULATOR` | ❌ Absent | ✅ `JetBrains-JediTerm` | Terminal plugin |
| `TERM_SESSION_ID` | ❌ Absent | ✅ Present (UUID) | Terminal plugin |
| `INTELLIJ_TERMINAL_COMMAND_BLOCKS` | ❌ Absent | ✅ `1` | Terminal plugin |
| `IDEA_INITIAL_DIRECTORY` | ❌ Does not exist | ❌ Does not exist | — |
| `JETBRAINS_IDE` | ❌ Does not exist | ❌ Does not exist | — |
| `INTELLIJ_ENVIRONMENT_READER` | ❌ Absent | ❌ Absent | Startup-only |

## JediTerm definitively lacks CBT support

Source code analysis of `JediEmulator.java`'s `processControlSequence` method confirms that **case `'Z'` (CBT) is absent from the switch statement**. The handled CSI final characters are: `@` (ICH), `A` (CUU), `B` (CUD), `C` (CUF), `D` (CUB), `E` (CNL), `F` (CPL), `G`/`` ` `` (CHA), `H`/`f` (CUP), `J` (ED), `K` (EL), `L` (IL), `M` (DL), `T` (SD), `c` (DA), `d` (VPA), `g` (TBC), `h`/`l` (SM/RM), `m` (SGR), `q` (DECSCUSR), `r` (DECSTBM). Neither `Z` (CBT) nor `I` (CHT, Cursor Forward Tabulation) is implemented. Unrecognized sequences trigger a logged error and are silently discarded. Since JediTerm reports `TERM=xterm-256color`, whose terminfo entry advertises `cbt=\E[Z`, this creates a **capability mismatch**—the terminal claims CBT support but silently drops it.

## The "Emulate terminal" gap and a proposed fix

The "Emulate terminal in output console" Run/Debug option is the critical gap. In this mode:

- A real PTY is allocated via Pty4J, so `isatty()` returns `true`
- JediTerm is actively interpreting the output stream
- CBT sequences will be silently discarded
- But `TERMINAL_EMULATOR` is **not set** because `TerminalExecutionConsole` takes a different code path than `LocalOptionsConfigurer`

**This is technically trivial to fix on JetBrains' side.** Pty4J's `PtyProcessBuilder` accepts environment variables via its `setEnvironment()` method before spawning the child process. JetBrains already uses this exact mechanism in the Terminal tool window code path. The `TerminalExecutionConsole` code path simply needs to make the same call.

### Proposed YouTrack feature request

**Title:** Set `TERMINAL_EMULATOR=JetBrains-JediTerm` when "Emulate terminal in output console" is enabled in Run/Debug configurations

**Product:** IntelliJ Platform (affects all JetBrains IDEs: GoLand, PyCharm, CLion, IntelliJ IDEA, WebStorm, etc.)

**Type:** Feature Request

**Description:**

When "Emulate terminal in output console" is enabled in a Run/Debug configuration, the child process is spawned via Pty4J with JediTerm interpreting the output stream. However, unlike the Terminal tool window, no `TERMINAL_EMULATOR` environment variable is set in this code path.

The Terminal tool window correctly sets `TERMINAL_EMULATOR=JetBrains-JediTerm` via `LocalOptionsConfigurer`, but `TerminalExecutionConsole` (used by "Emulate terminal in output console") does not make an equivalent `setEnvironment()` call on the `PtyProcessBuilder`.

**Why this matters:**

JediTerm does not implement all terminal capabilities advertised by its `TERM=xterm-256color` setting. For example, CBT (Cursor Backward Tabulation, CSI Z) is absent from `JediEmulator.processControlSequence()` — the sequence is silently discarded. Terminal rendering libraries need to detect JediTerm to disable unsupported capabilities, but there is currently no environment signal in the "Emulate terminal" Run/Debug path to enable this detection.

On macOS, `__CFBundleIdentifier` provides an indirect signal. On Windows and Linux, there is no reliable automatic detection method for this code path.

**Proposed fix:**

In `TerminalExecutionConsole`, set `TERMINAL_EMULATOR=JetBrains-JediTerm` on the `PtyProcessBuilder` environment before spawning the child process, consistent with what `LocalOptionsConfigurer` already does for the Terminal tool window.

This is a one-line change — the Pty4J API already supports it via `PtyProcessBuilder.setEnvironment()`, and the precedent is established in the Terminal plugin code.

**Related issues:**

- PY-4853 — PyCharm should set an environment variable to allow executed scripts to detect it
- JetBrains/jediterm#253 — Set the TERMINAL_EMULATOR environment variable (resolved for Terminal tool window only)

## Parent process inspection: fallback for standalone IDE installs

For users who install JetBrains IDEs without Toolbox (where `TOOLBOX_VERSION` is absent), **walking the parent process tree** is the remaining automatic detection method.

On **Windows**, the direct parent process is the IDE launcher executable: `idea64.exe` (IntelliJ IDEA), `goland64.exe` (GoLand), `pycharm64.exe` (PyCharm), `clion64.exe` (CLion), `webstorm64.exe` (WebStorm), `rider64.exe` (Rider), `phpstorm64.exe` (PhpStorm), `rubymine64.exe` (RubyMine), or `datagrip64.exe` (DataGrip). In Go, use `github.com/shirou/gopsutil/v3/process` to read the parent name via `p.Name()`. Caveat: Windows aggressively reuses PIDs, so always verify the parent process still exists before trusting its name.

On **Linux**, the parent is a Java process. Inspecting `/proc/<ppid>/cmdline` reveals classpath entries containing product-specific strings like `idea`, `goland`, or `com.intellij.idea.Main`. Walking further up the tree may find the wrapper script or JetBrains Toolbox.

A practical implementation would walk the process tree upward (capped at ~10 levels) and check each ancestor's executable name or command line against a list of known JetBrains identifiers.

## What the charmbracelet ecosystem does

The charmbracelet Go terminal ecosystem (bubbletea, lipgloss, termenv) is the most prominent Go project dealing with JediTerm quirks, but **none of these libraries implement explicit JetBrains detection**. The most notable effort is bubbletea PR #1028 ("fix: wait for GoLand terminal size"), addressing JediTerm's ~400ms delay before reporting real window dimensions. The charmbracelet/crush v0.2.2 release notes mention improved rendering in JetBrains products. `muesli/termenv` checks `TERM`, `COLORTERM`, `TERM_PROGRAM`, and `TERM_PROGRAM_VERSION` but has no `TERMINAL_EMULATOR` check.

## Recommended detection strategy

For the specific use case of disabling CBT when JediTerm is the active terminal emulator, a layered detection approach covers all scenarios:

**Layer 0: `isatty()` guard.** Only perform JediTerm detection when stdout is a TTY. In the standard Run/Debug console (no "Emulate terminal"), `isatty()` returns `false`, and a terminal library shouldn't emit escape sequences at all—making CBT detection unnecessary. This eliminates the most common Run/Debug scenario from requiring any detection logic.

**Layer 1: `TERMINAL_EMULATOR == "JetBrains-JediTerm"` (all platforms).** Handles the Terminal tool window—the most common interactive scenario—with zero false positives.

**Layer 2: macOS bundle identifier (already solved).** Check `__CFBundleIdentifier` or `XPC_SERVICE_NAME` for the `com.jetbrains.` prefix. Covers all macOS execution contexts including "Emulate terminal" Run/Debug.

**Layer 3: `TOOLBOX_VERSION` is set (Windows and Linux).** Covers Run/Debug configs (including "Emulate terminal") when the IDE was launched via Toolbox. Cheap `os.Getenv` call, covers the majority of users since Toolbox is the default install method.

**Layer 4: Parent process inspection (Windows and Linux fallback).** Walk the process tree checking for JetBrains executables. Only needed for standalone IDE installs without Toolbox. Cache after first check.

**Long-term fix: YouTrack feature request.** Request that `TerminalExecutionConsole` set `TERMINAL_EMULATOR=JetBrains-JediTerm` via `PtyProcessBuilder.setEnvironment()`. This is a one-line change on JetBrains' side that would unify detection across all code paths and make Layers 3 and 4 unnecessary for the "Emulate terminal" case.
