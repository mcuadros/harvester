package format

import (
	"bytes"
	"strings"

	"github.com/mcuadros/harvester/src/intf"
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

func (f *CSV) SetConfig(config *CSVConfig) {
	f.format = NewFormatHelper(config.Format)
	f.parseFieldConfig(config.Fields)

	f.quoted = !config.NotQuoted
	f.trim = config.Trim

	f.quote = config.Quote
	if f.quote == 0 {
		f.quote = '"'
	}

	f.separator = config.Separator
	if f.separator == 0 {
		f.separator = ','
	}
}

func (f *CSV) parseFieldConfig(fields string) {
	for _, field := range strings.Split(fields, ",") {
		f.fields = append(f.fields, field)
	}
}

func (f *CSV) Parse(line string) intf.Record {
	record := make(intf.Record)
	chars := []byte(line)

	max := len(f.fields)
	index := 0
	quoted := false
	value := make([]byte, 0)
	for _, char := range chars {
		if f.quoted && char == f.quote {
			if !quoted {
				quoted = true
			} else {
				quoted = false
			}
		} else if !quoted && char == f.separator {
			if f.trim {
				value = trim(value)
			}

			field := f.fields[index]
			index++

			if field != "_" {
				record[field] = f.format.Format(field, string(value))
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
		if f.trim {
			value = trim(value)
		}

		field := f.fields[index]
		record[field] = f.format.Format(field, string(value))
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
