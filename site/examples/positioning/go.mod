module github.com/mikeschinkel/go-tealeaves/site/examples/positioning

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teautils v0.0.0
	github.com/mikeschinkel/go-tealeaves/teafields v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teafields => ../../../teafields
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => ../../../../go-dt
)
