package experimental_test

import (
	"context"
	"testing"

	"github.com/InfluxCommunity/flux/dependencies/dependenciestest"
	"github.com/InfluxCommunity/flux/dependency"
	"github.com/InfluxCommunity/flux/execute"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/stdlib/experimental"
	"github.com/InfluxCommunity/flux/values"
)

var table1 = `
import "csv"

data = "#datatype,string,long,dateTime:RFC3339,double,string,string
#group,false,false,false,false,true,true
#default,_result,,,,,
,result,table,_time,_value,_measurement,user
,,0,2018-05-22T19:53:26Z,0,CPU,user1
,,0,2018-05-22T19:53:36Z,1,CPU,user1
,,1,2018-05-22T19:53:26Z,4,CPU,user2
,,1,2018-05-22T19:53:36Z,20,CPU,user2
,,1,2018-05-22T19:53:46Z,7,CPU,user2
,,2,2018-05-22T19:53:26Z,1,RAM,user1
"

inj = csv.from(csv: data)

`

var table2 = `
import "csv"

data = "#datatype,string,long,dateTime:RFC3339,double,string,string
#group,false,false,false,false,true,true
#default,_result,,,,,
,result,table,_time,_value,_measurement,user
,,0,2018-05-22T19:53:26Z,0,RAM,user1
,,0,2018-05-22T19:53:36Z,1,RAM,user1
,,1,2018-05-22T19:53:26Z,4,RAM,user2
,,1,2018-05-22T19:53:36Z,20,RAM,user2
,,1,2018-05-22T19:53:46Z,7,RAM,user2
,,2,2018-05-22T19:53:26Z,1,CPU,user1
"

inj = csv.from(csv: data)

`

func makeArgs(first values.Value, second values.Value) values.Object {
	argMap := map[string]values.Value{
		"first":  first,
		"second": second,
	}
	args := values.NewObjectWithValues(argMap)
	return args
}

func TestChain(t *testing.T) {
	ctx, deps := dependency.Inject(
		context.Background(),
		dependenciestest.Default(),
		execute.DefaultExecutionDependencies(),
	)
	defer deps.Finish()

	_, scope, err := runtime.Eval(ctx, table1)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	table1, ok := scope.Lookup("inj")
	if !ok {
		t.Fatal("unable to find input in table1 script")
	}

	_, scope, err = runtime.Eval(ctx, table2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	table2, ok := scope.Lookup("inj")
	if !ok {
		t.Fatal("unable to find input in table1 script")
	}

	testcases := []struct {
		name     string
		args     values.Object
		expected values.Value
	}{
		{
			name:     "chain success",
			args:     makeArgs(table1, table2),
			expected: table2,
		},
	}

	for _, testcase := range testcases {

		chain := experimental.MakeChainFunction()
		result, err := chain.Call(
			ctx,
			testcase.args,
		)

		if err != nil {
			t.Error(err.Error())
		} else if result != testcase.expected {
			t.Errorf("expected %s, got %s", testcase.expected, result)
		}
	}
}
