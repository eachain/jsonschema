package basic

import (
	"math/big"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func MaxItems(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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
	max, _ := num.Int64()

	return &MaxItemsValidator{
		Number: orinum,
		Value:  int(max),
	}, nil
}

func ValidateMaxItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type MaxItemsValidator struct {
	Number string
	Value  int
}

func (m *MaxItemsValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeArray {
		return nil
	}

	n := js.(jsi.Array).Len()
	if n > m.Value {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "array length should be lte " + m.Number,
		})
	}
	return nil
}
