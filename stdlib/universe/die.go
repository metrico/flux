package universe

import (
	"context"

	"github.com/InfluxCommunity/flux"
	"github.com/InfluxCommunity/flux/codes"
	"github.com/InfluxCommunity/flux/interpreter"
	"github.com/InfluxCommunity/flux/runtime"
	"github.com/InfluxCommunity/flux/values"
)

const DieKind = "die"

func init() {
	runtime.RegisterPackageValue("universe", DieKind, Die())
}

func Die() values.Function {
	return values.NewFunction(
		DieKind,
		runtime.MustLookupBuiltinType("universe", DieKind),
		func(ctx context.Context, args values.Object) (values.Value, error) {
			return interpreter.DoFunctionCallContext(func(ctx context.Context, args interpreter.Arguments) (values.Value, error) {
				msg, err := args.GetRequiredString("msg")
				if err != nil {
					return nil, err
				} else {
					return nil, &flux.Error{
						Code: codes.Invalid,
						Msg:  msg,
					}
				}
			}, ctx, args)
		}, false,
	)
}
