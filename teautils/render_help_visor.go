package teautils

// EXTRACTION-BOUNDARY: teahelp

import (
	"sort"
	"strings"
	"unicode"

	"charm.land/lipgloss/v2"
)

// HelpVisorStyle holds styling configuration for help visor rendering
type HelpVisorStyle struct {
	TitleStyle    lipgloss.Style // Style for title (flush left, no centering)
	CategoryStyle lipgloss.Style // Style for category headers (no margin — model handles spacing)
	KeyStyle      lipgloss.Style // Style for key labels (Width is set dynamically; do not set Width here)
	DescStyle     lipgloss.Style // Style for descriptions

	KeyColumnGap  int      // Extra spaces between longest key and description (default 4)
	CategoryOrder []string // Preferred order for categories (unspecified go last)
}

// ThemedHelpVisorStyle returns a HelpVisorStyle derived from the given theme.
func ThemedHelpVisorStyle(theme Theme) HelpVisorStyle {
	return HelpVisorStyle{
		TitleStyle:    theme.HelpVisor.TitleStyle,
		CategoryStyle: theme.HelpVisor.CategoryStyle,
		KeyStyle:      theme.HelpVisor.KeyStyle,
		DescStyle:     theme.HelpVisor.DescStyle,
		KeyColumnGap:  4,
		CategoryOrder: []string{"Navigation", "Actions", "System", "Other"},
	}
}

// DefaultHelpVisorStyle returns the default styling for help visor
func DefaultHelpVisorStyle() HelpVisorStyle {
	return HelpVisorStyle{
		TitleStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("99")), // Purple — flush left
		CategoryStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("178")), // Gold — distinct from cyan keys
		KeyStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")). // Cyan
			Bold(true),
		DescStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")), // Light gray
		KeyColumnGap:  4,
		CategoryOrder: []string{"Navigation", "Actions", "System", "Other"},
	}
}

// FormatKeyDisplay formats key bindings for display.
// Uses DisplayKeys if provided, otherwise formats binding keys with special handling.
// Multiple keys are joined with "/" separator.
// Duplicate display keys (e.g., "a" and "A" both mapping to "A") are collapsed.
func FormatKeyDisplay(k KeyMeta) string {
	var displayKeys []string

	if len(k.DisplayKeys) > 0 {
		// Use custom display names
		displayKeys = k.DisplayKeys
	} else {
		// Get keys from binding and format them
		bindingKeys := k.Binding.Keys()
		displayKeys = make([]string, len(bindingKeys))
		for i, bk := range bindingKeys {
			displayKeys[i] = formatSingleKey(bk)
		}
	}

	// Deduplicate (e.g., "a" and "A" both title-case to "A")
	seen := make(map[string]bool)
	var unique []string
	for _, dk := range displayKeys {
		if !seen[dk] {
			seen[dk] = true
			unique = append(unique, dk)
		}
	}
	displayKeys = unique

	if len(displayKeys) == 0 {
		return ""
	}

	// Join with "/" separator: [Up]/[K] format
	var parts []string
	for _, dk := range displayKeys {
		parts = append(parts, "["+dk+"]")
	}
	return strings.Join(parts, "/")
}

// formatSingleKey converts a raw key string to title-cased display format.
// Handles special keys like space, arrows, etc.
func formatSingleKey(key string) string {
	switch key {
	case " ":
		return "Space"
	case "backspace":
		return "Bksp"
	case "delete":
		return "Del"
	case "pgup":
		return "PgUp"
	case "pgdown":
		return "PgDn"
	}
	return ProperCaseShortcut(key)
}

// ProperCaseShortcut title-cases each part of a shortcut key string.
// Handles compound keys separated by "+": "ctrl+c" -> "Ctrl+C", "alt+v" -> "Alt+V"
// Single keys: "esc" -> "Esc", "a" -> "A"
func ProperCaseShortcut(s string) string {
	parts := strings.Split(s, "+")
	for i, p := range parts {
		if p == "" {
			continue
		}
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, "+")
}

// GetSortedCategories returns categories sorted by preferred order.
// Categories in preferredOrder come first (in that order), then alphabetical.
func GetSortedCategories(keysByCategory map[string][]KeyMeta, preferredOrder []string) []string {
	// Create order map
	orderMap := make(map[string]int)
	for i, cat := range preferredOrder {
		orderMap[cat] = i
	}

	// Collect all categories
	categories := make([]string, 0, len(keysByCategory))
	for cat := range keysByCategory {
		categories = append(categories, cat)
	}

	// Sort by preferred order, then alphabetically
	sort.Slice(categories, func(i, j int) bool {
		cat1, cat2 := categories[i], categories[j]
		order1, has1 := orderMap[cat1]
		order2, has2 := orderMap[cat2]

		// Both in preferred order - use that order
		if has1 && has2 {
			return order1 < order2
		}

		// Only one in preferred order - it goes first
		if has1 {
			return true
		}
		if has2 {
			return false
		}

		// Neither in preferred order - alphabetical
		return cat1 < cat2
	})

	return categories
}
