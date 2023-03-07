package basic

import (
	"fmt"
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Enum(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be array",
		})
	}

	enum := js.(jsi.Array)
	if enum.Len() == 0 {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should NOT have fewer than 1 items",
		})
	}

	values := make([]jsi.JSON, enum.Len())
	var result *schema.Result
	for i := 0; i < enum.Len(); i++ {
		values[i] = enum.Index(i)
		for j := 0; j < i; j++ {
			if jsi.Equal(values[j], values[i]) {
				result = result.WithError(schema.Error{
					Field: ctx.Array(i).Field(),
					Type:  js.Type(),
					Value: js,
					Msg:   fmt.Sprintf("should NOT have duplicate items (items ## %v and %v are identical)", j, i),
				})
			}
		}
	}

	if !result.Valid() {
		return nil, result
	}

	return &EnumValidator{Values: values}, result
}

func ValidateEnum(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (val schema.Validator, result *schema.Result) {
		val, result = cmp.Compile(ctx, js)
		if js.Type() != jsi.TypeArray {
			return
		}
		enum := js.(jsi.Array)
		if enum.Len() == 0 {
			return
		}

		typs := siblingOfType(js)
		if len(typs) == 0 {
			return
		}

		for i := 0; i < enum.Len(); i++ {
			if !matchOneOf(typs, enum.Index(i)) {
				result = result.WithWarning(schema.Error{
					Field: ctx.Array(i).Field(),
					Type:  js.Type(),
					Value: js,
					Msg: fmt.Sprintf("items ## %v useless for type(s): %v",
						i, strings.Join(typs, ", ")),
				})
			}
		}
		return
	}
}

type EnumValidator struct {
	Values []jsi.JSON
}

func (enum *EnumValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	for i := 0; i < len(enum.Values); i++ {
		if jsi.Equal(enum.Values[i], js) {
			return nil
		}
	}

	return schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should be equal to one of the allowed values",
	})
}
