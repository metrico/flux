package arrow

import (
	"github.com/InfluxCommunity/flux/array"
	"github.com/InfluxCommunity/flux/memory"
)

func NewString(vs []string, alloc memory.Allocator) *array.String {
	b := NewStringBuilder(alloc)
	b.Resize(len(vs))
	sz := 0
	for _, v := range vs {
		sz += len(v)
	}
	b.ReserveData(sz)
	for _, v := range vs {
		b.Append(v)
	}
	a := b.NewStringArray()
	b.Release()
	return a
}

func StringSlice(arr *array.String, i, j int) *array.String {
	return Slice(arr, int64(i), int64(j)).(*array.String)
}

func NewStringBuilder(a memory.Allocator) *array.StringBuilder {
	if a == nil {
		a = memory.DefaultAllocator
	}
	return array.NewStringBuilder(a)
}
