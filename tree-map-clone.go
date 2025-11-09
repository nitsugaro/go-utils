package goutils

import "reflect"

func (d *TreeMap) Clone() TreeMapImpl {
	if d.err != nil {
		return &TreeMap{err: d.err}
	}
	return &TreeMap{
		value: deepClone(d.value),
		root:  d.root,
	}
}

func deepClone(v any) any {
	if v == nil {
		return nil
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		newMap := reflect.MakeMap(rv.Type())
		for _, key := range rv.MapKeys() {
			newMap.SetMapIndex(key, reflect.ValueOf(deepClone(rv.MapIndex(key).Interface())))
		}
		return newMap.Interface()
	case reflect.Slice:
		newSlice := reflect.MakeSlice(rv.Type(), rv.Len(), rv.Len())
		for i := 0; i < rv.Len(); i++ {
			newSlice.Index(i).Set(reflect.ValueOf(deepClone(rv.Index(i).Interface())))
		}
		return newSlice.Interface()
	case reflect.Array:
		newArray := reflect.New(rv.Type()).Elem()
		for i := 0; i < rv.Len(); i++ {
			newArray.Index(i).Set(reflect.ValueOf(deepClone(rv.Index(i).Interface())))
		}
		return newArray.Interface()
	case reflect.Ptr:
		val := reflect.New(rv.Elem().Type())
		val.Elem().Set(reflect.ValueOf(deepClone(rv.Elem().Interface())))
		return val.Interface()
	default:
		return v
	}
}
