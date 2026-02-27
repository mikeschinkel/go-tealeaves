## TeaModal — A Modal Dialog Component for Bubble Tea Apps

A TUI modal dialog component for **FULL SCREEN** [Bubble Tea](https://github.com/charmbracelet/bubbletea), with centered positioning and clean message consumption pattern.

**Note:** This component uses absolute positioning and overlay rendering, designed for full-screen terminal UIs (with `tea.WithAltScreen()`). It is not suitable for inline CLI prompts.

## Features

- **Full-Screen TUI Design**: Built for terminal UIs with absolute positioning and overlay rendering
- **Centered Positioning**: Automatically centers modal dialog in the screen
- **Multiple Modal Types**: Yes/No confirmations and OK alerts
- **Button Focus**: Tab to switch between buttons, Enter to confirm
- **Mouse Support**: Click buttons to select, hover for visual feedback
- **Customizable Styling**: Full control over borders, title, message, and button appearance via lipgloss
- **Customizable Keys**: Rebind keyboard shortcuts via Keys field
- **Modal Behavior**: Clean message consumption pattern - doesn't "infect" parent Update()
- **String Compositing**: Overlay approach for seamless integration with parent views
- **Fluent API**: Wither methods for runtime property updates

## Installation

```bash
# Currently part of gommod module
# Future standalone: go get github.com/mikeschinkel/go-tealeaves/teamodal
```

## Quick Start

### Yes/No Confirmation Dialog

```go
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

type model struct {
	confirmDialog teamodal.ModalModel
	confirmed     bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var modal tea.Model

	// Let modal handle message first
	modal, cmd = m.confirmDialog.Update(msg)
	if cmd != nil {
		m.confirmDialog = modal.(teamodal.ModalModel)
		return m, cmd
	}

	// Modal didn't handle - parent processes
	switch msg := msg.(type) {
	case teamodal.AnsweredYesMsg:
		m.confirmed = true
		return m, tea.Quit

	case teamodal.AnsweredNoMsg, teamodal.ClosedMsg:
		return m, tea.Quit

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case " ":
			m.confirmDialog, cmd = m.confirmDialog.Open()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.confirmDialog = m.confirmDialog.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() string {
	baseView := "Press Space to open confirmation, q to quit\n\n"

	if m.confirmed {
		baseView += "User confirmed!"
	}

	// Composite modal if open (automatic positioning)
	return m.confirmDialog.OverlayModal(baseView)
}

func main() {
	modal := teamodal.NewYesNoModal("Do you want to proceed?", &teamodal.ModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
		Title:        "Confirmation",
		DefaultYes:   true,
	})

	p := tea.NewProgram(model{
		confirmDialog: modal,
	}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
```

### OK Alert Dialog

```go
modal := teamodal.NewOKModal("Operation completed successfully!", &teamodal.ModelArgs{
	ScreenWidth:  80,
	ScreenHeight: 24,
	Title:        "Success",
})
```

## API Reference

### Types

#### ModalType

```go
type ModalType int

const (
	ModalTypeOK     ModalType = iota // Single OK button (alert)
	ModalTypeYesNo                   // Yes and No buttons (confirmation)
)
```

#### ModalKeyMap

```go
type ModalKeyMap struct {
	Confirm      key.Binding // enter - confirm selection
	Cancel       key.Binding // esc - cancel/close
	NextButton   key.Binding // tab - move to next button (YesNo only)
	PrevButton   key.Binding // shift+tab - move to previous button (YesNo only)
	SelectLeft   key.Binding // left - select left button (Yes) (YesNo only)
	SelectRight  key.Binding // right - select right button (No) (YesNo only)
}
```

Keyboard bindings for modal dialogs. Access via `modal.Keys` field. Use `DefaultModalKeyMap()` to get default bindings.

#### ModalModel

```go
type ModalModel struct {
	Keys ModalKeyMap // Keyboard bindings (customizable)

	// All fields below are private - access via getter methods
	// Content: Title(), Message(), Type()
	// Button labels: YesLabel(), NoLabel(), OKLabel()
	// State: IsOpen(), FocusButton(), ScreenWidth(), ScreenHeight()
	// Styling: BorderStyle(), TitleStyle(), MessageStyle(), ButtonStyle(), FocusedButtonStyle()
}
```

**Note**: All model fields are private. Use getter methods like `m.Title()`, `m.IsOpen()`, `m.ScreenWidth()`, etc. to access values. Use wither methods like `m.WithTitle("New Title")` to create modified copies.

#### Messages

```go
// Sent when user confirms with Enter on Yes button
type AnsweredYesMsg struct{}

// Sent when user selects No button and presses Enter
type AnsweredNoMsg struct{}

// Sent when user closes alert (OK) or cancels with Esc
type ClosedMsg struct{}
```

### Constructors

```go
func NewOKModal(message string, args *ModelArgs) ModalModel

func NewYesNoModal(message string, args *ModelArgs) ModalModel

type ModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string
	DefaultYes   bool   // For YesNo: default focus to Yes (true) or No (false)
	YesLabel     string // Custom Yes button label
	NoLabel      string // Custom No button label
	OKLabel      string // Custom OK button label

	// Alignment (optional - defaults to lipgloss.Center)
	TextAlign    lipgloss.Position // Horizontal alignment for title and message
	TitleAlign   lipgloss.Position // Horizontal alignment for title (overrides TextAlign)
	MessageAlign lipgloss.Position // Horizontal alignment for message (overrides TextAlign)
	ButtonAlign  lipgloss.Position // Horizontal alignment for buttons (defaults to Center)

	// Optional styling
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
}
```

### Methods

```go
func (m ModalModel) Init() tea.Cmd
func (m ModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m ModalModel) View() string
```

Standard Bubble Tea model interface.

**Important**: `Update()` returns a non-nil `tea.Cmd` when it handles a message, `nil` when it doesn't. Parent should check `if cmd != nil` to determine if message was consumed.

```go
func (m ModalModel) Open() (ModalModel, tea.Cmd)
```

Opens the modal and returns the updated model.

```go
func (m ModalModel) Close() (ModalModel, tea.Cmd)
```

Closes the modal and returns the updated model.

```go
func (m ModalModel) SetSize(width, height int) ModalModel
```

Updates screen dimensions. Call when handling `tea.WindowSizeMsg`.

**Getter Methods** (access private fields):
```go
func (m ModalModel) Title() string
func (m ModalModel) Message() string
func (m ModalModel) Type() ModalType
func (m ModalModel) YesLabel() string
func (m ModalModel) NoLabel() string
func (m ModalModel) OKLabel() string
func (m ModalModel) IsOpen() bool
func (m ModalModel) FocusButton() int
func (m ModalModel) ScreenWidth() int
func (m ModalModel) ScreenHeight() int
func (m ModalModel) BorderStyle() lipgloss.Style
func (m ModalModel) TitleStyle() lipgloss.Style
func (m ModalModel) MessageStyle() lipgloss.Style
func (m ModalModel) ButtonStyle() lipgloss.Style
func (m ModalModel) FocusedButtonStyle() lipgloss.Style
func (m ModalModel) TitleAlign() lipgloss.Position
func (m ModalModel) MessageAlign() lipgloss.Position
func (m ModalModel) ButtonAlign() lipgloss.Position
```

**Wither Methods** (create modified copies):
```go
func (m ModalModel) WithTitle(title string) ModalModel
func (m ModalModel) WithMessage(message string) ModalModel
func (m ModalModel) WithYesLabel(label string) ModalModel
func (m ModalModel) WithNoLabel(label string) ModalModel
func (m ModalModel) WithOKLabel(label string) ModalModel
func (m ModalModel) WithBorderStyle(style lipgloss.Style) ModalModel
func (m ModalModel) WithTitleStyle(style lipgloss.Style) ModalModel
func (m ModalModel) WithMessageStyle(style lipgloss.Style) ModalModel
func (m ModalModel) WithButtonStyle(style lipgloss.Style) ModalModel
func (m ModalModel) WithFocusedButtonStyle(style lipgloss.Style) ModalModel
func (m ModalModel) WithTextAlign(align lipgloss.Position) ModalModel
func (m ModalModel) WithTitleAlign(align lipgloss.Position) ModalModel
func (m ModalModel) WithMessageAlign(align lipgloss.Position) ModalModel
func (m ModalModel) WithButtonAlign(align lipgloss.Position) ModalModel
```

Wither methods return a new `ModalModel` with the specified property modified.

### Overlay Functions

#### Method: OverlayModal (Recommended)

```go
func (m ModalModel) OverlayModal(background string) string
```

Automatically overlays the modal centered on the background view. Handles positioning internally.

**Example**:
```go
func (m model) View() string {
	baseView := renderYourView()
	return m.modal.OverlayModal(baseView) // Automatic centering
}
```

#### Function: OverlayModal (Advanced)

```go
func OverlayModal(background, foreground string, row, col int) string
```

Manually overlays modal at specified position. Use when you need full control over positioning.

**Parameters**:
- `background`: The base view (fully rendered string with ANSI codes)
- `foreground`: The modal view (fully rendered string with ANSI codes)
- `row`: Line number in background where foreground row 0 should appear (0-indexed)
- `col`: Display column in background where foreground col 0 should appear (0-indexed)

**Example**:
```go
func (m model) View() string {
	baseView := renderYourView()
	if m.modal.IsOpen() {
		modalView := m.modal.View()
		row, col := calculateCustomPosition()
		return teamodal.OverlayModal(baseView, modalView, row, col)
	}
	return baseView
}
```

```go
func EnsureTermGetSize(fd uintptr) (w int, h int, ok bool)
```

Robust terminal size detection with IDE fallbacks. Call before starting your Bubble Tea program if experiencing terminal size detection issues.

### Alignment

Control horizontal alignment of text and buttons using standard `lipgloss.Position` constants.

**Buttons default to centered** (standard modal UX). Text content can be easily aligned:

```go
// Left-align text content (title and message), buttons stay centered
modal := teamodal.NewYesNoModal("Proceed?", &teamodal.ModelArgs{
	ScreenWidth:  80,
	ScreenHeight: 24,
	Title:        "Confirmation",
	TextAlign:    lipgloss.Left, // Affects title and message
})

// Or control title and message individually
modal := teamodal.NewYesNoModal("Proceed?", &teamodal.ModelArgs{
	ScreenWidth:  80,
	ScreenHeight: 24,
	Title:        "Confirmation",
	TitleAlign:   lipgloss.Center, // Center title
	MessageAlign: lipgloss.Left,   // Left-align message
	// ButtonAlign defaults to Center
})

// Using fluent methods
modal = modal.WithTextAlign(lipgloss.Left)
// Or individually
modal = modal.WithTitleAlign(lipgloss.Center).
              WithMessageAlign(lipgloss.Left)

// Buttons can be overridden if needed (rare)
modal = modal.WithButtonAlign(lipgloss.Right)
```

**Alignment Options:**
- `lipgloss.Left` - Align to left edge
- `lipgloss.Center` - Center align (default)
- `lipgloss.Right` - Align to right edge

### Keyboard Customization

Customize keyboard shortcuts via the `Keys` field:

```go
modal := teamodal.NewYesNoModal("Proceed?", &teamodal.ModelArgs{
	ScreenWidth:  80,
	ScreenHeight: 24,
})

// Customize key bindings
modal.Keys.Confirm = key.NewBinding(key.WithKeys("y"))     // 'y' to confirm
modal.Keys.Cancel = key.NewBinding(key.WithKeys("n"))      // 'n' to cancel
modal.Keys.NextButton = key.NewBinding(key.WithKeys("tab"))
```

**Default Key Bindings:**
- `enter` - Confirm selection
- `esc` - Cancel/close
- `tab` - Next button (YesNo modals)
- `shift+tab` - Previous button (YesNo modals)
- `left` - Select Yes button (YesNo modals)
- `right` - Select No button (YesNo modals)

### Mouse Support

The modal supports full mouse interaction:

- **Click buttons** - Click any button to select and confirm
- **Hover feedback** - Button focus follows mouse hover (YesNo modals)
- **Automatic detection** - Modal tracks its position for accurate click detection

No additional code needed - mouse support works automatically.

### Wither Methods

Update modal properties at runtime using fluent wither methods:

```go
// Update content
modal = modal.WithTitle("New Title").
              WithMessage("Updated message")

// Update button labels
modal = modal.WithYesLabel("Accept").
              WithNoLabel("Decline").
              WithOKLabel("Got it")

// Update styling
modal = modal.WithBorderStyle(myBorderStyle).
              WithTitleStyle(myTitleStyle).
              WithFocusedButtonStyle(myFocusStyle)
```

**Available Withers:**
- `WithTitle(string)`, `WithMessage(string)`
- `WithYesLabel(string)`, `WithNoLabel(string)`, `WithOKLabel(string)`
- `WithTextAlign(lipgloss.Position)` - Sets both title and message alignment
- `WithTitleAlign(lipgloss.Position)`, `WithMessageAlign(lipgloss.Position)`, `WithButtonAlign(lipgloss.Position)`
- `WithBorderStyle(lipgloss.Style)`, `WithTitleStyle(lipgloss.Style)`, `WithMessageStyle(lipgloss.Style)`
- `WithButtonStyle(lipgloss.Style)`, `WithFocusedButtonStyle(lipgloss.Style)`

All wither methods return a new `ModalModel` instance (functional style).

### Styling

Customize modal appearance by providing styles in `ModelArgs`:

```go
modal := teamodal.NewYesNoModal("Proceed?", &teamodal.ModelArgs{
	ScreenWidth:  80,
	ScreenHeight: 24,
	Title:        "Confirmation",
	BorderStyle: lipgloss.NewStyle().
		BorderStyle(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("62")),
	FocusedButtonStyle: lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("230")).
		Background(lipgloss.Color("63")),
})
```

### Default Styles

```go
func DefaultBorderStyle() lipgloss.Style
func DefaultTitleStyle() lipgloss.Style
func DefaultMessageStyle() lipgloss.Style
func DefaultButtonStyle() lipgloss.Style
func DefaultFocusedButtonStyle() lipgloss.Style
```

Returns the default styles used when no options are provided.

## Modal Behavior Pattern

The modal implements clean message consumption without "infecting" the parent:

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var modal tea.Model

	// Let modal handle first
	modal, cmd = m.confirmDialog.Update(msg)
	if cmd != nil {
		m.confirmDialog = modal.(teamodal.ModalModel)
		return m, cmd  // Modal consumed message - we're done
	}

	// Modal didn't handle - parent processes
	switch msg := msg.(type) {
	case teamodal.AnsweredYesMsg:
		// Handle yes response
		return m, m.handleConfirmation()

	case teamodal.AnsweredNoMsg, teamodal.ClosedMsg:
		// Handle no/cancel response
		return m, nil
	}

	return m, nil
}
```

**Key insight**: The modal returns `nil` cmd when it doesn't handle the message, non-nil when it does. This prevents "modal infection" - parent doesn't need to check modal state for every key.

## Design

TeaModal follows the same proven patterns as teadd:

- **Modal Message Consumption**: Non-nil cmd signals "I handled this", nil signals "not for me"
- **String Compositing**: ANSI-aware overlay using `OverlayModal()`
- **ClearPath Style**: Single return, goto end, no else
- **Functional Options**: Methods return new model values

See [teadd documentation](../teadd/README.md) for more on these patterns.

## Examples

See `example/` directory for a working demonstration.

## License

MIT License - See LICENSE.txt
