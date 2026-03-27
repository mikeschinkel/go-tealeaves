/*
Package tealayout provides a declarative layout engine for Bubble Tea v2 and
LipGloss v2 terminal UIs.

tealayout replaces manual width/height arithmetic with dimension-based
resolution: Percent(n), Fixed(n), Flex(weight), and Fit() dimensions within
Row and Column components. The engine handles space distribution, integer
rounding, min/max clamping, optional child collapse, and the "two-width
problem" (total vs content dimensions) automatically.

Basic usage:

	root := tealayout.NewRow(tealayout.Percent100,
		tealayout.NewColumn(tealayout.Fit(), tree),
		tealayout.NewColumn(tealayout.Flex(1.0), code),
		tealayout.NewColumn(tealayout.Flex(1.0), diff),
	)
	layout := tealayout.NewLayout(root)
	layout.SetSize(termWidth, termHeight)
	output, err := layout.Render()

# Stability

This package is provisional as of v0.2.0. The public API may change in
minor releases until promoted to stable.
*/
package tealayout
