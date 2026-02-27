# Plan: Add ChoiceModel to teamodal

## Overview

Add a general-purpose multi-option modal (`ChoiceModel`) to the teamodal package. This modal presents N choices (2-5 buttons) with hotkey support, used for confirmation dialogs with more than 2 options.

**Use case example:** Exit confirmation with 3 options: "Reorganize & Exit", "Save & Exit", "Cancel"

## Files to Create

### 1. `choice_messages.go`

```go
package teamodal

// ChoiceSelectedMsg is sent when a choice is selected (Enter on focused button or hotkey)
type ChoiceSelectedMsg struct {
	OptionID string // The ID of the selected option
	Index    int    // The 0-based index of the selected option
}

// ChoiceCancelledMsg is sent when the modal is cancelled (Esc key)
type ChoiceCancelledMsg struct{}
```

### 2. `choice_keymap.go`

```go
package teamodal

import "github.com/charmbracelet/bubbles/key"

// ChoiceKeyMap defines key bindings for ChoiceModel
type ChoiceKeyMap struct {
	NextButton key.Binding // Tab, Right arrow
	PrevButton key.Binding // Shift+Tab, Left arrow
	Confirm    key.Binding // Enter
	Cancel     key.Binding // Esc
}

// DefaultChoiceKeyMap returns default key bindings for ChoiceModel
func DefaultChoiceKeyMap() ChoiceKeyMap {
	return ChoiceKeyMap{
		NextButton: key.NewBinding(
			key.WithKeys("tab", "right"),
			key.WithHelp("tab/→", "next"),
		),
		PrevButton: key.NewBinding(
			key.WithKeys("shift+tab", "left"),
			key.WithHelp("shift+tab/←", "prev"),
		),
		Confirm: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
	}
}
```

### 3. `choice_model.go`

**Types:**

```go
// ChoiceOption represents a single choice in the modal
type ChoiceOption struct {
	Label  string // Display text: "Reorganize & Exit"
	Hotkey rune   // Optional: 'r' (case-insensitive, triggers without Tab+Enter)
	ID     string // Returned in ChoiceSelectedMsg to identify selection
}

// ChoiceModelArgs contains initialization arguments for ChoiceModel
type ChoiceModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string         // Optional title above message
	Message      string         // Main message text
	Options      []ChoiceOption // 2-5 options
	DefaultIndex int            // Which button is focused initially (0-based)

	// Style overrides (optional - use same pattern as ModalModel)
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
	HotkeyStyle        lipgloss.Style // Style for hotkey character in button label
}

// ChoiceModel is a Bubble Tea model for multi-option confirmation dialogs
type ChoiceModel struct {
	Keys ChoiceKeyMap

	// Content
	title   string
	message string
	options []ChoiceOption

	// State
	isOpen       bool
	focusButton  int // Index of focused button (0-based)
	screenWidth  int
	screenHeight int

	// Calculated dimensions (for overlay positioning)
	width   int
	height  int
	lastRow int
	lastCol int

	// Styles
	borderStyle        lipgloss.Style
	titleStyle         lipgloss.Style
	messageStyle       lipgloss.Style
	buttonStyle        lipgloss.Style
	focusedButtonStyle lipgloss.Style
	hotkeyStyle        lipgloss.Style
}
```

**Constructor:**

```go
// NewChoiceModel creates a new multi-option choice modal
func NewChoiceModel(args *ChoiceModelArgs) ChoiceModel
```

**Key Behaviors in Update():**

1. **Tab / Right arrow:** Move focus to next button (wrap around)
2. **Shift+Tab / Left arrow:** Move focus to previous button (wrap around)
3. **Enter:** Select focused button, close modal, return `ChoiceSelectedMsg`
4. **Esc:** Close modal, return `ChoiceCancelledMsg`
5. **Hotkey press:** If any option has a matching hotkey (case-insensitive), select it immediately

**Hotkey matching logic:**
```go
// In Update(), after checking standard keys:
if keyMsg.Type == tea.KeyRunes && len(keyMsg.Runes) == 1 {
    pressedRune := unicode.ToLower(keyMsg.Runes[0])
    for i, opt := range m.options {
        if opt.Hotkey != 0 && unicode.ToLower(opt.Hotkey) == pressedRune {
            m.isOpen = false
            return m, func() tea.Msg {
                return ChoiceSelectedMsg{OptionID: opt.ID, Index: i}
            }
        }
    }
}
```

**Button rendering:**

- Buttons displayed horizontally with 2-space gap between them
- Format: `[ Label ]` (unfocused) vs highlighted `[ Label ]` (focused)
- If option has a hotkey, render that character with hotkeyStyle (e.g., underlined or different color)
- Example: `[ (R)eorganize & Exit ]` or `[ Reorganize & Exit ]` with 'R' highlighted

**Hotkey display in button:**
```go
// renderButton renders a single button, highlighting the hotkey character if present
func (m ChoiceModel) renderButton(opt ChoiceOption, isFocused bool) string {
    label := opt.Label
    style := m.buttonStyle
    if isFocused {
        style = m.focusedButtonStyle
    }

    // If hotkey exists, find and highlight it in the label
    if opt.Hotkey != 0 {
        // Find first occurrence of hotkey (case-insensitive) and wrap with hotkeyStyle
        // Implementation: find index, split label, apply hotkeyStyle to that char
    }

    return style.Render("[ " + label + " ]")
}
```

**Standard methods to implement:**
- `Init() tea.Cmd` - returns nil
- `Update(msg tea.Msg) (tea.Model, tea.Cmd)` - handles keys, window resize
- `View() string` - renders modal content
- `Open() (ChoiceModel, tea.Cmd)` - opens modal, calculates position
- `Close() (ChoiceModel, tea.Cmd)` - closes modal
- `SetSize(width, height int) ChoiceModel` - updates screen dimensions
- `IsOpen() bool` - returns open state
- `OverlayModal(background string) string` - renders modal over background

**Layout (renderModal):**
```
┌─────────────────────────────────────────┐
│                                         │
│              Title (if set)             │
│                                         │
│    Message text goes here, can be       │
│    multiple lines if needed.            │
│                                         │
│  [ Option 1 ]  [ Option 2 ]  [ Cancel ] │
│                                         │
└─────────────────────────────────────────┘
```

## Style Defaults

Use the same default styles as `ModalModel` for consistency:
- `DefaultBorderStyle()` - from styles.go
- `DefaultTitleStyle()` - from styles.go
- `DefaultMessageStyle()` - from styles.go
- `DefaultButtonStyle()` - from styles.go
- `DefaultFocusedButtonStyle()` - from styles.go

Add new default for hotkey:
```go
// DefaultHotkeyStyle returns the default style for hotkey characters in buttons
func DefaultHotkeyStyle() lipgloss.Style {
    return lipgloss.NewStyle().
        Underline(true).
        Bold(true)
}
```

## Implementation Notes

### Following House Rules

1. **ClearPath style:** Use `goto end` pattern, single return at end, no `else` statements
2. **Error handling:** Use doterr pattern (though this model likely won't have errors)
3. **No compound if statements:** No `if err := foo(); err != nil`
4. **Named return values:** Use named returns in function signatures

### Example ClearPath Update():

```go
func (m ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var keyMsg tea.KeyMsg
    var ok bool
    var sizeMsg tea.WindowSizeMsg

    if !m.isOpen {
        goto end
    }

    keyMsg, ok = msg.(tea.KeyMsg)
    if ok {
        switch {
        case key.Matches(keyMsg, m.Keys.NextButton):
            m.focusButton = (m.focusButton + 1) % len(m.options)
            cmd = func() tea.Msg { return nil }
            goto end

        case key.Matches(keyMsg, m.Keys.PrevButton):
            m.focusButton = (m.focusButton - 1 + len(m.options)) % len(m.options)
            cmd = func() tea.Msg { return nil }
            goto end

        case key.Matches(keyMsg, m.Keys.Confirm):
            m.isOpen = false
            opt := m.options[m.focusButton]
            cmd = func() tea.Msg {
                return ChoiceSelectedMsg{OptionID: opt.ID, Index: m.focusButton}
            }
            goto end

        case key.Matches(keyMsg, m.Keys.Cancel):
            m.isOpen = false
            cmd = func() tea.Msg { return ChoiceCancelledMsg{} }
            goto end
        }

        // Check for hotkey press
        if keyMsg.Type == tea.KeyRunes && len(keyMsg.Runes) == 1 {
            pressedRune := unicode.ToLower(keyMsg.Runes[0])
            for i, opt := range m.options {
                if opt.Hotkey == 0 {
                    continue
                }
                if unicode.ToLower(opt.Hotkey) != pressedRune {
                    continue
                }
                m.isOpen = false
                selectedOpt := opt
                selectedIdx := i
                cmd = func() tea.Msg {
                    return ChoiceSelectedMsg{OptionID: selectedOpt.ID, Index: selectedIdx}
                }
                goto end
            }
        }
    }

    sizeMsg, ok = msg.(tea.WindowSizeMsg)
    if ok {
        m.screenWidth = sizeMsg.Width
        m.screenHeight = sizeMsg.Height
        cmd = func() tea.Msg { return nil }
        goto end
    }

end:
    return m, cmd
}
```

## Testing

Create `choice_model_test.go` with tests for:

1. **Navigation:** Tab cycles focus forward, Shift+Tab cycles backward, wraps around
2. **Selection:** Enter selects focused option, returns correct ChoiceSelectedMsg
3. **Cancel:** Esc returns ChoiceCancelledMsg
4. **Hotkeys:** Pressing hotkey character selects correct option (test both cases)
5. **Default focus:** DefaultIndex sets initial focus correctly

## Usage Example (for gomtui)

```go
// In batch_assignment_model.go

// Create the exit confirmation modal
exitModal := teamodal.NewChoiceModel(&teamodal.ChoiceModelArgs{
    Title:   "Exit with Pending Changes",
    Message: "Some files have been reassigned but not reorganized.\nHow would you like to proceed?",
    Options: []teamodal.ChoiceOption{
        {Label: "Reorganize & Exit", Hotkey: 'r', ID: "reorganize"},
        {Label: "Save & Exit", Hotkey: 's', ID: "save"},
        {Label: "Cancel", Hotkey: 0, ID: "cancel"},
    },
    DefaultIndex:  2, // Cancel is default
    ScreenWidth:   m.terminalWidth,
    ScreenHeight:  m.terminalHeight,
})

// Handle the response
case teamodal.ChoiceSelectedMsg:
    switch msg.OptionID {
    case "reorganize":
        m = m.reorganizeTree()
        cmd = tea.Batch(m.requestSaveCmd(), func() tea.Msg { return DrillUpMsg{} })
    case "save":
        cmd = tea.Batch(m.requestSaveCmd(), func() tea.Msg { return DrillUpMsg{} })
    case "cancel":
        // Do nothing, stay in UI
    }

case teamodal.ChoiceCancelledMsg:
    // Same as "cancel" option - stay in UI
```

## Checklist

- [ ] Create `choice_messages.go`
- [ ] Create `choice_keymap.go`
- [ ] Create `choice_model.go`
- [ ] Add `DefaultHotkeyStyle()` to `styles.go`
- [ ] Create `choice_model_test.go`
- [ ] Verify ClearPath style compliance
- [ ] Test hotkey functionality (both upper and lower case)
- [ ] Test navigation wrapping
- [ ] Test overlay positioning
