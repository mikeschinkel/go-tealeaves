module github.com/mikeschinkel/go-tealeaves/site/examples/choice-dialog

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teamodal v0.0.0
	github.com/mikeschinkel/go-tealeaves/teautils v0.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teamodal => ../../../teamodal
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
