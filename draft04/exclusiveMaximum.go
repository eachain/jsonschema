package draft04

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func exclusiveMaximum(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	maximum := jsi.SiblingOf(js, "maximum")
	if maximum == nil {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should have property maximum when property exclusiveMinimum is present",
		})
	}
	if maximum.Type() != jsi.TypeNumber {
		return nil, nil
	}

	num := maximum.(jsi.Number).Value().String()
	return schema.ValidateFunc(func(ctx *schema.Context, js jsi.JSON) *schema.Result {
		if js.Type() != jsi.TypeNumber {
			return nil
		}
		if jsi.Equal(maximum, js) {
			return schema.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should lt " + num,
			})
		}
		return nil
	}), nil
}
