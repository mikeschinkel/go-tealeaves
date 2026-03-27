package teaguide

// GuidePriority determines how a section is displayed in the guide.
type GuidePriority int

const (
	// PriorityRecommended highlights 1-2 items with prose explanations.
	PriorityRecommended GuidePriority = iota

	// PriorityAvailable shows other actions the user can take now.
	PriorityAvailable

	// PriorityBlocked shows a collapsed section with block reasons.
	PriorityBlocked
)

// GuideItem represents a single action or informational entry in a guide section.
type GuideItem struct {
	// ActionKey is the key string for dispatch (e.g., "enter", "v").
	// Empty means info-only (no dispatch).
	ActionKey string

	// KeyDisplay is the rendered label shown to the user (e.g., "[Enter]").
	KeyDisplay string

	// Label is the action name (e.g., "Select Module").
	Label string

	// Prose is the explanation text (e.g., "This module has uncommitted changes").
	Prose string

	// BlockReason explains why this item is blocked (only for PriorityBlocked items).
	BlockReason string
}

// GuideSection groups related guide items under a heading with a priority level.
type GuideSection struct {
	Priority GuidePriority
	Heading  string
	Items    []GuideItem
}

// GuideData is the complete data for one guide display. The host application
// builds this based on current state and passes it to [GuideModel.Open].
type GuideData struct {
	Title    string
	Sections []GuideSection // Ordered: Recommended, Available, Blocked
}
