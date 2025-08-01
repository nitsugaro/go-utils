package goutils

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ------------------- Get / Set / Delete -------------------
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

	return &TreeMap{value: current, root: d.root}
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

func (d *TreeMap) Set(path string, value any) *TreeMap {
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
