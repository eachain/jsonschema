package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Items(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() == jsi.TypeArray {
		return PrefixItems(ctx, js)
	}

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	val, result := root.Compile(ctx, js)
	if !result.Valid() {
		return nil, result
	}

	// var prefix int
	// if prefixItems := jsi.SiblingOf(js, "prefixItems"); prefixItems != nil {
	// 	if prefixItems.Type() == jsi.TypeArray {
	// 		prefix = prefixItems.(jsi.Array).Len()
	// 	}
	// }

	return &ItemTypeValidator{Validator: val}, nil
}

func ValidateItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type ItemTypeValidator struct {
	Validator schema.Validator
}

func (it *ItemTypeValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return
	}

	arr := js.(jsi.Array)
	for i := 0; i < arr.Len(); i++ {
		result = result.Merge(it.Validator.Validate(ctx, arr.Index(i)))
	}
	return
}
