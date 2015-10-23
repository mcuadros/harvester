package format

import (
	"encoding/json"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
)

type JSONConfig struct {
	Empty string
}

type JSON struct {
}

func NewJSON(config *JSONConfig) *JSON {
	format := new(JSON)
	format.SetConfig(config)

	return format
}

func (f *JSON) SetConfig(config *JSONConfig) {

}

func (f *JSON) Parse(line string) intf.Record {
	var record intf.Record

	err := json.Unmarshal([]byte(line), &record)
	if err != nil {
		Warning("error: %s", err)
	}

	return record
}
