package jsonschema

import (
	"fmt"
	"strings"
)

type Type = string

const (
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeString  Type = "string"
	TypeNumber  Type = "number"
	TypeInteger Type = "integer"
	TypeBoolean Type = "boolean"
	TypeNULL    Type = "null"
)

var AllTypes = []Type{
	TypeObject,
	TypeArray,
	TypeString,
	TypeNumber,
	TypeInteger,
	TypeBoolean,
	TypeNULL,
}

var (
	MsgNoDraftRegister      = "no draft register"
	MsgRootCompilerNotFound = "root compiler not found"
)

type Error struct {
	Field string
	Type  Type
	Value interface{}
	Msg   string
}

func (e Error) Error() string {
	if e.Field == "" {
		e.Field = "."
	}

	val := fmt.Sprint(e.Value)
	if val == "" {
		return fmt.Sprintf("path %v type '%v': %v", e.Field, e.Type, e.Msg)
	}
	return fmt.Sprintf("path %v type '%v': %v, value: %v",
		e.Field, e.Type, e.Msg, e.Value)
}

type Result struct {
	Warnings []Error
	Errors   []Error
}

func (r *Result) Valid() bool {
	return r == nil || len(r.Errors) == 0
}

func (r *Result) Warned() bool {
	return r != nil && len(r.Warnings) > 0
}

func (r *Result) Error() string {
	if r.Valid() {
		return ""
	}
	errs := make([]string, len(r.Errors))
	for i, e := range r.Errors {
		errs[i] = e.Error()
	}
	return strings.Join(errs, "; ")
}

func (r *Result) Warning() string {
	if r == nil || len(r.Warnings) == 0 {
		return ""
	}
	errs := make([]string, len(r.Warnings))
	for i, e := range r.Warnings {
		errs[i] = e.Error()
	}
	return strings.Join(errs, "; ")
}

func WithWarning(e Error) *Result {
	return &Result{Warnings: []Error{e}}
}

func WithError(e Error) *Result {
	return &Result{Errors: []Error{e}}
}

func (r *Result) WithWarning(e Error) *Result {
	if r == nil {
		return &Result{Warnings: []Error{e}}
	}
	r.Warnings = append(r.Warnings, e)
	return r
}

func (r *Result) WithError(e Error) *Result {
	if r == nil {
		return &Result{Errors: []Error{e}}
	}
	r.Errors = append(r.Errors, e)
	return r
}

func (r *Result) Merge(t *Result) *Result {
	if t == nil {
		return r
	}
	if r == nil {
		return t
	}
	if len(r.Warnings) == 0 {
		r.Warnings = t.Warnings
	} else {
		if len(t.Warnings) > 0 {
			r.Warnings = append(r.Warnings, t.Warnings...)
		}
	}

	if len(r.Errors) == 0 {
		r.Errors = t.Errors
	} else {
		if len(t.Errors) > 0 {
			r.Errors = append(r.Errors, t.Errors...)
		}
	}
	return r
}
