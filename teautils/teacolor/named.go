package teacolor

import "image/color"

// Curated semantic color aliases — descriptive names for commonly used colors.
// These use either ANSI 256 indices or hex values as appropriate.
var (
	// Grays (ANSI 256 grayscale ramp)
	DarkGray    color.Color = Color("240")
	Gray        color.Color = Color("245")
	LightGray   color.Color = Color("252")
	DimGray     color.Color = Color("238")
	SilverGray  color.Color = Color("250")
	CharcoalGray color.Color = Color("236")

	// Reds
	Coral    color.Color = Color("#FF7F50")
	Crimson  color.Color = Color("#DC143C")
	Salmon   color.Color = Color("#FA8072")
	Tomato   color.Color = Color("#FF6347")
	FireBrick color.Color = Color("#B22222")
	DarkRed  color.Color = Color("#8B0000")

	// Oranges
	Orange     color.Color = Color("#FFA500")
	DarkOrange color.Color = Color("#FF8C00")
	Tangerine  color.Color = Color("#FF9966")

	// Yellows
	Gold       color.Color = Color("#FFD700")
	Amber      color.Color = Color("#FFBF00")
	Khaki      color.Color = Color("#F0E68C")

	// Greens
	Lime       color.Color = Color("#00FF00")
	Emerald    color.Color = Color("#50C878")
	ForestGreen color.Color = Color("#228B22")
	Olive      color.Color = Color("#808000")
	Mint       color.Color = Color("#98FF98")
	SeaGreen   color.Color = Color("#2E8B57")

	// Blues
	SkyBlue    color.Color = Color("#87CEEB")
	DodgerBlue color.Color = Color("#1E90FF")
	RoyalBlue  color.Color = Color("#4169E1")
	Navy       color.Color = Color("#000080")
	SteelBlue  color.Color = Color("#4682B4")
	CornflowerBlue color.Color = Color("#6495ED")

	// Purples
	Plum       color.Color = Color("#DDA0DD")
	Indigo     color.Color = Color("#4B0082")
	Violet     color.Color = Color("#EE82EE")
	Orchid     color.Color = Color("#DA70D6")
	Purple     color.Color = Color("#800080")
	Lavender   color.Color = Color("#E6E6FA")

	// Cyans / Teals
	Teal       color.Color = Color("#008080")
	Turquoise  color.Color = Color("#40E0D0")
	Aquamarine color.Color = Color("#7FFFD4")

	// Pinks
	HotPink    color.Color = Color("#FF69B4")
	Rose       color.Color = Color("#FF007F")
	Blush      color.Color = Color("#DE5D83")

	// Neutrals
	SlateGray  color.Color = Color("#708090")
	Ivory      color.Color = Color("#FFFFF0")
	Beige      color.Color = Color("#F5F5DC")
)
