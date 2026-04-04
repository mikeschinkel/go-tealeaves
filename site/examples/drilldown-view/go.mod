module github.com/mikeschinkel/go-tealeaves/site/examples/drilldown-view

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teatree v0.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teatree => ../../../teatree
	github.com/mikeschinkel/go-tealeaves/teafields => ../../../teafields
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => ../../../../go-dt
)
