package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func ValidateArrayType(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeArray)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}
