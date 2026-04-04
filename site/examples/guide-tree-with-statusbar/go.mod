module github.com/mikeschinkel/go-tealeaves/site/examples/guide-tree-with-statusbar

go 1.25.7

require (
	github.com/mikeschinkel/go-tealeaves/teatree v0.0.0
	github.com/mikeschinkel/go-tealeaves/teastatus v0.0.0
	github.com/mikeschinkel/go-dt v0.5.0
	charm.land/bubbles/v2 v2.0.0
	charm.land/lipgloss/v2 v2.0.0
)

replace (
	github.com/mikeschinkel/go-tealeaves/teatree => ../../../teatree
	github.com/mikeschinkel/go-tealeaves/teastatus => ../../../teastatus
	github.com/mikeschinkel/go-tealeaves/teafields => ../../../teafields
	github.com/mikeschinkel/go-tealeaves/teautils => ../../../teautils
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../../teacolor
	github.com/mikeschinkel/go-dt => /Users/mikeschinkel/Projects/go-pkgs/go-dt
)
