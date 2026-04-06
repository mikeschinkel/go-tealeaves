package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// SelectableItem implements teamodal.MultiSelectItem
type SelectableItem struct {
	id    string
	label string
}

func (s SelectableItem) ID() string    { return s.id }
func (s SelectableItem) Label() string { return s.label }

type model struct {
	packageModal    teamodal.MultiSelectModel[SelectableItem]
	permissionModal teamodal.MultiSelectModel[SelectableItem]
	cleanupModal    teamodal.MultiSelectModel[SelectableItem]
	statusMsg       string
	width           int
	height          int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	// Let active modal handle message first
	if m.packageModal.IsOpen() {
		m.packageModal, cmd = m.packageModal.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	if m.permissionModal.IsOpen() {
		m.permissionModal, cmd = m.permissionModal.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	if m.cleanupModal.IsOpen() {
		m.cleanupModal, cmd = m.cleanupModal.Update(msg)
		if cmd != nil {
			return m, cmd
		}
	}

	// Handle messages that bubble up from modals
	switch msg := msg.(type) {
	case teamodal.MultiSelectButtonPressedMsg[SelectableItem]:
		labels := make([]string, 0, len(msg.Selected))
		for _, item := range msg.Selected {
			labels = append(labels, item.Label())
		}
		m.statusMsg = fmt.Sprintf("Action: %s, Selected: [%s]", msg.ButtonID, strings.Join(labels, ", "))
		return m, nil

	case teamodal.MultiSelectCancelledMsg:
		m.statusMsg = "Cancelled (Esc)"
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "p":
			m.statusMsg = ""
			m.packageModal, cmd = m.packageModal.Open()
			return m, cmd

		case "r":
			m.statusMsg = ""
			m.permissionModal, cmd = m.permissionModal.Open()
			return m, cmd

		case "c":
			m.statusMsg = ""
			m.cleanupModal, cmd = m.cleanupModal.Open()
			return m, cmd
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.packageModal = m.packageModal.SetSize(msg.Width, msg.Height)
		m.permissionModal = m.permissionModal.SetSize(msg.Width, msg.Height)
		m.cleanupModal = m.cleanupModal.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	var view string
	var baseView strings.Builder
	var lines []string

	baseView.WriteString("TeaModal MultiSelect Example\n")
	baseView.WriteString("================================\n\n")
	baseView.WriteString("Commands:\n")
	baseView.WriteString("  p - Package installer (all checked, scrollbar)\n")
	baseView.WriteString("  r - Permission manager (none checked)\n")
	baseView.WriteString("  c - Cleanup tool (footer warning, 3 buttons)\n")
	baseView.WriteString("  q - Quit\n\n")

	if m.statusMsg != "" {
		baseView.WriteString(fmt.Sprintf("Status: %s\n", m.statusMsg))
	}

	view = baseView.String()

	// Pad view to fill screen
	lines = strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	// Composite modal if open
	if m.packageModal.IsOpen() {
		view = m.packageModal.OverlayModal(view)
		goto end
	}

	if m.permissionModal.IsOpen() {
		view = m.permissionModal.OverlayModal(view)
		goto end
	}

	if m.cleanupModal.IsOpen() {
		view = m.cleanupModal.OverlayModal(view)
		goto end
	}

end:
	v := tea.NewView(view)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func main() {
	teamodal.EnsureTermGetSize(os.Stdout.Fd())

	// Use case 1: Package installer (all checked by default, 6 items with MaxVisible=5)
	packageModal := teamodal.NewMultiSelectModel([]SelectableItem{
		{id: "fmt", label: "fmt"},
		{id: "net-http", label: "net/http"},
		{id: "encoding-json", label: "encoding/json"},
		{id: "os", label: "os"},
		{id: "io", label: "io"},
		{id: "context", label: "context"},
	}, &teamodal.MultiSelectModelArgs{
		Title:      "Install Packages",
		Message:    "Select packages to install:",
		AllChecked: true,
		MaxVisible: 5,
		Buttons: []teamodal.MultiSelectButton{
			{Label: "Install Selected", Hotkey: 'i', ID: "install"},
			{Label: "Skip", Hotkey: 's', ID: "skip"},
		},
	})

	// Use case 2: Permission manager (none checked by default, 4 items)
	permissionModal := teamodal.NewMultiSelectModel([]SelectableItem{
		{id: "read", label: "Read"},
		{id: "write", label: "Write"},
		{id: "execute", label: "Execute"},
		{id: "admin", label: "Admin"},
	}, &teamodal.MultiSelectModelArgs{
		Title:      "Grant Permissions",
		Message:    "Select permissions to grant to new user:",
		AllChecked: false,
		Buttons: []teamodal.MultiSelectButton{
			{Label: "Grant", Hotkey: 'g', ID: "grant"},
		},
	})

	// Use case 3: Cleanup tool (all checked, footer warning, 3 buttons)
	cleanupModal := teamodal.NewMultiSelectModel([]SelectableItem{
		{id: "tmp-logs", label: "/tmp/app-logs/"},
		{id: "cache", label: "~/.cache/build/"},
		{id: "coverage", label: "coverage.out"},
		{id: "dist", label: "dist/"},
		{id: "node-modules", label: "node_modules/"},
	}, &teamodal.MultiSelectModelArgs{
		Title:      "Cleanup Temporary Files",
		Message:    "The following files can be removed:",
		Footer:     "Warning: This action cannot be undone",
		AllChecked: true,
		Buttons: []teamodal.MultiSelectButton{
			{Label: "Delete Selected", Hotkey: 'd', ID: "delete"},
			{Label: "Keep All", Hotkey: 'k', ID: "keep"},
		},
	})

	m := model{
		packageModal:    packageModal,
		permissionModal: permissionModal,
		cleanupModal:    cleanupModal,
		width:           80,
		height:          24,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
