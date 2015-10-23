package format

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	. "github.com/mcuadros/harvester/src/logger"
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

func (h *FormatHelper) SetConfig(config string) {
	h.parseFieldsConfig(config)
}

func (h *FormatHelper) parseFieldsConfig(config string) {
	for _, fieldConfig := range helperConfigGroupRegExp.FindAllStringSubmatch(config, -1) {
		field, format, layout := h.parseField(fieldConfig[0])
		h.fields[field] = &fieldFormat{format: format, layout: layout}
	}
}

func (h *FormatHelper) parseField(fieldConfig string) (field string, format string, layout string) {
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

func (h *FormatHelper) GetFields() map[string]*fieldFormat {
	return h.fields
}

func (h *FormatHelper) Format(field, value string) interface{} {
	if _, ok := h.fields[field]; !ok {
		return value
	}

	switch h.fields[field].format {
	case "int":
		return h.toInt(value)
	case "float":
		return h.toFloat(value)
	case "bool":
		return h.toBool(value)
	case "string":
		return h.toString(value)
	case "time":
		return h.toTime(value, h.fields[field].layout)
	}

	return value
}

func (h *FormatHelper) toInt(original string) int {
	value, err := strconv.Atoi(strings.Trim(original, " \"'"))
	if err != nil {
		return 0
	}

	return value
}

func (h *FormatHelper) toFloat(original string) float64 {
	original = strings.Trim(original, " \"'")
	original = strings.Replace(original, ",", ".", -1)

	value, err := strconv.ParseFloat(original, 64)
	if err != nil {
		return 0
	}

	return value
}

func (h *FormatHelper) toBool(original string) interface{} {
	value, err := strconv.ParseBool(strings.Trim(original, " \"'"))
	if err != nil {
		return nil
	}

	return value
}

func (h *FormatHelper) toString(original string) string {
	return strings.Trim(original, " ")
}

func (h *FormatHelper) toTime(original string, layout string) interface{} {
	value, err := time.Parse(layout, original)
	if err != nil {
		return nil
	}

	return value
}
