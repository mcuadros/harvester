package format

import (
	. "harvesterd/intf"
	"regexp"
)

type RegExpConfig struct {
	Pattern string
}

type RegExp struct {
	fields []string
	regexp *regexp.Regexp
}

func NewRegExp(config *RegExpConfig) *RegExp {
	format := new(RegExp)
	format.SetConfig(config)

	return format
}

func (self *RegExp) SetConfig(config *RegExpConfig) {
	self.regexp = regexp.MustCompile(config.Pattern)
}

func (self *RegExp) Parse(line string) Record {
	names := self.regexp.SubexpNames()
	values := self.regexp.FindStringSubmatch(line)

	record := make(Record)
	for index, value := range values {
		if names[index] != "" {
			record[names[index]] = value
		}
	}

	return record
}
