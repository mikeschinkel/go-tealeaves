# Repo: `go-tealeaves`
- Path: `~/Projects/go-pkgs/go-tealeaves`
- Module: `github.com/mikeschinkel/go-tealeaves`

## Module: `./cmd/color-viewer`

### Package: `color-viewer`
- Path: `./cmd/color-viewer`

## Module: `./teadrpdwn`

### Package: `teadrpdwn`
- Path: `./teadrpdwn`

#### Vars
- `ErrDropdown = errors.New("dropdown error")`
- `ErrEmptyOptions = errors.New("empty options error")`
- `ErrInvalidBounds = errors.New("invalid bounds error")`
- `ErrInvalidCol = errors.New("invalid column error")`
- `ErrInvalidIndex = errors.New("invalid index error")`
- `ErrInvalidRow = errors.New("invalid row error")`

#### Funcs
- `DefaultBorderStyle() lipgloss.Style`
- `DefaultItemStyle() lipgloss.Style`
- `DefaultSelectedStyle() lipgloss.Style`
- `EnsureTermGetSize(fd uintptr) (w int, h int, ok bool)`
- `OverlayDropdown(background string, foreground string, row int, col int) string`

#### Types

- `DropdownCancelledMsg struct{}`

- `DropdownKeyMap struct{}`
  - Properties
    - `Cancel key.Binding`
    - `Down key.Binding`
    - `Select key.Binding`
    - `Up key.Binding`

- `DropdownModel struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `BottomMargin int`
    - `Col int`
    - `DisplayAbove bool`
    - `FieldCol int`
    - `FieldRow int`
    - `IsOpen bool`
    - `ItemStyle lipgloss.Style`
    - `Keys DropdownKeyMap`
    - `Options []Option`
    - `Row int`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `ScrollOffset int`
    - `Selected int`
    - `SelectedStyle lipgloss.Style`
    - `TopMargin int`
  - Methods
    - `Close() (DropdownModel, tea.Cmd)`
    - `Init() tea.Cmd`
    - `Open() (DropdownModel, tea.Cmd)`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() (view string)`
    - `WithBottomMargin(margin int) DropdownModel`
    - `WithOptions(items []Option) DropdownModel`
    - `WithPosition(fieldRow int, fieldCol int) DropdownModel`
    - `WithScreenSize(width int, height int) DropdownModel`
    - `WithTopMargin(margin int) DropdownModel`

- `DropdownModelArgs struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `BottomMargin int`
    - `FieldCol int`
    - `FieldRow int`
    - `ItemStyle lipgloss.Style`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `SelectedStyle lipgloss.Style`
    - `TopMargin int`

- `Option struct{}`
  - Properties
    - `Text string`
    - `Value interface{}`

- `OptionSelectedMsg struct{}`
  - Properties
    - `Index int`
    - `Text string`
    - `Value interface{}`

- `Position int`

## Module: `./teadepview`

### Package: `teadepview`
- Path: `./teadepview`

#### Vars
- `ErrDependency = errors.New("dependency error")`
- `ErrEmptyPath = errors.New("empty path error")`
- `ErrInvalidLevel = errors.New("invalid level error")`
- `ErrInvalidNode = errors.New("invalid node error")`

#### Funcs
- `DefaultBorderStyle() lipgloss.Style`
- `DefaultPathStyle() lipgloss.Style`
- `DefaultSelectedStyle() lipgloss.Style`
- `EnsureTermGetSize(fd uintptr) (w int, h int, ok bool)`

#### Types

- `BaseNode struct{}`
  - Methods
    - `Dependencies() []Node`
    - `DisplayName() string`
    - `SetDependencies(nodes []Node)`
    - `SetDisplayName(name string)`

- `BaseNodeArgs struct{}`
  - Properties
    - `Dependencies []Node`

- `ChangeNodeMsg struct{}`
  - Properties
    - `Level int`
    - `Tree *Tree`

- `FocusNodeMsg struct{}`
  - Properties
    - `Level int`
    - `Tree *Tree`

- `Node interface{}`
  - Methods
    - `Dependencies() []Node`
    - `DisplayName() string`
    - `SetDependencies([]Node)`
    - `SetDisplayName(string)`

- `PathViewerArgs struct{}`
  - Properties
    - `Prompt string`
    - `SelectorFunc SelectorFunc`

- `PathViewerKeyMap struct{}`
  - Properties
    - `Down key.Binding`
    - `OpenDropdown key.Binding`
    - `Select key.Binding`
    - `Up key.Binding`

- `PathViewerModel struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `Dropdown teadrpdwn.DropdownModel`
    - `DropdownOpen bool`
    - `Height int`
    - `InsertLine bool`
    - `Keys PathViewerKeyMap`
    - `Path []*Tree`
    - `PathStyle lipgloss.Style`
    - `Prompt string`
    - `Root *Tree`
    - `SelectedLevel int`
    - `SelectedStyle lipgloss.Style`
    - `SelectorFunc SelectorFunc`
    - `Width int`
  - Methods
    - `Init() tea.Cmd`
    - `Initialize() (model PathViewerModel, err error)`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() (view string)`
    - `WithBorderStyle(style lipgloss.Style) PathViewerModel`
    - `WithPathStyle(style lipgloss.Style) PathViewerModel`
    - `WithSelectedStyle(style lipgloss.Style) PathViewerModel`

- `SelectNodeMsg struct{}`
  - Properties
    - `Tree *Tree`

- `SelectorFunc func(parent *Tree, children []*Tree) (best *Tree, err error)`

- `Tree struct{}`
  - Properties
    - `Children []*Tree`
    - `Node Node`
    - `Parent *Tree`
  - Methods
    - `Alternatives() []*Tree`
    - `HasAlternatives() bool`
    - `IsLeaf() bool`

## Module: `./teagrid`

### Package: `teagrid`
- Path: `./teagrid`

#### Types

- `Border struct{}`
  - Properties
    - `Bottom string`
    - `BottomJunction string`
    - `BottomLeft string`
    - `BottomRight string`
    - `InnerDivider string`
    - `InnerJunction string`
    - `Left string`
    - `LeftJunction string`
    - `Right string`
    - `RightJunction string`
    - `Top string`
    - `TopJunction string`
    - `TopLeft string`
    - `TopRight string`

- `Column struct{}`
  - Methods
    - `Filterable() bool`
    - `FlexFactor() int`
    - `FmtString() string`
    - `IsFlex() bool`
    - `Key() string`
    - `Style() lipgloss.Style`
    - `Title() string`
    - `Width() int`
    - `WithFiltered(filterable bool) Column`
    - `WithFormatString(fmtString string) Column`
    - `WithStyle(style lipgloss.Style) Column`

- `FilterFunc func(FilterFuncInput) bool`

- `FilterFuncInput struct{}`
  - Properties
    - `Columns []Column`
    - `Filter string`
    - `GlobalMetadata map[string]any`
    - `Row Row`

- `KeyMap struct{}`
  - Properties
    - `Filter key.Binding`
    - `FilterBlur key.Binding`
    - `FilterClear key.Binding`
    - `PageDown key.Binding`
    - `PageFirst key.Binding`
    - `PageLast key.Binding`
    - `PageUp key.Binding`
    - `RowDown key.Binding`
    - `RowSelectToggle key.Binding`
    - `RowUp key.Binding`
    - `ScrollLeft key.Binding`
    - `ScrollRight key.Binding`

- `Model struct{}`
  - Methods
    - `Border(border Border) Model`
    - `BorderDefault() Model`
    - `BorderRounded() Model`
    - `CurrentPage() int`
    - `Filtered(filtered bool) Model`
    - `Focused(focused bool) Model`
    - `FullHelp() [][]key.Binding`
    - `GetCanFilter() bool`
    - `GetCellCursorColumnIndex() int`
    - `GetCellCursorMode() bool`
    - `GetColumnSorting() []SortColumn`
    - `GetCurrentFilter() string`
    - `GetFocused() bool`
    - `GetFooterVisibility() bool`
    - `GetHeaderVisibility() bool`
    - `GetHighlightedRowIndex() int`
    - `GetHorizontalScrollColumnOffset() int`
    - `GetIsFilterActive() bool`
    - `GetIsFilterInputFocused() bool`
    - `GetLastUpdateUserEvents() []UserEvent`
    - `GetPaginationWrapping() bool`
    - `GetVisibleColumnRange() (start int, end int)`
    - `GetVisibleRows() []Row`
    - `HeaderStyle(style lipgloss.Style) Model`
    - `HighlightStyle(style lipgloss.Style) Model`
    - `HighlightedRow() Row`
    - `Init() tea.Cmd`
    - `KeyMap() KeyMap`
    - `MaxPages() int`
    - `PageDown() Model`
    - `PageFirst() Model`
    - `PageLast() Model`
    - `PageSize() int`
    - `PageUp() Model`
    - `ScrollLeft() Model`
    - `ScrollRight() Model`
    - `SelectableRows(selectable bool) Model`
    - `SelectedRows() []Row`
    - `ShortHelp() []key.Binding`
    - `SortByAsc(columnKey string) Model`
    - `SortByDesc(columnKey string) Model`
    - `StartFilterTyping() Model`
    - `ThenSortByAsc(columnKey string) Model`
    - `ThenSortByDesc(columnKey string) Model`
    - `TotalRows() int`
    - `Update(msg tea.Msg) (Model, tea.Cmd)`
    - `View() string`
    - `VisibleIndices() (start int, end int)`
    - `WithAdditionalFullHelpKeys(keys []key.Binding) Model`
    - `WithAdditionalShortHelpKeys(keys []key.Binding) Model`
    - `WithAllRowsDeselected() Model`
    - `WithBaseStyle(style lipgloss.Style) Model`
    - `WithCellCursorMode(enabled bool) Model`
    - `WithColumns(columns []Column) Model`
    - `WithCurrentPage(currentPage int) Model`
    - `WithFilterFunc(shouldInclude FilterFunc) Model`
    - `WithFilterInput(input textinput.Model) Model`
    - `WithFilterInputValue(value string) Model`
    - `WithFooterVisibility(visibility bool) Model`
    - `WithFuzzyFilter() Model`
    - `WithGlobalMetadata(metadata map[string]any) Model`
    - `WithHeaderVisibility(visibility bool) Model`
    - `WithHighlightedRow(index int) Model`
    - `WithHorizontalFreezeColumnCount(columnsToFreeze int) Model`
    - `WithKeyMap(keyMap KeyMap) Model`
    - `WithMaxTotalWidth(maxTotalWidth int) Model`
    - `WithMinimumHeight(minimumHeight int) Model`
    - `WithMissingDataIndicator(str string) Model`
    - `WithMissingDataIndicatorStyled(styled StyledCell) Model`
    - `WithMultiline(multiline bool) Model`
    - `WithNoPagination() Model`
    - `WithPageSize(pageSize int) Model`
    - `WithPaginationWrapping(wrapping bool) Model`
    - `WithRowStyleFunc(f func(RowStyleFuncInput) lipgloss.Style) Model`
    - `WithRows(rows []Row) Model`
    - `WithSelectedText(unselected string, selected string) Model`
    - `WithStaticFooter(footer string) Model`
    - `WithTargetWidth(totalWidth int) Model`

- `Row struct{}`
  - Properties
    - `Data RowData`
    - `Style lipgloss.Style`
  - Methods
    - `Selected(selected bool) Row`
    - `WithStyle(style lipgloss.Style) Row`

- `RowData map[string]any`

- `RowStyleFuncInput struct{}`
  - Properties
    - `Index int`
    - `IsHighlighted bool`
    - `Row Row`

- `SortColumn struct{}`
  - Properties
    - `ColumnKey string`
    - `Direction SortDirection`

- `SortDirection int`

- `StyledCell struct{}`
  - Properties
    - `Data any`
    - `Style lipgloss.Style`
    - `StyleFunc StyledCellFunc`

- `StyledCellFunc = func(input StyledCellFuncInput) lipgloss.Style`

- `StyledCellFuncInput struct{}`
  - Properties
    - `Column Column`
    - `Data any`
    - `GlobalMetadata map[string]any`
    - `Row Row`

- `UserEvent any`

- `UserEventFilterInputFocused struct{}`

- `UserEventFilterInputUnfocused struct{}`

- `UserEventHighlightedIndexChanged struct{}`
  - Properties
    - `PreviousRowIndex int`
    - `SelectedRowIndex int`

- `UserEventRowSelectToggled struct{}`
  - Properties
    - `IsSelected bool`
    - `RowIndex int`

## Module: `./teamodal`

### Package: `teamodal`
- Path: `./teamodal`

#### Vars
- `ErrCancelled = errors.New("cancelled error")`
- `ErrInvalidBounds = errors.New("invalid bounds error")`
- `ErrModal = errors.New("modal error")`

#### Funcs
- `DefaultActiveItemStyle() lipgloss.Style`
- `DefaultBorderStyle() lipgloss.Style`
- `DefaultButtonStyle() lipgloss.Style`
- `DefaultCancelKeyStyle() lipgloss.Style`
- `DefaultCancelTextStyle() lipgloss.Style`
- `DefaultEditItemStyle() lipgloss.Style`
- `DefaultFocusedButtonStyle() lipgloss.Style`
- `DefaultListFooterStyle() lipgloss.Style`
- `DefaultListItemStyle() lipgloss.Style`
- `DefaultListScrollbarStyle() lipgloss.Style`
- `DefaultListScrollbarThumbStyle() lipgloss.Style`
- `DefaultMessageStyle() lipgloss.Style`
- `DefaultSelectedItemStyle() lipgloss.Style`
- `DefaultStatusStyle() lipgloss.Style`
- `DefaultTitleStyle() lipgloss.Style`
- `EnsureTermGetSize(fd uintptr) (w int, h int, ok bool)`
- `OverlayModal(background string, foreground string, row int, col int) string`

#### Types

- `AnsweredNoMsg struct{}`

- `AnsweredYesMsg struct{}`

- `ChoiceCancelledMsg struct{}`

- `ChoiceKeyMap struct{}`
  - Properties
    - `Cancel key.Binding`
    - `Confirm key.Binding`
    - `NextButton key.Binding`
    - `PrevButton key.Binding`

- `ChoiceModel struct{}`
  - Properties
    - `Keys ChoiceKeyMap`
  - Methods
    - `Close() (ChoiceModel, tea.Cmd)`
    - `FocusButton() int`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Open() (ChoiceModel, tea.Cmd)`
    - `OverlayModal(background string) (view string)`
    - `SetSize(width int, height int) ChoiceModel`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() (view string)`

- `ChoiceModelArgs struct{}`
  - Properties
    - `AllowCancel *bool`
    - `BorderStyle lipgloss.Style`
    - `ButtonStyle lipgloss.Style`
    - `CancelKeyStyle lipgloss.Style`
    - `CancelTextStyle lipgloss.Style`
    - `DefaultIndex int`
    - `FocusedButtonStyle lipgloss.Style`
    - `Message string`
    - `MessageStyle lipgloss.Style`
    - `Options []ChoiceOption`
    - `Orientation Orientation`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `ShowBrackets *bool`
    - `ShowCancelHint *bool`
    - `Title string`
    - `TitleStyle lipgloss.Style`

- `ChoiceOption struct{}`
  - Properties
    - `Hotkey rune`
    - `ID string`
    - `Label string`

- `ChoiceSelectedMsg struct{}`
  - Properties
    - `Index int`
    - `OptionID string`

- `ClosedMsg struct{}`

- `DeleteItemRequestedMsg struct{}`
  - Properties
    - `Item T`

- `EditCompletedMsg struct{}`
  - Properties
    - `Item T`
    - `NewLabel string`

- `ItemSelectedMsg struct{}`
  - Properties
    - `Item T`

- `ListAcceptedMsg struct{}`
  - Properties
    - `Item T`

- `ListCancelledMsg struct{}`

- `ListItem interface{}`
  - Methods
    - `ID() string`
    - `IsActive() bool`
    - `Label() string`

- `ListKeyMap struct{}`
  - Properties
    - `Accept key.Binding`
    - `Cancel key.Binding`
    - `Delete key.Binding`
    - `Down key.Binding`
    - `Edit key.Binding`
    - `Help key.Binding`
    - `New key.Binding`
    - `Preview key.Binding`
    - `Up key.Binding`

- `ListModel struct{}`
  - Properties
    - `Keys ListKeyMap`
    - `Logger *slog.Logger`
  - Methods
    - `ActiveItem() (item T)`
    - `ActiveItemStyle() lipgloss.Style`
    - `BorderStyle() lipgloss.Style`
    - `ClearStatus() ListModel[T]`
    - `Close() ListModel[T]`
    - `Cursor() int`
    - `FooterStyle() lipgloss.Style`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `ItemStyle() lipgloss.Style`
    - `Items() []T`
    - `LabelWidth() int`
    - `Offset() int`
    - `Open() ListModel[T]`
    - `OverlayModal(background string) (view string)`
    - `SelectedItem() (item T)`
    - `SelectedItemStyle() lipgloss.Style`
    - `SetCursor(index int) ListModel[T]`
    - `SetCursorToLast() ListModel[T]`
    - `SetItems(items []T) ListModel[T]`
    - `SetSize(width int, height int) ListModel[T]`
    - `SetStatus(msg string) ListModel[T]`
    - `Title() string`
    - `TitleStyle() lipgloss.Style`
    - `Update(msg tea.Msg) (ListModel[T], tea.Cmd)`
    - `View() (view string)`
    - `WithActiveItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithBorderStyle(style lipgloss.Style) ListModel[T]`
    - `WithFooterStyle(style lipgloss.Style) ListModel[T]`
    - `WithItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithLabelWidth(width int) ListModel[T]`
    - `WithLogger(logger *slog.Logger) ListModel[T]`
    - `WithMaxVisible(max int) ListModel[T]`
    - `WithSelectedItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithTitle(title string) ListModel[T]`
    - `WithTitleStyle(style lipgloss.Style) ListModel[T]`

- `ListModelArgs struct{}`
  - Properties
    - `ActiveItemStyle lipgloss.Style`
    - `BorderStyle lipgloss.Style`
    - `FooterStyle lipgloss.Style`
    - `ItemStyle lipgloss.Style`
    - `LabelWidth int`
    - `MaxVisible int`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `SelectedItemStyle lipgloss.Style`
    - `Title string`
    - `TitleStyle lipgloss.Style`

- `ModalKeyMap struct{}`
  - Properties
    - `Cancel key.Binding`
    - `Confirm key.Binding`
    - `NextButton key.Binding`
    - `PrevButton key.Binding`
    - `SelectLeft key.Binding`
    - `SelectRight key.Binding`

- `ModalModel struct{}`
  - Properties
    - `Keys ModalKeyMap`
  - Methods
    - `BorderStyle() lipgloss.Style`
    - `ButtonAlign() lipgloss.Position`
    - `ButtonStyle() lipgloss.Style`
    - `Close() (ModalModel, tea.Cmd)`
    - `FocusButton() int`
    - `FocusedButtonStyle() lipgloss.Style`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Message() string`
    - `MessageAlign() lipgloss.Position`
    - `MessageStyle() lipgloss.Style`
    - `NoLabel() string`
    - `OKLabel() string`
    - `Open() (ModalModel, tea.Cmd)`
    - `OverlayModal(background string) (view string)`
    - `ScreenHeight() int`
    - `ScreenWidth() int`
    - `SetSize(width int, height int) ModalModel`
    - `Title() string`
    - `TitleAlign() lipgloss.Position`
    - `TitleStyle() lipgloss.Style`
    - `Type() ModalType`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() (view string)`
    - `WithBorderStyle(style lipgloss.Style) ModalModel`
    - `WithButtonAlign(align lipgloss.Position) ModalModel`
    - `WithButtonStyle(style lipgloss.Style) ModalModel`
    - `WithFocusedButtonStyle(style lipgloss.Style) ModalModel`
    - `WithMessage(message string) ModalModel`
    - `WithMessageAlign(align lipgloss.Position) ModalModel`
    - `WithMessageStyle(style lipgloss.Style) ModalModel`
    - `WithNoLabel(label string) ModalModel`
    - `WithOKLabel(label string) ModalModel`
    - `WithTextAlign(align lipgloss.Position) ModalModel`
    - `WithTitle(title string) ModalModel`
    - `WithTitleAlign(align lipgloss.Position) ModalModel`
    - `WithTitleStyle(style lipgloss.Style) ModalModel`
    - `WithYesLabel(label string) ModalModel`
    - `YesLabel() string`

- `ModalType int`

- `ModelArgs struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `ButtonAlign lipgloss.Position`
    - `ButtonStyle lipgloss.Style`
    - `DefaultYes bool`
    - `FocusedButtonStyle lipgloss.Style`
    - `MessageAlign lipgloss.Position`
    - `MessageStyle lipgloss.Style`
    - `NoLabel string`
    - `OKLabel string`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `TextAlign lipgloss.Position`
    - `Title string`
    - `TitleAlign lipgloss.Position`
    - `TitleStyle lipgloss.Style`
    - `YesLabel string`

- `NewItemRequestedMsg struct{}`

- `Orientation int`

- `ProgressBackgroundMsg struct{}`

- `ProgressCancelledMsg struct{}`

- `ProgressModal struct{}`
  - Properties
    - `Keys ProgressModalKeyMap`
  - Methods
    - `Close() ProgressModal`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Open() ProgressModal`
    - `OverlayModal(background string) (view string)`
    - `SetBackgroundEnabled(enabled bool) ProgressModal`
    - `SetSize(width int, height int) ProgressModal`
    - `SetTitle(title string) ProgressModal`
    - `Update(msg tea.Msg) (ProgressModal, tea.Cmd)`
    - `View() (view string)`

- `ProgressModalArgs struct{}`
  - Properties
    - `BackgroundEnabled bool`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `Title string`

- `ProgressModalKeyMap struct{}`
  - Properties
    - `Background key.Binding`
    - `Cancel key.Binding`

## Module: `./teanotify`

### Package: `teanotify`
- Path: `./teanotify`

#### Consts
- `BackColor = "#000000"`
- `DebugASCIIPrefix = "(?)"`
- `DebugColor = "#FF00FF"`
- `DebugNerdSymbol = "󰃤 "`
- `DebugUnicodePrefix = "\u003F"`
- `DefaultLerpIncrement = 0.18`
- `ErrorASCIIPrefix = "[!!]"`
- `ErrorColor = "#FF0000"`
- `ErrorNerdSymbol = "󰬅 "`
- `ErrorUnicodePrefix = "\u2718"`
- `InfoASCIIPrefix = "(i)"`
- `InfoColor = "#00FF00"`
- `InfoNerdSymbol = " "`
- `InfoUnicodePrefix = "\u24D8 "`
- `WarnColor = "#FFFF00"`
- `WarnNerdSymbol = "󱈸 "`
- `WarningASCIIPrefix = "(!)"`
- `WarningUnicodePrefix = "\u26A0"`

#### Vars
- `ErrInvalidColor = errors.New("invalid color")`
- `ErrInvalidDuration = errors.New("invalid duration")`
- `ErrInvalidNoticeKey = errors.New("invalid notice key")`
- `ErrInvalidWidth = errors.New("invalid width")`
- `ErrNotify = errors.New("notify")`

#### Types

- `NoticeDefinition struct{}`
  - Properties
    - `ForeColor string`
    - `Key NoticeKey`
    - `Prefix string`
    - `Style lipgloss.Style`

- `NoticeKey string`

- `NotifyModel struct{}`
  - Methods
    - `HasActiveNotice() (active bool)`
    - `Init() (cmd tea.Cmd)`
    - `Initialize() (out NotifyModel, err error)`
    - `NewNotifyCmd(noticeType NoticeKey, message string) (cmd tea.Cmd)`
    - `RegisterNoticeType(def NoticeDefinition) (out NotifyModel, err error)`
    - `Render(content string) (result string)`
    - `Update(msg tea.Msg) (out NotifyModel, cmd tea.Cmd)`
    - `View() (s string)`
    - `WithAllowEscToClose() (out NotifyModel)`
    - `WithMinWidth(min int) (out NotifyModel)`
    - `WithPosition(pos Position) (out NotifyModel)`
    - `WithUnicodePrefix() (out NotifyModel)`

- `NotifyOpts struct{}`
  - Properties
    - `AllowEscToClose bool`
    - `CustomNotices []NoticeDefinition`
    - `Duration time.Duration`
    - `MinWidth int`
    - `NoDefaultNotices bool`
    - `Position Position`
    - `UseNerdFont bool`
    - `UseUnicodePrefix bool`
    - `Width int`

- `Position string`

## Module: `./teastatus`

### Package: `teastatus`
- Path: `./teastatus`

#### Vars
- `ErrInvalidWidth = errors.New("invalid width")`
- `ErrStatusBar = errors.New("status bar")`

#### Funcs
- `RenderMenuLine(items []MenuItem, styles Styles) string`

#### Types

- `MenuItem struct{}`
  - Properties
    - `Binding key.Binding`
    - `Key string`
    - `Label string`

- `Model struct{}`
  - Properties
    - `Styles Styles`
  - Methods
    - `Init() tea.Cmd`
    - `SetIndicators(indicators []StatusIndicator) Model`
    - `SetMenuItems(items []MenuItem) Model`
    - `SetSize(width int) Model`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() (view string)`
    - `WithStyles(styles Styles) Model`

- `SeparatorKind int`

- `SetIndicatorsMsg struct{}`
  - Properties
    - `Indicators []StatusIndicator`

- `SetMenuItemsMsg struct{}`
  - Properties
    - `Items []MenuItem`

- `StatusIndicator struct{}`
  - Properties
    - `Style lipgloss.Style`
    - `Text string`
  - Methods
    - `WithStyle(style lipgloss.Style) StatusIndicator`

- `Styles struct{}`
  - Properties
    - `BarStyle lipgloss.Style`
    - `IndicatorSepStyle lipgloss.Style`
    - `IndicatorStyle lipgloss.Style`
    - `MenuKeyStyle lipgloss.Style`
    - `MenuLabelStyle lipgloss.Style`
    - `MenuSeparator string`
    - `SeparatorKind SeparatorKind`

## Module: `./teatxtsnip`

### Package: `teatxtsnip`
- Path: `./teatxtsnip`

#### Vars
- `SelectionStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("39")).
	Foreground(lipgloss.Color("232"))`

#### Types

- `Model struct{}`
  - Properties
    - `Logger *slog.Logger`
    - `textarea.Model`
  - Methods
    - `ClearSelection() Model`
    - `Copy() Model`
    - `Cut() Model`
    - `HasSelection() bool`
    - `Init() tea.Cmd`
    - `IsSingleLine() bool`
    - `Paste() (Model, tea.Cmd)`
    - `SelectedText() string`
    - `Selection() Selection`
    - `SelectionKeyMap() SelectionKeyMap`
    - `SetLogger(logger *slog.Logger) Model`
    - `SetSelection(sel Selection) Model`
    - `SetSelectionKeyMap(km SelectionKeyMap) Model`
    - `Update(msg tea.Msg) (Model, tea.Cmd)`
    - `View() string`

- `Position struct{}`
  - Properties
    - `Col int`
    - `Row int`
  - Methods
    - `After(other Position) bool`
    - `Before(other Position) bool`
    - `Equal(other Position) bool`

- `Selection struct{}`
  - Properties
    - `Active bool`
    - `End Position`
    - `Start Position`
  - Methods
    - `Begin(pos Position) Selection`
    - `Clear() Selection`
    - `Contains(pos Position) bool`
    - `Extend(pos Position) Selection`
    - `IsEmpty() bool`
    - `Normalized() (start Position, end Position)`

- `SelectionKeyMap struct{}`
  - Properties
    - `ClearSelection key.Binding`
    - `Copy key.Binding`
    - `Cut key.Binding`
    - `Paste key.Binding`
    - `SelectAll key.Binding`
    - `SelectDown key.Binding`
    - `SelectLeft key.Binding`
    - `SelectRight key.Binding`
    - `SelectToEnd key.Binding`
    - `SelectToLineEnd key.Binding`
    - `SelectToLineStart key.Binding`
    - `SelectToStart key.Binding`
    - `SelectUp key.Binding`
    - `SelectWordLeft key.Binding`
    - `SelectWordRight key.Binding`

## Module: `./teatree`

### Package: `teatree`
- Path: `./teatree`

#### Vars
- `NoExpanderControls = ExpanderControls{
	Expand:        "",
	Collapse:      "",
	NotApplicable: "",
}`
- `PlusExpanderControls = ExpanderControls{
	Expand:        "+",
	Collapse:      "─",
	NotApplicable: "",
}`
- `TriangleExpanderControls = ExpanderControls{
	Expand:        "▶",
	Collapse:      "▼",
	NotApplicable: "",
}`

#### Types

- `BranchStyle struct{}`
  - Properties
    - `EmptySpace string`
    - `ExpanderControls ExpanderControls`
    - `Horizontal string`
    - `LastChild string`
    - `MiddleChild string`
    - `PreExpanderIndent string`
    - `PreIconIndent string`
    - `PreSuffixIndent string`
    - `PreTextIndent string`
    - `Vertical string`

- `BuildFileTreeArgs struct{}`
  - Properties
    - `RootPath dt.PathSegment`

- `CompactNodeProvider struct{}`
  - Methods
    - `BranchStyle() BranchStyle`
    - `ExpanderControl(node *Node[T]) string`
    - `Icon(node *Node[T]) string`
    - `Style(node *Node[T], tree *Tree[T]) lipgloss.Style`
    - `Suffix(node *Node[T]) string`
    - `Text(node *Node[T]) string`

- `ExpanderControls struct{}`
  - Properties
    - `Collapse string`
    - `Expand string`
    - `NotApplicable string`

- `File struct{}`
  - Properties
    - `Path dt.RelFilepath`
    - `YOffset int`
  - Methods
    - `Content() string`
    - `Data() any`
    - `HasContent() bool`
    - `HasData() bool`
    - `HasMeta() bool`
    - `IsEmpty() bool`
    - `LoadMeta(root dt.DirPath) (err error)`
    - `Meta() *FileMeta`
    - `SetContent(content string)`
    - `WithData(data any) *File`
    - `WithMeta(meta *FileMeta) *File`
    - `WithYOffset(yOfs int) *File`

- `FileMeta struct{}`
  - Properties
    - `Data any`
    - `EntryStatus dt.EntryStatus`
    - `ModTime time.Time`
    - `Permissions os.FileMode`
    - `Size int64`

- `FileNode = Node[File]`

- `Model struct{}`
  - Properties
    - `Keys TreeKeyMap`
  - Methods
    - `FocusedNode() (node *Node[T])`
    - `Init() tea.Cmd`
    - `MaxLineWidth() int`
    - `SetFocusedNode(nodeID string) Model[T]`
    - `SetSize(width int, height int) Model[T]`
    - `Tree() *Tree[T]`
    - `Update(msg tea.Msg) (Model[T], tea.Cmd)`
    - `View() string`

- `Node struct{}`
  - Methods
    - `AddChild(child *Node[T])`
    - `AncestorIsLastChild() []bool`
    - `Children() []*Node[T]`
    - `Collapse()`
    - `Data() *T`
    - `Depth() int`
    - `Expand()`
    - `FindByID(id string) *Node[T]`
    - `HasChildren() bool`
    - `HasGrandChildren() bool`
    - `ID() string`
    - `InsertChildSorted(child *Node[T], less func(a, b *Node[T]) bool)`
    - `IsExpanded() bool`
    - `IsLastChild() bool`
    - `IsRoot() bool`
    - `IsVisible() bool`
    - `Name() string`
    - `Parent() *Node[T]`
    - `RemoveChild(id string) bool`
    - `SetChildren(children []*Node[T])`
    - `SetExpanded(expanded bool)`
    - `SetName(name string)`
    - `SetText(text string)`
    - `SetVisible(visible bool)`
    - `Text() string`
    - `Toggle()`

- `NodeProvider interface{}`
  - Methods
    - `BranchStyle() BranchStyle`
    - `ExpanderControl(node *Node[T]) string`
    - `Icon(node *Node[T]) string`
    - `Style(node *Node[T], tree *Tree[T]) lipgloss.Style`
    - `Suffix(node *Node[T]) string`
    - `Text(node *Node[T]) string`

- `Renderer struct{}`
  - Methods
    - `GetMaxLineWidth() int`
    - `Render() string`
    - `RenderToLines() []string`

- `SimpleNodeProvider struct{}`
  - Methods
    - `Format(node *Node[T]) string`
    - `Icon(node *Node[T]) string`
    - `Style(node *Node[T], isFocused bool) lipgloss.Style`

- `Tree struct{}`
  - Methods
    - `CollapseAll()`
    - `CollapseFocused() bool`
    - `ExpandAll()`
    - `ExpandFocused() bool`
    - `FindByID(nodeID string) *Node[T]`
    - `FirstNode() *Node[T]`
    - `FocusedNode() (node *Node[T])`
    - `IsFocusedNode(node *Node[T]) bool`
    - `MoveDown() bool`
    - `MoveUp() bool`
    - `Nodes() []*Node[T]`
    - `Provider() NodeProvider[T]`
    - `SetFocusedNode(nodeID string) bool`
    - `SetNodes(nodes []*Node[T])`
    - `SetProvider(provider NodeProvider[T])`
    - `ToggleFocused() bool`
    - `VisibleNodes() []*Node[T]`

- `TreeArgs struct{}`
  - Properties
    - `ExpanderControls *ExpanderControls`
    - `FocusedNode *Node[T]`
    - `NodeProvider NodeProvider[T]`

- `TreeKeyMap struct{}`
  - Properties
    - `CollapseOrUp key.Binding`
    - `Down key.Binding`
    - `ExpandOrEnter key.Binding`
    - `Toggle key.Binding`
    - `Up key.Binding`

- `TreeOption func(*Tree[T])`

## Module: `./teautils`

### Package: `teautils`
- Path: `./teautils`

#### Vars
- `ErrEmptyKeyID = errors.New("cannot register key with empty identifier")`
- `ErrEmptyKeyIdentifier = errors.New("empty key identifier")`
- `ErrKeyIdentifierEmptyPart = errors.New("key identifier contains empty part")`
- `ErrKeyIdentifierInvalidPart = errors.New("key identifier contains invalid part")`
- `ErrKeyIdentifierMissingDot = errors.New("key identifier must contain at least one dot separator")`
- `ErrKeyNotFound = errors.New("key not found in registry")`

#### Funcs
- `ApplyBoxBorder(borderStyle lipgloss.Style, content string) string`
- `CalculateCenter(screenW int, screenH int, modalW int, modalH int) (row int, col int)`
- `CenterModal(renderedView string, screenW int, screenH int) (width int, height int, row int, col int)`
- `FormatKeyDisplay(k KeyMeta) string`
- `GetSortedCategories(keysByCategory map[string][]KeyMeta, preferredOrder []string) []string`
- `MeasureRenderedView(renderedView string) (width int, height int)`
- `ProperCaseShortcut(s string) string`
- `RenderAlignedLine(text string, style lipgloss.Style, width int, align lipgloss.Position) string`
- `RenderCenteredLine(text string, style lipgloss.Style, width int) string`

#### Types

- `HelpVisorStyle struct{}`
  - Properties
    - `CategoryOrder []string`
    - `CategoryStyle lipgloss.Style`
    - `DescStyle lipgloss.Style`
    - `KeyColumnGap int`
    - `KeyStyle lipgloss.Style`
    - `TitleStyle lipgloss.Style`

- `KeyIdentifier string`

- `KeyMeta struct{}`
  - Properties
    - `Binding key.Binding`
    - `Category string`
    - `DisplayKeys []string`
    - `HelpModal bool`
    - `HelpText string`
    - `ID KeyIdentifier`
    - `StatusBar bool`
    - `StatusBarLabel string`

- `KeyRegistry struct{}`
  - Methods
    - `ByCategory() map[string][]KeyMeta`
    - `Clear()`
    - `ForHelpModal() []KeyMeta`
    - `ForStatusBar() []KeyMeta`
    - `Get(id KeyIdentifier) (meta *KeyMeta, err error)`
    - `MustRegister(meta KeyMeta)`
    - `MustRegisterMany(metas []KeyMeta)`
    - `Register(meta KeyMeta) (err error)`
    - `RegisterMany(metas []KeyMeta) (err error)`

## Module: `./examples/teadrpdwn/demo`

### Package: `demo`
- Path: `./examples/teadrpdwn/demo`

## Module: `./examples/teadrpdwn/simple`

### Package: `simple`
- Path: `./examples/teadrpdwn/simple`

## Module: `./examples/teadepview/treenav`

### Package: `treenav`
- Path: `./examples/teadepview/treenav`

#### Funcs
- `ExampleTree() *teadepview.Tree`

## Module: `./examples/teagrid/filtering`

### Package: `filtering`
- Path: `./examples/teagrid/filtering`

## Module: `./examples/teagrid/scrolling`

### Package: `scrolling`
- Path: `./examples/teagrid/scrolling`

## Module: `./examples/teagrid/simplest`

### Package: `simplest`
- Path: `./examples/teagrid/simplest`

## Module: `./examples/teagrid/sorting`

### Package: `sorting`
- Path: `./examples/teagrid/sorting`

## Module: `./examples/teamodal/choices`

### Package: `choices`
- Path: `./examples/teamodal/choices`

## Module: `./examples/teamodal/editlist`

### Package: `editlist`
- Path: `./examples/teamodal/editlist`

#### Types

- `Task struct{}`
  - Methods
    - `ID() string`
    - `IsActive() bool`
    - `Label() string`
    - `String() string`

## Module: `./examples/teamodal/various`

### Package: `various`
- Path: `./examples/teamodal/various`

## Module: `./examples/teamodal/vertical`

### Package: `vertical`
- Path: `./examples/teamodal/vertical`

## Module: `./examples/teanotify/simple`

### Package: `simple`
- Path: `./examples/teanotify/simple`

## Module: `./examples/teastatus/statusbar`

### Package: `statusbar`
- Path: `./examples/teastatus/statusbar`

## Module: `./examples/teatxtsnip/editor`

### Package: `editor`
- Path: `./examples/teatxtsnip/editor`

## Module: `./examples/teatree/filetree`

### Package: `filetree`
- Path: `./examples/teatree/filetree`

## Module: `./examples/teautils/keyhelp`

### Package: `keyhelp`
- Path: `./examples/teautils/keyhelp`

