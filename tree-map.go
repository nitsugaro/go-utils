package goutils

import (
	"encoding/json"
)

type DefaultMap = map[string]any

type TreeMap struct {
	value any
	root  *TreeMap
	err   error
}

// ------------------- Constructors -------------------
func NewTreeMap(data ...any) *TreeMap {
	var val any = make(map[string]any)

	if len(data) > 0 && data[0] != nil {
		val = normalizeToDefault(data[0])
	}

	root := &TreeMap{value: val}
	root.root = root
	return root
}

func (d *TreeMap) ToJsonString(pretty bool) string {
	if d.err != nil {
		return "{}"
	}
	var (
		bytes []byte
		err   error
	)
	if pretty {
		bytes, err = json.MarshalIndent(d.value, "", "  ")
	} else {
		bytes, err = json.Marshal(d.value)
	}
	if err != nil {
		d.err = err
		return "{}"
	}
	return string(bytes)
}
