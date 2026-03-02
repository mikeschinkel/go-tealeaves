package teamodal

import (
	"log/slog"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// ListItem is the constraint interface for items displayed in ListModel
type ListItem interface {
	// ID returns a unique identifier for the item
	ID() string
	// Label returns the display text for the item
	Label() string
	// IsActive returns true if this item is the currently active/selected item
	// (distinct from cursor position - this marks the "in use" item)
	IsActive() bool
}

// ListModelArgs provides optional configuration for ListModel
type ListModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	Title        string
	MaxVisible   int // Default: 8
	LabelWidth   int // Fixed width for labels/edit field. 0 = use max item label width

	// Optional style overrides
	BorderStyle       lipgloss.Style
	TitleStyle        lipgloss.Style
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	ActiveItemStyle   lipgloss.Style
	FooterStyle       lipgloss.Style
}

// ListModel is a Bubble Tea model for list modal dialogs with CRUD operations
type ListModel[T ListItem] struct {
	// Public
	Keys   ListKeyMap
	Logger *slog.Logger

	// Content
	title         string
	items         []T
	statusMessage string // Feedback message shown above footer

	// Navigation state
	cursor int // Currently highlighted item index
	offset int // Scroll offset for viewport

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

	// Inline editing state
	isEditing     bool   // Whether we're in edit mode
	editBuffer    string // Current text being edited
	editCursor    int    // Cursor position in edit buffer
	editOriginal  string // Original value (for cancel)
	editOverwrite bool   // True if first char typed should overwrite all

	// Label rendering
	labelWidth int // Fixed width for labels and edit field

	// Help visor
	showHelp bool // Whether help visor is visible

	// Styles
	borderStyle       lipgloss.Style
	titleStyle        lipgloss.Style
	itemStyle         lipgloss.Style
	selectedItemStyle lipgloss.Style
	activeItemStyle   lipgloss.Style
	footerStyle       lipgloss.Style
	statusStyle       lipgloss.Style
	editItemStyle     lipgloss.Style // Style for item being edited
}

// NewListModel creates a new ListModel with the given items
func NewListModel[T ListItem](items []T, args *ListModelArgs) (m ListModel[T]) {
	if args == nil {
		args = &ListModelArgs{}
	}

	maxVisible := args.MaxVisible
	if maxVisible <= 0 {
		maxVisible = 8
	}

	m = ListModel[T]{
		Keys:              DefaultListKeyMap(),
		title:             args.Title,
		items:             items,
		cursor:            0,
		offset:            0,
		isOpen:            false,
		screenWidth:       args.ScreenWidth,
		screenHeight:      args.ScreenHeight,
		maxVisible:        maxVisible,
		borderStyle:       DefaultBorderStyle(),
		titleStyle:        DefaultTitleStyle(),
		itemStyle:         DefaultListItemStyle(),
		selectedItemStyle: DefaultSelectedItemStyle(),
		activeItemStyle:   DefaultActiveItemStyle(),
		footerStyle:       DefaultListFooterStyle(),
		statusStyle:       DefaultStatusStyle(),
		editItemStyle:     DefaultEditItemStyle(),
	}

	// Apply custom styles if provided (check if non-zero)
	if args.BorderStyle.String() != "" {
		m.borderStyle = args.BorderStyle
	}
	if args.TitleStyle.String() != "" {
		m.titleStyle = args.TitleStyle
	}
	if args.ItemStyle.String() != "" {
		m.itemStyle = args.ItemStyle
	}
	if args.SelectedItemStyle.String() != "" {
		m.selectedItemStyle = args.SelectedItemStyle
	}
	if args.ActiveItemStyle.String() != "" {
		m.activeItemStyle = args.ActiveItemStyle
	}
	if args.FooterStyle.String() != "" {
		m.footerStyle = args.FooterStyle
	}

	// Calculate label width: use provided value or calculate from max item label
	m.labelWidth = args.LabelWidth
	if m.labelWidth <= 0 {
		m.labelWidth = m.maxItemLabelWidth()
	}

	return m
}

// Init implements tea.Model - returns nil (no initial command)
func (m ListModel[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
func (m ListModel[T]) Update(msg tea.Msg) (ListModel[T], tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var sizeMsg tea.WindowSizeMsg

	if !m.isOpen {
		goto end // Not open = nil cmd = didn't handle
	}

	// Handle editing mode separately
	if m.isEditing {
		return m.updateEditing(msg)
	}

	// Try as KeyPressMsg first
	keyMsg, ok = msg.(tea.KeyPressMsg)
	if ok {
		switch {
		case key.Matches(keyMsg, m.Keys.Up):
			m = m.moveCursorUp()
			m.statusMessage = "" // Clear status on navigation
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Down):
			m = m.moveCursorDown()
			m.statusMessage = "" // Clear status on navigation
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Preview):
			// Space: preview-select (mark active, stay open)
			if len(m.items) > 0 {
				item := m.items[m.cursor]
				cmd = func() tea.Msg { return ItemSelectedMsg[T]{Item: item} }
			}
			goto end

		case key.Matches(keyMsg, m.Keys.Accept):
			// Enter: select cursor item if not already active, then close
			if len(m.items) > 0 {
				cursorItem := m.items[m.cursor]
				// If cursor item is not active, send selection first
				if !cursorItem.IsActive() {
					// Send both ItemSelectedMsg and ListAcceptedMsg
					m.isOpen = false
					cmd = tea.Batch(
						func() tea.Msg { return ItemSelectedMsg[T]{Item: cursorItem} },
						func() tea.Msg { return ListAcceptedMsg[T]{Item: cursorItem} },
					)
				} else {
					// Already active, just close with accept
					m.isOpen = false
					cmd = func() tea.Msg { return ListAcceptedMsg[T]{Item: cursorItem} }
				}
			} else {
				// No items, just close
				m.isOpen = false
				cmd = func() tea.Msg { return ListAcceptedMsg[T]{} }
			}
			goto end

		case key.Matches(keyMsg, m.Keys.New):
			cmd = func() tea.Msg { return NewItemRequestedMsg{} }
			goto end

		case key.Matches(keyMsg, m.Keys.Edit):
			if len(m.items) > 0 {
				// Enter inline edit mode
				item := m.items[m.cursor]
				m.isEditing = true
				m.editBuffer = item.Label()
				m.editCursor = 0 // Cursor at start
				m.editOriginal = item.Label()
				m.editOverwrite = true // First keystroke overwrites
				cmd = func() tea.Msg { return nil }
			}
			goto end

		case key.Matches(keyMsg, m.Keys.Delete):
			if len(m.items) > 0 {
				item := m.items[m.cursor]
				cmd = func() tea.Msg { return DeleteItemRequestedMsg[T]{Item: item} }
			}
			goto end

		case key.Matches(keyMsg, m.Keys.Help):
			// Toggle help visor
			m.showHelp = !m.showHelp
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Cancel):
			// If help is showing, close help instead of modal
			if m.showHelp {
				m.showHelp = false
				cmd = func() tea.Msg { return nil }
				goto end
			}
			m.isOpen = false
			cmd = func() tea.Msg { return ListCancelledMsg{} }
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

end:
	return m, cmd
}

// updateEditing handles key events when in inline edit mode
func (m ListModel[T]) updateEditing(msg tea.Msg) (ListModel[T], tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var inputRunes []rune

	keyMsg, ok = msg.(tea.KeyPressMsg)
	if !ok {
		goto end
	}

	switch keyMsg.Code {
	case tea.KeyEnter:
		// Complete the edit
		item := m.items[m.cursor]
		newLabel := m.editBuffer
		m.isEditing = false
		m.editBuffer = ""
		m.editCursor = 0
		m.editOriginal = ""
		m.editOverwrite = false
		cmd = func() tea.Msg { return EditCompletedMsg[T]{Item: item, NewLabel: newLabel} }
		goto end

	case tea.KeyEscape:
		// Cancel the edit
		m.isEditing = false
		m.editBuffer = ""
		m.editCursor = 0
		m.editOriginal = ""
		m.editOverwrite = false
		cmd = func() tea.Msg { return nil }
		goto end

	case tea.KeyLeft:
		// Move cursor left, disable overwrite mode
		m.editOverwrite = false
		if m.editCursor > 0 {
			m.editCursor--
		}
		cmd = func() tea.Msg { return nil }
		goto end

	case tea.KeyRight:
		// Move cursor right, disable overwrite mode
		m.editOverwrite = false
		runes := []rune(m.editBuffer)
		if m.editCursor < len(runes) {
			m.editCursor++
		}
		cmd = func() tea.Msg { return nil }
		goto end

	case tea.KeyBackspace:
		// Delete character before cursor
		m.editOverwrite = false
		if m.editCursor > 0 {
			runes := []rune(m.editBuffer)
			m.editBuffer = string(runes[:m.editCursor-1]) + string(runes[m.editCursor:])
			m.editCursor--
		}
		cmd = func() tea.Msg { return nil }
		goto end

	case tea.KeyDelete:
		// Delete character at cursor
		m.editOverwrite = false
		runes := []rune(m.editBuffer)
		if m.editCursor < len(runes) {
			m.editBuffer = string(runes[:m.editCursor]) + string(runes[m.editCursor+1:])
		}
		cmd = func() tea.Msg { return nil }
		goto end
	}

	// Handle text input (runes + space, unified in v2)
	if keyMsg.Text != "" {
		inputRunes = []rune(keyMsg.Text)
		if m.editOverwrite {
			// First keystroke: overwrite entire buffer
			m.editBuffer = keyMsg.Text
			m.editCursor = len(inputRunes)
			m.editOverwrite = false
		} else {
			// Insert at cursor position
			runes := []rune(m.editBuffer)
			newRunes := make([]rune, 0, len(runes)+len(inputRunes))
			newRunes = append(newRunes, runes[:m.editCursor]...)
			newRunes = append(newRunes, inputRunes...)
			newRunes = append(newRunes, runes[m.editCursor:]...)
			m.editBuffer = string(newRunes)
			m.editCursor += len(inputRunes)
		}
		cmd = func() tea.Msg { return nil }
		goto end
	}

end:
	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m ListModel[T]) View() tea.View {
	var view string

	if !m.isOpen {
		goto end
	}

	view = m.renderModal()

end:
	return tea.NewView(view)
}

// Open opens the modal and returns updated model
func (m ListModel[T]) Open() ListModel[T] {
	m.isOpen = true

	// Position cursor on the active item if one exists, otherwise start at 0
	m.cursor = 0
	for i, item := range m.items {
		if item.IsActive() {
			m.cursor = i
			break
		}
	}

	// Ensure cursor is in bounds
	if m.cursor >= len(m.items) {
		m.cursor = 0
	}

	// Adjust offset to ensure cursor is visible
	m.offset = 0
	m = m.adjustOffset()

	// Pre-calculate modal dimensions and position
	modalView := m.renderModal()

	// Use helper to measure and center the modal
	m.width, m.height, m.lastRow, m.lastCol = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)

	return m
}

// Close closes the modal and returns updated model
func (m ListModel[T]) Close() ListModel[T] {
	m.isOpen = false
	return m
}

// IsOpen returns whether the modal is currently open
func (m ListModel[T]) IsOpen() bool {
	return m.isOpen
}

// SetSize sets screen dimensions
func (m ListModel[T]) SetSize(width, height int) ListModel[T] {
	m.screenWidth = width
	m.screenHeight = height
	return m
}

// SetItems updates the list items
func (m ListModel[T]) SetItems(items []T) ListModel[T] {
	m.items = items
	// Adjust cursor if it's out of bounds
	if m.cursor >= len(m.items) {
		m.cursor = len(m.items) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
	// Adjust offset if needed
	m = m.adjustOffset()

	// Recalculate labelWidth based on new items
	// This is critical when items are set after construction (e.g., nil at NewListModel)
	newLabelWidth := m.maxItemLabelWidth()
	if newLabelWidth > m.labelWidth {
		m.labelWidth = newLabelWidth
	}

	return m
}

// Items returns the current list items
func (m ListModel[T]) Items() []T {
	return m.items
}

// SetCursor sets the cursor to the specified index
func (m ListModel[T]) SetCursor(index int) ListModel[T] {
	if index < 0 {
		index = 0
	}
	if index >= len(m.items) {
		index = len(m.items) - 1
	}
	if index < 0 {
		index = 0
	}
	m.cursor = index
	m = m.adjustOffset()
	return m
}

// SetCursorToLast sets the cursor to the last item
func (m ListModel[T]) SetCursorToLast() ListModel[T] {
	return m.SetCursor(len(m.items) - 1)
}

// SetStatus sets a status/feedback message to display above the footer
func (m ListModel[T]) SetStatus(msg string) ListModel[T] {
	m.statusMessage = msg
	return m
}

// ClearStatus clears the status message
func (m ListModel[T]) ClearStatus() ListModel[T] {
	m.statusMessage = ""
	return m
}

// SelectedItem returns the item at the cursor position
func (m ListModel[T]) SelectedItem() (item T) {
	if len(m.items) == 0 {
		goto end
	}
	if m.cursor >= len(m.items) {
		goto end
	}
	item = m.items[m.cursor]

end:
	return item
}

// ActiveItem returns the item where IsActive() returns true
func (m ListModel[T]) ActiveItem() (item T) {
	for _, i := range m.items {
		if i.IsActive() {
			item = i
			goto end
		}
	}

end:
	return item
}

// OverlayModal renders the modal centered over the background view
func (m ListModel[T]) OverlayModal(background string) (view string) {
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
// Internal helpers
// =============================================================================

// moveCursorUp moves cursor up and adjusts offset if needed
func (m ListModel[T]) moveCursorUp() ListModel[T] {
	if m.cursor > 0 {
		m.cursor--
	}
	m = m.adjustOffset()
	return m
}

// moveCursorDown moves cursor down and adjusts offset if needed
func (m ListModel[T]) moveCursorDown() ListModel[T] {
	if m.cursor < len(m.items)-1 {
		m.cursor++
	}
	m = m.adjustOffset()
	return m
}

// adjustOffset ensures cursor is visible in viewport
func (m ListModel[T]) adjustOffset() ListModel[T] {
	// Cursor above viewport
	if m.cursor < m.offset {
		m.offset = m.cursor
	}
	// Cursor below viewport
	if m.cursor >= m.offset+m.maxVisible {
		m.offset = m.cursor - m.maxVisible + 1
	}
	// Ensure offset is not negative
	if m.offset < 0 {
		m.offset = 0
	}
	return m
}

// visibleItemCount returns how many items are visible in the viewport
func (m ListModel[T]) visibleItemCount() int {
	count := len(m.items)
	if count > m.maxVisible {
		count = m.maxVisible
	}
	return count
}

// needsScrollbar returns true if there are more items than can be displayed
func (m ListModel[T]) needsScrollbar() bool {
	return len(m.items) > m.maxVisible
}

// renderModal creates the modal box view (FOLLOWS ClearPath)
func (m ListModel[T]) renderModal() (view string) {
	var content strings.Builder
	var titleLine string
	var itemsContent string
	var footerLine string
	var maxWidth int
	var titleWidth int
	var footerWidth int
	var maxItemWidth int
	var leftPad int

	// DEBUG: Log label widths for each item to diagnose emoji width issues
	if m.Logger != nil {
		for i, item := range m.items {
			label := item.Label()
			m.Logger.Info("DEBUG item width",
				"index", i,
				"label", label,
				"byteLen", len(label),
				"runeLen", len([]rune(label)),
				"ansiWidth", ansi.StringWidth(label),
				"labelWidth", m.labelWidth)
		}
	}

	// Calculate content widths (use ansi.StringWidth for proper emoji/wide char handling)
	titleWidth = ansi.StringWidth(m.title)

	// Calculate max item width (labelWidth + " [ACTIVE]" + "▶ " prefix)
	// The prefix "▶ " is 2 chars, " [ACTIVE]" is 9 chars
	maxItemWidth = m.labelWidth + 2 + 9

	// Add scrollbar width if scrollbar will be shown (space + vertical bar = 2 chars)
	if m.needsScrollbar() {
		maxItemWidth += 2
	}

	// Render footer and get its width
	// Calculate width using the NORMAL footer (not edit footer) to ensure consistent dialog size
	footerLine = m.renderFooter()
	footerWidth = ansi.StringWidth(footerLine)

	// Always calculate based on normal footer width for consistency
	normalFooterWidth := m.normalFooterWidth()
	if normalFooterWidth > footerWidth {
		footerWidth = normalFooterWidth
	}

	// Calculate left padding to balance items with footer
	if footerWidth > maxItemWidth {
		leftPad = (footerWidth - maxItemWidth) / 2
	}

	// Render items with calculated padding
	itemsContent = m.renderItemsWithPadding(leftPad)
	itemsLines := strings.Split(itemsContent, "\n")
	itemsWidth := 0
	for i, line := range itemsLines {
		lineWidth := ansi.StringWidth(line)
		if lineWidth > itemsWidth {
			itemsWidth = lineWidth
		}
		// DEBUG: See actual line widths
		if m.Logger != nil {
			m.Logger.Info("DEBUG line width", "index", i, "ansiWidth", lineWidth)
		}
	}
	if m.Logger != nil {
		m.Logger.Info("DEBUG widths summary",
			"maxItemWidth", maxItemWidth,
			"itemsWidth", itemsWidth,
			"labelWidth", m.labelWidth)
	}

	// Determine max width (add padding: 2 for left/right margins)
	maxWidth = titleWidth
	if itemsWidth > maxWidth {
		maxWidth = itemsWidth
	}
	if footerWidth > maxWidth {
		maxWidth = footerWidth
	}
	maxWidth = maxWidth + 4 // Add padding

	// Ensure minimum width
	if maxWidth < 40 {
		maxWidth = 40
	}

	// Render title (if present)
	// Extra blank line after title provides minimum height for help visor overlay
	if m.title != "" {
		titleLine = teautils.RenderCenteredLine(m.title, m.titleStyle, maxWidth)
		content.WriteString(titleLine)
		content.WriteString("\n\n\n") // Extra blank line for visor headroom
	}

	// Render items
	content.WriteString(itemsContent)

	// Add spacing and status/feedback line (3 lines total: blank, status or blank, blank)
	content.WriteString("\n") // Blank line after items
	if m.statusMessage != "" {
		// Center the status message
		statusLine := teautils.RenderCenteredLine(m.statusMessage, m.statusStyle, maxWidth)
		content.WriteString(statusLine)
	}
	content.WriteString("\n") // Blank line before footer

	// Render footer with key hints (centered)
	footerLine = teautils.RenderCenteredLine(footerLine, m.footerStyle, maxWidth)
	content.WriteString(footerLine)

	// Apply border
	view = teautils.ApplyBoxBorder(m.borderStyle, content.String())

	// Overlay help visor if visible
	if m.showHelp {
		helpVisor := m.renderHelpVisor()
		view = m.overlayHelpVisor(view, helpVisor)
	}

	return view
}

// maxItemLabelWidth returns the display width of the longest item label
// Uses ansi.StringWidth to properly handle emojis and wide characters
func (m ListModel[T]) maxItemLabelWidth() (maxWidth int) {
	for _, item := range m.items {
		w := ansi.StringWidth(item.Label())
		if w > maxWidth {
			maxWidth = w
		}
	}
	return maxWidth
}

// renderItemsWithPadding renders the visible items with specified left padding
func (m ListModel[T]) renderItemsWithPadding(leftPad int) (view string) {
	// lineData holds pre-computed parts of each line for two-phase rendering
	type lineData struct {
		prefix     string
		styledText string
		suffix     string
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
	for i := 0; i < leftPad; i++ {
		padding += " "
	}

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
		scrollbarHeight = m.maxVisible * m.visibleItemCount() / len(m.items)
		if scrollbarHeight < 1 {
			scrollbarHeight = 1
		}
		scrollbarPos = m.offset * m.visibleItemCount() / len(m.items)
	}

	// PHASE 1: Build all lines WITHOUT padding, measure max width
	// This avoids emoji width calculation issues by measuring complete styled output
	for i := m.offset; i < endIdx; i++ {
		item := m.items[i]
		var ld lineData

		// Cursor prefix
		if i == m.cursor {
			ld.prefix = "▶ "
		} else {
			ld.prefix = "  "
		}

		// Active badge suffix
		if item.IsActive() && !(m.isEditing && i == m.cursor) {
			ld.suffix = " " + m.activeItemStyle.Render("[ACTIVE]")
		} else {
			ld.suffix = "         " // 9 spaces to match " [ACTIVE]" width
		}

		// Styled label (no padding yet)
		if i == m.cursor && m.isEditing {
			ld.styledText = m.renderEditField()
		} else if i == m.cursor {
			ld.styledText = m.selectedItemStyle.Render(item.Label())
		} else {
			ld.styledText = m.itemStyle.Render(item.Label())
		}

		// Measure this line's width (leftPad + prefix + styledText + suffix)
		lineWidth := leftPad + ansi.StringWidth(ld.prefix) + ansi.StringWidth(ld.styledText) + ansi.StringWidth(ld.suffix)
		if lineWidth > maxLineWidth {
			maxLineWidth = lineWidth
		}

		lineDataList = append(lineDataList, ld)
	}

	// PHASE 2: Build final lines, padding styledText area to match max width
	for i, ld := range lineDataList {
		lineIdx := i
		itemIdx := m.offset + i
		item := m.items[itemIdx]

		// Calculate how much padding this line needs after the styled text
		currentWidth := leftPad + ansi.StringWidth(ld.prefix) + ansi.StringWidth(ld.styledText) + ansi.StringWidth(ld.suffix)
		paddingNeeded := maxLineWidth - currentWidth
		if paddingNeeded < 0 {
			paddingNeeded = 0
		}

		// For selected item, include padding INSIDE the styled area so background extends
		var line string
		if itemIdx == m.cursor && !m.isEditing {
			// Re-render with padding inside the style
			paddedLabel := item.Label() + strings.Repeat(" ", paddingNeeded)
			styledWithPadding := m.selectedItemStyle.Render(paddedLabel)
			line = padding + ld.prefix + styledWithPadding + ld.suffix
		} else if itemIdx == m.cursor && m.isEditing {
			// Edit field + trailing padding
			line = padding + ld.prefix + ld.styledText + strings.Repeat(" ", paddingNeeded) + ld.suffix
		} else {
			// Non-selected: padding goes after styled text (invisible, that's fine)
			line = padding + ld.prefix + ld.styledText + strings.Repeat(" ", paddingNeeded) + ld.suffix
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

// padLabel pads a label to the fixed labelWidth with trailing spaces
// Uses ansi.StringWidth to properly handle emojis and wide characters
func (m ListModel[T]) padLabel(label string) string {
	labelWidth := ansi.StringWidth(label)
	if labelWidth >= m.labelWidth {
		return label
	}
	// Pad with spaces to reach labelWidth
	padding := strings.Repeat(" ", m.labelWidth-labelWidth)
	return label + padding
}

// renderEditField renders the inline edit field with cursor, padded to labelWidth
func (m ListModel[T]) renderEditField() string {
	runes := []rune(m.editBuffer)

	// Build the display with cursor
	var before, cursorChar, after string

	if m.editCursor < len(runes) {
		before = string(runes[:m.editCursor])
		cursorChar = string(runes[m.editCursor])
		after = string(runes[m.editCursor+1:])
	} else {
		before = m.editBuffer
		cursorChar = " " // Cursor at end shows as space
		after = ""
	}

	// Calculate how much padding we need after the text to reach labelWidth
	// Total content: before + cursorChar + after
	contentLen := len(runes)
	if m.editCursor >= len(runes) {
		contentLen++ // Account for cursor space at end
	}

	paddingNeeded := m.labelWidth - contentLen
	if paddingNeeded < 0 {
		paddingNeeded = 0
	}
	trailingPad := strings.Repeat(" ", paddingNeeded)

	// Cursor position shown with inverted colors
	cursorStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("226")). // Yellow background for cursor
		Foreground(lipgloss.Color("11"))   // Bright yellow text for cursor

	return m.editItemStyle.Render(before) +
		cursorStyle.Render(cursorChar) +
		m.editItemStyle.Render(after+trailingPad)
}

// normalFooterWidth returns the width of the normal (non-edit) footer for consistent sizing
func (m ListModel[T]) normalFooterWidth() int {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("43"))
	footer := keyStyle.Render("[?]") + " Help  " +
		keyStyle.Render("[a]") + " Add  " +
		keyStyle.Render("[e]") + " Edit  " +
		keyStyle.Render("[d]") + " Delete  " +
		keyStyle.Render("[esc]") + " Cancel"
	return ansi.StringWidth(footer)
}

// renderFooter renders the key hints footer with colorized keys
func (m ListModel[T]) renderFooter() string {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("43"))

	if m.isEditing {
		// Edit mode footer (users can figure out [esc] cancels)
		return keyStyle.Render("[←|→]") + " move  " +
			keyStyle.Render("[⌫]") + " bksp  " +
			keyStyle.Render("[del]") + " del  " +
			keyStyle.Render("[enter]") + " accept"
	}

	// Shortened footer - full keys available in help visor
	return keyStyle.Render("[?]") + " Help  " +
		keyStyle.Render("[a]") + " Add  " +
		keyStyle.Render("[e]") + " Edit  " +
		keyStyle.Render("[d]") + " Delete  " +
		keyStyle.Render("[esc]") + " Cancel"
}

// =============================================================================
// Help Visor
// =============================================================================

// buildHelpKeys returns KeyMeta entries for keys NOT shown in footer.
// Footer already shows: [?] Help [a] Add [e] Edit [d] Delete [esc] Cancel
// Help visor shows only the "hidden" keys: navigation and selection.
func (m ListModel[T]) buildHelpKeys() map[string][]teautils.KeyMeta {
	result := make(map[string][]teautils.KeyMeta)

	// Navigation keys (not in footer)
	result["Navigation"] = []teautils.KeyMeta{
		{Binding: m.Keys.Up, HelpModal: true},
		{Binding: m.Keys.Down, HelpModal: true},
	}

	// Selection keys (not in footer) - Preview uses space key, needs DisplayKeys override
	result["Selection"] = []teautils.KeyMeta{
		{Binding: m.Keys.Preview, HelpModal: true, DisplayKeys: []string{"space"}},
		{Binding: m.Keys.Accept, HelpModal: true},
	}

	return result
}

// renderHelpVisor renders the help visor with 3-edge border for overlay
func (m ListModel[T]) renderHelpVisor() string {
	keysByCategory := m.buildHelpKeys()
	categoryOrder := []string{"Navigation", "Selection"}

	// First pass: calculate max key width for alignment
	maxKeyWidth := 0
	for _, category := range categoryOrder {
		keys, ok := keysByCategory[category]
		if !ok {
			continue
		}
		for _, k := range keys {
			keyDisplay := teautils.FormatKeyDisplay(k)
			w := len(keyDisplay)
			if w > maxKeyWidth {
				maxKeyWidth = w
			}
		}
	}

	// Calculate visor content width for centering title
	// First, calculate max line width from keys + descriptions
	maxLineWidth := 0
	keyIndent := "   " // 3 spaces
	for _, category := range categoryOrder {
		keys, ok := keysByCategory[category]
		if !ok {
			continue
		}
		// Category line width
		if len(category) > maxLineWidth {
			maxLineWidth = len(category)
		}
		for _, k := range keys {
			desc := k.HelpText
			if desc == "" {
				desc = k.Binding.Help().Desc
			}
			// Line: indent + key (padded to max) + 2 spaces + desc
			lineWidth := len(keyIndent) + maxKeyWidth + 2 + len(desc)
			if lineWidth > maxLineWidth {
				maxLineWidth = lineWidth
			}
		}
	}

	var content strings.Builder

	// Title - centered within visor, bright green
	titleText := "Keyboard Shortcuts"
	titlePadding := (maxLineWidth - len(titleText)) / 2
	if titlePadding < 0 {
		titlePadding = 0
	}
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("46")) // Bright green
	centeredTitle := strings.Repeat(" ", titlePadding) + titleStyle.Render(titleText)
	content.WriteString(centeredTitle)
	content.WriteString("\n\n") // Blank line after title

	// Category and key styles
	categoryStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("43"))
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("43")).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// Render categories (keyIndent already defined above)
	firstCategory := true
	for _, category := range categoryOrder {
		keys, ok := keysByCategory[category]
		if !ok || len(keys) == 0 {
			continue
		}

		// Blank line between categories (not before first)
		if !firstCategory {
			content.WriteString("\n")
		}
		firstCategory = false

		content.WriteString(categoryStyle.Render(category))
		content.WriteString("\n")

		for _, k := range keys {
			keyDisplay := teautils.FormatKeyDisplay(k)
			desc := k.HelpText
			if desc == "" {
				desc = k.Binding.Help().Desc
			}

			// Left-align key with padding to max width
			paddedKey := keyDisplay + strings.Repeat(" ", maxKeyWidth-len(keyDisplay))

			keyPart := keyStyle.Render(paddedKey)
			descPart := descStyle.Render(desc)

			// 3-space indent, key, 2 spaces, description
			content.WriteString(keyIndent)
			content.WriteString(keyPart)
			content.WriteString("  ")
			content.WriteString(descPart)
			content.WriteString("\n")
		}
	}

	// Apply 3-edge border (top, left, right - no bottom)
	// Using custom border with bottom edges as continuation of sides
	threeEdgeBorder := lipgloss.Border{
		Top:         "─",
		Bottom:      "",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "",
		BottomRight: "",
	}

	borderStyle := lipgloss.NewStyle().
		Border(threeEdgeBorder).
		BorderForeground(lipgloss.Color("99")). // Purple to contrast with dialog
		BorderBottom(false).
		PaddingTop(0).
		PaddingBottom(0).
		PaddingLeft(2).
		PaddingRight(2)

	return borderStyle.Render(content.String())
}

// overlayHelpVisor composites the help visor on top of the dialog
func (m ListModel[T]) overlayHelpVisor(dialogView, helpVisor string) string {
	dialogLines := strings.Split(dialogView, "\n")
	visorLines := strings.Split(helpVisor, "\n")

	dialogHeight := len(dialogLines)
	visorHeight := len(visorLines)

	// Calculate dialog and visor widths
	dialogWidth := 0
	for _, line := range dialogLines {
		w := ansi.StringWidth(line)
		if w > dialogWidth {
			dialogWidth = w
		}
	}
	visorWidth := 0
	for _, line := range visorLines {
		w := ansi.StringWidth(line)
		if w > visorWidth {
			visorWidth = w
		}
	}

	// Position visor: bottom should be above dialog footer
	// This leaves footer row visible below the visor
	// Dialog structure: border(1) + padding(1) + title + blank + content... + blank + footer + padding(1) + border(1)
	// We want visor bottom to end 5 rows before dialog end to show footer below
	bottomMargin := 5 // Rows from bottom of dialog to visor bottom
	startRow := dialogHeight - visorHeight - bottomMargin

	// Ensure visor doesn't overlap top border (minimum row 1, right after border)
	if startRow < 1 {
		startRow = 1
	}

	// Position visor: center - 2 to align left border with footer's "[?]"
	startCol := (dialogWidth-visorWidth)/2 - 2
	if startCol < 1 {
		startCol = 1
	}

	return OverlayModal(dialogView, helpVisor, startRow, startCol)
}

// =============================================================================
// Style Getters
// =============================================================================

// BorderStyle returns the border style
func (m ListModel[T]) BorderStyle() lipgloss.Style {
	return m.borderStyle
}

// TitleStyle returns the title style
func (m ListModel[T]) TitleStyle() lipgloss.Style {
	return m.titleStyle
}

// ItemStyle returns the item style
func (m ListModel[T]) ItemStyle() lipgloss.Style {
	return m.itemStyle
}

// SelectedItemStyle returns the selected item style
func (m ListModel[T]) SelectedItemStyle() lipgloss.Style {
	return m.selectedItemStyle
}

// ActiveItemStyle returns the active item style
func (m ListModel[T]) ActiveItemStyle() lipgloss.Style {
	return m.activeItemStyle
}

// FooterStyle returns the footer style
func (m ListModel[T]) FooterStyle() lipgloss.Style {
	return m.footerStyle
}

// =============================================================================
// Style Withers
// =============================================================================

// WithBorderStyle returns a copy with the specified border style
func (m ListModel[T]) WithBorderStyle(style lipgloss.Style) ListModel[T] {
	m.borderStyle = style
	return m
}

// WithTitleStyle returns a copy with the specified title style
func (m ListModel[T]) WithTitleStyle(style lipgloss.Style) ListModel[T] {
	m.titleStyle = style
	return m
}

// WithItemStyle returns a copy with the specified item style
func (m ListModel[T]) WithItemStyle(style lipgloss.Style) ListModel[T] {
	m.itemStyle = style
	return m
}

// WithSelectedItemStyle returns a copy with the specified selected item style
func (m ListModel[T]) WithSelectedItemStyle(style lipgloss.Style) ListModel[T] {
	m.selectedItemStyle = style
	return m
}

// WithActiveItemStyle returns a copy with the specified active item style
func (m ListModel[T]) WithActiveItemStyle(style lipgloss.Style) ListModel[T] {
	m.activeItemStyle = style
	return m
}

// WithFooterStyle returns a copy with the specified footer style
func (m ListModel[T]) WithFooterStyle(style lipgloss.Style) ListModel[T] {
	m.footerStyle = style
	return m
}

// WithLogger returns a copy with the specified logger for debugging
func (m ListModel[T]) WithLogger(logger *slog.Logger) ListModel[T] {
	m.Logger = logger
	return m
}

// =============================================================================
// Content Withers
// =============================================================================

// WithTitle returns a copy with the specified title
func (m ListModel[T]) WithTitle(title string) ListModel[T] {
	m.title = title
	return m
}

// WithMaxVisible returns a copy with the specified max visible items
func (m ListModel[T]) WithMaxVisible(max int) ListModel[T] {
	m.maxVisible = max
	return m
}

// WithLabelWidth returns a copy with the specified label width for consistent sizing
func (m ListModel[T]) WithLabelWidth(width int) ListModel[T] {
	m.labelWidth = width
	return m
}

// LabelWidth returns the current label width
func (m ListModel[T]) LabelWidth() int {
	return m.labelWidth
}

// Title returns the modal title
func (m ListModel[T]) Title() string {
	return m.title
}

// Cursor returns the current cursor position
func (m ListModel[T]) Cursor() int {
	return m.cursor
}

// Offset returns the current scroll offset
func (m ListModel[T]) Offset() int {
	return m.offset
}
