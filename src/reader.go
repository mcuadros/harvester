package harvesterd

import (
	. "harvesterd/intf"
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

func NewReader() *Reader {
	reader := new(Reader)

	return reader
}

func (self *Reader) SetInputs(inputs []Input) {
	self.inputs = inputs
}

func (self *Reader) GoReadIntoChannel(channel chan Record) {
	go self.doReadIntoChannel(channel)
}

func (self *Reader) doReadIntoChannel(channel chan Record) {
	defer close(channel)

	for _, input := range self.inputs {
		self.wait.Add(1)
		go self.readInputIntoChannel(input, channel)
	}

	self.wait.Wait()
}

func (self *Reader) readInputIntoChannel(input Input, channel chan Record) {
	for !input.IsEOF() {
		record := input.GetRecord()
		self.emitRecord(channel, record)
	}

	self.wait.Done()
}

func (self *Reader) emitRecord(channel chan Record, row Record) {
	if len(row) > 0 {
		channel <- row
		self.counter++
	}
}

func (self *Reader) Finish() {
	for _, input := range self.inputs {
		input.Finish()
	}
}
