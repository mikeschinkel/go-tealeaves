# teadep - Bubble Tea Dependency Path Viewer

A standalone Bubble Tea component for visualizing dependency paths with interactive branch selection.

## Overview

teadep displays **one dependency path at a time** (from project to leaf dependency) with dropdown navigation for alternate dependencies at each level. It's designed to help users navigate complex dependency trees without overwhelming them with the full graph.

## Features

- **Linear path display**: Shows one path from root → leaf (vertical list, leaf at bottom)
- **Interactive navigation**: Navigate between levels with arrow keys
- **Alternative selection**: Choose different dependencies at any level via dropdown
- **Leaf selection**: Start with leaf node selected (common use case)
- **Customizable strategy**: Caller defines what makes a "best" path via ChildSelector function
- **Minimal dependencies**: Only requires bubbletea and lipgloss

## Installation

```bash
go get github.com/mikeschinkel/go-tealeaves/teadep
```

## Quick Start

```go
package main

import (
	"github.com/mikeschinkel/go-tealeaves/teadep"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Build your dependency tree
	root := teadep.DependencyNode{
		DisplayName: "[R] ~/Projects/myproject",
		Children: []teadep.DependencyNode{
			{
				DisplayName: "[RM] ~/Projects/go-pkgs/mylib",
				Alternatives: []teadep.DependencyNode{
					{DisplayName: "[RM] ~/Projects/go-pkgs/altlib"},
				},
			},
		},
	}

	// Define selection strategy
	selectChild := func(parent teadep.DependencyNode, children []teadep.DependencyNode) (teadep.DependencyNode, error) {
		// Return first child (simplest)
		// Or implement custom logic: longest path, most in-flux, etc.
		return children[0], nil
	}

	// Create viewer
	viewer, _ := teadep.NewPathViewer(root, selectChild)

	// Run with Bubble Tea
	tea.NewProgram(viewer).Run()
}
```

## API Reference

### Types

#### DependencyNode

```go
type DependencyNode struct {
    DisplayName  string            // Pre-formatted: "[RM] ~/Projects/go-pkgs/go-dt"
    Alternatives []DependencyNode  // Siblings - other options at this level
    Children     []DependencyNode  // Next level down in path
}
```

**Methods**:
- `IsLeaf() bool` - Returns true if node has no children
- `HasAlternatives() bool` - Returns true if node has siblings

#### ChildSelector

```go
type ChildSelector func(parent DependencyNode, children []DependencyNode) (best DependencyNode, err error)
```

Strategy function for choosing which child to follow at each level. Allows caller to define what "best" means (longest path, deepest, most in-flux, etc.).

#### PathViewerModel

```go
type PathViewerModel struct {
    Root          DependencyNode
    Path          []DependencyNode
    SelectedLevel int
    SelectChild   ChildSelector
    // ... (styling and display fields)
}
```

**Constructor**:
```go
func NewPathViewer(root DependencyNode, selector ChildSelector, opts ...Option) (PathViewerModel, error)
```

**Methods**:
- `SetSize(width, height int)` - Update display dimensions

### Messages

#### SelectDependencyMsg

Sent when user confirms selection with Enter on a leaf node.

```go
type SelectDependencyMsg struct {
    Node DependencyNode
}
```

#### ChangeDependencyMsg

Sent when user picks a different dependency from dropdown.

```go
type ChangeDependencyMsg struct {
    Level int
    Node  DependencyNode
}
```

### Styling

```go
// Functional options for customization
teadep.WithPathStyle(style lipgloss.Style)
teadep.WithSelectedStyle(style lipgloss.Style)
```

## Path Selection Algorithm

The component uses a caller-provided `ChildSelector` function to determine the "best" child at each level:

1. **Initial path**: Built by walking from root to leaf, using selector at each level
2. **Alternative selection**: When user selects an alternative, path is rebuilt from that point using the same selector
3. **Leaf selection**: Initial view has leaf node selected (bottom of path)

## Navigation

- `↑↓` / `k``j`: Move selection between levels
- `Space` / `→`: Open dropdown for alternatives (if available)
- `Enter`: Open dropdown on non-leaf, select on leaf
- `Esc`: Close dropdown or return to previous view

## Visual Indicators

- **▶ indicator**: Shown when node has alternatives (siblings)
  - Means: "There are other options besides THIS node at this level"
- **Highlighting**: Selected level is highlighted
- **DisplayName**: Rendered as-is (pre-formatted by caller)

## Example

See `example/main.go` for a complete working demo with hardcoded xmlui/cli dependency tree.

```bash
cd demo
go run main.go
```

## License

MIT License - See LICENSE.txt
