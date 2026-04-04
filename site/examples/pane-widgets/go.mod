module github.com/mikeschinkel/go-tealeaves/site/examples/pane-widgets

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teapane v0.0.0
	github.com/mikeschinkel/go-tealeaves/tealayout v0.0.0
	github.com/mikeschinkel/go-tealeaves/teacolor v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teapane => ../../../teapane
	github.com/mikeschinkel/go-tealeaves/tealayout => ../../../tealayout
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-tealeaves/teacrumbs => ../../../teacrumbs
	github.com/mikeschinkel/go-dt => ../../../../go-dt
)
