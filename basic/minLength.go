package basic

import (
	"math/big"
	"unicode/utf8"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func MinLength(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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
	if !ok || !num.IsInt() {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be integer",
		})
	}
	min, _ := num.Int64()

	return &MinLengthValidator{
		Number: orinum,
		Value:  int(min),
	}, nil
}

func ValidateMinLength(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeString)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type MinLengthValidator struct {
	Number string
	Value  int
}

func (m *MinLengthValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}

	val := js.(jsi.String).Value()
	if utf8.RuneCountInString(val) < m.Value {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "string length should be gte " + m.Number,
		})
	}
	return nil
}
