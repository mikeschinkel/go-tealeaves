module github.com/mikeschinkel/go-tealeaves/site/examples/guide-notification-after-action

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teanotify v0.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teanotify => ../../../teanotify
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => /Users/mikeschinkel/Projects/go-pkgs/go-dt
)
