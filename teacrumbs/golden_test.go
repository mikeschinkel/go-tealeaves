package teacrumbs

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/charmbracelet/x/exp/golden"
)

func TestBreadcrumbsModel_Golden_ThreeCrumbs(t *testing.T) {
	m := NewBreadcrumbsModel().WithStyles(DefaultStyles())
	m = m.SetSize(80)
	m = m.Push(Crumb{Text: "Home", Short: "~"})
	m = m.Push(Crumb{Text: "Projects", Short: "proj"})
	m = m.Push(Crumb{Text: "gomion", Short: "gom"})
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}

func TestBreadcrumbsModel_Golden_SingleCrumb(t *testing.T) {
	m := NewBreadcrumbsModel().WithStyles(DefaultStyles())
	m = m.SetSize(80)
	m = m.Push(Crumb{Text: "Root"})
	output := ansi.Strip(m.View().Content)
	golden.RequireEqual(t, []byte(output))
}
