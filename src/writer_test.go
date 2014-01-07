package harvesterd

import (
	"fmt"
	"harvesterd/intf"
	"strconv"
	"time"
)

import . "launchpad.net/gocheck"

type WriterSuite struct{}

var _ = Suite(&WriterSuite{})

func (s *WriterSuite) TestWriteFromChannelSingleOutput(c *C) {
	output := new(MockOutput)
	output.Return = true

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{output})
	writer.SetThreads(1)
	writer.Setup()
	writer.Boot()

	recordChan, closeChan := writer.GetChannels()
	go func() {
		for i := 0; i < 10; i++ {
			recordChan <- intf.Record{"foo": fmt.Sprintf("%d", i)}
		}

		close(recordChan)
		close(closeChan)
	}()

	for {
		time.Sleep(100 * time.Microsecond)
		if !writer.IsAlive() {
			break
		}
	}

	created, failed, _, _ := writer.GetCounters()
	c.Assert(created, Equals, int32(10))
	c.Assert(failed, Equals, int32(0))
	c.Assert(output.Count, Equals, 45)

}

func (s *WriterSuite) TestWriteFromChannelMultipleOutput(c *C) {
	outputw := new(MockOutput)
	outputw.Return = true
	outputf := new(MockOutput)
	outputf.Return = false

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{outputw, outputf})
	writer.SetThreads(1)
	writer.Setup()
	writer.Boot()

	recordChan, closeChan := writer.GetChannels()
	go func() {

		for i := 0; i < 10; i++ {
			recordChan <- intf.Record{"foo": fmt.Sprintf("%d", i)}
		}

		close(recordChan)
		close(closeChan)
	}()

	for {
		time.Sleep(100 * time.Microsecond)
		if !writer.IsAlive() {
			break
		}
	}

	created, failed, _, _ := writer.GetCounters()
	c.Assert(created, Equals, int32(10))
	c.Assert(failed, Equals, int32(10))
	c.Assert(outputw.Count, Equals, 45)
	c.Assert(outputf.Count, Equals, 45)

}

func (s *WriterSuite) TestWriteIsAlive(c *C) {
	outputw := new(MockOutput)
	outputw.Return = true
	outputf := new(MockOutput)
	outputf.Return = false

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{outputw, outputf})
	writer.SetThreads(1)
	writer.Setup()
	writer.Boot()

	channel, _ := writer.GetChannels()

	c.Assert(writer.IsAlive(), Equals, true)

	time.Sleep(100 * time.Microsecond)
	c.Assert(writer.IsAlive(), Equals, true)
	close(channel)

	time.Sleep(100 * time.Microsecond)
	c.Assert(writer.IsAlive(), Equals, false)
}

type MockOutput struct {
	Count  int
	Return bool
}

func (self *MockOutput) PutRecord(record intf.Record) bool {
	number, _ := strconv.Atoi(record["foo"].(string))
	self.Count += number
	return self.Return
}
