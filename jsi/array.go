package jsi

import "encoding/json"

type jsArray struct {
	p JSON
	l []JSON
	r json.RawMessage
}

func (*jsArray) Type() Type {
	return TypeArray
}

func (a *jsArray) Parent() JSON {
	return a.p
}

func (a *jsArray) String() string {
	a.MarshalJSON()
	return bytes2str(a.r)
}

func (a *jsArray) Len() int {
	return len(a.l)
}

func (a *jsArray) Index(i int) JSON {
	return a.l[i]
}

func (a *jsArray) MarshalJSON() ([]byte, error) {
	if len(a.r) > 0 {
		return a.r, nil
	}
	b, err := json.Marshal(a.l)
	if err != nil {
		return nil, err
	}
	a.r = b
	return b, nil
}
