module github.com/mikeschinkel/go-tealeaves/teafields/examples

go 1.25.7

replace (
	github.com/mikeschinkel/go-dt => ../../../go-dt
	github.com/mikeschinkel/go-dt/dtx => ../../../go-dt/dtx
	github.com/mikeschinkel/go-tealeaves/teacolor => ../../teacolor
	github.com/mikeschinkel/go-tealeaves/teafields => ../
	github.com/mikeschinkel/go-tealeaves/teautils => ../../teautils
)

require (
	charm.land/bubbletea/v2 v2.0.2
	charm.land/lipgloss/v2 v2.0.2
	github.com/charmbracelet/x/ansi v0.11.6
	github.com/mikeschinkel/go-tealeaves/teafields v0.0.0-00010101000000-000000000000
)

require (
	charm.land/bubbles/v2 v2.0.0 // indirect
	github.com/charmbracelet/colorprofile v0.4.3 // indirect
	github.com/charmbracelet/ultraviolet v0.0.0-20260223171050-89c142e4aa73 // indirect
	github.com/charmbracelet/x/term v0.2.2 // indirect
	github.com/charmbracelet/x/termios v0.1.1 // indirect
	github.com/charmbracelet/x/windows v0.2.2 // indirect
	github.com/clipperhouse/displaywidth v0.11.0 // indirect
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.4.0 // indirect
	github.com/mattn/go-runewidth v0.0.22 // indirect
	github.com/mikeschinkel/go-dt v0.7.0 // indirect
	github.com/mikeschinkel/go-dt/dtx v0.2.1 // indirect
	github.com/mikeschinkel/go-tealeaves/teacolor v0.0.0 // indirect
	github.com/mikeschinkel/go-tealeaves/teautils v0.2.0 // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/rivo/uniseg v0.4.7 // indirect
	github.com/xo/terminfo v0.0.0-20220910002029-abceb7e1c41e // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/term v0.40.0 // indirect
)
