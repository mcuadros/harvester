package harvesterd

import (
	. "harvesterd/intf"
	"sync"
)

type ReaderConfig struct {
	Input     []string
	Processor []string
}

type Reader struct {
	counter       int32
	inputs        []Input
	processors    []PostProcessor
	hasProcessors bool
	wait          sync.WaitGroup
}

func NewReader() *Reader {
	reader := new(Reader)

	return reader
}

func (self *Reader) SetInputs(inputs []Input) {
	self.inputs = inputs
}

func (self *Reader) SetProcessors(processors []PostProcessor) {
	if len(processors) > 0 {
		self.hasProcessors = true
	}

	self.processors = processors
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

func (self *Reader) emitRecord(channel chan Record, record Record) {
	if len(record) > 0 {
		self.applyProcessors(record)

		channel <- record
		self.counter++
	}
}

func (self *Reader) applyProcessors(record Record) {
	if self.hasProcessors {
		for _, proc := range self.processors {
			proc.Do(record)
		}
	}
}

func (self *Reader) Finish() {
	for _, input := range self.inputs {
		input.Finish()
	}
}
