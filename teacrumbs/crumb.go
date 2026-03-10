package teacrumbs

import "charm.land/lipgloss/v2"

// Crumb represents a single breadcrumb in the crumbs.
type Crumb struct {
	// Text is the primary display text (full form, e.g., "github.com:mikeschinkel/gomion").
	Text string

	// Short is the compact text for non-current crumbs (e.g., "gomion").
	// Falls back to Text if empty.
	Short string

	// Style is a per-crumb style override. When non-nil, it takes precedence
	// over the model's ParentStyle/CurrentStyle.
	Style *lipgloss.Style

	// Renderer generates dynamic display text for this crumb. When non-nil,
	// it is called instead of the standard styling pipeline.
	Renderer Renderer

	// Data carries an app-defined payload available in click/hover messages.
	// Same pattern as context.Value — the breadcrumb model never inspects it.
	Data any
}

type CrumbArgs struct {
	// Short is the compact text for non-current crumbs (e.g., "gomion").
	// Falls back to Text if empty.
	Short string

	// Style is a per-crumb style override. When non-nil, it takes precedence
	// over the model's ParentStyle/CurrentStyle.
	Style *lipgloss.Style

	// Renderer generates dynamic display text for this crumb. When non-nil,
	// it is called instead of the standard styling pipeline.
	Renderer Renderer

	// Data carries an app-defined payload available in click/hover messages.
	// Same pattern as context.Value — the breadcrumb model never inspects it.
	Data any
}

var nullArgs = &CrumbArgs{}

func NewCrumb(text string, args *CrumbArgs) Crumb {
	if args == nil {
		args = nullArgs
	}
	return Crumb{
		Text:     text,
		Short:    args.Short,
		Style:    args.Style,
		Renderer: args.Renderer,
		Data:     args.Data,
	}
}
