package teagrid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultKeyMap(t *testing.T) {
	km := DefaultKeyMap()

	assert.NotEmpty(t, km.RowDown.Keys())
	assert.NotEmpty(t, km.RowUp.Keys())
	assert.NotEmpty(t, km.RowSelectToggle.Keys())
	assert.NotEmpty(t, km.PageDown.Keys())
	assert.NotEmpty(t, km.PageUp.Keys())
	assert.NotEmpty(t, km.PageFirst.Keys())
	assert.NotEmpty(t, km.PageLast.Keys())
	assert.NotEmpty(t, km.ColLeft.Keys())
	assert.NotEmpty(t, km.ColRight.Keys())
	assert.NotEmpty(t, km.ColSelect.Keys())
	assert.NotEmpty(t, km.Filter.Keys())
	assert.NotEmpty(t, km.FilterBlur.Keys())
	assert.NotEmpty(t, km.FilterClear.Keys())
	assert.NotEmpty(t, km.ScrollRight.Keys())
	assert.NotEmpty(t, km.ScrollLeft.Keys())
}

func TestDefaultKeyMapColNavSeparateFromPagination(t *testing.T) {
	km := DefaultKeyMap()

	// ColLeft/ColRight should use arrow keys, not pgup/pgdown
	colLeftKeys := km.ColLeft.Keys()
	colRightKeys := km.ColRight.Keys()
	pageDownKeys := km.PageDown.Keys()
	pageUpKeys := km.PageUp.Keys()

	// Verify no overlap between col nav and pagination
	for _, ck := range colLeftKeys {
		for _, pk := range pageUpKeys {
			assert.NotEqual(t, ck, pk, "ColLeft should not share keys with PageUp")
		}
	}
	for _, ck := range colRightKeys {
		for _, pk := range pageDownKeys {
			assert.NotEqual(t, ck, pk, "ColRight should not share keys with PageDown")
		}
	}
}
