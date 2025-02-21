package testing_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/execute/executetest"
	"github.com/InfluxCommunity/flux/plan"
	fluxtesting "github.com/InfluxCommunity/flux/stdlib/testing"
	"github.com/google/go-cmp/cmp"
)

func TestAssertEquals_Process(t *testing.T) {
	testCases := []struct {
		name    string
		spec    *fluxtesting.AssertEqualsProcedureSpec
		data0   []*executetest.Table // data from parent 0
		data1   []*executetest.Table // data from parent 1
		want    []*executetest.Table
		wantErr bool
	}{
		{
			name: "simple equality",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			want: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
		},
		{
			name: "left empty",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "right empty",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "left short",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "right short",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "multiple Tables",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 100.0},
						{execute.Time(2), 200.0},
						{execute.Time(3), 300.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 100.0},
						{execute.Time(2), 200.0},
						{execute.Time(3), 300.0},
					},
				},
			},
			want: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
		},
		{
			name: "left short one table",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 100.0},
						{execute.Time(2), 200.0},
						{execute.Time(3), 300.0},
					},
				},
			},
			want: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
		},
		{
			name: "right short one table",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 100.0},
						{execute.Time(2), 200.0},
						{execute.Time(3), 300.0},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 10.0},
						{execute.Time(2), 20.0},
						{execute.Time(3), 30.0},
					},
				},
			},
			want: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), 1.0},
						{execute.Time(2), 2.0},
						{execute.Time(3), 3.0},
					},
				},
			},
		},
		{
			name: "null equality",
			spec: &fluxtesting.AssertEqualsProcedureSpec{
				DefaultCost: plan.DefaultCost{},
				Name:        "simple",
			},
			data0: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), nil},
					},
				},
			},
			data1: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), nil},
					},
				},
			},
			want: []*executetest.Table{
				{
					ColMeta: []flux.ColMeta{
						{Label: "_time", Type: flux.TTime},
						{Label: "_value", Type: flux.TFloat},
					},
					Data: [][]interface{}{
						{execute.Time(1), nil},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			id0 := executetest.RandomDatasetID()
			id1 := executetest.RandomDatasetID()

			parents := []execute.DatasetID{
				execute.DatasetID(id0),
				execute.DatasetID(id1),
			}

			d := executetest.NewDataset(executetest.RandomDatasetID())
			c := execute.NewTableBuilderCache(executetest.UnlimitedAllocator)
			c.SetTriggerSpec(plan.DefaultTriggerSpec)
			jt := fluxtesting.NewAssertEqualsTransformation(d, c, tc.spec, parents[0], parents[1], executetest.UnlimitedAllocator)

			executetest.NormalizeTables(tc.data0)
			executetest.NormalizeTables(tc.data1)
			l := len(tc.data0)
			if len(tc.data1) > l {
				l = len(tc.data1)
			}
			var err error
			for i := 0; i < l; i++ {
				if i < len(tc.data0) {
					if err = jt.Process(parents[0], tc.data0[i]); err != nil {
						break
					}
				}
				if i < len(tc.data1) {
					if err = jt.Process(parents[1], tc.data1[i]); err != nil {
						break
					}
				}
			}

			if tc.wantErr {
				if err == nil {
					t.Fatal(fmt.Errorf("case %s expected an error, got none", tc.name))
				} else {
					return
				}

			}

			got, err := executetest.TablesFromCache(c)
			if err != nil {
				t.Fatal(err)
			}

			executetest.NormalizeTables(got)
			executetest.NormalizeTables(tc.want)

			sort.Sort(executetest.SortedTables(got))
			sort.Sort(executetest.SortedTables(tc.want))

			if !cmp.Equal(tc.want, got) {
				t.Errorf("unexpected tables -want/+got\n%s", cmp.Diff(tc.want, got))
			}
		})
	}
}
