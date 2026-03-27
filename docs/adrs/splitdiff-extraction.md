# Planning Brief: SplitDiffPaneModel Extraction

## Objective
Extract the SplitDiffPaneModel from gomion/gommod/gomtui into go-tealeaves/teadiffr as an interactive side-by-side diff viewer component.

## Context
- The existing teadiffr module provides condensed/summary diff rendering (TUIRenderer). The SplitDiffPaneModel is a full interactive line-by-line diff viewer that would coexist alongside it.
- User considers this a prerequisite before further site updates.
- Syntax highlighting extraction is a separate parallel effort (see syntax-highlight-foundation.md).

## Source Files (in gomion)
All in `/Users/mikeschinkel/Projects/gomion/gommod/gomtui/`:

| File | Lines | Purpose |
|------|-------|---------|
| `split_diff_pane_model.go` | 1153 | Main tea.Model — 50+ exported methods, cursor/selection/scroll/gutter state |
| `diff_builder.go` | 265 | Converts diff data into aligned SplitPaneRows |
| `dff_line.go` | 84 | PaneLine interface + TextLine, BlockMarker, PlaceholderLine types |
| `diff_pane_colors.go` | 37 | ANSI color constants |
| `syntax_highlighter.go` | 102 | Chroma-based highlighting (DO NOT extract here — separate effort) |
| `diff_pane_model_test.go` | ~50 | Tests |

## Target
`/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/teadiffr/`

## Key Decisions Needed

### 1. Abstracting gitutils.FileChanges
The main blocker. The model currently depends on `gitutils.FileChanges` (gomion-specific). Options:
- **Option A**: Create a `DiffPaneData` interface in teadiffr that gomion's FileChanges implements
- **Option B**: Define concrete types in teadiffr (DiffFile, ChangeBlock, etc.) and have gomion convert

### 2. Syntax Highlighting
Currently calls `DetectChromaLexerName()` and `HighlightCode()` from gomtui's syntax_highlighter.go.
- Do NOT extract highlighting as part of this effort (separate Foundation effort)
- Instead: accept an optional `HighlightFunc func(text, language string) string` in constructor args
- Or: accept pre-highlighted text lines

### 3. Naming
- Rename to `SplitDiffModel` (drop "Pane" — conciseness)
- Constructor: `NewSplitDiffModel(args *SplitDiffModelArgs) SplitDiffModel`

### 4. Dependencies After Extraction
gomion/gommod/gomtui will import from teadiffr instead. Changes needed:
- `batch_assignment_model.go` — creates and uses SplitDiffPaneModel
- Possibly `file_intent_model.go` and `commit_target_model.go`
- Import path changes only; public API stays the same

## Charmbracelet Dependencies (already in go-tealeaves)
- bubbletea/v2
- lipgloss/v2
- charmbracelet/x/ansi

## Non-Goals
- Syntax highlighting extraction (separate brief)
- Changing the condensed diff renderer (TUIRenderer)
- Site page updates (downstream of this)

## Go House Rules
- Follow ClearPath production style, doterr error handling
- No compound if-init statements
- Use go-dt types where applicable (RelFilepath already used)
