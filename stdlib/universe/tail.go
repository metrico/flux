package universe

import (
	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/arrow"
	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/internal/errors"
	"github.com/InfluxCommunity/flux/plan"
	"github.com/InfluxCommunity/flux/runtime"
)

const TailKind = "tail"

// TailOpSpec tails the number of rows returned per table.
type TailOpSpec struct {
	N      int64 `json:"n"`
	Offset int64 `json:"offset"`
}

func init() {
	tailSignature := runtime.MustLookupBuiltinType("universe", "tail")

	runtime.RegisterPackageValue("universe", TailKind, flux.MustValue(flux.FunctionValue(TailKind, createTailOpSpec, tailSignature)))
	plan.RegisterProcedureSpec(TailKind, newTailProcedure, TailKind)
	execute.RegisterTransformation(TailKind, createTailTransformation)
}

func createTailOpSpec(args flux.Arguments, a *flux.Administration) (flux.OperationSpec, error) {
	if err := a.AddParentFromArgs(args); err != nil {
		return nil, err
	}

	spec := new(TailOpSpec)

	n, err := args.GetRequiredInt("n")
	if err != nil {
		return nil, err
	}
	spec.N = n

	if offset, ok, err := args.GetInt("offset"); err != nil {
		return nil, err
	} else if ok {
		spec.Offset = offset
	}

	return spec, nil
}

func (s *TailOpSpec) Kind() flux.OperationKind {
	return TailKind
}

type TailProcedureSpec struct {
	plan.DefaultCost
	N      int64 `json:"n"`
	Offset int64 `json:"offset"`
}

func newTailProcedure(qs flux.OperationSpec, pa plan.Administration) (plan.ProcedureSpec, error) {
	spec, ok := qs.(*TailOpSpec)
	if !ok {
		return nil, errors.Newf(codes.Internal, "invalid spec type %T", qs)
	}
	return &TailProcedureSpec{
		N:      spec.N,
		Offset: spec.Offset,
	}, nil
}

func (s *TailProcedureSpec) Kind() plan.ProcedureKind {
	return TailKind
}
func (s *TailProcedureSpec) Copy() plan.ProcedureSpec {
	ns := new(TailProcedureSpec)
	*ns = *s
	return ns
}

// TriggerSpec implements plan.TriggerAwareProcedureSpec
func (s *TailProcedureSpec) TriggerSpec() plan.TriggerSpec {
	return plan.NarrowTransformationTriggerSpec{}
}

func createTailTransformation(id execute.DatasetID, mode execute.AccumulationMode, spec plan.ProcedureSpec, a execute.Administration) (execute.Transformation, execute.Dataset, error) {
	s, ok := spec.(*TailProcedureSpec)
	if !ok {
		return nil, nil, errors.Newf(codes.Internal, "invalid spec type %T", spec)
	}
	cache := execute.NewTableBuilderCache(a.Allocator())
	d := execute.NewDataset(id, mode, cache)
	t := NewTailTransformation(d, cache, s)
	return t, d, nil
}

type tailTransformation struct {
	execute.ExecutionNode
	d     execute.Dataset
	cache execute.TableBuilderCache

	n, offset int
}

func NewTailTransformation(d execute.Dataset, cache execute.TableBuilderCache, spec *TailProcedureSpec) *tailTransformation {
	return &tailTransformation{
		d:      d,
		cache:  cache,
		n:      int(spec.N),
		offset: int(spec.Offset),
	}
}

func (t *tailTransformation) RetractTable(id execute.DatasetID, key flux.GroupKey) error {
	return t.d.RetractTable(key)
}

func (t *tailTransformation) Process(id execute.DatasetID, tbl flux.Table) error {
	builder, created := t.cache.TableBuilder(tbl.Key())
	if !created {
		return errors.Newf(codes.FailedPrecondition, "tail found duplicate table with key: %v", tbl.Key())
	}
	if err := execute.AddTableCols(tbl, builder); err != nil {
		return err
	}

	n := t.n
	offset := t.offset
	readers := make([]flux.ColReader, 0)
	numRecords := 0

	var finished bool
	if err := tbl.Do(func(cr flux.ColReader) error {
		if n <= 0 {
			// Returning an error terminates iteration
			finished = true
			return errors.New(codes.Canceled)
		}

		cr.Retain()
		readers = append(readers, cr)
		numRecords += cr.Len()

		for numRecords-readers[0].Len() >= n+offset {
			numRecords -= readers[0].Len()
			readers[0].Release()
			readers = readers[1:]
		}

		return nil
	}); err != nil && !finished {
		return err
	}

	endIndex := numRecords
	offsetIndex := endIndex - offset
	startIndex := offsetIndex - n

	curr := 0
	for _, cr := range readers {
		var start, end int

		if startIndex > curr && startIndex < cr.Len() {
			start = startIndex
		} else {
			start = 0
		}

		if offsetIndex > curr && offsetIndex < curr+cr.Len() {
			end = offsetIndex - curr
		} else if offsetIndex <= curr {
			break
		} else {
			end = cr.Len()
		}

		if err := appendSlicedCols(cr, builder, start, end); err != nil {
			return err
		}

		curr += cr.Len()

		cr.Release()
	}

	return nil
}

func (t *tailTransformation) UpdateWatermark(id execute.DatasetID, mark execute.Time) error {
	return t.d.UpdateWatermark(mark)
}
func (t *tailTransformation) UpdateProcessingTime(id execute.DatasetID, pt execute.Time) error {
	return t.d.UpdateProcessingTime(pt)
}
func (t *tailTransformation) Finish(id execute.DatasetID, err error) {
	t.d.Finish(err)
}

func appendSlicedCols(reader flux.ColReader, builder execute.TableBuilder, start, stop int) error {
	for j, c := range reader.Cols() {
		if j > len(builder.Cols()) {
			return errors.New(codes.Internal, "builder index out of bounds")
		}

		switch c.Type {
		case flux.TBool:
			s := arrow.BoolSlice(reader.Bools(j), start, stop)
			if err := builder.AppendBools(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		case flux.TInt:
			s := arrow.IntSlice(reader.Ints(j), start, stop)
			if err := builder.AppendInts(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		case flux.TUInt:
			s := arrow.UintSlice(reader.UInts(j), start, stop)
			if err := builder.AppendUInts(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		case flux.TFloat:
			s := arrow.FloatSlice(reader.Floats(j), start, stop)
			if err := builder.AppendFloats(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		case flux.TString:
			s := arrow.StringSlice(reader.Strings(j), start, stop)
			if err := builder.AppendStrings(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		case flux.TTime:
			s := arrow.IntSlice(reader.Times(j), start, stop)
			if err := builder.AppendTimes(j, s); err != nil {
				s.Release()
				return err
			}
			s.Release()
		default:
			execute.PanicUnknownType(c.Type)
		}
	}

	return nil
}
