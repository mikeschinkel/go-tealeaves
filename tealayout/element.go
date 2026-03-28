package tealayout

import (
	"fmt"

	lipgloss "charm.land/lipgloss/v2"
)

// element is the unexported interface that both *Element[T] and *Pane satisfy.
// It allows Pane to hold heterogeneous children in []element.
type element interface {
	view() string
	setSize(width, height int)
	style() (lipgloss.Style, bool)
	sizeHint(availW, availH int) (SizeHint, bool)
}

// Element wraps any value with auto-detected capabilities for use in layouts.
// Capabilities (View, SetSize, Style, SizeHint) are detected once at
// construction time via type switches.
type Element[T any] struct {
	value   T
	viewFn  func(T) string
	sizeFn  func(T, int, int)
	styleFn func(T) lipgloss.Style
	hintFn  func(T, int, int) SizeHint
}

// NewElement creates an Element wrapping v and auto-detects capabilities.
// Detected interfaces: Viewer, SetSizer, Styler, SizeHinter, fmt.Stringer.
// For plain string values, view() returns the string itself.
func NewElement[T any](v T) *Element[T] {
	e := &Element[T]{value: v}

	// Detect view capability
	switch w := any(v).(type) {
	case Viewer:
		e.viewFn = func(_ T) string { return w.View() }
	case fmt.Stringer:
		e.viewFn = func(_ T) string { return w.String() }
	default:
		if s, ok := any(v).(string); ok {
			e.viewFn = func(_ T) string { return s }
		}
	}

	// Detect SetSize capability
	if ss, ok := any(v).(SetSizer); ok {
		e.sizeFn = func(_ T, w, h int) { ss.SetSize(w, h) }
	}

	// Detect Style capability
	if st, ok := any(v).(Styler); ok {
		e.styleFn = func(_ T) lipgloss.Style { return st.Style() }
	}

	// Detect SizeHint capability
	if sh, ok := any(v).(SizeHinter); ok {
		e.hintFn = func(_ T, w, h int) SizeHint { return sh.SizeHint(w, h) }
	}

	return e
}

// Widget returns the underlying value for type-safe access.
func (e *Element[T]) Widget() T {
	return e.value
}

// view satisfies the element interface.
func (e *Element[T]) view() string {
	if e.viewFn != nil {
		return e.viewFn(e.value)
	}
	return ""
}

// setSize satisfies the element interface.
func (e *Element[T]) setSize(width, height int) {
	if e.sizeFn != nil {
		e.sizeFn(e.value, width, height)
	}
}

// style satisfies the element interface.
func (e *Element[T]) style() (lipgloss.Style, bool) {
	if e.styleFn != nil {
		return e.styleFn(e.value), true
	}
	return lipgloss.Style{}, false
}

// sizeHint satisfies the element interface.
func (e *Element[T]) sizeHint(availW, availH int) (SizeHint, bool) {
	if e.hintFn != nil {
		return e.hintFn(e.value, availW, availH), true
	}
	return SizeHint{}, false
}

// StringElement is a shorthand for NewElement wrapping a plain string.
func StringElement(s string) *Element[string] {
	return &Element[string]{
		value:  s,
		viewFn: func(v string) string { return v },
	}
}

// Compile-time interface checks.
var _ element = (*Element[any])(nil)
