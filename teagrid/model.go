package teagrid

import (
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// GridModel is the main grid model. Create using NewGridModel().
type GridModel struct {
	// Data
	columns  []Column
	rows     []Row
	metadata map[string]any

	// Caches
	visibleRowCacheUpdated bool
	visibleRowCache        []Row

	// Missing data indicator
	missingDataIndicator any

	// Interaction
	focused               bool
	keyMap                KeyMap
	selectableRows        bool
	selectColumn          bool
	rowCursorIndex        int
	cellCursorMode        bool
	cellCursorColumnIndex int

	// Events
	lastUpdateUserEvents []UserEvent

	// Styles
	baseStyle       lipgloss.Style
	highlightStyle  lipgloss.Style
	cellCursorStyle lipgloss.Style
	headerStyle     lipgloss.Style
	footerStyle     lipgloss.Style
	rowStyleFunc    func(RowStyleFuncInput) lipgloss.Style
	border          BorderConfig

	selectedText   string
	unselectedText string

	// Header
	headerVisible bool

	// Footer
	footerVisible bool
	staticFooter  string

	// Pagination
	pageSize           int
	currentPage        int
	paginationWrapping bool

	// Sorting
	sortOrder []SortColumn

	// Filter
	filtered        bool
	filterTextInput textinput.Model
	filterFunc      FilterFunc

	// Dimensions
	viewportWidth  int
	viewportHeight int

	// Scrolling
	horizontalScrollOffsetCol          int
	horizontalScrollFreezeColumnsCount int
	maxHorizontalColumnIndex           int

	// Height
	minimumHeight int

	// Editing stubs (v0.3.0)
	editable      bool
	cellValidator CellValidatorFunc

	// Help keys
	additionalShortHelpKeys func() []key.Binding
	additionalFullHelpKeys  func() []key.Binding

	// Overflow indicator
	overflowIndicator bool
}

const selectColumnKey = "___select___"

var (
	defaultHighlightStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#874BFD"))
	defaultCellCursorStyle = lipgloss.NewStyle().Reverse(true)
)

// NewGridModel creates a new grid with the given columns.
// Defaults: left-aligned text, visible highlight (purple), cell cursor is
// Reverse, rounded borders, no right-align on baseStyle (fixes v0.1.0 #1).
func NewGridModel(columns []Column) GridModel {
	filterInput := textinput.New()
	filterInput.Prompt = "/"

	m := GridModel{
		columns:            make([]Column, len(columns)),
		metadata:           make(map[string]any),
		keyMap:             DefaultKeyMap(),
		border:             BorderRounded(),
		headerVisible:      true,
		footerVisible:      true,
		highlightStyle:     defaultHighlightStyle,
		cellCursorStyle:    defaultCellCursorStyle,
		baseStyle:          lipgloss.NewStyle(),
		filterTextInput:    filterInput,
		filterFunc:         filterFuncContains,
		selectedText:       "[x]",
		unselectedText:     "[ ]",
		paginationWrapping: true,
	}

	copy(m.columns, columns)
	m.recalculateWidth()

	return m
}

// Init initializes the grid per the Bubble Tea architecture.
func (m GridModel) Init() tea.Cmd {
	return nil
}

// SetSize sets the viewport dimensions and auto-configures the grid.
// Width triggers flex column resolution and fill/scroll mode.
// Height triggers automatic page size computation.
func (m GridModel) SetSize(width, height int) GridModel {
	m.viewportWidth = width
	m.viewportHeight = height
	m.recalculateWidth()
	m.recalculatePageSize()
	return m
}
