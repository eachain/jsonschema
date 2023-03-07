package draft04

import (
	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/basic"
)

const Version = "draft-04"

func init() {
	schema.RegisterKeyword(Version, schema.RootKeyword, schema.CompileFunc(basic.RootObject))

	basic.RegisterPointer(Version, basic.PointerKeywords{
		Id:  "id",
		Ref: "$ref",
		// Anchor: "$anchor", // $anchor not support for draft-04
		Immunity: []string{"definitions"},
	})

	// Validation keywords for number and integer
	schema.RegisterKeyword(Version, "multipleOf", basic.ValidateMultipleOf(schema.CompileFunc(basic.MultipleOf)))
	schema.RegisterKeyword(Version, "maximum", basic.ValidateCompare(schema.CompileFunc(basic.Maximum)))
	schema.RegisterKeyword(Version, "exclusiveMaximum", schema.CompileFunc(exclusiveMaximum))
	schema.RegisterKeyword(Version, "minimum", basic.ValidateCompare(schema.CompileFunc(basic.Minimum)))
	schema.RegisterKeyword(Version, "exclusiveMinimum", schema.CompileFunc(exclusiveMinimum))

	// Validation keywords for strings
	schema.RegisterKeyword(Version, "maxLength", basic.ValidateMaxLength(schema.CompileFunc(basic.MaxLength)))
	schema.RegisterKeyword(Version, "minLength", basic.ValidateMinLength(schema.CompileFunc(basic.MinLength)))
	schema.RegisterKeyword(Version, "pattern", basic.ValidatePattern(schema.CompileFunc(basic.Pattern)))
	schema.RegisterKeyword(Version, "format", basic.ValidateFormat(schema.CompileFunc(basic.Format)))

	// Validation keywords for arrays
	schema.RegisterKeyword(Version, "additionalItems", basic.ValidateAdditionalItems(schema.CompileFunc(basic.AdditionalItems)))
	schema.RegisterKeyword(Version, "items", basic.ValidateItems(schema.CompileFunc(basic.Items)))
	// schema.RegisterKeyword(Version, "prefixItems", basic.ValidatePrefixItems(schema.CompileFunc(basic.PrefixItems)))
	schema.RegisterKeyword(Version, "maxItems", basic.ValidateMaxItems(schema.CompileFunc(basic.MaxItems)))
	schema.RegisterKeyword(Version, "minItems", basic.ValidateMinItems(schema.CompileFunc(basic.MinItems)))
	schema.RegisterKeyword(Version, "uniqueItems", basic.ValidateUniqueItems(schema.CompileFunc(basic.UniqueItems)))

	// Validation keywords for objects
	schema.RegisterKeyword(Version, "maxProperties", basic.ValidateMaxProperties(schema.CompileFunc(basic.MaxProperties)))
	schema.RegisterKeyword(Version, "minProperties", basic.ValidateMinProperties(schema.CompileFunc(basic.MinProperties)))
	schema.RegisterKeyword(Version, "required", basic.ValidateRequired(schema.CompileFunc(basic.Required)))
	schema.RegisterKeyword(Version, "additionalProperties", basic.ValidateAdditionalProperties(schema.CompileFunc(basic.AdditionalProperties)))
	schema.RegisterKeyword(Version, "properties", basic.ValidateProperties(schema.CompileFunc(basic.Properties)))
	schema.RegisterKeyword(Version, "patternProperties", basic.ValidatePatternProperties(schema.CompileFunc(basic.PatternProperties)))
	schema.RegisterKeyword(Version, "dependencies", basic.ValidateDependencies(schema.CompileFunc(basic.Dependencies)))

	// Validation keywords for any instance type
	schema.RegisterKeyword(Version, "enum", basic.ValidateEnum(schema.CompileFunc(basic.Enum)))
	schema.RegisterKeyword(Version, "type", schema.CompileFunc(basic.Type))
	schema.RegisterKeyword(Version, "allOf", basic.ValidateAllOf(schema.CompileFunc(basic.AllOf)))
	schema.RegisterKeyword(Version, "anyOf", basic.ValidateAnyOf(schema.CompileFunc(basic.AnyOf)))
	schema.RegisterKeyword(Version, "oneOf", basic.ValidateOneOf(schema.CompileFunc(basic.OneOf)))
	schema.RegisterKeyword(Version, "not", basic.ValidateNot(schema.CompileFunc(basic.Not)))
	schema.RegisterKeyword(Version, "definitions", basic.ValidateDefinitions(schema.CompileFunc(basic.Definitions)))

	// Metadata keywords
	schema.RegisterKeyword(Version, "$schema", schema.CompileFunc(basic.Schema))
	schema.RegisterKeyword(Version, "title", schema.CompileFunc(basic.Title))
	schema.RegisterKeyword(Version, "description", schema.CompileFunc(basic.Description))
	schema.RegisterKeyword(Version, "default", basic.ValidateDefault(schema.CompileFunc(basic.Default)))
}
