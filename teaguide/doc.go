// Package teaguide provides a Bubble Tea v2 component for context-aware
// workflow guidance overlays.
//
// The guide is a modal overlay that shows prioritized, actionable next steps
// with prose explanations. It features direct key dispatch: pressing an
// action key inside the guide closes it AND triggers the action via
// [ActionSelectedMsg].
//
// Guide content is organized into three priority sections:
//
//   - Recommended: 1-2 highlighted items with prose explanations
//   - Available: other actions the user can take now
//   - Blocked: collapsed section showing actions not yet available, with reasons
//
// The host application provides [GuideData] describing current state;
// teaguide handles rendering, scrolling, and key dispatch.
//
// Basic usage:
//
//	guide := teaguide.NewGuideModel()
//	guide = guide.SetSize(width, height)
//
//	// Open with context-specific data
//	guide, cmd = guide.Open(teaguide.GuideData{
//	    Title: "What's Next?",
//	    Sections: []teaguide.GuideSection{
//	        {Priority: teaguide.PriorityRecommended, Heading: "Recommended", Items: items},
//	    },
//	})
//
//	// In View(), overlay on background
//	view = guide.OverlayModal(backgroundView)
//
// # Stability
//
// This package is provisional as of v0.3.0. The public API may change in
// minor releases until promoted to stable.
package teaguide
