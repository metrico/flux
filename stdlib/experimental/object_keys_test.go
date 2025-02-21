package experimental_test

import (
	"context"
	"testing"

	"github.com/InfluxCommunity/flux/dependencies/dependenciestest"
	"github.com/InfluxCommunity/flux/dependency"
	"github.com/InfluxCommunity/flux/runtime"
)

func TestObjectKeys(t *testing.T) {
	script := `
import "experimental"
import "internal/testutil"

o = {a: 1, b: 2, c: 3}
experimental.objectKeys(o: o) == ["a", "b", "c"] or testutil.fail()
`
	ctx, deps := dependency.Inject(context.Background(), dependenciestest.Default())
	defer deps.Finish()
	if _, _, err := runtime.Eval(ctx, script); err != nil {
		t.Fatal("evaluation of objectKeys failed: ", err)
	}
}
