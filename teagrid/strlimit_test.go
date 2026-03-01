package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLimitStr(t *testing.T) {
	t.Run("short string unchanged", func(t *testing.T) {
		assert.Equal(t, "hello", limitStr("hello", 10))
	})

	t.Run("truncates long string", func(t *testing.T) {
		result := limitStr("hello world", 5)
		assert.LessOrEqual(t, len(result), 8) // may include ellipsis + ANSI
	})

	t.Run("zero max returns empty", func(t *testing.T) {
		assert.Equal(t, "", limitStr("hello", 0))
	})

	t.Run("newline replaced with ellipsis", func(t *testing.T) {
		result := limitStr("line1\nline2", 20)
		assert.NotContains(t, result, "\n")
		assert.Contains(t, result, "…")
	})
}
