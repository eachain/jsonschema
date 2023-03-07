package basic

import (
	"regexp"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func AdditionalProperties(ctx *schema.Context, js jsi.JSON) (val schema.Validator, result *schema.Result) {
	var validateAdditional schema.Validator
	if js.Type() == jsi.TypeBoolean {
		if js.(jsi.Boolean).Value() {
			return
		}
		validateAdditional = schema.ValidateFunc(func(ctx *schema.Context, js jsi.JSON) *schema.Result {
			return schema.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should NOT have additional properties",
			})
		})
	} else {
		root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
		validateAdditional, result = root.Compile(ctx, js)
		if !result.Valid() {
			return
		}
	}

	inProp := make(map[string]bool)
	if pp := jsi.SiblingOf(js, "properties"); pp != nil && pp.Type() == jsi.TypeObject {
		prop := pp.(jsi.Object).Iter()
		for prop.Next() {
			key, _ := prop.Entry()
			inProp[key] = true
		}
	}

	var regs []*regexp.Regexp
	if pp := jsi.SiblingOf(js, "patternProperties"); pp != nil && pp.Type() == jsi.TypeObject {
		prop := pp.(jsi.Object).Iter()
		for prop.Next() {
			expr, _ := prop.Entry()
			if re, err := regexp.Compile(expr); err == nil {
				regs = append(regs, re)
			}
		}
	}

	isProperty := func(key string) bool {
		if inProp[key] {
			return true
		}
		for _, re := range regs {
			if re.MatchString(key) {
				return true
			}
		}
		return false
	}

	val = &AdditionalPropertiesValidator{
		IsProperty: isProperty,
		Additional: validateAdditional,
	}
	return
}

func ValidateAdditionalProperties(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type AdditionalPropertiesValidator struct {
	IsProperty func(key string) bool
	Additional schema.Validator
}

func (ap *AdditionalPropertiesValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeObject {
		return
	}

	iter := js.(jsi.Object).Iter()
	for iter.Next() {
		key, js := iter.Entry()
		if ap.IsProperty(key) {
			continue
		}

		result = result.Merge(ap.Additional.Validate(ctx.Object(key), js))
	}
	return
}
