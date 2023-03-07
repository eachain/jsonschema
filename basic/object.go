package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func RootObject(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be object",
		})
	}

	var result *schema.Result
	var vals []schema.Validator

	obj := js.(jsi.Object)
	iter := obj.Iter()
	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	for iter.Next() {
		key, js := iter.Entry()
		cmp := schema.GetKeyword(ctx.Draft(), key)
		if cmp == nil {
			_, res := root.Compile(ctx.Object(key), js)
			result = result.Merge(res)
			if result.Valid() {
				result = result.WithWarning(schema.Error{
					Field: ctx.Field(),
					Type:  js.Type(),
					Value: js,
					Msg:   "keyword not defined: " + key,
				})
			}
			continue
		}
		val, res := cmp.Compile(ctx.Object(key), js)
		result = result.Merge(res)
		if val != nil {
			vals = append(vals, val)
		}
	}
	// }
	if !result.Valid() {
		return nil, result
	}

	return &RootObjectValidator{Validators: vals}, result
}

func ValidateObjectType(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeObject)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type RootObjectValidator struct {
	Validators []schema.Validator
}

func (obj *RootObjectValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	for _, v := range obj.Validators {
		result = result.Merge(v.Validate(ctx, js))
	}
	return
}
