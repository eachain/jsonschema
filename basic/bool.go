package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func RootBoolean(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeBoolean {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be boolean",
		})
	}
	if js.(jsi.Boolean).Value() {
		return nil, nil
	}

	return schema.ValidateFunc(func(ctx *schema.Context, js jsi.JSON) *schema.Result {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "invalid schema",
		})
	}), nil
}
