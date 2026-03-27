package teaguide

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// GuideModel is the Bubble Tea model for the guide overlay.
type GuideModel struct {
	Keys   GuideKeyMap
	Styles GuideStyles

	isOpen          bool
	data            GuideData
	screenWidth     int
	screenHeight    int
	scrollOffset    int
	blockedExpanded bool
	actionMap       map[string]bool // keys that trigger ActionSelectedMsg
	lastRow         int             // cached overlay position
	lastCol         int
}

// NewGuideModel creates a new guide model with default key bindings and styles.
func NewGuideModel() GuideModel {
	return GuideModel{
		Keys:   DefaultGuideKeyMap(),
		Styles: DefaultGuideStyles(),
	}
}

// Init implements tea.Model.
func (m GuideModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m GuideModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var sizeMsg tea.WindowSizeMsg
	var ok bool

	if !m.isOpen {
		goto end
	}

	keyMsg, ok = msg.(tea.KeyPressMsg)
	if !ok {
		// Handle window resize
		sizeMsg, ok = msg.(tea.WindowSizeMsg)
		if ok {
			m.screenWidth = sizeMsg.Width
			m.screenHeight = sizeMsg.Height
		}
		goto end
	}

	// Close bindings
	if key.Matches(keyMsg, m.Keys.Close) {
		m.isOpen = false
		cmd = func() tea.Msg { return GuideDismissedMsg{} }
		goto end
	}

	// Scroll bindings
	if key.Matches(keyMsg, m.Keys.ScrollUp) {
		if m.scrollOffset > 0 {
			m.scrollOffset--
		}
		goto end
	}
	if key.Matches(keyMsg, m.Keys.ScrollDown) {
		m.scrollOffset++
		goto end
	}

	// Toggle blocked section
	if key.Matches(keyMsg, m.Keys.ToggleBlock) {
		m.blockedExpanded = !m.blockedExpanded
		goto end
	}

	// Action key dispatch — check against actionMap
	if m.actionMap[keyMsg.String()] {
		m.isOpen = false
		actionKey := keyMsg.String()
		cmd = func() tea.Msg { return ActionSelectedMsg{ActionKey: actionKey} }
		goto end
	}

end:
	return m, cmd
}

// View implements tea.Model.
func (m GuideModel) View() tea.View {
	var view string

	if !m.isOpen {
		goto end
	}

	view = m.renderGuide()

end:
	return tea.NewView(view)
}

// Open opens the guide with the provided data and builds the action key lookup.
func (m GuideModel) Open(data GuideData) (GuideModel, tea.Cmd) {
	m.isOpen = true
	m.data = data
	m.scrollOffset = 0
	m.blockedExpanded = false
	m.actionMap = make(map[string]bool)

	// Build action key map from Recommended and Available sections
	for _, section := range data.Sections {
		if section.Priority == PriorityBlocked {
			continue
		}
		for _, item := range section.Items {
			if item.ActionKey != "" {
				m.actionMap[item.ActionKey] = true
			}
		}
	}

	return m, nil
}

// Close closes the guide.
func (m GuideModel) Close() GuideModel {
	m.isOpen = false
	return m
}

// IsOpen returns whether the guide is currently visible.
func (m GuideModel) IsOpen() bool {
	return m.isOpen
}

// SetSize sets the screen dimensions for overlay positioning.
func (m GuideModel) SetSize(w, h int) GuideModel {
	m.screenWidth = w
	m.screenHeight = h
	return m
}

// WithStyles sets custom styles on the guide model.
func (m GuideModel) WithStyles(styles GuideStyles) GuideModel {
	m.Styles = styles
	return m
}

// OverlayModal composites the guide over a background view.
// Returns the background unchanged if the guide is not open.
func (m GuideModel) OverlayModal(background string) (view string) {
	var modalView string
	var row, col int

	view = background
	if !m.isOpen {
		goto end
	}

	modalView = m.renderGuide()
	_, _, row, col = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)
	m.lastRow = row
	m.lastCol = col

	view = overlayModal(background, modalView, row, col)

end:
	return view
}

// renderGuide builds the full guide modal content.
func (m GuideModel) renderGuide() string {
	var content strings.Builder
	var maxContentWidth int

	// Determine max content width (limit to 60% of screen, min 30, max 50)
	maxContentWidth = m.screenWidth * 60 / 100
	if maxContentWidth < 30 {
		maxContentWidth = 30
	}
	if maxContentWidth > 50 {
		maxContentWidth = 50
	}

	// Title
	content.WriteString(m.Styles.Title.Render(m.data.Title))
	content.WriteString("\n")

	// Sections
	for _, section := range m.data.Sections {
		content.WriteString("\n")

		if section.Priority == PriorityBlocked {
			m.renderBlockedSection(&content, section)
			continue
		}

		// Section heading
		content.WriteString(m.Styles.SectionHeading.Render(section.Heading))
		content.WriteString("\n")

		// Items
		for _, item := range section.Items {
			m.renderItem(&content, item, maxContentWidth)
		}
	}

	// Footer
	content.WriteString("\n")
	content.WriteString(m.Styles.Footer.Render("[Esc] Close  [↑↓] Scroll"))

	// Apply border
	bordered := m.Styles.Border.Render(content.String())

	// Apply scrolling viewport
	bordered = m.applyScroll(bordered)

	return bordered
}

// renderItem renders a single guide item with key, label, and optional prose.
func (m GuideModel) renderItem(b *strings.Builder, item GuideItem, maxWidth int) {
	var line strings.Builder

	if item.KeyDisplay != "" {
		line.WriteString(m.Styles.ItemKey.Render(item.KeyDisplay))
		line.WriteString(" ")
	}
	line.WriteString(m.Styles.ItemLabel.Render(item.Label))
	b.WriteString(line.String())
	b.WriteString("\n")

	if item.Prose != "" {
		// Indent prose and wrap if needed
		prose := item.Prose
		if maxWidth > 6 {
			prose = wrapText(prose, maxWidth-4) // indent allowance
		}
		for _, pLine := range strings.Split(prose, "\n") {
			b.WriteString("  ")
			b.WriteString(m.Styles.ItemProse.Render(pLine))
			b.WriteString("\n")
		}
	}
}

// renderBlockedSection renders the collapsed/expanded blocked section.
func (m GuideModel) renderBlockedSection(b *strings.Builder, section GuideSection) {
	var indicator string
	var heading string

	count := len(section.Items)
	if count == 0 {
		goto end
	}

	// Heading with expand/collapse indicator
	indicator = "▶"
	if m.blockedExpanded {
		indicator = "▼"
	}
	heading = fmt.Sprintf("%s %s (%d)", indicator, section.Heading, count)
	b.WriteString(m.Styles.BlockedHeading.Render(heading))
	b.WriteString("\n")

	if !m.blockedExpanded {
		goto end
	}

	// Expanded items
	for _, item := range section.Items {
		var line strings.Builder
		line.WriteString("  ")
		line.WriteString(m.Styles.BlockedItem.Render(item.Label))
		if item.BlockReason != "" {
			line.WriteString(" — ")
			line.WriteString(m.Styles.BlockReason.Render(item.BlockReason))
		}
		b.WriteString(line.String())
		b.WriteString("\n")
	}

end:
}

// applyScroll clips the rendered content to fit within screen height,
// applying scroll offset.
func (m GuideModel) applyScroll(rendered string) (result string) {
	var lines []string
	var maxVisible int
	var maxOffset int
	var offset int
	var visible []string

	result = rendered
	lines = strings.Split(rendered, "\n")
	maxVisible = m.screenHeight - 2 // leave margin
	if maxVisible <= 0 || len(lines) <= maxVisible {
		goto end
	}

	// Clamp scroll offset
	maxOffset = len(lines) - maxVisible
	offset = m.scrollOffset
	if offset > maxOffset {
		offset = maxOffset
	}

	visible = lines[offset:]
	if len(visible) > maxVisible {
		visible = visible[:maxVisible]
	}

	result = strings.Join(visible, "\n")

end:
	return result
}

// wrapText performs simple word wrapping at the given width.
func wrapText(text string, maxWidth int) (result string) {
	var words []string
	var lines []string
	var current strings.Builder
	var currentWidth int

	result = text
	if maxWidth <= 0 {
		goto end
	}

	words = strings.Fields(text)
	if len(words) == 0 {
		result = ""
		goto end
	}

	for _, word := range words {
		wordWidth := ansi.StringWidth(word)
		if currentWidth > 0 && currentWidth+1+wordWidth > maxWidth {
			lines = append(lines, current.String())
			current.Reset()
			currentWidth = 0
		}
		if currentWidth > 0 {
			current.WriteString(" ")
			currentWidth++
		}
		current.WriteString(word)
		currentWidth += wordWidth
	}
	if current.Len() > 0 {
		lines = append(lines, current.String())
	}

	result = strings.Join(lines, "\n")

end:
	return result
}

// overlayModal composites foreground on background at the given position.
// This follows the proven pattern from teamodal.OverlayModal.
func overlayModal(background, foreground string, row, col int) string {
	var result strings.Builder

	bgLines := strings.Split(background, "\n")
	fgLines := strings.Split(foreground, "\n")

	for i, bgLine := range bgLines {
		fgRow := i - row

		if fgRow < 0 || fgRow >= len(fgLines) {
			result.WriteString(bgLine)
			result.WriteString("\n")
			continue
		}

		fgLine := fgLines[fgRow]
		composited := overlayLine(bgLine, fgLine, col)
		result.WriteString(composited)
		result.WriteString("\n")
	}

	output := result.String()
	if len(output) > 0 && output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}

	return output
}

// overlayLine overlays foreground onto background at column position (ANSI-aware).
func overlayLine(background, foreground string, col int) string {
	if col < 0 {
		col = 0
	}

	bgWidth := ansi.StringWidth(background)
	fgWidth := ansi.StringWidth(foreground)

	var result strings.Builder

	if col > 0 {
		result.WriteString(overlayLeft(background, bgWidth, col))
	}

	result.WriteString(foreground)

	endCol := col + fgWidth
	if endCol < bgWidth {
		remaining := ansi.TruncateLeft(background, endCol, "")
		result.WriteString(remaining)
	}

	return result.String()
}

// overlayLeft returns the left portion of background to place before the overlay.
func overlayLeft(background string, bgWidth, col int) string {
	if col <= bgWidth {
		return ansi.Truncate(background, col, "")
	}
	return background + strings.Repeat(" ", col-bgWidth)
}
