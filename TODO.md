# go-tealeaves TODO

## PANE-NAME — Introduce typed PaneName across tealayout

Replace stringly-typed pane name parameters with a `type PaneName string` domain type (following the go-dt philosophy). This would catch pane name typos at compile time and make APIs self-documenting.

Scope: `PaneDef.Name`, `PaneNames()`, `FocusPane()`, `TogglePane()`, `SoloPane()`, `ShowPane()`, `HidePane()`, `SetPaneVisible()`, `VisibilityRotator` combos, and TCL's hardcoded `"tree"`/`"content"` strings. Consider whether `PaneName` belongs in tealayout or teautils.

## GRID-ALIGN — teagrid cell alignment support

Add alignment support to teagrid cells using tealayout's `Alignment` bitmap type.

Consider whether `Alignment` should live in teautils (shared foundation) rather than tealayout, since both teagrid and tealayout need it.
