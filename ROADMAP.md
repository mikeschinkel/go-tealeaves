# Roadmap

## Planned Components

### teaform — Bubble Tea Form Component
**Status:** 🔴 Not Started
**Scope:** Large — general-purpose reusable component

**Problem**: There is no general-purpose form component for Bubble Tea. Building any settings UI, wizard, or multi-field input requires ad-hoc assembly of individual inputs with custom layout and navigation logic each time.

**Proposed Solution**: A reusable `teaform` package that provides declarative form building with field types, validation, layout, and keyboard navigation.

**Potential Field Types**:
- Text input (single-line, with optional validation)
- Text area (multi-line)
- Select / dropdown (single choice)
- Multi-select (checkboxes)
- Toggle (boolean)
- Number input (with min/max)
- File/directory path (with completion)

**Key Design Considerations**:
- Declarative field definitions (define fields, get a working form)
- Tab/Shift+Tab navigation between fields
- Per-field validation with inline error display
- Support for field groups / sections
- Scrollable when form exceeds viewport height
- Composable — embeddable within other Bubble Tea models

**Open Questions**:
- Should it use struct tags (like HTML form libraries) or a builder API?
- How to handle dynamic forms (fields that appear/disappear based on other field values)?
- What's the right abstraction for custom field types?
- Should it integrate with go-cfgstore schema definitions?

**Known Use Cases**:
- Gomion configuration editor TUI
- Commit message templates
- Project scaffolding wizards
- Filter/search dialogs

**When to Implement**:
- When a concrete use case (like the Gomion Configuration Editor) justifies the investment
