package draft04

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func exclusiveMinimum(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	minimum := jsi.SiblingOf(js, "minimum")
	if minimum == nil {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should have property minimum when property exclusiveMinimum is present",
		})
	}
	if minimum.Type() != jsi.TypeNumber {
		return nil, nil
	}

	num := minimum.(jsi.Number).Value().String()
	return schema.ValidateFunc(func(ctx *schema.Context, js jsi.JSON) *schema.Result {
		if js.Type() != jsi.TypeNumber {
			return nil
		}
		if jsi.Equal(minimum, js) {
			return schema.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should gt " + num,
			})
		}
		return nil
	}), nil
}
