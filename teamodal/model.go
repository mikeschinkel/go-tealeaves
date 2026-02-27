package teamodal

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// ModalType represents the type of modal dialog
type ModalType int

const (
	// ModalTypeOK shows a single OK button (alert dialog)
	ModalTypeOK ModalType = iota
	// ModalTypeYesNo shows Yes and No buttons (confirmation dialog)
	ModalTypeYesNo
)

// ModelArgs contains initialization arguments for ModalModel
type ModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string
	DefaultYes   bool // For YesNo modals: default focus to Yes (true) or No (false)
	YesLabel     string
	NoLabel      string
	OKLabel      string

	// Alignment (optional - defaults to lipgloss.Center)
	// Use TextAlign to set both title and message, or individual fields for fine control
	TextAlign    lipgloss.Position // Horizontal alignment for title and message
	TitleAlign   lipgloss.Position // Horizontal alignment for title (overrides TextAlign)
	MessageAlign lipgloss.Position // Horizontal alignment for message (overrides TextAlign)
	ButtonAlign  lipgloss.Position // Horizontal alignment for buttons (defaults to Center)

	// Styling (optional - defaults will be used if not provided)
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
}

// ModalModel is a Bubble Tea model for modal dialogs
type ModalModel struct {
	Keys ModalKeyMap // Keyboard bindings

	// Content
	title   string
	message string
	typ     ModalType

	// Button labels
	yesLabel string
	noLabel  string
	okLabel  string

	// State
	isOpen       bool
	focusButton  int // Index of focused button (0 = first button, 1 = second button for YesNo)
	screenWidth  int
	screenHeight int

	// Calculated dimensions (for overlay positioning)
	width   int
	height  int
	lastRow int // Row where modal was last rendered
	lastCol int // Column where modal was last rendered

	// Alignment (internal pointers to distinguish unset from Left=0)
	titleAlign   *lipgloss.Position
	messageAlign *lipgloss.Position
	buttonAlign  *lipgloss.Position

	// Styling (all private with getters/withers)
	borderStyle        lipgloss.Style
	titleStyle         lipgloss.Style
	messageStyle       lipgloss.Style
	buttonStyle        lipgloss.Style
	focusedButtonStyle lipgloss.Style
}

// NewOKModal creates a new OK-only modal (alert dialog)
func NewOKModal(message string, args *ModelArgs) (m ModalModel) {
	m = newModalModel(message, ModalTypeOK, args)
	return m
}

// NewYesNoModal creates a new Yes/No confirmation modal
func NewYesNoModal(message string, args *ModelArgs) (m ModalModel) {
	m = newModalModel(message, ModalTypeYesNo, args)
	return m
}

func newModalModel(message string, modalType ModalType, args *ModelArgs) (m ModalModel) {
	if args == nil {
		args = &ModelArgs{}
	}

	m = ModalModel{
		Keys:               DefaultModalKeyMap(),
		title:              args.Title,
		message:            message,
		typ:                modalType,
		yesLabel:           "Yes",
		noLabel:            "No",
		okLabel:            "OK",
		isOpen:             false,
		focusButton:        0,
		screenWidth:        args.ScreenWidth,
		screenHeight:       args.ScreenHeight,
		borderStyle:        DefaultBorderStyle(),
		titleStyle:         DefaultTitleStyle(),
		messageStyle:       DefaultMessageStyle(),
		buttonStyle:        DefaultButtonStyle(),
		focusedButtonStyle: DefaultFocusedButtonStyle(),
	}
	// Note: alignment fields are nil by default, getters return Center

	// Apply custom labels if provided
	if args.YesLabel != "" {
		m.yesLabel = args.YesLabel
	}
	if args.NoLabel != "" {
		m.noLabel = args.NoLabel
	}
	if args.OKLabel != "" {
		m.okLabel = args.OKLabel
	}

	// Set default focus for YesNo modals
	if modalType == ModalTypeYesNo && !args.DefaultYes {
		m.focusButton = 1 // Focus No button
	}

	// Apply text alignment to title and message (if provided)
	if args.TextAlign != 0 {
		align := args.TextAlign
		m.titleAlign = &align
		m.messageAlign = &align
	}

	// Then apply specific alignments (these override TextAlign if both provided)
	if args.TitleAlign != 0 {
		align := args.TitleAlign
		m.titleAlign = &align
	}
	if args.MessageAlign != 0 {
		align := args.MessageAlign
		m.messageAlign = &align
	}
	if args.ButtonAlign != 0 {
		align := args.ButtonAlign
		m.buttonAlign = &align
	}

	// Apply custom styles if provided (check if non-zero)
	borderStyleStr := args.BorderStyle.String()
	if borderStyleStr != "" {
		m.borderStyle = args.BorderStyle
	}
	titleStyleStr := args.TitleStyle.String()
	if titleStyleStr != "" {
		m.titleStyle = args.TitleStyle
	}
	messageStyleStr := args.MessageStyle.String()
	if messageStyleStr != "" {
		m.messageStyle = args.MessageStyle
	}
	buttonStyleStr := args.ButtonStyle.String()
	if buttonStyleStr != "" {
		m.buttonStyle = args.ButtonStyle
	}
	focusedButtonStyleStr := args.FocusedButtonStyle.String()
	if focusedButtonStyleStr != "" {
		m.focusedButtonStyle = args.FocusedButtonStyle
	}

	return m
}

// Init implements tea.Model - returns nil (no initial command)
func (m ModalModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
func (m ModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyMsg
	var ok bool
	var sizeMsg tea.WindowSizeMsg
	var mouseMsg tea.MouseMsg

	if !m.isOpen {
		goto end // Not open = nil cmd = didn't handle
	}

	// Try as KeyMsg first
	keyMsg, ok = msg.(tea.KeyMsg)
	if ok {
		switch {
		case key.Matches(keyMsg, m.Keys.NextButton, m.Keys.PrevButton):
			// Switch focus between buttons (only for YesNo modals)
			if m.typ == ModalTypeYesNo {
				m.focusButton = 1 - m.focusButton // Toggle between 0 and 1
			}
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.SelectLeft):
			// Move focus to left button (Yes = 0) for YesNo modals
			if m.typ == ModalTypeYesNo {
				m.focusButton = 0
			}
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.SelectRight):
			// Move focus to right button (No = 1) for YesNo modals
			if m.typ == ModalTypeYesNo {
				m.focusButton = 1
			}
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Confirm):
			// Confirm selection
			m.isOpen = false
			if m.typ == ModalTypeOK {
				cmd = func() tea.Msg { return ClosedMsg{} }
			} else if m.typ == ModalTypeYesNo {
				if m.focusButton == 0 {
					cmd = func() tea.Msg { return AnsweredYesMsg{} }
				} else {
					cmd = func() tea.Msg { return AnsweredNoMsg{} }
				}
			}
			goto end

		case key.Matches(keyMsg, m.Keys.Cancel):
			// Cancel/close
			m.isOpen = false
			if m.typ == ModalTypeOK {
				cmd = func() tea.Msg { return ClosedMsg{} }
			} else {
				cmd = func() tea.Msg { return AnsweredNoMsg{} }
			}
			goto end
		}
	}

	// Try as WindowSizeMsg
	sizeMsg, ok = msg.(tea.WindowSizeMsg)
	if ok {
		m.screenWidth = sizeMsg.Width
		m.screenHeight = sizeMsg.Height
		cmd = func() tea.Msg { return nil }
		goto end
	}

	// Try as MouseMsg
	mouseMsg, ok = msg.(tea.MouseMsg)
	if ok {
		switch mouseMsg.Type {
		case tea.MouseLeft:
			// Handle button clicks
			if m.isClickOnButton(mouseMsg.X, mouseMsg.Y) {
				// Close modal and emit appropriate message
				m.isOpen = false
				cmd = m.handleButtonClick(mouseMsg.X, mouseMsg.Y)
				goto end
			}

		case tea.MouseMotion:
			// Handle mouse hover for visual feedback (only for YesNo modals)
			if m.typ == ModalTypeYesNo && m.isClickOnButton(mouseMsg.X, mouseMsg.Y) {
				// Update focus to match hovered button
				hoveredButton := m.getHoveredButton(mouseMsg.X, mouseMsg.Y)
				if hoveredButton >= 0 && hoveredButton != m.focusButton {
					m.focusButton = hoveredButton
					cmd = func() tea.Msg { return nil }
				}
			}
		}
	}

end:
	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m ModalModel) View() (view string) {
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
func (m ModalModel) Open() (ModalModel, tea.Cmd) {
	m.isOpen = true

	// Pre-calculate modal dimensions and position for mouse click detection
	modalView, _ := m.renderModal()

	// Use helper to measure and center the modal
	m.width, m.height, m.lastRow, m.lastCol = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)

	return m, nil
}

// Close closes the modal and returns updated model
func (m ModalModel) Close() (ModalModel, tea.Cmd) {
	m.isOpen = false
	return m, nil
}

// SetSize sets screen dimensions
func (m ModalModel) SetSize(width, height int) ModalModel {
	m.screenWidth = width
	m.screenHeight = height
	return m
}

// TitleAlign returns the title alignment (defaults to Center if not set)
func (m ModalModel) TitleAlign() lipgloss.Position {
	if m.titleAlign == nil {
		return lipgloss.Center
	}
	return *m.titleAlign
}

// MessageAlign returns the message alignment (defaults to Center if not set)
func (m ModalModel) MessageAlign() lipgloss.Position {
	if m.messageAlign == nil {
		return lipgloss.Center
	}
	return *m.messageAlign
}

// ButtonAlign returns the button alignment (defaults to Center if not set)
func (m ModalModel) ButtonAlign() lipgloss.Position {
	if m.buttonAlign == nil {
		return lipgloss.Center
	}
	return *m.buttonAlign
}

// WithTextAlign sets the horizontal alignment for text content (title and message, not buttons)
func (m ModalModel) WithTextAlign(align lipgloss.Position) ModalModel {
	m.titleAlign = &align
	m.messageAlign = &align
	return m
}

// WithTitleAlign sets the horizontal alignment for the title
func (m ModalModel) WithTitleAlign(align lipgloss.Position) ModalModel {
	m.titleAlign = &align
	return m
}

// WithMessageAlign sets the horizontal alignment for the message
func (m ModalModel) WithMessageAlign(align lipgloss.Position) ModalModel {
	m.messageAlign = &align
	return m
}

// WithButtonAlign sets the horizontal alignment for the buttons
func (m ModalModel) WithButtonAlign(align lipgloss.Position) ModalModel {
	m.buttonAlign = &align
	return m
}

// OverlayModal renders the modal centered over the background view.
// Handles positioning automatically based on screen dimensions.
func (m ModalModel) OverlayModal(background string) (view string) {
	var row, col int
	var modalView string

	if !m.isOpen {
		view = background
		goto end
	}

	// Render modal view
	modalView = m.View()

	// Use pre-calculated position from Open() method
	// (Position is stored when modal opens, not during rendering)
	row = m.lastRow
	col = m.lastCol

	// Overlay the rendered modal at calculated position
	view = OverlayModal(background, modalView, row, col)

end:
	return view
}

// renderModal creates the modal box view (FOLLOWS ClearPath)
func (m ModalModel) renderModal() (view string, err error) {
	var content strings.Builder
	var titleLine string
	var messageLine string
	var buttonLine string
	var maxWidth int
	var titleWidth int
	var messageWidth int
	var buttonWidth int
	var messageLines int

	if !m.isOpen {
		goto end
	}

	// Calculate content widths
	titleWidth = len([]rune(m.title))
	messageWidth = m.calculateMessageWidth()
	buttonWidth = m.calculateButtonWidth()

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
	m.height = messageLines + 7 // message + spacing after (2) + buttons (1) + padding (2) + borders (2)
	if m.title != "" {
		m.height += 3 // title (1) + spacing after title (2)
	}

	// Render title (if present)
	if m.title != "" {
		titleLine = teautils.RenderAlignedLine(m.title, m.titleStyle, maxWidth, m.TitleAlign())
		content.WriteString(titleLine)
		content.WriteString("\n\n")
	}

	// Render message
	messageLine = teautils.RenderAlignedLine(m.message, m.messageStyle, maxWidth, m.MessageAlign())
	content.WriteString(messageLine)
	content.WriteString("\n\n")

	// Render buttons (already styled by renderButtons, just align)
	buttonLine = teautils.RenderAlignedLine(m.renderButtons(), lipgloss.NewStyle(), maxWidth, m.ButtonAlign())
	content.WriteString(buttonLine)

	// Apply border
	view = teautils.ApplyBoxBorder(m.borderStyle, content.String())

end:
	return view, err
}

// calculateMessageWidth returns the width needed for the message.
// For multi-line messages (containing \n), returns the max line width.
func (m ModalModel) calculateMessageWidth() (width int) {
	lines := strings.Split(m.message, "\n")
	for _, line := range lines {
		lineWidth := len([]rune(line))
		if lineWidth > width {
			width = lineWidth
		}
	}
	return width
}

// calculateButtonWidth returns the total width needed for buttons
func (m ModalModel) calculateButtonWidth() (width int) {
	if m.typ == ModalTypeOK {
		width = len([]rune(m.okLabel)) + 4 // +4 for button padding
	} else if m.typ == ModalTypeYesNo {
		yesWidth := len([]rune(m.yesLabel)) + 4
		noWidth := len([]rune(m.noLabel)) + 4
		width = yesWidth + noWidth + 2 // +2 for space between buttons
	}
	return width
}

// renderButtons renders the button line (FOLLOWS ClearPath)
func (m ModalModel) renderButtons() (line string) {
	var yesButton string
	var noButton string
	var okButton string

	if m.typ == ModalTypeOK {
		okButton = m.buttonStyle.Render("[ " + m.okLabel + " ]")
		if m.focusButton == 0 {
			okButton = m.focusedButtonStyle.Render("[ " + m.okLabel + " ]")
		}
		line = okButton
		goto end
	}

	// YesNo buttons
	yesButton = m.buttonStyle.Render("[ " + m.yesLabel + " ]")
	if m.focusButton == 0 {
		yesButton = m.focusedButtonStyle.Render("[ " + m.yesLabel + " ]")
	}

	noButton = m.buttonStyle.Render("[ " + m.noLabel + " ]")
	if m.focusButton == 1 {
		noButton = m.focusedButtonStyle.Render("[ " + m.noLabel + " ]")
	}

	line = yesButton + "  " + noButton

end:
	return line
}

// =============================================================================
// Getters
// =============================================================================

// Title returns the modal title
func (m ModalModel) Title() string {
	return m.title
}

// Message returns the modal message
func (m ModalModel) Message() string {
	return m.message
}

// Type returns the modal type
func (m ModalModel) Type() ModalType {
	return m.typ
}

// YesLabel returns the Yes button label
func (m ModalModel) YesLabel() string {
	return m.yesLabel
}

// NoLabel returns the No button label
func (m ModalModel) NoLabel() string {
	return m.noLabel
}

// OKLabel returns the OK button label
func (m ModalModel) OKLabel() string {
	return m.okLabel
}

// IsOpen returns whether the modal is currently open
func (m ModalModel) IsOpen() bool {
	return m.isOpen
}

// FocusButton returns the index of the currently focused button
func (m ModalModel) FocusButton() int {
	return m.focusButton
}

// ScreenWidth returns the screen width
func (m ModalModel) ScreenWidth() int {
	return m.screenWidth
}

// ScreenHeight returns the screen height
func (m ModalModel) ScreenHeight() int {
	return m.screenHeight
}

// BorderStyle returns the border style
func (m ModalModel) BorderStyle() lipgloss.Style {
	return m.borderStyle
}

// TitleStyle returns the title style
func (m ModalModel) TitleStyle() lipgloss.Style {
	return m.titleStyle
}

// MessageStyle returns the message style
func (m ModalModel) MessageStyle() lipgloss.Style {
	return m.messageStyle
}

// ButtonStyle returns the unfocused button style
func (m ModalModel) ButtonStyle() lipgloss.Style {
	return m.buttonStyle
}

// FocusedButtonStyle returns the focused button style
func (m ModalModel) FocusedButtonStyle() lipgloss.Style {
	return m.focusedButtonStyle
}

// =============================================================================
// Withers
// =============================================================================

// WithTitle returns a copy with the specified title
func (m ModalModel) WithTitle(title string) ModalModel {
	m.title = title
	return m
}

// WithMessage returns a copy with the specified message
func (m ModalModel) WithMessage(message string) ModalModel {
	m.message = message
	return m
}

// WithYesLabel returns a copy with the specified Yes button label
func (m ModalModel) WithYesLabel(label string) ModalModel {
	m.yesLabel = label
	return m
}

// WithNoLabel returns a copy with the specified No button label
func (m ModalModel) WithNoLabel(label string) ModalModel {
	m.noLabel = label
	return m
}

// WithOKLabel returns a copy with the specified OK button label
func (m ModalModel) WithOKLabel(label string) ModalModel {
	m.okLabel = label
	return m
}

// WithBorderStyle returns a copy with the specified border style
func (m ModalModel) WithBorderStyle(style lipgloss.Style) ModalModel {
	m.borderStyle = style
	return m
}

// WithTitleStyle returns a copy with the specified title style
func (m ModalModel) WithTitleStyle(style lipgloss.Style) ModalModel {
	m.titleStyle = style
	return m
}

// WithMessageStyle returns a copy with the specified message style
func (m ModalModel) WithMessageStyle(style lipgloss.Style) ModalModel {
	m.messageStyle = style
	return m
}

// WithButtonStyle returns a copy with the specified button style
func (m ModalModel) WithButtonStyle(style lipgloss.Style) ModalModel {
	m.buttonStyle = style
	return m
}

// WithFocusedButtonStyle returns a copy with the specified focused button style
func (m ModalModel) WithFocusedButtonStyle(style lipgloss.Style) ModalModel {
	m.focusedButtonStyle = style
	return m
}

// =============================================================================
// Mouse Click Detection
// =============================================================================

// isClickOnButton checks if click coordinates are within the button row
func (m ModalModel) isClickOnButton(x, y int) bool {
	// Button row is: lastRow + height - 3
	// -3 accounts for: border (1) + padding (1) + button line offset (1)
	buttonRow := m.lastRow + m.height - 3
	return y == buttonRow
}

// getHoveredButton determines which button (0=Yes, 1=No) the mouse is hovering over
// Returns -1 if not hovering over a button
func (m ModalModel) getHoveredButton(x, y int) int {
	if m.typ != ModalTypeYesNo {
		return -1
	}

	// Calculate button X positions based on alignment
	buttonLine := m.renderButtons() // "[ Yes ]  [ No ]"
	yesWidth := len([]rune("[ " + m.yesLabel + " ]"))

	// Click detection based on alignment
	// Note: Initial implementation assumes center alignment
	buttonStartCol := m.lastCol + (m.width-ansi.StringWidth(buttonLine))/2

	if x >= buttonStartCol && x < buttonStartCol+yesWidth {
		return 0 // Yes button
	} else if x >= buttonStartCol+yesWidth {
		return 1 // No button
	}

	return -1 // Not on a button
}

// handleButtonClick determines which button was clicked and returns appropriate message
func (m ModalModel) handleButtonClick(x, y int) tea.Cmd {
	// For OK modals, any click on button row closes the modal
	if m.typ == ModalTypeOK {
		return func() tea.Msg { return ClosedMsg{} }
	}

	// For YesNo modals, determine which button was clicked
	hoveredButton := m.getHoveredButton(x, y)
	if hoveredButton == 0 {
		// Clicked Yes button
		return func() tea.Msg { return AnsweredYesMsg{} }
	}
	// Clicked No button (or anywhere else on button row)
	return func() tea.Msg { return AnsweredNoMsg{} }
}
