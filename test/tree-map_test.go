package test

import (
	"fmt"
	"sync"
	"testing"

	goutils "github.com/nitsugaro/go-utils"
)

// ------------------- TreeMap Core Tests -------------------

func TestTreeMap_BasicUsage(t *testing.T) {
	root := map[string]interface{}{
		"initial": "my-value",
	}
	mapTree := goutils.NewTreeMap(root)

	mapTree.Set("array", []any{"1", "1"})
	if got := mapTree.Get("array.0").AsStringOr(""); got != "1" {
		t.Errorf("expected array[0]='1', got %v", got)
	}

	copy := mapTree.Clone()
	copy.Set("obj", map[string]string{"x": "2"})
	if mapTree.IsDefined("obj.x") {
		t.Errorf("clone must not affect original map")
	}

	if v, _ := copy.Get("obj.x").AsString(); v != "2" {
		t.Errorf("expected obj.x='2', got %v", v)
	}

	mapTree.Set("sub.key.slice", []any{1, 2, 3})
	if !mapTree.IsDefined("sub.key.slice.2") {
		t.Errorf("expected sub.key.slice.2 to exist")
	}
	mapTree.Delete("sub.key.slice")
	if mapTree.IsDefined("sub.key.slice") {
		t.Errorf("expected sub.key.slice to be deleted")
	}

	mapTree.Set("sub.key.another_key", "what")
	v := mapTree.Get("unreference_key").Or("sub.key.another_key").AsStringOr("")
	if v != "what" {
		t.Errorf("expected Or fallback to return 'what', got %v", v)
	}

	mapTree.TryDelete("sub")
	if mapTree.IsDefined("sub") {
		t.Errorf("expected 'sub' to be deleted")
	}
}

// ------------------- Value Conversion Tests -------------------

func TestTreeMap_Conversions(t *testing.T) {
	m := goutils.NewTreeMap()
	m.Set("nums", []any{1, "2", 3.0})
	m.Set("bools", []any{"true", false})
	m.Set("str", "10")

	// slices
	if ints := m.Get("nums").AsIntSlice(); len(ints) != 3 || ints[1] != 2 {
		t.Errorf("expected ints [1 2 3], got %v", ints)
	}

	if bools := m.Get("bools").AsBoolSlice(); len(bools) != 2 || !bools[0] {
		t.Errorf("expected bools [true false], got %v", bools)
	}

	if s, _ := m.Get("str").AsInt(); s != 10 {
		t.Errorf("expected str parsed as 10, got %v", s)
	}
}

// ------------------- Clone & DeepClone Tests -------------------

func TestTreeMap_CloneIsolation(t *testing.T) {
	m := goutils.NewTreeMap(map[string]any{"x": []any{1, 2}})
	clone := m.Clone()
	clone.Set("x.1", 99)

	if m.Get("x.1").AsIntOr(0) == 99 {
		t.Errorf("clone modifications must not affect original")
	}
}

// ------------------- SafeTreeMap Race & Stress Tests -------------------

func TestSafeTreeMapRace(t *testing.T) {
	safe := goutils.NewSyncTreeMap()
	safe.Set("user", map[string]any{"name": "Agus", "age": 27})

	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(id int) {
			if id%2 == 0 {
				safe.Get("user").Set(fmt.Sprintf("role-%d", id), id)
			} else {
				_ = safe.Get("user.name").AsStringOr("")
			}
			done <- true
		}(i)
	}
	for i := 0; i < 100; i++ {
		<-done
	}
}

// ------------------- SafeTreeMap Extended Stress -------------------

func TestSafeTreeMap_ExtendedStress(t *testing.T) {
	safe := goutils.NewSyncTreeMap()
	safe.Set("root", map[string]any{"list": []any{1, 2, 3}, "active": true})

	wg := sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%4 == 0 {
				safe.Get("root").Set(fmt.Sprintf("k%d", id), id)
			} else if id%4 == 1 {
				_ = safe.Get("root.list").AsIntSlice()
			} else if id%4 == 2 {
				_ = safe.Get("root.active").AsBoolOr(false)
			} else {
				_ = safe.Clone().ToJsonString(false)
			}
		}(i)
	}
	wg.Wait()

	if !safe.Get("root.active").AsBoolOr(false) {
		t.Errorf("expected root.active to remain true after concurrent ops")
	}
}

// ------------------- SafeTreeMap Snapshot Clone -------------------

func TestSafeTreeMap_CloneIsolation(t *testing.T) {
	safe := goutils.NewSyncTreeMap()
	safe.Set("ctx", map[string]any{"x": 1, "y": 2})

	clone := safe.Clone()
	clone.Set("ctx.x", 999)

	if safe.Get("ctx.x").AsIntOr(0) == 999 {
		t.Errorf("clone modifications must not affect original SafeTreeMap")
	}
}
