package basic

import (
	"math/big"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func MultipleOf(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeNumber {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be number",
		})
	}

	orinum := string(js.(jsi.Number).Value())
	num, ok := new(big.Float).SetString(orinum)
	if !ok {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be number",
		})
	}

	return &MultipleOfValidator{
		Number: orinum,
		Value:  num,
	}, nil
}

func ValidateMultipleOf(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeNumber, schema.TypeInteger)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type MultipleOfValidator struct {
	Number string
	Value  *big.Float
}

func (m *MultipleOfValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeNumber {
		return nil
	}

	val, ok := new(big.Float).SetString(string(js.(jsi.Number).Value()))
	if !ok {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be number",
		})
	}

	if !val.Quo(val, m.Value).IsInt() {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be multiple of " + m.Number,
		})
	}
	return nil
}
