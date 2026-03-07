package teatree

import (
	"fmt"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teadrpdwn"
)

//DEBUG THIS

// DrillDownModel is a Bubble Tea model for drill-down path navigation.
// It displays a vertical breadcrumb trail from root to leaf, allowing the user
// to navigate levels and switch between sibling nodes via a dropdown.
type DrillDownModel[T any] struct {
	Keys DrillDownKeyMap // Keyboard bindings

	// Tree structure
	root *Node[T] // Complete tree root

	// Path is a slice of pointers to nodes in the tree
	Path          []*Node[T]               `json:"path,omitempty"`
	SelectedLevel int                      `json:"selected_level,omitempty"`
	SelectorFunc  DrillDownSelectorFunc[T] `json:"-"`

	// Display state
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	// Prompt is an optional value displayed above the path tree
	Prompt string `json:"prompt,omitempty"`
	// InsertLine when true inserts a line between the prompt and the first item.
	InsertLine bool `json:"insert_line,omitempty"`

	// Dropdown for alternatives
	Dropdown     teadrpdwn.DropdownModel `json:"dropdown"`
	DropdownOpen bool                    `json:"dropdown_open,omitempty"`

	// Styling (public for customization)
	PathStyle     lipgloss.Style `json:"path_style"`
	SelectedStyle lipgloss.Style `json:"selected_style"`
	BorderStyle   lipgloss.Style `json:"border_style"`
}

// DrillDownArgs contains initialization arguments for NewDrillDownModel
type DrillDownArgs[T any] struct {
	SelectorFunc DrillDownSelectorFunc[T]
	Prompt       string
}

// NewDrillDownModel creates a new DrillDownModel.
// The root should be a complete Node tree.
// The selector determines which child to follow at each level to build the path.
// Call Initialize() after construction to validate inputs and build the initial path.
func NewDrillDownModel[T any](root *Node[T], args DrillDownArgs[T]) DrillDownModel[T] {
	return DrillDownModel[T]{
		Keys:          DefaultDrillDownKeyMap(),
		root:          root,
		SelectorFunc:  args.SelectorFunc,
		Prompt:        args.Prompt,
		PathStyle:     DefaultDrillDownPathStyle(),
		SelectedStyle: DefaultDrillDownSelectedStyle(),
		BorderStyle:   DefaultDrillDownBorderStyle(),
	}
}

// Initialize validates inputs and builds the initial path.
// This must be called after NewDrillDownModel() before using the model.
func (m DrillDownModel[T]) Initialize() (model DrillDownModel[T], err error) {
	var path []*Node[T]

	model = m

	if model.root == nil {
		err = NewErr(ErrDrillDown, ErrInvalidNode, "reason", "nil root")
		goto end
	}

	if model.SelectorFunc == nil {
		err = NewErr(ErrDrillDown, ErrInvalidNode, "reason", "nil selector")
		goto end
	}

	// Build initial path using selector
	path, err = buildDrillDownPath(model.root, model.SelectorFunc)
	if err != nil {
		goto end
	}
	model.Path = path

	// Start with leaf node selected (bottom of list - common case)
	model.SelectedLevel = len(path) - 1

end:
	return model, err
}

// Root returns the root node
func (m DrillDownModel[T]) Root() *Node[T] {
	return m.root
}

// WithPathStyle sets custom path item style
func (m DrillDownModel[T]) WithPathStyle(style lipgloss.Style) DrillDownModel[T] {
	m.PathStyle = style
	return m
}

// WithSelectedStyle sets custom selected item style
func (m DrillDownModel[T]) WithSelectedStyle(style lipgloss.Style) DrillDownModel[T] {
	m.SelectedStyle = style
	return m
}

// WithBorderStyle sets custom border style
func (m DrillDownModel[T]) WithBorderStyle(style lipgloss.Style) DrillDownModel[T] {
	m.BorderStyle = style
	return m
}

// Init implements tea.Model
func (m DrillDownModel[T]) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model (FOLLOWS ClearPath)
//
//goland:noinspection GoAssignmentToReceiver
func (m DrillDownModel[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var keyMsg tea.KeyPressMsg
	var ok bool
	var handled bool

	switch t := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = t.Width
		m.Height = t.Height
		goto end

	case teadrpdwn.OptionSelectedMsg:
		if !m.DropdownOpen {
			goto end
		}
		m, cmd = m.handleDropdownSelection(t)
		goto end
	case teadrpdwn.DropdownCancelledMsg:
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
func (m DrillDownModel[T]) View() tea.View {
	var view string
	var lines []string
	var line string
	var i int
	var node *Node[T]
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
	for i, node = range m.Path {
		switch {
		case i == m.SelectedLevel:
			style = m.SelectedStyle
		default:
			style = m.PathStyle
		}

		// Format: "▶ node-name" or "  node-name"
		prefix := "  "
		if hasAlternatives(node) {
			prefix = "▶ "
		}
		line = fmt.Sprintf(" %s%s ", prefix, node.Name())

		// Apply style with full width for full-width highlight
		lines = append(lines, style.Width(contentWidth).Render(line))
	}

	baseView = lipgloss.JoinVertical(lipgloss.Left, lines...)

	// Apply border
	view = m.BorderStyle.
		Width(m.Width).
		Height(m.Height).
		Render(baseView)

	// Overlay dropdown if open
	if m.DropdownOpen {
		dropdownView = m.Dropdown.View().Content
		view = teadrpdwn.OverlayDropdown(view, dropdownView, m.Dropdown.Row, m.Dropdown.Col)
	}

	return tea.NewView(view)
}

// handleDropdownSelection processes dropdown item selection
func (m DrillDownModel[T]) handleDropdownSelection(selectedMsg teadrpdwn.OptionSelectedMsg) (model DrillDownModel[T], cmd tea.Cmd) {
	var currentNode *Node[T]
	var alts []*Node[T]
	var newNode *Node[T]
	var newPath []*Node[T]
	var err error

	model = m
	model.DropdownOpen = false

	// Get alternatives for current node
	currentNode = model.Path[model.SelectedLevel]
	alts = alternatives(currentNode)

	// Check if selection is different from current node
	if selectedMsg.Index >= 0 && selectedMsg.Index < len(alts) {
		newNode = alts[selectedMsg.Index]

		// Skip rebuild if user selected the same node
		if newNode == currentNode {
			goto end
		}

		// Rebuild path: keep Path[0:m.SelectedLevel], rebuild from newNode onwards
		newPath, err = rebuildDrillDownPath(model.Path, model.SelectedLevel, newNode, model.SelectorFunc)
		if err == nil {
			model.Path = newPath
			cmd = func() tea.Msg {
				return DrillDownChangeMsg[T]{
					Level: model.SelectedLevel,
					Node:  newNode,
				}
			}
		}
	}

end:
	return model, cmd
}

// handleDropdownCancellation closes the dropdown without making a selection
func (m DrillDownModel[T]) handleDropdownCancellation() (model DrillDownModel[T]) {
	model = m
	model.DropdownOpen = false
	return model
}

// handleDropdownInput delegates input to dropdown when it's open
func (m DrillDownModel[T]) handleDropdownInput(msg tea.Msg) (model DrillDownModel[T], cmd tea.Cmd, handled bool) {
	var dropdown tea.Model

	model = m
	dropdown, cmd = model.Dropdown.Update(msg)
	model.Dropdown = dropdown.(teadrpdwn.DropdownModel)
	handled = cmd != nil

	return model, cmd, handled
}

// handleUpKey moves selection up in the path
func (m DrillDownModel[T]) handleUpKey() (model DrillDownModel[T], cmd tea.Cmd) {
	var currentNode *Node[T]

	model = m

	if model.SelectedLevel > 0 {
		model.SelectedLevel--
		currentNode = model.Path[model.SelectedLevel]
		cmd = func() tea.Msg {
			return DrillDownFocusMsg[T]{
				Level: model.SelectedLevel,
				Node:  currentNode,
			}
		}
	}

	return model, cmd
}

// handleDownKey moves selection down in the path
func (m DrillDownModel[T]) handleDownKey() (model DrillDownModel[T], cmd tea.Cmd) {
	var currentNode *Node[T]

	model = m

	if model.SelectedLevel < len(model.Path)-1 {
		model.SelectedLevel++
		currentNode = model.Path[model.SelectedLevel]
		cmd = func() tea.Msg {
			return DrillDownFocusMsg[T]{
				Level: model.SelectedLevel,
				Node:  currentNode,
			}
		}
	}

	return model, cmd
}

// handleOpenDropdown opens the dropdown with alternatives for the current node
func (m DrillDownModel[T]) handleOpenDropdown() (model DrillDownModel[T], cmd tea.Cmd) {
	var currentNode *Node[T]
	var alts []*Node[T]
	var items []string
	var i int
	var alt *Node[T]
	var row, col int

	model = m
	currentNode = model.Path[model.SelectedLevel]

	if !hasAlternatives(currentNode) {
		goto end
	}

	// Get all alternatives at this level (includes current node)
	alts = alternatives(currentNode)

	// Build dropdown items from alternatives
	items = make([]string, len(alts))
	for i, alt = range alts {
		items[i] = alt.Name()
	}

	// Position dropdown one row below current path level
	row = model.SelectedLevel + 3 // Account for title/padding + 1 row below
	col = 2                       // Indent

	model.Dropdown = teadrpdwn.NewDropdownModel(teadrpdwn.ToOptions(items), &teadrpdwn.DropdownModelArgs{
		FieldRow:     row,
		FieldCol:     col,
		ScreenWidth:  model.Width,
		ScreenHeight: model.Height,
	})

	// Set initial selection to current node
	for i, alt = range alts {
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
func (m DrillDownModel[T]) handleEnterKey() (model DrillDownModel[T], cmd tea.Cmd) {
	var currentNode *Node[T]

	model = m
	currentNode = model.Path[model.SelectedLevel]

	// Confirm selection on leaf node - send message to parent
	if !currentNode.HasChildren() {
		cmd = func() tea.Msg {
			return DrillDownSelectMsg[T]{
				Node: currentNode,
			}
		}
	}

	return model, cmd
}
