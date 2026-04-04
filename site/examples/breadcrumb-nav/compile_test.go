package examples_test

import (
	"strings"
	"testing"

	"github.com/mikeschinkel/go-tealeaves/teacrumbs"
	lipgloss "charm.land/lipgloss/v2"
)

// TestCompile_BreadcrumbsQuickExample verifies the quick example from breadcrumb-nav.mdx compiles.
func TestCompile_BreadcrumbsQuickExample(t *testing.T) {
	bc := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil)).
		Push(teacrumbs.NewCrumb("Settings", nil)).
		Push(teacrumbs.NewCrumb("Network", nil)).
		SetSize(80)

	_ = bc.View()
}

// TestCompile_CrumbType verifies the Crumb type and CrumbArgs from breadcrumb-nav.mdx.
func TestCompile_CrumbType(t *testing.T) {
	type somePayload struct{ val string }
	payload := somePayload{"test"}

	// Full form
	crumb := teacrumbs.NewCrumb("github.com/mikeschinkel/myapp", &teacrumbs.CrumbArgs{
		Short: "myapp",
		Data:  payload,
	})
	_ = crumb

	// Minimal form
	crumb2 := teacrumbs.NewCrumb("Home", nil)
	_ = crumb2
}

// TestCompile_DefaultStyles verifies DefaultStyles and WithStyles from breadcrumb-nav.mdx.
func TestCompile_DefaultStyles(t *testing.T) {
	bc := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil))

	styles := teacrumbs.DefaultStyles()
	styles.CurrentStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	bc = bc.WithStyles(styles)
	_ = bc
}

// TestCompile_RendererFunc verifies the RendererFunc type from breadcrumb-nav.mdx.
func TestCompile_RendererFunc(t *testing.T) {
	crumb := teacrumbs.NewCrumb("settings", &teacrumbs.CrumbArgs{
		Renderer: teacrumbs.RendererFunc(func(i int, m teacrumbs.BreadcrumbsModel) string {
			if i == m.Len()-1 {
				return strings.ToUpper(m.Crumbs()[i].Text)
			}
			return m.Crumbs()[i].Text
		}),
	})
	_ = crumb
}

// TestCompile_MouseSupport verifies SetPosition from breadcrumb-nav.mdx.
func TestCompile_MouseSupport(t *testing.T) {
	bc := teacrumbs.NewBreadcrumbsModel().
		Push(teacrumbs.NewCrumb("Home", nil))

	bc = bc.SetPosition(3, 0)
	_ = bc
}
