package tealayout

import (
	"os"
	"strconv"
)

// Layout wraps a root Pane and provides the top-level API for
// BubbleTea integration: SetSize from WindowSizeMsg, Render for View.
type Layout struct {
	root      *Pane
	cfg       config
	width     int
	height    int
	sizeKnown bool
	panes     map[string]*Pane // name → pane registry (lazy)
}

// NewLayout creates a Layout wrapping the given root pane.
func NewLayout(root *Pane, opts ...Option) *Layout {
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

// MarkDirty forces the root pane to re-resolve and re-render.
func (l *Layout) MarkDirty() {
	l.root.MarkDirty()
}

// Root returns the root pane.
func (l *Layout) Root() *Pane {
	return l.root
}

// Pane returns the named pane, or nil if not found.
func (l *Layout) Pane(name string) *Pane {
	l.ensureRegistry()
	return l.panes[name]
}

// ensureRegistry lazily builds the name→pane registry by walking the tree.
func (l *Layout) ensureRegistry() {
	if l.panes != nil {
		return
	}
	l.panes = make(map[string]*Pane)
	l.walkPanes(l.root)
}

// walkPanes recursively registers named panes.
func (l *Layout) walkPanes(p *Pane) {
	if p.name != "" {
		l.panes[p.name] = p
	}
	for _, ch := range p.children {
		if cp, ok := ch.elem.(*Pane); ok {
			l.walkPanes(cp)
		}
	}
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
