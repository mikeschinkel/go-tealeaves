package teagrid

import (
	"fmt"
	"strings"
)

// FilterFuncInput is the input to a FilterFunc.
type FilterFuncInput struct {
	Columns        []Column
	Row            Row
	GlobalMetadata map[string]any
	Filter         string
}

// FilterFunc returns true if the row should be visible, false to hide it.
type FilterFunc func(FilterFuncInput) bool

func (m Model) getFilteredRows(rows []Row) []Row {
	filterValue := m.filterTextInput.Value()
	if !m.filtered || filterValue == "" {
		return rows
	}

	filterFunc := m.filterFunc
	if filterFunc == nil {
		filterFunc = filterFuncContains
	}

	filtered := make([]Row, 0)
	for _, row := range rows {
		if filterFunc(FilterFuncInput{
			Columns:        m.columns,
			Row:            row,
			Filter:         filterValue,
			GlobalMetadata: m.metadata,
		}) {
			filtered = append(filtered, row)
		}
	}

	return filtered
}

// filterFuncContains performs case-insensitive substring matching
// across all filterable columns.
func filterFuncContains(input FilterFuncInput) bool {
	if input.Filter == "" {
		return true
	}

	checkedAny := false
	filterLower := strings.ToLower(input.Filter)

	for _, column := range input.Columns {
		if !column.filterable {
			continue
		}

		checkedAny = true

		data, ok := input.Row.Data[column.key]
		if !ok {
			continue
		}

		// Extract CellValue data
		if cv, ok := data.(CellValue); ok {
			data = cv.Data
		}

		var target string
		switch v := data.(type) {
		case string:
			target = v
		case fmt.Stringer:
			target = v.String()
		default:
			target = fmt.Sprintf("%v", data)
		}

		if strings.Contains(strings.ToLower(target), filterLower) {
			return true
		}
	}

	return !checkedAny
}

// filterFuncFuzzy performs case-insensitive fuzzy (subsequence) matching
// across all filterable columns concatenated.
func filterFuncFuzzy(input FilterFuncInput) bool {
	filter := strings.TrimSpace(input.Filter)
	if filter == "" {
		return true
	}

	var builder strings.Builder
	for _, col := range input.Columns {
		if !col.filterable {
			continue
		}

		value, ok := input.Row.Data[col.key]
		if !ok {
			continue
		}

		if cv, ok := value.(CellValue); ok {
			value = cv.Data
		}

		builder.WriteString(fmt.Sprint(value))
		builder.WriteByte(' ')
	}

	haystack := strings.ToLower(builder.String())
	if haystack == "" {
		return false
	}

	for _, token := range strings.Fields(strings.ToLower(filter)) {
		if !fuzzySubsequenceMatch(haystack, token) {
			return false
		}
	}

	return true
}

// fuzzySubsequenceMatch returns true if all runes in needle appear in order
// within haystack. Case must be normalized by caller.
func fuzzySubsequenceMatch(haystack, needle string) bool {
	if needle == "" {
		return true
	}

	haystackRunes := []rune(haystack)
	needleRunes := []rune(needle)
	hi, ni := 0, 0

	for hi < len(haystackRunes) && ni < len(needleRunes) {
		if haystackRunes[hi] == needleRunes[ni] {
			ni++
		}
		hi++
	}

	return ni == len(needleRunes)
}
