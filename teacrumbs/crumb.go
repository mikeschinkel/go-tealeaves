package teacrumbs

// Crumb represents a single breadcrumb in the trail.
type Crumb struct {
	// Text is the display text for this breadcrumb.
	Text string

	// PreStyled indicates that Text already contains ANSI styling
	// and should not be styled further by the renderer.
	PreStyled bool
}
