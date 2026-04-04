module github.com/mikeschinkel/go-tealeaves/site/examples/guide-quick-start

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teamodal v0.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teamodal => ../../../teamodal
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => /Users/mikeschinkel/Projects/go-pkgs/go-dt
)
