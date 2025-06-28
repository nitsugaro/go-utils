package goutils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type DefaultMap map[string]interface{}

type TreeMap struct {
	value interface{}
	err   error
}

func NewTreeMap(data ...DefaultMap) *TreeMap {
	if len(data) != 0 && data[0] != nil {
		return &TreeMap{value: data[0]}
	} else {
		return &TreeMap{value: DefaultMap{}}
	}
}

func (d *TreeMap) Get(path string) *TreeMap {
	if d.err != nil || d.value == nil {
		return d
	}

	parts := strings.Split(path, ".")
	current := d.value

	for _, part := range parts {
		switch val := current.(type) {
		case DefaultMap:
			current = val[part]
		default:
			rv := reflect.ValueOf(current)
			if rv.Kind() == reflect.Slice {
				i, err := strconv.Atoi(part)
				if err != nil || i < 0 || i >= rv.Len() {
					return &TreeMap{err: fmt.Errorf("index out of bounds: %s", part)}
				}
				current = rv.Index(i).Interface()
			} else {
				return &TreeMap{err: fmt.Errorf("invalid access at %s", part)}
			}
		}
	}

	return &TreeMap{value: current}
}

func (d *TreeMap) Set(path string, value interface{}) *TreeMap {
	if d.err != nil {
		return d
	}

	root, ok := d.value.(DefaultMap)
	if !ok {
		return &TreeMap{err: fmt.Errorf("root is not a map")}
	}

	parts := strings.Split(path, ".")
	last := len(parts) - 1
	current := root

	for i, part := range parts {
		if i == last {
			current[part] = value
			return d
		}

		next, ok := current[part]
		if !ok {
			newMap := make(DefaultMap)
			current[part] = newMap
			current = newMap
			continue
		}

		nestedMap, ok := next.(DefaultMap)
		if !ok {
			return &TreeMap{err: fmt.Errorf("intermediate value at %s is not a map", part)}
		}

		current = nestedMap
	}

	return d
}

func (d *TreeMap) Delete(path string) *TreeMap {
	if d.err != nil {
		return d
	}

	root, ok := d.value.(DefaultMap)
	if !ok {
		return &TreeMap{err: fmt.Errorf("root is not a map")}
	}

	parts := strings.Split(path, ".")
	last := len(parts) - 1
	current := root

	for i, part := range parts {
		if i == last {
			delete(current, part)
			return d
		}

		next, ok := current[part]
		if !ok {
			return d
		}

		nestedMap, ok := next.(DefaultMap)
		if !ok {
			return &TreeMap{err: fmt.Errorf("intermediate path '%s' is not a map", part)}
		}
		current = nestedMap
	}

	return d
}

func (d *TreeMap) IsDefined(path string) bool {
	v := d.Get(path)
	return v.err == nil && v.value != nil
}

func (d *TreeMap) AsString() (string, error) {
	if d.err != nil {
		return "", d.err
	}
	switch v := d.value.(type) {
	case string:
		return v, nil
	case fmt.Stringer:
		return v.String(), nil
	case float64:
		return fmt.Sprintf("%v", v), nil
	case int:
		return fmt.Sprintf("%v", v), nil
	case bool:
		return strconv.FormatBool(v), nil
	default:
		return "", fmt.Errorf("cannot convert to string: %T", v)
	}
}

func (d *TreeMap) AsInt() (int, error) {
	if d.err != nil {
		return 0, d.err
	}
	switch v := d.value.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
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

func (d *TreeMap) AsSlice() ([]*TreeMap, error) {
	if d.err != nil {
		return nil, d.err
	}
	arr, ok := d.value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a slice")
	}
	var result []*TreeMap
	for _, v := range arr {
		result = append(result, &TreeMap{value: v})
	}
	return result, nil
}

func (d *TreeMap) AsMap() (DefaultMap, error) {
	if d.err != nil {
		return nil, d.err
	}
	switch v := d.value.(type) {
	case DefaultMap:
		return d.value.(DefaultMap), nil
	default:
		return nil, fmt.Errorf("cannot convert to map: %T", v)
	}
}

func (d *TreeMap) AsSliceOf(target interface{}) error {
	if d.err != nil {
		return d.err
	}

	data, err := json.Marshal(d.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (d *TreeMap) AsStruct(target interface{}) error {
	if d.err != nil {
		return d.err
	}
	bytes, err := json.Marshal(d.value)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func (d *TreeMap) ToJsonString(pretty bool) (string, error) {
	if d.err != nil {
		return "", d.err
	}
	if pretty {
		bytes, err := json.MarshalIndent(d.value, "", "  ")
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	} else {
		bytes, err := json.Marshal(d.value)
		if err != nil {
			return "", err
		}
		return string(bytes), nil
	}
}
