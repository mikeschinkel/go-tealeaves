package teahilite

import "errors"

var (
	// ErrHighlight is the layer sentinel for syntax highlighting errors.
	ErrHighlight = errors.New("highlight")

	// ErrTokenize indicates a Chroma tokenization failure.
	ErrTokenize = errors.New("tokenize")

	// ErrFormat indicates a Chroma formatter failure.
	ErrFormat = errors.New("format")
)
