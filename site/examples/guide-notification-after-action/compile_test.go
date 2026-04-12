// Source: site/src/content/docs/cookbook/notification-after-action.mdx:181,196,206
package main_test

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/mikeschinkel/go-tealeaves/teanotify"
)

// TestCompile_NotifyInitialization verifies that NewNotifyModel accepts
// NotifyOpts with all fields from the cookbook snippet and that Initialize()
// returns (NotifyModel, error).
// Source line 181.
func TestCompile_NotifyInitialization(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:           40,
		Duration:        3 * time.Second,
		Position:        teanotify.TopRightPosition,
		AllowEscToClose: true,
	})
	var err error
	notify, err = notify.Initialize()
	if err != nil {
		t.Fatalf("unexpected error from Initialize: %v", err)
	}
	_ = notify
}

// TestCompile_NotifyCmd verifies that NewNotifyCmd returns a tea.Cmd and that
// it can be passed to tea.Batch.
// Source line 196.
func TestCompile_NotifyCmd(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
	})
	notify, _ = notify.Initialize()

	var notifyCmd tea.Cmd // placeholder for the notify update cmd
	cmd := notify.NewNotifyCmd(teanotify.InfoKey, "Saved successfully")
	_ = tea.Batch(notifyCmd, cmd)
}

// TestCompile_NotifyRender verifies that Render(string) string exists on
// NotifyModel and accepts a content string.
// Source line 206.
func TestCompile_NotifyRender(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
	})
	notify, _ = notify.Initialize()

	fullContent := "some content"
	result := notify.Render(fullContent)
	_ = result
}
