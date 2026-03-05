package teagrid

import (
	"fmt"
	"sort"
)

// SortDirection indicates sort order.
type SortDirection int

const (
	SortDirectionAsc SortDirection = iota
	SortDirectionDesc
)

// SortColumn describes which column should be sorted and how.
type SortColumn struct {
	ColumnKey string
	Direction SortDirection
}

// SortByAsc sets the primary sort column in ascending order.
func (m GridModel) SortByAsc(columnKey string) GridModel {
	m.sortOrder = []SortColumn{
		{ColumnKey: columnKey, Direction: SortDirectionAsc},
	}
	m.visibleRowCacheUpdated = false
	return m
}

// SortByDesc sets the primary sort column in descending order.
func (m GridModel) SortByDesc(columnKey string) GridModel {
	m.sortOrder = []SortColumn{
		{ColumnKey: columnKey, Direction: SortDirectionDesc},
	}
	m.visibleRowCacheUpdated = false
	return m
}

// ThenSortByAsc adds a secondary ascending sort.
func (m GridModel) ThenSortByAsc(columnKey string) GridModel {
	m.sortOrder = append([]SortColumn{
		{ColumnKey: columnKey, Direction: SortDirectionAsc},
	}, m.sortOrder...)
	m.visibleRowCacheUpdated = false
	return m
}

// ThenSortByDesc adds a secondary descending sort.
func (m GridModel) ThenSortByDesc(columnKey string) GridModel {
	m.sortOrder = append([]SortColumn{
		{ColumnKey: columnKey, Direction: SortDirectionDesc},
	}, m.sortOrder...)
	m.visibleRowCacheUpdated = false
	return m
}

type sortableTable struct {
	rows     []Row
	byColumn SortColumn
}

func (s *sortableTable) Len() int      { return len(s.rows) }
func (s *sortableTable) Swap(i, j int) { s.rows[i], s.rows[j] = s.rows[j], s.rows[i] }

func (s *sortableTable) extractString(i int, column string) string {
	iData, exists := s.rows[i].Data[column]
	if !exists {
		return ""
	}

	switch v := iData.(type) {
	case CellValue:
		if v.SortValue != nil {
			return fmt.Sprintf("%v", v.SortValue)
		}
		return fmt.Sprintf("%v", v.Data)

	case string:
		return v

	default:
		return fmt.Sprintf("%v", v)
	}
}

func (s *sortableTable) extractNumber(i int, column string) (float64, bool) {
	iData, exists := s.rows[i].Data[column]
	if !exists {
		return 0, false
	}

	return asNumber(iData)
}

func (s *sortableTable) Less(first, second int) bool {
	firstNum, firstOK := s.extractNumber(first, s.byColumn.ColumnKey)
	secondNum, secondOK := s.extractNumber(second, s.byColumn.ColumnKey)

	if firstOK && secondOK {
		if s.byColumn.Direction == SortDirectionAsc {
			return firstNum < secondNum
		}
		return firstNum > secondNum
	}

	firstVal := s.extractString(first, s.byColumn.ColumnKey)
	secondVal := s.extractString(second, s.byColumn.ColumnKey)

	if s.byColumn.Direction == SortDirectionAsc {
		return firstVal < secondVal
	}
	return firstVal > secondVal
}

func getSortedRows(sortOrder []SortColumn, rows []Row) []Row {
	if len(sortOrder) == 0 {
		return rows
	}

	sortedRows := make([]Row, len(rows))
	copy(sortedRows, rows)

	for _, byColumn := range sortOrder {
		sorted := &sortableTable{
			rows:     sortedRows,
			byColumn: byColumn,
		}
		sort.Stable(sorted)
		sortedRows = sorted.rows
	}

	return sortedRows
}
