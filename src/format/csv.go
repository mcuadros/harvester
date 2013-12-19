package format

import (
	"bytes"
	"strings"
)

type CSVConfig struct {
	Fields    string
	Quoted    bool
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

	self.quoted = true
	self.separator = ','
	self.quote = '"'
	self.trim = false
}

func (self *CSV) Parse(line string) map[string]string {
	record := make(map[string]string)
	chars := []byte(line)

	max := len(self.fields)
	index := 0
	quoted := false
	field := make([]byte, max)
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
			field = make([]byte, max)
			index++

			if index >= max {
				break
			}
		} else {
			field = append(field, char)
		}
	}

	if max > index {
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
