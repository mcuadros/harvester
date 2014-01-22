package util

import (
	"harvesterd/intf"
	"regexp"
	"strconv"
	"strings"
)

import "github.com/stretchr/objx"

type Template struct {
	template string
	tokens   []string
	isEmpty  bool
}

var tokenRegexp = regexp.MustCompile(`%{([\w.]+)}`)

func NewTemplate(template string) *Template {
	tmpl := &Template{template: template}
	tmpl.ParseTokens()

	return tmpl
}

func (self *Template) ParseTokens() {
	tokens := tokenRegexp.FindAll([]byte(self.template), -1)
	self.tokens = make([]string, len(tokens))

	if len(tokens) == 0 {
		self.isEmpty = true
	}

	for index, token := range tokens {
		self.tokens[index] = string(token[2 : len(token)-1])
	}
}

func (self *Template) Apply(record intf.Record) string {
	if self.isEmpty {
		return self.template
	}

	return self.replaceTokens(record)
}

func (self *Template) replaceTokens(record intf.Record) string {
	output := self.template
	mapper := objx.Map(record)

	for _, token := range self.tokens {
		value := self.castValueToString(mapper.Get(token))
		output = strings.Replace(output, `%{`+token+`}`, value, -1)
	}

	return output
}

func (self *Template) castValueToString(value *objx.Value) string {
	switch {
	case value.IsStr():
		return value.Str()
	case value.IsBool():
		return strconv.FormatBool(value.Bool())
	case value.IsFloat32():
		return strconv.FormatFloat(float64(value.Float32()), 'f', -1, 32)
	case value.IsFloat64():
		return strconv.FormatFloat(value.Float64(), 'f', -1, 64)
	case value.IsInt():
		return strconv.FormatInt(int64(value.Int()), 10)
	case value.IsInt():
		return strconv.FormatInt(int64(value.Int()), 10)
	case value.IsInt8():
		return strconv.FormatInt(int64(value.Int8()), 10)
	case value.IsInt16():
		return strconv.FormatInt(int64(value.Int16()), 10)
	case value.IsInt32():
		return strconv.FormatInt(int64(value.Int32()), 10)
	case value.IsInt64():
		return strconv.FormatInt(value.Int64(), 10)
	case value.IsUint():
		return strconv.FormatUint(uint64(value.Uint()), 10)
	case value.IsUint8():
		return strconv.FormatUint(uint64(value.Uint8()), 10)
	case value.IsUint16():
		return strconv.FormatUint(uint64(value.Uint16()), 10)
	case value.IsUint32():
		return strconv.FormatUint(uint64(value.Uint32()), 10)
	case value.IsUint64():
		return strconv.FormatUint(value.Uint64(), 10)
	}

	return ""
}
