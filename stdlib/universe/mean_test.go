package universe_test

import (
	"testing"

	"github.com/InfluxCommunity/flux/array"
	"github.com/InfluxCommunity/flux/arrow"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/memory"
	"github.com/InfluxCommunity/flux/stdlib/universe"
)

func TestMean_Process(t *testing.T) {
	testCases := []struct {
		name string
		data func() *array.Float
		want interface{}
	}{
		{
			name: "zero",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{0, 0, 0}, nil)
			},
			want: 0.0,
		},
		{
			name: "nonzero",
			data: func() *array.Float {
				return arrow.NewFloat([]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}, nil)
			},
			want: 4.5,
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
				b.AppendValues([]float64{0, 1, 2, 3}, nil)
				b.AppendNull()
				b.AppendValues([]float64{5, 6}, nil)
				b.AppendNull()
				b.AppendValues([]float64{8, 9}, nil)
				return b.NewFloatArray()
			},
			want: 4.25,
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
			data := tc.data()
			defer data.Release()

			executetest.AggFuncTestHelper(
				t,
				new(universe.MeanAgg),
				data,
				tc.want,
			)
		})
	}
}

func BenchmarkMean(b *testing.B) {
	data := arrow.NewFloat(NormalData, &memory.ResourceAllocator{})
	executetest.AggFuncBenchmarkHelper(
		b,
		new(universe.MeanAgg),
		data,
		9.99847267384332,
	)
}
