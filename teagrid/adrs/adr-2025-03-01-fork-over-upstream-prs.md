# ADR-2025-03-01: Fork from bubble-table Rather Than Upstream PRs

## Status
**Accepted**

## Date
2025-03-01

## Context
teagrid v0.2.0 originated as a fork of `github.com/evertras/bubble-table` by Brandon Fulljames (MIT License, Copyright 2022). During development, 17 structural problems were identified in the upstream project (documented in `CHARM_V1_TABLE_PROBLEMS.md`). A decision was needed on whether to address these problems by submitting pull requests to bubble-table or by maintaining a separate fork under a new name.

The 17 problems span the full surface area of the component:

1. Default right-alignment
2. No cell padding
3. Footer inherits baseStyle alignment
4. Invisible highlight default
5. No filter match highlighting
6. Confusing auto-select column
7. No column-level margin abstraction
8. FormatString doesn't apply to headers
9. baseStyle leaks into footer
10. No per-cell highlight in cell cursor mode
11. Heavy-handed overflow indicator
12. No cell/row editing
13. No sort key / display value separation
14. No rich/mixed coloration within a cell
15. No per-region border customization
16. Header cannot be fully removed
17. No borderless mode

## Decision
Fork the project and develop it independently as teagrid rather than submitting upstream pull requests to bubble-table.

## Rationale

### 1. Scope of changes is too large for incremental PRs
The rewrite addresses 17 documented structural problems. Many require changes to the core rendering pipeline, border system, and data model simultaneously. These cannot be submitted as independent PRs because they are architecturally intertwined.

### 2. Different design philosophy
bubble-table pre-computes 14+ border styles at init and bakes cursor styling into row data (O(n) rebuild per keystroke). teagrid v0.2.0 uses region-based borders with literal string rendering and render-time cursor styling (zero row rebuilding). These are fundamentally different architectural approaches, not bug fixes.

### 3. Charm v2 migration
The rewrite targets Charm v2 exclusively (`charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`). bubble-table is built on Charm v1. A PR that migrates to v2 while simultaneously restructuring the entire rendering pipeline would be unreviewable.

### 4. New features not in bubble-table's scope
CellValue with SortValue separation, Span-based rich text, CellStyleFunc with cursor awareness, independent footer zones, configurable overflow -- these represent a different vision for what a table component should be.

### 5. Breaking API changes
v0.2.0 changes method signatures (Update returns concrete Model, View returns tea.View), renames types (StyledCell to CellValue), changes defaults (left-aligned, paddingLeft=1, no auto-select column), and removes deprecated patterns. These would break all existing bubble-table users.

### 6. Maintainer relationship
Submitting 17+ PRs that collectively rewrite the project would effectively ask the maintainer to adopt someone else's vision for their project. A fork is more respectful of the original maintainer's design choices.

## Consequences

### Positive
- Freedom to make sweeping architectural changes without coordinating with upstream.
- Clean migration to Charm v2 without backward-compatibility constraints.
- Ability to introduce breaking API changes that improve ergonomics (left-aligned defaults, required explicit column selection, render-time styling).
- Independent release cadence aligned with teagrid's own roadmap.
- The original project remains intact for its existing users.

### Requirements
- Must maintain MIT license attribution to Brandon Fulljames (Copyright 2022) as required by the original license.
- Must clearly document the fork's origin in README and LICENSE.
- Must document the 17 upstream problems to justify the fork's existence and guide development priorities.

## Alternatives Considered

### Submit individual PRs for each problem
Rejected. Many of the 17 problems are architecturally intertwined. Fixing default alignment requires changes to the rendering pipeline; fixing border behavior requires restructuring the border system; fixing cursor styling requires rethinking the data model. These cannot be isolated into independent, reviewable PRs.

### Submit a single massive PR
Rejected. A PR that rewrites the rendering pipeline, migrates to Charm v2, restructures the border system, changes default behaviors, and renames types would be unreviewable. It would also be disrespectful to the maintainer -- effectively asking them to replace their project with someone else's rewrite.

### Wrap bubble-table with an adapter layer
Rejected. The rendering pipeline issues (border pre-computation, cursor styling baked into row data, baseStyle leaking into footer) cannot be fixed from outside the component. An adapter could paper over API ergonomics but not structural rendering problems.

### Contribute to Charm's official table (bubbles/table)
Rejected. The official Charm table component is even more limited than bubble-table and suffers from the same category of structural issues. It would require the same scope of changes with the additional overhead of working within the Charm organization's review process.

## Summary
The scope, depth, and interconnected nature of the required changes -- combined with a Charm v2 migration, fundamentally different rendering architecture, and breaking API changes -- make a fork the only practical path. Submitting upstream PRs would produce either an unreviewable monolithic PR or a series of intertwined PRs that cannot be merged independently. A fork respects the original maintainer's design choices while giving teagrid the freedom to pursue a different architectural vision.
