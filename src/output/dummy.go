package output

import (
	"fmt"
	"harvesterd/intf"
)

type DummyConfig struct {
	Print bool
}

type Dummy struct {
	printInfo bool
}

func NewDummy(config *DummyConfig) *Dummy {
	output := new(Dummy)
	output.SetConfig(config)
	return output
}

func (o *Dummy) SetConfig(config *DummyConfig) {
	o.printInfo = config.Print
}

func (o *Dummy) PutRecord(record intf.Record) bool {
	if o.printInfo {
		fmt.Println(record)
	}

	return true
}
