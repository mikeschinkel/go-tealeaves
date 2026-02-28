package teagrid

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple3x3WithSelectableDefaults(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).SelectableRows(true)

	const expectedTable = `┏━━━┳━━━━┳━━━━┳━━━━┓
┃[x]┃   1┃   2┃   3┃
┣━━━╋━━━━╋━━━━╋━━━━┫
┃[ ]┃ 1,1┃ 2,1┃ 3,1┃
┃[ ]┃ 1,2┃ 2,2┃ 3,2┃
┃[ ]┃ 1,3┃ 2,3┃ 3,3┃
┗━━━┻━━━━┻━━━━┻━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimple3x3WithCustomSelectableText(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).
		SelectableRows(true).
		WithSelectedText(" ", "✓")

	const expectedTable = `┏━┳━━━━┳━━━━┳━━━━┓
┃✓┃   1┃   2┃   3┃
┣━╋━━━━╋━━━━╋━━━━┫
┃ ┃ 1,1┃ 2,1┃ 3,1┃
┃ ┃ 1,2┃ 2,2┃ 3,2┃
┃ ┃ 1,3┃ 2,3┃ 3,3┃
┗━┻━━━━┻━━━━┻━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestSimple3x3WithCustomSelectableTextAndFooter(t *testing.T) {
	model := New([]Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	})

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		WithStaticFooter("Footer")

	const expectedTable = `┏━┳━━━━┳━━━━┳━━━━┓
┃✓┃   1┃   2┃   3┃
┣━╋━━━━╋━━━━╋━━━━┫
┃ ┃ 1,1┃ 2,1┃ 3,1┃
┃ ┃ 1,2┃ 2,2┃ 3,2┃
┃ ┃ 1,3┃ 2,3┃ 3,3┃
┣━┻━━━━┻━━━━┻━━━━┫
┃          Footer┃
┗━━━━━━━━━━━━━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}

func TestRegeneratingColumnsKeepsSelectableText(t *testing.T) {
	columns := []Column{
		NewColumn("1", "1", 4),
		NewColumn("2", "2", 4),
		NewColumn("3", "3", 4),
	}

	model := New(columns)

	rows := []Row{}

	for rowIndex := 1; rowIndex <= 3; rowIndex++ {
		rowData := RowData{}

		for columnIndex := 1; columnIndex <= 3; columnIndex++ {
			id := fmt.Sprintf("%d", columnIndex)

			rowData[id] = fmt.Sprintf("%d,%d", columnIndex, rowIndex)
		}

		rows = append(rows, NewRow(rowData))
	}

	model = model.WithRows(rows).
		SelectableRows(true).
		WithSelectedText(" ", "✓").
		WithColumns(columns)

	const expectedTable = `┏━┳━━━━┳━━━━┳━━━━┓
┃✓┃   1┃   2┃   3┃
┣━╋━━━━╋━━━━╋━━━━┫
┃ ┃ 1,1┃ 2,1┃ 3,1┃
┃ ┃ 1,2┃ 2,2┃ 3,2┃
┃ ┃ 1,3┃ 2,3┃ 3,3┃
┗━┻━━━━┻━━━━┻━━━━┛`

	rendered := model.View()

	assert.Equal(t, expectedTable, rendered)
}
