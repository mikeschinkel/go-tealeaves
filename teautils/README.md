# teautils

Utility components for Bubble Tea applications.

## Key Registry System

The key registry provides a way for apps to control how key bindings are presented across different contexts (status bar, help modal, etc.) while allowing components to use standard `key.Binding` types.

### Design Philosophy

- **Components define keys** using standard `key.Binding` - keeps them compatible with the Bubble Tea ecosystem
- **App controls presentation** via the registry - decides what shows in status bar vs help modal
- **Separation of concerns** - presentation is an app-level decision, not a key property

### Components

#### KeyIdentifier

A validated identifier for keys in the registry.

- Format: `"namespace.key"` (e.g., `"app.help"`, `"file-intent.commit"`)
- Allows alphanumeric, hyphens, underscores
- Must contain at least one dot separator

```go
id := teautils.MustParseKeyIdentifier("file-intent.commit")
```

#### KeyRegistry

Manages key bindings and their display metadata.

```go
registry := teautils.NewKeyRegistry()

registry.MustRegister(teautils.KeyMeta{
    ID:        teautils.MustParseKeyIdentifier("app.help"),
    Binding:   appKeys.Help,
    StatusBar: true,
    HelpModal: true,
    Category:  "System",
    HelpText:  "Toggle this help screen",
})

// Get keys for different contexts
statusBarKeys := registry.ForStatusBar()
helpKeys := registry.ByCategory()
```

#### Help Visor Rendering

Renders a comprehensive help overlay organized by category.

```go
helpVisor := teautils.RenderHelpVisorDefault(registry.ByCategory())
```

### Usage Pattern

**Component (uses standard key.Binding):**
```go
type FileIntentKeyMap struct {
    Back   key.Binding
    Commit key.Binding
}

func (m FileIntentModel) ShortHelpBindings() []key.Binding {
    return []key.Binding{m.Keys.Back, m.Keys.Commit}
}
```

**App (controls presentation):**
```go
func (m *AppModel) setupKeys() {
    m.keyRegistry = teautils.NewKeyRegistry()

    m.keyRegistry.MustRegisterMany([]teautils.KeyMeta{
        {
            ID:        teautils.MustParseKeyIdentifier("app.help"),
            Binding:   m.Keys.Help,
            StatusBar: true,
            HelpModal: true,
            Category:  "System",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.commit"),
            Binding:   m.fileIntent.Keys.Commit,
            StatusBar: false, // Not in status bar
            HelpModal: true,
            Category:  "File Actions",
            HelpText:  "Mark file for inclusion in next commit",
        },
    })
}
```

### Integration Example (gomtui)

This shows how to integrate the KeyRegistry system into gomtui to replace the current help implementation with the visor-up model.

#### Step 1: Add KeyRegistry to AppModel

```go
// In gomtui/app_model.go

type AppModel struct {
    // ... existing fields ...

    keyRegistry *teautils.KeyRegistry
    showHelp    bool // Toggle for help visor
}
```

#### Step 2: Initialize Registry in NewAppModel

```go
func NewAppModel(args AppModelArgs) AppModel {
    m := AppModel{
        // ... existing initialization ...
        keyRegistry: teautils.NewKeyRegistry(),
        showHelp:    false,
    }

    // Register app-level keys
    m.registerAppKeys()

    return m
}

func (m *AppModel) registerAppKeys() {
    m.keyRegistry.MustRegisterMany([]teautils.KeyMeta{
        {
            ID:        teautils.MustParseKeyIdentifier("app.help"),
            Binding:   m.Keys.Help,
            StatusBar: true,  // Always visible
            HelpModal: true,
            Category:  "System",
            HelpText:  "Toggle this help screen",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("app.quit"),
            Binding:   m.Keys.Quit,
            StatusBar: false, // Not in status bar
            HelpModal: true,
            Category:  "System",
            HelpText:  "Quit application (from top level)",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("app.force-quit"),
            Binding:   m.Keys.ForceQuit,
            StatusBar: false,
            HelpModal: true,
            Category:  "System",
            HelpText:  "Force quit immediately",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("app.refresh"),
            Binding:   m.Keys.RefreshData,
            StatusBar: false,
            HelpModal: true,
            Category:  "System",
            HelpText:  "Refresh cached data",
        },
    })
}
```

#### Step 3: Register View-Specific Keys When View Changes

```go
// In gomtui/app_model.go

func (m *AppModel) pushView(view UIViewer) tea.Cmd {
    // ... existing push logic ...

    // Register keys for this view
    m.registerViewKeys(view)

    return nil
}

func (m *AppModel) registerViewKeys(view UIViewer) {
    // Clear previous view's keys (keep app-level keys)
    m.keyRegistry.Clear()
    m.registerAppKeys() // Re-register app keys

    // Register view-specific keys
    switch v := view.(type) {
    case *FileIntentModel:
        m.registerFileIntentKeys(v)
    case *CommitTargetModel:
        m.registerCommitTargetKeys(v)
    }
}

func (m *AppModel) registerFileIntentKeys(v *FileIntentModel) {
    m.keyRegistry.MustRegisterMany([]teautils.KeyMeta{
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.back"),
            Binding:   v.Keys.Back,
            StatusBar: true, // Essential navigation
            HelpModal: true,
            Category:  "Navigation",
            HelpText:  "Return to commit targets view",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.tab"),
            Binding:   v.Keys.TogglePane,
            StatusBar: true, // Essential navigation
            HelpModal: true,
            Category:  "Navigation",
            HelpText:  "Switch between tree and content panes",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.commit"),
            Binding:   v.Keys.Commit,
            StatusBar: false, // Too many for status bar
            HelpModal: true,
            Category:  "File Actions",
            HelpText:  "Mark file for inclusion in next commit",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.omit"),
            Binding:   v.Keys.Omit,
            StatusBar: false,
            HelpModal: true,
            Category:  "File Actions",
            HelpText:  "Omit file from commit (no action)",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.ignore"),
            Binding:   v.Keys.Ignore,
            StatusBar: false,
            HelpModal: true,
            Category:  "File Actions",
            HelpText:  "Add file pattern to .gitignore",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.exclude"),
            Binding:   v.Keys.Exclude,
            StatusBar: false,
            HelpModal: true,
            Category:  "File Actions",
            HelpText:  "Add file pattern to .git/info/exclude",
        },
        {
            ID:        teautils.MustParseKeyIdentifier("file-intent.save"),
            Binding:   v.Keys.Save,
            StatusBar: false,
            HelpModal: true,
            Category:  "Actions",
            HelpText:  "Manually save commit plan to disk",
        },
    })
}
```

#### Step 4: Handle Help Toggle in Update

```go
// In gomtui/app_model.go Update()

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Handle help toggle BEFORE view updates
        if key.Matches(msg, m.Keys.Help) {
            m.showHelp = !m.showHelp
            return m, nil
        }

        // ... rest of update logic ...
    }

    // ... rest of update ...
}
```

#### Step 5: Render Status Bar and Help Visor

```go
// In gomtui/app_model.go View()

func (m AppModel) View() string {
    if m.Err != nil {
        return fmt.Sprintf("Error: %v\n", m.Err)
    }

    // Get main view from current viewer
    var mainView string
    if m.viewStack != nil && m.viewStack.Current() != nil {
        mainView = m.viewStack.Current().View()
    }

    // Render help visor if active
    if m.showHelp {
        helpVisor := teautils.RenderHelpVisorDefault(m.keyRegistry.ByCategory())
        mainView = m.overlayHelpVisor(mainView, helpVisor)
    }

    return mainView
}

// overlayHelpVisor overlays the help visor on the main view
func (m AppModel) overlayHelpVisor(mainView, helpVisor string) string {
    mainLines := strings.Split(mainView, "\n")
    helpLines := strings.Split(helpVisor, "\n")

    // Calculate position to center visor
    helpHeight := len(helpLines)
    mainHeight := len(mainLines)
    startLine := (mainHeight - helpHeight) / 2
    if startLine < 0 {
        startLine = 0
    }

    // Overlay help lines
    for i, helpLine := range helpLines {
        lineIdx := startLine + i
        if lineIdx < len(mainLines) {
            // Center horizontally
            helpWidth := lipgloss.Width(helpLine)
            leftPadding := (m.Width - helpWidth) / 2
            if leftPadding < 0 {
                leftPadding = 0
            }
            mainLines[lineIdx] = strings.Repeat(" ", leftPadding) + helpLine
        }
    }

    return strings.Join(mainLines, "\n")
}
```

#### Step 6: Remove Old Help System

Remove these files (now replaced by KeyRegistry):
- `gomtui/help_bindings.go` - No longer needed
- `gomtui/app_help.go` - Replaced by KeyRegistry methods
- `gomtui/file_intent_help.go` - Replaced by registration
- `gomtui/commit_target_help.go` - Replaced by registration

Remove from models:
- `help help.Model` field - No longer needed
- `ShortHelpBindings()` method - No longer needed
- `FullHelpBindings()` method - No longer needed

#### Benefits

1. **Cleaner status bar** - Only shows `[?] Help  [tab] Switch pane  [esc] Back`
2. **Better UX** - Full help available on demand via `?` key
3. **Centralized control** - App decides what goes where
4. **Reusable components** - Components still use standard `key.Binding`
5. **Colorful and engaging** - Styled keys and categories
6. **Context-aware** - Different keys for different views

#### Before/After

**Before** (current):
```
esc back to commit targets • tab switch pane • c mark for commit • o omit from commit • i add to .gitignore • x add to .git/info/exclude • ctrl+s save commit plan
```

**After** (with KeyRegistry):
```
[?] Help  [tab] Switch pane  [esc] Back
```

Press `?` to see full help organized by category with detailed descriptions.

### Future Extraction

This package is designed to be eventually extracted as `go-teautils` for use by other projects.
