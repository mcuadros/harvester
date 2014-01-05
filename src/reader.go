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
	wait          sync.WaitGroup
	recordsChan   RecordsChan
	closeChan     CloseChan
	counter       int32
	inputs        []Input
	hasProcessors bool
	processors    []PostProcessor
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

func (self *Reader) SetChannels(recordsChan RecordsChan, closeChan CloseChan) {
	self.recordsChan = recordsChan
	self.closeChan = closeChan
}

func (self *Reader) GoRead() {
	self.setChannelToProcessors()
	go self.doReadIntoChannel()
}

func (self *Reader) doReadIntoChannel() {
	for _, input := range self.inputs {
		self.wait.Add(1)
		go self.readInputIntoChannel(input)
	}

	self.wait.Wait()

	self.teardownProcessors()
	self.closeChan <- true
}

func (self *Reader) readInputIntoChannel(input Input) {
	for !input.IsEOF() {
		record := input.GetRecord()
		self.emitRecord(record)
	}

	self.wait.Done()
}

func (self *Reader) emitRecord(record Record) {
	if len(record) > 0 {
		if self.applyProcessors(record) {
			self.recordsChan <- record
		}

		self.counter++
	}
}

func (self *Reader) setChannelToProcessors() {
	if self.hasProcessors {
		for _, proc := range self.processors {
			proc.SetChannel(self.recordsChan)
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

func (self *Reader) teardownProcessors() {
	for _, proc := range self.processors {
		proc.Teardown()
	}
}

func (self *Reader) teardownInputs() {
	for _, input := range self.inputs {
		input.Teardown()
	}
}

func (self *Reader) Teardown() {
	self.teardownInputs()
}
