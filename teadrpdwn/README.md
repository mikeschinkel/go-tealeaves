## TeaDD — A Dropdown Component for Bubble Tea Apps

A TUI dropdown/popup selection component for **FULL SCREEN** [Bubble Tea](https://github.com/charmbracelet/bubbletea), with intelligent automatic positioning.

**Note:** This component uses absolute positioning and overlay rendering, designed for full-screen terminal UIs (with `tea.WithAltScreen()`). It is not suitable for inline CLI prompts.

## Features

- **Full-Screen TUI Design**: Built for terminal UIs with absolute positioning and overlay rendering
- **Intelligent Positioning**: Automatically displays below or above the field based on available space
- **Scrolling Support**: Handles lists of any size with automatic scrolling
- **Truncation**: Long items are automatically truncated with ellipsis
- **Customizable Styling**: Full control over borders, items, and selection appearance via lipgloss
- **Modal Behavior**: Clean message consumption pattern - doesn't "infect" parent Update()
- **String Compositing**: Overlay approach for seamless integration with parent views

## Installation

```bash
# Currently part of gommod module
# Future standalone: go get github.com/mikeschinkel/go-tealeaves/teadd
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/mikeschinkel/go-tealeaves/teadd"
)

type model struct {
    dropdown teadd.DropdownModel
}

func (m model) Init() tea.Cmd {
    return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var dropdown tea.Model

    // Let dropdown handle message first
    dropdown, cmd = m.dropdown.Update(msg)
    if cmd != nil {
        m.dropdown = dropdown.(teadd.DropdownModel)
        return m, cmd
    }

    // Dropdown didn't handle - parent processes
    switch msg := msg.(type) {
    case teadd.OptionSelectedMsg:
        fmt.Printf("Selected: %s\n", msg.Text)
        // msg.Value contains the underlying value (same as Text when using ToOptions)
        return m, tea.Quit

    case teadd.DropdownCancelledMsg:
        return m, tea.Quit

    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case " ":
            if m.dropdown.IsOpen {
                m.dropdown, cmd = m.dropdown.Close()
            } else {
                m.dropdown, cmd = m.dropdown.Open()
            }
            return m, cmd
        }

    case tea.WindowSizeMsg:
        m.dropdown = m.dropdown.WithScreenSize(msg.Width, msg.Height)
        return m, nil
    }

    return m, nil
}

func (m model) View() string {
    baseView := "Press Space to open dropdown, q to quit\n\n"

    // Add field at position
    field := "▼ Select Option"
    lines := make([]string, 20)
    lines[0] = baseView
    lines[5] = "    " + field  // Field at row 5, col 4
    baseView = strings.Join(lines, "\n")

    // Composite dropdown if open
    if m.dropdown.IsOpen {
        dropdownView := m.dropdown.View()
        return teadd.OverlayDropdown(baseView, dropdownView, m.dropdown.Row, m.dropdown.Col)
    }

    return baseView
}

func main() {
    items := teadd.ToOptions([]string{"Apple", "Banana", "Cherry", "Date", "Elderberry"})

    dropdown := teadd.NewModel(items, 5, 4, &teadd.ModelArgs{
        ScreenWidth:  80,
        ScreenHeight: 24,
        TopMargin:    1,  // Reserve top row for menu bar
        BottomMargin: 1,  // Reserve bottom row for status bar
    })

    p := tea.NewProgram(model{
        dropdown: dropdown,
    })

    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

## API Reference

### Types

#### Option

```go
type Option struct {
    Text  string      // Display text shown in dropdown
    Value interface{} // Underlying value (e.g., DB primary key, ID, any type)
}
```

Represents a dropdown item with separate display text and value.

**Helper function:**
```go
func ToOptions(strings []string) []Option
```

Converts a string slice to Options where Text and Value are the same. Useful for simple dropdowns where the display text is the value.

**Example:**
```go
// Simple case - strings
items := teadd.ToOptions([]string{"Apple", "Banana", "Cherry"})

// Advanced case - with values
items := []teadd.Option{
    {Text: "Apple", Value: 1},
    {Text: "Banana (Fresh)", Value: 2},
    {Text: "Record #123", Value: dbRecord},
}
```

#### DropdownModel

```go
type DropdownModel struct {
    // Position in parent view
    Row int
    Col int

    // Field position (reference point for dropdown positioning)
    FieldRow int
    FieldCol int

    // Options and selection
    Options        []Option
    Selected     int
    ScrollOffset int // First visible item index (for scrolling)

    // Display state
    IsOpen       bool
    DisplayAbove bool // True if dropdown is displayed above field, false if below
    ScreenWidth  int
    ScreenHeight int

    // Margins - exclude screen areas from dropdown positioning
    TopMargin    int // Don't position dropdown above this row (e.g., 1 to avoid menu bar)
    BottomMargin int // Don't position dropdown below screenHeight - bottomMargin (e.g., 1 to avoid status bar)

    // Styling (public for customization)
    BorderStyle   lipgloss.Style
    ItemStyle     lipgloss.Style
    SelectedStyle lipgloss.Style
}
```

#### Messages

```go
// Sent when user confirms selection with Enter
type OptionSelectedMsg struct {
    Index int
    Text  string      // Display text
    Value interface{} // Underlying value
}

// Sent when user cancels with Esc
type DropdownCancelledMsg struct{}
```

#### Position Constants

```go
type Position int

const (
    TopLeft Position = iota
    TopMiddle
    TopRight
    MiddleLeft
    Middle
    MiddleRight
    BottomLeft
    BottomMiddle
    BottomRight
)
```

### Constructor

```go
func NewModel(items []Option, fieldRow, fieldCol int, args *ModelArgs) DropdownModel

type ModelArgs struct {
    ScreenWidth  int
    ScreenHeight int
    TopMargin    int            // Don't position dropdown above this row
    BottomMargin int            // Don't position dropdown below screenHeight - bottomMargin
    BorderStyle   lipgloss.Style // Optional - uses default if not provided
    ItemStyle     lipgloss.Style // Optional - uses default if not provided
    SelectedStyle lipgloss.Style // Optional - uses default if not provided
}
```

Creates a new dropdown for a field at the specified position.

**Parameters**:
- `items`: List of selectable items (use `ToOptions()` to convert from `[]string`)
- `fieldRow`: Row position of the field in parent view (dropdown calculates its own position)
- `fieldCol`: Column position of the field in parent view
- `args`: Configuration arguments (screen size, margins, optional styling). Pass `nil` to use all defaults (screen size will be set from `tea.WindowSizeMsg`)

### Methods

```go
func (m DropdownModel) Init() tea.Cmd
func (m DropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m DropdownModel) View() string
```

Standard Bubble Tea model interface.

**Important**: `Update()` returns a non-nil `tea.Cmd` when it handles a message, `nil` when it doesn't. Parent should check `if cmd != nil` to determine if message was consumed.

```go
func (m DropdownModel) Open() (DropdownModel, tea.Cmd)
```

Opens the dropdown and returns the updated model.

```go
func (m DropdownModel) Close() (DropdownModel, tea.Cmd)
```

Closes the dropdown and returns the updated model.

### Fluent Configuration Methods

These methods return a new `DropdownModel` with the specified field updated, allowing method chaining:

```go
func (m DropdownModel) WithPosition(fieldRow, fieldCol int) DropdownModel
```

Updates field position. Dropdown position is recalculated when opened.

```go
func (m DropdownModel) WithOptions(items []Option) DropdownModel
```

Updates items list. Adjusts `Selected` if necessary. Use `ToOptions()` to convert from `[]string`.

```go
func (m DropdownModel) WithScreenSize(width, height int) DropdownModel
```

Updates screen dimensions. Call when handling `tea.WindowSizeMsg`.

```go
func (m DropdownModel) WithTopMargin(margin int) DropdownModel
```

Sets top margin. Dropdown won't position above this row (useful for menu bars).

```go
func (m DropdownModel) WithBottomMargin(margin int) DropdownModel
```

Sets bottom margin. Dropdown won't extend below `screenHeight - margin` (useful for status bars).

### Functions

```go
func OverlayDropdown(background, foreground string, row, col int) string
```

Overlays dropdown view onto base view at specified position. Use this in your parent's `View()` method when dropdown is open.

**Parameters**:
- `background`: The base view (fully rendered string with ANSI codes)
- `foreground`: The overlay/dropdown view (fully rendered string with ANSI codes)
- `row`: Line number in background where foreground row 0 should appear (0-indexed)
- `col`: Display column in background where foreground col 0 should appear (0-indexed)

```go
func EnsureTermGetSize(fd uintptr) (w int, h int, ok bool)
```

Robust terminal size detection with IDE fallbacks. This utility helps get accurate terminal dimensions even in environments where standard `term.GetSize()` may initially return (0,0), such as GoLand's terminal during debugging.

**How it works**:
1. Retries `term.GetSize()` briefly (250ms) with small delays for terminals that initialize slowly
2. Falls back to `/dev/tty` (controlling terminal) if available
3. Falls back to `COLUMNS` and `LINES` environment variables
4. Returns `ok=false` if no valid size can be determined

**Usage**:
```go
func main() {
    // Call before starting Bubble Tea program to ensure term.GetSize() is ready
    teadd.EnsureTermGetSize(os.Stdout.Fd())

    p := tea.NewProgram(initialModel(), tea.WithAltScreen())
    // ...
}
```

**When to use**: Call this before starting your Bubble Tea program if you're experiencing issues with terminal size detection in IDEs or specific terminal emulators. Not required for most environments, but harmless if called unnecessarily.

### Styling

Customize dropdown appearance by providing styles in `ModelArgs`:

**Example**:
```go
dropdown := teadd.NewModel(items, 5, 10, &teadd.ModelArgs{
    ScreenWidth:  80,
    ScreenHeight: 24,
    BorderStyle: lipgloss.NewStyle().
        BorderStyle(lipgloss.DoubleBorder()).
        BorderForeground(lipgloss.Color("62")),
    SelectedStyle: lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("230")).
        Background(lipgloss.Color("63")),
})
```

### Default Styles

```go
func DefaultBorderStyle() lipgloss.Style
func DefaultItemStyle() lipgloss.Style
func DefaultSelectedStyle() lipgloss.Style
```

Returns the default styles used when no options are provided.

## Positioning

The dropdown uses intelligent automatic positioning:

1. **Vertical Placement**: Compares available space below vs above the field, displays dropdown in whichever location has more space
2. **Horizontal**: Left border starts 2 columns left of field's first character
3. **Overflow Handling**:
   - **Vertical**: Uses all available space (respects margins), enable scrolling if items don't fit
   - **Horizontal**: Shifts left if extends beyond screen width, truncates to screen width if necessary
   - **Truncation**: Options exceeding width are truncated with ellipsis (…)

### Positioning Rules

- **Vertical placement**: Compare available space below vs above field, prefer whichever has more space
- **Available space calculation**: Respects `TopMargin` and `BottomMargin` settings
- **Visible items**: Uses all available space up to total item count (no artificial limit)
- **Scrolling**: Automatically enabled when not all items fit in available space
- **Scroll behavior**: Highlight moves to top/bottom edge of visible area before scrolling starts (intuitive navigation)
- **Scroll indicators**: ▲ appears at top when items above, ▼ appears at bottom when items below

## Styling

All visual aspects are customizable via lipgloss styles:

```go
import "github.com/charmbracelet/lipgloss"

dropdown := teadd.NewModel(items, row, col, &teadd.ModelArgs{
    ScreenWidth:  80,
    ScreenHeight: 24,
    BorderStyle: lipgloss.NewStyle().
        BorderStyle(lipgloss.ThickBorder()).
        BorderForeground(lipgloss.Color("205")).
        Padding(0, 1),
    ItemStyle: lipgloss.NewStyle().
        Foreground(lipgloss.Color("252")),
    SelectedStyle: lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("230")).
        Background(lipgloss.Color("99")),
})
```

## Examples

### Basic Example

See `examples/basic/` for a simple implementation:

```bash
cd examples/basic
go run main.go
```

Shows the essential pattern: create dropdown, handle messages, overlay on view.

### Demo App

See `demo/` for an interactive app that exercises all positioning scenarios:

- 9 field positions (horizontal × vertical combinations)
- 3 item counts (3, 7, 25 items with scrolling)
- 3 item widths (short, medium, long with truncation)
- 2 start states (open or closed)

```bash
cd demo
go run main.go
```

The demo uses an interactive configuration UI with progressive disclosure, allowing you to test any specific scenario.

## Modal Behavior Pattern

The dropdown implements clean message consumption without "infecting" the parent:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var dropdown tea.Model

    // Let dropdown handle first
    dropdown, cmd = m.dropdown.Update(msg)
    if cmd != nil {
        m.dropdown = dropdown.(teadd.DropdownModel)
        return m, cmd  // Dropdown consumed message - we're done
    }

    // Dropdown didn't handle - parent processes
    switch msg := msg.(type) {
    case teadd.OptionSelectedMsg:
        // Handle selection
        m.selectedValue = msg.Value
        m.dropdown, cmd = m.dropdown.Close()
        return m, cmd

    case tea.KeyMsg:
        switch msg.String() {
        case " ":
            if !m.dropdown.IsOpen {
                m.dropdown, cmd = m.dropdown.Open()
                return m, cmd
            }
        }
    }

    return m, nil
}
```

**Key insight**: The dropdown returns `nil` cmd when it doesn't handle the message, non-nil when it does. This prevents "modal infection" - parent doesn't need to check dropdown state for every key. The state check for opening (`if !m.dropdown.IsOpen`) is a UI decision (don't open when already open), not modal behavior enforcement. All methods that modify the model return a new model value (functional style).

## Testing

### Test Scenarios

The example app supports testing **162 configuration permutations** (9 positions × 3 item counts × 3 item widths × 2 start states):

**Field Positions (9)**:
- Top: Left, Middle, Right
- Middle: Left, Center, Right
- Bottom: Left, Middle, Right

**Option Counts (3)**:
- 3 items (small list, fits in minimal space)
- 7 items (moderate list, tests typical dropdown size)
- 25 items (large list, requires scrolling)

**Option Widths (3)**:
- Short (10 chars)
- Medium (calculated at 70% of screen width)
- Long (calculated at 125% of screen width - tests truncation)

**Start States (2)**:
- Open (dropdown visible on load)
- Closed (dropdown hidden on load)

### Validation Checklist

**Positioning**:
- [ ] Displays below field when space available
- [ ] Displays above field when insufficient space below
- [ ] Never renders off-screen (all 9 positions tested)
- [ ] Correct horizontal alignment (2 columns left of field)
- [ ] Shifts left if would overflow right edge
- [ ] Shifts right if would start at negative column

**Rendering**:
- [ ] Options truncate with ellipsis (…) when too wide
- [ ] Scrolling works when items exceed available space (25-item test)
- [ ] Selected item highlighted correctly
- [ ] Border renders with correct dimensions
- [ ] Overlay composition with parent view works (no overlap issues)
- [ ] Scroll indicators (▲/▼) appear when scrolling active

**Functionality**:
- [ ] Dropdown opens/closes correctly
- [ ] Options navigate with up/down/k/j
- [ ] Enter sends OptionSelectedMsg with correct index/value
- [ ] Esc sends DropdownCancelledMsg
- [ ] Selection clamped at boundaries (doesn't wrap)
- [ ] Screen resize updates positioning correctly

**Modal Behavior**:
- [ ] Modal keys (q) don't reach parent when dropdown open
- [ ] Parent keys work when dropdown closed
- [ ] Cmd return signals message consumption correctly

**Edge Cases**:
- [ ] Very small terminal (< 20×20)
- [ ] Very large terminal (> 200×80)
- [ ] Resize during dropdown open
- [ ] Empty items list (should error gracefully)
- [ ] Single item list
- [ ] Options longer than screen width (125% test)
- [ ] Field at screen edges (all 9 positions)

## Planned Features

### Container-Aware Positioning

**Feature**: `AllowOverflow` option to control whether dropdown can extend beyond container boundaries.

**Use Case**: When dropdown is within a bordered/constrained container (like a screenArea), should it:
- **Stay inside** (AllowOverflow=false): Respect container bounds, constrain dropdown width/position
- **Pop out** (AllowOverflow=true): Overlay container borders for maximum space

**Design Considerations**:
- **Coordinate system**: Row/Col are currently absolute screen coordinates, but container bounds need to be relative
- **Options**:
  - A) Switch to relative coordinates within container (simpler, requires app changes)
  - B) Add ContainerLeft/Top/Width/Height fields (complex, backward compatible)
- **Trade-offs**: Needs real-world usage patterns to determine best approach

**Status**: Deferred until usage patterns are clearer

### Field Text Management

**Feature**: Smart field text truncation in demo/example apps

**Current Behavior**: Long field values (like file paths) wrap or extend off-screen

**Desired Behavior**:
- Calculate available space for field text based on position
- Truncate with ellipsis if needed
- Show full value in dropdown

**Status**: Low priority - this is app-level responsibility, not component responsibility

## License

MIT License - See LICENSE.txt
