package collector

import (
	"collector/intf"
	"fmt"
	"strconv"
	"time"
)

import . "launchpad.net/gocheck"

type WriterSuite struct{}

var _ = Suite(&WriterSuite{})

func (s *WriterSuite) TestWriteFromChannelSingleOutput(c *C) {
	channel := make(chan map[string]string, 1)
	output := new(MockOutput)
	output.Return = true
	outputs := []intf.Output{output}

	writer := NewWriter(outputs)
	go writer.WriteFromChannel(channel)
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
	channel := make(chan map[string]string, 1)

	outputw := new(MockOutput)
	outputw.Return = true
	outputf := new(MockOutput)
	outputf.Return = false
	outputs := []intf.Output{outputw, outputf}

	writer := NewWriter(outputs)
	go writer.WriteFromChannel(channel)
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

type MockOutput struct {
	Count  int
	Return bool
}

func (self *MockOutput) PutRecord(record map[string]string) bool {
	number, _ := strconv.Atoi(record["foo"])
	self.Count += number
	return self.Return
}
