package vld

import (
	"errors"
	"fmt"
	"github.com/zzztttkkk/ink/internal/utils"
	"html"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type RuleType int

const (
	RuleTypeInt = RuleType(iota)
	RuleTypeDouble
	RuleTypeBool
	RuleTypeString
	RuleTypeFile
	RuleTypeTime
	RuleTypeBinder
)

type Rule struct {
	Name     string
	RuleType RuleType
	Gotype   reflect.Type
	IsSlice  bool
	Index    []int

	Optional bool

	MaxInt    *int64
	MinInt    *int64
	MaxDouble *float64
	MinDouble *float64

	MaxLen *int
	MinLen *int

	MaxRuneCount *int
	MinRuneCount *int
	NoTrim       bool
	NoEscape     bool
	Regexp       *regexp.Regexp

	TimeLayout string
	TimeUnit   string
}

func (rule *Rule) intOk(num int64) ErrorType {
	if rule.MinInt != nil && num < *rule.MinInt {
		return ErrorTypeNumOutOfRange
	}
	if rule.MaxInt != nil && num > *rule.MaxInt {
		return ErrorTypeNumOutOfRange
	}
	return ErrorTypeNil
}

func (rule *Rule) floatOk(num float64) ErrorType {
	if rule.MinDouble != nil && num < *rule.MinDouble {
		return ErrorTypeNumOutOfRange
	}
	if rule.MaxDouble != nil && num > *rule.MaxDouble {
		return ErrorTypeNumOutOfRange
	}
	return ErrorTypeNil
}

func (rule *Rule) string2Int(v string) (any, ErrorType) {
	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, ErrorTypeCanNotCastToNum
	}

	if er := rule.intOk(num); er != ErrorTypeNil {
		return nil, er
	}

	switch rule.Gotype.Kind() {
	case reflect.Int:
		return int(num), ErrorTypeNil
	case reflect.Int8:
		return int8(num), ErrorTypeNil
	case reflect.Int16:
		return int16(num), ErrorTypeNil
	case reflect.Int32:
		return int32(num), ErrorTypeNil
	case reflect.Int64:
		return num, ErrorTypeNil
	case reflect.Uint:
		return uint(num), ErrorTypeNil
	case reflect.Uint8:
		return uint8(num), ErrorTypeNil
	case reflect.Uint16:
		return uint16(num), ErrorTypeNil
	case reflect.Uint32:
		return uint32(num), ErrorTypeNil
	case reflect.Uint64:
		return uint64(num), ErrorTypeNil
	}
	return num, ErrorTypeNil
}

func (rule *Rule) string2Double(v string) (any, ErrorType) {
	num, err := strconv.ParseFloat(v, 10)
	if err != nil {
		return nil, ErrorTypeCanNotCastToNum
	}

	if er := rule.floatOk(num); er != ErrorTypeNil {
		return nil, er
	}

	switch rule.Gotype.Kind() {
	case reflect.Float32:
		return float32(num), ErrorTypeNil
	default:
		return num, ErrorTypeNil
	}
}

func (rule *Rule) string2Time(v string) (any, ErrorType) {
	if len(rule.TimeLayout) > 0 {
		t, e := time.Parse(rule.TimeLayout, v)
		if e != nil {
			return nil, ErrorTypeBadTimeValue
		}
		return t, ErrorTypeNil
	}

	num, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return nil, ErrorTypeNumOutOfRange
	}

	switch rule.TimeUnit {
	case "", "s":
		return time.Unix(num, 0), ErrorTypeNil
	default:
		return time.UnixMilli(num), ErrorTypeNil
	}
}

func (rule *Rule) stringOk(v string) ErrorType {
	var runeCount = -1
	if rule.MaxRuneCount != nil {
		runeCount = utf8.RuneCount(utils.B(v))
		if runeCount > *rule.MaxRuneCount {
			return ErrorTypeLengthOutOfRange
		}
	}

	if rule.MinRuneCount != nil {
		if runeCount < 0 {
			runeCount = utf8.RuneCount(utils.B(v))
		}
		if runeCount < *rule.MinRuneCount {
			return ErrorTypeLengthOutOfRange
		}
	}

	if rule.Regexp != nil && !rule.Regexp.Match(utils.B(v)) {
		return ErrorTypeNotMatchRegexp
	}
	return ErrorTypeNil
}

func (rule *Rule) sliceLenOk(lenv int) ErrorType {
	if rule.MinLen != nil && lenv < *rule.MinLen {
		return ErrorTypeLengthOutOfRange
	}
	if rule.MaxLen != nil && lenv > *rule.MaxLen {
		return ErrorTypeLengthOutOfRange
	}
	return ErrorTypeNil
}

func (rule *Rule) fileOk(f *multipart.FileHeader) ErrorType {
	if rule.MinLen != nil && int(f.Size) < *rule.MinLen {
		return ErrorTypeLengthOutOfRange
	}
	if rule.MaxLen != nil && int(f.Size) > *rule.MaxLen {
		return ErrorTypeLengthOutOfRange
	}
	return ErrorTypeNil
}

func (rule *Rule) one(raw string) (any, ErrorType) {
	switch rule.RuleType {
	case RuleTypeString:
		{
			if !rule.NoTrim {
				raw = strings.TrimSpace(raw)
			}

			if !rule.NoEscape {
				raw = html.EscapeString(raw)
			}

			if et := rule.stringOk(raw); et != ErrorTypeNil {
				return raw, et
			}
			return raw, ErrorTypeNil
		}
	case RuleTypeInt:
		{
			return rule.string2Int(raw)
		}
	case RuleTypeDouble:
		{
			return rule.string2Double(raw)
		}
	case RuleTypeBool:
		{
			bol, err := strconv.ParseBool(raw)
			if err != nil {
				return nil, ErrorTypeCanNotCastToBool
			}
			return bol, ErrorTypeNil
		}
	case RuleTypeTime:
		{
			return rule.string2Time(raw)
		}
	}
	panic(errors.New("unreachable error"))
}

func getFiles(req *http.Request, name string) []*multipart.FileHeader {
	_ = req.ParseForm()
	if req.MultipartForm == nil || req.MultipartForm.File == nil {
		return nil
	}
	return req.MultipartForm.File[name]
}

func (rule *Rule) get(req *http.Request) (any, ErrorType) {
	if rule.IsSlice {
		if rule.RuleType == RuleTypeFile {
			fhs := getFiles(req, rule.Name)
			if len(fhs) < 1 {
				if rule.Optional {
					return nil, ErrorTypeNil
				}
				return nil, ErrorTypeMissRequired
			}

			if et := rule.sliceLenOk(len(fhs)); et != ErrorTypeNil {
				return nil, et
			}

			for _, fh := range fhs {
				if et := rule.fileOk(fh); et != ErrorTypeNil {
					return nil, et
				}
			}

			nfhs := make([]*multipart.FileHeader, len(fhs), len(fhs))
			copy(nfhs, fhs)
			return nfhs, ErrorTypeNil
		} else {
			_ = req.ParseForm()
			svs := req.Form[rule.Name]
			if len(svs) < 1 {
				if rule.Optional {
					return nil, ErrorTypeNil
				}
				return nil, ErrorTypeMissRequired
			}

			if et := rule.sliceLenOk(len(svs)); et != ErrorTypeNil {
				return nil, et
			}

			sliceVal := reflect.MakeSlice(reflect.SliceOf(rule.Gotype), 0, len(svs))

			for _, sv := range svs {
				ele, et := rule.one(sv)
				if et != ErrorTypeNil {
					return nil, et
				}
				sliceVal = reflect.Append(sliceVal, reflect.ValueOf(ele))
			}
			return sliceVal.Interface(), ErrorTypeNil
		}
	}

	if rule.RuleType == RuleTypeFile {
		fhs := getFiles(req, rule.Name)
		if len(fhs) < 1 {
			if rule.Optional {
				return nil, ErrorTypeNil
			}
			return nil, ErrorTypeMissRequired
		}

		fp := fhs[0]
		if et := rule.fileOk(fp); et != ErrorTypeNil {
			return nil, et
		}
		return fp, ErrorTypeNil
	}

	sv := req.FormValue(rule.Name)
	if len(sv) < 1 {
		if rule.Optional {
			return nil, ErrorTypeNil
		}
		return nil, ErrorTypeMissRequired
	}
	return rule.one(sv)
}

type Rules struct {
	Gotype reflect.Type
	Data   []*Rule
}

func fromReq(t reflect.Type, req *http.Request) (reflect.Value, *Error) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	ptr := reflect.New(t)
	err := ptr.Interface().(Binder).FromRequest(req)
	if err != nil {
		return reflect.Value{}, err
	}
	return ptr.Elem(), nil
}

func (rules *Rules) BindAndValidate(req *http.Request, dist reflect.Value) error {
	var err Error
	for _, rule := range rules.Data {
		if len(rule.Index) > 1 {
			continue
		}

		if rule.RuleType == RuleTypeBinder {
			vv, ep := fromReq(rule.Gotype, req)
			if ep != nil {
				return ep
			}
			fv := dist.FieldByIndex(rule.Index)
			for fv.Type() != vv.Type() {
				vv = vv.Addr()
			}
			fv.Set(vv)
			continue
		}

		v, et := rule.get(req)
		if et != ErrorTypeNil {
			err.Rule = rule
			err.Type = et
			return &err
		}
		if v == nil {
			continue
		}
		vv := reflect.ValueOf(v)
		if !vv.IsValid() {
			continue
		}
		dist.FieldByIndex(rule.Index).Set(vv)
	}
	return nil
}

func (rules *Rules) Validate(v any) error {
	vv := reflect.ValueOf(v)
	for vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}

	for _, rule := range rules.Data {
		fv, err := vv.FieldByIndexErr(rule.Index)
		if err != nil {
			if rule.Optional {
				continue
			}
			return &Error{Type: ErrorTypeMissRequired, Rule: rule}
		}
		fmt.Println(rule.Name, fv)
	}
	return nil
}
