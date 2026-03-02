package teamodal

import (
	"strings"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// ProgressModalKeyMap defines the key bindings for the progress modal
type ProgressModalKeyMap struct {
	Cancel     key.Binding
	Background key.Binding
}

// DefaultProgressModalKeyMap returns the default key bindings
func DefaultProgressModalKeyMap() ProgressModalKeyMap {
	return ProgressModalKeyMap{
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Background: key.NewBinding(
			key.WithKeys("b", "B"),
			key.WithHelp("b", "background"),
		),
	}
}

// ProgressModal is a modal for showing an in-progress operation with Cancel and Background options.
// Used for AI-powered operations, long-running tasks, etc.
type ProgressModal struct {
	Keys ProgressModalKeyMap

	// Content
	title string

	// State
	isOpen            bool
	backgroundEnabled bool // Whether to show the Background option
	screenWidth       int
	screenHeight      int

	// Cached dimensions for overlay positioning
	width   int
	height  int
	lastRow int
	lastCol int

	// Styles
	borderStyle lipgloss.Style
	titleStyle  lipgloss.Style
}

// ProgressModalArgs contains initialization arguments for ProgressModal
type ProgressModalArgs struct {
	ScreenWidth       int
	ScreenHeight      int
	Title             string
	BackgroundEnabled bool // Whether to show [b] Background option
}

// NewProgressModal creates a new progress modal.
// The title should describe what is in progress (e.g., "Commit Message").
func NewProgressModal(args *ProgressModalArgs) ProgressModal {
	if args == nil {
		args = &ProgressModalArgs{}
	}

	return ProgressModal{
		Keys:              DefaultProgressModalKeyMap(),
		title:             args.Title,
		backgroundEnabled: args.BackgroundEnabled,
		screenWidth:       args.ScreenWidth,
		screenHeight:      args.ScreenHeight,
		borderStyle:       DefaultBorderStyle(),
		titleStyle:        DefaultTitleStyle(),
	}
}

// Init implements tea.Model
func (m ProgressModal) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m ProgressModal) Update(msg tea.Msg) (ProgressModal, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var sizeMsg tea.WindowSizeMsg

	if !m.isOpen {
		goto end
	}

	keyMsg, ok = msg.(tea.KeyPressMsg)
	if ok {
		switch {
		case key.Matches(keyMsg, m.Keys.Cancel):
			m.isOpen = false
			cmd = func() tea.Msg { return ProgressCancelledMsg{} }
			goto end

		case key.Matches(keyMsg, m.Keys.Background):
			if m.backgroundEnabled {
				m.isOpen = false
				cmd = func() tea.Msg { return ProgressBackgroundMsg{} }
				goto end
			}
		}
		// Ignore all other keys during progress
		goto end
	}

	sizeMsg, ok = msg.(tea.WindowSizeMsg)
	if ok {
		m.screenWidth = sizeMsg.Width
		m.screenHeight = sizeMsg.Height
		cmd = func() tea.Msg { return nil }
	}

end:
	return m, cmd
}

// View renders the modal
func (m ProgressModal) View() tea.View {
	var view string

	if !m.isOpen {
		goto end
	}

	view = m.renderModal()

end:
	return tea.NewView(view)
}

// Open opens the modal and returns updated model
func (m ProgressModal) Open() ProgressModal {
	m.isOpen = true

	// Pre-calculate modal dimensions and position
	modalView := m.renderModal()
	m.width, m.height, m.lastRow, m.lastCol = teautils.CenterModal(modalView, m.screenWidth, m.screenHeight)

	return m
}

// Close closes the modal and returns updated model
func (m ProgressModal) Close() ProgressModal {
	m.isOpen = false
	return m
}

// IsOpen returns whether the modal is currently open
func (m ProgressModal) IsOpen() bool {
	return m.isOpen
}

// SetSize sets screen dimensions
func (m ProgressModal) SetSize(width, height int) ProgressModal {
	m.screenWidth = width
	m.screenHeight = height
	return m
}

// SetTitle sets the modal title
func (m ProgressModal) SetTitle(title string) ProgressModal {
	m.title = title
	return m
}

// SetBackgroundEnabled sets whether the Background option is shown
func (m ProgressModal) SetBackgroundEnabled(enabled bool) ProgressModal {
	m.backgroundEnabled = enabled
	return m
}

// OverlayModal renders the modal centered over the background view
func (m ProgressModal) OverlayModal(background string) (view string) {
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

// renderModal creates the modal box view
func (m ProgressModal) renderModal() string {
	var content strings.Builder
	var titleLine string
	var helpLine string
	var maxWidth int
	var keyStyle lipgloss.Style
	var helpStyle lipgloss.Style

	// Fixed width for progress modal
	maxWidth = 50

	// Title with "Generating..." prefix - centered
	titleLine = teautils.RenderCenteredLine("Generating "+m.title+"...", m.titleStyle, maxWidth)
	content.WriteString(titleLine)
	content.WriteString("\n\n")

	// Help text with highlighted keys
	keyStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")). // Cyan
		Bold(true)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	// Build help line based on enabled options
	if m.backgroundEnabled {
		helpLine = keyStyle.Render("[esc]") + helpStyle.Render(" Cancel  ") +
			keyStyle.Render("[b]") + helpStyle.Render(" Background")
	} else {
		helpLine = helpStyle.Render("Press ") + keyStyle.Render("[esc]") + helpStyle.Render(" to cancel")
	}
	helpLine = teautils.RenderCenteredLine(helpLine, lipgloss.NewStyle(), maxWidth)
	content.WriteString(helpLine)

	// Apply border
	return teautils.ApplyBoxBorder(m.borderStyle, content.String())
}
