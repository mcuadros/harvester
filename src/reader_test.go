package harvesterd

import (
	"harvesterd/intf"
)

import . "launchpad.net/gocheck"

type ReaderSuite struct{}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) TestReadIntoChannelSingleInput(c *C) {
	channel := make(chan map[string]string, 1)
	input := new(MockInput)
	inputs := []intf.Input{input}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.GoReadIntoChannel(channel)

	count := 0
	for record := range channel {
		c.Check(record["line"], Equals, "foo")
		count++
	}

	reader.Finish()
	c.Check(count, Equals, 4)
	c.Check(input.Finished, Equals, true)
}

func (s *ReaderSuite) TestReadIntoChannelMultipleInputs(c *C) {
	channel := make(chan map[string]string, 1)
	inputA := new(MockInput)
	inputB := new(MockInput)
	inputC := new(MockInput)
	inputD := new(MockInput)

	inputs := []intf.Input{inputA, inputB, inputC, inputD}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.GoReadIntoChannel(channel)

	count := 0
	for record := range channel {
		c.Check(record["line"], Equals, "foo")
		count++
	}

	reader.Finish()
	c.Check(count, Equals, 16)
	c.Check(inputA.Finished, Equals, true)
	c.Check(inputB.Finished, Equals, true)
	c.Check(inputC.Finished, Equals, true)
	c.Check(inputD.Finished, Equals, true)

}

type MockInput struct {
	Current  int
	Finished bool
}

func (self *MockInput) GetLine() string {
	self.Current++

	return string("foo")
}

func (self *MockInput) GetRecord() map[string]string {
	line := self.GetLine()

	return map[string]string{"line": line}
}

func (self *MockInput) IsEOF() bool {
	if self.Current > 3 {
		return true
	}

	return false
}
func (self *MockInput) Finish() {
	self.Finished = true
}
