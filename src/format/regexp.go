package format

import (
	"harvesterd/intf"
	"regexp"
)

type RegExpConfig struct {
	Pattern string
	Format  string
}

type RegExp struct {
	fields []string
	regexp *regexp.Regexp
	format *FormatHelper
}

func NewRegExp(config *RegExpConfig) *RegExp {
	format := new(RegExp)
	format.SetConfig(config)

	return format
}

func (f *RegExp) SetConfig(config *RegExpConfig) {
	f.regexp = regexp.MustCompile(config.Pattern)
	f.format = NewFormatHelper(config.Format)
}

func (f *RegExp) Parse(line string) intf.Record {
	names := f.regexp.SubexpNames()
	values := f.regexp.FindStringSubmatch(line)

	record := make(intf.Record)
	for index, value := range values {
		if names[index] != "" {
			field := names[index]
			record[field] = f.format.Format(field, value)
		}
	}

	return record
}
