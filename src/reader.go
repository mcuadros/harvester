package collector

import (
	. "collector/intf"
	"fmt"
	"sync"
)

type ReaderConfig struct {
	Input []string
}

type Reader struct {
	counter int32
	inputs  []Input
	wait    sync.WaitGroup
}

func NewReader(inputs []Input) *Reader {
	reader := new(Reader)
	reader.SetInputs(inputs)

	return reader
}

func (self *Reader) SetInputs(inputs []Input) {
	self.inputs = inputs
}

func (self *Reader) GoReadIntoChannel(channel chan map[string]string) {
	go self.doReadIntoChannel(channel)
}

func (self *Reader) doReadIntoChannel(channel chan map[string]string) {
	defer close(channel)

	for _, input := range self.inputs {
		self.wait.Add(1)
		go self.readInputIntoChannel(input, channel)
	}

	self.wait.Wait()
}

func (self *Reader) readInputIntoChannel(input Input, channel chan map[string]string) {
	for !input.IsEOF() {
		record := input.GetRecord()
		self.emitRecord(channel, record)
	}

	self.wait.Done()
}

func (self *Reader) emitRecord(channel chan map[string]string, row map[string]string) {
	if len(row) > 0 {
		channel <- row
		self.counter++
		if self.counter%1000 == 0 {
			fmt.Println(fmt.Sprintf("%d", self.counter))
		}
	}
}
