package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Not(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	val, result := root.Compile(ctx, js)
	if !result.Valid() {
		return nil, result
	}

	return &NotValidator{Cond: val}, result
}

func ValidateNot(cmp schema.Compiler) schema.CompileFunc {
	return cmp.Compile
}

type NotValidator struct {
	Cond schema.Validator
}

func (a *NotValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	rs := a.Cond.Validate(ctx, js)
	if !rs.Valid() {
		return nil
	}
	return schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should NOT be valid",
	})
}
