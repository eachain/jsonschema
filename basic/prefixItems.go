package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func PrefixItems(ctx *schema.Context, js jsi.JSON) (val schema.Validator, result *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be array",
		})
	}

	arr := js.(jsi.Array)
	vs := make([]schema.Validator, arr.Len())
	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	for i := 0; i < arr.Len(); i++ {
		v, res := root.Compile(ctx.Array(i), arr.Index(i))
		result = result.Merge(res)
		if v != nil {
			vs[i] = v
		}
	}
	if !result.Valid() {
		return
	}

	val = &PrefixItemsValidator{Items: vs}
	return
}

func ValidatePrefixItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type PrefixItemsValidator struct {
	Items []schema.Validator
}

func (pi *PrefixItemsValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return
	}

	arr := js.(jsi.Array)
	n := arr.Len()
	if n > len(pi.Items) {
		n = len(pi.Items)
	}

	for i := 0; i < n; i++ {
		if pi.Items[i] != nil {
			result = result.Merge(pi.Items[i].Validate(ctx, arr.Index(i)))
		}
	}
	return
}
