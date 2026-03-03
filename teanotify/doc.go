// Package teanotify provides a Bubble Tea model for rendering overlay
// notification messages on top of terminal UI content.
//
// It is a successor to go.dalton.dog/bubbleup with a redesigned API
// and renamed terminology (alert → notify/notice).
//
// Usage:
//
//	model := teanotify.NewNotifyModel(teanotify.NotifyOpts{
//	    Width:    60,
//	    Duration: 3 * time.Second,
//	})
//	model, err := model.Initialize()
//
//	// In Update():
//	cmd := model.NewNotifyCmd(teanotify.InfoKey, "File saved")
//
//	// In View():
//	return model.Render(content)
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teanotify
