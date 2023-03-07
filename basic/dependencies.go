package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Dependencies(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be object",
		})
	}

	var result *schema.Result

	depOf := make(map[string]schema.Validator)
	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		subctx := ctx.Object(key)

		var val schema.Validator
		var rs *schema.Result
		if js.Type() == jsi.TypeArray {
			val, rs = Required(subctx, js)
		} else {
			val, rs = root.Compile(subctx, js)
		}
		result = result.Merge(rs)
		if rs.Valid() && val != nil {
			depOf[key] = val
		}
	}

	if !result.Valid() {
		return nil, result
	}

	return &DependenciesValidator{Dependency: depOf}, result
}

func ValidateDependencies(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type DependenciesValidator struct {
	Dependency map[string]schema.Validator
}

func (dv *DependenciesValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return
	}

	obj := js.(jsi.Object)
	for key, val := range dv.Dependency {
		if obj.Index(key) == nil {
			continue
		}
		result = result.Merge(val.Validate(ctx.Object(key), js))
	}
	return
}
