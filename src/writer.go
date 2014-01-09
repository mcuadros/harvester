package harvesterd

import (
	"harvesterd/intf"
	"sync"
	"sync/atomic"
)

type RecordsChan chan intf.Record
type CloseChan chan bool

type WriterConfig struct {
	Output  []string
	Reader  []string
	Threads int
}

type Writer struct {
	outputs     OutputsFactory
	readers     []*Reader
	failed      int32
	created     int32
	transferred int32
	maxThreads  int32
	threads     int32
	mutex       sync.Mutex
	recordsChan RecordsChan
	closeChan   CloseChan
}

func NewWriter() *Writer {
	writer := new(Writer)

	return writer
}

func (self *Writer) SetReaders(readers []*Reader) {
	self.readers = readers
}

func (self *Writer) SetOutputsFactory(factory OutputsFactory) {
	self.outputs = factory
}

func (self *Writer) SetThreads(threads int) {
	self.maxThreads = int32(threads)
}

func (self *Writer) GetChannels() (RecordsChan, CloseChan) {
	return self.recordsChan, self.closeChan
}

func (self *Writer) IsAlive() bool {
	return atomic.LoadInt32(&self.threads) != 0
}

func (self *Writer) Setup() {
	self.createChannels()
	self.setupReaders()
}

func (self *Writer) createChannels() {
	self.recordsChan = make(RecordsChan, self.maxThreads)
	self.closeChan = make(CloseChan, 1)
}

func (self *Writer) setupReaders() {
	for _, reader := range self.readers {
		reader.SetChannels(self.recordsChan, self.closeChan)
		reader.GoRead()
	}
}

func (self *Writer) Boot() {
	self.goWaitForReadersClose()
	self.goWriteFromChannel()
}

func (self *Writer) goWriteFromChannel() {
	for i := int32(0); i < self.maxThreads; i++ {
		atomic.AddInt32(&self.threads, 1)
		go self.doWriteFromChannel()
	}
}

func (self *Writer) goWaitForReadersClose() {
	go func() {
		readersClosed := 0
		readersCount := len(self.readers)
		for _ = range self.closeChan {
			readersClosed++

			if readersClosed >= readersCount {
				close(self.recordsChan)
				break
			}
		}
	}()
}

func (self *Writer) doWriteFromChannel() {
	outputs := self.outputs()
	for record := range self.recordsChan {
		self.writeRecordFromChannel(outputs, record)
	}

	atomic.AddInt32(&self.threads, -1)
}

func (self *Writer) writeRecordFromChannel(outputs []intf.Output, record intf.Record) {
	for _, output := range outputs {
		self.writeRecordIntoOutput(output, record)
	}
}

func (self *Writer) writeRecordIntoOutput(output intf.Output, record intf.Record) {
	if output.PutRecord(record) {
		self.created++
	} else {
		self.failed++
	}
}

func (self *Writer) GetCounters() (int32, int32, int32, int32) {
	return self.created, self.failed, self.transferred, self.threads
}

func (self *Writer) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *Writer) Teardown() {
	self.teardownReaders()
}

func (self *Writer) teardownReaders() {
	for _, reader := range self.readers {
		reader.Teardown()
	}
}
