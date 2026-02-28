package teanotify

// Position identifies where a notification overlay appears on screen.
type Position string

const (
	TopLeftPosition      Position = "TL"
	TopCenterPosition    Position = "TC"
	TopRightPosition     Position = "TR"
	BottomLeftPosition   Position = "BL"
	BottomCenterPosition Position = "BC"
	BottomRightPosition  Position = "BR"
	UnspecifiedPosition  Position = ""
)

// IsValid reports whether p is a recognized position constant.
func (p Position) IsValid() (valid bool) {
	valid = p.String() != "unknown"
	return valid
}

// String returns the kebab-case name of the position.
func (p Position) String() (s string) {
	switch p {
	case TopLeftPosition:
		s = "top-left"
	case TopCenterPosition:
		s = "top-center"
	case TopRightPosition:
		s = "top-right"
	case BottomLeftPosition:
		s = "bottom-left"
	case BottomCenterPosition:
		s = "bottom-center"
	case BottomRightPosition:
		s = "bottom-right"
	default:
		s = "unknown"
	}
	return s
}

// Label returns the human-readable label for the position.
func (p Position) Label() (label string) {
	switch p {
	case TopLeftPosition:
		label = "Top Left"
	case TopCenterPosition:
		label = "Top Center"
	case TopRightPosition:
		label = "Top Right"
	case BottomLeftPosition:
		label = "Bottom Left"
	case BottomCenterPosition:
		label = "Bottom Center"
	case BottomRightPosition:
		label = "Bottom Right"
	default:
		label = "Unknown"
	}
	return label
}
