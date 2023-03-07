package jsi

import "encoding/json"

type Type = string

const (
	TypeObject  Type = "object"
	TypeArray   Type = "array"
	TypeString  Type = "string"
	TypeNumber  Type = "number"
	TypeBoolean Type = "boolean"
	TypeNULL    Type = "null"
)

type JSON interface {
	Type() Type
	Parent() JSON
}

type ObjectIter interface {
	Next() bool
	Entry() (key string, val JSON)
}

type Object interface {
	Iter() ObjectIter
	Index(key string) JSON
	Len() int
}

type Array interface {
	Len() int
	Index(i int) JSON
}

type String interface {
	Value() string
}

type Number interface {
	Value() json.Number
}

type Boolean interface {
	Value() bool
}

type NULL interface {
	Value()
}

type Parser interface {
	Parse() (JSON, error)
}
