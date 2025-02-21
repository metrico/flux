package universe_test

import (
	"math"
	"testing"

	"github.com/InfluxCommunity/flux/array"
	"github.com/InfluxCommunity/flux/arrow"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/memory"
	"github.com/InfluxCommunity/flux/stdlib/universe"
)

func TestSkew_Process(t *testing.T) {
	testCases := []struct {
		name string
		data func() *array.Float
		want interface{}
	}{
		{
			name: "zero",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{1, 2, 3}, nil)
			},
			want: 0.0,
		},
		{
			name: "nonzero",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{2, 2, 3}, nil)
			},
			want: 0.7071067811865475,
		},
		{
			name: "nonzero 2",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{2, 2, 3, 4}, nil)
			},
			want: 0.49338220021815854,
		},
		{
			name: "NaN short",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{1}, nil)
			},
			want: math.NaN(),
		},
		{
			name: "NaN divide by zero",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{1, 1, 1}, nil)
			},
			want: math.NaN(),
		},
		{
			name: "empty",
			data: func() *array.Float {
				return arrow.NewFloat(nil, nil)
			},
			want: nil,
		},
		{
			name: "with nulls",
			data: func() *array.Float {
				b := arrow.NewFloatBuilder(nil)
				defer b.Release()
				b.Append(2)
				b.AppendNull()
				b.Append(2)
				b.AppendNull()
				b.Append(3)
				return b.NewFloatArray()
			},
			want: 0.7071067811865475,
		},
		{
			name: "only nulls",
			data: func() *array.Float {
				b := arrow.NewFloatBuilder(nil)
				defer b.Release()
				b.AppendNull()
				b.AppendNull()
				return b.NewFloatArray()
			},
			want: nil,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			executetest.AggFuncTestHelper(
				t,
				new(universe.SkewAgg),
				tc.data(),
				tc.want,
			)
		})
	}
}

func BenchmarkSkew(b *testing.B) {
	data := arrow.NewFloat(NormalData, &memory.ResourceAllocator{})
	executetest.AggFuncBenchmarkHelper(
		b,
		new(universe.SkewAgg),
		data,
		-0.0019606823191321435,
	)
}
