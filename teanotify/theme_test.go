package teanotify

import (
	"testing"
	"time"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

func TestWithTheme_UpdatesDefaultNoticeColors(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: time.Second,
	})
	m, err := m.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	origInfo := m.noticeTypes[InfoKey].ForeColor

	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	newInfo := m.noticeTypes[InfoKey].ForeColor
	// The color should have been updated (palette StatusInfo maps to ANSI 86,
	// which is different from the default #00FF00)
	if newInfo == origInfo {
		t.Errorf("InfoKey ForeColor not updated by theme: still %q", origInfo)
	}
}

func TestWithTheme_PreservesCustomNotices(t *testing.T) {
	custom := NoticeDefinition{
		Key:       "Custom",
		ForeColor: "#AABBCC",
		Prefix:    "[C]",
	}
	m := NewNotifyModel(NotifyOpts{
		Width:         40,
		Duration:      time.Second,
		CustomNotices: []NoticeDefinition{custom},
	})
	m, err := m.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	// Custom notice should be unchanged
	if m.noticeTypes["Custom"].ForeColor != "#AABBCC" {
		t.Errorf("custom notice color changed: %s", m.noticeTypes["Custom"].ForeColor)
	}
}

func TestWithTheme_UninitializedModel(t *testing.T) {
	m := NewNotifyModel(NotifyOpts{
		Width:    40,
		Duration: time.Second,
	})
	// Not initialized — noticeTypes is nil
	theme := teautils.NewTheme(teautils.DarkSystemPalette(nil))
	m = m.WithTheme(theme)

	// Should not panic, noticeTypes stays nil
	if m.noticeTypes != nil {
		t.Error("expected nil noticeTypes on uninitialized model")
	}
}
