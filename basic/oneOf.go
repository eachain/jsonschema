package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func OneOf(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
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
	var oneOf []schema.Validator
	for i := 0; i < arr.Len(); i++ {
		elt := arr.Index(i)
		subctx := ctx.Array(i)
		v, res := root.Compile(subctx, elt)
		result = result.Merge(res)
		if v != nil {
			oneOf = append(oneOf, v)
		}
	}
	if !result.Valid() {
		return nil, result
	}

	return &OneOfValidator{Conds: oneOf}, result
}

func ValidateOneOf(cmp schema.Compiler) schema.CompileFunc {
	return cmp.Compile
}

type OneOfValidator struct {
	Conds []schema.Validator
}

func (a *OneOfValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	matched := 0
	for _, v := range a.Conds {
		rs := v.Validate(ctx, js)
		if !rs.Valid() {
			continue
		}
		matched++
		if matched > 1 {
			break
		}
	}
	if matched == 1 {
		return nil
	}
	println("oneof matched:", ctx.Field(), matched)
	return schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should match exactly one schema in oneOf",
	})
}
