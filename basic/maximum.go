package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Maximum(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	return Compare(func(result int) bool {
		return result <= 0
	}, "lte").Compile(ctx, js)
}
