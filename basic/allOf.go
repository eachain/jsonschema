package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func AllOf(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be array",
		})
	}

	arr := js.(jsi.Array)
	if arr.Len() == 0 {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should NOT have fewer than 1 items",
		})
	}

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)

	var result *schema.Result
	var allOf []schema.Validator
	for i := 0; i < arr.Len(); i++ {
		elt := arr.Index(i)
		subctx := ctx.Array(i)
		v, res := root.Compile(subctx, elt)
		result = result.Merge(res)
		if v != nil {
			allOf = append(allOf, v)
		}
	}
	if !result.Valid() {
		return nil, result
	}

	return &AllOfValidator{Conds: allOf}, result
}

func ValidateAllOf(cmp schema.Compiler) schema.CompileFunc {
	return cmp.Compile
}

type AllOfValidator struct {
	Conds []schema.Validator
}

func (a *AllOfValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	for _, v := range a.Conds {
		result = result.Merge(v.Validate(ctx, js))
	}
	return
}
