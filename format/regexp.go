package format

import (
	"regexp"
)

type RegExpConfig struct {
	Fields  string
	Pattern string
}

type RegExp struct {
	fields []string
	regexp *regexp.Regexp
}

func NewRegExp(config RegExpConfig) *RegExp {
	format := new(RegExp)
	format.SetConfig(config)

	return format
}

func (self *RegExp) SetConfig(config RegExpConfig) {
	self.regexp = regexp.MustCompile(config.Pattern)
}

func (self *RegExp) Parse(line string) map[string]string {
	names := self.regexp.SubexpNames()
	values := self.regexp.FindStringSubmatch(line)

	record := make(map[string]string)
	for index, value := range values {
		record[names[index]] = value
	}

	return record
}
