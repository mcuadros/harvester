package format

import (
	. "harvesterd/logger"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validFormats = []string{"int", "float", "bool", "string", "time"}
var helperConfigRegExp = regexp.MustCompile("^\\((\\w+):?\"?([^\"]+)?\"?\\)(\\w+),?$")
var helperConfigGroupRegExp = regexp.MustCompile("(\\(\\w+(?::\"[^\"]+\")?\\)\\w+),?")

type FormatHelper struct {
	fields map[string]*fieldFormat
}

type fieldFormat struct {
	format string
	layout string
}

func NewFormatHelper(config string) *FormatHelper {
	helper := &FormatHelper{fields: make(map[string]*fieldFormat)}
	helper.SetConfig(config)

	return helper
}

func (self *FormatHelper) SetConfig(config string) {
	self.parseFieldsConfig(config)
}

func (self *FormatHelper) parseFieldsConfig(config string) {
	for _, fieldConfig := range helperConfigGroupRegExp.FindAllStringSubmatch(config, -1) {
		field, format, layout := self.parseField(fieldConfig[0])
		self.fields[field] = &fieldFormat{format: format, layout: layout}
	}
}

func (self *FormatHelper) parseField(fieldConfig string) (field string, format string, layout string) {
	config := helperConfigRegExp.FindStringSubmatch(fieldConfig)
	if len(config) != 4 {
		Critical("Malformed format config \"%s\"", fieldConfig)
	}

	isValid := false
	for _, valid := range validFormats {
		if valid == config[1] {
			isValid = true
		}
	}

	if !isValid {
		Critical("Unknown format config \"%s\", valid: %s", fieldConfig, validFormats)
	}

	return config[3], config[1], config[2]
}

func (self *FormatHelper) GetFields() map[string]*fieldFormat {
	return self.fields
}

func (self *FormatHelper) Format(field, value string) interface{} {
	if _, ok := self.fields[field]; !ok {
		return value
	}

	switch self.fields[field].format {
	case "int":
		return self.toInt(value)
	case "float":
		return self.toFloat(value)
	case "bool":
		return self.toBool(value)
	case "string":
		return self.toString(value)
	case "time":
		return self.toTime(value, self.fields[field].layout)
	}

	return value
}

func (self *FormatHelper) toInt(original string) int {
	value, err := strconv.Atoi(strings.Trim(original, " \"'"))
	if err != nil {
		return 0
	}

	return value
}

func (self *FormatHelper) toFloat(original string) float64 {
	original = strings.Trim(original, " \"'")
	original = strings.Replace(original, ",", ".", -1)

	value, err := strconv.ParseFloat(original, 64)
	if err != nil {
		return 0
	}

	return value
}

func (self *FormatHelper) toBool(original string) interface{} {
	value, err := strconv.ParseBool(strings.Trim(original, " \"'"))
	if err != nil {
		return nil
	}

	return value
}

func (self *FormatHelper) toString(original string) string {
	return strings.Trim(original, " ")
}

func (self *FormatHelper) toTime(original string, layout string) interface{} {
	value, err := time.Parse(layout, original)
	if err != nil {
		return nil
	}

	return value
}
