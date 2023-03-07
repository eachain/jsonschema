package basic

import (
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"

	schema "github.com/eachain/jsonschema"
	"github.com/eachain/jsonschema/jsi"
)

type FormatOf map[string]schema.Validator

var DefaultFormatOf = FormatOf{
	"date-time": schema.ValidateFunc(DateTimeFormatValidator),
	"email":     schema.ValidateFunc(EmailFormatValidator),
	"hostname":  schema.ValidateFunc(HostnameFormatValidator),
	"ipv4":      schema.ValidateFunc(IPv4FormatValidator),
	"ipv6":      schema.ValidateFunc(IPv6FormatValidator),
	"uri":       schema.ValidateFunc(URIFormatValidator),
	"regex":     schema.ValidateFunc(RegexFormatValidator),
}

func GenFormat(formatOf FormatOf) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		if js.Type() != jsi.TypeString {
			return nil, schema.WithError(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "should be string",
			})
		}
		format := js.(jsi.String).Value()
		validator := formatOf[format]
		if validator == nil {
			return nil, schema.WithWarning(schema.Error{
				Field: ctx.Field(),
				Type:  js.Type(),
				Value: js,
				Msg:   "format '" + format + "' not found",
			})
		}

		return &FormatValidator{
			Format:    format,
			Validator: validator,
		}, nil
	}
}

func Format(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
	return GenFormat(DefaultFormatOf)(ctx, js)
}

func ValidateFormat(cmp schema.Compiler) schema.CompileFunc {
	return func(ctx *schema.Context, js jsi.JSON) (schema.Validator, *schema.Result) {
		val, result := cmp.Compile(ctx, js)
		err := checkSiblingOfType(ctx, js, jsi.TypeString)
		if err != nil {
			return val, result.WithWarning(*err)
		}
		return val, result
	}
}

type FormatValidator struct {
	Format    string
	Validator schema.Validator
}

func (m *FormatValidator) Validate(ctx *schema.Context, js jsi.JSON) *schema.Result {
	return m.Validator.Validate(ctx, js)
}

// formats

func DateTimeFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	var format string
	val := js.(jsi.String).Value()
	switch len(val) {
	case 8:
		format = "15:04:05"
	case 10:
		format = "2006-01-02"
	case 11:
		format = "15:04:05Z07"
	case 14:
		format = "15:04:05Z07:00"
	case 19:
		format = "2006-01-02T15:04:05"
	case 22:
		format = "2006-01-02T15:04:05Z07"
	case len(time.RFC3339): // "2006-01-02T15:04:05Z07:00"
		format = time.RFC3339
	default: // "2006-01-02T15:04:05.999999999Z07:00"
		if zone := strings.LastIndexByte(val, '+'); zone > 0 {
			if colon := strings.LastIndexByte(val, ':'); colon > zone {
				format = time.RFC3339Nano
			} else {
				format = "2006-01-02T15:04:05.999999999Z07"
			}
		} else {
			format = "2006-01-02T15:04:05.999999999"
		}
	}
	_, err := time.Parse(format, val)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be date-time format",
		})
	}
	return nil
}

func DateFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	format := "2006-01-02"
	val := js.(jsi.String).Value()
	_, err := time.Parse(format, val)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be date format",
		})
	}
	return nil
}

func TimeFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	var format string
	val := js.(jsi.String).Value()
	switch len(val) {
	case 8:
		format = "15:04:05"
	case 11:
		format = "15:04:05Z07"
	case 14:
		format = "15:04:05Z07:00"
	}
	var err error
	if format != "" {
		_, err = time.Parse(format, val)
	}
	if format == "" || err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be time format",
		})
	}
	return nil
}

func EmailFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	_, err := mail.ParseAddress(js.(jsi.String).Value())
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be email format",
		})
	}
	return nil
}

var regexpHostname = regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])(\.([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}[a-zA-Z0-9]))*$`)

func HostnameFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	if !regexpHostname.MatchString(val) {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be hostname format",
		})
	}
	return nil
}

func IPv4FormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	if idx := strings.IndexByte(val, '.'); idx <= 0 {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be ipv4 format",
		})
	}
	ip := net.ParseIP(val)
	if ip == nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be ipv4 format",
		})
	}
	return nil
}

func IPv6FormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	if idx := strings.IndexByte(val, ':'); idx < 0 {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be ipv6 format",
		})
	}
	ip := net.ParseIP(val)
	if ip == nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be ipv6 format",
		})
	}
	return nil
}

func URIFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	if strings.IndexByte(val, '\\') >= 0 {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri format",
		})
	}

	u, err := url.Parse(val)
	if err != nil || u.Scheme == "" {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri format",
		})
	}
	return nil
}

func URIReferenceFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	if strings.IndexByte(val, '\\') >= 0 {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri format",
		})
	}

	_, err := url.Parse(val)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri format",
		})
	}
	return nil
}

func RegexFormatValidator(ctx *schema.Context, js jsi.JSON) *schema.Result {
	if js.Type() != jsi.TypeString {
		return nil
	}
	val := js.(jsi.String).Value()
	_, err := regexp.Compile(val)
	if err != nil {
		return schema.WithError(schema.Error{
			Field: ctx.Field(),
			Type:  js.Type(),
			Value: js,
			Msg:   "should be uri format",
		})
	}
	return nil
}
