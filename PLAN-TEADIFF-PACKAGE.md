# Plan: Create teadiff Package

## Goal

Create a new `teadiff` package in go-tealeaves that provides TUI rendering for multi-file condensed diffs. This package knows about Bubble Tea rendering (lipgloss, ANSI) but has NO dependency on gitutils or gomion.

## Background

Currently `gomtui/context_info_pane_model.go` has monolithic rendering logic for condensed diffs tightly coupled to gomion's context pane. The rendering (file headers, block headers, context lines, changed lines, truncation, status tags, background tints) is general-purpose and reusable.

After this extraction, any Bubble Tea application can render multi-file diffs using `teadiff`, without importing gomion.

## Input Types

`teadiff` must NOT import `gitutils`. It receives pre-built condensed diff data. Define its own input types that mirror the gitutils types:

```go
// FileStatus indicates whether a file is new, deleted, or modified
type FileStatus int

const (
    FileModified FileStatus = iota
    FileNew
    FileDeleted
)

// CondensedBlock represents a condensed view of a change block
type CondensedBlock struct {
    Type          string   // "added" or "deleted"
    LineCount     int
    ContextBefore []string
    ChangedLines  []string
    ContextAfter  []string
    IsTruncated   bool
}

// FileDiff holds condensed diffs for a single file
type FileDiff struct {
    Path   string
    Status FileStatus
    Blocks []CondensedBlock
}
```

These are intentionally simple (string-based path, no `dt.RelFilepath`) to avoid coupling to go-dt internals. The caller (gomtui) converts from `gitutils.FileCondensedDiff` to `teadiff.FileDiff` at the integration boundary.

Alternatively, if go-tealeaves already depends on go-dt, use `dt.RelFilepath` directly and avoid the conversion overhead.

## DiffRenderer Interface

```go
// DiffRenderer formats diff content for a specific output medium
type DiffRenderer interface {
    RenderFileHeader(path string, status FileStatus, width int) string
    RenderBlockHeader(blockType string, lineCount int) string
    RenderContextLine(line string, status FileStatus, width int) string
    RenderAddedLine(line string, status FileStatus, width int) string
    RenderDeletedLine(line string, status FileStatus, width int) string
    RenderTruncation(status FileStatus) string
    RenderSeparator() string
}
```

The `status` parameter lets renderers apply file-level visual treatment (green tint for new, red tint for deleted) to individual lines.

## Built-in TUI Renderer

```go
// TUIRenderer renders diffs using lipgloss for terminal output
type TUIRenderer struct {
    // Style configuration (optional overrides)
    FileHeaderColor   lipgloss.Color // default: "244"
    BlockHeaderColor  lipgloss.Color // default: "135"
    ContextColor      lipgloss.Color // default: "240"
    AddedColor        lipgloss.Color // default: "34"
    DeletedColor      lipgloss.Color // default: "160"
    NewStatusColor    lipgloss.Color // default: "46"
    DeletedStatusColor lipgloss.Color // default: "168"
    NewBgColor        lipgloss.Color // default: "22"
    DeletedBgColor    lipgloss.Color // default: "52"
}
```

Default colors match the current gomtui implementation. All are overridable.

## RenderFileDiffs Function

Top-level rendering function that produces `[]string` (one per output line):

```go
// RenderFileDiffs renders multiple file diffs into styled lines
func RenderFileDiffs(files []FileDiff, renderer DiffRenderer, width int) []string
```

This replaces the file-rendering loop currently in `viewDirectory()`. It:
1. Iterates files, calling `RenderFileHeader` for each
2. Iterates blocks within each file, calling block/line renderers
3. Returns the complete set of rendered lines

The caller (gomtui) handles scrolling, padding, and viewport management.

## Package Structure

```
go-tealeaves/teadiff/
    types.go          — FileDiff, CondensedBlock, FileStatus, DiffRenderer interface
    tui_renderer.go   — TUIRenderer implementation
    render.go         — RenderFileDiffs top-level function
```

## Dependencies

```
teadiff → lipgloss, x/ansi
          NO dependency on gitutils, gomion, or go-dt (unless go-tealeaves already uses go-dt)
```

Follows the existing tea* pattern (teadd, teamodal, teastatus, etc.).

## Verification

1. `go build ./teadiff/...`
2. `go test ./teadiff/...` — unit tests rendering known inputs
3. Verify go-tealeaves module compiles: `go build ./...`

## Prerequisites

- The gitutils condensed diff extraction (Layer 1) should be complete first so the type design is stable
- However, since teadiff defines its own input types, it CAN be developed in parallel

## Future Extensions

- **Markdown renderer**: Implements `DiffRenderer` producing markdown output
- **HTML renderer**: Implements `DiffRenderer` producing HTML with CSS classes
- **teadiff.Model**: Optional Bubble Tea model wrapping rendered content with scrolling (if reuse demand warrants it)
