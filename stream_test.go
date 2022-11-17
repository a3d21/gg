package gg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBaseStream(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	want := 5
	got := Reduce(0, func(acc, it int) int {
		return acc + it
	}, Map(func(it int) int {
		return it * it
	}, Filter(func(it int) bool {
		return it < 3
	}, FromSlice(src))))

	assert.Equal(t, want, got)
}

func TestStreamToSLice(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}

	want := []int{9, 16, 25}
	got := ToSlice(
		Filter(func(it int) bool {
			return it > 5
		}, Map(func(it int) int {
			return it * it
		}, FromSlice(src))))

	assert.Equal(t, want, got)
}

func TestStreamToMap(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	want := map[string]int{"1": 1, "2": 2}
	got := ToMap(
		Map(func(it int) KeyValue[string, int] {
			return KeyValue[string, int]{
				Key:   fmt.Sprintf("%v", it),
				Value: it,
			}
		}, Filter(func(it int) bool {
			return it < 3
		}, FromSlice(src))))

	assert.Equal(t, want, got)
}

func TestFromMapToMap(t *testing.T) {
	src := map[string]int{"foo": 1, "bar": 2}

	want := map[int]string{1: "foo", 2: "bar"}
	got := ToMap(
		Map(func(kv KeyValue[string, int]) KeyValue[int, string] {
			return KeyValue[int, string]{
				Key:   kv.Value,
				Value: kv.Key,
			}
		}, FromMap(src)))

	assert.Equal(t, want, got)
}
