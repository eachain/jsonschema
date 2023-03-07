package jsi

import (
	"bytes"
	"encoding/json"
)

type jsObject struct {
	p JSON
	m map[string]JSON
	r json.RawMessage
	k []string
}

type jsObjectIter struct {
	o *jsObject
	i int
}

func (*jsObject) Type() Type {
	return TypeObject
}

func (o *jsObject) Parent() JSON {
	return o.p
}

func (o *jsObject) String() string {
	o.MarshalJSON()
	return bytes2str(o.r)
}

func (o *jsObject) Len() int {
	return len(o.k)
}

func (o *jsObject) Index(key string) JSON {
	return o.m[key]
}

func (o *jsObject) MarshalJSON() ([]byte, error) {
	if len(o.r) > 0 {
		return o.r, nil
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)

	b.WriteByte('{')
	for i, k := range o.k {
		if i > 0 {
			b.WriteByte(',')
		}
		err := enc.Encode(k)
		if err != nil {
			return nil, err
		}
		n := b.Len() - 1
		if b.Bytes()[n] == '\n' {
			b.Truncate(n)
		}
		b.WriteByte(':')
		err = enc.Encode(o.m[k])
		if err != nil {
			return nil, err
		}
		if b.Bytes()[n] == '\n' {
			b.Truncate(n)
		}
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func (o *jsObject) Iter() ObjectIter {
	return &jsObjectIter{o: o, i: -1}
}

func (i *jsObjectIter) Next() bool {
	i.i++
	return i.i < len(i.o.k)
}

func (i *jsObjectIter) Entry() (key string, val JSON) {
	key = i.o.k[i.i]
	val = i.o.m[key]
	return
}
