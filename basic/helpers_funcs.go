package basic

import (
	"math/big"
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func siblingOfType(js jsi.JSON) []jsi.Type {
	typ := jsi.SiblingOf(js, "type")
	if typ == nil {
		return nil
	}

	if typ.Type() == jsi.TypeString {
		return []jsi.Type{jsi.Type(typ.(jsi.String).Value())}
	}
	if typ.Type() != jsi.TypeArray {
		return nil
	}

	arr := typ.(jsi.Array)
	var typs []jsi.Type
	for i := 0; i < arr.Len(); i++ {
		elt := arr.Index(i)
		if elt.Type() == jsi.TypeString {
			typs = append(typs, jsi.Type(elt.(jsi.String).Value()))
		}
	}
	return typs
}

func inTypes(ls []jsi.Type, s jsi.Type) bool {
	for _, v := range ls {
		if v == s {
			return true
		}
	}
	return false
}

func oneOfTypes(ls []jsi.Type, ss ...jsi.Type) bool {
	for _, s := range ss {
		if inTypes(ls, s) {
			return true
		}
	}
	return false
}

func isInUnion(ctx *schema.Context, js jsi.JSON) bool {
	switch ctx.Parent().Parent().Index() {
	case "allOf", "oneOf", "anyOf":
		return true
	}

	return jsi.SiblingOf(js, "allOf") != nil ||
		jsi.SiblingOf(js, "oneOf") != nil ||
		jsi.SiblingOf(js, "anyOf") != nil
}

func checkSiblingOfType(ctx *schema.Context, js jsi.JSON, oneof ...jsi.Type) *schema.Error {
	typs := siblingOfType(js)
	if len(typs) == 0 {
		if isInUnion(ctx, js) {
			return nil
		}
		return &schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "sibling field 'type' not found",
		}
	}
	return checkOneOfType(ctx, js, typs, oneof...)
}

func checkOneOfType(ctx *schema.Context, js jsi.JSON, typs []jsi.Type, oneof ...jsi.Type) *schema.Error {
	if !oneOfTypes(typs, oneof...) {
		if len(typs) == 1 {
			return &schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "useless for 'type' " + string(typs[0]),
			}
		}
		return &schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "useless for any 'type' of " + strings.Join(typs, ", "),
		}
	}
	return nil
}

func matchOneOf(typs []jsi.Type, js jsi.JSON) bool {
	if inTypes(typs, js.Type()) {
		return true
	}

	if inTypes(typs, schema.TypeInteger) && js.Type() == jsi.TypeNumber {
		n, ok := new(big.Float).SetString(js.(jsi.Number).Value().String())
		if !ok {
			return false
		}
		return n.IsInt()
	}

	return false
}
