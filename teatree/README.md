# teatree - A Generic Tree Widget for Bubble Tea

An interactive, keyboard-navigable tree component for [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI applications. Fully generic — works with any data type.

## Features

- **Generic nodes**: `Node[T]`, `Tree[T]`, `Model[T]` — parameterized by your data type
- **Keyboard navigation**: Vim-style (`hjkl`) and arrow keys, with expand/collapse/toggle
- **Customizable rendering**: `NodeProvider` interface for complete control over text, icons, styles, and tree connectors
- **Predefined branch styles**: Compact, wide, ASCII, and minimal — or define your own
- **Viewport scrolling**: Automatic vertical scrolling with focus tracking
- **File tree helpers**: `BuildFileTree()` builds hierarchical nodes from flat file lists
- **Focus management**: Single-focus model with visible-node-aware up/down movement

## Installation

```bash
go get github.com/mikeschinkel/go-tealeaves/teatree
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/mikeschinkel/go-tealeaves/teatree"
)

type ItemData struct {
    Description string
}

func main() {
    // Build tree
    root := teatree.NewNode("root", "Project", ItemData{Description: "Root"})
    root.SetExpanded(true)

    src := teatree.NewNode("src", "src/", ItemData{Description: "Source"})
    src.AddChild(teatree.NewNode("main", "main.go", ItemData{Description: "Entry point"}))
    src.AddChild(teatree.NewNode("util", "util.go", ItemData{Description: "Utilities"}))
    root.AddChild(src)

    root.AddChild(teatree.NewNode("readme", "README.md", ItemData{Description: "Docs"}))

    // Create tree and model
    tree := teatree.NewTree([]*teatree.Node[ItemData]{root}, nil)
    model := teatree.NewModel(tree, 20)

    // Run
    p := tea.NewProgram(model)
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}
```

Output:
```
└─▼ Project
   ├─▶ src/
   └─── README.md
```

## API Reference

### Node[T]

A single tree node with generic payload data.

```go
// Create a node
node := teatree.NewNode[MyData](id, displayName, data)
```

| Method | Description |
|--------|-------------|
| `ID() string` | Unique identifier |
| `Name() string` | Display name |
| `SetName(string)` | Change display name |
| `Data() *T` | Payload data (pointer) |
| `Parent() *Node[T]` | Parent node (nil for roots) |
| `Children() []*Node[T]` | Child nodes |
| `HasChildren() bool` | Whether node has children |
| `HasGrandChildren() bool` | Whether any child has children (cached) |
| `IsRoot() bool` | Whether parent is nil |
| `IsExpanded() bool` | Expansion state |
| `SetExpanded(bool)` | Set expansion state |
| `Expand()` / `Collapse()` / `Toggle()` | Expansion shortcuts |
| `IsVisible() bool` | Visibility state |
| `SetVisible(bool)` | Set visibility |
| `AddChild(*Node[T])` | Add child (sets parent pointer) |
| `SetChildren([]*Node[T])` | Replace all children |
| `RemoveChild(id string) bool` | Remove child by ID |
| `InsertChildSorted(*Node[T], less)` | Insert in sorted position |
| `FindByID(string) *Node[T]` | Recursive search by ID |
| `Depth() int` | Depth from root (0-based) |
| `IsLastChild() bool` | Whether last among siblings |
| `AncestorIsLastChild() []bool` | Ancestor last-child flags (for rendering) |

### Tree[T]

Collection of root nodes with focus management.

```go
tree := teatree.NewTree(roots, &teatree.TreeArgs[T]{
    NodeProvider: myProvider,   // Custom rendering (optional)
    FocusedNode:  initialNode,  // Initial focus (optional)
})
```

| Method | Description |
|--------|-------------|
| `Nodes() []*Node[T]` | Root nodes |
| `SetNodes([]*Node[T])` | Replace roots (revalidates focus) |
| `Provider() NodeProvider[T]` | Current provider |
| `FocusedNode() *Node[T]` | Currently focused node |
| `SetFocusedNode(id string) bool` | Focus by ID |
| `FindByID(string) *Node[T]` | Find node anywhere in tree |
| `VisibleNodes() []*Node[T]` | All visible nodes in order |
| `MoveUp() bool` / `MoveDown() bool` | Move focus |
| `ExpandFocused()` / `CollapseFocused()` / `ToggleFocused()` | Modify focused node |
| `ExpandAll()` / `CollapseAll()` | Bulk expand/collapse |

### Model[T]

Bubble Tea model wrapping a Tree with keyboard input and viewport scrolling.

```go
model := teatree.NewModel(tree, height)
```

| Method | Description |
|--------|-------------|
| `Init() tea.Cmd` | Bubble Tea init |
| `Update(tea.Msg) (Model[T], tea.Cmd)` | Bubble Tea update (handles keys + resize) |
| `View() string` | Renders visible portion of tree |
| `Tree() *Tree[T]` | Underlying tree |
| `SetSize(w, h int) Model[T]` | Update dimensions |
| `MaxLineWidth() int` | Width needed for all content |
| `FocusedNode() *Node[T]` | Currently focused node |
| `SetFocusedNode(id string) Model[T]` | Focus by ID and scroll into view |

**Default key bindings** (`Model.Keys`):

| Key | Action |
|-----|--------|
| `Up` / `k` | Move focus up |
| `Down` / `j` | Move focus down |
| `Right` / `l` | Expand or enter first child |
| `Left` / `h` | Collapse or move to parent |
| `Enter` / `Space` | Toggle expand/collapse |

### NodeProvider[T] (Interface)

Customizes how nodes are rendered.

```go
type NodeProvider[T any] interface {
    Icon(node *Node[T]) string
    Text(node *Node[T]) string
    Suffix(node *Node[T]) string
    Style(node *Node[T], tree *Tree[T]) lipgloss.Style
    ExpanderControl(node *Node[T]) string
    BranchStyle() BranchStyle
}
```

**Built-in providers:**
- `CompactNodeProvider[T]` — Default provider using `CompactBranchStyle`. Embed this and override individual methods.
- `SimpleNodeProvider[T]` — Minimal provider for flat lists.

```go
// Custom provider — embed CompactNodeProvider, override what you need
type MyProvider struct {
    teatree.NodeProvider[MyData]
}

func NewMyProvider() *MyProvider {
    return &MyProvider{
        NodeProvider: teatree.NewCompactNodeProvider[MyData](teatree.TriangleExpanderControls),
    }
}

func (p *MyProvider) Style(node *teatree.Node[MyData], tree *teatree.Tree[MyData]) lipgloss.Style {
    style := lipgloss.NewStyle().Foreground(lipgloss.Color("76"))
    if tree.IsFocusedNode(node) {
        return style.Reverse(true)
    }
    return style
}
```

### BranchStyle

Controls tree connector characters and spacing.

```go
type BranchStyle struct {
    Vertical          string // Continuation line (e.g. "│")
    Horizontal        string // Horizontal line (e.g. "─")
    MiddleChild       string // Branch connector (e.g. "├─")
    LastChild         string // Last child connector (e.g. "└─")
    EmptySpace        string // Blank indent (e.g. "  ")
    PreExpanderIndent string // Space before expander symbol
    PreIconIndent     string // Space before icon
    PreTextIndent     string // Space before text
    PreSuffixIndent   string // Space before suffix
    ExpanderControls  ExpanderControls
}
```

**Predefined styles:**

| Style | Connectors | Indent per level | Example |
|-------|-----------|-----------------|---------|
| `CompactBranchStyle` | `├─` `└─` `│` | 2 chars | `├─▶ item` |
| `WideBranchStyle` | `├── ` `└── ` `│   ` | 4 chars | `├── ▶ item` |
| `ASCIIBranchStyle` | `+- ` `` `- `` `\| ` | 2 chars | `+- item` |
| `MinimalBranchStyle` | `├` `└` `│` | 1 char | `├item` |
| `DefaultBranchStyle` | (none) | spacing only | ` ▶ item` |

### ExpanderControls

Symbols for expand/collapse indicators.

| Preset | Expand | Collapse | Leaf |
|--------|--------|----------|------|
| `TriangleExpanderControls` | `▶` | `▼` | (empty) |
| `PlusExpanderControls` | `+` | `─` | (empty) |
| `NoExpanderControls` | (empty) | (empty) | (empty) |

### File Tree Helpers

Build a tree from a flat list of files:

```go
files := []*teatree.File{
    teatree.NewFile("src/main.go", nil),
    teatree.NewFile("src/util.go", nil),
    teatree.NewFile("README.md", nil),
}

nodes := teatree.BuildFileTree(files, teatree.BuildFileTreeArgs{
    RootPath: "myproject",
})
```

This creates synthetic folder nodes and wires up the parent/child hierarchy. Files are sorted by path.

## Patterns

### Phantom Parent for Root-Child Indentation

By default, children of root nodes (nodes passed to `NewTree`) render at column 0, without indentation relative to their parent. This is because `AncestorIsLastChild()` stops walking when it reaches a node whose parent is `nil` — which is the case for tree roots.

If you need children of a visible root node to be indented (e.g., a tree with a single root that acts as a header), use the **phantom parent** technique:

```go
// Phantom node — never passed to the tree, never rendered.
// Exists only to give root a non-nil parent pointer.
phantom := teatree.NewNode("phantom", "", MyData{})

// Visible root — this IS passed to the tree
root := teatree.NewNode("root", "My Project", MyData{})
root.SetExpanded(true)
phantom.AddChild(root) // sets root.parent = phantom

// Add children
root.AddChild(teatree.NewNode("child-1", "First", MyData{}))
root.AddChild(teatree.NewNode("child-2", "Second", MyData{}))

// Pass root (NOT phantom) as the tree root
tree := teatree.NewTree([]*teatree.Node[MyData]{root}, nil)
```

**Without phantom parent:**
```
└─▼ My Project
├─── First
└─── Second
```
Children render at the same column as the root — no indentation.

**With phantom parent:**
```
└─▼ My Project
   ├─── First
   └─── Second
```
Children are properly indented under the root.

**How it works:**

The renderer calls `AncestorIsLastChild()` to determine indentation depth. That method walks up the parent chain:

```go
// node.go — simplified
func (n *Node[T]) AncestorIsLastChild() []bool {
    var result []bool
    current := n.parent
    for current != nil && current.parent != nil {
        result = append([]bool{current.IsLastChild()}, result...)
        current = current.parent
    }
    return result
}
```

The loop condition `current.parent != nil` stops at tree roots (parent is nil). Without the phantom, `root.parent == nil`, so the loop immediately stops and children get zero indentation entries.

With the phantom, the walk goes: child's parent (`root`) has `root.parent == phantom` (not nil), so one indentation entry is added. Then `phantom.parent == nil` stops the loop. The phantom is never passed to `NewTree`, so it is never rendered — it exists purely to extend the parent chain by one level.

**When to use this:**
- Single visible root node whose children should be indented
- Tree header/title node wrapping all content
- Any case where you want root-level children to appear nested

**When you do NOT need this:**
- Multiple root nodes at the same level (typical file trees)
- Flat lists where indentation is unnecessary
- Trees where root children should align with the root

## License

MIT License - See LICENSE.txt
