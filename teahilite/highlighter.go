package teahilite

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

const (
	DefaultStyleName     = "monokai"
	DefaultFormatterName = "terminal256"
)

// HighlighterArgs configures a Highlighter instance.
type HighlighterArgs struct {
	StyleName     string // Chroma style name (default "monokai")
	FormatterName string // Chroma formatter name (default "terminal256")
}

// Highlighter renders syntax-highlighted code as ANSI-styled terminal output.
type Highlighter struct {
	style     *chroma.Style
	formatter chroma.Formatter
}

// NewHighlighter creates a Highlighter with the given configuration.
// Unknown style or formatter names fall back to Chroma defaults.
func NewHighlighter(args HighlighterArgs) *Highlighter {
	styleName := args.StyleName
	if styleName == "" {
		styleName = DefaultStyleName
	}
	formatterName := args.FormatterName
	if formatterName == "" {
		formatterName = DefaultFormatterName
	}

	style := styles.Get(styleName)
	if style == nil {
		style = styles.Fallback
	}

	formatter := formatters.Get(formatterName)
	if formatter == nil {
		formatter = formatters.Fallback
	}

	return &Highlighter{
		style:     style,
		formatter: formatter,
	}
}

// Highlight returns syntax-highlighted code as an ANSI-styled string.
func (h *Highlighter) Highlight(code, language string) (result string, err error) {
	var iterator chroma.Iterator
	var buf bytes.Buffer

	lexer := lexers.Get(language)
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	iterator, err = lexer.Tokenise(nil, code)
	if err != nil {
		err = NewErr(ErrHighlight, ErrTokenize, err)
		goto end
	}

	err = h.formatter.Format(&buf, h.style, iterator)
	if err != nil {
		err = NewErr(ErrHighlight, ErrFormat, err)
		goto end
	}

	result = buf.String()

end:
	return result, err
}

// HighlightLines returns highlighted code split into individual lines.
// Each line is a separate ANSI-styled string.
func (h *Highlighter) HighlightLines(code, language string) (lines []string, err error) {
	var highlighted string

	highlighted, err = h.Highlight(code, language)
	if err != nil {
		goto end
	}

	lines = strings.Split(highlighted, "\n")

end:
	return lines, err
}
