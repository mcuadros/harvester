package output

import (
	"fmt"
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

func (self *Dummy) SetConfig(config *DummyConfig) {
	self.printInfo = config.Print
}

func (self *Dummy) PutRecord(record map[string]string) bool {
	if self.printInfo {
		fmt.Println(record)
	}

	return true
}
