module github.com/mikeschinkel/go-tealeaves/site/examples/diff-renderer

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teadiffr v0.0.0
	github.com/mikeschinkel/go-tealeaves/teautils v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teadiffr => ../../../teadiffr
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
