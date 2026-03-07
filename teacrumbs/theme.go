package teacrumbs

import (
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with styles derived from the given theme.
// Individual WithStyles() calls take precedence over theme if called after.
func (m BreadcrumbsModel) WithTheme(theme teautils.Theme) BreadcrumbsModel {
	m.Styles = Styles{
		ParentStyle:    theme.Breadcrumb.ParentStyle,
		CurrentStyle:   theme.Breadcrumb.CurrentStyle,
		SeparatorStyle: theme.Breadcrumb.SeparatorStyle,
	}
	return m
}
