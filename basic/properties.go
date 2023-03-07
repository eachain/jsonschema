package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Properties(ctx *schema.Context, js jsi.JSON) (validator schema.Validator, result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be object",
		})
	}

	property := make(map[string]schema.Validator)

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		val, rs := root.Compile(ctx.Object(key), js)
		result = result.Merge(rs)
		if val != nil {
			property[key] = val
		}
	}

	if !result.Valid() {
		return
	}

	validator = &PropertiesValidator{Prop: property}
	return
}

func ValidateProperties(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type PropertiesValidator struct {
	Prop map[string]schema.Validator
}

func (p *PropertiesValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return
	}

	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		if val := p.Prop[key]; val != nil {
			result = result.Merge(val.Validate(ctx.Object(key), js))
		}
	}

	return
}
