package teadep

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teadrpdwn"
)

// PathViewerModel is a Bubble Tea model for dependency path visualization
type PathViewerModel struct {
	Keys PathViewerKeyMap // Keyboard bindings

	// Tree structure
	Root *Tree `json:"root,omitempty"` // Complete tree built upfront

	// Path is a slice of pointers to nodes in the tree
	// When rebuilding: keep Path[0:n-1], rebuild from Path[n] onwards
	Path          []*Tree      `json:"path,omitempty"`
	SelectedLevel int          `json:"selected_level,omitempty"` // Which level in path is currently selected
	SelectorFunc  SelectorFunc `json:"selector_func,omitempty"`

	// Display state
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	// Prompt is an optional value to be displayed above the path tree
	Prompt string `json:"prompt,omitempty"`
	// InsertLine when true inserts a line between the prompt and the first item.
	InsertLine bool `json:"insert_line,omitempty"`

	// Dropdown for alternatives
	Dropdown     teadrpdwn.DropdownModel `json:"dropdown"`
	DropdownOpen bool                `json:"dropdown_open,omitempty"`

	// Styling (public for customization)
	PathStyle     lipgloss.Style `json:"path_style"`
	SelectedStyle lipgloss.Style `json:"selected_style"`
	BorderStyle   lipgloss.Style `json:"border_style"`
}

type PathViewerArgs struct {
	SelectorFunc SelectorFunc
	Prompt       string
}

// NewPathViewer creates a new PathViewerModel
// The root should be a complete Tree structure (built via BuildTree)
// The selector determines which child to follow at each level to build the best path
// The selector can capture any necessary context (metadata, etc.) in a closure
// Call Initialize() after construction to validate inputs and build the initial path
func NewPathViewer(root *Tree, args PathViewerArgs) (m PathViewerModel) {
	return PathViewerModel{
		Keys:          DefaultPathViewerKeyMap(),
		Root:          root,
		SelectorFunc:  args.SelectorFunc,
		Prompt:        args.Prompt,
		PathStyle:     DefaultPathStyle(),
		SelectedStyle: DefaultSelectedStyle(),
		BorderStyle:   DefaultBorderStyle(),
	}
}

// Initialize validates inputs and builds the initial path
// This must be called after NewPathViewer() before using the model
func (m PathViewerModel) Initialize() (model PathViewerModel, err error) {
	var path []*Tree

	model = m

	if model.Root == nil {
		err = NewErr(ErrDependency, ErrInvalidNode, "reason", "nil root")
		goto end
	}

	if model.SelectorFunc == nil {
		err = NewErr(ErrDependency, ErrInvalidNode, "reason", "nil selector")
		goto end
	}

	// Build initial path using selector (which may capture metadata in closure)
	path, err = model.Root.buildPath(model)
	if err != nil {
		goto end
	}
	model.Path = path

	// Start with leaf node selected (bottom of list - common case)
	model.SelectedLevel = len(path) - 1

end:
	return model, err
}

// WithPathStyle sets custom path item style
func (m PathViewerModel) WithPathStyle(style lipgloss.Style) PathViewerModel {
	m.PathStyle = style
	return m
}

// WithSelectedStyle sets custom selected item style
func (m PathViewerModel) WithSelectedStyle(style lipgloss.Style) PathViewerModel {
	m.SelectedStyle = style
	return m
}

// WithBorderStyle sets custom border style
func (m PathViewerModel) WithBorderStyle(style lipgloss.Style) PathViewerModel {
	m.BorderStyle = style
	return m
}

// Init implements tea.Model
func (m PathViewerModel) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
//
//goland:noinspection GoAssignmentToReceiver
func (m PathViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var handled bool

	switch t := msg.(type) {

	case tea.WindowSizeMsg:
		// Handle window size messages first
		m.Width = t.Width
		m.Height = t.Height
		goto end

	case teadrpdwn.OptionSelectedMsg:
		// Handle dropdown selection
		if !m.DropdownOpen {
			goto end
		}
		m, cmd = m.handleDropdownSelection(t)
		goto end
	case teadrpdwn.DropdownCancelledMsg:
		// Handle dropdown cancellation
		if !m.DropdownOpen {
			goto end
		}
		m = m.handleDropdownCancellation()
		goto end
	default:
		// Let dropdown handle input if open
		if m.DropdownOpen {
			m, cmd, handled = m.handleDropdownInput(t)
			if handled {
				goto end
			}
		}
	}

	// Process key messages
	keyMsg, ok = msg.(tea.KeyPressMsg)
	if !ok {
		goto end
	}

	switch {
	case key.Matches(keyMsg, m.Keys.Up):
		m, cmd = m.handleUpKey()
		goto end

	case key.Matches(keyMsg, m.Keys.Down):
		m, cmd = m.handleDownKey()
		goto end

	case key.Matches(keyMsg, m.Keys.OpenDropdown):
		m, cmd = m.handleOpenDropdown()
		goto end

	case key.Matches(keyMsg, m.Keys.Select):
		m, cmd = m.handleEnterKey()
		goto end
	}

end:
	return m, cmd
}

// View implements tea.Model (FOLLOWS ClearPath)
func (m PathViewerModel) View() tea.View {
	var view string
	var lines []string
	var line string
	var i int
	var tree *Tree
	var style lipgloss.Style
	var baseView string
	var contentWidth int
	var dropdownView string

	// Calculate content width (border takes 2 columns)
	contentWidth = m.Width - 2
	if contentWidth < 0 {
		contentWidth = 0
	}

	// Title (full width)
	if m.Prompt != "" {
		line = " " + m.Prompt
		lines = append(lines, lipgloss.NewStyle().Bold(true).Width(contentWidth).Render(line))
		if m.InsertLine {
			lines = append(lines, lipgloss.NewStyle().Width(contentWidth).Render(""))
		}
	}

	// Render each level in path with full-width highlight
	for i, tree = range m.Path {
		switch {
		case i == m.SelectedLevel:
			style = m.SelectedStyle
		default:
			style = m.PathStyle
		}

		// Format: "▶ ~/Projects/go-pkgs/go-dt" or "   ~/Projects/xmlui/cli"
		prefix := "  "
		if tree.HasAlternatives() {
			prefix = "▶ "
		}
		line = fmt.Sprintf(" %s%s ", prefix, tree.Node.DisplayName())

		// Apply style with full width for full-width highlight
		lines = append(lines, style.Width(contentWidth).Render(line))
	}

	baseView = lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Apply border
	view = m.BorderStyle.
		Width(m.Width - 2).
		Height(m.Height - 2).
		Render(baseView)

	// Overlay dropdown if open
	if m.DropdownOpen {
		dropdownView = m.Dropdown.View().Content
		view = teadrpdwn.OverlayDropdown(view, dropdownView, m.Dropdown.Row, m.Dropdown.Col)
		goto end
	}

	goto end

end:
	return tea.NewView(view)
}

// handleDropdownSelection processes dropdown item selection
func (m PathViewerModel) handleDropdownSelection(selectedMsg teadrpdwn.OptionSelectedMsg) (model PathViewerModel, cmd tea.Cmd) {
	var currentNode *Tree
	var alternatives []*Tree
	var newNode *Tree
	var newPath []*Tree
	var err error

	model = m
	model.DropdownOpen = false

	// Get alternatives for current tree node
	currentNode = model.Path[model.SelectedLevel]
	alternatives = currentNode.Alternatives()

	// Check if selection is different from current node
	if selectedMsg.Index >= 0 && selectedMsg.Index < len(alternatives) {
		newNode = alternatives[selectedMsg.Index]

		// Skip rebuild if user selected the same node
		if newNode == currentNode {
			goto end
		}

		// Rebuild path: keep Path[0:m.SelectedLevel], rebuild from newNode onwards using model state
		newPath, err = newNode.rebuildPath(model)
		if err == nil {
			model.Path = newPath
			// Keep selection on the node that was just changed
			cmd = func() tea.Msg {
				return ChangeNodeMsg{
					Level: model.SelectedLevel,
					Tree:  newNode,
				}
			}
		}
	}

end:
	return model, cmd
}

// handleDropdownCancellation closes the dropdown without making a selection
func (m PathViewerModel) handleDropdownCancellation() (model PathViewerModel) {
	model = m
	model.DropdownOpen = false
	return model
}

// handleDropdownInput delegates input to dropdown when it's open
func (m PathViewerModel) handleDropdownInput(msg tea.Msg) (model PathViewerModel, cmd tea.Cmd, handled bool) {
	var dropdown tea.Model

	model = m
	dropdown, cmd = model.Dropdown.Update(msg)
	model.Dropdown = dropdown.(teadrpdwn.DropdownModel)
	handled = cmd != nil

	return model, cmd, handled
}

// handleUpKey moves selection up in the path
func (m PathViewerModel) handleUpKey() (model PathViewerModel, cmd tea.Cmd) {
	var currentNode *Tree

	model = m

	if model.SelectedLevel > 0 {
		model.SelectedLevel--
		currentNode = model.Path[model.SelectedLevel]
		cmd = func() tea.Msg {
			return FocusNodeMsg{
				Level: model.SelectedLevel,
				Tree:  currentNode,
			}
		}
	}

	return model, cmd
}

// handleDownKey moves selection down in the path
func (m PathViewerModel) handleDownKey() (model PathViewerModel, cmd tea.Cmd) {
	var currentNode *Tree

	model = m

	if model.SelectedLevel < len(model.Path)-1 {
		model.SelectedLevel++
		currentNode = model.Path[model.SelectedLevel]
		cmd = func() tea.Msg {
			return FocusNodeMsg{
				Level: model.SelectedLevel,
				Tree:  currentNode,
			}
		}
	}

	return model, cmd
}

// handleOpenDropdown opens the dropdown with alternatives for current node
func (m PathViewerModel) handleOpenDropdown() (model PathViewerModel, cmd tea.Cmd) {
	var currentNode *Tree
	var alternatives []*Tree
	var items []string
	var i int
	var alt *Tree
	var row, col int

	model = m
	currentNode = model.Path[model.SelectedLevel]

	if !currentNode.HasAlternatives() {
		goto end
	}

	// Get all alternatives at this level (includes current node)
	alternatives = currentNode.Alternatives()

	// Build dropdown items from alternatives
	items = make([]string, len(alternatives))
	for i, alt = range alternatives {
		items[i] = alt.Node.DisplayName()
	}

	// Position dropdown one row below current path level
	row = model.SelectedLevel + 3 // Account for title/padding + 1 row below
	col = 2                       // Indent

	model.Dropdown = teadrpdwn.NewModel(teadrpdwn.ToOptions(items), row, col, &teadrpdwn.ModelArgs{
		ScreenWidth:  model.Width,
		ScreenHeight: model.Height,
	})

	// Set initial selection to current node
	for i, alt = range alternatives {
		if alt == currentNode {
			model.Dropdown.Selected = i
			break
		}
	}

	model.DropdownOpen = true
	model.Dropdown, cmd = model.Dropdown.Open()

end:
	return model, cmd
}

// handleEnterKey confirms selection on leaf nodes
func (m PathViewerModel) handleEnterKey() (model PathViewerModel, cmd tea.Cmd) {
	var currentNode *Tree

	model = m
	currentNode = model.Path[model.SelectedLevel]

	// Confirm selection on leaf node - send message to parent
	if currentNode.IsLeaf() {
		cmd = func() tea.Msg {
			return SelectNodeMsg{
				Tree: currentNode,
			}
		}
	}

	return model, cmd
}
