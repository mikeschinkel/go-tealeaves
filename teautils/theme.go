package teautils

import (
	"charm.land/lipgloss/v2"
)

// Theme holds derived lipgloss.Style values for each component, built from
// a Palette. Components accept Theme via WithTheme() to override their defaults.
type Theme struct {
	Palette Palette

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
	StatusBar StatusBarTheme
	HelpVisor HelpVisorTheme
	Modal     ModalTheme
	Dropdown  DropdownTheme
	List      ListTheme
	Grid      GridTheme
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

// NewTheme builds a Theme from a Palette, deriving all styles from the
// palette's semantic color slots.
func NewTheme(p Palette) Theme {
	return Theme{
		Palette: p,

		// Common styles
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Border),
		BorderAccent: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.BorderAccent),
		Title: lipgloss.NewStyle().
			Foreground(p.Accent).
			Bold(true),
		Message: lipgloss.NewStyle().
			Foreground(p.TextSecondary),
		Button: lipgloss.NewStyle().
			Foreground(p.ButtonFg),
		FocusedButton: lipgloss.NewStyle().
			Foreground(p.ButtonFocusFg).
			Background(p.ButtonFocusBg),
		Item: lipgloss.NewStyle().
			Foreground(p.TextPrimary),
		SelectedItem: lipgloss.NewStyle().
			Foreground(p.SelectionFg).
			Background(p.SelectionBg),
		ActiveItem: lipgloss.NewStyle().
			Foreground(p.Accent),

		// Status bar
		StatusBar: StatusBarTheme{
			MenuKeyStyle: lipgloss.NewStyle().
				Foreground(p.AccentSubtle).
				Bold(true),
			MenuLabelStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
			IndicatorStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
			IndicatorSepStyle: lipgloss.NewStyle().
				Foreground(p.TextDim),
			BarStyle: lipgloss.NewStyle(),
		},

		// Help visor
		HelpVisor: HelpVisorTheme{
			TitleStyle: lipgloss.NewStyle().
				Bold(true).
				Foreground(p.AccentAlt),
			CategoryStyle: lipgloss.NewStyle().
				Bold(true).
				Foreground(p.AccentAlt),
			KeyStyle: lipgloss.NewStyle().
				Foreground(p.AccentSubtle).
				Bold(true),
			DescStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
		},

		// Modal
		Modal: ModalTheme{
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(p.BorderAccent),
			TitleStyle: lipgloss.NewStyle().
				Foreground(p.Accent).
				Bold(true),
			MessageStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
			ButtonStyle: lipgloss.NewStyle().
				Foreground(p.ButtonFg),
			FocusedButtonStyle: lipgloss.NewStyle().
				Foreground(p.ButtonFocusFg).
				Background(p.ButtonFocusBg),
			CancelKeyStyle: lipgloss.NewStyle().
				Foreground(p.Accent),
			CancelTextStyle: lipgloss.NewStyle().
				Foreground(p.TextMuted),
		},

		// Dropdown
		Dropdown: DropdownTheme{
			BorderStyle: lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(p.Border),
			ItemStyle: lipgloss.NewStyle().
				Foreground(p.TextPrimary),
			SelectedStyle: lipgloss.NewStyle().
				Foreground(p.SelectionFg).
				Background(p.SelectionBg),
		},

		// List
		List: ListTheme{
			ItemStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
			SelectedItemStyle: lipgloss.NewStyle().
				Foreground(p.SelectionFg).
				Background(p.SelectionBg),
			ActiveItemStyle: lipgloss.NewStyle().
				Foreground(p.Accent),
			FooterStyle: lipgloss.NewStyle().
				Foreground(p.TextDim),
			StatusStyle: lipgloss.NewStyle().
				Foreground(p.StatusWarn),
			EditItemStyle: lipgloss.NewStyle().
				Foreground(p.EditFg).
				Background(p.EditBg),
			ScrollbarStyle: lipgloss.NewStyle().
				Foreground(p.ScrollTrack),
			ScrollThumbStyle: lipgloss.NewStyle().
				Foreground(p.ScrollThumb),
		},

		// Grid
		Grid: GridTheme{
			HeaderStyle: lipgloss.NewStyle().
				Bold(true).
				Foreground(p.TextPrimary),
			BaseStyle: lipgloss.NewStyle().
				Foreground(p.TextSecondary),
			HighlightStyle: lipgloss.NewStyle().
				Background(p.SelectionBg).
				Foreground(p.SelectionFg),
			BorderStyle: lipgloss.NewStyle().
				Foreground(p.Border),
		},
	}
}

// DefaultTheme returns NewTheme(DefaultPalette()).
func DefaultTheme() Theme {
	return NewTheme(DefaultPalette())
}
