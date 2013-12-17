package collector

import (
	"testing"
)

func TestReadIntoChannel(t *testing.T) {
	channel := make(chan map[string]string, 1)

	format := new(MockFormat)
	input := new(MockInput)

	config := NewReader(format, input)
	go config.ReadIntoChannel(channel)

	count := 0
	for record := range channel {
		if record["foo"] != "foo" {
			t.Errorf("FAIL")
		}
		count++
	}

	if count != 4 {
		t.Errorf("FAIL")
	}
}

type MockFormat struct {
}

func (self *MockFormat) Parse(line string) map[string]string {
	record := make(map[string]string)
	record["foo"] = line

	return record
}

type MockInput struct {
	current int
}

func (self *MockInput) GetLine() string {
	self.current++

	return string("foo")
}

func (self *MockInput) IsEOF() bool {
	if self.current > 3 {
		return true
	}

	return false
}
