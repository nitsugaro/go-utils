package goutils

type TreeMapImpl interface {
	Get(path string) TreeMapImpl
	IsDefined(path string) bool
	Or(path string) TreeMapImpl
	Set(path string, value any) TreeMapImpl
	Delete(path string) TreeMapImpl
	TryDelete(path string) TreeMapImpl
	Clone() TreeMapImpl
	ToJsonString(pretty bool) string
	AsMap() (DefaultMap, error)
	AsSlice() ([]TreeMapImpl, error)
	AsString() (string, error)
	AsInt() (int64, error)
	AsFloat() (float64, error)
	AsBool() (bool, error)
	AsAny() (any, error)
	AsSliceOf(target []any) error
	AsStruct(target any) error
	AsStringOr(def string) string
	AsIntOr(def int64) int64
	AsFloatOr(def float64) float64
	AsBoolOr(def bool) bool
	AsAnyOr(def any) any
	AsStrSlice() []string
	AsIntSlice() []int64
	AsBoolSlice() []bool
	AsAnySlice() []any

	getValue() any
	getError() error
	getRoot() TreeMapImpl
}
