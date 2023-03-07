package jsonschema

import (
	"net/url"
	"path"
	"strconv"
	"strings"
)

type Pointer struct {
	Scheme string
	Host   string
	Path   string
	Frag   []string // frag path with prefix slash, eg. {"/$def", "/abc"}
}

func ParsePointer(id string) (*Pointer, error) {
	noPath := false
	if strings.HasPrefix(id, "#") {
		noPath = true
		id = "t" + id
	}
	u, err := url.Parse(id)
	if err != nil {
		return nil, err
	}
	if noPath {
		u.Path = ""
	}
	return &Pointer{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   u.Path,
		Frag:   parseFrag(u.Fragment),
	}, nil
}

func parseFrag(frag string) (ps []string) {
	var rs []rune
	slash := false
	for i, r := range frag {
		if slash {
			slash = false
			rs = append(rs, r)
			continue
		}

		if r == '\\' {
			slash = true
			continue
		}

		if r == '/' && i > 0 {
			if len(rs) == 1 && rs[0] == '/' {
			} else {
				ps = append(ps, string(rs))
			}
			rs = rs[:0]
		}

		rs = append(rs, r)
	}

	if len(rs) > 0 {
		if len(rs) == 1 && rs[0] == '/' {
		} else {
			ps = append(ps, string(rs))
		}
	}
	return ps
}

func (p *Pointer) clone() *Pointer {
	if p == nil {
		return nil
	}
	q := *p
	q.Frag = make([]string, len(p.Frag))
	copy(q.Frag, p.Frag)
	return &q
}

func (p *Pointer) Fix(q *Pointer) *Pointer {
	if p == nil {
		return q
	}
	if q == nil {
		return p.clone()
	}
	if q.Scheme != "" {
		return q
	}
	q.Scheme = p.Scheme
	if q.Host != "" {
		return q
	}
	q.Host = p.Host
	if q.Path != "" {
		qIsRoot := strings.HasPrefix(q.Path, "/")
		if qIsRoot {
			return q
		}
		pIsDir := strings.HasSuffix(p.Path, "/")
		qIsDir := strings.HasSuffix(q.Path, "/")
		if pIsDir {
			q.Path = path.Clean(path.Join(p.Path, q.Path))
		} else {
			q.Path = path.Clean(path.Join(path.Dir(p.Path), q.Path))
		}
		if qIsDir {
			q.Path += "/"
		}
		return q
	}
	q.Path = p.Path
	return q
}

func (p *Pointer) Array(i int) *Pointer {
	q := p.clone()
	q.Frag = append(q.Frag, "/"+strconv.Itoa(i))
	return q
}

func (p *Pointer) Object(key string) *Pointer {
	q := p.clone()
	if strings.HasPrefix(key, "/") {
		q.Frag = append(q.Frag, escape(key))
	} else {
		q.Frag = append(q.Frag, "/"+escape(key))
	}
	return q
}

func (p *Pointer) escapedIndex(key string) *Pointer {
	q := p.clone()
	if strings.HasPrefix(key, "/") {
		q.Frag = append(q.Frag, key)
	} else {
		q.Frag = append(q.Frag, "/"+key)
	}
	return q
}

func (p *Pointer) URL() *url.URL {
	return &url.URL{
		Scheme:   p.Scheme,
		Host:     p.Host,
		Path:     p.Path,
		Fragment: strings.Join(p.Frag, ""),
	}
}

func (p *Pointer) String() string {
	if p == nil {
		return "#"
	}
	b := new(strings.Builder)
	if p.Scheme != "" {
		b.WriteString(p.Scheme)
		b.WriteString("://")
	}
	if p.Host != "" {
		b.WriteString(p.Host)
	}
	if p.Path != "" {
		b.WriteString(p.Path)
	}
	b.WriteByte('#')
	for _, f := range p.Frag {
		b.WriteString(f)
	}
	return b.String()
}

func escape(f string) string {
	var escaped bool
	var rs []rune
	for i, r := range f {
		if r == '~' {
			rs = append(rs, '~', '0')
			escaped = true
		} else if r == '/' && i > 0 {
			rs = append(rs, '~', '1')
			escaped = true
		} else if r == '%' {
			rs = append(rs, '%', '2', '5')
		} else {
			rs = append(rs, r)
		}
	}
	if !escaped {
		return f
	}
	return string(rs)
}
