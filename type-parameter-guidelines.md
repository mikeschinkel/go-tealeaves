# Type Parameter Guidelines for go-tealeaves

## Decision Rule

### Use `Model[T]` when ALL of these are true:
1. The model holds **user-provided data types** (not fixed internal structs)
2. **Messages need to carry that data back type-safely** (e.g., `SelectionMsg[T]{Item: T}`)
3. **Renderers or providers need type access** (e.g., `NodeProvider[T].Style(node *Node[T])`)

### Use plain `Model` when ANY of these are true:
1. Selection returns an **ID/index**, not the object itself
2. Items are **fixed internal structs** (like `ChoiceOption`)
3. The model is **display-only** (no data export to parent)

## Current Inventory

| Model | Generic? | Rationale |
|-------|----------|-----------|
| `TreeModel[T]` | Yes | Nodes carry user payload T; NodeProvider needs T for rendering |
| `ListModel[T]` | Yes | Messages return selected item as T; parent needs type-safe access |
| `DrillDownModel[T]` | Yes | Same Node[T] as TreeModel; drill-down selections carry T |
| `ConfirmModel` | No | Returns Yes/No boolean answer, no data payload |
| `ChoiceModel` | No | Returns OptionID string; options are fixed ChoiceOption structs |
| `ProgressModal` | No | Display only — spinner with text |
| `GridModel` | No | Display/sort/filter; data is [][]string (fixed) |
| `StatusBarModel` | No | Display only — status text |
| `TextSnipModel` | No | Display only — text selection overlay |

## Quick Test

Before adding a type parameter, ask:

> "When the user selects/interacts, does the parent need the actual typed object back?"

- **Yes** → `Model[T]` with `SomeMsg[T]{Item: T}`
- **No** → Plain `Model` with `SomeMsg{ID: string}` or `SomeMsg{}`

## Anti-patterns

- Don't use `[T any]` when `T` is only used internally and never surfaces in messages
- Don't use generics just because the model accepts a slice — if items are a fixed struct, no generics needed
- Don't add a type parameter "for future flexibility" — add it when there's a concrete need
