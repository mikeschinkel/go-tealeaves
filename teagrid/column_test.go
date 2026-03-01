package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewColumnDefaults(t *testing.T) {
	col := NewColumn("name", "Name", 20)

	assert.Equal(t, "name", col.Key())
	assert.Equal(t, "Name", col.Title())
	assert.Equal(t, 20, col.Width())
	assert.Equal(t, 0, col.FlexFactor())
	assert.False(t, col.IsFlex())
	assert.Equal(t, lipgloss.Left, col.Alignment())
	assert.Equal(t, 1, col.PaddingLeft())
	assert.Equal(t, 0, col.PaddingRight())
	assert.False(t, col.Filterable())
	assert.Equal(t, "", col.FmtString())
}

func TestNewFlexColumnDefaults(t *testing.T) {
	col := NewFlexColumn("desc", "Description", 3)

	assert.Equal(t, "desc", col.Key())
	assert.Equal(t, "Description", col.Title())
	assert.Equal(t, 0, col.Width())
	assert.Equal(t, 3, col.FlexFactor())
	assert.True(t, col.IsFlex())
	assert.Equal(t, lipgloss.Left, col.Alignment())
}

func TestNewFlexColumnMinFactor(t *testing.T) {
	col := NewFlexColumn("x", "X", 0)
	assert.Equal(t, 1, col.FlexFactor(), "flex factor should be at least 1")
}

func TestColumnRenderWidth(t *testing.T) {
	t.Run("default padding", func(t *testing.T) {
		col := NewColumn("x", "X", 10)
		// paddingLeft=1 + width=10 + paddingRight=0 = 11
		assert.Equal(t, 11, col.RenderWidth())
	})

	t.Run("custom padding", func(t *testing.T) {
		col := NewColumn("x", "X", 10).WithPadding(2, 3)
		// paddingLeft=2 + width=10 + paddingRight=3 = 15
		assert.Equal(t, 15, col.RenderWidth())
	})
}

func TestColumnWithPadding(t *testing.T) {
	col := NewColumn("x", "X", 10).WithPadding(3, 2)
	assert.Equal(t, 3, col.PaddingLeft())
	assert.Equal(t, 2, col.PaddingRight())
}

func TestColumnWithPaddingLeft(t *testing.T) {
	col := NewColumn("x", "X", 10).WithPaddingLeft(5)
	assert.Equal(t, 5, col.PaddingLeft())
	assert.Equal(t, 0, col.PaddingRight())
}

func TestColumnWithPaddingRight(t *testing.T) {
	col := NewColumn("x", "X", 10).WithPaddingRight(5)
	assert.Equal(t, 1, col.PaddingLeft())
	assert.Equal(t, 5, col.PaddingRight())
}

func TestColumnWithAlignment(t *testing.T) {
	col := NewColumn("x", "X", 10).WithAlignment(lipgloss.Right)
	assert.Equal(t, lipgloss.Right, col.Alignment())
}

func TestColumnWithFiltered(t *testing.T) {
	col := NewColumn("x", "X", 10).WithFiltered(true)
	assert.True(t, col.Filterable())
}

func TestColumnWithFormatString(t *testing.T) {
	col := NewColumn("x", "X", 10).WithFormatString("%.2f")
	assert.Equal(t, "%.2f", col.FmtString())
}

func TestColumnWithStyle(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)
	col := NewColumn("x", "X", 10).WithStyle(style)
	assert.Equal(t, style, col.Style())
}

func TestColumnImmutability(t *testing.T) {
	original := NewColumn("x", "X", 10)
	modified := original.WithPaddingLeft(5)

	assert.Equal(t, 1, original.PaddingLeft(), "original should be unchanged")
	assert.Equal(t, 5, modified.PaddingLeft())
}
