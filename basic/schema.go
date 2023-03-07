package basic

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

func Schema(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	switch schema := js.(jsi.String).Value(); schema {
	default:
		return nil, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "schema unrecognized",
		})
	case "http://json-schema.org/draft-04/schema", "http://json-schema.org/draft-04/schema#":
		ctx.UseDraft("draft-04")
	case "http://json-schema.org/draft-06/schema", "http://json-schema.org/draft-06/schema#":
		ctx.UseDraft("draft-06")
	case "http://json-schema.org/draft-07/schema", "http://json-schema.org/draft-07/schema#":
		ctx.UseDraft("draft-07")
	case "https://json-schema.org/draft/2019-09/schema":
		ctx.UseDraft("draft-201909")
	case "https://json-schema.org/draft/2020-12/schema":
		ctx.UseDraft("draft-202012")
	}

	return nil, nil
}
