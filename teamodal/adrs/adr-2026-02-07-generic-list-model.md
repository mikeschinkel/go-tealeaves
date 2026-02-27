# ADR: Generic ListModel[T] for Modal List Selection

**Date**: 2026-02-07
**Status**: Accepted

## Context

The teamodal package provides modal dialogs for Bubble Tea applications. We needed a reusable modal component for managing lists of items with CRUD operations, specifically to support "Commit Alternatives" in the Gomion TUI - allowing users to select between different AI-generated commit message alternatives.

The component needed to be generic enough for reuse across different list types while providing a consistent modal experience.

## Decision

### 1. Naming Convention: "Commit Alternatives"

**Decision**: Use "Commit Alternatives" as the terminology for commit groupings/options.

**Rationale**:
- Clear and descriptive - alternatives to choose from
- Neutral term that works for AI-generated and user-created options
- Consistent with common UX patterns (alternative views, alternative selections)

### 2. Type Parameter Approach: Interface Constraint

**Decision**: Use a Go generic type parameter with an interface constraint rather than reflection or `any`.

**Interface**:
```go
type ListItem interface {
    ID() string      // Unique identifier
    Label() string   // Display text
    IsActive() bool  // Currently active/in-use item
}
```

**Rationale**:
- Type safety at compile time
- No runtime reflection overhead
- Clear contract for what list items must provide
- IDE autocompletion works correctly
- Familiar pattern for Go developers (similar to constraints in sort, slices packages)

**Trade-off**: Callers must implement the interface, but this is minimal (3 methods) and ensures consistent behavior.

### 3. Sub-operations: Callbacks to Parent via Messages

**Decision**: Emit Bubble Tea messages for CRUD operations rather than handling them internally.

**Messages**:
```go
type ItemSelectedMsg[T ListItem] struct { Item T }  // Space: preview-select
type ListAcceptedMsg[T ListItem] struct { Item T }  // Enter: accept and close
type NewItemRequestedMsg struct{}
type EditCompletedMsg[T ListItem] struct { Item T; NewLabel string }  // Inline edit completed
type DeleteItemRequestedMsg[T ListItem] struct { Item T }
type ListCancelledMsg struct{}  // Esc: cancel and close
```

**Rationale**:
- Follows Bubble Tea's message-passing architecture
- Parent model controls all business logic (API calls, validation, confirmation dialogs)
- Component stays presentation-focused
- Easy to add confirmation modals before destructive actions
- Parent can ignore messages it doesn't want to handle
- Inline editing (press 'e') is handled internally by ListModel, emitting EditCompletedMsg when done

**Pattern**:
```go
// Parent handles messages
case teamodal.EditCompletedMsg[CommitAlternative]:
    // Update the item with the new label
    updateItem(msg.Item.ID(), msg.NewLabel)

case teamodal.DeleteItemRequestedMsg[CommitAlternative]:
    // Show confirmation modal, then delete if confirmed
    m.confirmModal = teamodal.NewYesNoModal("Delete this alternative?", nil)
```

### 4. Selection Mode: Single-Select Only

**Decision**: Support only single-item selection, not multi-select.

**Rationale**:
- Simpler implementation and UX
- Primary use case (commit alternatives) only needs single selection
- Multi-select adds complexity (shift-click, ctrl-click, select all, etc.)
- Can be added later if needed without breaking changes

### 5. Active vs Selected Distinction

**Decision**: Distinguish between "cursor position" (navigation) and "active item" (currently in use).

**Implementation**:
- `cursor` field tracks which item is highlighted (keyboard navigation)
- `IsActive()` method on items marks the item currently in use
- Visual distinction: `>` prefix for cursor, `[ACTIVE]` badge for active item

**Rationale**:
- Users need to see both where they are navigating AND which item is currently applied
- Common pattern in settings UIs (current selection vs browsing options)
- Prevents confusion when navigating away from active item

### 6. Scrollbar for Long Lists

**Decision**: Show a scrollbar when items exceed `maxVisible` (default 8).

**Implementation**:
- Vertical scrollbar on right edge using Unicode block characters
- `offset` tracks scroll position
- Cursor movement automatically adjusts offset to keep cursor visible

**Rationale**:
- Clear indication that more items exist above/below
- Familiar scrolling metaphor
- Works in terminal environment without mouse scroll support

### 7. Reuse Existing Overlay Infrastructure

**Decision**: Reuse `OverlayModal()` from overlay_modal.go for positioning.

**Rationale**:
- Consistent positioning with existing modals
- ANSI-aware string composition already solved
- No code duplication
- Same visual appearance across modal types

### 8. Help Visor Integration

**Decision**: Use teautils KeyRegistry pattern for help visor display.

**Implementation**:
- Press `?` to toggle inline help visor above footer
- Help visor shows all keys organized by category (Navigation, Selection, Actions, System)
- Uses `teautils.RenderHelpVisor()` for consistent styling
- Shortened footer shows only `[?] Help [a] Add [e] Edit [d] Delete [esc] Cancel`
- Full key reference available in help visor

**Rationale**:
- Keeps footer compact while providing discoverability
- Reuses teautils infrastructure
- Consistent with other Gomion components
- User-configurable keys supported via `ListKeyMap`

### 9. Fixed Label Width for Consistent Layout

**Decision**: Use configurable `LabelWidth` for consistent dialog sizing.

**Implementation**:
- `LabelWidth` field specifies fixed width for labels and edit field
- Default: calculated from longest item label
- All active items padded to width so `[ACTIVE]` badge aligns
- Edit field maintains width during editing (no dialog shift)

**Rationale**:
- Prevents dialog width changes when entering edit mode
- Ensures `[ACTIVE]` badge always aligns consistently
- Predictable layout regardless of content changes

## Consequences

### Positive
- Type-safe generic component usable for any list type
- Clean separation between presentation (ListModel) and business logic (parent)
- Consistent with existing teamodal patterns
- Minimal interface requirement for list items
- Future-proof: can add multi-select without breaking API

### Negative
- Callers must implement ListItem interface
- No built-in persistence or validation
- Single-select only (may need extension later)

### Trade-offs
- **Interface constraint** vs `any`: Type safety wins over flexibility
- **Message callbacks** vs internal handling: Separation of concerns wins over convenience
- **Single-select** vs multi-select: Simplicity wins over completeness

## Files Created

| File | Purpose |
|------|---------|
| `list_model.go` | Main ListModel[T] implementation |
| `list_keymap.go` | ListKeyMap with default key bindings |
| `list_messages.go` | Message types for list events |
| `list_styles.go` | Default styles for list items |

## Usage Example

```go
// Implement ListItem interface
type CommitAlternative struct {
    id     string
    label  string
    active bool
}

func (c CommitAlternative) ID() string     { return c.id }
func (c CommitAlternative) Label() string  { return c.label }
func (c CommitAlternative) IsActive() bool { return c.active }

// Create and use ListModel
alternatives := []CommitAlternative{...}
listModal := teamodal.NewListModel(alternatives, &teamodal.ListModelArgs{
    Title:        "Commit Alternatives",
    ScreenWidth:  80,
    ScreenHeight: 24,
})

// Open modal
listModal = listModal.Open()

// Handle in parent Update()
case teamodal.ItemSelectedMsg[CommitAlternative]:
    // Apply selected alternative
```

## Alternatives Considered

### Alternative 1: Reflection-Based Generic List
**Rejected**: Runtime type checking is error-prone and loses IDE support.

### Alternative 2: Internal CRUD Handling
**Rejected**: Would couple presentation to business logic. Different use cases need different handling (API calls, local state, confirmation dialogs).

### Alternative 3: Multi-Select Support
**Deferred**: Not needed for initial use case. Can be added later with `multiSelect` flag.

### Alternative 4: Lazy Loading Items
**Rejected**: Over-engineering for typical list sizes. Can pre-load all alternatives.

## References

- Existing modal implementation: `model.go`
- Overlay positioning: `overlay_modal.go`
- Key binding pattern: `keymap.go`
