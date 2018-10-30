package transformations

import (
	"github.com/influxdata/flux"
)

func init() {
	flux.RegisterBuiltIn("helpers", helpersBuiltIn)
}

var helpersBuiltIn = `
// AggregateWindow applies an aggregate function to fixed windows of time.
aggregateWindow = (every, fn, columns=["_value"], timeSrc="_stop",timeDst="_time", table=<-) =>
	table
		|> window(every:every)
		|> fn(columns:columns)
		|> duplicate(column:timeSrc,as:timeDst)
		|> window(every:inf, timeCol:timeDst)
`
