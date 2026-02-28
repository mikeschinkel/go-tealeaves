# Plan: go-tealeaves Theming System

## Context

go-tealeaves has 9 modules with styles fragmented across 5 different patterns. Colors are hardcoded as raw string literals (`lipgloss.Color("240")`) in `Default*Style()` functions with no coordination. Applications like gomion cannot theme all components from a single configuration point.

This plan introduces a layered theming system: a `Color` type with named constants, a central color palette, and a full Theme struct on top. Components remain usable standalone — theming is opt-in, never required.

This plan is a hard-gated dependency of:
1. `/Users/mikeschinkel/Projects/go-pkgs/go-tealeaves/UPGRADE_V2_PLAN.md`

## Locked Requirements

1. **Opt-in theming:** Components work without theming via their existing `Default*Style()` functions. Theming is additive, never required.
2. **No magic color strings:** A `Color` type replaces raw `lipgloss.Color("240")` usage. Named constants are provided for all ANSI 256 colors plus curated semantic aliases.
3. **Layered architecture:** Color type → Palette (semantic slots) → Theme (component styles). Use whichever layer fits.
4. **Embeddable palette:** Apps extend `Palette` via Go struct embedding to add domain-specific color slots.
5. **Adaptive backgrounds:** Auto-detect dark/light terminal background. Manual override available. Charm v2 improves detection significantly (`tea.BackgroundColorMsg`, `lipgloss.HasDarkBackground` with explicit I/O).
6. **Consumer-driven:** Apps like gomion can define a custom palette or theme and have all go-tealeaves components respect it.
7. **Lives in teautils module:** Palette and Theme types live in the `teautils` package. Color type and named constants live in `teautils/teacolor` subpackage. No new module — `teacolor` is part of the teautils module. teautils is the shared infrastructure layer for things multiple modules need.
8. **No breaking changes:** Existing `Default*Style()` functions and `With*Style()` methods continue to work. Theme integration is additive.
9. **teagrid v2 supports theming:** The Charm v2 reimagining of teagrid should support `WithTheme()` natively (but not require it).

## Current State: Style Inventory

### Color Palette (Implicit, Distributed)

| Semantic Role | Current Color | Used In |
|---|---|---|
| Primary text | 15 (white) | teadd items, modal buttons |
| Secondary text | 252 (light gray) | teamodal messages, teastatus, list items |
| Muted text | 244-246 (medium gray) | teastatus labels, cancel text |
| Dim text | 240 (dark gray) | borders, separators, scrollbars |
| Accent/title | 46 (bright green) | teamodal titles, teastatus active |
| Selection background | 62 (purple) | teadd, teadep, teamodal focused buttons |
| Selection foreground | 230 (bright yellow) | teadd, teadep, teamodal focused buttons |
| Border | 240 (dark gray) | teadd, teadep borders |
| Border accent | 51 (cyan) | teamodal borders |
| Key display | 86 (cyan) | teautils help visor, teastatus keys |
| Category header | 178 (gold) | teautils help visor |
| Help title | 99 (purple) | teautils help visor |
| Status/warning | 214 (orange) | teamodal status messages |
| Active item | 43 (bright green) | teamodal list active |
| Edit highlight | 11 (bright yellow bg) | teamodal edit mode |

### Style Customization Patterns (5 patterns, inconsistent)

| Pattern | Used By | Mechanism |
|---|---|---|
| Constructor args with defaults | teadd, teadep, teamodal | `ModelArgs` struct fields |
| Public fields | teadd, teadep | Direct assignment |
| Immutable With* methods | teamodal, teadep, teastatus | `WithBorderStyle()` returns copy |
| Styles struct | teastatus | `WithStyles(Styles)` all-or-nothing |
| No customization | teatree (colors), teagrid, teatextsel, teanotify | Hardcoded |

### Gomion Color Audit Summary

Gomion uses ~50+ distinct color values across:
- **UI chrome** (borders, text, selections, focus, buttons) — maps to palette slots
- **File intent states** (commit=green, omit=gray, ignore=red, exclude=orange) — app domain
- **Git staging** (staged=green, unstaged=yellow, both=cyan) — app domain
- **Diff display** (added=green, deleted=red, context tints) — partially palette, partially app
- **Commit group palettes** (8 distinct colors for visual grouping) — app domain

## Design

### Design Principle: UI Chrome vs Domain Colors

1. **UI chrome colors** — borders, text hierarchy, selections, focus indicators, modal styling, status states. Universal to any TUI app. **These belong in the palette.**

2. **Domain-specific colors** — file intents, git staging, diff lines, commit groups. App business logic. **Apps define these themselves**, optionally deriving from palette slots via embedding.

The palette provides the vocabulary. Apps write the sentences.

### Layer 0: Color Type (`teautils/teacolor` subpackage)

Package `teacolor` (import path `github.com/mikeschinkel/go-tealeaves/teautils/teacolor`) provides a `Color` type and named constants. This is a subpackage of the teautils module — no separate `go.mod`.

```go
package teacolor

// Color represents a terminal color value.
type Color = lipgloss.Color

// Full ANSI 256 palette — every slot has a constant.
const (
    Color0   Color = "0"   // Black
    Color1   Color = "1"   // Red (dark)
    Color2   Color = "2"   // Green (dark)
    // ... through Color255
)

// Standard ANSI names (0-15).
const (
    Black       Color = "0"
    Red         Color = "1"
    Green       Color = "2"
    Yellow      Color = "3"
    Blue        Color = "4"
    Magenta     Color = "5"
    Cyan        Color = "6"
    White       Color = "7"
    BrightBlack Color = "8"
    // ... bright variants 9-15
)

// Curated semantic aliases — descriptive names for commonly used colors.
const (
    Coral       Color = "#FF7F50"
    SkyBlue     Color = "#87CEEB"
    Gold        Color = "#FFD700"
    Crimson     Color = "#DC143C"
    DodgerBlue  Color = "#1E90FF"
    Teal        Color = "#008080"
    Salmon      Color = "#FA8072"
    Olive       Color = "#808000"
    Plum        Color = "#DDA0DD"
    SlateGray   Color = "#708090"
    Indigo      Color = "#4B0082"
    DarkGray    Color = "240"
    LightGray   Color = "252"
    // ... ~40-60 curated names total
)
```

Usage:
```go
import "github.com/mikeschinkel/go-tealeaves/teautils/teacolor"

// Before (magic strings):
lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

// After (named constants):
lipgloss.NewStyle().Foreground(teacolor.Color240)
// or with semantic alias:
lipgloss.NewStyle().Foreground(teacolor.DarkGray)
```

Since `Color` is a type alias for `lipgloss.Color`, apps can still create custom colors from hex strings when needed. But named constants are the encouraged path.

### Layer 1: Palette

A `Palette` struct holds semantic color slots for UI chrome. Uses `teacolor.Color` for all fields. Designed for embedding — apps add domain-specific fields.

```go
import "github.com/mikeschinkel/go-tealeaves/teautils/teacolor"

type Palette struct {
    // Text hierarchy (4 levels)
    TextPrimary    teacolor.Color // Main content text
    TextSecondary  teacolor.Color // Labels, descriptions
    TextMuted      teacolor.Color // Metadata, hints
    TextDim        teacolor.Color // Disabled, placeholders

    // Accent colors (for highlights, titles, active elements)
    Accent         teacolor.Color // Primary accent (titles, active items)
    AccentAlt      teacolor.Color // Secondary accent (category headers, links)
    AccentSubtle   teacolor.Color // Subtle accent (key display, help text)

    // Selection and focus
    SelectionBg    teacolor.Color // Selected item background
    SelectionFg    teacolor.Color // Selected item foreground
    FocusBorder    teacolor.Color // Focused pane/component border
    FocusBg        teacolor.Color // Focused row/cell background

    // Borders and chrome
    Border         teacolor.Color // Default border color
    BorderAccent   teacolor.Color // Emphasized borders (modals, active panes)
    Separator      teacolor.Color // Dividers, scrollbar tracks

    // Status indicators
    StatusSuccess  teacolor.Color // Success, confirmations
    StatusInfo     teacolor.Color // Informational
    StatusWarn     teacolor.Color // Warnings, caution
    StatusError    teacolor.Color // Errors, destructive

    // Interactive elements
    ButtonBg       teacolor.Color // Unfocused button background
    ButtonFg       teacolor.Color // Unfocused button foreground
    ButtonFocusBg  teacolor.Color // Focused button background
    ButtonFocusFg  teacolor.Color // Focused button foreground

    // Edit mode
    EditBg         teacolor.Color // Inline edit background
    EditFg         teacolor.Color // Inline edit foreground

    // Diff context (tints for code/content backgrounds)
    TintPositive   teacolor.Color // Added content background tint
    TintNegative   teacolor.Color // Removed content background tint
    TintNeutral    teacolor.Color // Unchanged content background tint

    // Scrollbar
    ScrollTrack    teacolor.Color // Scrollbar track
    ScrollThumb    teacolor.Color // Scrollbar thumb
}
```

Factory functions:
- `DarkPalette() Palette` — Current implicit colors mapped to named constants
- `LightPalette() Palette` — Adjusted for light terminal backgrounds
- `AdaptivePalette() Palette` — Auto-detects background and returns appropriate palette
- `DefaultPalette() Palette` — Returns `AdaptivePalette()` result

### How Apps Extend the Palette

Apps embed `Palette` and add domain-specific color slots. go-tealeaves components accept `Palette` (the embedded portion). App code uses the full struct.

```go
// gomion defines its own extended palette
type GomionPalette struct {
    teautils.Palette // embedded — go-tealeaves components see this

    // File intent states (derived from palette status colors)
    IntentCommit  teacolor.Color
    IntentOmit    teacolor.Color
    IntentIgnore  teacolor.Color
    IntentExclude teacolor.Color

    // Git file status
    FileAdded     teacolor.Color
    FileModified  teacolor.Color
    FileDeleted   teacolor.Color
    FileUntracked teacolor.Color

    // Commit group visual identity (app-specific)
    GroupColors   []teacolor.Color
}

func NewGomionPalette() GomionPalette {
    p := teautils.DarkPalette()
    return GomionPalette{
        Palette:       p,
        IntentCommit:  p.StatusSuccess,
        IntentOmit:    p.TextMuted,
        IntentIgnore:  p.StatusError,
        IntentExclude: p.StatusWarn,
        FileAdded:     p.StatusSuccess,
        FileModified:  p.StatusWarn,
        FileDeleted:   p.StatusError,
        FileUntracked: p.TextMuted,
        GroupColors: []teacolor.Color{
            teacolor.Coral,
            teacolor.SkyBlue,
            teacolor.LightGreen,
            teacolor.Peach,
            teacolor.Plum,
        },
    }
}
```

go-tealeaves components accept the embedded `Palette`. Gomion code uses `gp.IntentCommit`, `gp.FileModified`, etc. No magic strings anywhere.

### Layer 2: Theme

A `Theme` struct holds derived `lipgloss.Style` values for each component's style slots, built from a `Palette`.

```go
type Theme struct {
    Palette Palette

    // Common styles derived from palette
    Border           lipgloss.Style
    BorderAccent     lipgloss.Style
    Title            lipgloss.Style
    Message          lipgloss.Style
    Button           lipgloss.Style
    FocusedButton    lipgloss.Style
    Item             lipgloss.Style
    SelectedItem     lipgloss.Style
    ActiveItem       lipgloss.Style

    // Component-specific style groups
    StatusBar    StatusBarTheme
    HelpVisor    HelpVisorTheme
    Modal        ModalTheme
    Dropdown     DropdownTheme
    List         ListTheme
    Grid         GridTheme
}

type StatusBarTheme struct {
    MenuKeyStyle      lipgloss.Style
    MenuLabelStyle    lipgloss.Style
    IndicatorStyle    lipgloss.Style
    IndicatorSepStyle lipgloss.Style
    BarStyle          lipgloss.Style
}

type HelpVisorTheme struct {
    TitleStyle    lipgloss.Style
    CategoryStyle lipgloss.Style
    KeyStyle      lipgloss.Style
    DescStyle     lipgloss.Style
}

type ModalTheme struct {
    BorderStyle        lipgloss.Style
    TitleStyle         lipgloss.Style
    MessageStyle       lipgloss.Style
    ButtonStyle        lipgloss.Style
    FocusedButtonStyle lipgloss.Style
    CancelKeyStyle     lipgloss.Style
    CancelTextStyle    lipgloss.Style
}

type DropdownTheme struct {
    BorderStyle   lipgloss.Style
    ItemStyle     lipgloss.Style
    SelectedStyle lipgloss.Style
}

type ListTheme struct {
    ItemStyle         lipgloss.Style
    SelectedItemStyle lipgloss.Style
    ActiveItemStyle   lipgloss.Style
    FooterStyle       lipgloss.Style
    StatusStyle       lipgloss.Style
    EditItemStyle     lipgloss.Style
    ScrollbarStyle    lipgloss.Style
    ScrollThumbStyle  lipgloss.Style
}

type GridTheme struct {
    HeaderStyle    lipgloss.Style
    BaseStyle      lipgloss.Style
    HighlightStyle lipgloss.Style
    BorderStyle    lipgloss.Style
}
```

Factory functions:
- `NewTheme(palette Palette) Theme` — Builds all styles from palette
- `DefaultTheme() Theme` — `NewTheme(DefaultPalette())`

### Integration: How Components Use Themes

Each component gains an optional `WithTheme(Theme)` method. When a theme is set, it overrides the component's default styles. When no theme is set, existing `Default*Style()` behavior is preserved.

```go
// Before (still works — no theme required):
m := teamodal.NewOKModal("Hello", &teamodal.ModelArgs{
    ScreenWidth: 80, ScreenHeight: 24,
})

// Opt-in theming:
theme := teautils.DefaultTheme()
m := teamodal.NewOKModal("Hello", &teamodal.ModelArgs{
    ScreenWidth: 80, ScreenHeight: 24,
}).WithTheme(theme)

// Custom palette with named constants (no magic strings):
palette := teautils.DarkPalette()
palette.Accent = teacolor.DodgerBlue
palette.BorderAccent = teacolor.DodgerBlue
theme := teautils.NewTheme(palette)
m := teamodal.NewOKModal("Hello", &teamodal.ModelArgs{
    ScreenWidth: 80, ScreenHeight: 24,
}).WithTheme(theme)
```

Individual `With*Style()` methods continue to work and take precedence over theme — they act as per-instance overrides on top of the theme.

### Integration: How `Default*Style()` Functions Change

The existing `Default*Style()` functions stay as-is for backward compatibility. Internally, they are updated to use named `Color` constants instead of string literals (non-breaking since `Color` is a type alias).

New palette-aware variants are added:

```go
// Existing — updated to use named constants (same behavior)
func DefaultBorderStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(teacolor.DarkGray)
}

// New — palette-aware
func ThemedBorderStyle(p teautils.Palette) lipgloss.Style {
    return lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(p.Border)
}
```

Components with `WithTheme()` call the themed variants internally.

### Integration: Gomion Consumer Pattern

```go
// 1. Create extended palette with domain colors
gp := NewGomionPalette()
// Optionally customize UI chrome:
gp.Accent = teacolor.DodgerBlue
gp.BorderAccent = teacolor.DodgerBlue

// 2. Build theme from the embedded Palette
theme := teautils.NewTheme(gp.Palette)

// 3. Pass theme to all go-tealeaves components
modal := teamodal.NewOKModal("Done", args).WithTheme(theme)
statusBar := teastatus.New().WithTheme(theme)
tree := teatree.NewModel(tree, 20).WithTheme(theme)
dropdown := teadd.NewModel(opts, row, col, args).WithTheme(theme)

// 4. Use domain colors in gomion's own rendering
style := lipgloss.NewStyle().Foreground(gp.IntentCommit) // green
```

All go-tealeaves components share the visual identity via the embedded palette. Gomion's domain colors derive from the same foundation. No magic strings anywhere.

## Scope

In scope:
1. `Color` type (alias for `lipgloss.Color`) with full ANSI 256 constants + ~50 curated semantic aliases.
2. `Palette` struct with ~30 semantic UI chrome color slots, designed for embedding.
3. Built-in dark/light/adaptive palette factories using named constants.
4. `Theme` struct with derived `lipgloss.Style` values for all component style slots, built from a `Palette`.
5. `WithTheme()` methods on all go-tealeaves components that have styles (including teagrid v2).
6. Themed style factory functions alongside existing default factories.
7. Adaptive background detection with manual override.
8. Documentation showing the embedding pattern for app-specific domain colors.

Out of scope:
1. Domain-specific color roles (file intents, git status, commit groups) — consumer defines via embedding.
2. Configuration file format for themes — consumer responsibility.
3. Theme hot-reloading at runtime.
4. Theming for tree branch characters in teatree (structural text, not colors).

---

## PHASE `COLOR` — Color Type and Named Constants (`teautils/teacolor`)

1. Create `teautils/teacolor/` subpackage directory.
2. Define `Color` type alias in `teautils/teacolor/color.go`.
3. Define all 256 ANSI color constants (`Color0` through `Color255`) in `teautils/teacolor/ansi256.go`.
4. Define ~16 standard ANSI name aliases (`Black`, `Red`, `Green`, ..., `BrightWhite`) in `teautils/teacolor/ansi_names.go`.
5. Define ~40-50 curated semantic aliases (`Coral`, `SkyBlue`, `Gold`, `Crimson`, `DodgerBlue`, `DarkGray`, `LightGray`, etc.) in `teautils/teacolor/named.go`.
6. Add tests verifying constants produce valid lipgloss colors.
7. No separate `go.mod` — `teacolor` is part of the `teautils` module.

Acceptance gate: `teacolor` subpackage compiles; all constants defined; importable as `github.com/mikeschinkel/go-tealeaves/teautils/teacolor`; tests pass.

## PHASE `PALETTE` — Palette Type and Built-in Palettes

1. Define `Palette` struct in `teautils/palette.go` using `Color` type for all fields.
2. Implement `DarkPalette()` mapping current implicit colors to named constants.
3. Implement `LightPalette()` with appropriate colors for light terminals.
4. Implement `AdaptivePalette()` using `lipgloss.HasDarkBackground()` (v1) or `tea.BackgroundColorMsg` (v2).
5. Implement `DefaultPalette()` as alias for `AdaptivePalette()`.
6. Add tests: palette construction, embedding pattern works, field access.

Acceptance gate: Palette types compile; dark/light/adaptive factories return valid palettes; embedding verified; tests pass.

## PHASE `THEME` — Theme Type and Style Derivation

1. Define `Theme` struct with common styles and component-specific style groups.
2. Define component-specific theme structs: `StatusBarTheme`, `HelpVisorTheme`, `ModalTheme`, `DropdownTheme`, `ListTheme`, `GridTheme`.
3. Implement `NewTheme(Palette) Theme` — derives all styles from palette semantics.
4. Implement `DefaultTheme()`.
5. Add tests: theme built from palette has expected style properties.

Acceptance gate: Theme construction from palette works; all style slots populated; tests pass.

## PHASE `INTEGRATE` — Component WithTheme() Methods

Per-component integration. Each component gains:
1. A `WithTheme(Theme)` method that stores the theme and applies its styles.
2. Internal logic: if theme is set, use themed styles; otherwise use existing defaults.
3. Existing `Default*Style()` functions updated to use named `Color` constants (non-breaking).

Component order (by complexity):
1. **teastatus** — Already has `Styles` struct; `WithTheme` maps `StatusBarTheme` fields.
2. **teamodal** — Has private style fields with With* methods; `WithTheme` calls them internally.
3. **teadd** — Has public style fields; `WithTheme` sets them.
4. **teadep** — Has public style fields; `WithTheme` sets them.
5. **teatree** — Needs new style support for node colors (focused, expanded indicators).
6. **teatextsel** — Selection highlight style; `WithTheme` sets `SelectionStyle`.
7. **teanotify** — Notice type colors; `WithTheme` adjusts default notice definitions.
8. **teagrid** — v1: partial integration (header/highlight/base). v2 reimagining: full native theming with `GridTheme`.
9. **teautils HelpVisor** — `HelpVisorStyle` already exists; `WithTheme` maps from `HelpVisorTheme`.

Each integration:
- Must NOT break existing `Default*Style()` or `With*Style()` methods.
- Individual `With*Style()` methods take precedence over theme (per-instance override).
- Add `_test.go` coverage for themed vs unthemed behavior.
- Verify visual output is equivalent when theme uses `DarkPalette()` (matches current hardcoded colors).

Acceptance gate: All components accept `WithTheme()`; unthemed behavior is unchanged; themed behavior derives from palette; tests pass.

## PHASE `ADAPTIVE` — Background Detection and Override

1. Implement adaptive detection in `AdaptivePalette()`:
   - v1: Use `lipgloss.HasDarkBackground()`.
   - v2: Use `lipgloss.HasDarkBackground(os.Stdin, os.Stderr)` for standalone usage; in Bubble Tea apps, listen for `tea.BackgroundColorMsg` in Update.
   - Fall back to dark palette if detection fails or is inconclusive.
2. Add override mechanism: `func SetDarkBackground(isDark bool)` or similar to let apps that provide user configuration override detection.
3. Document detection behavior, Bubble Tea integration, and override pattern.
4. Test adaptive behavior with mocked detection.

Acceptance gate: Auto-detection works; manual override works; fallback is safe; Bubble Tea integration path documented.

## PHASE `DOCS` — Documentation

1. Add theming section to teautils README.
2. Document Color constants — full list with categories.
3. Document palette color slots and their semantic meaning.
4. Document the embedding pattern for app-specific domain colors (with gomion-style example).
5. Document consumer pattern (palette → theme → WithTheme on components).
6. Document adaptive behavior, Bubble Tea integration, and manual override.
7. Add before/after examples showing themed vs unthemed usage.

Acceptance gate: A consumer can theme all go-tealeaves components and extend with domain colors using docs alone.

## PHASE `VALIDATE` — Cross-Component Verification

1. `go build ./...` and `go test ./...` across all modules.
2. Verify existing examples still compile and render correctly (unthemed).
3. Build a themed example (or modify an existing one) demonstrating cross-component theming.
4. Verify gomion can integrate theming with embedding pattern.

Acceptance gate: All builds and tests pass; themed and unthemed paths work; at least one example demonstrates theming.

---

## Risks

1. **Import cycle risk:** Components import teautils for `Theme`; teautils must not import components.
   Mitigation: Theme types are defined in teautils; components import teautils. One-way dependency — same as today with `KeyRegistry`.

2. **Palette slot coverage:** New UI chrome needs may not map to existing palette slots.
   Mitigation: Adding a slot is non-breaking (zero value is valid). App domain colors use embedding, never require new palette slots.

3. **Adaptive detection on v1:** `lipgloss.HasDarkBackground()` v1 has implicit I/O behavior that can be unreliable.
   Mitigation: Dark palette as fallback; manual override always available. v2 detection is explicit and reliable.

4. **Charm v2 color type change:** `lipgloss.Color` type may change in v2 (return type of `lipgloss.Color()` becomes `color.Color`).
   Mitigation: `Color` type alias isolates the change. Update the alias during BUBLIP phase of UPGRADE_V2_PLAN.md.

5. **Color constant API surface:** ~300 constants is a large API surface.
   Mitigation: Organize in `color.go` with clear sections. Constants are simple, self-documenting, and cannot be wrong.

## Completion Checklist

1. `COLOR` gate satisfied.
2. `PALETTE` gate satisfied.
3. `THEME` gate satisfied.
4. `INTEGRATE` gate satisfied.
5. `ADAPTIVE` gate satisfied.
6. `DOCS` gate satisfied.
7. `VALIDATE` gate satisfied.
8. This plan no longer blocks RELEASE phase of UPGRADE_V2_PLAN.md.
