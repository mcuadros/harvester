package format

import (
	"bytes"
	"strings"
)

type CSVConfig struct {
	Fields    string
	NotQuoted bool
	Quote     byte
	Trim      bool
	Separator byte
}

type CSV struct {
	fields    []string
	quoted    bool
	quote     byte
	trim      bool
	separator byte
}

func NewCSV(config *CSVConfig) *CSV {
	format := new(CSV)
	format.SetConfig(config)

	return format
}

func (self *CSV) SetConfig(config *CSVConfig) {
	for _, field := range strings.Split(config.Fields, ",") {
		self.fields = append(self.fields, field)
	}

	self.quoted = !config.NotQuoted
	self.trim = config.Trim

	self.quote = config.Quote
	if self.quote == 0 {
		self.quote = '"'
	}

	self.separator = config.Separator
	if self.separator == 0 {
		self.separator = ','
	}
}

func (self *CSV) Parse(line string) map[string]string {
	record := make(map[string]string)
	chars := []byte(line)

	max := len(self.fields)
	index := 0
	quoted := false
	field := make([]byte, 0)
	for _, char := range chars {
		if self.quoted && char == self.quote {
			if !quoted {
				quoted = true
			} else {
				quoted = false
			}
		} else if !quoted && char == self.separator {
			if self.trim {
				field = trim(field)
			}

			record[self.fields[index]] = string(field)
			field = make([]byte, 0)
			index++

			if index >= max {
				break
			}
		} else {
			field = append(field, char)
		}
	}

	if max > index {
		if self.trim {
			field = trim(field)
		}

		record[self.fields[index]] = string(field)
	}

	return record
}

// bytes.TrimSpace may return nil...
func trim(s []byte) []byte {
	t := bytes.TrimSpace(s)
	if t == nil {
		return s[0:0]
	}
	return t
}
