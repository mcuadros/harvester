package format

import (
	"encoding/json"
	. "harvesterd/intf"
	. "harvesterd/logger"
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

func (self *JSON) SetConfig(config *JSONConfig) {

}

func (self *JSON) Parse(line string) Record {
	var record Record

	err := json.Unmarshal([]byte(line), &record)
	if err != nil {
		Warning("error: %s", err)
	}

	return record
}
