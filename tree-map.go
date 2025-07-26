package goutils

import (
	"encoding/json"
)

type DefaultMap = map[string]interface{}

// ------------------- Core Struct -------------------
type TreeMap struct {
	value interface{}
	err   error
}

// ------------------- Constructors -------------------
func NewTreeMap(data ...interface{}) *TreeMap {
	if len(data) != 0 && data[0] != nil {
		return &TreeMap{value: data[0]}
	}
	return &TreeMap{value: DefaultMap{}}
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
