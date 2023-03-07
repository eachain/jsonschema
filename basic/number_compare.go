package basic

import (
	"math/big"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

// cmp: func to judge the compare result
//
//	-1 if doc.val <  schema.val
//	 0 if doc.val == schema.val
//	+1 if doc.val >  schema.val
func Compare(cmp func(int) bool, op string) schema.Compiler {
	return schema.CompileFunc(func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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

		return &CompareValidator{
			HandleCompareResultFunc: cmp,
			Op:                      op,
			Number:                  orinum,
			Value:                   num,
		}, nil
	})
}

func ValidateCompare(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeNumber, schema.TypeInteger)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type CompareValidator struct {
	// HandleCompareResultFunc: func to judge the compare result
	//
	//	-1 if doc.val <  schema.val
	//	 0 if doc.val == schema.val
	//	+1 if doc.val >  schema.val
	HandleCompareResultFunc func(int) bool
	Op                      string
	Number                  string
	Value                   *big.Float
}

func (c *CompareValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
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

	if !c.HandleCompareResultFunc(val.Cmp(c.Value)) {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should " + c.Op + " " + c.Number,
		})
	}
	return nil
}
