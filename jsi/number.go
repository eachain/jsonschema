package jsi

import "encoding/json"

type jsNumber struct {
	p JSON
	n json.Number
}

func (*jsNumber) Type() Type {
	return TypeNumber
}

func (n *jsNumber) Parent() JSON {
	return n.p
}

func (n *jsNumber) Value() json.Number {
	return n.n
}

func (n *jsNumber) String() string {
	return string(n.n)
}

func (n *jsNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.n)
}
