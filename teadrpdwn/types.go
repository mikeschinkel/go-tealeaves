package teadrpdwn

// Option represents a dropdown item with display text and optional value
type Option struct {
	Text  string
	Value interface{} // Actual value (e.g., DB primary key, ID, or any type)
}

// ToOptions converts a string slice to Options (Text and Value are the same)
func ToOptions(strings []string) []Option {
	items := make([]Option, len(strings))
	for i, s := range strings {
		items[i] = Option{Text: s, Value: s}
	}
	return items
}

// OptionSelectedMsg is sent when user confirms selection with Enter
type OptionSelectedMsg struct {
	Index int
	Text  string      // Display text
	Value interface{} // Underlying value
}

// DropdownCancelledMsg is sent when user cancels with Esc
type DropdownCancelledMsg struct{}
