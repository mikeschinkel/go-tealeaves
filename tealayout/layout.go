package tealayout

import (
	"os"
	"strconv"
)

// Layout wraps a root Component and provides the top-level API for
// BubbleTea integration: SetSize from WindowSizeMsg, Render for View.
type Layout struct {
	root      *Component
	cfg       config
	width     int
	height    int
	sizeKnown bool
}

// NewLayout creates a Layout wrapping the given root component.
func NewLayout(root *Component, opts ...Option) *Layout {
	l := &Layout{root: root}
	for _, opt := range opts {
		opt(&l.cfg)
	}
	return l
}

// SetSize updates the layout dimensions (typically from tea.WindowSizeMsg).
func (l *Layout) SetSize(width, height int) {
	l.width = width
	l.height = height
	l.sizeKnown = true
	l.root.SetSize(width, height)
}

// Render resolves the layout and returns the composed output string.
// If auto-detect is enabled and SetSize has not been called, attempts
// to detect terminal size.
func (l *Layout) Render() (string, error) {
	if !l.sizeKnown {
		if l.cfg.autoDetectSize {
			w, h, err := l.detectSize()
			if err != nil {
				return "", NewErr(ErrZeroDimensions, "reason", "auto-detect failed", err)
			}
			l.SetSize(w, h)
		} else {
			return "", NewErr(ErrZeroDimensions, "reason", "SetSize not called and auto-detect disabled")
		}
	}
	return l.root.Render()
}

// MarkDirty forces the root component to re-resolve and re-render.
func (l *Layout) MarkDirty() {
	l.root.MarkDirty()
}

// Root returns the root component.
func (l *Layout) Root() *Component {
	return l.root
}

// Width returns the current layout width.
func (l *Layout) Width() int {
	return l.width
}

// Height returns the current layout height.
func (l *Layout) Height() int {
	return l.height
}

// detectSize uses the configured size source or falls back to env vars.
func (l *Layout) detectSize() (int, int, error) {
	if l.cfg.sizeSource != nil {
		return l.cfg.sizeSource()
	}
	return defaultSizeDetect()
}

// defaultSizeDetect tries COLUMNS/LINES environment variables.
func defaultSizeDetect() (int, int, error) {
	cols := os.Getenv("COLUMNS")
	lines := os.Getenv("LINES")
	if cols == "" || lines == "" {
		return 0, 0, NewErr(ErrZeroDimensions, "reason", "COLUMNS/LINES not set")
	}
	w, err := strconv.Atoi(cols)
	if err != nil {
		return 0, 0, NewErr(ErrZeroDimensions, "key", "COLUMNS", "value", cols, err)
	}
	h, err := strconv.Atoi(lines)
	if err != nil {
		return 0, 0, NewErr(ErrZeroDimensions, "key", "LINES", "value", lines, err)
	}
	if w <= 0 || h <= 0 {
		return 0, 0, NewErr(ErrZeroDimensions, "width", w, "height", h)
	}
	return w, h, nil
}
