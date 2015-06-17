package harvesterd

import (
	"fmt"
	"harvesterd/intf"
	"strconv"
)

import . "gopkg.in/check.v1"

type ReaderSuite struct{}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) TestReadIntoChannelSingleInput(c *C) {
	recordsChan := make(RecordsChan, 1)
	closeChan := make(CloseChan, 1)

	input := new(MockInput)
	inputs := []intf.Input{input}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.SetChannels(recordsChan, closeChan)
	reader.GoRead()

	go func() {
		for _ = range closeChan {
			close(recordsChan)
		}
	}()

	count := 0
	for record := range recordsChan {
		c.Assert(record["line"], Equals, "foo")
		count++
	}

	reader.Teardown()
	c.Assert(count, Equals, 4)
	c.Assert(input.Finished, Equals, true)
}

func (s *ReaderSuite) TestReadIntoChannelWithProcessors(c *C) {
	recordsChan := make(RecordsChan, 1)
	closeChan := make(CloseChan, 1)

	input := new(MockInput)
	inputs := []intf.Input{input}

	processor := new(MockProcessor)
	processor.Value = 10

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.SetProcessors([]intf.PostProcessor{processor})
	reader.SetChannels(recordsChan, closeChan)
	reader.GoRead()

	go func() {
		for _ = range closeChan {
			close(recordsChan)
		}
	}()

	count := 0
	for record := range recordsChan {
		c.Assert(record["line"], Equals, "10")
		count++
	}

	reader.Teardown()
	c.Assert(count, Equals, 2)
	c.Assert(input.Finished, Equals, true)
}

func (s *ReaderSuite) TestReadIntoChannelMultipleInputs(c *C) {
	recordsChan := make(RecordsChan, 1)
	closeChan := make(CloseChan, 1)

	inputA := new(MockInput)
	inputB := new(MockInput)
	inputC := new(MockInput)
	inputD := new(MockInput)

	inputs := []intf.Input{inputA, inputB, inputC, inputD}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.SetChannels(recordsChan, closeChan)
	reader.GoRead()

	go func() {
		for _ = range closeChan {
			close(recordsChan)
		}
	}()

	count := 0
	for record := range recordsChan {
		c.Assert(record["line"], Equals, "foo")
		count++
	}

	reader.Teardown()
	c.Assert(count, Equals, 16)
	c.Assert(inputA.Finished, Equals, true)
	c.Assert(inputB.Finished, Equals, true)
	c.Assert(inputC.Finished, Equals, true)
	c.Assert(inputD.Finished, Equals, true)
}

type MockInput struct {
	Current  int
	Finished bool
}

func (self *MockInput) GetLine() string {
	self.Current++

	return string("foo")
}

func (self *MockInput) GetRecord() intf.Record {
	line := self.GetLine()

	return intf.Record{"line": line}
}

func (self *MockInput) IsEOF() bool {
	if self.Current > 3 {
		return true
	}

	return false
}
func (self *MockInput) Teardown() {
	self.Finished = true
}

type MockProcessor struct {
	Value int
	Count int
}

func (self *MockProcessor) SetChannel(recordsChan chan intf.Record) {
}

func (self *MockProcessor) Teardown() {
}

func (self *MockProcessor) Do(record intf.Record) bool {
	self.Count++

	number, _ := strconv.Atoi(record["line"].(string))
	record["line"] = fmt.Sprintf("%d", number+self.Value)

	if self.Count%2 == 0 {
		return true
	}

	return false
}
