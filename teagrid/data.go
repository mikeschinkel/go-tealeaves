package teagrid

import "time"

// asInt attempts to convert any numeric type to int64.
//
//nolint:cyclop
func asInt(data any) (int64, bool) {
	switch val := data.(type) {
	case int:
		return int64(val), true

	case int8:
		return int64(val), true

	case int16:
		return int64(val), true

	case int32:
		return int64(val), true

	case int64:
		return val, true

	case uint:
		// #nosec: G115
		return int64(val), true

	case uint8:
		return int64(val), true

	case uint16:
		return int64(val), true

	case uint32:
		return int64(val), true

	case uint64:
		// #nosec: G115
		return int64(val), true

	case time.Duration:
		return int64(val), true

	case CellValue:
		if val.SortValue != nil {
			return asInt(val.SortValue)
		}
		return asInt(val.Data)
	}

	return 0, false
}

// asNumber attempts to convert any numeric type to float64.
func asNumber(data any) (float64, bool) {
	switch val := data.(type) {
	case float32:
		return float64(val), true

	case float64:
		return val, true

	case CellValue:
		if val.SortValue != nil {
			return asNumber(val.SortValue)
		}
		return asNumber(val.Data)
	}

	intVal, isInt := asInt(data)

	return float64(intVal), isInt
}
