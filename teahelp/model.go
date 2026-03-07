package teahelp

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// helpPage holds pre-rendered content lines for a single page (no border, no footer).
type helpPage struct {
	lines []string
}

// helpItem is a single content line with metadata for pagination context.
type helpItem struct {
	line    string // rendered line text
	catLine string // category header line this item belongs to (for context on page 2+)
	isCat   bool   // is this a category header?
	isBlank bool   // is this a blank separator?
}

// HelpVisorModel is a stateful model for the keyboard shortcuts help overlay.
// Supports paging with left/right navigation and colorized footer.
//
// View() returns the rendered visor content; callers handle overlay positioning
// (e.g., via teamodal.OverlayModal or their own logic).
type HelpVisorModel struct {
	Keys   HelpVisorKeyMap
	Styles HelpVisorStyles

	contentStyle    teautils.HelpVisorStyle
	keysByCategory  map[string][]teautils.KeyMeta
	page            int
	totalPages      int
	width           int
	height          int
	maxContentWidth int // widest content line (for centering footer)
	isOpen          bool
	pages           []helpPage
}

// NewHelpVisorModel creates a new, closed HelpVisorModel with default styles and keys.
func NewHelpVisorModel() HelpVisorModel {
	return HelpVisorModel{
		Keys:         DefaultHelpVisorKeyMap(),
		Styles:       DefaultHelpVisorStyles(),
		contentStyle: teautils.DefaultHelpVisorStyle(),
	}
}

// WithContentStyle returns a copy with the given content style (title, category, key, desc).
func (m HelpVisorModel) WithContentStyle(style teautils.HelpVisorStyle) HelpVisorModel {
	m.contentStyle = style
	return m
}

// WithStyles returns a copy with the given chrome styles (border, footer).
func (m HelpVisorModel) WithStyles(styles HelpVisorStyles) HelpVisorModel {
	m.Styles = styles
	return m
}

// WithKeys returns a copy with the given key bindings.
func (m HelpVisorModel) WithKeys(keys HelpVisorKeyMap) HelpVisorModel {
	m.Keys = keys
	return m
}

// Open opens the help visor with the given keys and builds pages.
func (m HelpVisorModel) Open(keysByCategory map[string][]teautils.KeyMeta) HelpVisorModel {
	m.keysByCategory = keysByCategory
	m.isOpen = true
	m.page = 0
	m = m.buildPages()
	return m
}

// Close closes the help visor and clears state.
func (m HelpVisorModel) Close() HelpVisorModel {
	m.isOpen = false
	m.page = 0
	m.pages = nil
	m.keysByCategory = nil
	return m
}

// SetSize updates the available terminal dimensions and rebuilds pages if open.
func (m HelpVisorModel) SetSize(width, height int) HelpVisorModel {
	m.width = width
	m.height = height
	if m.isOpen && m.keysByCategory != nil {
		m = m.buildPages()
	}
	return m
}

// IsOpen returns whether the help visor is currently displayed.
func (m HelpVisorModel) IsOpen() bool {
	return m.isOpen
}

// Page returns the current page index (0-based).
func (m HelpVisorModel) Page() int {
	return m.page
}

// TotalPages returns the total number of pages.
func (m HelpVisorModel) TotalPages() int {
	return m.totalPages
}

// Init implements tea.Model — returns nil (no initial command).
func (m HelpVisorModel) Init() tea.Cmd {
	return nil
}

// Update handles key events when the visor is open.
// Uses key.Matches for configurable bindings. Also supports direct page
// navigation via digit keys 1-9. All keys are consumed (not passed through).
func (m HelpVisorModel) Update(msg tea.Msg) (HelpVisorModel, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool

	if !m.isOpen {
		goto end
	}

	keyMsg, ok = msg.(tea.KeyPressMsg)
	if !ok {
		goto end
	}

	switch {
	case key.Matches(keyMsg, m.Keys.Close):
		m = m.Close()
		cmd = func() tea.Msg { return ClosedMsg{} }
		goto end

	case key.Matches(keyMsg, m.Keys.PrevPage):
		if m.page > 0 {
			m.page--
		}
		goto end

	case key.Matches(keyMsg, m.Keys.NextPage):
		if m.page < m.totalPages-1 {
			m.page++
		}
		goto end
	}

	// Direct page navigation via digit keys
	if len(keyMsg.Text) == 1 {
		ch := keyMsg.Text[0]
		if ch >= '1' && ch <= '9' {
			n := int(ch - '0')
			if n >= 1 && n <= m.totalPages {
				m.page = n - 1
			}
		}
	}

end:
	return m, cmd
}

// View renders the help visor with open-bottom border.
// Returns empty view when closed.
func (m HelpVisorModel) View() tea.View {
	if !m.isOpen || len(m.pages) == 0 {
		return tea.NewView("")
	}

	var pg helpPage
	var inner strings.Builder
	var footer string
	var borderStyle lipgloss.Style

	// Get current page lines
	pg = m.pages[m.page]

	// Build inner content: page lines + optional paging footer
	inner.WriteString(strings.Join(pg.lines, "\n"))

	// Append paging line with blank separator above
	footer = m.renderFooter()
	if footer != "" {
		inner.WriteString("\n\n")
		inner.WriteString(footer)
	}

	// Open-bottom border: no bottom line, PaddingBottom(1) adds a blank row
	// with left/right vertical borders continuing (the "visor-up" look).
	// PaddingLeft(0): content handles its own indentation (1/3/5 spaces).
	// PaddingRight(3): matches the 3-space category indent on the left.
	// Width must account for both padding and border to prevent lipgloss word-wrapping.
	borderStyle = m.Styles.BorderStyle.
		Width(m.maxContentWidth + m.Styles.BorderStyle.GetHorizontalPadding() + m.Styles.BorderStyle.GetHorizontalBorderSize())

	return tea.NewView(borderStyle.Render(inner.String()))
}

// buildPages pre-renders content items and splits them into pages.
// Every page gets the title. Pages that continue mid-category get category context.
//
//goland:noinspection GoAssignmentToReceiver
func (m HelpVisorModel) buildPages() HelpVisorModel {
	var categories []string
	var gap int
	var maxKeyWidth int
	var currentCatLine string

	gap = m.contentStyle.KeyColumnGap
	if gap <= 0 {
		gap = 4
	}

	// Compute max key display width across all categories
	for _, keys := range m.keysByCategory {
		for _, k := range keys {
			w := ansi.StringWidth(teautils.FormatKeyDisplay(k))
			if w > maxKeyWidth {
				maxKeyWidth = w
			}
		}
	}

	keyStyle := m.contentStyle.KeyStyle.Width(maxKeyWidth + gap)

	// Title line (shown on every page)
	titleLine := " " + m.contentStyle.TitleStyle.Render("Keyboard Shortcuts")

	// Build content items (categories + keys, without title)
	var items []helpItem

	categories = teautils.GetSortedCategories(m.keysByCategory, m.contentStyle.CategoryOrder)

	for _, category := range categories {
		keys := m.keysByCategory[category]
		if len(keys) == 0 {
			continue
		}

		// Blank line before category
		items = append(items, helpItem{isBlank: true, catLine: currentCatLine})

		// Category header (3-space indent)
		catLine := "   " + m.contentStyle.CategoryStyle.Render(category)
		currentCatLine = catLine
		items = append(items, helpItem{line: catLine, catLine: currentCatLine, isCat: true})

		// Key lines (5-space indent)
		for _, k := range keys {
			keyDisplay := teautils.FormatKeyDisplay(k)
			desc := k.HelpText
			if desc == "" {
				desc = k.Binding.Help().Desc
			}
			keyPart := keyStyle.Render(keyDisplay)
			descPart := m.contentStyle.DescStyle.Render(desc)
			items = append(items, helpItem{
				line:    "     " + keyPart + descPart,
				catLine: currentCatLine,
			})
		}
	}

	// Compute max content width for footer centering (title + all items)
	m.maxContentWidth = ansi.StringWidth(titleLine)
	for _, item := range items {
		w := ansi.StringWidth(item.line)
		if w > m.maxContentWidth {
			m.maxContentWidth = w
		}
	}

	// Calculate max lines per page
	// Overhead: 1 top border + 1 bottom padding + 1 status bar = 3
	// Multi-page adds: 1 blank + 1 footer = 2 more -> 5 total
	capacity := m.height - 5
	if capacity < 5 {
		capacity = 5
	}

	// Build pages incrementally
	m.pages = nil
	m.page = 0
	idx := 0
	pageNum := 0

	for idx < len(items) {
		var pageLines []string

		// Every page starts with title
		pageLines = append(pageLines, titleLine)

		// Page 2+: add category context if needed
		if pageNum > 0 && idx < len(items) {
			item := items[idx]
			switch {
			case item.isBlank:
				// Natural category break — blank line is already in items
			case item.isCat:
				// Category header — just add a blank separator after title
				pageLines = append(pageLines, "")
			default:
				// Continuing mid-category — add blank + category context
				if item.catLine != "" {
					pageLines = append(pageLines, "")
					pageLines = append(pageLines, item.catLine)
				}
			}
		}

		// Fill page up to capacity
		for idx < len(items) && len(pageLines) < capacity {
			// Category header bump: don't put category as last line on page
			if len(pageLines) == capacity-1 && items[idx].isCat {
				break
			}
			pageLines = append(pageLines, items[idx].line)
			idx++
		}

		m.pages = append(m.pages, helpPage{lines: pageLines})
		pageNum++
	}

	// Ensure at least one page
	if len(m.pages) == 0 {
		m.pages = []helpPage{{lines: []string{titleLine}}}
	}

	m.totalPages = len(m.pages)
	if m.page >= m.totalPages {
		m.page = m.totalPages - 1
	}

	return m
}

// renderFooter builds a centered paging line for multi-page visor.
// Returns empty string for single-page (no footer needed).
// Self-contained: does not depend on teastatus.
func (m HelpVisorModel) renderFooter() string {
	var parts []string
	var pageIndicator string
	var combined string
	var contentWidth int
	var combinedWidth int
	var pad int

	if m.totalPages <= 1 {
		return ""
	}

	// Prev (if not first page)
	if m.page > 0 {
		parts = append(parts,
			m.Styles.FooterKeyStyle.Render("\u2190")+
				" "+
				m.Styles.FooterLabelStyle.Render("Prev"),
		)
	}

	// Next (if not last page)
	if m.page < m.totalPages-1 {
		parts = append(parts,
			m.Styles.FooterKeyStyle.Render("\u2192")+
				" "+
				m.Styles.FooterLabelStyle.Render("Next"),
		)
	}

	// Page indicator
	pageIndicator = m.Styles.FooterLabelStyle.Render(fmt.Sprintf("Page %d/%d", m.page+1, m.totalPages))

	// Build combined string
	if len(parts) > 0 {
		combined = strings.Join(parts, "  ") + "  " + pageIndicator
	} else {
		combined = pageIndicator
	}

	// Center within the actual content width of the visor
	contentWidth = m.maxContentWidth
	if contentWidth < 20 {
		contentWidth = 20
	}
	combinedWidth = ansi.StringWidth(combined)
	pad = (contentWidth - combinedWidth) / 2
	if pad < 0 {
		pad = 0
	}

	return strings.Repeat(" ", pad) + combined
}
