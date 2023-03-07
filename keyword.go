package jsonschema

import (
	"sync"

	"github.com/eachain/jsonschema/jsi"
)

const RootKeyword = ""

type Validator interface {
	Validate(ctx *Context, js jsi.JSON) *Result
}

type ValidateFunc func(ctx *Context, js jsi.JSON) *Result

func (fn ValidateFunc) Validate(ctx *Context, js jsi.JSON) *Result {
	return fn(ctx, js)
}

type Compiler interface {
	Compile(ctx *Context, js jsi.JSON) (Validator, *Result)
}

type CompileFunc func(ctx *Context, js jsi.JSON) (Validator, *Result)

func (fn CompileFunc) Compile(ctx *Context, js jsi.JSON) (Validator, *Result) {
	return fn(ctx, js)
}

var (
	keywordsMu sync.RWMutex
	keywordsOf map[string]map[string]Compiler
)

func RegisterKeyword(draft string, name string, cmp Compiler) {
	keywordsMu.Lock()
	defer keywordsMu.Unlock()

	keywords := keywordsOf[draft]
	if keywords == nil {
		keywords = make(map[string]Compiler)
		if keywordsOf == nil {
			keywordsOf = make(map[string]map[string]Compiler)
		}
		keywordsOf[draft] = keywords
	}

	keywords[name] = cmp
}

func supportDrafts() []string {
	keywordsMu.RLock()
	defer keywordsMu.RUnlock()
	var drafts []string
	for draft := range keywordsOf {
		drafts = append(drafts, draft)
	}
	return drafts
}

func GetKeyword(draft string, keyword string) Compiler {
	keywordsMu.RLock()
	defer keywordsMu.RUnlock()
	return keywordsOf[draft][keyword]
}
