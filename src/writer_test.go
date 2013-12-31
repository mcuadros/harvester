package harvesterd

import (
	"harvesterd/intf"
	"fmt"
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

	channel := writer.GoWriteFromChannel()
	go func(channel chan map[string]string) {
		for i := 0; i < 10; i++ {
			channel <- map[string]string{"foo": fmt.Sprintf("%d", i)}
		}

		close(channel)
	}(channel)

	for {
		time.Sleep(100 * time.Microsecond)
		if !writer.IsAlive() {
			break
		}
	}

	created, failed, _ := writer.GetCounters()
	c.Check(created, Equals, int32(10))
	c.Check(failed, Equals, int32(0))
	c.Check(output.Count, Equals, 45)

}

func (s *WriterSuite) TestWriteFromChannelMultipleOutput(c *C) {
	outputw := new(MockOutput)
	outputw.Return = true
	outputf := new(MockOutput)
	outputf.Return = false

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{outputw, outputf})
	writer.SetThreads(1)

	channel := writer.GoWriteFromChannel()
	go func(channel chan map[string]string) {
		for i := 0; i < 10; i++ {
			channel <- map[string]string{"foo": fmt.Sprintf("%d", i)}
		}

		close(channel)
	}(channel)

	for {
		time.Sleep(100 * time.Microsecond)
		if !writer.IsAlive() {
			break
		}
	}

	created, failed, _ := writer.GetCounters()
	c.Check(created, Equals, int32(10))
	c.Check(failed, Equals, int32(10))
	c.Check(outputw.Count, Equals, 45)
	c.Check(outputf.Count, Equals, 45)

}

func (s *WriterSuite) TestWriteIsAlive(c *C) {
	outputw := new(MockOutput)
	outputw.Return = true
	outputf := new(MockOutput)
	outputf.Return = false

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{outputw, outputf})
	writer.SetThreads(1)

	channel := writer.GoWriteFromChannel()

	c.Check(writer.IsAlive(), Equals, true)

	time.Sleep(100 * time.Microsecond)
	c.Check(writer.IsAlive(), Equals, true)
	close(channel)

	time.Sleep(100 * time.Microsecond)
	c.Check(writer.IsAlive(), Equals, false)
}

func (s *WriterSuite) TestWriteProcessor(c *C) {
	output := new(MockOutput)
	output.Return = true

	processor := new(MockProcessor)
	processor.Value = 1

	writer := NewWriter()
	writer.SetOutputs([]intf.Output{output})
	writer.SetProcessors([]intf.PostProcessor{processor})
	writer.SetThreads(1)

	channel := writer.GoWriteFromChannel()
	go func(channel chan map[string]string) {
		channel <- map[string]string{"foo": "1"}
		close(channel)
	}(channel)

	for {
		time.Sleep(100 * time.Microsecond)
		if !writer.IsAlive() {
			break
		}
	}

	c.Check(output.Count, Equals, 2)
}

type MockOutput struct {
	Count  int
	Return bool
}

func (self *MockOutput) PutRecord(record map[string]string) bool {
	number, _ := strconv.Atoi(record["foo"])
	self.Count += number
	return self.Return
}

type MockProcessor struct {
	Value int
}

func (self *MockProcessor) Do(record map[string]string) {
	number, _ := strconv.Atoi(record["foo"])
	record["foo"] = fmt.Sprintf("%d", number+self.Value)
}
