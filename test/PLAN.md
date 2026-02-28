# Test Plan: go-tealeaves v1 Comprehensive Test Suite

## Strategy

Three-layer testing approach covering logic, rendering, and full program lifecycle.

### Layer 1: Direct Model Unit Tests

Test `Update()` state transitions, message production, and constructor/option behavior by calling model methods directly with constructed `tea.Msg` values. No framework beyond `testing`.

**What this catches:** Logic bugs, key handling regressions, state machine errors, message contract violations.

### Layer 2: View() Output Assertions

Call `View()` on models in known states and assert on the rendered string. Check border geometry, alignment, content placement, and overlay positioning. Use string containment, line counting, and ANSI-aware width measurement.

**What this catches:** Rendering regressions — borders too wide/narrow, incorrect indentation, misaligned content, broken overlay compositing.

### Layer 3: teatest Program Tests

Use `github.com/charmbracelet/x/exp/teatest` to spin up a real `tea.Program`, send messages through the full event loop, and capture rendered output. Use `RequireEqualOutput()` with golden files for regression detection.

**What this catches:** Full lifecycle rendering bugs, Init/Update/View integration, timing-dependent behavior, rendering pipeline issues that only manifest through the real renderer.

### Golden File Convention

- Golden files live in `<module>/testdata/` directories
- Named `<TestName>.golden`
- Updated via `go test -update` flag (teatest convention)
- Committed to version control

### V2 Migration Path

- Tests written now use `github.com/charmbracelet/x/exp/teatest` (Charm v1)
- During MIGRATE_V2 phase, swap to `github.com/charmbracelet/x/exp/teatest/v2`
- Golden files will be regenerated after v2 migration to capture new rendering baseline
- Test logic (state transitions, message contracts) should survive v2 migration with minimal changes to key message construction

---

## Module Test Specifications

### TEAUTILS — `teautils/*_test.go`

Pure utility functions — Layer 1 only (no models, no teatest needed).

#### `positioning_test.go`

| Test | What It Verifies |
|---|---|
| `TestCalculateCenter` | Centers modal in screen; clamps negative to 0; handles modal larger than screen |
| `TestCalculateCenter_OddDimensions` | Correct rounding for odd width/height |
| `TestMeasureRenderedView` | Returns correct width/height for multi-line plain text |
| `TestMeasureRenderedView_ANSI` | ANSI escape sequences don't inflate width measurement |
| `TestMeasureRenderedView_Empty` | Returns 0,0 for empty string |
| `TestMeasureRenderedView_SingleLine` | Width correct, height=1 |
| `TestCenterModal` | Combines measure+center; verifies shift-up-by-1 behavior |

#### `render_styled_test.go`

| Test | What It Verifies |
|---|---|
| `TestRenderAlignedLine_Left` | Text aligned left within width |
| `TestRenderAlignedLine_Center` | Text centered within width |
| `TestRenderAlignedLine_Right` | Text right-aligned within width |
| `TestRenderCenteredLine` | Convenience wrapper matches RenderAlignedLine center |
| `TestApplyBoxBorder` | Output contains border characters; content is inside; padding applied |
| `TestApplyBoxBorder_MultiLine` | Multi-line content renders correctly inside border |

#### `render_help_visor_test.go`

| Test | What It Verifies |
|---|---|
| `TestFormatKeyDisplay_SingleKey` | Single binding renders as-is |
| `TestFormatKeyDisplay_MultipleKeys` | Multiple keys joined with "/" |
| `TestFormatKeyDisplay_Deduplication` | Duplicate keys collapsed |
| `TestFormatKeyDisplay_CustomDisplayKeys` | DisplayKeys override binding keys |
| `TestProperCaseShortcut` | "ctrl+c" → "Ctrl+C", "esc" → "Esc", "shift+tab" → "Shift+Tab" |
| `TestProperCaseShortcut_AlreadyProper` | Already-proper-cased input unchanged |
| `TestGetSortedCategories_PreferredOrder` | Preferred order respected; others sorted alphabetically after |
| `TestGetSortedCategories_NoPreference` | All categories sorted alphabetically |
| `TestDefaultHelpVisorStyle` | Returns non-zero styles (smoke check) |

#### `key_identifier_test.go`

| Test | What It Verifies |
|---|---|
| `TestParseKeyIdentifier_Valid` | "app.help", "tree.nav.up" accepted |
| `TestParseKeyIdentifier_Empty` | Returns ErrEmptyKeyIdentifier |
| `TestParseKeyIdentifier_NoDot` | Returns ErrKeyIdentifierMissingDot |
| `TestParseKeyIdentifier_EmptyPart` | "app..help" returns ErrKeyIdentifierEmptyPart |
| `TestParseKeyIdentifier_InvalidPart` | Invalid characters return ErrKeyIdentifierInvalidPart |
| `TestMustParseKeyIdentifier_Panics` | Panics on invalid input |
| `TestKeyIdentifier_String` | String() returns original input |

#### `key_registry_test.go`

| Test | What It Verifies |
|---|---|
| `TestKeyRegistry_Register` | Key registered and retrievable via Get |
| `TestKeyRegistry_Register_DefaultsHelpText` | Empty HelpText defaults to binding help description |
| `TestKeyRegistry_Register_EmptyID` | Returns ErrEmptyKeyID |
| `TestKeyRegistry_RegisterMany` | Multiple keys registered; all retrievable |
| `TestKeyRegistry_Get_NotFound` | Returns ErrKeyNotFound |
| `TestKeyRegistry_Clear` | After clear, previously registered keys return ErrKeyNotFound |
| `TestKeyRegistry_ForStatusBar` | Returns only keys with StatusBar=true, in registration order |
| `TestKeyRegistry_ForHelpModal` | Returns only keys with HelpModal=true |
| `TestKeyRegistry_ByCategory` | Groups keys by Category; uncategorized go to "Other"; preserves order within category |

---

### TEADD — `teadd/*_test.go`

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewModel_Defaults` | 1 | Default styles applied; IsOpen=false; Selected=0 |
| `TestNewModel_WithOptions` | 1 | Options stored; count matches input |
| `TestNewModel_EmptyOptions` | 1 | Handles empty option slice gracefully |
| `TestDropdownModel_Open` | 1 | IsOpen becomes true; position calculated |
| `TestDropdownModel_Close` | 1 | IsOpen becomes false |
| `TestDropdownModel_UpdateWhenClosed` | 1 | Returns model unchanged, nil cmd |
| `TestDropdownModel_KeyUp` | 1 | Selected decrements; stops at 0; ScrollOffset adjusts |
| `TestDropdownModel_KeyDown` | 1 | Selected increments; stops at len-1 |
| `TestDropdownModel_KeySelect` | 1 | Returns OptionSelectedMsg with correct Index, Text, Value; closes dropdown |
| `TestDropdownModel_KeyCancel` | 1 | Returns DropdownCancelledMsg; closes dropdown |
| `TestDropdownModel_WindowSizeMsg` | 1 | ScreenWidth/ScreenHeight updated |
| `TestDropdownModel_ScrollOffset` | 1 | Scrolling works when options exceed visible area |
| `TestDropdownModel_View_Closed` | 2 | Returns empty string when closed |
| `TestDropdownModel_View_Open` | 2 | Non-empty output; contains option text; selected item visually distinct |
| `TestDropdownModel_WithPosition` | 1 | FieldRow/FieldCol updated |
| `TestDropdownModel_WithOptions` | 1 | Options replaced; Selected reset to 0 |
| `TestDropdownModel_WithScreenSize` | 1 | ScreenWidth/ScreenHeight set |

#### `overlay_dropdown_test.go` — Layer 2

| Test | What It Verifies |
|---|---|
| `TestOverlayDropdown_BasicOverlay` | Foreground composited at correct row/col |
| `TestOverlayDropdown_AtOrigin` | Row=0, Col=0 overlay works |
| `TestOverlayDropdown_OutOfBounds` | Overlay beyond background dimensions handled gracefully |
| `TestOverlayDropdown_EmptyForeground` | Returns background unchanged |
| `TestOverlayDropdown_ANSIContent` | ANSI-styled foreground rendered correctly |

#### `types_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestToOptions` | Converts strings; Text and Value match |
| `TestToOptions_Empty` | Empty input returns empty slice |

#### `teatest_test.go` — Layer 3

| Test | What It Verifies |
|---|---|
| `TestDropdownModel_FullLifecycle` | Open → navigate → select → verify final model via teatest |
| `TestDropdownModel_RenderGolden` | Golden file comparison of rendered dropdown at various states |

---

### TEAMODAL — `teamodal/*_test.go`

Existing `choice_model_test.go` covers ChoiceModel Layer 1. Extend with remaining models and layers.

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewOKModal` | 1 | Type is ModalTypeOK; message stored; not open initially |
| `TestNewYesNoModal` | 1 | Type is ModalTypeYesNo; focusButton starts at 0 (Yes) |
| `TestModalModel_Open` | 1 | IsOpen becomes true |
| `TestModalModel_Close` | 1 | IsOpen becomes false |
| `TestModalModel_SetSize` | 1 | ScreenWidth/ScreenHeight updated |
| `TestOKModal_EnterClosesAndSendsClosedMsg` | 1 | Enter key → ClosedMsg, modal closed |
| `TestOKModal_EscClosesAndSendsClosedMsg` | 1 | Esc key → ClosedMsg, modal closed |
| `TestYesNoModal_TabTogglesFocus` | 1 | Tab alternates between Yes(0) and No(1) |
| `TestYesNoModal_EnterOnYes` | 1 | Enter when focus=0 → AnsweredYesMsg |
| `TestYesNoModal_EnterOnNo` | 1 | Enter when focus=1 → AnsweredNoMsg |
| `TestYesNoModal_EscSendsAnsweredNo` | 1 | Esc → AnsweredNoMsg |
| `TestYesNoModal_ArrowKeysFocus` | 1 | Left → focus Yes; Right → focus No |
| `TestYesNoModal_MouseClickYes` | 1 | MouseLeft on Yes button → AnsweredYesMsg |
| `TestYesNoModal_MouseClickNo` | 1 | MouseLeft on No button → AnsweredNoMsg |
| `TestYesNoModal_MouseMotionHover` | 1 | MouseMotion updates focus button on YesNo modal |
| `TestModalModel_ClosedIgnoresInput` | 1 | Closed modal returns nil cmd for all input |
| `TestModalModel_View_Closed` | 2 | Returns empty string |
| `TestModalModel_View_OKOpen` | 2 | Contains message text, OK button label, border characters |
| `TestModalModel_View_YesNoOpen` | 2 | Contains message text, Yes and No labels, border characters |
| `TestModalModel_OverlayModal` | 2 | Open modal overlays on background; closed returns background unchanged |
| `TestModalModel_Alignment` | 2 | Title/message/button alignment respected in View output |
| `TestModalModel_CustomLabels` | 2 | WithYesLabel/WithNoLabel/WithOKLabel reflected in View output |
| `TestModalModel_CustomStyles` | 1 | With*Style methods update corresponding getters |

#### `progress_modal_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewProgressModal` | 1 | Not open initially; title stored |
| `TestProgressModal_Open` | 1 | IsOpen becomes true |
| `TestProgressModal_Close` | 1 | IsOpen becomes false |
| `TestProgressModal_EscCancels` | 1 | Esc → ProgressCancelledMsg |
| `TestProgressModal_BackgroundKey` | 1 | 'b' key → ProgressBackgroundMsg when enabled |
| `TestProgressModal_BackgroundDisabled` | 1 | 'b' key → no msg when disabled |
| `TestProgressModal_ClosedIgnoresInput` | 1 | Returns nil cmd when closed |
| `TestProgressModal_View_Open` | 2 | Contains title, cancel hint; border characters present |
| `TestProgressModal_View_BackgroundHint` | 2 | Background hint visible when enabled, absent when disabled |
| `TestProgressModal_OverlayModal` | 2 | Overlay compositing correct |

#### `list_model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewListModel` | 1 | Items stored; cursor at 0; not open |
| `TestListModel_Open_FocusesActiveItem` | 1 | Open positions cursor on active item |
| `TestListModel_Open_NoActiveItem` | 1 | Open positions cursor at 0 |
| `TestListModel_KeyUp` | 1 | Cursor decrements; stops at 0 |
| `TestListModel_KeyDown` | 1 | Cursor increments; stops at len-1 |
| `TestListModel_KeySpace` | 1 | Sends ItemSelectedMsg; modal stays open |
| `TestListModel_KeyEnter` | 1 | Sends ListAcceptedMsg; sends ItemSelectedMsg if not already active; closes |
| `TestListModel_KeyNew` | 1 | Sends NewItemRequestedMsg |
| `TestListModel_KeyEdit` | 1 | Enters edit mode |
| `TestListModel_KeyDelete` | 1 | Sends DeleteItemRequestedMsg with correct item |
| `TestListModel_KeyHelp` | 1 | Toggles help visor visibility |
| `TestListModel_HelpVisibleEscClosesHelp` | 1 | Esc closes help visor (not the modal) when help is visible |
| `TestListModel_KeyCancel` | 1 | Sends ListCancelledMsg; closes modal |
| `TestListModel_EditEnterCompletes` | 1 | Enter in edit mode sends EditCompletedMsg with new label |
| `TestListModel_EditEscCancels` | 1 | Esc in edit mode exits edit, restores original |
| `TestListModel_EditFirstKeystrokeOverwrites` | 1 | First typed character replaces entire buffer |
| `TestListModel_EditSubsequentInsertion` | 1 | After first keystroke, typing inserts at cursor |
| `TestListModel_EditCursorMovement` | 1 | Left/Right move edit cursor; Backspace/Delete work |
| `TestListModel_SetItems` | 1 | Items replaced; cursor clamped; label width recalculated |
| `TestListModel_SetCursor` | 1 | Cursor set; clamped to valid range |
| `TestListModel_SetCursorToLast` | 1 | Cursor set to last item |
| `TestListModel_Scrolling` | 1+2 | Offset adjusts when cursor moves beyond visible window |
| `TestListModel_View_Open` | 2 | Contains title, item labels, scrollbar if needed; border geometry correct |
| `TestListModel_View_SelectedItem` | 2 | Selected item visually distinct from others |
| `TestListModel_View_ActiveItem` | 2 | Active item visually distinct |
| `TestListModel_View_EditMode` | 2 | Edit indicator/cursor visible |
| `TestListModel_View_StatusMessage` | 2 | Status message rendered in footer |
| `TestListModel_View_HelpVisor` | 2 | Help visor overlay renders with 3-edge border |
| `TestListModel_OverlayModal` | 2 | Overlay compositing correct |

#### `choice_model_test.go` — Extend Existing

Already covers: navigation, selection, cancel, hotkeys, default index, closed modal, overlay.

| Test | Layer | What It Verifies |
|---|---|---|
| `TestChoiceModel_View_Horizontal` | 2 | Buttons laid out horizontally |
| `TestChoiceModel_View_Vertical` | 2 | Buttons stacked vertically |
| `TestChoiceModel_View_FocusedButton` | 2 | Focused button visually distinct |
| `TestChoiceModel_View_BorderGeometry` | 2 | Border width/height consistent |

#### `overlay_modal_test.go` — Layer 2

| Test | What It Verifies |
|---|---|
| `TestOverlayModal_BasicOverlay` | Foreground composited at correct row/col |
| `TestOverlayModal_ANSIContent` | ANSI-styled content handled correctly |
| `TestOverlayModal_BoundaryConditions` | Row/col at edges; oversized foreground |

#### `teamodal_teatest_test.go` — Layer 3

| Test | What It Verifies |
|---|---|
| `TestOKModal_GoldenRender` | Golden file for OK modal rendering |
| `TestYesNoModal_GoldenRender` | Golden file for YesNo modal rendering |
| `TestChoiceModel_GoldenRender` | Golden file for choice modal rendering |
| `TestListModel_GoldenRender` | Golden file for list modal rendering |

---

### TEASTATUS — `teastatus/*_test.go`

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNew` | 1 | Returns valid model with default styles |
| `TestModel_SetSize` | 1 | Width stored |
| `TestModel_SetMenuItems` | 1 | Items stored; retrievable |
| `TestModel_SetIndicators` | 1 | Indicators stored |
| `TestModel_Update_SetMenuItemsMsg` | 1 | Message updates menu items |
| `TestModel_Update_SetIndicatorsMsg` | 1 | Message updates indicators |
| `TestModel_Update_UnknownMsg` | 1 | Returns model unchanged, nil cmd |
| `TestModel_View_Empty` | 2 | No items/indicators → renders bar background only |
| `TestModel_View_MenuItems` | 2 | Contains key labels in correct format; left-aligned |
| `TestModel_View_Indicators` | 2 | Contains indicator text; right-aligned |
| `TestModel_View_BothZones` | 2 | Menu left + indicators right; total width matches SetSize |
| `TestModel_View_SeparatorKinds` | 2 | Pipe, Space, and Bracket separators render differently |
| `TestModel_View_Truncation` | 2 | Content truncated gracefully when width too narrow |

#### `render_test.go` — Layer 2

| Test | What It Verifies |
|---|---|
| `TestRenderMenuLine` | Renders "[key] label" format; multiple items separated |
| `TestRenderMenuLine_Empty` | Empty input → empty output |

#### `types_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestNewMenuItemFromBinding` | Extracts first key from binding as display text |
| `TestNewStatusIndicator` | Text stored |
| `TestStatusIndicator_WithStyle` | Returns copy with style override |

---

### TEADEP — `teadep/*_test.go`

#### `tree_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestNewTree` | Builds tree from Node; children populated |
| `TestTree_Alternatives` | Returns siblings from parent |
| `TestTree_IsLeaf` | Leaf node returns true; non-leaf returns false |
| `TestTree_HasAlternatives` | Node with siblings returns true |
| `TestNewTree_CircularDeps` | Circular dependencies detected (does not infinite loop) |

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewPathViewer` | 1 | Model created; not yet initialized |
| `TestPathViewer_Initialize` | 1 | Path built using SelectorFunc; SelectedLevel set |
| `TestPathViewer_Initialize_NilSelector` | 1 | Error returned for nil selector |
| `TestPathViewer_KeyUp` | 1 | SelectedLevel decrements; stops at 0; sends FocusNodeMsg |
| `TestPathViewer_KeyDown` | 1 | SelectedLevel increments; stops at path length-1; sends FocusNodeMsg |
| `TestPathViewer_OpenDropdown` | 1 | Dropdown opens at current level when alternatives exist |
| `TestPathViewer_OpenDropdown_NoAlternatives` | 1 | No-op when leaf has no alternatives |
| `TestPathViewer_EnterOnLeaf` | 1 | Sends SelectNodeMsg with leaf tree |
| `TestPathViewer_EnterOnNonLeaf` | 1 | No selection sent (only leaves are selectable) |
| `TestPathViewer_DropdownSelection` | 1 | OptionSelectedMsg rebuilds path from selected alternative |
| `TestPathViewer_DropdownCancellation` | 1 | Dropdown closes; path unchanged |
| `TestPathViewer_WindowSizeMsg` | 1 | Width/Height updated |
| `TestPathViewer_View` | 2 | Path levels rendered in order; selected level visually distinct |
| `TestPathViewer_View_DropdownOpen` | 2 | Dropdown overlay visible over path |
| `TestPathViewer_View_BorderGeometry` | 2 | Border width matches content; no overflow |

---

### TEATREE — `teatree/*_test.go`

#### `node_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestNewNode` | ID, Name, Data stored correctly |
| `TestNode_AddChild` | Child added; parent pointer set |
| `TestNode_SetChildren` | Children replaced; parent pointers correct |
| `TestNode_RemoveChild` | Child removed by ID; returns false for unknown ID |
| `TestNode_InsertChildSorted` | Inserts in correct order per comparator |
| `TestNode_FindByID` | Finds node in subtree; returns nil for unknown ID |
| `TestNode_Depth` | Root=0; child=1; grandchild=2 |
| `TestNode_IsLastChild` | Last child returns true; others false; root false |
| `TestNode_AncestorIsLastChild` | Returns correct boolean slice for tree structure prefixes |
| `TestNode_ExpandCollapse` | Expand/Collapse/Toggle change IsExpanded state |
| `TestNode_HasGrandChildren` | True when children have children; cached |
| `TestNode_Text_FallsBackToName` | Text() returns name when text not set |

#### `tree_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestNewTree` | Nodes stored; first node focused |
| `TestTree_MoveUp` | Focus moves to previous visible node; returns false at top |
| `TestTree_MoveDown` | Focus moves to next visible node; returns false at bottom |
| `TestTree_ExpandFocused` | Focused node expanded; children become visible |
| `TestTree_CollapseFocused` | Focused node collapsed; children hidden |
| `TestTree_ToggleFocused` | Toggles expansion state |
| `TestTree_ExpandAll` | All nodes expanded |
| `TestTree_CollapseAll` | All nodes collapsed; only roots visible |
| `TestTree_VisibleNodes` | Returns correct visible nodes after expand/collapse operations |
| `TestTree_SetFocusedNode` | Focus set by ID; returns false for unknown ID |
| `TestTree_FindByID` | Finds node anywhere in tree |

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewModel` | 1 | Tree stored; height set |
| `TestModel_KeyUp` | 1 | Focus moves up; viewport adjusts |
| `TestModel_KeyDown` | 1 | Focus moves down; viewport adjusts |
| `TestModel_KeyRight_Expand` | 1 | Expands focused node; children become visible |
| `TestModel_KeyRight_EnterChild` | 1 | If already expanded, focus moves to first child |
| `TestModel_KeyLeft_Collapse` | 1 | Collapses focused node |
| `TestModel_KeyLeft_MoveToParent` | 1 | If already collapsed, focus moves to parent |
| `TestModel_KeyToggle` | 1 | Enter/Space toggles expansion |
| `TestModel_WindowSizeMsg` | 1 | Width/Height updated; viewport resized |
| `TestModel_SetSize` | 1 | Dimensions stored |
| `TestModel_SetFocusedNode` | 1 | Focus changed; viewport scrolls to show focused |
| `TestModel_View_BasicTree` | 2 | Root nodes rendered; tree prefixes (├──, └──) correct |
| `TestModel_View_ExpandedTree` | 2 | Children indented; connection lines correct |
| `TestModel_View_FocusedNode` | 2 | Focused node visually distinct |
| `TestModel_View_Scrolling` | 2 | Viewport shows correct slice when tree exceeds height |

#### `teatree_teatest_test.go` — Layer 3

| Test | What It Verifies |
|---|---|
| `TestTreeModel_NavigationGolden` | Golden file: expand tree, navigate, verify rendered output |
| `TestTreeModel_ExpandCollapseGolden` | Golden file: expand/collapse sequence rendering |

---

### TEATEXTSEL — `teatextsel/*_test.go`

#### `selection_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestNewSelection` | Not active; Start/End zero |
| `TestSelection_Begin` | Active; Start set; End equals Start |
| `TestSelection_Extend` | End updated; Active remains true |
| `TestSelection_Clear` | Active becomes false |
| `TestSelection_Normalized` | Start always before End regardless of extend direction |
| `TestSelection_Contains` | Positions inside range return true; outside false; boundary cases |
| `TestSelection_IsEmpty` | Empty when Start==End; not empty otherwise |
| `TestSelectAll` | Spans entire content; Active true |
| `TestSelection_BeforeAfterEqual` | Position comparison functions correct |

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNew` | 1 | Multi-line model created; no selection active |
| `TestNewSingleLine` | 1 | Single-line mode; IsSingleLine() true |
| `TestNewFromTextarea` | 1 | Wraps existing textarea; selection state separate |
| `TestModel_SelectLeft` | 1 | Selection extends left by 1 character |
| `TestModel_SelectRight` | 1 | Selection extends right by 1 character |
| `TestModel_SelectUp` | 1 | Selection extends up one line |
| `TestModel_SelectDown` | 1 | Selection extends down one line |
| `TestModel_SelectWordLeft` | 1 | Selection extends to previous word boundary |
| `TestModel_SelectWordRight` | 1 | Selection extends to next word boundary |
| `TestModel_SelectToLineStart` | 1 | Selection extends to beginning of line |
| `TestModel_SelectToLineEnd` | 1 | Selection extends to end of line |
| `TestModel_SelectToStart` | 1 | Selection extends to document start |
| `TestModel_SelectToEnd` | 1 | Selection extends to document end |
| `TestModel_SelectAll` | 1 | Entire content selected |
| `TestModel_ClearSelection` | 1 | Selection cleared; HasSelection() false |
| `TestModel_CursorMovementClearsSelection` | 1 | Arrow keys without shift clear active selection |
| `TestModel_TypingReplacesSelection` | 1 | Typing with selection active deletes selection then inserts |
| `TestModel_SingleLine_BlocksEnter` | 1 | Enter key blocked in single-line mode |
| `TestModel_SingleLine_BlocksVerticalSelection` | 1 | SelectUp/SelectDown blocked in single-line mode |
| `TestModel_View_NoSelection` | 2 | Renders textarea without selection highlighting |
| `TestModel_View_WithSelection` | 2 | Selected text has highlight styling applied |
| `TestModel_View_MultiLineSelection` | 2 | Multi-line selection highlighting spans lines correctly |

#### `clipboard_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestSelectedText_SingleLine` | Returns correct substring |
| `TestSelectedText_MultiLine` | Returns correct multi-line text |
| `TestSelectedText_NoSelection` | Returns empty string |
| `TestCopy` | SelectedText stored in internal clipboard |
| `TestCut` | Text removed from textarea; stored in clipboard |
| `TestPaste_NoSelection` | Clipboard text inserted at cursor |
| `TestPaste_WithSelection` | Selection replaced by clipboard text |

---

### TEANOTIFY — `teanotify/*_test.go`

Tests written after EXPAND-NOTIFY creates the module.

#### `model_test.go` — Layer 1 + Layer 2

| Test | Layer | What It Verifies |
|---|---|---|
| `TestNewNotifyModel_Defaults` | 1 | Default position, duration, width applied |
| `TestNewNotifyModel_CustomOpts` | 1 | All NotifyOpts fields respected |
| `TestNotifyModel_WithPosition` | 1 | Position updated via wither |
| `TestNotifyModel_WithMinWidth` | 1 | MinWidth updated via wither |
| `TestNotifyModel_RegisterNoticeType` | 1 | Custom notice type registered; retrievable |
| `TestNotifyModel_DefaultNoticeTypes` | 1 | Info, Warn, Error, Debug registered by default |
| `TestNotifyModel_NewNotifyCmd` | 1 | Returns non-nil cmd; cmd produces notifyMsg |
| `TestNotifyModel_Update_NotifyMsg` | 1 | Notice becomes active; HasActiveNotice() true |
| `TestNotifyModel_Update_TickExpiration` | 1 | After sufficient ticks, notice expires; HasActiveNotice() false |
| `TestNotifyModel_Update_EscClose` | 1 | Esc closes notice when AllowEscToClose enabled |
| `TestNotifyModel_Update_EscNoClose` | 1 | Esc ignored when AllowEscToClose disabled |
| `TestNotifyModel_HasActiveNotice` | 1 | False before activation; true during; false after expiry |
| `TestNotifyModel_Render_TopLeft` | 2 | Notice overlaid at top-left of content |
| `TestNotifyModel_Render_TopCenter` | 2 | Notice overlaid at top-center |
| `TestNotifyModel_Render_TopRight` | 2 | Notice overlaid at top-right |
| `TestNotifyModel_Render_BottomLeft` | 2 | Notice overlaid at bottom-left |
| `TestNotifyModel_Render_BottomCenter` | 2 | Notice overlaid at bottom-center |
| `TestNotifyModel_Render_BottomRight` | 2 | Notice overlaid at bottom-right |
| `TestNotifyModel_Render_NoActiveNotice` | 2 | Returns content unchanged |
| `TestNotifyModel_Render_DynamicWidth` | 2 | Width clamps between min and max |
| `TestNotifyModel_Render_BorderGeometry` | 2 | Border width/height consistent; no overflow |

#### `util_test.go` — Layer 1 + Layer 2

| Test | What It Verifies |
|---|---|
| `TestGetLines` | Splits on newlines; returns correct widest width |
| `TestGetLines_ANSI` | ANSI escapes don't inflate width |
| `TestCutLeft` | Removes correct printable chars; preserves ANSI |
| `TestCutLeft_WideChars` | Wide characters (CJK) handled correctly |
| `TestCutRight` | Keeps correct printable chars; adds reset if needed |
| `TestCutRight_WideChars` | Wide characters handled correctly |
| `TestHangingWrap` | Wraps at textWidth; continuation lines indented |
| `TestHangingWrap_ShortMessage` | No wrapping when message fits |

#### `position_test.go` — Layer 1

| Test | What It Verifies |
|---|---|
| `TestPosition_IsValid` | Valid positions return true; invalid return false |
| `TestPosition_String` | Returns readable string ("top-left", etc.) |
| `TestPosition_Label` | Returns display label ("Top Left", etc.) |

---

## Test Infrastructure

### Shared Test Helpers

Each module needing Layer 3 tests adds `github.com/charmbracelet/x/exp/teatest` to its `go.mod` (test dependency only).

No shared test helper module needed initially — each module's tests are self-contained. If patterns emerge during implementation, extract to a shared internal testhelper later.

### Golden File Management

Convention:
- Files in `<module>/testdata/<TestName>.golden`
- Generated by teatest's `RequireEqualOutput` with `-update` flag
- Terminal size fixed at 80x24 for reproducibility
- Golden files committed to git

### Test Data

Complex test fixtures (trees, node hierarchies, option lists) defined as helper functions within `_test.go` files — same pattern as existing `threeOptions()` in `choice_model_test.go`.

---

## Execution Order

Tests can be written in any order, but recommended priority:

1. **teautils** — Pure functions, no model dependencies, quick wins
2. **teadd** — Simple model, few interactions, establishes testing patterns
3. **teastatus** — Simple passive model
4. **teamodal** — Extend existing tests; high rendering regression risk
5. **teatree** — Medium complexity; node data structures + viewport
6. **teadep** — Depends on teadd; integration behavior
7. **teatextsel** — Highest complexity; textarea wrapping + selection
8. **teanotify** — Written after EXPAND-NOTIFY creates the module

---

## V2 Migration Impact on Tests

Tests that will need mechanical updates during MIGRATE_V2:

| Change | Affected Tests |
|---|---|
| `tea.KeyMsg` → `tea.KeyPressMsg` | All Layer 1 tests constructing key messages |
| `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}}` → v2 construction | Hotkey tests, typing tests |
| `tea.KeyMsg{Type: tea.KeyTab}` → v2 construction | Navigation tests |
| `View() string` → `View() tea.View` | All Layer 2 tests calling View() |
| `tea.MouseMsg` → type-switch variants | teamodal mouse tests |
| `teatest` → `teatest/v2` | All Layer 3 tests |
| Golden file regeneration | All Layer 3 golden file tests |

Tests that should survive v2 unchanged:
- teautils pure function tests
- teatree node/tree data structure tests
- Selection data structure tests
- Clipboard logic tests
