package format

import (
	"bytes"
	"harvesterd/intf"
	"strings"
)

type CSVConfig struct {
	Fields    string
	Format    string
	NotQuoted bool
	Quote     byte
	Trim      bool
	Separator byte
}

type CSV struct {
	format    *FormatHelper
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
	self.format = NewFormatHelper(config.Format)
	self.parseFieldConfig(config.Fields)

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

func (self *CSV) parseFieldConfig(fields string) {
	for _, field := range strings.Split(fields, ",") {
		self.fields = append(self.fields, field)
	}
}

func (self *CSV) Parse(line string) intf.Record {
	record := make(intf.Record)
	chars := []byte(line)

	max := len(self.fields)
	index := 0
	quoted := false
	value := make([]byte, 0)
	for _, char := range chars {
		if self.quoted && char == self.quote {
			if !quoted {
				quoted = true
			} else {
				quoted = false
			}
		} else if !quoted && char == self.separator {
			if self.trim {
				value = trim(value)
			}

			field := self.fields[index]
			index++

			if field != "_" {
				record[field] = self.format.Format(field, string(value))
			}

			value = make([]byte, 0)

			if index >= max {
				break
			}
		} else {
			value = append(value, char)
		}
	}

	if max > index {
		if self.trim {
			value = trim(value)
		}

		field := self.fields[index]
		record[field] = self.format.Format(field, string(value))
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
