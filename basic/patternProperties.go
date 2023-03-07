package basic

import (
	"regexp"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func PatternProperties(ctx *schema.Context, js jsi.JSON) (validator schema.Validator, result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be object",
		})
	}

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	var props []PatternPropertyValidator
	prop := js.(jsi.Object).Iter()
	for prop.Next() {
		expr, js := prop.Entry()
		subctx := ctx.Object(expr)
		re, err := regexp.Compile(expr)
		if err != nil {
			result = result.WithError(schema.Error{
				Field: subctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should be regex",
			})
		}

		val, rs := root.Compile(subctx, js)
		result = result.Merge(rs)
		if rs.Valid() && val != nil {
			props = append(props, PatternPropertyValidator{
				Regex:     re,
				Validator: val,
			})
		}
	}
	if !result.Valid() {
		return
	}

	validator = &PatternPropertiesValidator{
		Validators: props,
	}
	return
}

func ValidatePatternProperties(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type PatternPropertyValidator struct {
	Regex     *regexp.Regexp
	Validator schema.Validator
}

type PatternPropertiesValidator struct {
	Validators []PatternPropertyValidator
}

func (pp *PatternPropertiesValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return
	}

	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		subctx := ctx.Object(key)
		for _, pv := range pp.Validators {
			if pv.Regex.MatchString(key) {
				result = result.Merge(pv.Validator.Validate(subctx, js))
			}
		}
	}
	return
}
