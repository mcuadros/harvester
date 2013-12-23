package collector

import (
	"collector/intf"
)

import . "launchpad.net/gocheck"

type ReaderSuite struct{}

var _ = Suite(&ReaderSuite{})

func (s *ReaderSuite) TestReadIntoChannelSingleInput(c *C) {
	channel := make(chan map[string]string, 1)
	inputs := []intf.Input{new(MockInput)}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.GoReadIntoChannel(channel)

	count := 0
	for record := range channel {
		c.Check(record["line"], Equals, "foo")
		count++
	}

	c.Check(count, Equals, 4)
}

func (s *ReaderSuite) TestReadIntoChannelMultipleInputs(c *C) {
	channel := make(chan map[string]string, 1)
	inputs := []intf.Input{new(MockInput), new(MockInput), new(MockInput), new(MockInput)}

	reader := NewReader()
	reader.SetInputs(inputs)
	reader.GoReadIntoChannel(channel)

	count := 0
	for record := range channel {
		c.Check(record["line"], Equals, "foo")
		count++
	}

	c.Check(count, Equals, 16)
}

type MockInput struct {
	current int
}

func (self *MockInput) GetLine() string {
	self.current++

	return string("foo")
}

func (self *MockInput) GetRecord() map[string]string {
	line := self.GetLine()

	return map[string]string{"line": line}
}

func (self *MockInput) IsEOF() bool {
	if self.current > 3 {
		return true
	}

	return false
}
