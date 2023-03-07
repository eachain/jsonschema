package jsi

type jsBoolean struct {
	p JSON
	v bool
}

func (*jsBoolean) Type() Type {
	return TypeBoolean
}

func (b *jsBoolean) Parent() JSON {
	return b.p
}

func (b *jsBoolean) Value() bool {
	return b.v
}

func (b *jsBoolean) String() string {
	if b.v {
		return "true"
	}
	return "false"
}

func (b *jsBoolean) MarshalJSON() ([]byte, error) {
	if b.v {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}
