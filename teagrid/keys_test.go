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
	assert.NotEmpty(t, km.CellLeft.Keys())
	assert.NotEmpty(t, km.CellRight.Keys())
	assert.NotEmpty(t, km.CellSelect.Keys())
	assert.NotEmpty(t, km.Filter.Keys())
	assert.NotEmpty(t, km.FilterBlur.Keys())
	assert.NotEmpty(t, km.FilterClear.Keys())
	assert.NotEmpty(t, km.ScrollRight.Keys())
	assert.NotEmpty(t, km.ScrollLeft.Keys())
}

func TestDefaultKeyMapCellNavSeparateFromPagination(t *testing.T) {
	km := DefaultKeyMap()

	// CellLeft/CellRight should use arrow keys, not pgup/pgdown
	cellLeftKeys := km.CellLeft.Keys()
	cellRightKeys := km.CellRight.Keys()
	pageDownKeys := km.PageDown.Keys()
	pageUpKeys := km.PageUp.Keys()

	// Verify no overlap between cell nav and pagination
	for _, ck := range cellLeftKeys {
		for _, pk := range pageUpKeys {
			assert.NotEqual(t, ck, pk, "CellLeft should not share keys with PageUp")
		}
	}
	for _, ck := range cellRightKeys {
		for _, pk := range pageDownKeys {
			assert.NotEqual(t, ck, pk, "CellRight should not share keys with PageDown")
		}
	}
}
