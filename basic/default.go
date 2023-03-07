package basic

import (
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Default(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	return nil, nil
}

func ValidateDefault(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		typs := siblingOfType(js)
		if len(typs) == 0 {
			if isInUnion(ctx, js) {
				return val, result
			}
			return nil, result.WithWarning(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "sibling field 'type' not found",
			})
		}

		if !matchOneOf(typs, js) {
			if len(typs) == 1 {
				result = result.WithWarning(schema.Error{
					Field: ctx.Field(),
					Type:  js.Type(),
					Value: js,
					Msg:   "should be " + typs[0],
				})
			} else {
				result = result.WithWarning(schema.Error{
					Field: ctx.Field(),
					Type:  js.Type(),
					Value: js,
					Msg:   "should be one of " + strings.Join(typs, ", "),
				})
			}
		}

		return val, result
	}
}
