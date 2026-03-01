package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCD(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"both zero", 0, 0, 0},
		{"x zero", 0, 5, 5},
		{"y zero", 7, 0, 7},
		{"equal", 6, 6, 6},
		{"coprime", 7, 13, 1},
		{"common factor", 12, 8, 4},
		{"one divides other", 3, 9, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, gcd(tt.x, tt.y))
		})
	}
}
