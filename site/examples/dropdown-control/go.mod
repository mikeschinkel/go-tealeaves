module github.com/mikeschinkel/go-tealeaves/site/examples/dropdown-control

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teafields v0.0.0
	charm.land/bubbletea/v2 v2.0.0
	charm.land/bubbles/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teafields => ../../../teafields
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => ../../../../go-dt
)
