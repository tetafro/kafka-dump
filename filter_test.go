package main

import "testing"

func TestCheck(t *testing.T) {
	testCases := []struct {
		name   string
		msg    Message
		filter map[string]interface{}
		result bool
	}{
		{
			name: "contains correct field #1",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 10,
			}},
			filter: map[string]interface{}{
				"type": "foo",
			},
			result: true,
		},
		{
			name: "contains correct field #2",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 10,
			}},
			filter: map[string]interface{}{
				"value": 10,
			},
			result: true,
		},
		{
			name: "contains field with wrong value",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 10,
			}},
			filter: map[string]interface{}{
				"type": "bar",
			},
			result: false,
		},
		{
			name: "contains field with one of correct values #1",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 30,
			}},
			filter: map[string]interface{}{
				"value": []interface{}{10, 20, 30},
			},
			result: true,
		},
		{
			name: "contains field with one of correct values #2",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 30,
			}},
			filter: map[string]interface{}{
				"value": []interface{}{10, 20, 30},
			},
			result: true,
		},
		{
			name: "satisfies complex filter",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 30,
				"time":  1609336660,
			}},
			filter: map[string]interface{}{
				"type":  "foo",
				"value": []interface{}{10, 20, 30},
			},
			result: true,
		},
		{
			name: "doesn't satisfy complex filter",
			msg: Message{data: map[string]interface{}{
				"type":  "foo",
				"value": 30,
				"time":  1609336660,
			}},
			filter: map[string]interface{}{
				"type":  "foo",
				"value": []interface{}{10, 20},
			},
			result: false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			f := NewFieldFilter(tt.filter)
			if result := f.Check(tt.msg); result != tt.result {
				t.Fatalf("Expected %v, got %v", tt.result, result)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	testCases := []struct {
		name   string
		in1    interface{}
		in2    interface{}
		result bool
	}{
		{
			name:   "int, int => equal",
			in1:    int(10),
			in2:    int(10),
			result: true,
		},
		{
			name:   "int, float => equal",
			in1:    int(10),
			in2:    float64(10),
			result: true,
		},
		{
			name:   "int, bool => invalid",
			in1:    int(1),
			in2:    true,
			result: false,
		},
		{
			name:   "int, string => equal",
			in1:    int(20),
			in2:    "20",
			result: true,
		},
		{
			name:   "string, int => invalid",
			in1:    "hello",
			in2:    int(10),
			result: false,
		},
		{
			name:   "string, string => equal",
			in1:    "hello",
			in2:    "hello",
			result: true,
		},
		{
			name:   "string, string => not equal",
			in1:    "hello",
			in2:    "world",
			result: false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if result := equal(tt.in1, tt.in2); result != tt.result {
				t.Fatalf("Expected %v, got %v", tt.result, result)
			}
		})
	}
}

func TestContain(t *testing.T) {
	testCases := []struct {
		name    string
		slice   interface{}
		element interface{}
		result  bool
	}{
		{
			name:    "[]int, int => true",
			slice:   []interface{}{int(10), int(20), int(30)},
			element: int(20),
			result:  true,
		},
		{
			name:    "[]int, int => false",
			slice:   []interface{}{int(10), int(20), int(30)},
			element: int(40),
			result:  false,
		},
		{
			name:    "[]interface{}, int => true",
			slice:   []interface{}{10, "20", 30, struct{}{}},
			element: int(20),
			result:  true,
		},
		{
			name:    "[]interface{}, int => false",
			slice:   []interface{}{10, "20", 30, struct{}{}},
			element: int(25),
			result:  false,
		},
		{
			name:    "not a []interface{}",
			slice:   []int{10, 20, 30},
			element: int(10),
			result:  false,
		},
		{
			name:    "not a slice",
			slice:   map[int]int{10: 10},
			element: int(10),
			result:  false,
		},
		{
			name:    "nil slice",
			slice:   nil,
			element: int(10),
			result:  false,
		},
		{
			name:    "nil element",
			slice:   []int{10},
			element: nil,
			result:  false,
		},
		{
			name:    "all nils",
			slice:   nil,
			element: nil,
			result:  false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if result := contain(tt.slice, tt.element); result != tt.result {
				t.Fatalf("Expected %v, got %v", tt.result, result)
			}
		})
	}
}

func TestToFloat(t *testing.T) {
	testCases := []struct {
		name string
		in   interface{}
		out  float64
		ok   bool
	}{
		{
			name: "float32",
			in:   float32(10),
			out:  10,
			ok:   true,
		},
		{
			name: "float64",
			in:   float64(10),
			out:  10,
			ok:   true,
		},
		{
			name: "int",
			in:   int(10),
			out:  10,
			ok:   true,
		},
		{
			name: "int64",
			in:   int64(10),
			out:  10,
			ok:   true,
		},
		{
			name: "string",
			in:   "10",
			out:  10,
			ok:   true,
		},
		{
			name: "invalid string",
			in:   "hello",
			out:  0,
			ok:   false,
		},
		{
			name: "invalid type",
			in:   struct{}{},
			out:  0,
			ok:   false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out, ok := toFloat(tt.in)
			if ok != tt.ok {
				t.Fatalf("Expected %v, got %v", tt.ok, ok)
			}
			if out != tt.out {
				t.Fatalf("Expected %f, got %f", tt.out, out)
			}
		})
	}
}

func TestToBool(t *testing.T) {
	testCases := []struct {
		name string
		in   interface{}
		out  bool
		ok   bool
	}{
		{
			name: "bool",
			in:   true,
			out:  true,
			ok:   true,
		},
		{
			name: "string",
			in:   "true",
			out:  true,
			ok:   true,
		},
		{
			name: "invalid string",
			in:   "hello",
			out:  false,
			ok:   false,
		},
		{
			name: "invalid type",
			in:   struct{}{},
			out:  false,
			ok:   false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out, ok := toBool(tt.in)
			if ok != tt.ok {
				t.Fatalf("Expected %v, got %v", tt.ok, ok)
			}
			if out != tt.out {
				t.Fatalf("Expected %v, got %v", tt.out, out)
			}
		})
	}
}

func TestToString(t *testing.T) {
	testCases := []struct {
		name string
		in   interface{}
		out  string
		ok   bool
	}{
		{
			name: "string",
			in:   "hello",
			out:  "hello",
			ok:   true,
		},
		{
			name: "int",
			in:   int(10),
			out:  "",
			ok:   false,
		},
		{
			name: "invalid type",
			in:   struct{}{},
			out:  "",
			ok:   false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			out, ok := toString(tt.in)
			if ok != tt.ok {
				t.Fatalf("Expected %v, got %v", tt.ok, ok)
			}
			if out != tt.out {
				t.Fatalf("Expected %v, got %v", tt.out, out)
			}
		})
	}
}
