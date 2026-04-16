// Source: site/src/content/docs/components/theming.mdx:26#295a7232,43#5ff1e885,66#c34175e4,93#517778d7,115#3cd30429
package examples_test

import (
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teacolor"
	"github.com/mikeschinkel/go-tealeaves/teagrid"
	"github.com/mikeschinkel/go-tealeaves/teamodal"
	"github.com/mikeschinkel/go-tealeaves/teastatus"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// TestCompile_ThemingQuickStart verifies the quick start from theming.mdx.
func TestCompile_ThemingQuickStart(t *testing.T) {
	// Use adaptive palette (auto-detects dark/light terminal)
	theme := teautils.DefaultTheme()

	// Apply to any component
	columns := []teagrid.Column{teagrid.NewColumn("id", "ID", 10)}
	grid := teagrid.NewGridModel(columns).WithTheme(theme)
	_ = grid

	statusBar := teastatus.NewStatusBarModel().WithTheme(theme)
	_ = statusBar

	modal := teamodal.NewYesNoModal("Confirm?", nil).WithTheme(theme)
	_ = modal
}

// TestCompile_PaletteConstructors verifies palette constructors from theming.mdx.
func TestCompile_PaletteConstructors(t *testing.T) {
	dark := teautils.DarkSystemPalette(nil)
	_ = dark

	light := teautils.LightSystemPalette(nil)
	_ = light

	adaptive := teautils.AdaptiveSystemPalette(nil)
	_ = adaptive

	def := teautils.DefaultSystemPalette(nil)
	_ = def

	// With PaletteOpts
	opts := &teautils.PaletteOpts{Adaptive: true}
	withOpts := teautils.DarkSystemPalette(opts)
	_ = withOpts
}

// TestCompile_PaletteGeneric verifies generic Palette[T] from theming.mdx.
func TestCompile_PaletteGeneric(t *testing.T) {
	type AppColors struct {
		Brand   teacolor.SemanticColor
		Sidebar teacolor.SemanticColor
	}

	palette := teautils.Palette[AppColors]{
		System: teautils.DarkSystemPalette(nil),
		App: AppColors{
			Brand:   teacolor.NewSemanticColor(teacolor.Color46),
			Sidebar: teacolor.NewSemanticColor(teacolor.Color238),
		},
	}

	// Pass palette.System to NewTheme for component theming
	theme := teautils.NewTheme(palette.System)
	_ = theme

	// Access app colors directly
	_ = palette.App.Brand
	_ = palette.App.Sidebar
}

// TestCompile_ThemeConstructors verifies Theme constructors from theming.mdx.
func TestCompile_ThemeConstructors(t *testing.T) {
	palette := teautils.DarkSystemPalette(nil)
	theme := teautils.NewTheme(palette)
	_ = theme

	defaultTheme := teautils.DefaultTheme()
	_ = defaultTheme
}

// TestCompile_ThemeFields verifies Theme field access from theming.mdx.
func TestCompile_ThemeFields(t *testing.T) {
	theme := teautils.DefaultTheme()

	// Common styles
	_ = theme.System
	_ = theme.Border
	_ = theme.BorderAccent
	_ = theme.Title
	_ = theme.Message
	_ = theme.Button
	_ = theme.FocusedButton
	_ = theme.Item
	_ = theme.SelectedItem
	_ = theme.ActiveItem

	// Component-specific themes
	_ = theme.Breadcrumb
	_ = theme.StatusBar
	_ = theme.HelpVisor
	_ = theme.Modal
	_ = theme.Dropdown
	_ = theme.List
	_ = theme.Grid
}

// TestCompile_BreadcrumbTheme verifies BreadcrumbTheme fields from theming.mdx.
func TestCompile_BreadcrumbTheme(t *testing.T) {
	theme := teautils.DefaultTheme()
	_ = theme.Breadcrumb.ParentStyle
	_ = theme.Breadcrumb.CurrentStyle
	_ = theme.Breadcrumb.SeparatorStyle
	_ = theme.Breadcrumb.HoverStyle
}

// TestCompile_StatusBarTheme verifies StatusBarTheme fields from theming.mdx.
func TestCompile_StatusBarTheme(t *testing.T) {
	theme := teautils.DefaultTheme()
	_ = theme.StatusBar.MenuKeyStyle
	_ = theme.StatusBar.MenuLabelStyle
	_ = theme.StatusBar.IndicatorStyle
	_ = theme.StatusBar.IndicatorSepStyle
	_ = theme.StatusBar.BarStyle
}

// TestCompile_ModalTheme verifies ModalTheme fields from theming.mdx.
func TestCompile_ModalTheme(t *testing.T) {
	theme := teautils.DefaultTheme()
	_ = theme.Modal.BorderStyle
	_ = theme.Modal.TitleStyle
	_ = theme.Modal.MessageStyle
	_ = theme.Modal.ButtonStyle
	_ = theme.Modal.FocusedButtonStyle
	_ = theme.Modal.CancelKeyStyle
	_ = theme.Modal.CancelTextStyle
}

// TestCompile_GridTheme verifies GridTheme fields from theming.mdx.
func TestCompile_GridTheme(t *testing.T) {
	theme := teautils.DefaultTheme()
	_ = theme.Grid.HeaderStyle
	_ = theme.Grid.BaseStyle
	_ = theme.Grid.HighlightStyle
	_ = theme.Grid.BorderStyle
}
