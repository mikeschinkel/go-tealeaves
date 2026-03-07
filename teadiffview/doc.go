// Package teadiffview provides TUI diff rendering components for Bubble Tea.
//
// It includes two rendering modes:
//
//   - [SplitDiffModel]: An interactive side-by-side diff viewer with cursor,
//     selection, scrolling, and gutter annotation support. Accepts
//     [diffutils.DiffContent] for Git-agnostic operation.
//
//   - [TUIRenderer]: A condensed diff renderer using lipgloss for styled
//     +/- line output. Suitable for non-interactive display.
//
// # Migration from teadiffr
//
// This package replaces teadiffr, which is now deprecated. Import paths:
//
//	Old: github.com/mikeschinkel/go-tealeaves/teadiffr
//	New: github.com/mikeschinkel/go-tealeaves/teadiffview
//
// All teadiffr types are available here unchanged: [TUIRenderer],
// [DiffRenderer], [FileDiff], [CondensedBlock], [FileStatus], etc.
//
// # Stability
//
// This package is provisional as of v0.1.0. The public API may change in
// minor releases until promoted to stable.
package teadiffview
