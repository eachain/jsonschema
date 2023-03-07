package jsi

import (
	"bytes"
	"encoding/json"
)

type jsString struct {
	p JSON
	s string
	r json.RawMessage
}

func (*jsString) Type() Type {
	return TypeString
}

func (s *jsString) Parent() JSON {
	return s.p
}

func (s *jsString) Value() string {
	return s.s
}

func (s *jsString) String() string {
	s.MarshalJSON()
	return bytes2str(s.r)
}

func (s *jsString) MarshalJSON() ([]byte, error) {
	if len(s.r) > 0 {
		return s.r, nil
	}

	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(s.s)
	if err != nil {
		return nil, err
	}
	s.r = json.RawMessage(b.Bytes())
	n := len(s.r) - 1
	if n >= 0 && s.r[n] == '\n' {
		s.r = s.r[:n]
	}
	return s.r, nil
}
