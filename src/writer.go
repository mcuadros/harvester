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

func (w *Writer) SetReaders(readers []*Reader) {
	w.readers = readers
}

func (w *Writer) SetOutputsFactory(factory OutputsFactory) {
	w.outputs = factory
}

func (w *Writer) SetThreads(threads int) {
	w.maxThreads = int32(threads)
}

func (w *Writer) GetChannels() (RecordsChan, CloseChan) {
	return w.recordsChan, w.closeChan
}

func (w *Writer) IsAlive() bool {
	return atomic.LoadInt32(&w.threads) != 0
}

func (w *Writer) Setup() {
	w.createChannels()
	w.setupReaders()
}

func (w *Writer) createChannels() {
	w.recordsChan = make(RecordsChan, w.maxThreads)
	w.closeChan = make(CloseChan, 1)
}

func (w *Writer) setupReaders() {
	for _, reader := range w.readers {
		reader.SetChannels(w.recordsChan, w.closeChan)
		reader.GoRead()
	}
}

func (w *Writer) Boot() {
	w.goWaitForReadersClose()
	w.goWriteFromChannel()
}

func (w *Writer) goWriteFromChannel() {
	for i := int32(0); i < w.maxThreads; i++ {
		atomic.AddInt32(&w.threads, 1)
		go w.doWriteFromChannel()
	}
}

func (w *Writer) goWaitForReadersClose() {
	go func() {
		readersClosed := 0
		readersCount := len(w.readers)
		for _ = range w.closeChan {
			readersClosed++

			if readersClosed >= readersCount {
				close(w.recordsChan)
				break
			}
		}
	}()
}

func (w *Writer) doWriteFromChannel() {
	outputs := w.outputs()
	for record := range w.recordsChan {
		w.writeRecordFromChannel(outputs, record)
	}

	atomic.AddInt32(&w.threads, -1)
}

func (w *Writer) writeRecordFromChannel(outputs []intf.Output, record intf.Record) {
	for _, output := range outputs {
		w.writeRecordIntoOutput(output, record)
	}
}

func (w *Writer) writeRecordIntoOutput(output intf.Output, record intf.Record) {
	if output.PutRecord(record) {
		w.created++
	} else {
		w.failed++
	}
}

func (w *Writer) GetCounters() (int32, int32, int32, int32) {
	return w.created, w.failed, w.transferred, w.threads
}

func (w *Writer) ResetCounters() {
	w.created = 0
	w.failed = 0
	w.transferred = 0
}

func (w *Writer) Teardown() {
	w.teardownReaders()
}

func (w *Writer) teardownReaders() {
	for _, reader := range w.readers {
		reader.Teardown()
	}
}
