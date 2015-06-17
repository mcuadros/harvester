package harvesterd

import (
	"sync"

	"github.com/mcuadros/harvesterd/src/intf"
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
	inputs        []intf.Input
	hasProcessors bool
	processors    []intf.PostProcessor
}

func NewReader() *Reader {
	reader := new(Reader)

	return reader
}

func (r *Reader) SetInputs(inputs []intf.Input) {
	r.inputs = inputs
}

func (r *Reader) SetProcessors(processors []intf.PostProcessor) {
	if len(processors) > 0 {
		r.hasProcessors = true
	}

	r.processors = processors
}

func (r *Reader) SetChannels(recordsChan RecordsChan, closeChan CloseChan) {
	r.recordsChan = recordsChan
	r.closeChan = closeChan
}

func (r *Reader) GoRead() {
	r.setChannelToProcessors()
	go r.doReadIntoChannel()
}

func (r *Reader) doReadIntoChannel() {
	for _, input := range r.inputs {
		r.wait.Add(1)
		go r.readInputIntoChannel(input)
	}

	r.wait.Wait()

	r.teardownProcessors()
	r.closeChan <- true
}

func (r *Reader) readInputIntoChannel(input intf.Input) {
	for !input.IsEOF() {
		record := input.GetRecord()
		r.emitRecord(record)
	}

	r.wait.Done()
}

func (r *Reader) emitRecord(record intf.Record) {
	if len(record) > 0 {
		if r.applyProcessors(record) {
			r.recordsChan <- record
		}

		r.counter++
	}
}

func (r *Reader) setChannelToProcessors() {
	if r.hasProcessors {
		for _, proc := range r.processors {
			proc.SetChannel(r.recordsChan)
		}
	}
}

func (r *Reader) applyProcessors(record intf.Record) bool {
	if r.hasProcessors {
		for _, proc := range r.processors {
			if proc.Do(record) == false {
				return false
			}
		}
	}

	return true
}

func (r *Reader) teardownProcessors() {
	for _, proc := range r.processors {
		proc.Teardown()
	}
}

func (r *Reader) teardownInputs() {
	for _, input := range r.inputs {
		input.Teardown()
	}
}

func (r *Reader) Teardown() {
	r.teardownInputs()
}
