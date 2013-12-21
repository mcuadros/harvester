package collector

import (
	. "collector/intf"
	"sync"
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
	threads     int
	isAlive     bool
}

func NewWriter(output []Output, threads int) *Writer {
	writer := new(Writer)
	writer.SetOutputs(output)
	writer.SetThreads(threads)

	return writer
}

func (self *Writer) SetOutputs(outputs []Output) {
	self.outputs = outputs
}

func (self *Writer) SetThreads(threads int) {
	self.threads = threads
}

func (self *Writer) GetCounters() (int32, int32, int32) {
	return self.created, self.failed, self.transferred
}

func (self *Writer) IsAlive() bool {
	return self.isAlive
}

func (self *Writer) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *Writer) GoWriteFromChannel() chan map[string]string {
	channel := make(chan map[string]string, self.threads)

	for i := 0; i < self.threads; i++ {
		go self.doWriteFromChannel(channel)
	}

	return channel
}

func (self *Writer) doWriteFromChannel(channel chan map[string]string) {
	self.isAlive = true
	for record := range channel {
		self.writeRecordFromChannel(record)
	}

	self.isAlive = false
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
