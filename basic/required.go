package basic

import (
	"fmt"
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Required(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be array",
		})
	}

	var result *schema.Result
	arr := js.(jsi.Array)
	size := arr.Len()
	requires := make([]string, size)
	for i := 0; i < arr.Len(); i++ {
		elt := arr.Index(i)
		if elt.Type() != jsi.TypeString {
			result = result.WithError(schema.Error{
				Field: ctx.Array(i).Field(),
				Type:  elt.Type(),
				Value: elt,
				Msg:   "should be string",
			})
			continue
		}
		key := elt.(jsi.String).Value()
		requires[i] = key
		for j := 0; j < i; j++ {
			if requires[j] == requires[i] {
				result = result.WithError(schema.Error{
					Field: ctx.Array(i).Field(),
					Type:  elt.Type(),
					Value: elt,
					Msg:   fmt.Sprintf("should NOT have duplicate items (items ## %v and %v are identical)", j, i),
				})
				break
			}
		}
	}

	return &RequiredValidator{Keys: requires}, result
}

func ValidateRequired(cmp schema.Compiler) schema.CompileFunc {
	return ValidateObjectType(cmp)
}

type RequiredValidator struct {
	Keys []string
}

func (rv *RequiredValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeObject {
		return nil
	}

	obj := js.(jsi.Object)

	var required []string
	for _, key := range rv.Keys {
		if obj.Index(key) == nil {
			required = append(required, key)
		}
	}

	if len(required) > 0 {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "object properties required keys: " + strings.Join(required, ", "),
		})
	}
	return nil
}
