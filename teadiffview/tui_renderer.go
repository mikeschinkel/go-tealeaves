package teadiffview

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// TUIRenderer renders diffs using lipgloss for terminal output.
type TUIRenderer struct {
	FileHeaderColor    color.Color
	BlockHeaderColor   color.Color
	ContextColor       color.Color
	AddedColor         color.Color
	DeletedColor       color.Color
	NewStatusColor     color.Color
	DeletedStatusColor color.Color
	NewBgColor         color.Color
	DeletedBgColor     color.Color
}

// TUIRendererArgs holds optional configuration for NewTUIRenderer.
type TUIRendererArgs struct {
	FileHeaderColor    color.Color
	BlockHeaderColor   color.Color
	ContextColor       color.Color
	AddedColor         color.Color
	DeletedColor       color.Color
	NewStatusColor     color.Color
	DeletedStatusColor color.Color
	NewBgColor         color.Color
	DeletedBgColor     color.Color
}

// NewTUIRenderer creates a TUIRenderer with default colors.
// Pass non-nil args to override specific colors.
func NewTUIRenderer(args *TUIRendererArgs) *TUIRenderer {
	r := &TUIRenderer{
		FileHeaderColor:    lipgloss.Color("244"),
		BlockHeaderColor:   lipgloss.Color("135"),
		ContextColor:       lipgloss.Color("240"),
		AddedColor:         lipgloss.Color("34"),
		DeletedColor:       lipgloss.Color("160"),
		NewStatusColor:     lipgloss.Color("46"),
		DeletedStatusColor: lipgloss.Color("168"),
		NewBgColor:         lipgloss.Color("22"),
		DeletedBgColor:     lipgloss.Color("52"),
	}
	if args == nil {
		return r
	}
	if args.FileHeaderColor != nil {
		r.FileHeaderColor = args.FileHeaderColor
	}
	if args.BlockHeaderColor != nil {
		r.BlockHeaderColor = args.BlockHeaderColor
	}
	if args.ContextColor != nil {
		r.ContextColor = args.ContextColor
	}
	if args.AddedColor != nil {
		r.AddedColor = args.AddedColor
	}
	if args.DeletedColor != nil {
		r.DeletedColor = args.DeletedColor
	}
	if args.NewStatusColor != nil {
		r.NewStatusColor = args.NewStatusColor
	}
	if args.DeletedStatusColor != nil {
		r.DeletedStatusColor = args.DeletedStatusColor
	}
	if args.NewBgColor != nil {
		r.NewBgColor = args.NewBgColor
	}
	if args.DeletedBgColor != nil {
		r.DeletedBgColor = args.DeletedBgColor
	}
	return r
}

func (r *TUIRenderer) RenderFileHeader(path string, status FileStatus, width int) string {
	var statusTag string
	var statusStyle lipgloss.Style

	headerStyle := lipgloss.NewStyle().Foreground(r.FileHeaderColor).Bold(true)

	switch status {
	case FileNew:
		statusStyle = lipgloss.NewStyle().Foreground(r.NewStatusColor).Bold(true)
		statusTag = statusStyle.Render(" [NEW]")
	case FileDeleted:
		statusStyle = lipgloss.NewStyle().Foreground(r.DeletedStatusColor).Bold(true)
		statusTag = statusStyle.Render(" [DELETED]")
	}

	header := headerStyle.Render(path) + statusTag
	if width > 0 {
		headerWidth := ansi.StringWidth(header)
		if headerWidth < width {
			header += strings.Repeat(" ", width-headerWidth)
		}
	}
	return header
}

func (r *TUIRenderer) RenderBlockHeader(blockType string, lineCount int) string {
	style := lipgloss.NewStyle().Foreground(r.BlockHeaderColor)
	return style.Render(fmt.Sprintf("  %s (%d lines)", blockType, lineCount))
}

func (r *TUIRenderer) RenderContextLine(line string, status FileStatus, width int) string {
	style := lipgloss.NewStyle().Foreground(r.ContextColor)
	bg := r.bgForStatus(status)
	if bg != nil {
		style = style.Background(bg)
	}
	rendered := style.Render("    " + line)
	return r.padToWidth(rendered, width, style)
}

func (r *TUIRenderer) RenderAddedLine(line string, status FileStatus, width int) string {
	style := lipgloss.NewStyle().Foreground(r.AddedColor)
	bg := r.bgForStatus(status)
	if bg != nil {
		style = style.Background(bg)
	}
	rendered := style.Render("  + " + line)
	return r.padToWidth(rendered, width, style)
}

func (r *TUIRenderer) RenderDeletedLine(line string, status FileStatus, width int) string {
	style := lipgloss.NewStyle().Foreground(r.DeletedColor)
	bg := r.bgForStatus(status)
	if bg != nil {
		style = style.Background(bg)
	}
	rendered := style.Render("  - " + line)
	return r.padToWidth(rendered, width, style)
}

func (r *TUIRenderer) RenderTruncation(status FileStatus) string {
	style := lipgloss.NewStyle().Foreground(r.ContextColor)
	bg := r.bgForStatus(status)
	if bg != nil {
		style = style.Background(bg)
	}
	return style.Render("    ...")
}

func (r *TUIRenderer) RenderSeparator() string {
	return ""
}

func (r *TUIRenderer) bgForStatus(status FileStatus) color.Color {
	switch status {
	case FileNew:
		return r.NewBgColor
	case FileDeleted:
		return r.DeletedBgColor
	default:
		return nil
	}
}

func (r *TUIRenderer) padToWidth(rendered string, width int, style lipgloss.Style) string {
	if width <= 0 {
		return rendered
	}
	renderedWidth := ansi.StringWidth(rendered)
	if renderedWidth < width {
		padding := strings.Repeat(" ", width-renderedWidth)
		rendered += style.Render(padding)
	}
	return rendered
}
