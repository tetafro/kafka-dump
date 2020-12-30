package main

import (
	"strconv"
)

// Filter describes a way to decide whether a message should be saved or not.
type Filter interface {
	Check(Message) bool
}

// FieldFilter is a filter that makes a decision based on message's fields.
type FieldFilter struct {
	fields map[string]interface{}
}

// NewFieldFilter creates new field filter.
func NewFieldFilter(fields map[string]interface{}) *FieldFilter {
	return &FieldFilter{fields: fields}
}

// Check decides whether a message should be saved or not.
func (f *FieldFilter) Check(msg Message) bool {
	for k, v := range f.fields {
		data, ok := msg.data[k]
		if !ok {
			return false
		}
		if !equal(v, data) && !contain(v, data) {
			return false
		}
	}
	return true
}

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

func contain(a, b interface{}) bool {
	slice, ok := a.([]interface{})
	if !ok {
		return false
	}
	for i := range slice {
		res := equal(slice[i], b)
		if res {
			return true
		}
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
