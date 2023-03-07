package basic

import (
	"math/big"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func MinItems(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeNumber {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be integer",
		})
	}

	orinum := string(js.(jsi.Number).Value())
	num, ok := new(big.Float).SetString(orinum)
	if !ok || !num.IsInt() {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be integer",
		})
	}
	min, _ := num.Int64()

	return &MinItemsValidator{
		Number: orinum,
		Value:  int(min),
	}, nil
}

func ValidateMinItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type MinItemsValidator struct {
	Number string
	Value  int
}

func (m *MinItemsValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeArray {
		return nil
	}

	n := js.(jsi.Array).Len()
	if n < m.Value {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "array length should be gte " + m.Number,
		})
	}
	return nil
}
