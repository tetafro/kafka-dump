package main

import "strconv"

func equal(a, b interface{}) bool {
	f1, ok1 := toFloat(a)
	f2, ok2 := toFloat(b)
	if ok1 && ok2 && f1 == f2 {
		return true
	}

	b1, ok1 := toBool(a)
	b2, ok2 := toBool(b)
	if ok1 && ok2 && b1 == b2 {
		return true
	}

	s1, ok1 := toString(a)
	s2, ok2 := toString(b)
	if ok1 && ok2 && s1 == s2 {
		return true
	}

	return false
}

func toFloat(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

func toBool(val interface{}) (bool, bool) {
	switch v := val.(type) {
	case bool:
		return v, true
	case string:
		if b, err := strconv.ParseBool(v); err == nil {
			return b, true
		}
		return false, false
	default:
		return false, false
	}
}

func toString(val interface{}) (string, bool) {
	switch v := val.(type) {
	case string:
		return v, true
	default:
		return "", false
	}
}
