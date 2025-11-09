package goutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ------------------- Get / Set / Delete -------------------

// ------------------- Get -------------------
func (d *TreeMap) Get(path string) *TreeMap {
	if d.err != nil || d.value == nil {
		return d
	}
	parts := strings.Split(path, ".")
	current := d.value

	for _, part := range parts {
		switch v := current.(type) {
		case map[string]any:
			current = v[part]

		case []any:
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 0 || idx >= len(v) {
				return &TreeMap{err: fmt.Errorf("index out of bounds: %s", part)}
			}
			current = v[idx]

		default:
			// Fallback: si todavía no está normalizado, probamos una vez
			rv := reflect.ValueOf(current)
			switch rv.Kind() {
			case reflect.Map:
				if rv.Type().Key().Kind() != reflect.String {
					return &TreeMap{err: fmt.Errorf("map key is not string at %s", part)}
				}
				val := rv.MapIndex(reflect.ValueOf(part))
				if !val.IsValid() {
					return &TreeMap{err: fmt.Errorf("key not found: %s", part)}
				}
				current = val.Interface()

			case reflect.Slice, reflect.Array:
				i, err := strconv.Atoi(part)
				if err != nil || i < 0 || i >= rv.Len() {
					return &TreeMap{err: fmt.Errorf("index out of bounds: %s", part)}
				}
				current = rv.Index(i).Interface()

			default:
				return &TreeMap{err: fmt.Errorf("invalid access at %s", part)}
			}
		}
	}
	return &TreeMap{value: current, root: d.root}
}

// ------------------- Set -------------------
func (d *TreeMap) Set(path string, value any) *TreeMap {
	if d.err != nil {
		return d
	}

	root, ok := d.value.(map[string]any)
	if !ok {
		return &TreeMap{err: fmt.Errorf("root is not a map")}
	}

	parts := strings.Split(path, ".")
	last := len(parts) - 1
	curr := root

	for i, part := range parts {
		if i == last {
			curr[part] = normalizeToDefault(value)
			return d
		}
		next, ok := curr[part]
		if !ok || next == nil {
			nm := make(map[string]any)
			curr[part] = nm
			curr = nm
			continue
		}

		if mm, ok := next.(map[string]any); ok {
			curr = mm
			continue
		}

		nn := normalizeToDefault(next)
		mm, ok := nn.(map[string]any)
		if !ok {
			return &TreeMap{err: fmt.Errorf("intermediate value at %s is not a map", part)}
		}
		curr[part] = mm
		curr = mm
	}
	return d
}

func (d *TreeMap) Or(path string) *TreeMap {
	if d.err != nil || d.value == nil {
		if d.root != nil {
			return d.root.Get(path)
		}
		return d
	}
	return d
}

// delete path key and returns a new treemap with from path value
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
			val := current[part]
			delete(current, part)
			return &TreeMap{value: val}
		}
		next, ok := current[part]
		if !ok {
			return &TreeMap{err: fmt.Errorf("not found key '%s'", part)}
		}
		nestedMap, ok := next.(DefaultMap)
		if !ok {
			return &TreeMap{err: fmt.Errorf("intermediate path '%s' is not a map", part)}
		}
		current = nestedMap
	}
	return &TreeMap{err: fmt.Errorf("cannot delete path '%s'", path)}
}

// delete path key and returns the root treemap
func (d *TreeMap) TryDelete(path string) *TreeMap {
	if d.err != nil {
		return d
	}

	d.Delete(path)

	return d
}

// ------------------- Status -------------------
func (d *TreeMap) IsDefined(path string) bool {
	v := d.Get(path)
	return v.err == nil && v.value != nil
}

func (d *TreeMap) Exists() bool {
	return d.err == nil
}

func (d *TreeMap) IsEmpty() bool {
	return d.err == nil && d.value == nil
}

func normalizeToDefault(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			return v
		}
		out := make(map[string]any, rv.Len())
		for _, k := range rv.MapKeys() {
			out[k.String()] = normalizeToDefault(rv.MapIndex(k).Interface())
		}
		return out
	case reflect.Slice, reflect.Array:
		n := rv.Len()
		out := make([]any, n)
		for i := 0; i < n; i++ {
			out[i] = normalizeToDefault(rv.Index(i).Interface())
		}
		return out
	default:
		return v
	}
}
