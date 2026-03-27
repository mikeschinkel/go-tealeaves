package teamodal

import (
	"log/slog"
	"strings"
	"unicode"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// MultiSelectItem is the constraint interface for items displayed in MultiSelectModel
type MultiSelectItem interface {
	// ID returns a unique identifier for the item
	ID() string
	// Label returns the display text for the item
	Label() string
}

// MultiSelectButton represents a button in the multi-select modal
type MultiSelectButton struct {
	Label  string // Display text for the button
	Hotkey rune   // Optional single-char hotkey (case-insensitive)
	ID     string // Returned in message to identify which button
}

// MultiSelectModelArgs provides optional configuration for MultiSelectModel
type MultiSelectModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string               // Title text at top of modal
	Message      string               // Optional description above the list
	Footer       string               // Optional note below the list (e.g. warning text)
	Buttons      []MultiSelectButton  // Action buttons (2-4)
	AllChecked               bool    // Default: true — all items start checked
	MaxVisible               int     // Default: 8 — scrollable if more items
	NoCancel                 bool    // Don't auto-append Cancel button (default: false = auto-cancel)
	CancelEmitsButtonPressed bool    // Cancel button emits ButtonPressedMsg instead of CancelledMsg

	// Style overrides (optional — defaults applied)
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	FooterStyle        lipgloss.Style
	ItemStyle          lipgloss.Style
	SelectedItemStyle  lipgloss.Style
	CheckedStyle       lipgloss.Style
	UncheckedStyle     lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
}

// multiSelectFocus tracks which UI element has focus
type multiSelectFocus int

const (
	focusList   multiSelectFocus = iota // Focus is on the checkbox list
	focusButton                         // Focus is on a button (buttonIdx determines which)
)

// MultiSelectModel is a Bubble Tea model for multi-select checkbox list modals
type MultiSelectModel[T MultiSelectItem] struct {
	// Public
	Keys   MultiSelectKeyMap
	Logger *slog.Logger

	// Content
	title   string
	message string
	footer  string
	items   []T
	buttons []MultiSelectButton

	// Selection state
	checked                  map[string]bool // item.ID() → checked state
	cancelEmitsButtonPressed bool            // Cancel button emits ButtonPressedMsg instead of CancelledMsg

	// Navigation state
	cursor    int              // Currently highlighted item index
	offset    int              // Scroll offset for viewport
	focus     multiSelectFocus // Which area has focus
	buttonIdx int              // Which button is focused (when focus == focusButton)

	// Dimensions
	isOpen       bool
	screenWidth  int
	screenHeight int
	maxVisible   int // Max items visible in viewport

	// Rendering (calculated on Open)
	width   int
	height  int
	lastRow int // Row where modal was last rendered
	lastCol int // Column where modal was last rendered

	// Styles
	borderStyle        lipgloss.Style
	titleStyle         lipgloss.Style
	messageStyle       lipgloss.Style
	footerStyle        lipgloss.Style
	itemStyle          lipgloss.Style
	selectedItemStyle  lipgloss.Style
	checkedStyle       lipgloss.Style
	uncheckedStyle     lipgloss.Style
	buttonStyle        lipgloss.Style
	focusedButtonStyle lipgloss.Style
}

// NewMultiSelectModel creates a new MultiSelectModel with the given items
func NewMultiSelectModel[T MultiSelectItem](items []T, args *MultiSelectModelArgs) (m MultiSelectModel[T]) {
	if args == nil {
		args = &MultiSelectModelArgs{
			AllChecked: true,
		}
	}

	maxVisible := args.MaxVisible
	if maxVisible <= 0 {
		maxVisible = 8
	}

	m = MultiSelectModel[T]{
		Keys:               DefaultMultiSelectKeyMap(),
		title:              args.Title,
		message:            args.Message,
		footer:             args.Footer,
		items:              items,
		buttons:            args.Buttons,
		cursor:             0,
		offset:             0,
		focus:              focusList,
		buttonIdx:          0,
		isOpen:             false,
		screenWidth:        args.ScreenWidth,
		screenHeight:       args.ScreenHeight,
		maxVisible:         maxVisible,
		borderStyle:        DefaultBorderStyle(),
		titleStyle:         DefaultTitleStyle(),
		messageStyle:       DefaultMessageStyle(),
		footerStyle:        DefaultMultiSelectFooterStyle(),
		itemStyle:          DefaultListItemStyle(),
		selectedItemStyle:  DefaultSelectedItemStyle(),
		checkedStyle:       DefaultCheckedStyle(),
		uncheckedStyle:     DefaultUncheckedStyle(),
		buttonStyle:        DefaultButtonStyle(),
		focusedButtonStyle: DefaultFocusedButtonStyle(),
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
	if args.FooterStyle.String() != "" {
		m.footerStyle = args.FooterStyle
	}
	if args.ItemStyle.String() != "" {
		m.itemStyle = args.ItemStyle
	}
	if args.SelectedItemStyle.String() != "" {
		m.selectedItemStyle = args.SelectedItemStyle
	}
	if args.CheckedStyle.String() != "" {
		m.checkedStyle = args.CheckedStyle
	}
	if args.UncheckedStyle.String() != "" {
		m.uncheckedStyle = args.UncheckedStyle
	}
	if args.ButtonStyle.String() != "" {
		m.buttonStyle = args.ButtonStyle
	}
	if args.FocusedButtonStyle.String() != "" {
		m.focusedButtonStyle = args.FocusedButtonStyle
	}

	// Auto-append Cancel button unless opted out
	if !args.NoCancel {
		m.buttons = append(m.buttons, MultiSelectButton{
			Label: "Cancel",
			ID:    "cancel",
			// No Hotkey — Esc already serves as keyboard shortcut
		})
		m.cancelEmitsButtonPressed = args.CancelEmitsButtonPressed
	}

	// Initialize checked state
	m.checked = make(map[string]bool, len(items))
	for _, item := range items {
		m.checked[item.ID()] = args.AllChecked
	}

	return m
}

// Init implements tea.Model - returns nil (no initial command)
func (m MultiSelectModel[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
func (m MultiSelectModel[T]) Update(msg tea.Msg) (MultiSelectModel[T], tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var sizeMsg tea.WindowSizeMsg
	var pressedRune rune

	if !m.isOpen {
		goto end // Not open = nil cmd = didn't handle
	}

	keyMsg, ok = msg.(tea.KeyPressMsg)
	if !ok {
		// Try as WindowSizeMsg
		sizeMsg, ok = msg.(tea.WindowSizeMsg)
		if ok {
			m.screenWidth = sizeMsg.Width
			m.screenHeight = sizeMsg.Height
			cmd = func() tea.Msg { return nil }
		}
		goto end
	}

	// Cancel (esc) — always available
	if key.Matches(keyMsg, m.Keys.Cancel) {
		m.isOpen = false
		cmd = func() tea.Msg { return MultiSelectCancelledMsg{} }
		goto end
	}

	// Check for button hotkey press (works from any focus state)
	if len(keyMsg.Text) == 1 {
		pressedRune = unicode.ToLower([]rune(keyMsg.Text)[0])
		for _, btn := range m.buttons {
			if btn.Hotkey == 0 {
				continue
			}
			if unicode.ToLower(btn.Hotkey) != pressedRune {
				continue
			}
			m.isOpen = false
			if btn.ID == "cancel" && !m.cancelEmitsButtonPressed {
				cmd = func() tea.Msg { return MultiSelectCancelledMsg{} }
			} else {
				selected := m.selectedItems()
				cmd = func() tea.Msg {
					return MultiSelectButtonPressedMsg[T]{ButtonID: btn.ID, Selected: selected}
				}
			}
			goto end
		}
	}

	// Focus cycling (tab / shift+tab) — always available
	if key.Matches(keyMsg, m.Keys.NextFocus) {
		m = m.cycleFocusForward()
		cmd = func() tea.Msg { return nil }
		goto end
	}
	if key.Matches(keyMsg, m.Keys.PrevFocus) {
		m = m.cycleFocusBackward()
		cmd = func() tea.Msg { return nil }
		goto end
	}

	// Dispatch based on focus
	if m.focus == focusList {
		m, cmd = m.handleListKeys(keyMsg)
		goto end
	}

	m, cmd = m.handleButtonKeys(keyMsg)

end:
	return m, cmd
}

// handleListKeys handles key events when focus is on the list
func (m MultiSelectModel[T]) handleListKeys(keyMsg tea.KeyPressMsg) (MultiSelectModel[T], tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(keyMsg, m.Keys.Up):
		m = m.moveCursorUp()
		cmd = func() tea.Msg { return nil }

	case key.Matches(keyMsg, m.Keys.Down):
		m = m.moveCursorDown()
		cmd = func() tea.Msg { return nil }

	case key.Matches(keyMsg, m.Keys.Toggle):
		if len(m.items) > 0 {
			item := m.items[m.cursor]
			m.checked[item.ID()] = !m.checked[item.ID()]
			cmd = func() tea.Msg { return nil }
		}

	case key.Matches(keyMsg, m.Keys.Confirm):
		// Enter on list = toggle checkbox
		if len(m.items) > 0 {
			item := m.items[m.cursor]
			m.checked[item.ID()] = !m.checked[item.ID()]
			cmd = func() tea.Msg { return nil }
		}
	}

	return m, cmd
}

// handleButtonKeys handles key events when focus is on the buttons
func (m MultiSelectModel[T]) handleButtonKeys(keyMsg tea.KeyPressMsg) (MultiSelectModel[T], tea.Cmd) {
	var cmd tea.Cmd

	switch {
	case key.Matches(keyMsg, m.Keys.NextButton):
		if len(m.buttons) > 0 {
			m.buttonIdx = (m.buttonIdx + 1) % len(m.buttons)
			cmd = func() tea.Msg { return nil }
		}

	case key.Matches(keyMsg, m.Keys.PrevButton):
		if len(m.buttons) > 0 {
			m.buttonIdx = (m.buttonIdx - 1 + len(m.buttons)) % len(m.buttons)
			cmd = func() tea.Msg { return nil }
		}

	case key.Matches(keyMsg, m.Keys.Confirm):
		// Enter on button = activate
		if len(m.buttons) > 0 {
			m.isOpen = false
			btn := m.buttons[m.buttonIdx]
			if btn.ID == "cancel" && !m.cancelEmitsButtonPressed {
				cmd = func() tea.Msg { return MultiSelectCancelledMsg{} }
			} else {
				selected := m.selectedItems()
				cmd = func() tea.Msg {
					return MultiSelectButtonPressedMsg[T]{ButtonID: btn.ID, Selected: selected}
				}
			}
		}
	}

	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m MultiSelectModel[T]) View() tea.View {
	var view string

	if !m.isOpen {
		goto end
	}

	view = m.renderModal()

end:
	return tea.NewView(view)
}

// Open opens the modal, initializes state, and calculates position
func (m MultiSelectModel[T]) Open() (MultiSelectModel[T], tea.Cmd) {
	m.isOpen = true
	m.cursor = 0
	m.offset = 0
	m.focus = focusList
	m.buttonIdx = 0
	m = m.adjustOffset()

	// Pre-calculate modal dimensions and position
	modalView := m.renderModal()
	m.width, m.height, m.lastRow, m.lastCol = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)

	return m, nil
}

// Close closes the modal
func (m MultiSelectModel[T]) Close() (MultiSelectModel[T], tea.Cmd) {
	m.isOpen = false
	return m, nil
}

// IsOpen returns whether the modal is currently open
func (m MultiSelectModel[T]) IsOpen() bool {
	return m.isOpen
}

// SetSize sets screen dimensions
func (m MultiSelectModel[T]) SetSize(width, height int) MultiSelectModel[T] {
	m.screenWidth = width
	m.screenHeight = height
	return m
}

// SetItems updates the list items and resets selection state
func (m MultiSelectModel[T]) SetItems(items []T) MultiSelectModel[T] {
	m.items = items
	// Reset checked state — all checked by default
	m.checked = make(map[string]bool, len(items))
	for _, item := range items {
		m.checked[item.ID()] = true
	}
	// Adjust cursor if out of bounds
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	m = m.adjustOffset()
	return m
}

// Selected returns the currently checked items
func (m MultiSelectModel[T]) Selected() []T {
	return m.selectedItems()
}

// Cursor returns the current cursor position
func (m MultiSelectModel[T]) Cursor() int {
	return m.cursor
}

// OverlayModal renders the modal centered over the background view
func (m MultiSelectModel[T]) OverlayModal(background string) (view string) {
	var modalView string
	var row, col int

	if !m.isOpen {
		view = background
		goto end
	}

	modalView = m.View().Content
	row = m.lastRow
	col = m.lastCol

	view = OverlayModal(background, modalView, row, col)

end:
	return view
}

// =============================================================================
// Navigation helpers
// =============================================================================

// moveCursorUp moves cursor up and adjusts offset if needed
func (m MultiSelectModel[T]) moveCursorUp() MultiSelectModel[T] {
	if m.cursor > 0 {
		m.cursor--
	}
	m = m.adjustOffset()
	return m
}

// moveCursorDown moves cursor down and adjusts offset if needed
func (m MultiSelectModel[T]) moveCursorDown() MultiSelectModel[T] {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
	m = m.adjustOffset()
	return m
}

// adjustOffset ensures cursor is visible in viewport
func (m MultiSelectModel[T]) adjustOffset() MultiSelectModel[T] {
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	if m.cursor >= m.offset+m.maxVisible {
		m.offset = m.cursor - m.maxVisible + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
	return m
}

// visibleItemCount returns how many items are visible in the viewport
func (m MultiSelectModel[T]) visibleItemCount() int {
	count := len(m.items)
	if count > m.maxVisible {
		count = m.maxVisible
	}
	return count
}

// needsScrollbar returns true if there are more items than can be displayed
func (m MultiSelectModel[T]) needsScrollbar() bool {
	return len(m.items) > m.maxVisible
}

// cycleFocusForward cycles focus: list → button0 → button1 → ... → list
func (m MultiSelectModel[T]) cycleFocusForward() MultiSelectModel[T] {
	if m.focus == focusList {
		if len(m.buttons) > 0 {
			m.focus = focusButton
			m.buttonIdx = 0
		}
		return m
	}
	// Currently on a button
	if m.buttonIdx < len(m.buttons)-1 {
		m.buttonIdx++
	} else {
		m.focus = focusList
	}
	return m
}

// cycleFocusBackward reverses the focus cycle
func (m MultiSelectModel[T]) cycleFocusBackward() MultiSelectModel[T] {
	if m.focus == focusList {
		if len(m.buttons) > 0 {
			m.focus = focusButton
			m.buttonIdx = len(m.buttons) - 1
		}
		return m
	}
	// Currently on a button
	if m.buttonIdx > 0 {
		m.buttonIdx--
	} else {
		m.focus = focusList
	}
	return m
}

// selectedItems returns items that are currently checked
func (m MultiSelectModel[T]) selectedItems() []T {
	selected := make([]T, 0, len(m.items))
	for _, item := range m.items {
		if m.checked[item.ID()] {
			selected = append(selected, item)
		}
	}
	return selected
}

// =============================================================================
// Rendering
// =============================================================================

// renderModal creates the modal box view (FOLLOWS ClearPath)
func (m MultiSelectModel[T]) renderModal() (view string) {
	var content strings.Builder
	var titleLine string
	var messageLine string
	var itemsContent string
	var footerLine string
	var buttonLine string
	var maxWidth int
	var titleWidth int
	var messageWidth int
	var footerWidth int
	var itemsWidth int
	var buttonWidth int

	// Calculate content widths
	titleWidth = ansi.StringWidth(m.title)
	messageWidth = m.calculateMessageWidth()
	footerWidth = ansi.StringWidth(m.footer)

	// Render buttons first to measure width
	buttonLine = m.renderButtons()
	buttonWidth = ansi.StringWidth(buttonLine)

	// Render items to measure width
	itemsContent = m.renderItems(0)
	for _, line := range strings.Split(itemsContent, "\n") {
		lineWidth := ansi.StringWidth(line)
		if lineWidth > itemsWidth {
			itemsWidth = lineWidth
		}
	}

	// Determine max width (add padding: 4 for left/right margins)
	maxWidth = titleWidth
	if messageWidth > maxWidth {
		maxWidth = messageWidth
	}
	if footerWidth > maxWidth {
		maxWidth = footerWidth
	}
	if itemsWidth > maxWidth {
		maxWidth = itemsWidth
	}
	if buttonWidth > maxWidth {
		maxWidth = buttonWidth
	}
	maxWidth = maxWidth + 4

	// Ensure minimum width
	if maxWidth < 40 {
		maxWidth = 40
	}

	// Render title
	if m.title != "" {
		titleLine = teautils.RenderCenteredLine(m.title, m.titleStyle, maxWidth)
		content.WriteString(titleLine)
		content.WriteString("\n\n")
	}

	// Render message
	if m.message != "" {
		messageLine = teautils.RenderCenteredLine(m.message, m.messageStyle, maxWidth)
		content.WriteString(messageLine)
		content.WriteString("\n\n")
	}

	// Re-render items with proper left padding for centering within maxWidth
	leftPad := 0
	if itemsWidth < maxWidth-4 {
		leftPad = (maxWidth - 4 - itemsWidth) / 2
	}
	itemsContent = m.renderItems(leftPad)
	content.WriteString(itemsContent)
	content.WriteString("\n")

	// Render footer
	if m.footer != "" {
		content.WriteString("\n")
		footerLine = teautils.RenderCenteredLine(m.footer, m.footerStyle, maxWidth)
		content.WriteString(footerLine)
	}

	// Render buttons (centered)
	content.WriteString("\n\n")
	centeredButtons := teautils.RenderCenteredLine(buttonLine, lipgloss.NewStyle(), maxWidth)
	content.WriteString(centeredButtons)

	// Apply border
	view = teautils.ApplyBoxBorder(m.borderStyle, content.String())

	return view
}

// calculateMessageWidth returns the width needed for the message
func (m MultiSelectModel[T]) calculateMessageWidth() (width int) {
	lines := strings.Split(m.message, "\n")
	for _, line := range lines {
		lineWidth := ansi.StringWidth(line)
		if lineWidth > width {
			width = lineWidth
		}
	}
	return width
}

// renderItems renders the visible items with specified left padding
func (m MultiSelectModel[T]) renderItems(leftPad int) (view string) {
	// lineData holds pre-computed parts of each line for two-phase rendering
	type lineData struct {
		prefix     string // cursor indicator
		checkbox   string // [✓] or [ ]
		labelText  string // raw label text (unstyled, for re-rendering)
		styledText string // styled label
	}

	var lines []string
	var endIdx int
	var showScrollbar bool
	var scrollbarPos int
	var scrollbarHeight int
	var padding string
	var maxLineWidth int
	var lineDataList []lineData

	// Build padding string
	padding = strings.Repeat(" ", leftPad)

	if len(m.items) == 0 {
		view = padding + m.itemStyle.Render("(no items)")
		goto end
	}

	endIdx = m.offset + m.visibleItemCount()
	if endIdx > len(m.items) {
		endIdx = len(m.items)
	}

	showScrollbar = m.needsScrollbar()

	// Calculate scrollbar position and height
	if showScrollbar {
		scrollbarPos, scrollbarHeight = calcScrollbar(m.offset, m.maxVisible, len(m.items))
	}

	// PHASE 1: Build all lines WITHOUT padding, measure max width
	for i := m.offset; i < endIdx; i++ {
		item := m.items[i]
		var ld lineData

		// Cursor prefix
		if m.focus == focusList && i == m.cursor {
			ld.prefix = "▶ "
		} else {
			ld.prefix = "  "
		}

		// Checkbox
		if m.checked[item.ID()] {
			ld.checkbox = m.checkedStyle.Render("[✓]") + " "
		} else {
			ld.checkbox = m.uncheckedStyle.Render("[ ]") + " "
		}

		// Label
		ld.labelText = item.Label()
		if m.focus == focusList && i == m.cursor {
			ld.styledText = m.selectedItemStyle.Render(ld.labelText)
		} else {
			ld.styledText = m.itemStyle.Render(ld.labelText)
		}

		// Measure this line's width
		lineWidth := leftPad + ansi.StringWidth(ld.prefix) + ansi.StringWidth(ld.checkbox) + ansi.StringWidth(ld.styledText)
		if lineWidth > maxLineWidth {
			maxLineWidth = lineWidth
		}

		lineDataList = append(lineDataList, ld)
	}

	// PHASE 2: Build final lines, padding styled text area to match max width
	for i, ld := range lineDataList {
		lineIdx := i
		itemIdx := m.offset + i

		// Calculate how much padding this line needs
		currentWidth := leftPad + ansi.StringWidth(ld.prefix) + ansi.StringWidth(ld.checkbox) + ansi.StringWidth(ld.styledText)
		paddingNeeded := maxLineWidth - currentWidth
		if paddingNeeded < 0 {
			paddingNeeded = 0
		}

		// For cursor row, include padding INSIDE the style so background extends
		var line string
		if m.focus == focusList && itemIdx == m.cursor {
			paddedLabel := ld.labelText + strings.Repeat(" ", paddingNeeded)
			styledWithPadding := m.selectedItemStyle.Render(paddedLabel)
			line = padding + ld.prefix + ld.checkbox + styledWithPadding
		} else {
			line = padding + ld.prefix + ld.checkbox + ld.styledText + strings.Repeat(" ", paddingNeeded)
		}

		// Scrollbar character
		if showScrollbar {
			var scrollChar string
			if lineIdx >= scrollbarPos && lineIdx < scrollbarPos+scrollbarHeight {
				scrollChar = " \u2588" // Full block
			} else {
				scrollChar = " \u2502" // Light vertical bar
			}
			line = line + scrollChar
		}

		lines = append(lines, line)
	}

	view = strings.Join(lines, "\n")

end:
	return view
}

// renderButtons renders the horizontal button bar
func (m MultiSelectModel[T]) renderButtons() (line string) {
	parts := make([]string, 0, len(m.buttons))

	for i, btn := range m.buttons {
		isFocused := m.focus == focusButton && i == m.buttonIdx
		parts = append(parts, m.renderButton(btn, isFocused))
	}

	line = strings.Join(parts, "  ")
	return line
}

// renderButton renders a single button with brackets
func (m MultiSelectModel[T]) renderButton(btn MultiSelectButton, isFocused bool) (rendered string) {
	var style lipgloss.Style

	style = m.buttonStyle
	if isFocused {
		style = m.focusedButtonStyle
	}

	rendered = style.Render("[ " + btn.Label + " ]")
	return rendered
}

// =============================================================================
// Style Getters
// =============================================================================

// BorderStyle returns the border style
func (m MultiSelectModel[T]) BorderStyle() lipgloss.Style {
	return m.borderStyle
}

// TitleStyle returns the title style
func (m MultiSelectModel[T]) TitleStyle() lipgloss.Style {
	return m.titleStyle
}

// MessageStyle returns the message style
func (m MultiSelectModel[T]) MessageStyle() lipgloss.Style {
	return m.messageStyle
}

// FooterStyle returns the footer style
func (m MultiSelectModel[T]) FooterStyle() lipgloss.Style {
	return m.footerStyle
}

// ItemStyle returns the item style
func (m MultiSelectModel[T]) ItemStyle() lipgloss.Style {
	return m.itemStyle
}

// SelectedItemStyle returns the selected item style
func (m MultiSelectModel[T]) SelectedItemStyle() lipgloss.Style {
	return m.selectedItemStyle
}

// CheckedStyle returns the checked checkbox style
func (m MultiSelectModel[T]) CheckedStyle() lipgloss.Style {
	return m.checkedStyle
}

// UncheckedStyle returns the unchecked checkbox style
func (m MultiSelectModel[T]) UncheckedStyle() lipgloss.Style {
	return m.uncheckedStyle
}

// ButtonStyle returns the button style
func (m MultiSelectModel[T]) ButtonStyle() lipgloss.Style {
	return m.buttonStyle
}

// FocusedButtonStyle returns the focused button style
func (m MultiSelectModel[T]) FocusedButtonStyle() lipgloss.Style {
	return m.focusedButtonStyle
}

// =============================================================================
// Style Withers
// =============================================================================

// WithBorderStyle returns a copy with the specified border style
func (m MultiSelectModel[T]) WithBorderStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.borderStyle = style
	return m
}

// WithTitleStyle returns a copy with the specified title style
func (m MultiSelectModel[T]) WithTitleStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.titleStyle = style
	return m
}

// WithMessageStyle returns a copy with the specified message style
func (m MultiSelectModel[T]) WithMessageStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.messageStyle = style
	return m
}

// WithFooterStyle returns a copy with the specified footer style
func (m MultiSelectModel[T]) WithFooterStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.footerStyle = style
	return m
}

// WithItemStyle returns a copy with the specified item style
func (m MultiSelectModel[T]) WithItemStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.itemStyle = style
	return m
}

// WithSelectedItemStyle returns a copy with the specified selected item style
func (m MultiSelectModel[T]) WithSelectedItemStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.selectedItemStyle = style
	return m
}

// WithCheckedStyle returns a copy with the specified checked checkbox style
func (m MultiSelectModel[T]) WithCheckedStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.checkedStyle = style
	return m
}

// WithUncheckedStyle returns a copy with the specified unchecked checkbox style
func (m MultiSelectModel[T]) WithUncheckedStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.uncheckedStyle = style
	return m
}

// WithButtonStyle returns a copy with the specified button style
func (m MultiSelectModel[T]) WithButtonStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.buttonStyle = style
	return m
}

// WithFocusedButtonStyle returns a copy with the specified focused button style
func (m MultiSelectModel[T]) WithFocusedButtonStyle(style lipgloss.Style) MultiSelectModel[T] {
	m.focusedButtonStyle = style
	return m
}

// =============================================================================
// Content Withers
// =============================================================================

// WithTitle returns a copy with the specified title
func (m MultiSelectModel[T]) WithTitle(title string) MultiSelectModel[T] {
	m.title = title
	return m
}

// WithMessage returns a copy with the specified message
func (m MultiSelectModel[T]) WithMessage(message string) MultiSelectModel[T] {
	m.message = message
	return m
}

// WithFooter returns a copy with the specified footer
func (m MultiSelectModel[T]) WithFooter(footer string) MultiSelectModel[T] {
	m.footer = footer
	return m
}

// WithMaxVisible returns a copy with the specified max visible items
func (m MultiSelectModel[T]) WithMaxVisible(max int) MultiSelectModel[T] {
	m.maxVisible = max
	return m
}

// WithLogger returns a copy with the specified logger for debugging
func (m MultiSelectModel[T]) WithLogger(logger *slog.Logger) MultiSelectModel[T] {
	m.Logger = logger
	return m
}
