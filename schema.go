package jsonschema

import (
	"github.com/eachain/jsonschema/jsi"
)

type Schema struct {
	draft string
	val   Validator
}

func (s *Schema) Draft() string {
	return s.draft
}

func (s *Schema) Validate(js jsi.JSON) *Result {
	if s.val == nil {
		return nil
	}
	return s.val.Validate(newContext(s.draft), js)
}

func Compile(js jsi.JSON, drafts ...string) (*Schema, *Result) {
	if len(drafts) == 0 {
		drafts = supportDrafts()
	}
	if len(drafts) == 0 {
		return nil, &Result{Errors: []Error{{
			Field: "",
			Type:  js.Type(),
			Value: js,
			Msg:   MsgNoDraftRegister,
		}}}
	}

	draft := drafts[0]
	val, result := compile(drafts[0], js)
	warns := 0
	if result != nil {
		warns = len(result.Warnings)
	}
	if result.Valid() && warns == 0 {
		return &Schema{draft: draft, val: val}, nil
	}

	for i := 1; i < len(drafts); i++ {
		v, r := compile(drafts[i], js)
		if !r.Valid() {
			continue
		}
		ws := 0
		if r != nil {
			ws = len(r.Warnings)
		}

		if !result.Valid() || ws < warns {
			draft = drafts[i]
			val = v
			result = r
		}
		if ws == 0 {
			break
		}
	}
	return &Schema{draft: draft, val: val}, result
}

func compile(draft string, js jsi.JSON) (Validator, *Result) {
	root := GetKeyword(draft, RootKeyword)
	if root == nil {
		return nil, &Result{Errors: []Error{{
			Field: ".",
			Type:  js.Type(),
			Value: js,
			Msg:   MsgRootCompilerNotFound,
		}}}
	}
	ctx := newContext(draft)
	return root.Compile(ctx, js)
}
