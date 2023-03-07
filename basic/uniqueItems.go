package basic

import (
	"fmt"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func UniqueItems(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeBoolean {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be boolean",
		})
	}
	if !js.(jsi.Boolean).Value() {
		return nil, nil
	}

	return (*UniqueItemsValidator)(nil), nil
}

func ValidateUniqueItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type UniqueItemsValidator struct{}

func (*UniqueItemsValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return
	}

	arr := js.(jsi.Array)
	for i := 1; i < arr.Len(); i++ {
		for j := 0; j < i; j++ {
			if jsi.Equal(arr.Index(j), arr.Index(i)) {
				result = result.WithError(schema.Error{
					Field: ctx.Field(),
					Type:  js.Type(),
					Value: js,
					Msg:   fmt.Sprintf("should NOT have duplicate items (items ## %v and %v are identical)", i, j),
				})
			}
		}
	}
	return
}
