package basic

import (
	"math/big"
	"strconv"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func MaxProperties(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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

	return &MaxPropertiesValidator{Value: int(max)}, nil
}

func ValidateMaxProperties(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type MaxPropertiesValidator struct {
	Value int
}

func (m *MaxPropertiesValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeObject {
		return nil
	}

	if js.(jsi.Object).Len() > m.Value {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "object properties length should be lte " + strconv.Itoa(m.Value),
		})
	}
	return nil
}
