package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestBorderDefault(t *testing.T) {
	bc := BorderDefault()

	assert.True(t, bc.HasOuterBorder())
	assert.True(t, bc.HasHeaderSeparator())
	assert.True(t, bc.HasInnerDividers())
	assert.True(t, bc.HasFooterSeparator())
	assert.Equal(t, 2, bc.OuterWidth())
	assert.Equal(t, 1, bc.InnerDividerWidth())
	assert.Equal(t, "┏", bc.Chars.TopLeft)
	assert.Equal(t, "━", bc.Chars.Horizontal)
}

func TestBorderRounded(t *testing.T) {
	bc := BorderRounded()

	assert.True(t, bc.HasOuterBorder())
	assert.Equal(t, "╭", bc.Chars.TopLeft)
	assert.Equal(t, "─", bc.Chars.Horizontal)
	assert.Equal(t, "│", bc.Chars.InnerDivider)
}

func TestBorderless(t *testing.T) {
	bc := Borderless()

	assert.False(t, bc.HasOuterBorder())
	assert.False(t, bc.HasHeaderSeparator())
	assert.False(t, bc.HasInnerDividers())
	assert.False(t, bc.HasFooterSeparator())
	assert.Equal(t, 0, bc.OuterWidth())
	assert.Equal(t, 0, bc.InnerDividerWidth())
}

func TestBorderMinimal(t *testing.T) {
	bc := BorderMinimal()

	assert.False(t, bc.HasOuterBorder())
	assert.True(t, bc.HasHeaderSeparator())
	assert.False(t, bc.HasInnerDividers())
	assert.False(t, bc.HasFooterSeparator())
}

func TestBorderConfigWithRegions(t *testing.T) {
	bc := BorderDefault().
		WithOuter(BorderRegion{Visible: false}).
		WithInner(BorderRegion{Visible: false})

	assert.False(t, bc.HasOuterBorder())
	assert.True(t, bc.HasHeaderSeparator())
	assert.False(t, bc.HasInnerDividers())
	assert.True(t, bc.HasFooterSeparator())
}

func TestBorderConfigWithChars(t *testing.T) {
	custom := BorderChars{
		Horizontal: "=",
		Vertical:   "|",
	}
	bc := BorderDefault().WithChars(custom)
	assert.Equal(t, "=", bc.Chars.Horizontal)
	assert.Equal(t, "|", bc.Chars.Vertical)
}

func TestBorderRegionWithStyle(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	region := BorderRegion{Visible: true}.WithStyle(style)
	assert.Equal(t, style, region.Style)
}

func TestBorderRegionWithVisible(t *testing.T) {
	region := BorderRegion{Visible: true}.WithVisible(false)
	assert.False(t, region.Visible)
}

func TestBorderConfigImmutability(t *testing.T) {
	original := BorderDefault()
	modified := original.WithOuter(BorderRegion{Visible: false})

	assert.True(t, original.HasOuterBorder(), "original should be unchanged")
	assert.False(t, modified.HasOuterBorder())
}
