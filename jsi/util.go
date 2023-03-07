package jsi

import (
	"math/big"
	"unsafe"
)

func RootOf(js JSON) JSON {
	if js.Parent() == nil {
		return js
	}
	return RootOf(js.Parent())
}

func SiblingOf(js JSON, key string) JSON {
	p := js.Parent()
	if p == nil {
		return nil
	}
	if p.Type() != TypeObject {
		return nil
	}
	return p.(Object).Index(key)
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Equal(a, b JSON) bool {
	if a.Type() != b.Type() {
		return false
	}

	switch a.Type() {
	case TypeObject:
		oa := a.(Object)
		ob := b.(Object)
		if oa.Len() != ob.Len() {
			return false
		}
		iter := oa.Iter()
		for iter.Next() {
			key, x := iter.Entry()
			y := ob.Index(key)
			if y == nil {
				return false
			}
			if !Equal(x, y) {
				return false
			}
		}
		return true

	case TypeArray:
		aa := a.(Array)
		ab := b.(Array)
		if aa.Len() != ab.Len() {
			return false
		}
		for i := 0; i < aa.Len(); i++ {
			if !Equal(aa.Index(i), ab.Index(i)) {
				return false
			}
		}
		return true

	case TypeString:
		return a.(String).Value() == b.(String).Value()

	case TypeNumber:
		fa, ok1 := new(big.Float).SetString(a.(Number).Value().String())
		fb, ok2 := new(big.Float).SetString(b.(Number).Value().String())
		if !ok1 || !ok2 {
			return false
		}
		return fa.Cmp(fb) == 0

	case TypeBoolean:
		return a.(Boolean).Value() == b.(Boolean).Value()

	case TypeNULL:
		return true
	}

	return false
}
