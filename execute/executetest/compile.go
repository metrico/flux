package executetest

import (
	"context"
	"fmt"
	"testing"

	"github.com/InfluxCommunity/flux/interpreter"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/semantic"
	"github.com/InfluxCommunity/flux/values"
)

var (
	prelude values.Scope
	stdlib  interpreter.Importer
)

// FunctionExpression will take a function expression as a string
// and return the *semantic.FunctionExpression.
//
// This will cause a fatal error in the test on failure.
func FunctionExpression(t testing.TB, source string, args ...interface{}) *semantic.FunctionExpression {
	t.Helper()

	if len(args) > 0 {
		source = fmt.Sprintf(source, args...)
	}

	if stdlib == nil {
		stdlib = runtime.StdLib()
	}
	if prelude == nil {
		prelude = values.NewScope()
		for _, path := range runtime.PreludeList {
			p, err := stdlib.ImportPackageObject(path)
			if err != nil {
				t.Fatalf("error importing prelude package %q: %s", path, err)
			}
			p.Range(prelude.Set)
		}
	}

	ctx := context.Background()
	pkg, err := runtime.AnalyzeSource(ctx, source)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	// Interpret and resolve the function which will replace
	// variables with their values (notably identifiers "true"
	// and "false" will be replaced with boolean literals)
	itrp := interpreter.NewInterpreter(nil, nil)
	se, err := itrp.Eval(ctx, pkg, prelude, stdlib)
	if err != nil {
		t.Fatal(err)
	}

	if len(se) != 1 {
		t.Fatal("expected just one side effect")
	}

	f := se[0].Value.(values.Function)
	rf, err := interpreter.ResolveFunction(f)
	if err != nil {
		t.Fatal(err)
	}

	return rf.Fn
}
