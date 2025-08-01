package test

import (
	"fmt"
	"testing"

	goutils "github.com/nitsugaro/go-utils"
)

func TestMapTree(t *testing.T) {
	root := map[string]interface{}{
		"initial": "my-value",
	}

	mapTree := goutils.NewTreeMap(root)

	if value, err := mapTree.Get("initial").AsString(); err != nil || value != "my-value" {
		t.Errorf("expected 'initial' value be 'my-value' and got: %s", value)
	}

	fmt.Println(mapTree.Set("sub.key", goutils.DefaultMap{"slice": []interface{}{1, 2, 3}}))
	if !mapTree.IsDefined("sub.key") {
		t.Errorf("expected 'sub.key' be defined")
	}

	if mapSubKey, err := mapTree.Get("sub.key").AsMap(); err != nil || mapSubKey["slice"] == nil {
		t.Errorf("expected 'sub.key' be an object with a slice and got: %v", mapSubKey)
	}

	if item, err := mapTree.Get("sub.key.slice.0").AsInt(); err != nil || item != 1 {
		t.Errorf("expected 'sub.key.slice.0' be 1 and got: %v", item)
	}

	if items, err := mapTree.Get("sub.key.slice").AsSlice(); err != nil || len(items) != 3 {
		t.Errorf("expected len of 'sub.key.slice' be 3 and got: %v", items)
	}

	fmt.Println(mapTree.AsMap())

	mapTree.Delete("sub.key.slice")
	if mapTree.IsDefined("sub.key.slice") {
		t.Errorf("expected 'sub.key.slice' be removed from map tree")
	}

	mapTree.Set("sub.key.another_key", "what")

	fmt.Println(mapTree.ToJsonString(true))

	fmt.Println(mapTree.Get("unreference_key").Or("sub.key.another_key").AsString())
	fmt.Println(mapTree.Get("unreference_key").IsEmpty())

	fmt.Println(mapTree.Delete("sub.key").Get("another_key").AsString())

	fmt.Println(mapTree.ToJsonString(true))

	fmt.Println(mapTree.TryDelete("sub").ToJsonString(true))
}
