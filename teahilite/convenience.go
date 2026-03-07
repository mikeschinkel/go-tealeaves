package teahilite

var defaultHighlighter = NewHighlighter(HighlighterArgs{})

// Highlight returns syntax-highlighted code as an ANSI-styled string
// using the default Highlighter (monokai style, terminal256 formatter).
func Highlight(code, language string) (string, error) {
	return defaultHighlighter.Highlight(code, language)
}

// HighlightLines returns highlighted code split into individual lines
// using the default Highlighter (monokai style, terminal256 formatter).
func HighlightLines(code, language string) ([]string, error) {
	return defaultHighlighter.HighlightLines(code, language)
}
