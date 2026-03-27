# Repo: `go-tealeaves`
- Path: `~/Projects/go-pkgs/go-tealeaves`
- Module: `github.com/mikeschinkel/go-tealeaves`

## Module: `./cmd/color-viewer`

### Package: `color-viewer`
- Path: `./cmd/color-viewer`

## Module: `./cmd/jediterm-bug`

### Package: `jediterm-bug`
- Path: `./cmd/jediterm-bug`

## Module: `./teacrumbs`

### Package: `teacrumbs`
- Path: `./teacrumbs`

#### Vars
- `ErrBreadcrumbs = errors.New("breadcrumbs")`
- `ErrIndexOutOfRange = errors.New("index out of range")`

#### Types

- `BreadcrumbsModel struct{}`
  - Properties
    - `Styles Styles`
  - Methods
    - `Crumbs() []Crumb`
    - `HandleMouse(msg tea.MouseMsg) tea.Cmd`
    - `HitTest(x int, y int) int`
    - `Init() tea.Cmd`
    - `Len() int`
    - `Pop() BreadcrumbsModel`
    - `Position() (row int, col int)`
    - `Push(crumb Crumb) BreadcrumbsModel`
    - `Separator() string`
    - `SetCrumb(index int, crumb Crumb) BreadcrumbsModel`
    - `SetCrumbs(crumbs []Crumb) BreadcrumbsModel`
    - `SetPosition(row int, col int) BreadcrumbsModel`
    - `SetSize(width int) BreadcrumbsModel`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() tea.View`
    - `Width() int`
    - `WithSeparator(sep string) BreadcrumbsModel`
    - `WithStyles(styles Styles) BreadcrumbsModel`
    - `WithTheme(theme teautils.Theme) BreadcrumbsModel`

- `Crumb struct{}`
  - Properties
    - `Data any`
    - `Renderer Renderer`
    - `Short string`
    - `Style *lipgloss.Style`
    - `Text string`

- `CrumbArgs struct{}`
  - Properties
    - `Data any`
    - `Renderer Renderer`
    - `Short string`
    - `Style *lipgloss.Style`

- `CrumbClickedMsg struct{}`
  - Properties
    - `Button tea.MouseButton`
    - `Crumb Crumb`
    - `Index int`

- `CrumbHoverLeaveMsg struct{}`

- `CrumbHoverMsg struct{}`
  - Properties
    - `Crumb Crumb`
    - `Index int`

- `PopCrumbMsg struct{}`

- `PushCrumbMsg struct{}`
  - Properties
    - `Crumb Crumb`

- `Renderer interface{}`
  - Methods
    - `Render(index int, model BreadcrumbsModel) string`

- `RendererFunc func(index int, model BreadcrumbsModel) string`

- `SetCrumbMsg struct{}`
  - Properties
    - `Crumb Crumb`
    - `Index int`

- `SetCrumbsMsg struct{}`
  - Properties
    - `Crumbs []Crumb`

- `Styles struct{}`
  - Properties
    - `CurrentStyle lipgloss.Style`
    - `HoverStyle lipgloss.Style`
    - `ParentStyle lipgloss.Style`
    - `SeparatorStyle lipgloss.Style`

## Module: `./teadiffr`

### Package: `teadiffr`
- Path: `./teadiffr`

#### Vars
- `ErrDiff = errors.New("diff")`
- `ErrEmptyDiff = errors.New("empty diff")`
- `ErrInvalidBlock = errors.New("invalid block")`
- `ErrInvalidFile = errors.New("invalid file")`

#### Funcs
- `RenderFileDiffs(files []FileDiff, renderer DiffRenderer, width int) []string`

#### Types

- `CondensedBlock struct{}`
  - Properties
    - `ChangedLines []string`
    - `ContextAfter []string`
    - `ContextBefore []string`
    - `IsTruncated bool`
    - `LineCount int`
    - `Type string`

- `DiffRenderer interface{}`
  - Methods
    - `RenderAddedLine(line string, status FileStatus, width int) string`
    - `RenderBlockHeader(blockType string, lineCount int) string`
    - `RenderContextLine(line string, status FileStatus, width int) string`
    - `RenderDeletedLine(line string, status FileStatus, width int) string`
    - `RenderFileHeader(path string, status FileStatus, width int) string`
    - `RenderSeparator() string`
    - `RenderTruncation(status FileStatus) string`

- `FileDiff struct{}`
  - Properties
    - `Blocks []CondensedBlock`
    - `Path string`
    - `Status FileStatus`

- `FileStatus int`

- `TUIRenderer struct{}`
  - Properties
    - `AddedColor color.Color`
    - `BlockHeaderColor color.Color`
    - `ContextColor color.Color`
    - `DeletedBgColor color.Color`
    - `DeletedColor color.Color`
    - `DeletedStatusColor color.Color`
    - `FileHeaderColor color.Color`
    - `NewBgColor color.Color`
    - `NewStatusColor color.Color`
  - Methods
    - `RenderAddedLine(line string, status FileStatus, width int) string`
    - `RenderBlockHeader(blockType string, lineCount int) string`
    - `RenderContextLine(line string, status FileStatus, width int) string`
    - `RenderDeletedLine(line string, status FileStatus, width int) string`
    - `RenderFileHeader(path string, status FileStatus, width int) string`
    - `RenderSeparator() string`
    - `RenderTruncation(status FileStatus) string`

- `TUIRendererArgs struct{}`
  - Properties
    - `AddedColor color.Color`
    - `BlockHeaderColor color.Color`
    - `ContextColor color.Color`
    - `DeletedBgColor color.Color`
    - `DeletedColor color.Color`
    - `DeletedStatusColor color.Color`
    - `FileHeaderColor color.Color`
    - `NewBgColor color.Color`
    - `NewStatusColor color.Color`

## Module: `./teadiffview`

### Package: `teadiffview`
- Path: `./teadiffview`

#### Consts
- `ChangeBlockBgColor = "24"`
- `SelectionBgColor = "7"`

#### Vars
- `ANSIReset = "\x1b[0m"`
- `ChangeBlockBgANSI = "\x1b[48;5;" + ChangeBlockBgColor + "m"`
- `CommitGroupColors = []color.Color{
	lipgloss.Color("#4CAF50"),
	lipgloss.Color("#2196F3"),
	lipgloss.Color("#FF9800"),
	lipgloss.Color("#9C27B0"),
	lipgloss.Color("#00BCD4"),
	lipgloss.Color("#F44336"),
	lipgloss.Color("#FFEB3B"),
	lipgloss.Color("#795548"),
}`
- `ErrDiffView = errors.New("diffview")`
- `ErrEmptyDiff = errors.New("empty diff")`
- `ErrInvalidBlock = errors.New("invalid block")`
- `ErrInvalidContent = errors.New("invalid content")`
- `ErrInvalidFile = errors.New("invalid file")`
- `SelectionBgANSI = "\x1b[48;5;" + SelectionBgColor + "m"`
- `SingleCommitColor = lipgloss.Color("#4CAF50")`

#### Funcs
- `RenderFileDiffs(files []FileDiff, renderer DiffRenderer, width int) []string`

#### Types

- `BlockMarker struct{}`
  - Properties
    - `LineCount int`
  - Methods
    - `IsBlockStart() bool`
    - `LineNo() int`
    - `PaneLine()`

- `CondensedBlock struct{}`
  - Properties
    - `ChangedLines []string`
    - `ContextAfter []string`
    - `ContextBefore []string`
    - `IsTruncated bool`
    - `LineCount int`
    - `Type string`

- `DiffRenderer interface{}`
  - Methods
    - `RenderAddedLine(line string, status FileStatus, width int) string`
    - `RenderBlockHeader(blockType string, lineCount int) string`
    - `RenderContextLine(line string, status FileStatus, width int) string`
    - `RenderDeletedLine(line string, status FileStatus, width int) string`
    - `RenderFileHeader(path string, status FileStatus, width int) string`
    - `RenderSeparator() string`
    - `RenderTruncation(status FileStatus) string`

- `FileDiff struct{}`
  - Properties
    - `Blocks []CondensedBlock`
    - `Path string`
    - `Status FileStatus`

- `FileStatus int`

- `PaneLine interface{}`
  - Methods
    - `LineNo() int`
    - `PaneLine()`

- `PlaceholderLine struct{}`
  - Properties
    - `HunkLine int`
  - Methods
    - `IsWithinHunk() bool`
    - `LineNo() int`
    - `PaneLine()`

- `RowAnnotation struct{}`
  - Properties
    - `Char rune`
    - `Color color.Color`

- `SelectAction int`

- `SplitDiffModel struct{}`
  - Properties
    - `Logger *slog.Logger`
  - Methods
    - `Blur() SplitDiffModel`
    - `CenterOnFirstChangeBlock() SplitDiffModel`
    - `ClearSelection() SplitDiffModel`
    - `CursorIndex() int`
    - `ExtendSelectionDown() SplitDiffModel`
    - `ExtendSelectionUp() SplitDiffModel`
    - `Focus() SplitDiffModel`
    - `GetBlockIndexAtCursor() int`
    - `GetBlockRange() (start int, end int)`
    - `GetSelectedLines() (start int, end int)`
    - `GoToBottom() SplitDiffModel`
    - `GoToTop() SplitDiffModel`
    - `HasSelection() bool`
    - `Init() tea.Cmd`
    - `IsRowSelected(rowIdx int) bool`
    - `LeftLineNumWidth() int`
    - `LineCount() int`
    - `MoveCursorDown() SplitDiffModel`
    - `MoveCursorUp() SplitDiffModel`
    - `PageDown() SplitDiffModel`
    - `PageUp() SplitDiffModel`
    - `RightLineNumWidth() int`
    - `RowCount() int`
    - `Rows() []SplitPaneRow`
    - `ScrollLeft() SplitDiffModel`
    - `ScrollRight() SplitDiffModel`
    - `ScrollToColumn(col int) SplitDiffModel`
    - `ScrollToEnd() SplitDiffModel`
    - `SelectAllBlocks() SplitDiffModel`
    - `SelectCurrentBlock() SplitDiffModel`
    - `SetAnnotations(annotations map[int]RowAnnotation) SplitDiffModel`
    - `SetContent(content *diffutils.DiffContent) SplitDiffModel`
    - `SetGutter(chars []rune, colors []color.Color) SplitDiffModel`
    - `SetSize(splitContentWidth int, height int) SplitDiffModel`
    - `ToggleBlockSelection() SplitDiffModel`
    - `Update(msg tea.Msg) (SplitDiffModel, tea.Cmd)`
    - `View() (view tea.View)`

- `SplitDiffModelArgs struct{}`
  - Properties
    - `Height int`
    - `HighlightFunc func(text, language string) string`
    - `InlineHighlighter diffutils.InlineHighlighter`
    - `Logger *slog.Logger`
    - `Width int`

- `SplitPaneRow struct{}`
  - Properties
    - `ActualLine PaneLine`
    - `BlockIndex int`
    - `CommitLine PaneLine`
    - `LineOffset int`

- `TUIRenderer struct{}`
  - Properties
    - `AddedColor color.Color`
    - `BlockHeaderColor color.Color`
    - `ContextColor color.Color`
    - `DeletedBgColor color.Color`
    - `DeletedColor color.Color`
    - `DeletedStatusColor color.Color`
    - `FileHeaderColor color.Color`
    - `NewBgColor color.Color`
    - `NewStatusColor color.Color`
  - Methods
    - `RenderAddedLine(line string, status FileStatus, width int) string`
    - `RenderBlockHeader(blockType string, lineCount int) string`
    - `RenderContextLine(line string, status FileStatus, width int) string`
    - `RenderDeletedLine(line string, status FileStatus, width int) string`
    - `RenderFileHeader(path string, status FileStatus, width int) string`
    - `RenderSeparator() string`
    - `RenderTruncation(status FileStatus) string`

- `TUIRendererArgs struct{}`
  - Properties
    - `AddedColor color.Color`
    - `BlockHeaderColor color.Color`
    - `ContextColor color.Color`
    - `DeletedBgColor color.Color`
    - `DeletedColor color.Color`
    - `DeletedStatusColor color.Color`
    - `FileHeaderColor color.Color`
    - `NewBgColor color.Color`
    - `NewStatusColor color.Color`

- `TextLine struct{}`
  - Properties
    - `Text string`
  - Methods
    - `LineNo() int`
    - `PaneLine()`

## Module: `./teadiffview/examples/splitdiff`

### Package: `splitdiff`
- Path: `./teadiffview/examples/splitdiff`

## Module: `./teafields`

### Package: `teafields`
- Path: `./teafields`

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
    - `View() tea.View`
    - `WithBottomMargin(margin int) DropdownModel`
    - `WithOptions(items []Option) DropdownModel`
    - `WithPosition(fieldRow int, fieldCol int) DropdownModel`
    - `WithScreenSize(width int, height int) DropdownModel`
    - `WithTheme(theme teautils.Theme) DropdownModel`
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

- `ModelArgs = DropdownModelArgs`

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

## Module: `./teafields/examples/demo`

### Package: `demo`
- Path: `./teafields/examples/demo`

## Module: `./teafields/examples/simple`

### Package: `simple`
- Path: `./teafields/examples/simple`

## Module: `./teagrid`

### Package: `teagrid`
- Path: `./teagrid`

#### Vars
- `ErrGrid = errors.New("grid")`
- `ErrInvalidColumn = errors.New("invalid column")`
- `ErrInvalidRow = errors.New("invalid row")`

#### Types

- `BorderChars struct{}`
  - Properties
    - `BottomJunction string`
    - `BottomLeft string`
    - `BottomRight string`
    - `FreezeBottomJunction string`
    - `FreezeDivider string`
    - `FreezeInnerJunction string`
    - `FreezeTopJunction string`
    - `Horizontal string`
    - `InnerDivider string`
    - `InnerJunction string`
    - `LeftJunction string`
    - `OverflowVertical string`
    - `RightJunction string`
    - `TopJunction string`
    - `TopLeft string`
    - `TopRight string`
    - `Vertical string`

- `BorderConfig struct{}`
  - Properties
    - `Chars BorderChars`
    - `Footer BorderRegion`
    - `Header BorderRegion`
    - `Inner BorderRegion`
    - `Outer BorderRegion`
  - Methods
    - `HasFooterSeparator() bool`
    - `HasHeaderSeparator() bool`
    - `HasInnerDividers() bool`
    - `HasOuterBorder() bool`
    - `InnerDividerWidth() int`
    - `OuterWidth() int`
    - `WithChars(chars BorderChars) BorderConfig`
    - `WithFooter(region BorderRegion) BorderConfig`
    - `WithHeader(region BorderRegion) BorderConfig`
    - `WithInner(region BorderRegion) BorderConfig`
    - `WithOuter(region BorderRegion) BorderConfig`

- `BorderRegion struct{}`
  - Properties
    - `Style lipgloss.Style`
    - `Visible bool`
  - Methods
    - `WithStyle(style lipgloss.Style) BorderRegion`
    - `WithVisible(visible bool) BorderRegion`

- `CellEditStartedMsg struct{}`
  - Properties
    - `ColumnIndex int`
    - `ColumnKey string`
    - `RowIndex int`

- `CellEditedMsg struct{}`
  - Properties
    - `ColumnIndex int`
    - `ColumnKey string`
    - `NewValue any`
    - `OldValue any`
    - `RowIndex int`

- `CellStyleFunc func(CellStyleInput) lipgloss.Style`

- `CellStyleInput struct{}`
  - Properties
    - `Column Column`
    - `ColumnIndex int`
    - `Data any`
    - `GlobalMetadata map[string]any`
    - `IsColCursor bool`
    - `IsHighlightedRow bool`
    - `Row Row`
    - `RowIndex int`

- `CellValidatorFunc func(columnKey string, value any) error`

- `CellValue struct{}`
  - Properties
    - `Data any`
    - `SortValue any`
    - `Spans []Span`
    - `Style lipgloss.Style`
    - `StyleFunc CellStyleFunc`
  - Methods
    - `HasSpans() bool`
    - `SortableValue() any`

- `Column struct{}`
  - Methods
    - `Alignment() lipgloss.Position`
    - `Filterable() bool`
    - `FlexFactor() int`
    - `FmtString() string`
    - `IsFlex() bool`
    - `Key() string`
    - `MinWidth() int`
    - `PaddingLeft() int`
    - `PaddingRight() int`
    - `RenderWidth() int`
    - `Style() lipgloss.Style`
    - `Title() string`
    - `Width() int`
    - `WithAlignment(alignment lipgloss.Position) Column`
    - `WithFiltered(filterable bool) Column`
    - `WithFormatString(fmtString string) Column`
    - `WithMinWidth(n int) Column`
    - `WithPadding(left int, right int) Column`
    - `WithPaddingLeft(padding int) Column`
    - `WithPaddingRight(padding int) Column`
    - `WithStyle(style lipgloss.Style) Column`

- `FilterFunc func(FilterFuncInput) bool`

- `FilterFuncInput struct{}`
  - Properties
    - `Columns []Column`
    - `Filter string`
    - `GlobalMetadata map[string]any`
    - `Row Row`

- `FooterCell struct{}`
  - Properties
    - `Alignment lipgloss.Position`
    - `ColSpan int`
    - `ColumnKey string`
    - `Style lipgloss.Style`
    - `Value string`
  - Methods
    - `WithAlignment(a lipgloss.Position) FooterCell`
    - `WithStyle(s lipgloss.Style) FooterCell`

- `FooterRow struct{}`
  - Properties
    - `Cells []FooterCell`
    - `Style lipgloss.Style`
  - Methods
    - `WithStyle(s lipgloss.Style) FooterRow`

- `GridModel struct{}`
  - Methods
    - `Border() BorderConfig`
    - `BorderDefault() GridModel`
    - `BorderRounded() GridModel`
    - `CanFilter() bool`
    - `ColCursorColumnIndex() int`
    - `ColCursorMode() bool`
    - `ColCursorWrapping() bool`
    - `ColumnSorting() []SortColumn`
    - `CurrentFilter() string`
    - `CurrentPage() int`
    - `FillWidth() bool`
    - `FullHelp() [][]key.Binding`
    - `HeaderStyle(style lipgloss.Style) GridModel`
    - `HighlightStyle(style lipgloss.Style) GridModel`
    - `HighlightedRow() Row`
    - `HighlightedRowIndex() int`
    - `HorizontalScrollColumnOffset() int`
    - `Init() tea.Cmd`
    - `IsFilterActive() bool`
    - `IsFilterInputFocused() bool`
    - `IsFocused() bool`
    - `IsFooterVisible() bool`
    - `IsHeaderVisible() bool`
    - `IsPaginationWrapping() bool`
    - `KeyMap() KeyMap`
    - `LastUpdateUserEvents() []UserEvent`
    - `MaxPages() int`
    - `NaturalWidth() int`
    - `PageDown() GridModel`
    - `PageFirst() GridModel`
    - `PageLast() GridModel`
    - `PageSize() int`
    - `PageUp() GridModel`
    - `RowCursorWrapping() bool`
    - `ScrollLeft() GridModel`
    - `ScrollOffset() int`
    - `ScrollRight() GridModel`
    - `SelectableRows(selectable bool) GridModel`
    - `SelectedRows() []Row`
    - `SetBorder(border BorderConfig) GridModel`
    - `ShortHelp() []key.Binding`
    - `SortByAsc(columnKey string) GridModel`
    - `SortByDesc(columnKey string) GridModel`
    - `StartFilterTyping() GridModel`
    - `ThenSortByAsc(columnKey string) GridModel`
    - `ThenSortByDesc(columnKey string) GridModel`
    - `TotalRows() int`
    - `TotalWidth() int`
    - `Update(msg tea.Msg) (GridModel, tea.Cmd)`
    - `View() tea.View`
    - `VisibleIndices() (start int, end int)`
    - `VisibleRows() []Row`
    - `WithAdditionalFullHelpKeys(keys []key.Binding) GridModel`
    - `WithAdditionalShortHelpKeys(keys []key.Binding) GridModel`
    - `WithAllRowsDeselected() GridModel`
    - `WithBaseStyle(style lipgloss.Style) GridModel`
    - `WithBorder(border BorderConfig) GridModel`
    - `WithCellPadding(left int, right int) GridModel`
    - `WithCellValidator(fn CellValidatorFunc) GridModel`
    - `WithColCursorColumn(index int) GridModel`
    - `WithColCursorMode(enabled bool) GridModel`
    - `WithColCursorStyle(style lipgloss.Style) GridModel`
    - `WithColCursorWrapping(wrapping bool) GridModel`
    - `WithColumns(columns []Column) GridModel`
    - `WithCurrentPage(currentPage int) GridModel`
    - `WithCustomFilterFunc(fn FilterFunc) GridModel`
    - `WithDataRowCount(n int) GridModel`
    - `WithEditable(editable bool) GridModel`
    - `WithFillWidth(fill bool) GridModel`
    - `WithFilterFunc(fn FilterFunc) GridModel`
    - `WithFilterInput(input textinput.Model) GridModel`
    - `WithFilterInputValue(value string) GridModel`
    - `WithFiltered(filtered bool) GridModel`
    - `WithFocused(focused bool) GridModel`
    - `WithFooterAlignment(alignment lipgloss.Position) GridModel`
    - `WithFooterRows(rows ...FooterRow) GridModel`
    - `WithFooterStyle(style lipgloss.Style) GridModel`
    - `WithFooterVisibility(visible bool) GridModel`
    - `WithFuzzyFilter() GridModel`
    - `WithGlobalMetadata(metadata map[string]any) GridModel`
    - `WithHeaderStyle(style lipgloss.Style) GridModel`
    - `WithHeaderVisibility(visible bool) GridModel`
    - `WithHighlightStyle(style lipgloss.Style) GridModel`
    - `WithHighlightedRow(index int) GridModel`
    - `WithHorizontalFreezeColumnCount(count int) GridModel`
    - `WithKeyMap(keyMap KeyMap) GridModel`
    - `WithMaxTotalWidth(width int) GridModel`
    - `WithMetadata(metadata map[string]any) GridModel`
    - `WithMinimumHeight(height int) GridModel`
    - `WithMissingDataIndicator(str string) GridModel`
    - `WithMissingDataIndicatorStyled(styled StyledCell) GridModel`
    - `WithNoPagination() GridModel`
    - `WithOverflowIndicator(enabled bool) GridModel`
    - `WithPageSize(pageSize int) GridModel`
    - `WithPaginationWrapping(wrapping bool) GridModel`
    - `WithRowCursorWrapping(wrapping bool) GridModel`
    - `WithRowStyleFunc(f func(RowStyleFuncInput) lipgloss.Style) GridModel`
    - `WithRows(rows []Row) GridModel`
    - `WithSelectColumn(show bool) GridModel`
    - `WithSelectableRows(selectable bool) GridModel`
    - `WithSelectedText(unselected string, selected string) GridModel`
    - `WithSize(width int, height int) GridModel`
    - `WithStaticFooter(footer string) GridModel`
    - `WithTargetWidth(width int) GridModel`
    - `WithTheme(theme teautils.Theme) GridModel`

- `KeyMap struct{}`
  - Properties
    - `ColLeft key.Binding`
    - `ColRight key.Binding`
    - `ColSelect key.Binding`
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

- `OverflowConfig struct{}`
  - Properties
    - `LeftIndicator string`
    - `RightIndicator string`

- `Row struct{}`
  - Properties
    - `Data RowData`
    - `Style lipgloss.Style`
  - Methods
    - `ID() uint32`
    - `IsSelected() bool`
    - `Selected(selected bool) Row`
    - `WithStyle(style lipgloss.Style) Row`

- `RowData map[string]any`

- `RowEditCancelledMsg struct{}`
  - Properties
    - `RowIndex int`

- `RowEditedMsg struct{}`
  - Properties
    - `NewData RowData`
    - `OldData RowData`
    - `RowIndex int`

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

- `Span struct{}`
  - Properties
    - `Style lipgloss.Style`
    - `Text string`

- `StyledCell = CellValue`

- `StyledCellFunc = CellStyleFunc`

- `StyledCellFuncInput = CellStyleInput`

- `UserEvent any`

- `UserEventCellSelected struct{}`
  - Properties
    - `ColumnIndex int`
    - `ColumnKey string`
    - `Data any`
    - `RowIndex int`

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

## Module: `./teagrid/examples/filtering`

### Package: `filtering`
- Path: `./teagrid/examples/filtering`

## Module: `./teagrid/examples/panning`

### Package: `panning`
- Path: `./teagrid/examples/panning`

## Module: `./teagrid/examples/scrolling`

### Package: `scrolling`
- Path: `./teagrid/examples/scrolling`

## Module: `./teagrid/examples/simplest`

### Package: `simplest`
- Path: `./teagrid/examples/simplest`

## Module: `./teagrid/examples/sorting`

### Package: `sorting`
- Path: `./teagrid/examples/sorting`

## Module: `./teaguide`

### Package: `teaguide`
- Path: `./teaguide`

#### Vars
- `ErrGuide = errors.New("guide error")`

#### Types

- `ActionSelectedMsg struct{}`
  - Properties
    - `ActionKey string`

- `GuideData struct{}`
  - Properties
    - `Sections []GuideSection`
    - `Title string`

- `GuideDismissedMsg struct{}`

- `GuideItem struct{}`
  - Properties
    - `ActionKey string`
    - `BlockReason string`
    - `KeyDisplay string`
    - `Label string`
    - `Prose string`

- `GuideItemOpts struct{}`
  - Properties
    - `BlockReason string`
    - `Label string`
    - `Prose string`

- `GuideKeyMap struct{}`
  - Properties
    - `Close key.Binding`
    - `ScrollDown key.Binding`
    - `ScrollUp key.Binding`
    - `ToggleBlock key.Binding`

- `GuideModel struct{}`
  - Properties
    - `Keys GuideKeyMap`
    - `Styles GuideStyles`
  - Methods
    - `Close() GuideModel`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Open(data GuideData) (GuideModel, tea.Cmd)`
    - `OverlayModal(background string) (view string)`
    - `SetSize(w int, h int) GuideModel`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() tea.View`
    - `WithStyles(styles GuideStyles) GuideModel`

- `GuidePriority int`

- `GuideSection struct{}`
  - Properties
    - `Heading string`
    - `Items []GuideItem`
    - `Priority GuidePriority`

- `GuideStyles struct{}`
  - Properties
    - `BlockReason lipgloss.Style`
    - `BlockedHeading lipgloss.Style`
    - `BlockedItem lipgloss.Style`
    - `Border lipgloss.Style`
    - `Footer lipgloss.Style`
    - `ItemKey lipgloss.Style`
    - `ItemLabel lipgloss.Style`
    - `ItemProse lipgloss.Style`
    - `SectionHeading lipgloss.Style`
    - `Title lipgloss.Style`

## Module: `./teaguide/example`

### Package: `example`
- Path: `./teaguide/example`

## Module: `./teahelp`

### Package: `teahelp`
- Path: `./teahelp`

#### Vars
- `ErrHelpVisor = errors.New("help visor error")`

#### Types

- `ClosedMsg struct{}`

- `HelpVisorKeyMap struct{}`
  - Properties
    - `Close key.Binding`
    - `NextPage key.Binding`
    - `PrevPage key.Binding`

- `HelpVisorModel struct{}`
  - Properties
    - `Keys HelpVisorKeyMap`
    - `Styles HelpVisorStyles`
  - Methods
    - `Close() HelpVisorModel`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Open(keysByCategory map[string][]teautils.KeyMeta) HelpVisorModel`
    - `Page() int`
    - `SetSize(width int, height int) HelpVisorModel`
    - `TotalPages() int`
    - `Update(msg tea.Msg) (HelpVisorModel, tea.Cmd)`
    - `View() (view tea.View)`
    - `WithContentStyle(style teautils.HelpVisorStyle) HelpVisorModel`
    - `WithKeys(keys HelpVisorKeyMap) HelpVisorModel`
    - `WithStyles(styles HelpVisorStyles) HelpVisorModel`
    - `WithTheme(theme teautils.Theme) HelpVisorModel`

- `HelpVisorStyles struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `FooterKeyStyle lipgloss.Style`
    - `FooterLabelStyle lipgloss.Style`

## Module: `./teahilite`

### Package: `teahilite`
- Path: `./teahilite`

#### Consts
- `DefaultFormatterName = "terminal256"`
- `DefaultStyleName = "monokai"`

#### Vars
- `ErrFormat = errors.New("format")`
- `ErrHighlight = errors.New("highlight")`
- `ErrTokenize = errors.New("tokenize")`

#### Funcs
- `DetectLanguage[S ~string](path S) (name string)`
- `Highlight(code string, language string) (string, error)`
- `HighlightLines(code string, language string) ([]string, error)`

#### Types

- `Highlighter struct{}`
  - Methods
    - `Highlight(code string, language string) (result string, err error)`
    - `HighlightLines(code string, language string) (lines []string, err error)`

- `HighlighterArgs struct{}`
  - Properties
    - `FormatterName string`
    - `StyleName string`

## Module: `./tealayout`

### Package: `tealayout`
- Path: `./tealayout`

#### Vars
- `ErrZeroDimensions = errors.New("zero dimensions")`
- `Percent100 = Percent(100)`
- `Percent20 = Percent(20)`
- `Percent25 = Percent(25)`
- `Percent33 = Percent(33)`
- `Percent50 = Percent(50)`
- `Percent75 = Percent(75)`

#### Types

- `Component struct{}`
  - Methods
    - `ChildRect(i int) Rect`
    - `Direction() Direction`
    - `MarkDirty()`
    - `Render() (string, error)`
    - `Resolve() ([]int, error)`
    - `SetSize(width int, height int)`
    - `View() string`
    - `WithGap(n int) *Component`
    - `WithMaxSize(n int) *Component`
    - `WithMinSize(n int) *Component`
    - `WithOptional(b bool) *Component`

- `Dimension struct{}`

- `Direction int`

- `Element interface{}`

- `Layout struct{}`
  - Methods
    - `Height() int`
    - `MarkDirty()`
    - `Render() (string, error)`
    - `Root() *Component`
    - `SetSize(width int, height int)`
    - `Width() int`

- `Option func(*config)`

- `Rect struct{}`
  - Properties
    - `Height int`
    - `Width int`
    - `X int`
    - `Y int`

- `SetSizer interface{}`
  - Methods
    - `SetSize(width int, height int)`

- `Size struct{}`
  - Properties
    - `Height int`
    - `Width int`

- `SizeHint struct{}`
  - Properties
    - `Desired Size`
    - `Max Size`
    - `Min Size`

- `SizeHinter interface{}`
  - Methods
    - `SizeHint(availWidth int, availHeight int) SizeHint`

- `Styler interface{}`
  - Methods
    - `Style() lipgloss.Style`

- `Viewer interface{}`
  - Methods
    - `View() string`

## Module: `./tealayout/examples/collapsible`

### Package: `collapsible`
- Path: `./tealayout/examples/collapsible`

## Module: `./tealayout/examples/golden`

### Package: `golden`
- Path: `./tealayout/examples/golden`

## Module: `./tealayout/examples/threepane`

### Package: `threepane`
- Path: `./tealayout/examples/threepane`

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
- `DefaultCheckedStyle() lipgloss.Style`
- `DefaultEditItemStyle() lipgloss.Style`
- `DefaultFocusedButtonStyle() lipgloss.Style`
- `DefaultListFooterStyle() lipgloss.Style`
- `DefaultListItemStyle() lipgloss.Style`
- `DefaultListScrollbarStyle() lipgloss.Style`
- `DefaultListScrollbarThumbStyle() lipgloss.Style`
- `DefaultMessageStyle() lipgloss.Style`
- `DefaultMultiSelectFooterStyle() lipgloss.Style`
- `DefaultSelectedItemStyle() lipgloss.Style`
- `DefaultStatusStyle() lipgloss.Style`
- `DefaultTitleStyle() lipgloss.Style`
- `DefaultUncheckedStyle() lipgloss.Style`
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
    - `View() tea.View`
    - `WithTheme(theme teautils.Theme) ChoiceModel`

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

- `ConfirmModel struct{}`
  - Properties
    - `Keys ModalKeyMap`
  - Methods
    - `BorderStyle() lipgloss.Style`
    - `ButtonAlign() lipgloss.Position`
    - `ButtonStyle() lipgloss.Style`
    - `Close() (ConfirmModel, tea.Cmd)`
    - `FocusButton() int`
    - `FocusedButtonStyle() lipgloss.Style`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `Message() string`
    - `MessageAlign() lipgloss.Position`
    - `MessageStyle() lipgloss.Style`
    - `NoLabel() string`
    - `OKLabel() string`
    - `Open() (ConfirmModel, tea.Cmd)`
    - `OverlayModal(background string) (view string)`
    - `ScreenHeight() int`
    - `ScreenWidth() int`
    - `SetSize(width int, height int) ConfirmModel`
    - `Title() string`
    - `TitleAlign() lipgloss.Position`
    - `TitleStyle() lipgloss.Style`
    - `Type() ModalType`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() tea.View`
    - `WithBorderStyle(style lipgloss.Style) ConfirmModel`
    - `WithButtonAlign(align lipgloss.Position) ConfirmModel`
    - `WithButtonStyle(style lipgloss.Style) ConfirmModel`
    - `WithFocusedButtonStyle(style lipgloss.Style) ConfirmModel`
    - `WithMessage(message string) ConfirmModel`
    - `WithMessageAlign(align lipgloss.Position) ConfirmModel`
    - `WithMessageStyle(style lipgloss.Style) ConfirmModel`
    - `WithNoLabel(label string) ConfirmModel`
    - `WithOKLabel(label string) ConfirmModel`
    - `WithTextAlign(align lipgloss.Position) ConfirmModel`
    - `WithTheme(theme teautils.Theme) ConfirmModel`
    - `WithTitle(title string) ConfirmModel`
    - `WithTitleAlign(align lipgloss.Position) ConfirmModel`
    - `WithTitleStyle(style lipgloss.Style) ConfirmModel`
    - `WithYesLabel(label string) ConfirmModel`
    - `YesLabel() string`

- `ConfirmModelArgs struct{}`
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
    - `View() tea.View`
    - `WithActiveItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithBorderStyle(style lipgloss.Style) ListModel[T]`
    - `WithFooterStyle(style lipgloss.Style) ListModel[T]`
    - `WithItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithLabelWidth(width int) ListModel[T]`
    - `WithLogger(logger *slog.Logger) ListModel[T]`
    - `WithMaxVisible(max int) ListModel[T]`
    - `WithSelectedItemStyle(style lipgloss.Style) ListModel[T]`
    - `WithTheme(theme teautils.Theme) ListModel[T]`
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

- `ModalType int`

- `MultiSelectButton struct{}`
  - Properties
    - `Hotkey rune`
    - `ID string`
    - `Label string`

- `MultiSelectButtonPressedMsg struct{}`
  - Properties
    - `ButtonID string`
    - `Selected []T`

- `MultiSelectCancelledMsg struct{}`

- `MultiSelectItem interface{}`
  - Methods
    - `ID() string`
    - `Label() string`

- `MultiSelectKeyMap struct{}`
  - Properties
    - `Cancel key.Binding`
    - `Confirm key.Binding`
    - `Down key.Binding`
    - `NextButton key.Binding`
    - `NextFocus key.Binding`
    - `PrevButton key.Binding`
    - `PrevFocus key.Binding`
    - `Toggle key.Binding`
    - `Up key.Binding`

- `MultiSelectModel struct{}`
  - Properties
    - `Keys MultiSelectKeyMap`
    - `Logger *slog.Logger`
  - Methods
    - `BorderStyle() lipgloss.Style`
    - `ButtonStyle() lipgloss.Style`
    - `CheckedStyle() lipgloss.Style`
    - `Close() (MultiSelectModel[T], tea.Cmd)`
    - `Cursor() int`
    - `FocusedButtonStyle() lipgloss.Style`
    - `FooterStyle() lipgloss.Style`
    - `Init() tea.Cmd`
    - `IsOpen() bool`
    - `ItemStyle() lipgloss.Style`
    - `MessageStyle() lipgloss.Style`
    - `Open() (MultiSelectModel[T], tea.Cmd)`
    - `OverlayModal(background string) (view string)`
    - `Selected() []T`
    - `SelectedItemStyle() lipgloss.Style`
    - `SetItems(items []T) MultiSelectModel[T]`
    - `SetSize(width int, height int) MultiSelectModel[T]`
    - `TitleStyle() lipgloss.Style`
    - `UncheckedStyle() lipgloss.Style`
    - `Update(msg tea.Msg) (MultiSelectModel[T], tea.Cmd)`
    - `View() tea.View`
    - `WithBorderStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithButtonStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithCheckedStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithFocusedButtonStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithFooter(footer string) MultiSelectModel[T]`
    - `WithFooterStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithItemStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithLogger(logger *slog.Logger) MultiSelectModel[T]`
    - `WithMaxVisible(max int) MultiSelectModel[T]`
    - `WithMessage(message string) MultiSelectModel[T]`
    - `WithMessageStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithSelectedItemStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithTheme(theme teautils.Theme) MultiSelectModel[T]`
    - `WithTitle(title string) MultiSelectModel[T]`
    - `WithTitleStyle(style lipgloss.Style) MultiSelectModel[T]`
    - `WithUncheckedStyle(style lipgloss.Style) MultiSelectModel[T]`

- `MultiSelectModelArgs struct{}`
  - Properties
    - `AllChecked bool`
    - `BorderStyle lipgloss.Style`
    - `ButtonStyle lipgloss.Style`
    - `Buttons []MultiSelectButton`
    - `CancelEmitsButtonPressed bool`
    - `CheckedStyle lipgloss.Style`
    - `FocusedButtonStyle lipgloss.Style`
    - `Footer string`
    - `FooterStyle lipgloss.Style`
    - `ItemStyle lipgloss.Style`
    - `MaxVisible int`
    - `Message string`
    - `MessageStyle lipgloss.Style`
    - `NoCancel bool`
    - `ScreenHeight int`
    - `ScreenWidth int`
    - `SelectedItemStyle lipgloss.Style`
    - `Title string`
    - `TitleStyle lipgloss.Style`
    - `UncheckedStyle lipgloss.Style`

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
    - `View() tea.View`

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

## Module: `./teamodal/examples/choices`

### Package: `choices`
- Path: `./teamodal/examples/choices`

## Module: `./teamodal/examples/editlist`

### Package: `editlist`
- Path: `./teamodal/examples/editlist`

#### Types

- `Task struct{}`
  - Methods
    - `ID() string`
    - `IsActive() bool`
    - `Label() string`
    - `String() string`

## Module: `./teamodal/examples/multiselect`

### Package: `multiselect`
- Path: `./teamodal/examples/multiselect`

#### Types

- `SelectableItem struct{}`
  - Methods
    - `ID() string`
    - `Label() string`

## Module: `./teamodal/examples/various`

### Package: `various`
- Path: `./teamodal/examples/various`

## Module: `./teamodal/examples/vertical`

### Package: `vertical`
- Path: `./teamodal/examples/vertical`

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
    - `DismissCmd() (cmd tea.Cmd)`
    - `HasActiveNotice() (active bool)`
    - `Init() (cmd tea.Cmd)`
    - `Initialize() (out NotifyModel, err error)`
    - `NewNotifyCmd(noticeType NoticeKey, message string) (cmd tea.Cmd)`
    - `RegisterNoticeType(def NoticeDefinition) (out NotifyModel, err error)`
    - `Render(content string) (result string)`
    - `Update(msg tea.Msg) (out NotifyModel, cmd tea.Cmd)`
    - `View() tea.View`
    - `WithAllowEscToClose() (out NotifyModel)`
    - `WithMinWidth(min int) (out NotifyModel)`
    - `WithPosition(pos Position) (out NotifyModel)`
    - `WithTheme(theme teautils.Theme) (out NotifyModel)`
    - `WithUnicodePrefix() (out NotifyModel)`

- `NotifyOpts struct{}`
  - Properties
    - `AllowEscToClose bool`
    - `CustomNotices []NoticeDefinition`
    - `Duration time.Duration`
    - `MinWidth int`
    - `NoAnimation bool`
    - `NoDefaultNotices bool`
    - `Position Position`
    - `UseNerdFont bool`
    - `UseUnicodePrefix bool`
    - `Width int`

- `Position string`

## Module: `./teanotify/examples/simple`

### Package: `simple`
- Path: `./teanotify/examples/simple`

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
    - `KeyText string`
    - `Label string`

- `MenuItemOpts struct{}`
  - Properties
    - `Label string`

- `SeparatorKind int`

- `SetIndicatorsMsg struct{}`
  - Properties
    - `Indicators []StatusIndicator`

- `SetMenuItemsMsg struct{}`
  - Properties
    - `Items []MenuItem`

- `StatusBarModel struct{}`
  - Properties
    - `Styles Styles`
  - Methods
    - `Init() tea.Cmd`
    - `SetIndicators(indicators []StatusIndicator) StatusBarModel`
    - `SetMenuItems(items []MenuItem) StatusBarModel`
    - `SetSize(width int) StatusBarModel`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() tea.View`
    - `WithStyles(styles Styles) StatusBarModel`
    - `WithTheme(theme teautils.Theme) StatusBarModel`

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

## Module: `./teastatus/examples/statusbar`

### Package: `statusbar`
- Path: `./teastatus/examples/statusbar`

## Module: `./teatree`

### Package: `teatree`
- Path: `./teatree`

#### Vars
- `ErrDrillDown = errors.New("drilldown error")`
- `ErrEmptyPath = errors.New("empty path error")`
- `ErrInvalidLevel = errors.New("invalid level error")`
- `ErrInvalidNode = errors.New("invalid node error")`
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

#### Funcs
- `DefaultDrillDownBorderStyle() lipgloss.Style`
- `DefaultDrillDownPathStyle() lipgloss.Style`
- `DefaultDrillDownSelectedStyle() lipgloss.Style`

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

- `DrillDownArgs struct{}`
  - Properties
    - `Prompt string`
    - `SelectorFunc DrillDownSelectorFunc[T]`

- `DrillDownChangeMsg struct{}`
  - Properties
    - `Level int`
    - `Node *Node[T]`

- `DrillDownFocusMsg struct{}`
  - Properties
    - `Level int`
    - `Node *Node[T]`

- `DrillDownKeyMap struct{}`
  - Properties
    - `Down key.Binding`
    - `OpenDropdown key.Binding`
    - `Select key.Binding`
    - `Up key.Binding`

- `DrillDownModel struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `Dropdown teafields.DropdownModel`
    - `DropdownOpen bool`
    - `Height int`
    - `InsertLine bool`
    - `Keys DrillDownKeyMap`
    - `Path []*Node[T]`
    - `PathStyle lipgloss.Style`
    - `Prompt string`
    - `SelectedLevel int`
    - `SelectedStyle lipgloss.Style`
    - `SelectorFunc DrillDownSelectorFunc[T]`
    - `Width int`
  - Methods
    - `Init() tea.Cmd`
    - `Initialize() (model DrillDownModel[T], err error)`
    - `Root() *Node[T]`
    - `Update(msg tea.Msg) (tea.Model, tea.Cmd)`
    - `View() tea.View`
    - `WithBorderStyle(style lipgloss.Style) DrillDownModel[T]`
    - `WithDrillDownTheme(theme teautils.Theme) DrillDownModel[T]`
    - `WithPathStyle(style lipgloss.Style) DrillDownModel[T]`
    - `WithSelectedStyle(style lipgloss.Style) DrillDownModel[T]`

- `DrillDownSelectMsg struct{}`
  - Properties
    - `Node *Node[T]`

- `DrillDownSelectorFunc func(current *Node[T], children []*Node[T]) (*Node[T], error)`

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

- `TreeModel struct{}`
  - Properties
    - `Keys TreeKeyMap`
  - Methods
    - `FocusedNode() (node *Node[T])`
    - `Init() tea.Cmd`
    - `MaxLineWidth() int`
    - `RequiredWidth() int`
    - `SetBorderColor(c color.Color) TreeModel[T]`
    - `SetFocusedNode(nodeID string) TreeModel[T]`
    - `SetSize(width int, height int) TreeModel[T]`
    - `Theme() *teautils.Theme`
    - `Tree() *Tree[T]`
    - `Update(msg tea.Msg) (TreeModel[T], tea.Cmd)`
    - `View() (view tea.View)`
    - `WithFrame(style lipgloss.Style) TreeModel[T]`
    - `WithHeader(text string, style lipgloss.Style) TreeModel[T]`
    - `WithTheme(theme teautils.Theme) TreeModel[T]`

- `TreeOption func(*Tree[T])`

## Module: `./teatree/examples/drilldown`

### Package: `drilldown`
- Path: `./teatree/examples/drilldown`

## Module: `./teatree/examples/filetree`

### Package: `filetree`
- Path: `./teatree/examples/filetree`

## Module: `./teatxtsnip`

### Package: `teatxtsnip`
- Path: `./teatxtsnip`

#### Vars
- `ErrInvalidSelection = errors.New("invalid selection")`
- `ErrTextSnip = errors.New("text snip")`
- `SelectionStyle = lipgloss.NewStyle().
	Background(lipgloss.Color("39")).
	Foreground(lipgloss.Color("232"))`

#### Types

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

- `TextSnipModel struct{}`
  - Properties
    - `Logger *slog.Logger`
    - `textarea.Model`
  - Methods
    - `ClearSelection() TextSnipModel`
    - `Copy() TextSnipModel`
    - `Cut() TextSnipModel`
    - `HasSelection() bool`
    - `Init() tea.Cmd`
    - `IsSingleLine() bool`
    - `Paste() (TextSnipModel, tea.Cmd)`
    - `SelectedText() string`
    - `Selection() Selection`
    - `SelectionKeyMap() SelectionKeyMap`
    - `SetLogger(logger *slog.Logger) TextSnipModel`
    - `SetSelection(sel Selection) TextSnipModel`
    - `SetSelectionKeyMap(km SelectionKeyMap) TextSnipModel`
    - `Update(msg tea.Msg) (TextSnipModel, tea.Cmd)`
    - `View() (view tea.View)`
    - `WithTheme(theme teautils.Theme) TextSnipModel`

- `TextSnipModelArgs struct{}`
  - Properties
    - `SingleLine bool`
    - `Textarea *textarea.Model`

## Module: `./teatxtsnip/examples/editor`

### Package: `editor`
- Path: `./teatxtsnip/examples/editor`

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
- `SemanticBlack = NewSemanticColor(teacolor.Black)`
- `SemanticBlue = NewSemanticColor(teacolor.Blue)`
- `SemanticBrightBlack = NewSemanticColor(teacolor.BrightBlack)`
- `SemanticBrightBlue = NewSemanticColor(teacolor.BrightBlue)`
- `SemanticBrightCyan = NewSemanticColor(teacolor.BrightCyan)`
- `SemanticBrightGreen = NewSemanticColor(teacolor.BrightGreen)`
- `SemanticBrightMagenta = NewSemanticColor(teacolor.BrightMagenta)`
- `SemanticBrightRed = NewSemanticColor(teacolor.BrightRed)`
- `SemanticBrightWhite = NewSemanticColor(teacolor.BrightWhite)`
- `SemanticBrightYellow = NewSemanticColor(teacolor.BrightYellow)`
- `SemanticColor0 = NewSemanticColor(teacolor.Color0)`
- `SemanticColor124 = NewSemanticColor(teacolor.Color124)`
- `SemanticColor130 = NewSemanticColor(teacolor.Color130)`
- `SemanticColor15 = NewSemanticColor(teacolor.Color15)`
- `SemanticColor153 = NewSemanticColor(teacolor.Color153)`
- `SemanticColor157 = NewSemanticColor(teacolor.Color157)`
- `SemanticColor160 = NewSemanticColor(teacolor.Color160)`
- `SemanticColor166 = NewSemanticColor(teacolor.Color166)`
- `SemanticColor178 = NewSemanticColor(teacolor.Color178)`
- `SemanticColor214 = NewSemanticColor(teacolor.Color214)`
- `SemanticColor217 = NewSemanticColor(teacolor.Color217)`
- `SemanticColor22 = NewSemanticColor(teacolor.Color22)`
- `SemanticColor226 = NewSemanticColor(teacolor.Color226)`
- `SemanticColor230 = NewSemanticColor(teacolor.Color230)`
- `SemanticColor232 = NewSemanticColor(teacolor.Color232)`
- `SemanticColor238 = NewSemanticColor(teacolor.Color238)`
- `SemanticColor240 = NewSemanticColor(teacolor.Color240)`
- `SemanticColor244 = NewSemanticColor(teacolor.Color244)`
- `SemanticColor248 = NewSemanticColor(teacolor.Color248)`
- `SemanticColor250 = NewSemanticColor(teacolor.Color250)`
- `SemanticColor252 = NewSemanticColor(teacolor.Color252)`
- `SemanticColor28 = NewSemanticColor(teacolor.Color28)`
- `SemanticColor30 = NewSemanticColor(teacolor.Color30)`
- `SemanticColor33 = NewSemanticColor(teacolor.Color33)`
- `SemanticColor39 = NewSemanticColor(teacolor.Color39)`
- `SemanticColor46 = NewSemanticColor(teacolor.Color46)`
- `SemanticColor51 = NewSemanticColor(teacolor.Color51)`
- `SemanticColor52 = NewSemanticColor(teacolor.Color52)`
- `SemanticColor62 = NewSemanticColor(teacolor.Color62)`
- `SemanticColor86 = NewSemanticColor(teacolor.Color86)`
- `SemanticColorNil = NewSemanticColor(nil)`
- `SemanticCyan = NewSemanticColor(teacolor.Cyan)`
- `SemanticGreen = NewSemanticColor(teacolor.Green)`
- `SemanticMagenta = NewSemanticColor(teacolor.Magenta)`
- `SemanticRed = NewSemanticColor(teacolor.Red)`
- `SemanticWhite = NewSemanticColor(teacolor.White)`
- `SemanticYellow = NewSemanticColor(teacolor.Yellow)`

#### Funcs
- `ApplyBoxBorder(borderStyle lipgloss.Style, content string) string`
- `CalculateCenter(screenW int, screenH int, modalW int, modalH int) (row int, col int)`
- `CenterModal(renderedView string, screenW int, screenH int) (width int, height int, row int, col int)`
- `FormatKeyDisplay(k KeyMeta) string`
- `GetSortedCategories(keysByCategory map[string][]KeyMeta, preferredOrder []string) []string`
- `IsJediTerm() (isJedi bool)`
- `MeasureRenderedView(renderedView string) (width int, height int)`
- `ProperCaseShortcut(s string) string`
- `RenderAlignedLine(text string, style lipgloss.Style, width int, align lipgloss.Position) string`
- `RenderCenteredLine(text string, style lipgloss.Style, width int) string`

#### Types

- `BreadcrumbTheme struct{}`
  - Properties
    - `CurrentStyle lipgloss.Style`
    - `HoverStyle lipgloss.Style`
    - `ParentStyle lipgloss.Style`
    - `SeparatorStyle lipgloss.Style`

- `DropdownTheme struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `ItemStyle lipgloss.Style`
    - `SelectedStyle lipgloss.Style`

- `GridTheme struct{}`
  - Properties
    - `BaseStyle lipgloss.Style`
    - `BorderStyle lipgloss.Style`
    - `HeaderStyle lipgloss.Style`
    - `HighlightStyle lipgloss.Style`

- `HelpVisorStyle struct{}`
  - Properties
    - `CategoryOrder []string`
    - `CategoryStyle lipgloss.Style`
    - `DescStyle lipgloss.Style`
    - `KeyColumnGap int`
    - `KeyStyle lipgloss.Style`
    - `TitleStyle lipgloss.Style`

- `HelpVisorTheme struct{}`
  - Properties
    - `CategoryStyle lipgloss.Style`
    - `DescStyle lipgloss.Style`
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

- `ListTheme struct{}`
  - Properties
    - `ActiveItemStyle lipgloss.Style`
    - `EditItemStyle lipgloss.Style`
    - `FooterStyle lipgloss.Style`
    - `ItemStyle lipgloss.Style`
    - `ScrollThumbStyle lipgloss.Style`
    - `ScrollbarStyle lipgloss.Style`
    - `SelectedItemStyle lipgloss.Style`
    - `StatusStyle lipgloss.Style`

- `ModalTheme struct{}`
  - Properties
    - `BorderStyle lipgloss.Style`
    - `ButtonStyle lipgloss.Style`
    - `CancelKeyStyle lipgloss.Style`
    - `CancelTextStyle lipgloss.Style`
    - `FocusedButtonStyle lipgloss.Style`
    - `MessageStyle lipgloss.Style`
    - `TitleStyle lipgloss.Style`

- `Palette struct{}`
  - Properties
    - `App T`
    - `System SystemPalette`

- `PaletteOpts struct{}`
  - Properties
    - `Adaptive bool`

- `SemanticColor struct{}`
  - Methods
    - `Background() lipgloss.Style`
    - `BorderForeground() lipgloss.Style`
    - `Color() color.Color`
    - `Foreground() lipgloss.Style`
    - `IsZero() bool`
    - `RGBA() (r uint32, g uint32, b uint32, a uint32)`
    - `Render(text string) string`

- `StatusBarTheme struct{}`
  - Properties
    - `BarStyle lipgloss.Style`
    - `IndicatorSepStyle lipgloss.Style`
    - `IndicatorStyle lipgloss.Style`
    - `MenuKeyStyle lipgloss.Style`
    - `MenuLabelStyle lipgloss.Style`

- `SystemPalette struct{}`
  - Properties
    - `Accent SemanticColor`
    - `AccentAlt SemanticColor`
    - `AccentSubtle SemanticColor`
    - `Border SemanticColor`
    - `BorderAccent SemanticColor`
    - `ButtonBg SemanticColor`
    - `ButtonFg SemanticColor`
    - `ButtonFocusBg SemanticColor`
    - `ButtonFocusFg SemanticColor`
    - `EditBg SemanticColor`
    - `EditFg SemanticColor`
    - `FocusBg SemanticColor`
    - `FocusBorder SemanticColor`
    - `HighlightStyle string`
    - `ScrollThumb SemanticColor`
    - `ScrollTrack SemanticColor`
    - `SelectionBg SemanticColor`
    - `SelectionFg SemanticColor`
    - `Separator SemanticColor`
    - `StatusError SemanticColor`
    - `StatusInfo SemanticColor`
    - `StatusSuccess SemanticColor`
    - `StatusWarn SemanticColor`
    - `TextDim SemanticColor`
    - `TextMuted SemanticColor`
    - `TextPrimary SemanticColor`
    - `TextSecondary SemanticColor`
    - `TintNegative SemanticColor`
    - `TintNeutral SemanticColor`
    - `TintPositive SemanticColor`

- `Theme struct{}`
  - Properties
    - `ActiveItem lipgloss.Style`
    - `Border lipgloss.Style`
    - `BorderAccent lipgloss.Style`
    - `Breadcrumb BreadcrumbTheme`
    - `Button lipgloss.Style`
    - `Dropdown DropdownTheme`
    - `FocusedButton lipgloss.Style`
    - `Grid GridTheme`
    - `HelpVisor HelpVisorTheme`
    - `Item lipgloss.Style`
    - `List ListTheme`
    - `Message lipgloss.Style`
    - `Modal ModalTheme`
    - `SelectedItem lipgloss.Style`
    - `StatusBar StatusBarTheme`
    - `System SystemPalette`
    - `Title lipgloss.Style`

### Package: `teacolor`
- Path: `./teautils/teacolor`

#### Vars
- `Amber color.Color = Color("#FFBF00")`
- `Aquamarine color.Color = Color("#7FFFD4")`
- `Beige color.Color = Color("#F5F5DC")`
- `Black color.Color = Color("0")`
- `Blue color.Color = Color("4")`
- `Blush color.Color = Color("#DE5D83")`
- `BrightBlack color.Color = Color("8")`
- `BrightBlue color.Color = Color("12")`
- `BrightCyan color.Color = Color("14")`
- `BrightGreen color.Color = Color("10")`
- `BrightMagenta color.Color = Color("13")`
- `BrightRed color.Color = Color("9")`
- `BrightWhite color.Color = Color("15")`
- `BrightYellow color.Color = Color("11")`
- `CharcoalGray color.Color = Color("236")`
- `Color0 color.Color = Color("0")`
- `Color1 color.Color = Color("1")`
- `Color10 color.Color = Color("10")`
- `Color100 color.Color = Color("100")`
- `Color101 color.Color = Color("101")`
- `Color102 color.Color = Color("102")`
- `Color103 color.Color = Color("103")`
- `Color104 color.Color = Color("104")`
- `Color105 color.Color = Color("105")`
- `Color106 color.Color = Color("106")`
- `Color107 color.Color = Color("107")`
- `Color108 color.Color = Color("108")`
- `Color109 color.Color = Color("109")`
- `Color11 color.Color = Color("11")`
- `Color110 color.Color = Color("110")`
- `Color111 color.Color = Color("111")`
- `Color112 color.Color = Color("112")`
- `Color113 color.Color = Color("113")`
- `Color114 color.Color = Color("114")`
- `Color115 color.Color = Color("115")`
- `Color116 color.Color = Color("116")`
- `Color117 color.Color = Color("117")`
- `Color118 color.Color = Color("118")`
- `Color119 color.Color = Color("119")`
- `Color12 color.Color = Color("12")`
- `Color120 color.Color = Color("120")`
- `Color121 color.Color = Color("121")`
- `Color122 color.Color = Color("122")`
- `Color123 color.Color = Color("123")`
- `Color124 color.Color = Color("124")`
- `Color125 color.Color = Color("125")`
- `Color126 color.Color = Color("126")`
- `Color127 color.Color = Color("127")`
- `Color128 color.Color = Color("128")`
- `Color129 color.Color = Color("129")`
- `Color13 color.Color = Color("13")`
- `Color130 color.Color = Color("130")`
- `Color131 color.Color = Color("131")`
- `Color132 color.Color = Color("132")`
- `Color133 color.Color = Color("133")`
- `Color134 color.Color = Color("134")`
- `Color135 color.Color = Color("135")`
- `Color136 color.Color = Color("136")`
- `Color137 color.Color = Color("137")`
- `Color138 color.Color = Color("138")`
- `Color139 color.Color = Color("139")`
- `Color14 color.Color = Color("14")`
- `Color140 color.Color = Color("140")`
- `Color141 color.Color = Color("141")`
- `Color142 color.Color = Color("142")`
- `Color143 color.Color = Color("143")`
- `Color144 color.Color = Color("144")`
- `Color145 color.Color = Color("145")`
- `Color146 color.Color = Color("146")`
- `Color147 color.Color = Color("147")`
- `Color148 color.Color = Color("148")`
- `Color149 color.Color = Color("149")`
- `Color15 color.Color = Color("15")`
- `Color150 color.Color = Color("150")`
- `Color151 color.Color = Color("151")`
- `Color152 color.Color = Color("152")`
- `Color153 color.Color = Color("153")`
- `Color154 color.Color = Color("154")`
- `Color155 color.Color = Color("155")`
- `Color156 color.Color = Color("156")`
- `Color157 color.Color = Color("157")`
- `Color158 color.Color = Color("158")`
- `Color159 color.Color = Color("159")`
- `Color16 color.Color = Color("16")`
- `Color160 color.Color = Color("160")`
- `Color161 color.Color = Color("161")`
- `Color162 color.Color = Color("162")`
- `Color163 color.Color = Color("163")`
- `Color164 color.Color = Color("164")`
- `Color165 color.Color = Color("165")`
- `Color166 color.Color = Color("166")`
- `Color167 color.Color = Color("167")`
- `Color168 color.Color = Color("168")`
- `Color169 color.Color = Color("169")`
- `Color17 color.Color = Color("17")`
- `Color170 color.Color = Color("170")`
- `Color171 color.Color = Color("171")`
- `Color172 color.Color = Color("172")`
- `Color173 color.Color = Color("173")`
- `Color174 color.Color = Color("174")`
- `Color175 color.Color = Color("175")`
- `Color176 color.Color = Color("176")`
- `Color177 color.Color = Color("177")`
- `Color178 color.Color = Color("178")`
- `Color179 color.Color = Color("179")`
- `Color18 color.Color = Color("18")`
- `Color180 color.Color = Color("180")`
- `Color181 color.Color = Color("181")`
- `Color182 color.Color = Color("182")`
- `Color183 color.Color = Color("183")`
- `Color184 color.Color = Color("184")`
- `Color185 color.Color = Color("185")`
- `Color186 color.Color = Color("186")`
- `Color187 color.Color = Color("187")`
- `Color188 color.Color = Color("188")`
- `Color189 color.Color = Color("189")`
- `Color19 color.Color = Color("19")`
- `Color190 color.Color = Color("190")`
- `Color191 color.Color = Color("191")`
- `Color192 color.Color = Color("192")`
- `Color193 color.Color = Color("193")`
- `Color194 color.Color = Color("194")`
- `Color195 color.Color = Color("195")`
- `Color196 color.Color = Color("196")`
- `Color197 color.Color = Color("197")`
- `Color198 color.Color = Color("198")`
- `Color199 color.Color = Color("199")`
- `Color2 color.Color = Color("2")`
- `Color20 color.Color = Color("20")`
- `Color200 color.Color = Color("200")`
- `Color201 color.Color = Color("201")`
- `Color202 color.Color = Color("202")`
- `Color203 color.Color = Color("203")`
- `Color204 color.Color = Color("204")`
- `Color205 color.Color = Color("205")`
- `Color206 color.Color = Color("206")`
- `Color207 color.Color = Color("207")`
- `Color208 color.Color = Color("208")`
- `Color209 color.Color = Color("209")`
- `Color21 color.Color = Color("21")`
- `Color210 color.Color = Color("210")`
- `Color211 color.Color = Color("211")`
- `Color212 color.Color = Color("212")`
- `Color213 color.Color = Color("213")`
- `Color214 color.Color = Color("214")`
- `Color215 color.Color = Color("215")`
- `Color216 color.Color = Color("216")`
- `Color217 color.Color = Color("217")`
- `Color218 color.Color = Color("218")`
- `Color219 color.Color = Color("219")`
- `Color22 color.Color = Color("22")`
- `Color220 color.Color = Color("220")`
- `Color221 color.Color = Color("221")`
- `Color222 color.Color = Color("222")`
- `Color223 color.Color = Color("223")`
- `Color224 color.Color = Color("224")`
- `Color225 color.Color = Color("225")`
- `Color226 color.Color = Color("226")`
- `Color227 color.Color = Color("227")`
- `Color228 color.Color = Color("228")`
- `Color229 color.Color = Color("229")`
- `Color23 color.Color = Color("23")`
- `Color230 color.Color = Color("230")`
- `Color231 color.Color = Color("231")`
- `Color232 color.Color = Color("232")`
- `Color233 color.Color = Color("233")`
- `Color234 color.Color = Color("234")`
- `Color235 color.Color = Color("235")`
- `Color236 color.Color = Color("236")`
- `Color237 color.Color = Color("237")`
- `Color238 color.Color = Color("238")`
- `Color239 color.Color = Color("239")`
- `Color24 color.Color = Color("24")`
- `Color240 color.Color = Color("240")`
- `Color241 color.Color = Color("241")`
- `Color242 color.Color = Color("242")`
- `Color243 color.Color = Color("243")`
- `Color244 color.Color = Color("244")`
- `Color245 color.Color = Color("245")`
- `Color246 color.Color = Color("246")`
- `Color247 color.Color = Color("247")`
- `Color248 color.Color = Color("248")`
- `Color249 color.Color = Color("249")`
- `Color25 color.Color = Color("25")`
- `Color250 color.Color = Color("250")`
- `Color251 color.Color = Color("251")`
- `Color252 color.Color = Color("252")`
- `Color253 color.Color = Color("253")`
- `Color254 color.Color = Color("254")`
- `Color255 color.Color = Color("255")`
- `Color26 color.Color = Color("26")`
- `Color27 color.Color = Color("27")`
- `Color28 color.Color = Color("28")`
- `Color29 color.Color = Color("29")`
- `Color3 color.Color = Color("3")`
- `Color30 color.Color = Color("30")`
- `Color31 color.Color = Color("31")`
- `Color32 color.Color = Color("32")`
- `Color33 color.Color = Color("33")`
- `Color34 color.Color = Color("34")`
- `Color35 color.Color = Color("35")`
- `Color36 color.Color = Color("36")`
- `Color37 color.Color = Color("37")`
- `Color38 color.Color = Color("38")`
- `Color39 color.Color = Color("39")`
- `Color4 color.Color = Color("4")`
- `Color40 color.Color = Color("40")`
- `Color41 color.Color = Color("41")`
- `Color42 color.Color = Color("42")`
- `Color43 color.Color = Color("43")`
- `Color44 color.Color = Color("44")`
- `Color45 color.Color = Color("45")`
- `Color46 color.Color = Color("46")`
- `Color47 color.Color = Color("47")`
- `Color48 color.Color = Color("48")`
- `Color49 color.Color = Color("49")`
- `Color5 color.Color = Color("5")`
- `Color50 color.Color = Color("50")`
- `Color51 color.Color = Color("51")`
- `Color52 color.Color = Color("52")`
- `Color53 color.Color = Color("53")`
- `Color54 color.Color = Color("54")`
- `Color55 color.Color = Color("55")`
- `Color56 color.Color = Color("56")`
- `Color57 color.Color = Color("57")`
- `Color58 color.Color = Color("58")`
- `Color59 color.Color = Color("59")`
- `Color6 color.Color = Color("6")`
- `Color60 color.Color = Color("60")`
- `Color61 color.Color = Color("61")`
- `Color62 color.Color = Color("62")`
- `Color63 color.Color = Color("63")`
- `Color64 color.Color = Color("64")`
- `Color65 color.Color = Color("65")`
- `Color66 color.Color = Color("66")`
- `Color67 color.Color = Color("67")`
- `Color68 color.Color = Color("68")`
- `Color69 color.Color = Color("69")`
- `Color7 color.Color = Color("7")`
- `Color70 color.Color = Color("70")`
- `Color71 color.Color = Color("71")`
- `Color72 color.Color = Color("72")`
- `Color73 color.Color = Color("73")`
- `Color74 color.Color = Color("74")`
- `Color75 color.Color = Color("75")`
- `Color76 color.Color = Color("76")`
- `Color77 color.Color = Color("77")`
- `Color78 color.Color = Color("78")`
- `Color79 color.Color = Color("79")`
- `Color8 color.Color = Color("8")`
- `Color80 color.Color = Color("80")`
- `Color81 color.Color = Color("81")`
- `Color82 color.Color = Color("82")`
- `Color83 color.Color = Color("83")`
- `Color84 color.Color = Color("84")`
- `Color85 color.Color = Color("85")`
- `Color86 color.Color = Color("86")`
- `Color87 color.Color = Color("87")`
- `Color88 color.Color = Color("88")`
- `Color89 color.Color = Color("89")`
- `Color9 color.Color = Color("9")`
- `Color90 color.Color = Color("90")`
- `Color91 color.Color = Color("91")`
- `Color92 color.Color = Color("92")`
- `Color93 color.Color = Color("93")`
- `Color94 color.Color = Color("94")`
- `Color95 color.Color = Color("95")`
- `Color96 color.Color = Color("96")`
- `Color97 color.Color = Color("97")`
- `Color98 color.Color = Color("98")`
- `Color99 color.Color = Color("99")`
- `Coral color.Color = Color("#FF7F50")`
- `CornflowerBlue color.Color = Color("#6495ED")`
- `Crimson color.Color = Color("#DC143C")`
- `Cyan color.Color = Color("6")`
- `DarkGray color.Color = Color("240")`
- `DarkOrange color.Color = Color("#FF8C00")`
- `DarkRed color.Color = Color("#8B0000")`
- `DimGray color.Color = Color("238")`
- `DodgerBlue color.Color = Color("#1E90FF")`
- `Emerald color.Color = Color("#50C878")`
- `FireBrick color.Color = Color("#B22222")`
- `ForestGreen color.Color = Color("#228B22")`
- `Gold color.Color = Color("#FFD700")`
- `Gray color.Color = Color("245")`
- `Green color.Color = Color("2")`
- `HotPink color.Color = Color("#FF69B4")`
- `Indigo color.Color = Color("#4B0082")`
- `Ivory color.Color = Color("#FFFFF0")`
- `Khaki color.Color = Color("#F0E68C")`
- `Lavender color.Color = Color("#E6E6FA")`
- `LightGray color.Color = Color("252")`
- `LightPink color.Color = Color("#FFB6C1")`
- `Lime color.Color = Color("#00FF00")`
- `Magenta color.Color = Color("5")`
- `Mint color.Color = Color("#98FF98")`
- `Navy color.Color = Color("#000080")`
- `Olive color.Color = Color("#808000")`
- `Orange color.Color = Color("#FFA500")`
- `Orchid color.Color = Color("#DA70D6")`
- `PaleTurquoise color.Color = Color("#AFEEEE")`
- `Peach color.Color = Color("#FFDAB9")`
- `Plum color.Color = Color("#DDA0DD")`
- `Purple color.Color = Color("#800080")`
- `Red color.Color = Color("1")`
- `Rose color.Color = Color("#FF007F")`
- `RoyalBlue color.Color = Color("#4169E1")`
- `Salmon color.Color = Color("#FA8072")`
- `SeaGreen color.Color = Color("#2E8B57")`
- `Silver color.Color = Color("#C0C0C0")`
- `SilverGray color.Color = Color("250")`
- `SkyBlue color.Color = Color("#87CEEB")`
- `SlateGray color.Color = Color("#708090")`
- `SpringGreen color.Color = Color("#00FF7F")`
- `SteelBlue color.Color = Color("#4682B4")`
- `Tangerine color.Color = Color("#FF9966")`
- `Teal color.Color = Color("#008080")`
- `Tomato color.Color = Color("#FF6347")`
- `TrueBlack color.Color = Color("#000000")`
- `TrueBlue color.Color = Color("#0000FF")`
- `TrueCyan color.Color = Color("#00FFFF")`
- `TrueGreen color.Color = Color("#00FF00")`
- `TrueMagenta color.Color = Color("#FF00FF")`
- `TrueRed color.Color = Color("#FF0000")`
- `TrueWhite color.Color = Color("#FFFFFF")`
- `TrueYellow color.Color = Color("#FFFF00")`
- `Turquoise color.Color = Color("#40E0D0")`
- `Violet color.Color = Color("#EE82EE")`
- `White color.Color = Color("7")`
- `Yellow color.Color = Color("3")`

#### Funcs
- `Color(s string) color.Color`

## Module: `./teautils/examples/keyhelp`

### Package: `keyhelp`
- Path: `./teautils/examples/keyhelp`

## Module: `./teautils/examples/theming`

### Package: `theming`
- Path: `./teautils/examples/theming`

