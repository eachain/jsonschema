package jsonschema

type subSchema struct {
	path      *Pointer
	anchor    string
	validator Validator

	sub    map[string]*subSchema
	parent *subSchema
}

func newSubSchema(draft string) *subSchema {
	return &subSchema{
		path: new(Pointer),
		sub:  make(map[string]*subSchema),
	}
}

func (s *subSchema) loadOrNew(index string) *subSchema {
	m := s.sub[index]
	if m == nil {
		m = &subSchema{
			path:   s.path.Object(index),
			parent: s,
		}
		if s.sub == nil {
			s.sub = make(map[string]*subSchema)
		}
		s.sub[index] = m
	}
	return m
}

func (s *subSchema) setId(id *Pointer) {
	s.path = s.path.Fix(id.clone())
	s.path.Frag = nil
}

/*
func (s *subSchema) Print(prefix string) {
	fmt.Printf("%vpath: %v\n", prefix, s.path)
	fmt.Printf("%vanchor: %v\n", prefix, s.anchor)
	fmt.Printf("%vvalidator: %v\n", prefix, s.validator != nil)
	fmt.Printf("%vsub: {\n", prefix)
	for idx, sub := range s.sub {
		fmt.Printf("%v    %q: {\n", prefix, idx)
		sub.Print(prefix + "        ")
		fmt.Printf("%v    }\n", prefix)
	}
	fmt.Printf("%v}\n", prefix)
}
*/

func (s *subSchema) visit(impl map[string]Validator, id *Pointer) {
	if s.validator != nil {
		impl[id.String()] = s.validator
	}
	if len(s.sub) == 0 {
		return
	}

	for idx, sub := range s.sub {
		sub.visit(impl, id.escapedIndex(idx))
	}

	if s.anchor != "" {
		anchor := id.escapedIndex(s.anchor)
		if s.validator != nil {
			impl[anchor.String()] = s.validator
		}
		for idx, sub := range s.sub {
			sub.visit(impl, anchor.escapedIndex(idx))
		}
	}

	if s.path != nil {
		id2 := id.Fix(s.path.clone())
		if id2.String() == id.String() {
			return
		}
		if s.validator != nil {
			impl[id2.String()] = s.validator
		}
		for idx, sub := range s.sub {
			sub.visit(impl, id2.escapedIndex(idx))
		}
	}
}
