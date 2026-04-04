module github.com/mikeschinkel/go-tealeaves/site/examples/text-selection

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teatext v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teatext => ../../../teatext
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
