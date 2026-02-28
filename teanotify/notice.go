package teanotify

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
)

// NoticeKey is a typed key identifying a registered notice type.
type NoticeKey string

// Built-in notice keys for the default notice types.
const (
	InfoKey  NoticeKey = "Info"
	WarnKey  NoticeKey = "Warn"
	ErrorKey NoticeKey = "Error"
	DebugKey NoticeKey = "Debug"
)

// NerdFont prefix symbols. A NerdFont must be installed to render these.
// See https://www.nerdfonts.com/.
const (
	InfoNerdSymbol  = " "
	WarnNerdSymbol  = "󱈸 "
	ErrorNerdSymbol = "󰬅 "
	DebugNerdSymbol = "󰃤 "
)

// ASCII prefix strings for environments without special font support.
const (
	InfoASCIIPrefix    = "(i)"
	WarningASCIIPrefix = "(!)"
	ErrorASCIIPrefix   = "[!!]"
	DebugASCIIPrefix   = "(?)"
)

// Unicode prefix characters for environments that support Unicode.
const (
	InfoUnicodePrefix    = "\u24D8 "
	WarningUnicodePrefix = "\u26A0"
	ErrorUnicodePrefix   = "\u2718"
	DebugUnicodePrefix   = "\u003F"
)

// Default color hex codes for the built-in notice types.
const (
	InfoColor  = "#00FF00"
	WarnColor  = "#FFFF00"
	ErrorColor = "#FF0000"
	DebugColor = "#FF00FF"
	BackColor  = "#000000"
)

// DefaultLerpIncrement is the per-tick color interpolation step.
const DefaultLerpIncrement = 0.18

// Pre-parsed colors for the built-in defaults.
// The ignored errors are safe: these are hardcoded valid hex values.
var (
	infoColor, _  = colorful.Hex(InfoColor)
	warnColor, _  = colorful.Hex(WarnColor)
	errorColor, _ = colorful.Hex(ErrorColor)
	debugColor, _ = colorful.Hex(DebugColor)
	backColor, _  = colorful.Hex(BackColor)

	baseStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
)

// parsedColors caches parsed colorful.Color values by hex string.
var parsedColors = map[string]colorful.Color{
	InfoColor:  infoColor,
	WarnColor:  warnColor,
	ErrorColor: errorColor,
	DebugColor: debugColor,
	BackColor:  backColor,
}

// unicodePrefixes maps default notice keys to their Unicode prefix strings.
var unicodePrefixes = map[NoticeKey]string{
	InfoKey:  InfoUnicodePrefix,
	WarnKey:  WarningUnicodePrefix,
	ErrorKey: ErrorUnicodePrefix,
	DebugKey: DebugUnicodePrefix,
}

// NoticeDefinition contains all information needed to register a notice type.
type NoticeDefinition struct {
	// Key is the unique identifier for this notice type.
	Key NoticeKey

	// ForeColor is the hex color code (e.g. "#FF0000") for the notice.
	ForeColor string

	// Style is an optional lipgloss.Style used to render the notice.
	Style lipgloss.Style

	// Prefix is an optional string prepended to the notice message.
	Prefix string
}

// notifyMsg is the internal message type that activates a notice.
type notifyMsg struct {
	noticeKey NoticeKey
	msg       string
	dur       time.Duration
}

// notice represents an active notice instance with all data needed
// to render and expire itself.
type notice struct {
	message     string
	deathTime   time.Time
	prefix      string
	foreColor   colorful.Color
	style       lipgloss.Style
	width       int
	minWidth    int
	curLerpStep float64
	position    Position
}

// render produces the styled string representation of this notice,
// ready to be overlaid onto the main content.
func (n *notice) render() (result string) {
	newColor := backColor.BlendLab(n.foreColor, n.curLerpStep)
	lipColor := lipgloss.Color(newColor.Hex())

	actualWidth := n.width

	if n.minWidth > 0 {
		messageText := fmt.Sprintf("%v %v", n.prefix, n.message)
		messageWidth := lipgloss.Width(messageText)
		messageWidth += 3

		switch {
		case messageWidth < n.minWidth:
			actualWidth = n.minWidth
		case messageWidth > n.width:
			actualWidth = n.width
		default:
			actualWidth = messageWidth
		}
	}

	newStyle := baseStyle.
		Foreground(lipColor).
		BorderForeground(lipColor).
		Width(actualWidth).
		Padding(0, 1)

	textWidth := actualWidth - 2
	if textWidth < 1 {
		textWidth = 1
	}

	content := hangingWrap(n.prefix, n.message, textWidth)
	result = newStyle.Render(content)
	return result
}

// defaultNotices returns the four built-in notice definitions configured
// with the appropriate prefix style.
func defaultNotices(useNerdFont, useUnicodePrefix bool) (defs []NoticeDefinition) {
	var infoPref, warnPref, errPref, debugPref string

	switch {
	case useNerdFont:
		infoPref = InfoNerdSymbol
		warnPref = WarnNerdSymbol
		errPref = ErrorNerdSymbol
		debugPref = DebugNerdSymbol
	case useUnicodePrefix:
		infoPref = InfoUnicodePrefix
		warnPref = WarningUnicodePrefix
		errPref = ErrorUnicodePrefix
		debugPref = DebugUnicodePrefix
	default:
		infoPref = InfoASCIIPrefix
		warnPref = WarningASCIIPrefix
		errPref = ErrorASCIIPrefix
		debugPref = DebugASCIIPrefix
	}

	defs = []NoticeDefinition{
		{Key: InfoKey, Prefix: infoPref, ForeColor: InfoColor},
		{Key: WarnKey, Prefix: warnPref, ForeColor: WarnColor},
		{Key: ErrorKey, Prefix: errPref, ForeColor: ErrorColor},
		{Key: DebugKey, Prefix: debugPref, ForeColor: DebugColor},
	}
	return defs
}

// registerNotice validates a notice definition and adds it to the model's
// notice type registry. The model's noticeTypes map is modified in place;
// callers requiring immutable semantics must copy the map beforehand.
func registerNotice(m NotifyModel, def NoticeDefinition) (out NotifyModel, err error) {
	var fc colorful.Color

	if def.Key == "" {
		err = NewErr(ErrNotify, ErrInvalidNoticeKey,
			"key", string(def.Key),
		)
		goto end
	}

	fc, err = colorful.Hex(def.ForeColor)
	if err != nil {
		err = NewErr(ErrNotify, ErrInvalidColor,
			"color", def.ForeColor,
			"key", string(def.Key),
			err,
		)
		goto end
	}

	parsedColors[def.ForeColor] = fc
	m.noticeTypes[def.Key] = def
	out = m

end:
	return out, err
}
