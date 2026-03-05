package teatree

import (
	"charm.land/lipgloss/v2"
	"github.com/mikeschinkel/go-tealeaves/teautils"
)

// DefaultDrillDownPathStyle returns default path item styling
func DefaultDrillDownPathStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

// DefaultDrillDownSelectedStyle returns default selected item styling
func DefaultDrillDownSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Bold(true)
}

// DefaultDrillDownBorderStyle returns default border styling
func DefaultDrillDownBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
}

// WithDrillDownTheme applies a theme to the drill-down model's styles
func (m DrillDownModel[T]) WithDrillDownTheme(theme teautils.Theme) DrillDownModel[T] {
	m.SelectedStyle = theme.SelectedItem
	m.PathStyle = theme.Item
	m.BorderStyle = theme.Border
	return m
}
