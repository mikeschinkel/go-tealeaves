package teagrid

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewCellValue(t *testing.T) {
	style := lipgloss.NewStyle().Bold(true)
	cv := NewCellValue("hello", style)

	assert.Equal(t, "hello", cv.Data)
	assert.Nil(t, cv.SortValue)
	assert.Nil(t, cv.StyleFunc)
	assert.False(t, cv.HasSpans())
}

func TestNewCellValueWithStyleFunc(t *testing.T) {
	fn := func(input CellStyleInput) lipgloss.Style {
		return lipgloss.NewStyle().Bold(true)
	}
	cv := NewCellValueWithStyleFunc("data", fn)

	assert.Equal(t, "data", cv.Data)
	assert.NotNil(t, cv.StyleFunc)
}

func TestNewCellValueWithSortKey(t *testing.T) {
	cv := NewCellValueWithSortKey("Jan 1", 1704067200, lipgloss.NewStyle())

	assert.Equal(t, "Jan 1", cv.Data)
	assert.Equal(t, 1704067200, cv.SortValue)
}

func TestCellValueSortableValue(t *testing.T) {
	t.Run("returns Data when no SortValue", func(t *testing.T) {
		cv := NewCellValue("hello", lipgloss.NewStyle())
		assert.Equal(t, "hello", cv.SortableValue())
	})

	t.Run("returns SortValue when set", func(t *testing.T) {
		cv := NewCellValueWithSortKey("display", 42, lipgloss.NewStyle())
		assert.Equal(t, 42, cv.SortableValue())
	})
}

func TestNewCellValueWithSpans(t *testing.T) {
	spans := []Span{
		NewSpan("hello", lipgloss.NewStyle().Bold(true)),
		NewSpan(" world", lipgloss.NewStyle()),
	}
	cv := NewCellValueWithSpans(spans, lipgloss.NewStyle())

	assert.True(t, cv.HasSpans())
	assert.Len(t, cv.Spans, 2)
	assert.Equal(t, "hello", cv.Spans[0].Text)
	assert.Equal(t, " world", cv.Spans[1].Text)
}

func TestNewSpan(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	span := NewSpan("red text", style)

	assert.Equal(t, "red text", span.Text)
}

func TestCellValueHasSpans(t *testing.T) {
	t.Run("false when no spans", func(t *testing.T) {
		cv := NewCellValue("data", lipgloss.NewStyle())
		assert.False(t, cv.HasSpans())
	})

	t.Run("true with spans", func(t *testing.T) {
		cv := CellValue{Spans: []Span{{Text: "x"}}}
		assert.True(t, cv.HasSpans())
	})
}
