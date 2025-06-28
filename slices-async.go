package goutils

import "sync"

/* All functions are called async with waits for all to finish */
func MapAsync[T any, K any](collection []T, fn func(item T, index int) K) []K {
	result := make([]K, len(collection))

	var wg sync.WaitGroup
	wg.Add(len(collection))

	for i, item := range collection {
		go func(_item T, _i int) {
			res := fn(_item, _i)

			result[_i] = res

			wg.Done()
		}(item, i)
	}

	wg.Wait()

	return result
}

/* All functions are called async */
func ForEachAsync[T any](collection []T, fn func(item T, index int)) {
	for i, item := range collection {
		go func(_item T, _i int) { fn(_item, _i) }(item, i)
	}
}
