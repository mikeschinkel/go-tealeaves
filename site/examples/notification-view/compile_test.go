// Source: site/src/content/docs/components/notification-view.mdx:21,67,161,176
package examples_test

import (
	"testing"
	"time"

	"github.com/mikeschinkel/go-tealeaves/teanotify"
)

// TestCompile_NotifyQuickExample verifies the quick example from notification-view.mdx.
func TestCompile_NotifyQuickExample(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:           40,
		Duration:        3 * time.Second,
		AllowEscToClose: true,
	})
	var err error
	notify, err = notify.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	cmd := notify.NewNotifyCmd(teanotify.InfoKey, "File saved successfully")
	_ = cmd
}

// TestCompile_PositionConfig verifies position configuration from notification-view.mdx.
func TestCompile_PositionConfig(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Position: teanotify.BottomRightPosition,
	})
	notify, _ = notify.Initialize()

	notify = notify.WithPosition(teanotify.TopCenterPosition)
	_ = notify
}

// TestCompile_CustomNoticeType verifies RegisterNoticeType from notification-view.mdx.
func TestCompile_CustomNoticeType(t *testing.T) {
	notify := teanotify.NewNotifyModel(teanotify.NotifyOpts{
		Width:    40,
		Duration: 3 * time.Second,
	})
	var err error
	notify, err = notify.Initialize()
	if err != nil {
		t.Fatal(err)
	}

	notify, err = notify.RegisterNoticeType(teanotify.NoticeDefinition{
		Key:       "success",
		Prefix:    "[✓]",
		ForeColor: "#00ff00",
	})
	if err != nil {
		t.Fatal(err)
	}

	cmd := notify.NewNotifyCmd("success", "Build complete")
	_ = cmd
}
