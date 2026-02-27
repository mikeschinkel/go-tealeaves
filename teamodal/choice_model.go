package teamodal

import (
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// Orientation controls how choice buttons are laid out
type Orientation int

const (
	Horizontal Orientation = iota // Default: buttons in a row
	Vertical                      // Buttons stacked vertically
)

// ChoiceOption represents a single choice in the modal.
// Use capitalization in Label to indicate the hotkey visually (e.g. "Reorganize & Exit").
type ChoiceOption struct {
	Label  string // Display text, use caps for hotkey affordance: "Reorganize & Exit"
	Hotkey rune   // Optional: 'r' (case-insensitive, triggers without Tab+Enter)
	ID     string // Returned in ChoiceSelectedMsg to identify selection
}

// ChoiceModelArgs contains initialization arguments for ChoiceModel
type ChoiceModelArgs struct {
	ScreenWidth    int
	ScreenHeight   int
	Title          string         // Optional title above message
	Message        string         // Main message text
	Options        []ChoiceOption // 2-5 options
	DefaultIndex   int            // Which button is focused initially (0-based)
	Orientation    Orientation    // Horizontal (default) or Vertical
	ShowBrackets   *bool          // Show "[ ]" around labels; nil = true for Horizontal, false for Vertical
	AllowCancel    *bool          // Allow Esc to cancel; nil defaults to true
	ShowCancelHint *bool          // Show "[esc] Cancel" hint below options; nil defaults to AllowCancel

	// Style overrides (optional - defaults will be used if not provided)
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
	CancelKeyStyle     lipgloss.Style
	CancelTextStyle    lipgloss.Style
}

// ChoiceModel is a Bubble Tea model for multi-option confirmation dialogs
type ChoiceModel struct {
	Keys ChoiceKeyMap

	// Content
	title   string
	message string
	options []ChoiceOption

	// Layout
	orientation    Orientation
	showBrackets   bool
	allowCancel    bool
	showCancelHint bool

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
	cancelKeyStyle     lipgloss.Style
	cancelTextStyle    lipgloss.Style
}

// NewChoiceModel creates a new multi-option choice modal
func NewChoiceModel(args *ChoiceModelArgs) (m ChoiceModel) {
	if args == nil {
		args = &ChoiceModelArgs{}
	}

	m = ChoiceModel{
		Keys:               DefaultChoiceKeyMap(),
		title:              args.Title,
		message:            args.Message,
		options:            args.Options,
		orientation:        args.Orientation,
		focusButton:        args.DefaultIndex,
		screenWidth:        args.ScreenWidth,
		screenHeight:       args.ScreenHeight,
		borderStyle:        DefaultBorderStyle(),
		titleStyle:         DefaultTitleStyle(),
		messageStyle:       DefaultMessageStyle(),
		buttonStyle:        DefaultButtonStyle(),
		focusedButtonStyle: DefaultFocusedButtonStyle(),
		cancelKeyStyle:     DefaultCancelKeyStyle(),
		cancelTextStyle:    DefaultCancelTextStyle(),
	}

	// Clamp DefaultIndex to valid range
	if len(m.options) > 0 && (m.focusButton < 0 || m.focusButton >= len(m.options)) {
		m.focusButton = 0
	}

	// Resolve ShowBrackets: explicit *bool overrides, otherwise default by orientation
	if args.ShowBrackets != nil {
		m.showBrackets = *args.ShowBrackets
	}
	if args.ShowBrackets == nil {
		m.showBrackets = m.orientation == Horizontal
	}

	// Resolve AllowCancel: nil → true
	if args.AllowCancel != nil {
		m.allowCancel = *args.AllowCancel
	}
	if args.AllowCancel == nil {
		m.allowCancel = true
	}

	// Resolve ShowCancelHint: nil → follows AllowCancel
	if args.ShowCancelHint != nil {
		m.showCancelHint = *args.ShowCancelHint
	}
	if args.ShowCancelHint == nil {
		m.showCancelHint = m.allowCancel
	}

	// Apply custom styles if provided (check if non-zero)
	if args.BorderStyle.String() != "" {
		m.borderStyle = args.BorderStyle
	}
	if args.TitleStyle.String() != "" {
		m.titleStyle = args.TitleStyle
	}
	if args.MessageStyle.String() != "" {
		m.messageStyle = args.MessageStyle
	}
	if args.ButtonStyle.String() != "" {
		m.buttonStyle = args.ButtonStyle
	}
	if args.FocusedButtonStyle.String() != "" {
		m.focusedButtonStyle = args.FocusedButtonStyle
	}
	if args.CancelKeyStyle.String() != "" {
		m.cancelKeyStyle = args.CancelKeyStyle
	}
	if args.CancelTextStyle.String() != "" {
		m.cancelTextStyle = args.CancelTextStyle
	}

	return m
}

// Init implements tea.Model - returns nil (no initial command)
func (m ChoiceModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
func (m ChoiceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyMsg
	var ok bool
	var sizeMsg tea.WindowSizeMsg
	var pressedRune rune
	var i int
	var opt ChoiceOption

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
			opt = m.options[m.focusButton]
			i = m.focusButton
			cmd = func() tea.Msg {
				return ChoiceSelectedMsg{OptionID: opt.ID, Index: i}
			}
			goto end

		case key.Matches(keyMsg, m.Keys.Cancel):
			if !m.allowCancel {
				goto end
			}
			m.isOpen = false
			cmd = func() tea.Msg { return ChoiceCancelledMsg{} }
			goto end
		}

		// Check for hotkey press
		if keyMsg.Type == tea.KeyRunes && len(keyMsg.Runes) == 1 {
			pressedRune = unicode.ToLower(keyMsg.Runes[0])
			for i, opt = range m.options {
				if opt.Hotkey == 0 {
					continue
				}
				if unicode.ToLower(opt.Hotkey) != pressedRune {
					continue
				}
				m.isOpen = false
				cmd = func() tea.Msg {
					return ChoiceSelectedMsg{OptionID: opt.ID, Index: i}
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

// View implements tea.Model (FOLLOWS ClearPath)
func (m ChoiceModel) View() (view string) {
	var err error

	if !m.isOpen {
		goto end
	}

	view, err = m.renderModal()
	if err != nil {
		view = "Error: " + err.Error()
		goto end
	}

end:
	return view
}

// Open opens the modal and returns updated model
func (m ChoiceModel) Open() (ChoiceModel, tea.Cmd) {
	m.isOpen = true

	// Pre-calculate modal dimensions and position for overlay
	modalView, _ := m.renderModal()
	m.width, m.height, m.lastRow, m.lastCol = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)

	return m, nil
}

// Close closes the modal and returns updated model
func (m ChoiceModel) Close() (ChoiceModel, tea.Cmd) {
	m.isOpen = false
	return m, nil
}

// SetSize sets screen dimensions
func (m ChoiceModel) SetSize(width, height int) ChoiceModel {
	m.screenWidth = width
	m.screenHeight = height
	return m
}

// IsOpen returns whether the modal is currently open
func (m ChoiceModel) IsOpen() bool {
	return m.isOpen
}

// FocusButton returns the index of the currently focused button
func (m ChoiceModel) FocusButton() int {
	return m.focusButton
}

// OverlayModal renders the modal centered over the background view.
// Handles positioning automatically based on screen dimensions.
func (m ChoiceModel) OverlayModal(background string) (view string) {
	var row, col int
	var modalView string

	if !m.isOpen {
		view = background
		goto end
	}

	modalView = m.View()
	row = m.lastRow
	col = m.lastCol
	view = OverlayModal(background, modalView, row, col)

end:
	return view
}

// renderModal creates the modal box view (FOLLOWS ClearPath)
func (m ChoiceModel) renderModal() (view string, err error) {
	var content strings.Builder
	var titleLine string
	var messageLine string
	var renderedButtons string
	var buttonLine string
	var cancelHint string
	var maxWidth int
	var titleWidth int
	var messageWidth int
	var buttonWidth int
	var messageLines int
	var buttonLines int

	if !m.isOpen {
		goto end
	}

	// Render buttons first so we can measure actual styled width
	renderedButtons = m.renderButtons()

	// Calculate content widths
	titleWidth = len([]rune(m.title))
	messageWidth = m.calculateMessageWidth()
	buttonWidth = lipgloss.Width(renderedButtons)

	// Determine max width (add padding: 2 for left/right margins)
	maxWidth = titleWidth
	if messageWidth > maxWidth {
		maxWidth = messageWidth
	}
	if buttonWidth > maxWidth {
		maxWidth = buttonWidth
	}
	maxWidth = maxWidth + 4 // Add padding

	// Ensure minimum width
	if maxWidth < 30 {
		maxWidth = 30
	}

	// Store dimensions for overlay positioning
	// Width = content (maxWidth) + padding (4) + borders (2)
	m.width = maxWidth + 6

	// Height calculation for overlay positioning
	// Content: title (opt) + spacing + message lines + spacing + buttons
	// Chrome: padding top (1) + padding bottom (1) + border top (1) + border bottom (1) = 4 lines
	messageLines = strings.Count(m.message, "\n") + 1
	buttonLines = 1
	if m.orientation == Vertical {
		buttonLines = len(m.options)
	}
	m.height = messageLines + 6 + buttonLines // message + spacing after (2) + buttons + padding (2) + borders (2)
	if m.title != "" {
		m.height += 3 // title (1) + spacing after title (2)
	}
	if m.showCancelHint {
		m.height += 2 // blank line + hint line
	}

	// Render title (if present)
	if m.title != "" {
		titleLine = teautils.RenderAlignedLine(m.title, m.titleStyle, maxWidth, lipgloss.Center)
		content.WriteString(titleLine)
		content.WriteString("\n\n")
	}

	// Render message
	messageLine = teautils.RenderAlignedLine(m.message, m.messageStyle, maxWidth, lipgloss.Center)
	content.WriteString(messageLine)
	content.WriteString("\n\n")

	// Align buttons within content width
	if m.orientation == Vertical {
		// Center the bounding box of the button group, then left-align buttons within it.
		// Indent = (maxWidth - widestButton) / 2
		indent := (maxWidth - buttonWidth) / 2
		if indent < 0 {
			indent = 0
		}
		pad := strings.Repeat(" ", indent)
		for i, opt := range m.options {
			content.WriteString(pad)
			content.WriteString(m.renderButton(opt, i == m.focusButton))
			if i < len(m.options)-1 {
				content.WriteString("\n")
			}
		}
	}
	if m.orientation != Vertical {
		buttonLine = teautils.RenderAlignedLine(renderedButtons, lipgloss.NewStyle(), maxWidth, lipgloss.Center)
		content.WriteString(buttonLine)
	}

	if m.showCancelHint {
		cancelHint = m.cancelKeyStyle.Render("[esc]") + " " + m.cancelTextStyle.Render("Cancel")
		content.WriteString("\n\n")
		content.WriteString(teautils.RenderCenteredLine(cancelHint, lipgloss.NewStyle(), maxWidth))
	}

	// Apply border
	view = teautils.ApplyBoxBorder(m.borderStyle, content.String())

end:
	return view, err
}

// calculateMessageWidth returns the width needed for the message.
// For multi-line messages (containing \n), returns the max line width.
func (m ChoiceModel) calculateMessageWidth() (width int) {
	lines := strings.Split(m.message, "\n")
	for _, line := range lines {
		lineWidth := len([]rune(line))
		if lineWidth > width {
			width = lineWidth
		}
	}
	return width
}

// renderButtons renders the button row (horizontal) or block (vertical)
func (m ChoiceModel) renderButtons() (line string) {
	parts := make([]string, 0, len(m.options))

	for i, opt := range m.options {
		parts = append(parts, m.renderButton(opt, i == m.focusButton))
	}

	if m.orientation == Vertical {
		line = strings.Join(parts, "\n")
		goto done
	}

	line = strings.Join(parts, "  ")

done:
	return line
}

// renderButton renders a single button with the label as-is.
// Hotkey affordance is conveyed via capitalization in the label (developer's responsibility).
func (m ChoiceModel) renderButton(opt ChoiceOption, isFocused bool) (rendered string) {
	var style lipgloss.Style
	var label string

	style = m.buttonStyle
	if isFocused {
		style = m.focusedButtonStyle
	}

	label = opt.Label
	if m.showBrackets {
		label = "[ " + label + " ]"
	}

	rendered = style.Render(label)

	return rendered
}
