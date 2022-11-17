package gg

type Iterator[T any] func() (item T, ok bool)
type Fn[T, R any] func(T) R
type Accumulator[U, T any] func(U, T) U

type KeyValue[K comparable, V any] struct {
	Key   K
	Value V
}

type Stream[T any] struct {
	Iterate func() Iterator[T]
}

func FromSlice[T any](src []T) Stream[T] {
	srcLen := len(src)
	return Stream[T]{
		Iterate: func() Iterator[T] {
			index := 0

			return func() (item T, ok bool) {
				ok = index < srcLen
				if ok {
					item = src[index]
					index++
				}
				return
			}
		},
	}
}

func FromMap[K comparable, V any](src map[K]V) Stream[KeyValue[K, V]] {
	return Stream[KeyValue[K, V]]{
		Iterate: func() Iterator[KeyValue[K, V]] {
			var keys []K
			for k, _ := range src {
				keys = append(keys, k)
			}
			srcLen := len(keys)

			index := 0

			return func() (item KeyValue[K, V], ok bool) {
				ok = index < srcLen
				if ok {
					item = KeyValue[K, V]{
						Key:   keys[index],
						Value: src[keys[index]],
					}
					index++
				}
				return
			}
		},
	}
}

func Map[T, R any](fn Fn[T, R], s Stream[T]) Stream[R] {
	return Stream[R]{
		Iterate: func() Iterator[R] {
			next := s.Iterate()

			return func() (item R, ok bool) {
				it, ok := next()
				if ok {
					item = fn(it)
				}

				return
			}
		},
	}
}

func Filter[T any](predicate func(T) bool, s Stream[T]) Stream[T] {
	return Stream[T]{
		Iterate: func() Iterator[T] {
			next := s.Iterate()

			return func() (item T, ok bool) {
				for item, ok = next(); ok; item, ok = next() {
					if predicate(item) {
						return
					}
				}

				return
			}
		},
	}
}

func Reduce[U, T any](seed U, accumulator Accumulator[U, T], s Stream[T]) U {
	next := s.Iterate()
	result := seed

	for current, ok := next(); ok; current, ok = next() {
		result = accumulator(result, current)
	}

	return result
}

func (s Stream[T]) ForEach(action func(int, T)) {
	next := s.Iterate()
	index := 0

	for item, ok := next(); ok; item, ok = next() {
		action(index, item)
		index++
	}
}

func ToSlice[T any](s Stream[T]) []T {
	var result []T
	s.ForEach(func(_ int, it T) {
		result = append(result, it)
	})
	return result
}

func ToMap[K comparable, V any](s Stream[KeyValue[K, V]]) map[K]V {
	result := map[K]V{}
	s.ForEach(func(_ int, kv KeyValue[K, V]) {
		result[kv.Key] = kv.Value
	})
	return result
}
