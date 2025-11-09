package goutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ------------------- Value Conversions -------------------
func (d *TreeMap) AsString() (string, error) {
	switch v := d.value.(type) {
	case string:
		return v, nil
	case int, int64, float64, bool:
		return fmt.Sprint(v), nil
	default:
		return "", fmt.Errorf("cannot convert to string: %T", v)
	}
}

func (d *TreeMap) AsInt() (int64, error) {
	if d.err != nil {
		return 0, d.err
	}
	switch v := d.value.(type) {
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert to int: %T", v)
	}
}

func (d *TreeMap) AsFloat() (float64, error) {
	if d.err != nil {
		return 0, d.err
	}
	switch v := d.value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case string:
		return strconv.ParseFloat(v, 64)
	default:
		return 0, fmt.Errorf("cannot convert to float64: %T", v)
	}
}

func (d *TreeMap) AsBool() (bool, error) {
	if d.err != nil {
		return false, d.err
	}
	switch v := d.value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(v)
	default:
		return false, fmt.Errorf("cannot convert to bool: %T", v)
	}
}

func (d *TreeMap) AsAny() (any, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.value, nil
}

// ------------------- Struct / Slice Conversions -------------------
func (d *TreeMap) AsSliceOf(target []any) error {
	if d.err != nil {
		return d.err
	}
	data, err := json.Marshal(d.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

func (d *TreeMap) AsStruct(target any) error {
	if d.err != nil {
		return d.err
	}
	data, err := json.Marshal(d.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// ------------------- Default Fallbacks -------------------
func (d *TreeMap) AsStringOr(def string) string {
	s, err := d.AsString()
	if err != nil {
		return def
	}
	return s
}
func (d *TreeMap) AsIntOr(def int64) int64 {
	i, err := d.AsInt()
	if err != nil {
		return def
	}
	return i
}
func (d *TreeMap) AsFloatOr(def float64) float64 {
	f, err := d.AsFloat()
	if err != nil {
		return def
	}
	return f
}
func (d *TreeMap) AsBoolOr(def bool) bool {
	b, err := d.AsBool()
	if err != nil {
		return def
	}
	return b
}
func (d *TreeMap) AsAnyOr(def any) any {
	v, err := d.AsAny()
	if err != nil {
		return def
	}
	return v
}

func (d *TreeMap) AsMap() (DefaultMap, error) {
	if d.err != nil {
		return nil, d.err
	}
	if obj, ok := d.value.(DefaultMap); ok {
		return obj, nil
	}
	if m, ok := d.value.(map[string]any); ok {
		return m, nil
	}
	return nil, fmt.Errorf("cannot convert to map: %T", d.value)
}

func (d *TreeMap) AsSlice() ([]*TreeMap, error) {
	if d.err != nil {
		return nil, d.err
	}

	rv := reflect.ValueOf(d.value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("not a slice: %T", d.value)
	}

	var result []*TreeMap
	for i := range rv.Len() {
		result = append(result, &TreeMap{value: rv.Index(i).Interface(), root: d.root})
	}
	return result, nil
}

// ------------------- AsStrSlice / AsIntSlice / etc -------------------
func (d *TreeMap) AsStrSlice() []string {
	v, err := d.AsSlice()
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(v))
	for _, e := range v {
		s, _ := e.AsString()
		out = append(out, s)
	}
	return out
}

func (d *TreeMap) AsIntSlice() []int64 {
	v, err := d.AsSlice()
	if err != nil {
		return nil
	}
	out := make([]int64, 0, len(v))
	for _, e := range v {
		i, _ := e.AsInt()
		out = append(out, i)
	}
	return out
}

func (d *TreeMap) AsBoolSlice() []bool {
	v, err := d.AsSlice()
	if err != nil {
		return nil
	}
	out := make([]bool, 0, len(v))
	for _, e := range v {
		i, _ := e.AsBool()
		out = append(out, i)
	}
	return out
}

func (d *TreeMap) AsAnySlice() []any {
	v, err := d.AsSlice()
	if err != nil {
		return nil
	}
	out := make([]any, 0, len(v))
	for _, e := range v {
		out = append(out, e.AsAnyOr(nil))
	}
	return out
}
