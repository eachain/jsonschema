package jsonschema

import (
	"strconv"

	"github.com/eachain/jsonschema/jsi"
)

type Context struct {
	draft  string
	schema *subSchema
	id     *Pointer
	path   *Pointer
	index  string
	impl   map[string]Validator
	ref    map[string][]*Reference
	parent *Context
}

type Reference struct {
	Ref *Pointer
	*Context
	jsi.JSON
	Validator *Validator
}

func newContext(draft string) *Context {
	return &Context{
		draft:  draft,
		schema: newSubSchema(draft),
		id:     new(Pointer),
		path:   new(Pointer),
		impl:   make(map[string]Validator),
		ref:    make(map[string][]*Reference),
	}
}

func (ctx *Context) clone() *Context {
	return &Context{
		draft:  ctx.draft,
		schema: ctx.schema,
		id:     ctx.id.clone(),
		path:   ctx.path.clone(),
		index:  ctx.index,
		impl:   ctx.impl,
		ref:    ctx.ref,
		parent: ctx,
	}
}

func (ctx *Context) UseDraft(draft string) {
	ctx.draft = draft
}

func (ctx *Context) Draft() string {
	return ctx.draft
}

func (ctx *Context) Id() string {
	return ctx.id.String()
}

func (ctx *Context) Field() string {
	return ctx.path.String()
}

func (ctx *Context) Parent() *Context {
	if ctx == nil {
		return nil
	}
	return ctx.parent
}

func (ctx *Context) SetId(id *Pointer) {
	ctx.schema.setId(id)
}

func (ctx *Context) SetAnchor(anchor string) {
	ctx.schema.anchor = anchor
}

func (ctx *Context) Index() string {
	if ctx == nil {
		return ""
	}
	return ctx.index
}

func (ctx *Context) Array(i int) *Context {
	sc := ctx.clone()
	sc.schema = ctx.schema.loadOrNew(strconv.Itoa(i))
	sc.path = ctx.path.Array(i)
	sc.index = strconv.Itoa(i)
	return sc
}

func (ctx *Context) Object(key string) *Context {
	key = escape(key)
	sc := ctx.clone()
	sc.schema = ctx.schema.loadOrNew(key)
	sc.path = ctx.path.Object(key)
	sc.index = key
	return sc
}

func (ctx *Context) WithId(id *Pointer) *Context {
	sc := ctx.clone()
	sc.id = sc.id.Fix(id.clone())
	sc.id.Frag = nil
	sc.path = sc.id.clone()
	sc.schema.setId(id)
	return sc
}

func (ctx *Context) PushRef(js jsi.JSON, ref *Pointer, val *Validator) {
	ref = ctx.schema.path.Fix(ref.clone())
	r := ref.String()
	ctx.ref[r] = append(ctx.ref[r], &Reference{
		Ref:       ref,
		Context:   ctx,
		JSON:      js,
		Validator: val,
	})
}

func (ctx *Context) Impl(val Validator) {
	ctx.schema.validator = val
}

func (ctx *Context) FillRefs() []*Reference {
	ctx.spread()
	var refs []*Reference
	for uri, rs := range ctx.ref {
		impl := ctx.impl[uri]
		if impl == nil {
			refs = append(refs, rs...)
			continue
		}

		for _, v := range rs {
			*v.Validator = impl
		}
	}
	return refs
}

func (ctx *Context) spread() {
	ctx.schema.visit(ctx.impl, ctx.id)
}
