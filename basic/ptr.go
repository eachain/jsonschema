package basic

import (
	"fmt"
	"regexp"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

type PointerKeywords struct {
	Id       string
	Anchor   string
	Ref      string
	Immunity []string
}

func RegisterPointer(draft string, pk PointerKeywords) {
	root := schema.GetKeyword(draft, schema.RootKeyword)
	if root == nil {
		panic("root keyword not register")
	}

	schema.RegisterKeyword(draft, schema.RootKeyword, genCmpRoot(root, pk))
	if pk.Id != "" {
		id := schema.GetKeyword(draft, pk.Id)
		if id == nil {
			schema.RegisterKeyword(draft, pk.Id, schema.CompileFunc(placeholderCompileFunc))
		}
	}
	if pk.Ref != "" {
		ref := schema.GetKeyword(draft, pk.Ref)
		if ref == nil {
			schema.RegisterKeyword(draft, pk.Ref, schema.CompileFunc(placeholderCompileFunc))
		}
	}
	if pk.Anchor != "" {
		anchor := schema.GetKeyword(draft, pk.Anchor)
		if anchor == nil {
			schema.RegisterKeyword(draft, pk.Anchor, schema.CompileFunc(placeholderCompileFunc))
		}
	}
}

func placeholderCompileFunc(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	return nil, nil
}

func placeholderValidateFunc(ctx *schema.Context, js jsi.JSON) *schema.Result {
	return nil
}

type ReferenceValidator struct {
	Id        string
	Ref       string
	Validator schema.Validator
}

func (ref *ReferenceValidator) Validate(ctx *schema.Context, js jsi.JSON) (result *schema.Result) {
	if ref.Validator == nil {
		result = result.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "can't resolve reference " + ref.Ref + " from id " + ref.Id,
		})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			result = result.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   fmt.Sprint(r),
			})
		}
	}()
	result = ref.Validator.Validate(ctx, js)
	return
}

func genCmpRoot(root schema.Compiler, pk PointerKeywords) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (validator schema.Validator, result *schema.Result) {
		if js.Parent() == nil {
			defer func() { result = result.Merge(checkRef(ctx)) }()
		}

		defer func() {
			if validator != nil {
				ctx.Impl(validator)
			} else {
				ctx.Impl(schema.ValidateFunc(placeholderValidateFunc))
			}
		}()

		if js.Type() != jsi.TypeObject {
			validator, result = root.Compile(ctx, js)
			return
		}

		// check is reference
		obj := js.(jsi.Object)

		if pk.Ref != "" {
			ref := obj.Index(pk.Ref)
			if ref != nil {
				validator, result = cmpRef(ctx.Object(pk.Ref), ref, pk.Ref, pk.Immunity)
				return
			}
		}

		if pk.Id != "" {
			id := obj.Index(pk.Id)
			if id != nil {
				result = result.Merge(cmpId(ctx.Object(pk.Id), id))
			}
		}

		if pk.Anchor != "" {
			anchor := obj.Index(pk.Anchor)
			if anchor != nil {
				result = result.Merge(cmpAnchor(ctx.Object(pk.Anchor), anchor))
			}
		}

		validator, res := root.Compile(ctx, js)
		result = result.Merge(res)
		return
	}
}

func cmpRef(ctx *schema.Context, js jsi.JSON, keyword string, immunity []string) (validator schema.Validator, result *schema.Result) {
	if parent := js.Parent(); parent != nil && parent.Type() == jsi.TypeObject {
		obj := parent.(jsi.Object)
		root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
		count := 0
		for _, kw := range immunity {
			if sb := obj.Index(kw); sb != nil {
				_, res := root.Compile(ctx.Parent().Object(kw), sb)
				result = result.Merge(res)
				count++
			}
		}
		if count+1 < obj.Len() {
			result = result.WithWarning(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "property " + keyword + " here, any other properties will be ignored",
			})
		}
	}

	if js.Type() != jsi.TypeString {
		result = result.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
		return
	}

	ref := js.(jsi.String).Value()
	ptr, err := schema.ParsePointer(ref)
	if err != nil {
		result = result.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri",
		})
		return
	}

	refVal := &ReferenceValidator{
		Id:  ctx.Id(),
		Ref: ref,
	}
	ctx.PushRef(js, ptr, &refVal.Validator)

	validator = refVal
	return
}

func cmpId(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	val := js.(jsi.String).Value()
	id, err := schema.ParsePointer(val)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri",
		})
	}
	var result *schema.Result
	if len(id.Frag) > 0 {
		result = result.WithWarning(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "fragment is ignored",
		})
	}
	id.Frag = nil
	ctx.Parent().SetId(id)
	return result
}

var anchorFormat = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_:.]*`)

func cmpAnchor(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	anchor := js.(jsi.String).Value()
	if !anchorFormat.MatchString(anchor) {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should start with a letter followed by any number of letters, digits, -, _, :, or .",
		})
	}

	ctx.Parent().SetAnchor(anchor)
	return nil
}

func checkRef(ctx *schema.Context) (result *schema.Result) {
	for _, ref := range ctx.FillRefs() {
		if *ref.Validator == nil {
			result = result.Merge(loadSchemaFromURI(ref.Context, ref.JSON, ref.Ref))
		}
	}
	return
}

func loadSchemaFromURI(ctx *schema.Context, js jsi.JSON, ref *schema.Pointer) *schema.Result {
	u := ref.URL()
	u.Fragment = ""

	subjs, err := jsi.NewURLParser(u).Parse()
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "load reference from uri: " + err.Error(),
		})
	}

	root := schema.GetKeyword(ctx.Draft(), schema.RootKeyword)
	_, result := root.Compile(ctx.WithId(ref), subjs)
	return result
}
