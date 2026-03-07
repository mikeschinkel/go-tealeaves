# Planning Brief: Syntax Highlighting as a Foundation Feature

## Objective
Add syntax highlighting as a Foundation feature in go-tealeaves, making it available to any component that renders code (SplitDiffModel, future code viewers, etc.).

## Context
- Currently lives in gomion/gommod/gomtui/syntax_highlighter.go (~102 lines)
- Uses github.com/alecthomas/chroma/v2 for highlighting
- Two key functions: `DetectChromaLexerName(path)` and `HighlightCode(text, language)`
- User explicitly said this should be a Foundation feature (alongside Key Registry, Help Visor, Theming, Positioning)
- Foundation = lives in teautils, is NOT a tea.Model, provides infrastructure for other components
- Must be part of the website updates (site page needed)

## Source Reference
`/Users/mikeschinkel/Projects/gomion/gommod/gomtui/syntax_highlighter.go`:
- `DetectChromaLexerName(relPath dt.RelFilepath) string` — detects language from file extension
- `HighlightCode(code string, language string) string` — applies Chroma terminal highlighting

## Design Questions

### 1. Where Does It Live?
Options:
- **teautils** (alongside other Foundations) — keeps all infrastructure in one module
- **New module teahighlight** — separates the Chroma dependency (heavy) from teautils
- Recommendation: Likely teautils for consistency, but Chroma adds a non-trivial dependency

### 2. API Surface
What should the public API look like?
- Simple function API: `teautils.HighlightCode(code, language string) string`
- Language detection: `teautils.DetectLanguage(filepath string) string`
- Theme support: Should it integrate with teautils Theming (Palette/Theme)?
- Should there be a `Highlighter` struct that caches lexers?

### 3. Integration with SplitDiffModel
The SplitDiffModel extraction (separate effort) needs highlighting. Options:
- SplitDiffModel accepts `HighlightFunc func(text, lang string) string` in args
- SplitDiffModel imports teautils and calls highlighting directly
- Recommendation: Function argument for loose coupling

### 4. Terminal Format
Chroma has multiple terminal formatters:
- `terminal` (basic 8-color)
- `terminal256` (256-color)
- `terminal16m` (true color)
- Which to default to? Auto-detect?

## Site Page Needed
- `site/src/content/docs/components/syntax-highlighting.mdx`
- Grouped under "Foundations" in sidebar
- Template: info line (Module/Type/Install), Quick Example, Features, API Reference

## Dependencies
- github.com/alecthomas/chroma/v2
- github.com/mikeschinkel/go-dt (for RelFilepath if keeping that API)

## Non-Goals
- Full IDE-like highlighting with LSP integration
- Custom language grammar support
- Line-number rendering (that's the diff viewer's job)

## Go House Rules
- Follow ClearPath production style, doterr error handling
- No compound if-init statements
- Package-level var for compiled regexps
