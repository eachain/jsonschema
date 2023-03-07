package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Definitions(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be object",
		})
	}

	var result *schema.Result

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		subctx := ctx.Object(key)
		_, rs := root.Compile(subctx, js)
		result = result.Merge(rs)
	}

	return nil, result
}

func ValidateDefinitions(cmp schema.Compiler) schema.CompileFunc {
	return cmp.Compile
}
