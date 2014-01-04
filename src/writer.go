package harvesterd

import (
	. "harvesterd/intf"
	. "harvesterd/logger"
	"sync"
	"sync/atomic"
)

type WriterConfig struct {
	Output  []string
	Threads int
}

type Writer struct {
	outputs     []Output
	failed      int32
	created     int32
	transferred int32
	maxThreads  uint32
	threads     int32
	isAlive     bool
	mutex       sync.Mutex
}

func NewWriter() *Writer {
	writer := new(Writer)

	return writer
}

func (self *Writer) SetOutputs(outputs []Output) {
	self.outputs = outputs
}

func (self *Writer) SetThreads(threads int) {
	self.maxThreads = uint32(threads)
}

func (self *Writer) IsAlive() bool {
	return atomic.LoadInt32(&self.threads) != 0
}

func (self *Writer) GoWriteFromChannel() chan Record {
	channel := make(chan Record, self.maxThreads)
	for i := uint32(0); i < self.maxThreads; i++ {
		atomic.AddInt32(&self.threads, 1)
		go self.doWriteFromChannel(channel)
	}

	return channel
}

func (self *Writer) doWriteFromChannel(channel chan Record) {
	for record := range channel {
		self.writeRecordFromChannel(record)
	}

	atomic.AddInt32(&self.threads, -1)
}

func (self *Writer) writeRecordFromChannel(record Record) {
	var wait sync.WaitGroup

	for _, output := range self.outputs {
		wait.Add(1)
		go self.writeRecordIntoOutput(output, record, &wait)
	}

	wait.Wait()
}

func (self *Writer) writeRecordIntoOutput(output Output, record Record, wait *sync.WaitGroup) {
	if output.PutRecord(record) {
		self.created++
	} else {
		self.failed++
	}

	wait.Done()
}

func (self *Writer) GetCounters() (int32, int32, int32) {
	return self.created, self.failed, self.transferred
}

func (self *Writer) ResetCounters() {
	self.created = 0
	self.failed = 0
	self.transferred = 0
}

func (self *Writer) PrintCounters(elapsedSeconds int) {
	created, failed, _ := self.GetCounters()
	self.ResetCounters()

	logFormat := "Created %d document(s), failed %d times(s), %g doc/sec in %d thread(s)"

	rate := float64(created+failed) / float64(elapsedSeconds)
	Info(logFormat, created, failed, rate, self.threads)
}
