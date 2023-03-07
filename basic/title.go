package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Title(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	return nil, nil
}
