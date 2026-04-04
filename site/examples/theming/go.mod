module github.com/mikeschinkel/go-tealeaves/site/examples/theming

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teautils v0.0.0
	github.com/mikeschinkel/go-tealeaves/teacolor v0.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
