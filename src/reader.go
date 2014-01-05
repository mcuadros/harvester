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
	channel       chan Record
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

func (self *Reader) SetChannel(channel chan Record) {
	self.channel = channel
}

func (self *Reader) GoRead() {
	self.setChannelToProcessors(self.channel)
	go self.doReadIntoChannel(self.channel)
}

func (self *Reader) doReadIntoChannel(channel chan Record) {
	for _, input := range self.inputs {
		self.wait.Add(1)
		go self.readInputIntoChannel(input, channel)
	}

	self.wait.Wait()
	self.finishProcessors()
	close(channel)
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
		if self.applyProcessors(record) {
			channel <- record
		}

		self.counter++
	}
}

func (self *Reader) setChannelToProcessors(channel chan Record) {
	if self.hasProcessors {
		for _, proc := range self.processors {
			proc.SetChannel(channel)
		}
	}
}

func (self *Reader) applyProcessors(record Record) bool {
	if self.hasProcessors {
		for _, proc := range self.processors {
			if proc.Do(record) == false {
				return false
			}
		}
	}

	return true
}

func (self *Reader) finishProcessors() {
	for _, proc := range self.processors {
		proc.Finish()
	}
}

func (self *Reader) Finish() {
	for _, input := range self.inputs {
		input.Finish()
	}
}
