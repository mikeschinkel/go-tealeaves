module github.com/mikeschinkel/go-tealeaves/site/examples/breadcrumb-nav

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teacrumbs v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teacrumbs => ../../../teacrumbs
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
)
