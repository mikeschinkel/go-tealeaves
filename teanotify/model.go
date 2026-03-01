package teanotify

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/ansi"
	"github.com/lucasb-eyer/go-colorful"
)

// NotifyOpts holds all constructor-time configuration for a NotifyModel.
type NotifyOpts struct {
	// Width is the fixed (or maximum, when MinWidth > 0) width of the
	// notification overlay in terminal cells.
	Width int

	// MinWidth enables dynamic width mode when > 0. The notification
	// width will vary between MinWidth and Width based on message length.
	MinWidth int

	// Duration is how long each notification remains visible.
	Duration time.Duration

	// UseNerdFont selects NerdFont prefix symbols for default notices.
	UseNerdFont bool

	// UseUnicodePrefix selects Unicode prefix characters for default notices.
	UseUnicodePrefix bool

	// AllowEscToClose lets the user dismiss active notices with the Esc key.
	AllowEscToClose bool

	// Position determines where the notification overlay appears.
	// Defaults to TopLeftPosition if unspecified.
	Position Position

	// NoDefaultNotices prevents registration of the four built-in notice
	// types (Info, Warn, Error, Debug).
	NoDefaultNotices bool

	// CustomNotices are additional notice definitions registered during
	// construction, after defaults (if enabled).
	CustomNotices []NoticeDefinition
}

// NotifyModel maintains registered notice types and facilitates the display,
// animation, and clearing of overlay notifications.
type NotifyModel struct {
	useNerdFont      bool
	useUnicodePrefix bool
	allowEscToClose  bool
	noDefaultNotices bool
	customNotices    []NoticeDefinition
	noticeTypes      map[NoticeKey]NoticeDefinition
	activeNotice     *notice
	width            int
	minWidth         int
	duration         time.Duration
	position         Position
}

// NewNotifyModel creates a NotifyModel from the provided options. The model
// is not yet ready to use — call Initialize() to validate options, register
// default notices, and register any custom notices.
func NewNotifyModel(opts NotifyOpts) (m NotifyModel) {
	m = NotifyModel{
		width:            opts.Width,
		minWidth:         opts.MinWidth,
		duration:         opts.Duration,
		useNerdFont:      opts.UseNerdFont,
		useUnicodePrefix: opts.UseUnicodePrefix,
		allowEscToClose:  opts.AllowEscToClose,
		position:         opts.Position,
		noDefaultNotices: opts.NoDefaultNotices,
		customNotices:    opts.CustomNotices,
	}
	return m
}

// Initialize validates options, registers default notices (unless
// NoDefaultNotices), and registers any custom notices. Returns an error
// if validation or registration fails.
func (m NotifyModel) Initialize() (out NotifyModel, err error) {
	out = m

	if out.width <= 0 {
		err = NewErr(ErrNotify, ErrInvalidWidth,
			"width", out.width,
		)
		goto end
	}

	if out.duration <= 0 {
		err = NewErr(ErrNotify, ErrInvalidDuration,
			"duration", out.duration,
		)
		goto end
	}

	out.noticeTypes = make(map[NoticeKey]NoticeDefinition)

	if out.position == UnspecifiedPosition {
		out.position = TopLeftPosition
	}

	if out.minWidth > out.width {
		out.minWidth = out.width
	}

	if !out.noDefaultNotices {
		for _, def := range defaultNotices(out.useNerdFont, out.useUnicodePrefix) {
			out, err = registerNotice(out, def)
			if err != nil {
				goto end
			}
		}
	}

	for _, def := range out.customNotices {
		out, err = registerNotice(out, def)
		if err != nil {
			goto end
		}
	}

end:
	return out, err
}

// WithPosition returns a copy of the model with the specified position.
func (m NotifyModel) WithPosition(pos Position) (out NotifyModel) {
	m.position = pos
	out = m
	return out
}

// WithMinWidth returns a copy of the model with dynamic width enabled.
// If min exceeds the model's width, it is clamped.
func (m NotifyModel) WithMinWidth(min int) (out NotifyModel) {
	if min > m.width {
		min = m.width
	}
	m.minWidth = min
	out = m
	return out
}

// WithUnicodePrefix returns a copy of the model with Unicode prefixes
// applied to all registered default notice types.
func (m NotifyModel) WithUnicodePrefix() (out NotifyModel) {
	m.useNerdFont = false
	m.useUnicodePrefix = true
	newTypes := make(map[NoticeKey]NoticeDefinition, len(m.noticeTypes))
	for name, nt := range m.noticeTypes {
		prefix, ok := unicodePrefixes[nt.Key]
		if ok {
			nt.Prefix = prefix
		}
		newTypes[name] = nt
	}
	m.noticeTypes = newTypes
	out = m
	return out
}

// WithAllowEscToClose returns a copy of the model with Esc-to-dismiss
// enabled.
func (m NotifyModel) WithAllowEscToClose() (out NotifyModel) {
	m.allowEscToClose = true
	out = m
	return out
}

// Init satisfies the Bubble Tea model pattern. Returns nil.
func (m NotifyModel) Init() (cmd tea.Cmd) {
	return cmd
}

// Update processes a message and returns the updated model and command.
// Returns (NotifyModel, tea.Cmd) — a concrete type, not tea.Model.
func (m NotifyModel) Update(msg tea.Msg) (out NotifyModel, cmd tea.Cmd) {
	out = m
	switch msg := msg.(type) {
	case notifyMsg:
		out.activeNotice = out.newNotice(msg.noticeKey, msg.msg, msg.dur)
		cmd = tickCmd()
		goto end
	case tickMsg:
		if out.activeNotice == nil {
			goto end
		}
		if out.activeNotice.deathTime.Before(time.Time(msg)) {
			out.activeNotice = nil
			goto end
		}
		out.activeNotice.curLerpStep += DefaultLerpIncrement
		if out.activeNotice.curLerpStep > 1 {
			out.activeNotice.curLerpStep = 1
		}
		cmd = tickCmd()
		goto end
	case tea.KeyMsg:
		if out.activeNotice == nil {
			goto end
		}
		if msg.String() != "esc" {
			goto end
		}
		if !out.allowEscToClose {
			goto end
		}
		out.activeNotice = nil
		goto end
	default:
		if out.activeNotice != nil {
			cmd = tickCmd()
		}
	}

end:
	return out, cmd
}

// HasActiveNotice reports whether a notice is currently being displayed.
func (m NotifyModel) HasActiveNotice() (active bool) {
	active = m.activeNotice != nil
	return active
}

// View returns an empty string. NotifyModel is not meant to be rendered
// directly; use Render to overlay notifications onto your view content.
func (m NotifyModel) View() (s string) {
	return s
}

// Render overlays the active notification (if any) onto the provided
// content string. Call this as the final step of your parent model's View.
func (m NotifyModel) Render(content string) (result string) {
	var notifString string
	var notifSplit []string
	var contentSplit []string
	var notifWidth, contentWidth int
	var notifHeight, contentHeight int
	var builder strings.Builder

	result = content
	if m.activeNotice == nil {
		goto end
	}

	notifString = m.activeNotice.render()
	notifSplit, notifWidth = getLines(notifString)
	contentSplit, contentWidth = getLines(content)
	notifHeight = len(notifSplit)
	contentHeight = len(contentSplit)

	for i := range contentHeight {
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(m.buildLineForPosition(
			contentSplit[i],
			notifSplit,
			i,
			notifHeight,
			contentHeight,
			notifWidth,
			contentWidth,
		))
	}
	result = builder.String()

end:
	return result
}

// NewNotifyCmd constructs a tea.Cmd that triggers a notification of the
// given type with the provided message.
func (m NotifyModel) NewNotifyCmd(noticeType NoticeKey, message string) (cmd tea.Cmd) {
	cmd = func() tea.Msg {
		return notifyMsg{noticeKey: noticeType, msg: message, dur: m.duration}
	}
	return cmd
}

// RegisterNoticeType registers a custom notice definition with immutable
// value semantics. The returned model contains the new registration; the
// receiver is unchanged.
func (m NotifyModel) RegisterNoticeType(def NoticeDefinition) (out NotifyModel, err error) {
	newTypes := make(map[NoticeKey]NoticeDefinition, len(m.noticeTypes)+1)
	for k, v := range m.noticeTypes {
		newTypes[k] = v
	}
	m.noticeTypes = newTypes
	out, err = registerNotice(m, def)
	return out, err
}

// newNotice creates a notice instance from a registered notice type.
// Returns nil if key or msg is empty, or if key is not registered.
func (m NotifyModel) newNotice(key NoticeKey, msg string, dur time.Duration) (n *notice) {
	var noticeDef NoticeDefinition
	var foreColor colorful.Color
	var ok bool

	if msg == "" || key == "" {
		goto end
	}

	noticeDef, ok = m.noticeTypes[key]
	if !ok {
		goto end
	}

	foreColor, ok = parsedColors[noticeDef.ForeColor]
	if !ok {
		// Color was validated during registration; parse error is unreachable.
		foreColor, _ = colorful.Hex(noticeDef.ForeColor)
	}

	n = &notice{
		message:     msg,
		deathTime:   time.Now().Add(dur),
		prefix:      noticeDef.Prefix,
		foreColor:   foreColor,
		style:       noticeDef.Style,
		width:       m.width,
		minWidth:    m.minWidth,
		curLerpStep: 0.3,
		position:    m.position,
	}

end:
	return n
}

// buildLineForPosition overlays notification content onto a single content
// line based on the active notice's position.
func (m NotifyModel) buildLineForPosition(
	contentLine string,
	notifLines []string,
	lineIdx, notifHeight, contentHeight, notifWidth, contentWidth int,
) (result string) {
	var notifIdx int
	var showNotif bool
	var notifLine string
	var keepWidth int
	var contentLineWidth int
	var padding string

	result = contentLine

	switch m.activeNotice.position {
	case TopLeftPosition, TopCenterPosition, TopRightPosition:
		showNotif = lineIdx < notifHeight
		notifIdx = lineIdx
	case BottomLeftPosition, BottomCenterPosition, BottomRightPosition:
		startLine := contentHeight - notifHeight
		if startLine < 0 {
			startLine = 0
		}
		if lineIdx >= startLine {
			showNotif = true
			notifIdx = lineIdx - startLine
		}
	}

	if !showNotif {
		goto end
	}

	notifLine = notifLines[notifIdx]

	switch m.activeNotice.position {
	case TopLeftPosition, BottomLeftPosition:
		result = notifLine + cutLeft(contentLine, notifWidth)
	case TopRightPosition, BottomRightPosition:
		keepWidth = contentWidth - notifWidth
		if keepWidth < 0 {
			result = notifLine
			goto end
		}
		contentLineWidth = ansi.StringWidth(contentLine)
		if contentLineWidth < keepWidth {
			padding = strings.Repeat(" ", keepWidth-contentLineWidth)
			result = contentLine + padding + notifLine
			goto end
		}
		result = cutRight(contentLine, keepWidth) + notifLine
	case TopCenterPosition, BottomCenterPosition:
		result = m.overlayCenter(contentLine, notifLine, notifWidth, contentWidth)
	default:
		result = notifLine + cutLeft(contentLine, notifWidth)
	}

end:
	return result
}

// overlayCenter overlays a notification line at the horizontal center of
// a content line. Falls back to left overlay if content is too short.
func (m NotifyModel) overlayCenter(
	contentLine, notifLine string,
	notifWidth, contentWidth int,
) (result string) {
	var leftPad int
	var contentLen int
	var left string
	var rightStart int
	var right string

	leftPad = (contentWidth - notifWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	contentLen = ansi.StringWidth(contentLine)

	if contentLen < leftPad {
		result = notifLine + cutLeft(contentLine, notifWidth)
		goto end
	}

	left = cutRight(contentLine, leftPad)
	rightStart = leftPad + notifWidth

	if rightStart < contentLen {
		right = cutLeft(contentLine, rightStart)
	}

	result = left + notifLine + right

end:
	return result
}

// tickMsg is the timer message that drives notice animation and expiry.
type tickMsg time.Time

// tickCmd returns a tea.Cmd that ticks every 100ms for notice updates.
func tickCmd() (cmd tea.Cmd) {
	cmd = tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
	return cmd
}
