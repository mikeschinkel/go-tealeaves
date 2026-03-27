package teaguide

import (
	"unicode"
	"unicode/utf8"

	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// GuideItemOpts holds optional fields for NewGuideItem.
type GuideItemOpts struct {
	Label       string // Override Meta.StatusBarLabel
	Prose       string // Explanation text for Recommended section
	BlockReason string // Block reason for Blocked section
}

// NewGuideItem creates a GuideItem from KeyMeta registry data.
// Derives ActionKey and KeyDisplay from the binding's first key.
func NewGuideItem(meta teautils.KeyMeta, opts *GuideItemOpts) GuideItem {
	var o GuideItemOpts
	if opts != nil {
		o = *opts
	}
	keys := meta.Binding.Keys()
	actionKey := ""
	keyDisplay := ""
	if len(keys) > 0 {
		actionKey = keys[0]
		keyDisplay = "[" + capitalizeFirst(keys[0]) + "]"
	}
	label := meta.StatusBarLabel
	if o.Label != "" {
		label = o.Label
	}
	return GuideItem{
		ActionKey:   actionKey,
		KeyDisplay:  keyDisplay,
		Label:       label,
		Prose:       o.Prose,
		BlockReason: o.BlockReason,
	}
}

// capitalizeFirst uppercases the first rune of s ("tab" → "Tab").
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}
