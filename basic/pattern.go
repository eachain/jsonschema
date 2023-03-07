package basic

import (
	"regexp"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Pattern(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	expr := js.(jsi.String).Value()
	re, err := regexp.Compile(expr)
	if err != nil {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be regex",
		})
	}

	return &PatternValidator{Regex: re}, nil
}

func ValidatePattern(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeString)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type PatternValidator struct {
	Regex *regexp.Regexp
}

func (p *PatternValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}

	val := js.(jsi.String).Value()
	if !p.Regex.MatchString(val) {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "string should match regex: " + p.Regex.String(),
		})
	}
	return nil
}
