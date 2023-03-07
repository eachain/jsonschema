package basic

import (
	"fmt"
	"math/big"
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Type(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() == jsi.TypeString {
		return singleType(ctx, js)
	}

	if js.Type() == jsi.TypeArray {
		return multiTypes(ctx, js)
	}

	return nil, schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should be string or array",
	})
}

func singleType(ctx *schema.Context, js jsi.JSON) (*TypeValidator, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	typ := js.(jsi.String).Value()
	if !inTypes(schema.AllTypes, typ) {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be one of " + strings.Join(schema.AllTypes, ", "),
		})
	}

	return &TypeValidator{Types: []schema.Type{typ}}, nil
}

func multiTypes(ctx *schema.Context, js jsi.JSON) (*TypeValidator, *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be array",
		})
	}

	arr := js.(jsi.Array)
	if arr.Len() == 0 {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should NOT have fewer than 1 types",
		})
	}

	tv := new(TypeValidator)
	var result *schema.Result
	for i := 0; i < arr.Len(); i++ {
		v, res := singleType(ctx.Array(i), arr.Index(i))
		result = result.Merge(res)
		if v != nil {
			tv.Types = append(tv.Types, v.Types...)
		}
	}

	for i := 0; i < len(tv.Types); i++ {
		for j := 0; j < i; j++ {
			if tv.Types[i] == tv.Types[j] {
				result = result.WithError(schema.Error{
					Field: ctx.Array(i).Field(),
					Type:  arr.Index(i).Type(),
					Value: arr.Index(i),
					Msg:   fmt.Sprintf("should NOT have duplicate types (types ## %v and %v are identical)", j, i),
				})
			}
		}
	}

	return tv, result
}

type TypeValidator struct {
	Types []schema.Type
}

func (tv *TypeValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	var hasInt bool
	for _, typ := range tv.Types {
		if typ == schema.TypeInteger {
			hasInt = true
		}
		if js.Type() == typ {
			return nil
		}
	}
	if hasInt && js.Type() == jsi.TypeNumber {
		x, ok := new(big.Float).SetString(js.(jsi.Number).Value().String())
		if ok && x.IsInt() {
			return nil
		}
	}

	return schema.WithError(schema.Error{
		Field: ctx.Field(),
		Type:  js.Type(),
		Value: js,
		Msg:   "should be one of the allowed types",
	})
}
