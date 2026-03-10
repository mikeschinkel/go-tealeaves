package teautils

import (
	"charm.land/lipgloss/v2"
)

// Theme holds derived lipgloss.Style values for each component, built from
// a SystemPalette. Components accept Theme via WithTheme() to override their
// defaults.
type Theme struct {
	System SystemPalette

	// Common styles derived from palette
	Border        lipgloss.Style
	BorderAccent  lipgloss.Style
	Title         lipgloss.Style
	Message       lipgloss.Style
	Button        lipgloss.Style
	FocusedButton lipgloss.Style
	Item          lipgloss.Style
	SelectedItem  lipgloss.Style
	ActiveItem    lipgloss.Style

	// Component-specific style groups
	Breadcrumb BreadcrumbTheme
	StatusBar  StatusBarTheme
	HelpVisor  HelpVisorTheme
	Modal      ModalTheme
	Dropdown   DropdownTheme
	List       ListTheme
	Grid       GridTheme
}

// BreadcrumbTheme holds styles for the teacrumbs breadcrumb component.
type BreadcrumbTheme struct {
	ParentStyle    lipgloss.Style
	CurrentStyle   lipgloss.Style
	SeparatorStyle lipgloss.Style
	HoverStyle     lipgloss.Style
}

// StatusBarTheme holds styles for the teastatus status bar component.
type StatusBarTheme struct {
	MenuKeyStyle      lipgloss.Style
	MenuLabelStyle    lipgloss.Style
	IndicatorStyle    lipgloss.Style
	IndicatorSepStyle lipgloss.Style
	BarStyle          lipgloss.Style
}

// HelpVisorTheme holds styles for the help visor in teautils.
type HelpVisorTheme struct {
	TitleStyle    lipgloss.Style
	CategoryStyle lipgloss.Style
	KeyStyle      lipgloss.Style
	DescStyle     lipgloss.Style
}

// ModalTheme holds styles for the teamodal modal dialog component.
type ModalTheme struct {
	BorderStyle        lipgloss.Style
	TitleStyle         lipgloss.Style
	MessageStyle       lipgloss.Style
	ButtonStyle        lipgloss.Style
	FocusedButtonStyle lipgloss.Style
	CancelKeyStyle     lipgloss.Style
	CancelTextStyle    lipgloss.Style
}

// DropdownTheme holds styles for the teadrpdwn dropdown component.
type DropdownTheme struct {
	BorderStyle   lipgloss.Style
	ItemStyle     lipgloss.Style
	SelectedStyle lipgloss.Style
}

// ListTheme holds styles for the teamodal list component.
type ListTheme struct {
	ItemStyle         lipgloss.Style
	SelectedItemStyle lipgloss.Style
	ActiveItemStyle   lipgloss.Style
	FooterStyle       lipgloss.Style
	StatusStyle       lipgloss.Style
	EditItemStyle     lipgloss.Style
	ScrollbarStyle    lipgloss.Style
	ScrollThumbStyle  lipgloss.Style
}

// GridTheme holds styles for the teagrid grid component.
type GridTheme struct {
	HeaderStyle    lipgloss.Style
	BaseStyle      lipgloss.Style
	HighlightStyle lipgloss.Style
	BorderStyle    lipgloss.Style
}

// NewTheme builds a Theme from a SystemPalette, deriving all styles from the
// palette's semantic color slots. Single-color styles use SemanticColor's cached
// methods; multi-property styles chain from the cached base.
func NewTheme(sys SystemPalette) Theme {
	p := sys
	return Theme{
		System: sys,

		// Common styles
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border),
		BorderAccent: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.BorderAccent),
		Title:         p.Accent.Foreground().Bold(true),
		Message:       p.TextSecondary.Foreground(),
		Button:        p.ButtonFg.Foreground(),
		FocusedButton: p.ButtonFocusFg.Foreground().Background(p.ButtonFocusBg),
		Item:          p.TextPrimary.Foreground(),
		SelectedItem:  p.SelectionFg.Foreground().Background(p.SelectionBg),
		ActiveItem:    p.Accent.Foreground(),

		// Breadcrumb
		Breadcrumb: BreadcrumbTheme{
			ParentStyle:    p.TextMuted.Foreground(),
			CurrentStyle:   p.Accent.Foreground().Bold(true),
			SeparatorStyle: p.AccentSubtle.Foreground(),
			HoverStyle:     p.Accent.Foreground().Underline(true),
		},

		// Status bar
		StatusBar: StatusBarTheme{
			MenuKeyStyle:      p.AccentSubtle.Foreground().Bold(true),
			MenuLabelStyle:    p.TextSecondary.Foreground(),
			IndicatorStyle:    p.TextSecondary.Foreground(),
			IndicatorSepStyle: p.TextDim.Foreground(),
			BarStyle:          lipgloss.NewStyle(),
		},

		// Help visor
		HelpVisor: HelpVisorTheme{
			TitleStyle:    p.AccentAlt.Foreground().Bold(true),
			CategoryStyle: p.AccentAlt.Foreground().Bold(true),
			KeyStyle:      p.AccentSubtle.Foreground().Bold(true),
			DescStyle:     p.TextSecondary.Foreground(),
		},

		// Modal
		Modal: ModalTheme{
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(p.BorderAccent),
			TitleStyle:         p.Accent.Foreground().Bold(true),
			MessageStyle:       p.TextSecondary.Foreground(),
			ButtonStyle:        p.ButtonFg.Foreground(),
			FocusedButtonStyle: p.ButtonFocusFg.Foreground().Background(p.ButtonFocusBg),
			CancelKeyStyle:     p.Accent.Foreground(),
			CancelTextStyle:    p.TextMuted.Foreground(),
		},

		// Dropdown
		Dropdown: DropdownTheme{
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(p.Border),
			ItemStyle:     p.TextPrimary.Foreground(),
			SelectedStyle: p.SelectionFg.Foreground().Background(p.SelectionBg),
		},

		// List
		List: ListTheme{
			ItemStyle:         p.TextSecondary.Foreground(),
			SelectedItemStyle: p.SelectionFg.Foreground().Background(p.SelectionBg),
			ActiveItemStyle:   p.Accent.Foreground(),
			FooterStyle:       p.TextDim.Foreground(),
			StatusStyle:       p.StatusWarn.Foreground(),
			EditItemStyle:     p.EditFg.Foreground().Background(p.EditBg),
			ScrollbarStyle:    p.ScrollTrack.Foreground(),
			ScrollThumbStyle:  p.ScrollThumb.Foreground(),
		},

		// Grid
		Grid: GridTheme{
			HeaderStyle:    p.TextPrimary.Foreground().Bold(true),
			BaseStyle:      p.TextSecondary.Foreground(),
			HighlightStyle: p.SelectionFg.Foreground().Background(p.SelectionBg),
			BorderStyle:    p.Border.Foreground(),
		},
	}
}

// DefaultTheme returns NewTheme(DefaultSystemPalette()).
func DefaultTheme() Theme {
	return NewTheme(DefaultSystemPalette(nil))
}
