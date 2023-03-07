package jsi

type jsNULL struct {
	p JSON
}

func (*jsNULL) Type() Type {
	return TypeNULL
}

func (n *jsNULL) Parent() JSON {
	return n.p
}

func (*jsNULL) Value() {}

func (*jsNULL) String() string {
	return "null"
}

func (*jsNULL) MarshalJSON() ([]byte, error) {
	return []byte("null"), nil
}
