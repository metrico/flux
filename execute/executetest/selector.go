package executetest

import (
	"testing"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/google/go-cmp/cmp"
)

func RowSelectorFuncTestHelper(t *testing.T, selector execute.RowSelector, data flux.Table, want []execute.Row) {
	t.Helper()

	s := selector.NewFloatSelector()
	valueIdx := execute.ColIdx(execute.DefaultValueColLabel, data.Cols())
	if valueIdx < 0 {
		t.Fatal("no _value column found")
	}
	if err := data.Do(func(cr flux.ColReader) error {
		s.DoFloat(cr.Floats(valueIdx), cr)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	got := s.Rows()

	if !cmp.Equal(want, got) {
		t.Errorf("unexpected value -want/+got\n%s", cmp.Diff(want, got))
	}
}

//lint:ignore U1000 Not sure why we need this...someone write a better reason :) .
var rows []execute.Row

func RowSelectorFuncBenchmarkHelper(b *testing.B, selector execute.RowSelector, data flux.BufferedTable) {
	b.Helper()

	valueIdx := execute.ColIdx(execute.DefaultValueColLabel, data.Cols())
	if valueIdx < 0 {
		b.Fatal("no _value column found")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s := selector.NewFloatSelector()

		t := data.Copy()
		if err := t.Do(func(cr flux.ColReader) error {
			s.DoFloat(cr.Floats(valueIdx), cr)
			return nil
		}); err != nil {
			b.Fatal(err)
		}
		rows = s.Rows()
	}
}

func IndexSelectorFuncTestHelper(t *testing.T, selector execute.IndexSelector, data flux.Table, want [][]int) {
	t.Helper()

	var got [][]int
	s := selector.NewFloatSelector()
	valueIdx := execute.ColIdx(execute.DefaultValueColLabel, data.Cols())
	if valueIdx < 0 {
		t.Fatal("no _value column found")
	}
	if err := data.Do(func(cr flux.ColReader) error {
		var cpy []int
		selected := s.DoFloat(cr.Floats(valueIdx))
		if len(selected) > 0 {
			cpy = make([]int, len(selected))
			copy(cpy, selected)
		}
		got = append(got, cpy)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Errorf("unexpected value -want/+got\n%s", cmp.Diff(want, got))
	}
}

func IndexSelectorFuncBenchmarkHelper(b *testing.B, selector execute.IndexSelector, data flux.BufferedTable) {
	b.Helper()

	valueIdx := execute.ColIdx(execute.DefaultValueColLabel, data.Cols())
	if valueIdx < 0 {
		b.Fatal("no _value column found")
	}

	b.ResetTimer()
	var got [][]int
	for n := 0; n < b.N; n++ {
		s := selector.NewFloatSelector()

		t := data.Copy()
		if err := t.Do(func(cr flux.ColReader) error {
			got = append(got, s.DoFloat(cr.Floats(valueIdx)))
			return nil
		}); err != nil {
			b.Fatal(err)
		}
	}
}
