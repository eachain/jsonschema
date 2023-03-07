package jsi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type jsParser struct {
	dec *json.Decoder
}

type SyntaxError struct {
	msg    string // description of error
	Offset int64  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return e.msg }

func (p jsParser) Parse() (JSON, error) {
	t, err := p.dec.Token()
	if err == io.EOF {
		return nil, &SyntaxError{"unexpected end of JSON input", p.dec.InputOffset()}
	}
	return p.parse(nil, t)
}

func (p jsParser) parse(parent JSON, t json.Token) (JSON, error) {
	switch v := t.(type) {
	case json.Delim:
		switch v {
		case '{':
			return p.parseObject(parent)
		case '[':
			return p.parseArray(parent)
		default:
			return nil, &SyntaxError{fmt.Sprintf("invalid character '%v'", v), p.dec.InputOffset()}
		}
	case bool:
		return &jsBoolean{p: parent, v: v}, nil
	case json.Number:
		return &jsNumber{p: parent, n: v}, nil
	case string:
		return &jsString{p: parent, s: v}, nil
	case nil:
		return &jsNULL{p: parent}, nil
	default:
		return nil, &SyntaxError{fmt.Sprintf("invalid character '%v'", v), p.dec.InputOffset()}
	}
}

func (p jsParser) parseObject(parent JSON) (JSON, error) {
	obj := &jsObject{
		p: parent,
		m: make(map[string]JSON),
	}

	for p.dec.More() {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, err
		}
		key, ok := tok.(string)
		if !ok {
			return nil, &SyntaxError{fmt.Sprintf("invalid object key type '%T'", tok), p.dec.InputOffset()}
		}
		if _, ok = obj.m[key]; ok {
			return nil, &SyntaxError{fmt.Sprintf("duplicate object key '%v'", key), p.dec.InputOffset()}
		}

		tok, err = p.dec.Token()
		if err != nil {
			return nil, err
		}
		val, err := p.parse(obj, tok)
		if err != nil {
			return nil, err
		}
		obj.k = append(obj.k, key)
		obj.m[key] = val
	}

	tok, err := p.dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil, &SyntaxError{fmt.Sprintf("invalid token after object '%v'", tok), p.dec.InputOffset()}
	}
	if delim != '}' {
		return nil, &SyntaxError{fmt.Sprintf("invalid character '%v' after object", tok), p.dec.InputOffset()}
	}
	return obj, nil
}

func (p jsParser) parseArray(parent JSON) (JSON, error) {
	arr := &jsArray{p: parent}

	for p.dec.More() {
		tok, err := p.dec.Token()
		if err != nil {
			return nil, err
		}
		val, err := p.parse(arr, tok)
		if err != nil {
			return nil, err
		}
		arr.l = append(arr.l, val)
	}

	tok, err := p.dec.Token()
	if err != nil {
		return nil, err
	}
	delim, ok := tok.(json.Delim)
	if !ok {
		return nil, &SyntaxError{fmt.Sprintf("invalid token after array '%v'", tok), p.dec.InputOffset()}
	}
	if delim != ']' {
		return nil, &SyntaxError{fmt.Sprintf("invalid character '%v' after array", tok), p.dec.InputOffset()}
	}
	return arr, nil
}

func NewReaderParser(r io.Reader) Parser {
	p := jsParser{
		dec: json.NewDecoder(r),
	}
	p.dec.UseNumber()
	return p
}

func NewBytesParser(p []byte) Parser {
	return NewReaderParser(bytes.NewReader(p))
}

type goTypesParser struct {
	x any
}

func NewGoTypesParser(x any) Parser {
	return goTypesParser{x: x}
}

func (p goTypesParser) Parse() (JSON, error) {
	b := new(bytes.Buffer)
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	err := enc.Encode(p.x)
	if err != nil {
		return nil, err
	}
	return NewBytesParser(b.Bytes()).Parse()
}

type errorParser struct {
	err error
}

func (s errorParser) Parse() (JSON, error) {
	return nil, s.err
}

func NewURLParser(u *url.URL) Parser {
	switch u.Scheme {
	case "http", "https":
		return NewHTTPParser(u.String())
	case "file":
		return NewFileParser(u.Path)
	default:
		return errorParser{fmt.Errorf("unsupport url scheme '%v'", u.Scheme)}
	}
}

type readCloserParser struct {
	rc io.ReadCloser
	p  Parser
}

func NewReadCloserParser(rc io.ReadCloser) Parser {
	return readCloserParser{rc: rc, p: NewReaderParser(rc)}
}

func (p readCloserParser) Parse() (JSON, error) {
	defer func() {
		io.Copy(io.Discard, p.rc)
		p.rc.Close()
	}()
	return p.p.Parse()
}

func NewFileParser(path string) Parser {
	fp, err := os.Open(path)
	if err != nil {
		return errorParser{err}
	}
	return NewReadCloserParser(fp)
}

func NewHTTPParser(url string) Parser {
	resp, err := http.Get(url)
	if err != nil {
		return errorParser{err}
	}
	if resp.StatusCode != http.StatusOK {
		io.ReadAll(resp.Body)
		resp.Body.Close()
		return errorParser{fmt.Errorf("HTTP GET %v: %v", url, resp.Status)}
	}
	return NewReadCloserParser(resp.Body)
}
