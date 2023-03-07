//go:build ignore

package basic

import (
	"net/url"
	"regexp"
	"strings"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

type PointerKeywords struct {
	Id     string
	Anchor string
	Ref    string
}

func RegisterPointer(draft string, pk PointerKeywords) {
	root := schema.GetKeyword(draft, schema.RootKeyword)
	if root == nil {
		panic("root keyword not register")
	}

	schema.RegisterKeyword(draft, schema.RootKeyword, genCmpRoot(draft, root, pk))
	if pk.Id != "" {
		id := schema.GetKeyword(draft, pk.Id)
		if id == nil {
			schema.RegisterKeyword(draft, pk.Id, schema.CompileFunc(placeholder))
		}
	}
	if pk.Ref != "" {
		ref := schema.GetKeyword(draft, pk.Ref)
		if ref == nil {
			schema.RegisterKeyword(draft, pk.Ref, schema.CompileFunc(placeholder))
		}
	}
	if pk.Anchor != "" {
		anchor := schema.GetKeyword(draft, pk.Anchor)
		if anchor == nil {
			schema.RegisterKeyword(draft, pk.Anchor, schema.CompileFunc(placeholder))
		}
	}
}

func placeholder(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	return nil, nil
}

type refCtx struct {
	ctx   *schema.Context
	js    jsi.JSON
	ref   string
	inDoc bool
	val   *schema.Validator
}

type ptrCtx struct {
	draft  string
	root   schema.Compiler
	ref    map[string]*refCtx
	all    map[string]schema.Validator
	id     *url.URL
	exists map[string]bool // id exists
	pk     PointerKeywords
	inRef  bool
	parent *ptrCtx
}

func genCmpRoot(draft string, root schema.Compiler, pk PointerKeywords) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		pc := &ptrCtx{
			draft:  draft,
			root:   root,
			ref:    make(map[string]*refCtx),
			all:    make(map[string]schema.Validator),
			id:     new(url.URL),
			exists: make(map[string]bool),
			pk:     pk,
			inRef:  false,
		}
		return pc.cmpRoot(ctx, js)
	}
}

func (pc *ptrCtx) withId(id *url.URL) *ptrCtx {
	return &ptrCtx{
		draft:  pc.draft,
		root:   pc.root,
		ref:    pc.ref,
		all:    pc.all,
		id:     id,
		exists: pc.exists,
		pk:     pc.pk,
		inRef:  pc.inRef,
		parent: pc,
	}
}

func (pc *ptrCtx) cmpRoot(ctx *schema.Context, js jsi.JSON) (validator schema.Validator, result *schema.Result) {
	println("compile:", ctx.Path())
	if js == jsi.RootOf(js) {
		defer func() { result = result.Merge(pc.check()) }()
	}

	defer func() {
		if validator != nil {
			println("register validator:", ctx.Path())
			pc.all[ctx.Path()] = validator
		}
	}()

	if js.Type() != jsi.TypeObject {
		validator, result = pc.root.Compile(ctx, js)
		return
	}

	// check is reference
	obj := js.(jsi.Object)

	if pc.pk.Ref != "" {
		ref := obj.Index(pc.pk.Ref)
		if ref != nil {
			return pc.cmpRef(ctx.WithKeyword(pc.pk.Ref), ref)
		}
	}

	if pc.pk.Id != "" {
		id := obj.Index(pc.pk.Id)
		if id != nil {
			var res *schema.Result
			pc, res = pc.cmpId(ctx.WithKeyword(pc.pk.Id), id)
			result.Merge(res)
		}
	}

	if pc.pk.Anchor != "" {
		anchor := obj.Index(pc.pk.Anchor)
		if anchor != nil {
			ptr, res := pc.cmpAnchor(ctx.WithKeyword(pc.pk.Anchor), anchor)
			result = result.Merge(res)
			if ptr != "" {
				defer func() {
					if validator != nil {
						pc.all[ptr] = validator
					}
				}()
			}
		}
	}

	validator, res := pc.root.Compile(ctx, js)
	result = result.Merge(res)
	return
}

func (pc *ptrCtx) cmpRef(ctx *schema.Context, js jsi.JSON) (validator schema.Validator, result *schema.Result) {
	if parent := js.Parent(); parent != nil && parent.Type() == jsi.TypeObject {
		if parent.(jsi.Object).Len() > 1 {
			result = result.WithWarning(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "property " + pc.pk.Ref + " here, any other properties will be ignored",
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

	rc := &refCtx{
		ctx:   ctx,
		js:    js,
		ref:   js.(jsi.String).Value(),
		inDoc: pc.inRef,
		val:   new(schema.Validator),
	}

	var ref *url.URL
	if strings.HasPrefix(rc.ref, "#") {
		ref = new(url.URL)
	} else {
		var err error
		ref, err = url.Parse(rc.ref)
		if err != nil {
			result = result.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should be uri",
			})
			return
		}
		ref.RawFragment = ""
		ref.Fragment = ""
	}
	pc.diff(ref)

	// full uri
	if ref.Scheme != "" || ref.Host != "" || ref.Path != "" {
		ptr := ref.String() + "#"
		if ori := pc.ref[ptr]; ori == nil {
			pc.ref[ptr] = rc
		}
	} else if ref.Fragment != "" {
		path := "#" + ref.Fragment
		if ori := pc.ref[path]; ori == nil {
			pc.ref[path] = rc
		}
	}

	validator = &RefValidator{Id: pc.id.String() + "#", Ref: rc.ref, Validator: rc.val}
	return
}

func (pc *ptrCtx) cmpId(ctx *schema.Context, js jsi.JSON) (*ptrCtx, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return pc, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	val := js.(jsi.String).Value()
	u, err := url.ParseRequestURI(val)
	if err != nil {
		return pc, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri",
		})
	}
	if u.Fragment != "" {
		u.Fragment = ""
	}
	id := pc.diff(u).String()
	if pc.exists[id] {
		return pc, schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "id duplicated",
		})
	}
	pc.exists[id] = true
	return pc.withId(u), nil
}

var anchorFormat = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_:.]*`)

func (pc *ptrCtx) cmpAnchor(ctx *schema.Context, js jsi.JSON) (string, *schema.Result) {
	if js.Type() != jsi.TypeString {
		return "", schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be string",
		})
	}

	anchor := js.(jsi.String).Value()

	if !anchorFormat.MatchString(anchor) {
		return "", schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should start with a letter followed by any number of letters, digits, -, _, :, or .",
		})
	}

	ptr := pc.id.String() + "#" + anchor
	if pc.exists[ptr] {
		return "", schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "anchor duplicated",
		})
	}
	pc.exists[ptr] = true
	return ptr, nil
}

func (pc *ptrCtx) diff(ref *url.URL) *url.URL {
	if ref.Scheme != "" {
		if ref.Scheme == pc.id.Scheme {
			ref.Scheme = ""
		} else {
			return ref
		}
	}

	if ref.Host != "" {
		if ref.Host == pc.id.Host {
			ref.Host = ""
		} else {
			ref.Scheme = pc.id.Scheme
			return ref
		}
	}

	if ref.Path != "" {
		if ref.Path == pc.id.Path {
			ref.Path = ""
		} else {
			ref.Scheme = pc.id.Scheme
			ref.Host = pc.id.Host
			return ref
		}
	}

	if ref.RawQuery != "" {
		if ref.RawQuery == pc.id.RawQuery {
			ref.RawQuery = ""
		} else {
			ref.Scheme = pc.id.Scheme
			ref.Host = pc.id.Host
			ref.Path = pc.id.Path
			return ref
		}
	}

	if !pc.inRef {
		ref.Scheme = pc.id.Scheme
		ref.Host = pc.id.Host
		ref.Path = pc.id.Path
		ref.RawQuery = pc.id.RawQuery
	}
	return ref
}

func (pc *ptrCtx) check() (result *schema.Result) {
	pc.inRef = true
	for uri, ref := range pc.ref {
		println("- - - - - - ->", uri, *ref.val != nil)
		if *ref.val == nil {
			if v := pc.all[uri]; v != nil {
				*ref.val = v
			} else {
				result = result.Merge(pc.loadSchemaFromURI(ref.ctx, ref.js, uri))
			}
		}
	}
	return
}

func (pc *ptrCtx) loadSchemaFromURI(ctx *schema.Context, js jsi.JSON, uri string) *schema.Result {
	u, err := url.Parse(uri)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri",
		})
	}
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

	_, result := pc.withId(u).cmpRoot(ctx.WithPath(uri), subjs)
	return result
}
