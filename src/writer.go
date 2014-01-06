package harvesterd

import (
	. "harvesterd/intf"

	"sync"
	"sync/atomic"
)

type RecordsChan chan Record
type CloseChan chan bool

type WriterConfig struct {
	Output  []string
	Reader  []string
	Threads int
}

type BasicWriter struct {
	outputs     []Output
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

func NewWriter() *BasicWriter {
	writer := new(BasicWriter)

	return writer
}

func (self *BasicWriter) SetReaders(readers []*Reader) {
	self.readers = readers
}

func (self *BasicWriter) SetOutputs(outputs []Output) {
	self.outputs = outputs
}

func (self *BasicWriter) SetThreads(threads int) {
	self.maxThreads = int32(threads)
}

func (self *BasicWriter) GetChannels() (RecordsChan, CloseChan) {
	return self.recordsChan, self.closeChan
}

func (self *BasicWriter) IsAlive() bool {
	return atomic.LoadInt32(&self.threads) != 0
}

func (self *BasicWriter) Setup() {
	self.createChannels()
	self.setupReaders()
}

func (self *BasicWriter) createChannels() {
	self.recordsChan = make(RecordsChan, self.maxThreads)
	self.closeChan = make(CloseChan, 1)
}

func (self *BasicWriter) setupReaders() {
	for _, reader := range self.readers {
		reader.SetChannels(self.recordsChan, self.closeChan)
		reader.GoRead()
	}
}

func (self *BasicWriter) Boot() {
	self.goWaitForReadersClose()
	self.goWriteFromChannel()
}

func (self *BasicWriter) goWriteFromChannel() {
	for i := int32(0); i < self.maxThreads; i++ {
		atomic.AddInt32(&self.threads, 1)
		go self.doWriteFromChannel()
	}
}

func (self *BasicWriter) goWaitForReadersClose() {
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

func (self *BasicWriter) doWriteFromChannel() {
	for record := range self.recordsChan {
		self.writeRecordFromChannel(record)
	}

	atomic.AddInt32(&self.threads, -1)
}

func (self *BasicWriter) writeRecordFromChannel(record Record) {
	var wait sync.WaitGroup

	for _, output := range self.outputs {
		wait.Add(1)
		go self.writeRecordIntoOutput(output, record, &wait)
	}

	wait.Wait()
}

func (self *BasicWriter) writeRecordIntoOutput(output Output, record Record, wait *sync.WaitGroup) {
	if output.PutRecord(record) {
		self.created++
	} else {
		self.failed++
	}

	wait.Done()
}

func (self *BasicWriter) GetCounters() (int32, int32, int32, int32) {
	return self.created, self.failed, self.transferred, self.threads
}

func (self *BasicWriter) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *BasicWriter) Teardown() {
	self.teardownReaders()
}

func (self *BasicWriter) teardownReaders() {
	for _, reader := range self.readers {
		reader.Teardown()
	}
}
