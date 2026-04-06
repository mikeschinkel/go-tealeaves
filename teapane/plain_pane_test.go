package teapane_test

import (
	"strings"
	"testing"

	lipgloss "charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/tealayout"
	"github.com/mikeschinkel/go-tealeaves/teapane"
)

func TestPlainPane_InterfaceDetection(t *testing.T) {
	pp := teapane.NewPlainPane(func(w, h int, focused bool) string {
		return "test"
	})

	var _ tealayout.ContentProvider = pp
	var _ tealayout.SetSizer = pp
}

func TestPlainPane_View(t *testing.T) {
	pp := teapane.NewPlainPane(func(w, h int, focused bool) string {
		return "hello world"
	})

	pp.SetSize(40, 1)
	view := pp.Content()
	if !strings.Contains(view, "hello world") {
		t.Errorf("expected view to contain 'hello world', got: %s", view)
	}
}

func TestPlainPane_WithStyle(t *testing.T) {
	pp := teapane.NewPlainPane(func(w, h int, focused bool) string {
		return "styled"
	}).WithStyle(lipgloss.NewStyle().Bold(true))

	pp.SetSize(20, 1)
	view := pp.Content()
	if !strings.Contains(view, "styled") {
		t.Errorf("expected view to contain 'styled', got: %s", view)
	}
}
