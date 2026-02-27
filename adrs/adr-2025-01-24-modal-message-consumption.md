# ADR-2025-01-24: Modal Message Consumption via Command Return Pattern

## Status

**Accepted**

## Date

2025-01-24

## Context

The teadd dropdown component implements modal behavior: when open, it should consume keyboard messages (up/down/enter/esc) that would otherwise be processed by the parent.

The challenge: **How do we prevent "infection" of parent's `Update()` with modal state checks?**

Common anti-pattern:
```go
// Parent Update() - infected with modal checks
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    if !m.dropdown.IsOpen {  // ← State check "infection"
        switch msg := msg.(type) {
        case tea.KeyMsg:
            if msg.String() == "q" {
                return m, tea.Quit
            }
        }
    }
    return m, nil
}
```

This violates encapsulation - parent must know dropdown's internal state to decide whether to process messages.

## Decision

**Dropdown signals message consumption by returning a non-nil `tea.Cmd`.**

When dropdown doesn't handle a message (closed, or open but message not relevant), it returns `nil` cmd. Parent checks the returned cmd to determine if dropdown consumed the message.

### Dropdown Pattern

```go
func (m DropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var keyMsg tea.KeyMsg
    var ok bool

    if !m.IsOpen {
        goto end  // Not open = nil cmd = didn't handle
    }

    keyMsg, ok = msg.(tea.KeyMsg)
    if !ok {
        goto end  // Not a key message = nil cmd
    }

    switch keyMsg.String() {
    case "up", "k":
        m.Selected--
        if m.Selected < 0 {
            m.Selected = 0  // Clamp at top, don't wrap
        }
        cmd = func() tea.Msg { return nil }  // Non-nil cmd = handled
        goto end

    case "down", "j":
        m.Selected++
        if m.Selected >= len(m.Items) {
            m.Selected = len(m.Items) - 1  // Clamp at bottom, don't wrap
        }
        cmd = func() tea.Msg { return nil }  // Non-nil cmd = handled
        goto end

    case "enter":
        selected := ItemSelectedMsg{
            Index: m.Selected,
            Value: m.Items[m.Selected],
        }
        m.IsOpen = false
        cmd = func() tea.Msg { return selected }  // Non-nil cmd = handled
        goto end

    case "esc":
        m.IsOpen = false
        cmd = func() tea.Msg { return DropdownCancelledMsg{} }  // Non-nil cmd = handled
        goto end
    }

end:
    return m, cmd  // cmd = nil if not handled, non-nil if handled
}
```

### Parent Pattern

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd
    var dropdown tea.Model

    // Let dropdown handle first
    dropdown, cmd = m.dropdown.Update(msg)
    if cmd != nil {
        m.dropdown = dropdown.(teadd.DropdownModel)
        return m, cmd  // Dropdown consumed it - done!
    }

    // Dropdown didn't handle - parent processes
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q":
            return m, tea.Quit  // Only processed when dropdown didn't consume
        case " ":
            m.dropdown, cmd = m.dropdown.Open()
            return m, cmd
        }
    }

    return m, nil
}
```

## Rationale

### Why This Works

1. **True encapsulation**
   - Parent doesn't check `m.dropdown.IsOpen`
   - Dropdown decides what it consumes
   - Clean separation of concerns

2. **Explicit signal**
   - Non-nil cmd = "I handled this message"
   - Nil cmd = "This message isn't for me"
   - Simple, unambiguous contract

3. **Whitelist pattern**
   - Dropdown explicitly handles specific keys
   - Everything else falls through to parent
   - No blacklist needed

4. **Works with ClearPath**
   - `goto end` pattern makes delegation crystal clear
   - All variables declared upfront
   - Single return point

### Benefits Over State Checking

**Anti-pattern (state checking)**:
```go
if !m.dropdown.IsOpen {
    // Process parent keys
}
```

**Problems**:
* Parent must know dropdown state
* Infects parent's Update() with modal logic
* Violates encapsulation
* Harder to reason about

**This pattern (cmd checking)**:
```go
dropdown, cmd = m.dropdown.Update(msg)
if cmd != nil {
    return m, cmd  // Dropdown handled it
}
// Parent processes
```

**Advantages**:
* No state checking needed
* Parent code is clean
* Dropdown decides consumption
* Clear, obvious flow

## Consequences

### Positive

* **No infection** - Parent Update() remains clean
* **Encapsulation** - Dropdown's modal state is private
* **Obvious flow** - `if cmd != nil` clearly signals "message consumed"
* **Composable** - Pattern works with multiple modal components
* **ClearPath compatible** - `goto end` makes delegation explicit

### Requirements

* Dropdown must return non-nil cmd for **every** message it handles
* Even if cmd does nothing, it must be non-nil to signal consumption
* Parent must check `if cmd != nil` before processing message

## Alternatives Considered

### State Checking in Parent

```go
if !m.dropdown.IsOpen {
    // Process parent keys
}
```

**Rejected** because:
* Violates encapsulation
* Infects parent code
* Harder to maintain

### Returned Bool Flag

```go
func (m DropdownModel) Update(msg tea.Msg) (tea.Model, tea.Cmd, bool)
```

**Rejected** because:
* Requires API change (3-value return)
* `tea.Cmd` already serves as signal
* Less idiomatic

### Message Wrapper

```go
type ConsumedMsg struct{ Consumed bool }
```

**Rejected** because:
* Adds complexity
* Requires parent to check message type
* `tea.Cmd` return already available

## Summary

Using the returned `tea.Cmd` as a message consumption signal provides clean modal behavior without infecting the parent's Update() method with state checks. The pattern is simple, encapsulated, and works naturally with the ClearPath style.

---

*End of ADR*
