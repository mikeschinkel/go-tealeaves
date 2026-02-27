# ADR: teadep Component Architecture

**Date**: 2025-01-24
**Status**: Accepted

## Context

This component provides an interactive dependency path viewer for Bubble Tea applications. It displays one path at a time (root → leaf) with dropdown navigation for selecting alternative dependencies at each level, rather than overwhelming users with a full dependency graph.

## Decision

### 1. Module Structure: Flat Package With Separate Example

**Decision**: Component code lives in a single package, with the example app having its own go.mod.

**Rationale**:
- Keeps component simple - all code in one package
- Example app demonstrates standalone usage with its own go.mod
- Allows for future extraction to separate repository without restructuring

**Structure**:
```
teadep/
├── *.go              # Component code
└── example/
    ├── go.mod        # Separate module for demo
    └── main.go
```

### 2. SelectorFunc Abstraction Pattern

**Decision**: Use a caller-provided function to determine "best" child at each level, rather than hardcoding selection logic.

**Rationale**:
- Keeps component generic and reusable
- Different applications can define what "best" means (longest path, highest priority, most frequently used, etc.)
- Business logic stays in calling application, not in presentation component
- Allows different use cases without modifying the component

**Interface**:
```go
type SelectorFunc func(parent *Tree, children []*Tree) (best *Tree, err error)
```

The selector can capture any context (metadata, configuration) in a closure:
```go
selector := func(parent *Tree, children []*Tree) (*Tree, error) {
    return selectBestChild(parent, children, metadata)
}
```

### 3. Node Interface Design

**Decision**: Use a simple Node interface with BaseNode embedding pattern.

**Rationale**:
- Allows different node types (repos, modules, packages, etc.) without type assertions everywhere
- BaseNode provides common functionality for typical use cases
- Applications can embed BaseNode and add domain-specific fields
- Clean separation: component works with interface, applications use concrete types

**Pattern**:
```go
// Component defines minimal interface
type Node interface {
    DisplayName() string
    Dependencies() []Node
    SetDisplayName(string)
    SetDependencies([]Node)
}

// Component provides base implementation
type BaseNode struct {
    displayName  string
    dependencies []Node
}

// Applications embed and extend with domain-specific fields
type CustomNode struct {
    *BaseNode
    CustomField string  // Application-specific data
}
```

### 4. Display Name Pre-Formatting

**Decision**: DisplayName is pre-formatted by the calling application, not by the component.

**Rationale**:
- Component doesn't know about application-specific formatting conventions
- Different applications can format differently (type indicators, paths, colors, etc.)
- Keeps component generic and presentation-agnostic
- Calling application controls all presentation details

**Example**:
```go
// Application formats as needed:
node.SetDisplayName("[PKG] github.com/user/package")
// or
node.SetDisplayName("~/Projects/myproject")
// or any other format

// Component just renders:
line := prefix + tree.Node.DisplayName()
```

### 5. Tree Structure Over Graph Structure

**Decision**: Use a Tree wrapper around Node with separate tracking of alternatives (siblings).

**Rationale**:
- Path viewer shows one path at a time, not a full graph
- Alternatives are siblings at the same level, not children
- Tree structure makes path building/rebuilding straightforward
- Separation: Node has dependencies (children), Tree tracks alternatives (siblings)

**Structure**:
```go
type Tree struct {
    Node         Node     // Actual dependency node
    Children     []*Tree  // Children in tree (from Node.Dependencies)
    Parent       *Tree    // Parent in tree (for traversal)
    alternatives []*Tree  // Siblings at same level (set by BuildTree)
}
```

### 6. Integration with teadd

**Decision**: Use teadd as a dependency for dropdown behavior.

**Rationale**:
- Don't reinvent dropdown behavior
- teadd handles positioning, modal behavior, keyboard navigation
- Dropdown positioned dynamically based on selected path level
- Clean separation: PathViewerModel owns dropdown instance and lifecycle

**Integration Pattern**:
```go
// Position dropdown at current path level
row := model.SelectedLevel + 3  // Account for title/padding
col := 2                         // Indent

model.Dropdown = teadd.NewModel(items, row, col, args)
model.DropdownOpen = true
```

## Consequences

### Positive
- Clean, reusable component architecture
- Calling application controls business logic (what makes a "best" path)
- Component stays generic and publishable
- Clear separation of concerns
- Can be used for any tree-like dependency visualization

### Negative
- Caller must build full tree structure upfront (not lazy/on-demand)
- Caller must provide selector function (not built-in strategies)
- More setup code in calling application

### Trade-offs
- **Upfront tree building** vs lazy loading: Acceptable for most dependency visualization use cases (typically showing a filtered subset)
- **Caller-provided selector** vs built-in strategies: Flexibility wins over convenience for a generic component

## Alternatives Considered

### Alternative 1: Separate Go Module From Start
**Rejected**: Too complex to manage during active development. Can extract later.

### Alternative 2: Built-in Selection Strategies
**Rejected**: Would couple component to specific use cases. Selector pattern is more flexible.

### Alternative 3: Graph Structure Instead of Tree
**Rejected**: Over-engineered for showing one path at a time. Tree is simpler.

### Alternative 4: Component Formats Display Names
**Rejected**: Would require component to know about application-specific formatting conventions. Pre-formatting keeps component generic.

## References

- Original implementation plan: `PLAN.md` (deleted after ADR creation)
- Example usage: `example/main.go`
- API documentation: `README.md`
