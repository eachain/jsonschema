package basic

import (
	"strconv"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func AdditionalItems(ctx *schema.Context, js jsi.JSON) (val schema.Validator, result *schema.Result) {
	if js.Type() == jsi.TypeBoolean {
		if js.(jsi.Boolean).Value() {
			return nil, nil
		}
	} else {
		root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
		val, result = root.Compile(ctx, js)
		if !result.Valid() {
			return nil, result
		}
	}

	items := jsi.SiblingOf(js, "items")
	if items == nil || items.Type() != jsi.TypeArray {
		msg := "not defined"
		if items != nil {
			msg = string(items.Type())
		}
		result = result.WithWarning(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "useless when sibling 'items' is " + msg,
		})
		return
	}

	max := items.(jsi.Array).Len()
	if val == nil {
		return &MaxItemsValidator{
			Number: strconv.Itoa(max),
			Value:  max,
		}, result
	}
	return &AdditionalItemsValidator{
		ItemsCount: max,
		Additional: val,
	}, result
}

func ValidateAdditionalItems(cmp schema.Compiler) schema.CompileFunc {
	return ValidateArrayType(cmp)
}

type AdditionalItemsValidator struct {
	ItemsCount int
	Additional schema.Validator
}

func (ai *AdditionalItemsValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if js.Type() != jsi.TypeArray {
		return
	}

	arr := js.(jsi.Array)
	for i := ai.ItemsCount; i < arr.Len(); i++ {
		result = result.Merge(ai.Additional.Validate(ctx.Array(i), arr.Index(i)))
	}
	return
}
