package tealayout

// Option configures a Layout.
type Option func(*config)

type config struct {
	autoDetectSize bool
	sizeSource     func() (width, height int, err error)
}

// WithAutoDetectSize enables automatic terminal size detection on first
// Render() if SetSize() has not been called.
func WithAutoDetectSize(enabled bool) Option {
	return func(c *config) {
		c.autoDetectSize = enabled
	}
}

// WithSizeSource provides a custom function to detect terminal size.
// Used with WithAutoDetectSize.
func WithSizeSource(fn func() (width, height int, err error)) Option {
	return func(c *config) {
		c.sizeSource = fn
	}
}
