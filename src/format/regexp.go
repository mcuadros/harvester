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

func (self *RegExp) SetConfig(config *RegExpConfig) {
	self.regexp = regexp.MustCompile(config.Pattern)
	self.format = NewFormatHelper(config.Format)
}

func (self *RegExp) Parse(line string) intf.Record {
	names := self.regexp.SubexpNames()
	values := self.regexp.FindStringSubmatch(line)

	record := make(intf.Record)
	for index, value := range values {
		if names[index] != "" {
			field := names[index]
			record[field] = self.format.Format(field, value)
		}
	}

	return record
}
