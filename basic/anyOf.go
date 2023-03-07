package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func AnyOf(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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
	var anyOf []schema.Validator
	for i := 0; i < arr.Len(); i++ {
		elt := arr.Index(i)
		subctx := ctx.Array(i)
		v, res := root.Compile(subctx, elt)
		result = result.Merge(res)
		if v != nil {
			anyOf = append(anyOf, v)
		}
	}
	if !result.Valid() {
		return nil, result
	}

	return &AnyOfValidator{Conds: anyOf}, result
}

func ValidateAnyOf(cmp schema.Compiler) schema.CompileFunc {
	return cmp.Compile
}

type AnyOfValidator struct {
	Conds []schema.Validator
}

func (a *AnyOfValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	for _, v := range a.Conds {
		rs := v.Validate(ctx, js)
		if rs.Valid() {
			return nil
		}
	}
	return schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should match some schema in anyOf",
	})
}
