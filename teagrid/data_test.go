package teagrid

import (
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/stretchr/testify/assert"
)

func TestAsInt(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int64
		ok       bool
	}{
		{"int", 42, 42, true},
		{"int8", int8(8), 8, true},
		{"int16", int16(16), 16, true},
		{"int32", int32(32), 32, true},
		{"int64", int64(64), 64, true},
		{"uint", uint(10), 10, true},
		{"uint8", uint8(8), 8, true},
		{"uint16", uint16(16), 16, true},
		{"uint32", uint32(32), 32, true},
		{"uint64", uint64(64), 64, true},
		{"duration", time.Second, int64(time.Second), true},
		{"string fails", "hello", 0, false},
		{"float fails", 3.14, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := asInt(tt.input)
			assert.Equal(t, tt.ok, ok)
			if ok {
				assert.Equal(t, tt.expected, val)
			}
		})
	}
}

func TestAsIntCellValue(t *testing.T) {
	t.Run("uses Data when no SortValue", func(t *testing.T) {
		cv := NewCellValue(42, lipgloss.NewStyle())
		val, ok := asInt(cv)
		assert.True(t, ok)
		assert.Equal(t, int64(42), val)
	})

	t.Run("uses SortValue when set", func(t *testing.T) {
		cv := NewCellValueWithSortKey("display", 99, lipgloss.NewStyle())
		val, ok := asInt(cv)
		assert.True(t, ok)
		assert.Equal(t, int64(99), val)
	})
}

func TestAsNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		ok       bool
	}{
		{"float32", float32(3.14), float64(float32(3.14)), true},
		{"float64", 3.14, 3.14, true},
		{"int via asInt", 42, 42.0, true},
		{"string fails", "hello", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := asNumber(tt.input)
			assert.Equal(t, tt.ok, ok)
			if ok {
				assert.InDelta(t, tt.expected, val, 0.001)
			}
		})
	}
}

func TestAsNumberCellValue(t *testing.T) {
	t.Run("uses SortValue when set", func(t *testing.T) {
		cv := NewCellValueWithSortKey("display", 3.14, lipgloss.NewStyle())
		val, ok := asNumber(cv)
		assert.True(t, ok)
		assert.InDelta(t, 3.14, val, 0.001)
	})
}
