/*
Package tealayout provides a declarative layout engine for Bubble Tea v2 and
LipGloss v2 terminal UIs.

tealayout replaces manual width/height arithmetic with dimension-based
resolution: Percent(n), Fixed(n), Flex(weight), and Fit() dimensions within
Row and Column panes. The engine handles space distribution, integer
rounding, min/max clamping, optional child collapse, and the "two-width
problem" (total vs content dimensions) automatically.

Widgets are wrapped with NewElement[T] for type-safe access:

	tree := tealayout.NewElement(newTreeWidget())
	code := tealayout.NewElement(newCodeWidget())

	root := tealayout.NewColumn(tealayout.Percent100,
		tealayout.NewRow(tealayout.Fixed(1), header),
		tealayout.NewRow(tealayout.Flex(1),
			tealayout.NewColumn(tealayout.Flex(0.25), tree),
			tealayout.NewColumn(tealayout.Flex(0.75), code),
		),
	)
	layout := tealayout.NewLayout(root)
	layout.SetSize(termWidth, termHeight)
	output, err := layout.Render()

	// Type-safe widget access:
	tree.Widget().SomeMethod()

# Stability

This package is provisional as of v0.3.0. The public API may change in
minor releases until promoted to stable.
*/
package tealayout
