// Source: site/src/content/docs/reference/best-practices.mdx:29,55,74,96,131,151,168,186,215
package examples_test

import (
	"fmt"
	"log/slog"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/viewport"
	lipgloss "charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// bpModel is a minimal model used for best-practices examples that require
// Logger, width fields, and child components.
type bpModel struct {
	Logger          *slog.Logger
	terminalWidth   int
	treeContentWidth int
	treeTotalWidth  int
	contextViewport viewport.Model
	modal           teamodal.ConfirmModel
	diffViewport    *viewport.Model
	treePane        bpTreePane
	childModel      bpChildModel
}

// bpTreePane is a minimal stand-in for a tree pane component.
type bpTreePane struct{}

func (p bpTreePane) View() string { return "tree content" }

// bpChildModel is a minimal stand-in for a child component.
type bpChildModel struct{}

func (c bpChildModel) View() string { return "child content" }

// bpMyModel is a minimal model used for the parent-applies-borders example.
type bpMyModel struct {
	childModel bpChildModel
}

func (m bpMyModel) View() string {
	childContent := m.childModel.View()
	style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
	return style.Render(childContent)
}

// TestCompile_LipglossWidthCalculation verifies the lipgloss Width() usage from best-practices.mdx:29.
func TestCompile_LipglossWidthCalculation(t *testing.T) {
	style := lipgloss.NewStyle().
		Width(100).
		PaddingLeft(1).
		PaddingRight(1).
		Border(lipgloss.RoundedBorder())
	_ = style
}

// TestCompile_TotalRenderedWidth verifies the total-rendered-width calculation from best-practices.mdx:55.
func TestCompile_TotalRenderedWidth(t *testing.T) {
	const borderWidth = 2 // 1 left + 1 right
	totalWidth := 37
	widthForLipgloss := totalWidth - borderWidth // 35

	style := lipgloss.NewStyle().
		Width(widthForLipgloss). // 35 — includes padding
		PaddingLeft(1).
		PaddingRight(2).
		Border(lipgloss.RoundedBorder())

	// Result: 35 + 2 (border) = 37 total rendered width
	_ = style
}

// TestCompile_ParentAppliesBorders verifies that the parent model applies borders from best-practices.mdx:74.
func TestCompile_ParentAppliesBorders(t *testing.T) {
	m := bpMyModel{}
	result := m.View()
	_ = result
}

// TestCompile_WidthDebugging verifies the width debugging pattern from best-practices.mdx:96.
func TestCompile_WidthDebugging(t *testing.T) {
	logger := slog.Default()
	m := bpModel{
		Logger:           logger,
		terminalWidth:    80,
		treeContentWidth: 33,
		treeTotalWidth:   37,
	}

	treeStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

	m.Logger.Info("WIDTH DEBUG calculateLayout",
		"terminalWidth", m.terminalWidth,
		"treeContentWidth", m.treeContentWidth,
		"treeTotalWidth", m.treeTotalWidth,
	)

	// Measure what actually renders
	treeRendered := treeStyle.Render(m.treePane.View())
	if lines := strings.Split(treeRendered, "\n"); len(lines) > 0 {
		treeActualWidth := ansi.StringWidth(lines[0])
		m.Logger.Info("WIDTH DEBUG View",
			"treeCalculated", m.treeTotalWidth,
			"treeActual", treeActualWidth,
			"gap", m.treeTotalWidth-treeActualWidth,
		)
	}
}

// TestCompile_DebugLogging verifies the debug logging pattern from best-practices.mdx:131.
func TestCompile_DebugLogging(t *testing.T) {
	logger := slog.Default()
	m := bpModel{Logger: logger}
	value := 42

	// WRONG — corrupts TUI display (shown for documentation; suppressed here)
	_ = fmt.Sprintf("DEBUG: value=%d\n", value)

	// CORRECT — goes to log file
	m.Logger.Debug("debug message", "value", value)
}

// TestCompile_ViewportConcreteType verifies the viewport vs teamodal Update return types from best-practices.mdx:151.
func TestCompile_ViewportConcreteType(t *testing.T) {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))
	modal := teamodal.NewYesNoModal("Continue?", &teamodal.ConfirmModelArgs{
		ScreenWidth:  80,
		ScreenHeight: 24,
	})
	m := bpModel{
		contextViewport: vp,
		modal:           modal,
	}

	msg := tea.WindowSizeMsg{Width: 80, Height: 24}

	var cmd tea.Cmd

	// viewport returns concrete type — assign directly
	m.contextViewport, cmd = m.contextViewport.Update(msg)
	_ = cmd

	// teamodal returns tea.Model — type assert
	var updated tea.Model
	updated, cmd = m.modal.Update(msg)
	m.modal = updated.(teamodal.ConfirmModel)
	_ = cmd
	_ = m
}

// TestCompile_ViewportHorizontalScrolling verifies viewport horizontal scrolling from best-practices.mdx:168.
func TestCompile_ViewportHorizontalScrolling(t *testing.T) {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))
	vp.SetHorizontalStep(4) // Must be > 0 to enable; 0 = disabled (default)

	diffVp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))

	// Reset on content change
	newDiff := "some\tdiff\tcontent"
	diffVp.SetContent(newDiff)
	diffVp.GotoTop()
	diffVp.SetXOffset(0) // Reset horizontal scroll too

	_ = vp
}

// TestCompile_TabReplacement verifies tab-to-spaces replacement before viewport.SetContent from best-practices.mdx:186.
func TestCompile_TabReplacement(t *testing.T) {
	diff := "func\tmain()\t{}"
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(24))

	diff = strings.ReplaceAll(diff, "\t", "    ")
	vp.SetContent(diff)

	_ = vp
}

// TestCompile_StringAlignment verifies ASCII-based string alignment from best-practices.mdx:215.
func TestCompile_StringAlignment(t *testing.T) {
	// Unreliable alignment
	unreliable := []string{
		"  Commit Batch #1",
		"  Recommendation: 6 commits",
	}

	// Reliable alignment
	reliable := []string{
		"[U] Commit Batch #1",
		"[R] Recommendation: 6 commits",
	}

	_ = unreliable
	_ = reliable
}
