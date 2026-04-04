module github.com/mikeschinkel/go-tealeaves/site/examples/color-constants

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teacolor v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
