package main

import (
	"fmt"
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-cliutil"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
)

// Task implements teamodal.ListItem
type Task struct {
	id     string
	label  string
	active bool
}

func (t Task) ID() string     { return t.id }
func (t Task) Label() string  { return t.label }
func (t Task) IsActive() bool { return t.active }
func (t Task) String() string { return t.label }

// model is the main application model
type model struct {
	listModal     teamodal.ListModel[Task]
	confirmModal  teamodal.ModalModel
	tasks         []Task
	pendingDelete Task
	statusMsg     string
	width         int
	height        int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var listModel teamodal.ListModel[Task]
	var confirmModel tea.Model

	// Let confirm modal handle messages first if open
	if m.confirmModal.IsOpen() {
		confirmModel, cmd = m.confirmModal.Update(msg)
		if cmd != nil {
			m.confirmModal = confirmModel.(teamodal.ModalModel)
			return m, cmd
		}
	}

	// Let list modal handle messages if open
	if m.listModal.IsOpen() {
		listModel, cmd = m.listModal.Update(msg)
		if cmd != nil {
			m.listModal = listModel
			return m, cmd
		}
	}

	// Handle messages that bubble up from modals
	switch msg := msg.(type) {
	case teamodal.ItemSelectedMsg[Task]:
		// Preview select: mark as active but keep dialog open
		m.statusMsg = fmt.Sprintf("Preview selected: %s", msg.Item.Label())
		m = m.setActiveTask(msg.Item.ID())
		return m, nil

	case teamodal.ListAcceptedMsg[Task]:
		// Dialog accepted with final selection
		if msg.Item.ID() != "" {
			m.statusMsg = fmt.Sprintf("Accepted: %s", msg.Item.Label())
		} else {
			m.statusMsg = "Accepted (no items)"
		}
		return m, nil

	case teamodal.NewItemRequestedMsg:
		// Simulate adding a new task
		newID := fmt.Sprintf("task-%d", len(m.tasks)+1)
		newTask := Task{
			id:    newID,
			label: fmt.Sprintf("New Task %d", len(m.tasks)+1),
		}
		m.tasks = append(m.tasks, newTask)
		m.listModal = m.listModal.SetItems(m.tasks)
		m.statusMsg = fmt.Sprintf("Created: %s", newTask.Label())
		return m, nil

	case teamodal.EditCompletedMsg[Task]:
		// Inline edit completed - update the task label
		m = m.updateTaskLabel(msg.Item.ID(), msg.NewLabel)
		m.statusMsg = fmt.Sprintf("Edited: %s", msg.NewLabel)
		return m, nil

	case teamodal.DeleteItemRequestedMsg[Task]:
		// Close list and show confirmation
		m.listModal = m.listModal.Close()
		m.pendingDelete = msg.Item
		m.confirmModal, cmd = m.confirmModal.Open()
		return m, cmd

	case teamodal.ListCancelledMsg:
		m.statusMsg = "List cancelled"
		return m, nil

	case teamodal.AnsweredYesMsg:
		// User confirmed deletion
		m = m.deleteTask(m.pendingDelete.ID())
		m.statusMsg = fmt.Sprintf("Deleted: %s", m.pendingDelete.Label())
		m.pendingDelete = Task{}
		return m, nil

	case teamodal.AnsweredNoMsg:
		// User cancelled deletion, reopen list
		m.statusMsg = "Deletion cancelled"
		m.pendingDelete = Task{}
		m.listModal = m.listModal.Open()
		return m, nil

	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "l":
			// Open list modal
			m.listModal = m.listModal.Open()
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.listModal = m.listModal.SetSize(msg.Width, msg.Height)
		m.confirmModal = m.confirmModal.SetSize(msg.Width, msg.Height)
		return m, nil
	}

	return m, nil
}

func (m model) View() tea.View {
	var view string
	var baseView strings.Builder
	var lines []string

	baseView.WriteString("TeaModal List Example\n")
	baseView.WriteString("========================\n\n")
	baseView.WriteString("Commands:\n")
	baseView.WriteString("  l - Open task list\n")
	baseView.WriteString("  q - Quit\n\n")

	if m.statusMsg != "" {
		baseView.WriteString(fmt.Sprintf("Status: %s\n\n", m.statusMsg))
	}

	baseView.WriteString("Current Tasks:\n")
	for _, t := range m.tasks {
		marker := "  "
		if t.IsActive() {
			marker = "> "
		}
		baseView.WriteString(fmt.Sprintf("%s%s\n", marker, t.Label()))
	}

	view = baseView.String()

	// Pad view to fill screen
	lines = strings.Split(view, "\n")
	for len(lines) < m.height {
		lines = append(lines, "")
	}
	view = strings.Join(lines, "\n")

	// Composite confirm modal if open
	if m.confirmModal.IsOpen() {
		view = m.confirmModal.OverlayModal(view)
		goto end
	}

	// Composite list modal if open
	if m.listModal.IsOpen() {
		view = m.listModal.OverlayModal(view)
		goto end
	}

end:
	v := tea.NewView(view)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

// setActiveTask sets the specified task as active and clears active on others
func (m model) setActiveTask(id string) model {
	for i := range m.tasks {
		m.tasks[i].active = (m.tasks[i].id == id)
	}
	m.listModal = m.listModal.SetItems(m.tasks)
	return m
}

// updateTaskLabel updates the label of the task with the given ID
func (m model) updateTaskLabel(id, newLabel string) model {
	for i := range m.tasks {
		if m.tasks[i].id == id {
			m.tasks[i].label = newLabel
			break
		}
	}
	m.listModal = m.listModal.SetItems(m.tasks)
	return m
}

// deleteTask removes the task with the given ID
func (m model) deleteTask(id string) model {
	var newTasks []Task
	for _, t := range m.tasks {
		if t.id != id {
			newTasks = append(newTasks, t)
		}
	}
	m.tasks = newTasks
	m.listModal = m.listModal.SetItems(m.tasks)
	return m
}

func main() {
	// Ensure terminal size is available
	teamodal.EnsureTermGetSize(os.Stdout.Fd())

	// Sample tasks
	tasks := []Task{
		{id: "task-1", label: "Review pull request", active: true},
		{id: "task-2", label: "Write documentation"},
		{id: "task-3", label: "Fix bug in parser"},
		{id: "task-4", label: "Update dependencies"},
		{id: "task-5", label: "Refactor config module"},
	}

	listModal := teamodal.NewListModel(tasks, &teamodal.ListModelArgs{
		Title:      "Select Task",
		MaxVisible: 8,
	})

	confirmModal := teamodal.NewYesNoModal("", &teamodal.ModelArgs{
		Title: "Confirm Delete",
	})

	m := model{
		listModal:    listModal,
		confirmModal: confirmModal,
		tasks:        tasks,
		width:        80,
		height:       24,
	}

	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		cliutil.Stderr("Error: %v\n", err)
		os.Exit(1)
	}
}
