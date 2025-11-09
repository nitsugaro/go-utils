package goutils

import (
	"sync"
)

type SafeTreeMap struct {
	mu *sync.RWMutex
	tm *TreeMap
}

func NewSyncTreeMap(data ...any) *SafeTreeMap {
	mu := &sync.RWMutex{}
	return &SafeTreeMap{mu: mu, tm: NewTreeMap(data...)}
}

// ------------------- Core -------------------
func (s *SafeTreeMap) Get(path string) *SafeTreeMap {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SafeTreeMap{mu: s.mu, tm: s.tm.Get(path)}
}

func (s *SafeTreeMap) Set(path string, value any) *SafeTreeMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tm.Set(path, value)
	return s
}

func (s *SafeTreeMap) Delete(path string) *TreeMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tm.Delete(path)
}

func (s *SafeTreeMap) TryDelete(path string) *SafeTreeMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tm.TryDelete(path)
	return s
}

func (s *SafeTreeMap) Clone() *SafeTreeMap {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return &SafeTreeMap{mu: &sync.RWMutex{}, tm: s.tm.Clone()}
}

func (s *SafeTreeMap) ToJsonString(pretty bool) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.ToJsonString(pretty)
}

func (s *SafeTreeMap) AsMap() (DefaultMap, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsMap()
}

func (s *SafeTreeMap) AsSlice() ([]*TreeMap, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsSlice()
}

// ------------------- Value Conversions -------------------
func (s *SafeTreeMap) AsString() (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsString()
}

func (s *SafeTreeMap) AsInt() (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsInt()
}

func (s *SafeTreeMap) AsFloat() (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsFloat()
}

func (s *SafeTreeMap) AsBool() (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsBool()
}

func (s *SafeTreeMap) AsAny() (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsAny()
}

// ------------------- Struct / Slice Conversions -------------------
func (s *SafeTreeMap) AsSliceOf(target []any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsSliceOf(target)
}

func (s *SafeTreeMap) AsStruct(target any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsStruct(target)
}

// ------------------- Default Fallbacks -------------------
func (s *SafeTreeMap) AsStringOr(def string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsStringOr(def)
}

func (s *SafeTreeMap) AsIntOr(def int64) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsIntOr(def)
}

func (s *SafeTreeMap) AsFloatOr(def float64) float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsFloatOr(def)
}

func (s *SafeTreeMap) AsBoolOr(def bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsBoolOr(def)
}

func (s *SafeTreeMap) AsAnyOr(def any) any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsAnyOr(def)
}

// ------------------- Slice Helpers -------------------
func (s *SafeTreeMap) AsStrSlice() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsStrSlice()
}

func (s *SafeTreeMap) AsIntSlice() []int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsIntSlice()
}

func (s *SafeTreeMap) AsBoolSlice() []bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsBoolSlice()
}

func (s *SafeTreeMap) AsAnySlice() []any {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.tm.AsAnySlice()
}
