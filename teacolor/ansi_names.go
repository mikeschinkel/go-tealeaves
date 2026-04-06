package teacolor

import "image/color"

// Standard ANSI color names (colors 0-15).
// These map to the 16 standard terminal colors.
var (
	Black         color.Color = Color("0")
	Red           color.Color = Color("1")
	Green         color.Color = Color("2")
	Yellow        color.Color = Color("3")
	Blue          color.Color = Color("4")
	Magenta       color.Color = Color("5")
	Cyan          color.Color = Color("6")
	White         color.Color = Color("7")
	BrightBlack   color.Color = Color("8")
	BrightRed     color.Color = Color("9")
	BrightGreen   color.Color = Color("10")
	BrightYellow  color.Color = Color("11")
	BrightBlue    color.Color = Color("12")
	BrightMagenta color.Color = Color("13")
	BrightCyan    color.Color = Color("14")
	BrightWhite   color.Color = Color("15")
)
