package collector

import (
	. "collector/intf"
	"sync"
	"sync/atomic"
)

type WriterConfig struct {
	Output  []string
	Threads int
}

type Writer struct {
	failed      int32
	created     int32
	transferred int32
	outputs     []Output
	maxThreads  int32
	threads     int32
	isAlive     bool
}

func NewWriter(output []Output, threads int) *Writer {
	writer := new(Writer)
	writer.SetOutputs(output)
	writer.SetThreads(int32(threads))

	return writer
}

func (self *Writer) SetOutputs(outputs []Output) {
	self.outputs = outputs
}

func (self *Writer) SetThreads(threads int32) {
	self.maxThreads = threads
}

func (self *Writer) GetCounters() (int32, int32, int32) {
	return self.created, self.failed, self.transferred
}

func (self *Writer) IsAlive() bool {
	return atomic.LoadInt32(&self.threads) != 0
}

func (self *Writer) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *Writer) GoWriteFromChannel() chan map[string]string {
	channel := make(chan map[string]string, self.maxThreads)
	for i := int32(0); i < self.maxThreads; i++ {
		atomic.AddInt32(&self.threads, 1)
		go self.doWriteFromChannel(channel)
	}

	return channel
}

func (self *Writer) doWriteFromChannel(channel chan map[string]string) {
	for record := range channel {
		self.writeRecordFromChannel(record)
	}

	atomic.AddInt32(&self.threads, -1)
}

func (self *Writer) writeRecordFromChannel(record map[string]string) {
	var wait sync.WaitGroup

	for _, output := range self.outputs {
		wait.Add(1)
		go self.writeRecordIntoOutput(output, record, &wait)
	}

	wait.Wait()
}

func (self *Writer) writeRecordIntoOutput(output Output, record map[string]string, wait *sync.WaitGroup) {
	if output.PutRecord(record) {
		self.created++
	} else {
		self.failed++
	}

	wait.Done()
}
