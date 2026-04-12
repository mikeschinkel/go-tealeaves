// Source: site/src/content/docs/reference/overlay-compositing.mdx:15,41,78,116,130,152,159,166
package examples_test

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teafields"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// overlayLine composites a foreground string on top of a background string at a
// given column offset, using ANSI-aware string operations. This is illustrative
// of the internal algorithm used by Tea Leaves overlay functions.
func overlayLine(background, foreground string, col int) string {
	bgWidth := ansi.StringWidth(background)
	fgWidth := ansi.StringWidth(foreground)

	var result strings.Builder

	// Left part: background up to overlay column
	if col > 0 {
		if col <= bgWidth {
			left := ansi.Truncate(background, col, "")
			result.WriteString(left)
		} else {
			result.WriteString(background)
			result.WriteString(strings.Repeat(" ", col-bgWidth))
		}
	}

	// Middle: the overlay content
	result.WriteString(foreground)

	// Right part: background after the overlay ends
	endCol := col + fgWidth
	if endCol < bgWidth {
		remaining := ansi.TruncateLeft(background, endCol, "")
		result.WriteString(remaining)
	}

	return result.String()
}

// overlayDropdownExample illustrates the full multi-line overlay algorithm.
// Named to avoid collision with teafields.OverlayDropdown.
func overlayDropdownExample(background, foreground string, row, col int) string {
	var result strings.Builder

	bgLines := strings.Split(background, "\n")
	fgLines := strings.Split(foreground, "\n")

	for i, bgLine := range bgLines {
		fgRow := i - row

		if fgRow < 0 || fgRow >= len(fgLines) {
			// No overlay on this line
			result.WriteString(bgLine)
		} else {
			// Composite foreground onto background
			composited := overlayLine(bgLine, fgLines[fgRow], col)
			result.WriteString(composited)
		}
		result.WriteString("\n")
	}

	// Remove trailing newline
	output := result.String()
	if len(output) > 0 && output[len(output)-1] == '\n' {
		output = output[:len(output)-1]
	}

	return output
}

// ocMainContent is a minimal stand-in for a main content component.
type ocMainContent struct{}

func (c ocMainContent) View() tea.View { return tea.NewView("main content") }

// ocModalModel is used for the modal-overlay View() example.
type ocModalModel struct {
	modal       teamodal.ConfirmModel
	mainContent ocMainContent
}

func (m ocModalModel) Init() tea.Cmd { return nil }

func (m ocModalModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m ocModalModel) View() tea.View {
	mainView := m.mainContent.View().Content

	if m.modal.IsOpen() {
		return tea.NewView(m.modal.OverlayModal(mainView))
	}

	return tea.NewView(mainView)
}

// ocDropdownModel is used for the dropdown-overlay View() example.
type ocDropdownModel struct {
	dropdown teafields.DropdownModel
}

func (m ocDropdownModel) renderBase() string { return "base view content\nline 2\nline 3" }

func (m ocDropdownModel) Init() tea.Cmd { return nil }

func (m ocDropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m ocDropdownModel) View() tea.View {
	baseView := m.renderBase()

	if m.dropdown.IsOpen {
		dropdownView := m.dropdown.View()
		return tea.NewView(teafields.OverlayDropdown(
			baseView, dropdownView.Content,
			m.dropdown.Row+2, m.dropdown.Col+3,
		))
	}

	return tea.NewView(baseView)
}

// TestCompile_AnsiWidthComparison verifies ANSI vs len width comparison from overlay-compositing.mdx:15.
func TestCompile_AnsiWidthComparison(t *testing.T) {
	styled := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Hello")

	byteLen := len(styled)        // includes escape codes
	visWidth := ansi.StringWidth(styled) // visual width

	_ = byteLen
	_ = visWidth
}

// TestCompile_OverlayLineAlgorithm verifies the per-line overlay algorithm from overlay-compositing.mdx:41.
func TestCompile_OverlayLineAlgorithm(t *testing.T) {
	result := overlayLine("background text here", "OVERLAY", 5)
	_ = result
}

// TestCompile_FullOverlayFunction verifies the full overlay function from overlay-compositing.mdx:78.
func TestCompile_FullOverlayFunction(t *testing.T) {
	bg := "line one\nline two\nline three"
	fg := "FOO\nBAR"
	result := overlayDropdownExample(bg, fg, 1, 2)
	_ = result
}

// TestCompile_ModalOverlayView verifies the modal overlay View() pattern from overlay-compositing.mdx:116.
func TestCompile_ModalOverlayView(t *testing.T) {
	m := ocModalModel{
		modal: teamodal.NewYesNoModal("Continue?", &teamodal.ConfirmModelArgs{
			ScreenWidth:  80,
			ScreenHeight: 24,
		}),
	}
	view := m.View()
	_ = view
}

// TestCompile_DropdownOverlayView verifies the dropdown overlay View() pattern from overlay-compositing.mdx:130.
func TestCompile_DropdownOverlayView(t *testing.T) {
	m := ocDropdownModel{
		dropdown: teafields.NewDropdownModel(
			teafields.ToOptions([]string{"Alpha", "Beta", "Gamma"}),
			&teafields.DropdownModelArgs{},
		),
	}
	view := m.View()
	_ = view
}

// TestCompile_CenterPositioning verifies center-on-screen positioning from overlay-compositing.mdx:152.
func TestCompile_CenterPositioning(t *testing.T) {
	screenHeight := 24
	screenWidth := 80
	overlayHeight := 10
	overlayWidth := 40

	row := (screenHeight - overlayHeight) / 2
	col := (screenWidth - overlayWidth) / 2

	_, _ = row, col
}

// TestCompile_BelowTriggerPositioning verifies below-trigger positioning from overlay-compositing.mdx:159.
func TestCompile_BelowTriggerPositioning(t *testing.T) {
	triggerRow := 5
	triggerCol := 10

	row := triggerRow + 1
	col := triggerCol

	_, _ = row, col
}

// TestCompile_RightAlignedPositioning verifies right-aligned positioning from overlay-compositing.mdx:166.
func TestCompile_RightAlignedPositioning(t *testing.T) {
	screenWidth := 80
	overlayWidth := 30

	col := screenWidth - overlayWidth

	_ = col
}
