package teadd

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ModelArgs contains initialization arguments for DropdownModel
type ModelArgs struct {
	ScreenWidth  int
	ScreenHeight int
	TopMargin    int // Don't position dropdown above this row (e.g., 1 to avoid menu bar)
	BottomMargin int // Don't position dropdown below screenHeight - bottomMargin (e.g., 1 to avoid status bar)

	// Styling (optional - defaults will be used if not provided)
	BorderStyle lipgloss.Style
	ItemStyle   lipgloss.Style

	SelectedStyle lipgloss.Style
}

// DropdownModel is a Bubble Tea model for popup selection
type DropdownModel struct {
	Keys DropdownKeyMap // Keyboard bindings

	// Position in parent view
	Row int
	Col int

	// Field position (reference point for dropdown positioning)
	FieldRow int
	FieldCol int

	// Options and selection
	Options      []Option
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

// NewModel creates a new DropdownModel
// items: dropdown items to display
// fieldRow, fieldCol: field position (reference point for dropdown positioning)
// args: configuration arguments (screen size, margins, styling)
func NewModel(options []Option, fieldRow, fieldCol int, args *ModelArgs) (m DropdownModel) {
	if args == nil {
		args = &ModelArgs{}
	}
	m = DropdownModel{
		Keys:          DefaultDropdownKeyMap(),
		FieldRow:      fieldRow,
		FieldCol:      fieldCol,
		Row:           fieldRow + 1, // Initial position: below field
		Col:           fieldCol,
		Options:       options,
		Selected:      0,
		IsOpen:        false,
		ScreenWidth:   args.ScreenWidth,
		ScreenHeight:  args.ScreenHeight,
		TopMargin:     args.TopMargin,
		BottomMargin:  args.BottomMargin,
		BorderStyle:   DefaultBorderStyle(),
		ItemStyle:     DefaultItemStyle(),
		SelectedStyle: DefaultSelectedStyle(),
	}

	// Apply custom styles if provided (check if non-zero)
	// Note: We use a simple non-zero check since lipgloss.Style doesn't expose a reliable IsZero() method
	borderStyleStr := args.BorderStyle.String()
	if borderStyleStr != "" {
		m.BorderStyle = args.BorderStyle
	}
	itemStyleStr := args.ItemStyle.String()
	if itemStyleStr != "" {
		m.ItemStyle = args.ItemStyle
	}
	selectedStyleStr := args.SelectedStyle.String()
	if selectedStyleStr != "" {
		m.SelectedStyle = args.SelectedStyle
	}

	return m
}

// Init implements tea.Model - returns nil (no initial command)
func (m DropdownModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
func (m DropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var selected OptionSelectedMsg
	var sizeMsg tea.WindowSizeMsg

	if !m.IsOpen {
		goto end // Not open = nil cmd = didn't handle
	}

	// Try as KeyPressMsg first
	keyMsg, ok = msg.(tea.KeyPressMsg)
	if ok {
		switch {
		case key.Matches(keyMsg, m.Keys.Up):
			m.Selected--
			if m.Selected < 0 {
				m.Selected = 0
			}
			// Adjust scroll offset if needed (only when selected goes above visible area)
			if m.Selected < m.ScrollOffset {
				m.ScrollOffset = m.Selected
			}
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Down):
			m.Selected++
			if m.Selected >= len(m.Options) {
				m.Selected = len(m.Options) - 1
			}
			// Adjust scroll offset if needed (calculated later based on visible count)
			cmd = func() tea.Msg { return nil }
			goto end

		case key.Matches(keyMsg, m.Keys.Select):
			selected = OptionSelectedMsg{
				Index: m.Selected,
				Text:  m.Options[m.Selected].Text,
				Value: m.Options[m.Selected].Value,
			}
			m.IsOpen = false
			cmd = func() tea.Msg { return selected }
			goto end

		case key.Matches(keyMsg, m.Keys.Cancel):
			m.IsOpen = false
			cmd = func() tea.Msg { return DropdownCancelledMsg{} }
			goto end
		}
	}

	// Try as WindowSizeMsg
	sizeMsg, ok = msg.(tea.WindowSizeMsg)
	if ok {
		m.ScreenWidth = sizeMsg.Width
		m.ScreenHeight = sizeMsg.Height
		cmd = func() tea.Msg { return nil }
		goto end
	}

end:
	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m DropdownModel) View() tea.View {
	var view string
	var pos popupPosition
	var err error

	if !m.IsOpen {
		goto end
	}

	pos, err = m.calculatePosition()
	if err != nil {
		view = "Error: " + err.Error()
		goto end
	}

	view, err = m.renderDropdown(pos)
	if err != nil {
		view = "Error: " + err.Error()
		goto end
	}

end:
	return tea.NewView(view)
}

// Open opens the dropdown and calculates adjusted position
func (m DropdownModel) Open() (DropdownModel, tea.Cmd) {
	var pos popupPosition
	var err error

	// Don't open if screen size not yet set (wait for tea.WindowSizeMsg)
	if m.ScreenWidth <= 0 || m.ScreenHeight <= 0 {
		return m, nil
	}

	// Calculate adjusted position based on screen boundaries
	// This handles: shift left if exceeds right edge, place above if exceeds bottom
	pos, err = m.calculatePosition()
	if err != nil {
		goto end
	}

	// Update Row/Col to the calculated position
	m.Row = pos.y
	m.Col = pos.x
	m.DisplayAbove = pos.displayAbove

end:
	m.IsOpen = true
	return m, nil
}

// Close closes the dropdown
func (m DropdownModel) Close() (DropdownModel, tea.Cmd) {
	m.IsOpen = false
	return m, nil
}

// WithPosition sets the field position (dropdown position is recalculated on Open)
func (m DropdownModel) WithPosition(fieldRow, fieldCol int) DropdownModel {
	m.FieldRow = fieldRow
	m.FieldCol = fieldCol
	m.Row = fieldRow + 1 // Initial position below field
	m.Col = fieldCol
	return m
}

// WithOptions sets the dropdown items
func (m DropdownModel) WithOptions(items []Option) DropdownModel {
	m.Options = items
	if m.Selected >= len(items) {
		m.Selected = len(items) - 1
	}
	if m.Selected < 0 {
		m.Selected = 0
	}
	// Reset scroll offset when items change
	m.ScrollOffset = 0
	return m
}

// WithScreenSize sets screen dimensions
func (m DropdownModel) WithScreenSize(width, height int) DropdownModel {
	m.ScreenWidth = width
	m.ScreenHeight = height
	return m
}

// WithTopMargin sets the top margin (minimum row for dropdown positioning)
func (m DropdownModel) WithTopMargin(margin int) DropdownModel {
	m.TopMargin = margin
	return m
}

// WithBottomMargin sets the bottom margin (dropdown won't extend below screenHeight - margin)
func (m DropdownModel) WithBottomMargin(margin int) DropdownModel {
	m.BottomMargin = margin
	return m
}

// renderDropdown creates the popup box view (method on DropdownModel)
func (m DropdownModel) renderDropdown(pos popupPosition) (view string, err error) {
	var lines []string
	var itemText string
	var line string
	var i int
	var maxWidth int
	var itemIdx int
	var content string
	var hasScrolling bool
	var hasOptionsAbove bool
	var hasOptionsBelow bool
	var indicator string
	var fullLine string
	var itemMaxWidth int
	var itemWidth int
	var padding int

	if !m.IsOpen {
		goto end
	}

	// Check if scrolling is active
	hasScrolling = len(m.Options) > pos.visibleCount
	hasOptionsAbove = pos.visibleStart > 0
	hasOptionsBelow = pos.visibleStart+pos.visibleCount < len(m.Options)

	// Calculate interior width (between the two border characters)
	maxWidth = pos.width - 2

	// Render visible items
	for i = 0; i < pos.visibleCount; i++ {
		itemIdx = pos.visibleStart + i
		if itemIdx >= len(m.Options) {
			break
		}

		itemText = m.Options[itemIdx].Text

		if hasScrolling {
			// With scrolling: " " + itemText + padding + " " + indicator
			// Determine indicator for this line (default to space)
			indicator = " "
			if i == 0 && hasOptionsAbove {
				indicator = "▲"
			} else if i == pos.visibleCount-1 && hasOptionsBelow {
				indicator = "▼"
			}

			// Reserve 1 left + 1 space + 1 indicator
			itemMaxWidth = maxWidth - 3

			// Truncate and pad item to fill available space
			itemText = truncateWithEllipsis(itemText, itemMaxWidth)
			itemWidth = len([]rune(itemText))
			padding = itemMaxWidth - itemWidth
			if padding < 0 {
				padding = 0
			}

			// Build line with indicator
			fullLine = " " + itemText + strings.Repeat(" ", padding) + " " + indicator
		} else {
			// No scrolling: " " + itemText + padding + " "
			// Reserve 1 left + 1 right
			itemMaxWidth = maxWidth - 2

			// Truncate and pad item to fill available space
			itemText = truncateWithEllipsis(itemText, itemMaxWidth)
			itemWidth = len([]rune(itemText))
			padding = itemMaxWidth - itemWidth
			if padding < 0 {
				padding = 0
			}

			// Build line without indicator
			fullLine = " " + itemText + strings.Repeat(" ", padding) + " "
		}

		// Apply styling to complete plain text line
		if itemIdx == m.Selected {
			line = m.SelectedStyle.Render(fullLine)
		} else {
			line = m.ItemStyle.Render(fullLine)
		}

		lines = append(lines, line)
	}

	// Join lines
	content = strings.Join(lines, "\n")

	// Apply border
	view = m.BorderStyle.
		Width(pos.width - 2).
		Height(pos.height - 2).
		Render(content)

end:
	return view, err
}

// popupPosition holds calculated position information (unexported)
type popupPosition struct {
	displayAbove bool // true = above field, false = below
	x            int  // Left edge column
	y            int  // Top edge row (of popup box)
	width        int  // Popup width
	height       int  // Popup height (including borders)
	visibleStart int  // First visible item index (for scrolling)
	visibleCount int  // Number of items that fit
}

// calculatePosition computes optimal popup position (unexported, ClearPath style)
func (m DropdownModel) calculatePosition() (pos popupPosition, err error) {
	var maxOptionLen int
	var availableBelow, availableAbove int
	var requiredWidth int
	var useBelow bool

	items := m.Options
	screenWidth := m.ScreenWidth
	screenHeight := m.ScreenHeight

	// Validate inputs
	if screenWidth <= 0 || screenHeight <= 0 {
		err = NewErr(ErrDropdown, ErrInvalidBounds,
			"width", screenWidth,
			"height", screenHeight,
		)
		goto end
	}

	if len(items) == 0 {
		err = NewErr(ErrDropdown, ErrEmptyOptions)
		goto end
	}

	// Calculate required dimensions
	maxOptionLen = maxLength(items)

	// Calculate available space, respecting margins
	// availableBelow: rows available from (fieldRow+1) down to (screenHeight - bottomMargin - 1)
	availableBelow = screenHeight - m.BottomMargin - m.FieldRow
	// availableAbove: from topMargin to (fieldRow - 2) - need 1 row gap before field
	availableAbove = m.FieldRow - m.TopMargin - 1

	// Decide placement: prefer whichever side has more space
	useBelow = availableBelow >= availableAbove

	if useBelow {
		pos.displayAbove = false
		pos.y = m.FieldRow + 1
		// Use as much space as available, up to number of items
		pos.visibleCount = min(availableBelow-2, len(items))
	} else {
		pos.displayAbove = true
		// Use as much space as available, up to number of items (leave 1 row gap)
		pos.visibleCount = min(availableAbove-2, len(items))
		pos.y = m.FieldRow - pos.visibleCount - 3 // -visibleCount for items, -2 for borders, -1 for gap
	}

	// Calculate height with available space
	pos.height = pos.visibleCount + 2 // +2 for borders

	// Calculate width - add extra space for scroll indicator if needed
	requiredWidth = maxOptionLen + 4 // 2 for borders + 2 for padding
	if len(items) > pos.visibleCount {
		requiredWidth = requiredWidth + 2 // +2 for scroll indicator (space + arrow)
	}

	// Horizontal positioning (2 columns left of field)
	pos.x = m.FieldCol - 2
	if pos.x < 0 {
		pos.x = 0
	}

	pos.width = requiredWidth

	// Shift left if extends beyond screen
	if pos.x+pos.width > screenWidth {
		pos.x = screenWidth - pos.width
		if pos.x < 0 {
			pos.x = 0
			pos.width = screenWidth
		}
	}

	// Use stored scroll offset, but ensure selected is visible
	// If selected is below visible area, scroll down
	if m.Selected >= m.ScrollOffset+pos.visibleCount {
		m.ScrollOffset = m.Selected - pos.visibleCount + 1
	}
	// If selected is above visible area, scroll up (already handled in Update)
	// Clamp scroll offset to valid range
	if m.ScrollOffset < 0 {
		m.ScrollOffset = 0
	}
	if m.ScrollOffset+pos.visibleCount > len(items) {
		m.ScrollOffset = len(items) - pos.visibleCount
		if m.ScrollOffset < 0 {
			m.ScrollOffset = 0
		}
	}

	pos.visibleStart = m.ScrollOffset

end:
	return pos, err
}
