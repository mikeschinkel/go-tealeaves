# Planning Brief: GridModel[T] Generic Conversion

## Objective
Convert `teagrid.GridModel` from `[][]string`-based data to a generic `GridModel[T]` that operates on typed row data.

## Why
Currently the grid works with `[][]string` — every cell is a string. This means:
- When a user selects a row, the parent gets a string slice, not the business object
- Sorting requires string comparison; no type-safe field access
- Filtering operates on display strings, not structured data
- Column extraction is implicit (positional index), not explicit

With `GridModel[T]`, the parent would:
- Receive `RowSelectedMsg[T]{Item: T}` with the actual typed row
- Define column extractors as `func(T) string` for each column
- Sort with `func(a, b T) bool` comparators
- Filter with `func(T) bool` predicates

## Current API (simplified)
```go
type GridModel struct { /* uses [][]string internally */ }
func NewGridModel(columns []Column, rows [][]string) GridModel
```

## Proposed API (sketch)
```go
type GridModel[T any] struct { /* uses []T internally */ }

type ColumnDef[T any] struct {
    Title     string
    Width     int
    Extract   func(T) string        // How to get display text from T
    Compare   func(a, b T) int      // Optional: for sorting
}

func NewGridModel[T any](columns []ColumnDef[T], rows []T) GridModel[T]
```

## Scope
- Significant refactor of teagrid internals
- All 4 example apps need rewriting
- Sort, filter, pagination logic all needs generic adaptation
- ~145 references across 19 files

## Prerequisite
- Complete the Model → GridModel rename first (simple rename, no API change)
- Then plan the generic conversion separately

## Type Parameter Guideline Check
1. Model holds user-provided data types? YES — rows are user structs
2. Messages carry data back type-safely? YES — row selection returns T
3. Renderers need type access? YES — column extractors need T

All three criteria met → `GridModel[T]` is correct.

## Downstream: Gomion Refactor
Gomion uses teagrid and would need refactoring to use the generic `GridModel[T]` API.
This should be planned as part of the same effort — extract business object types,
define column extractors, and switch from `[][]string` row construction to typed rows.
