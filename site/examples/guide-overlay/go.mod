module github.com/mikeschinkel/go-tealeaves/site/examples/guide-overlay

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teaguide v0.0.0
	github.com/mikeschinkel/go-tealeaves/teautils v0.0.0
	charm.land/bubbles/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teaguide => ../../../teaguide
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
)
