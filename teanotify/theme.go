package teanotify

import (
	"github.com/lucasb-eyer/go-colorful"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// WithTheme returns a copy with default notice colors derived from the theme's
// palette status colors. Only affects registered default notice types (Info,
// Warn, Error, Debug). Custom notice types are unchanged.
func (m NotifyModel) WithTheme(theme teautils.Theme) (out NotifyModel) {
	out = m
	p := theme.System

	colorMap := map[NoticeKey]colorful.Color{
		InfoKey:  colorfulFromImageColor(p.StatusInfo),
		WarnKey:  colorfulFromImageColor(p.StatusWarn),
		ErrorKey: colorfulFromImageColor(p.StatusError),
		DebugKey: colorfulFromImageColor(p.AccentAlt),
	}

	if out.noticeTypes == nil {
		return out
	}

	newTypes := make(map[NoticeKey]NoticeDefinition, len(out.noticeTypes))
	for k, v := range out.noticeTypes {
		cf, ok := colorMap[k]
		if ok {
			hex := cf.Hex()
			v.ForeColor = hex
			parsedColors[hex] = cf
		}
		newTypes[k] = v
	}
	out.noticeTypes = newTypes

	return out
}

// colorfulFromImageColor converts an image/color.Color to a colorful.Color.
func colorfulFromImageColor(c interface{ RGBA() (uint32, uint32, uint32, uint32) }) colorful.Color {
	if c == nil {
		return colorful.Color{}
	}
	cf, _ := colorful.MakeColor(c)
	return cf
}
