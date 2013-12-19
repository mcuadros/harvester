package input

import (
	"testing"
)

import . "launchpad.net/gocheck"

type MockFormat struct{}

func (s *MockFormat) Parse(line string) map[string]string {
	return map[string]string{"foo": "bar"}
}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type InputFileSuite struct{}

var _ = Suite(&InputFileSuite{})

func (s *InputFileSuite) TestGetRecord(c *C) {
	config := FileConfig{Pattern: "../../tests/resources/plain.a.txt"}

	file := NewFile(&config, new(MockFormat))
	c.Check(file.IsEOF(), Equals, false)

	record := file.GetRecord()
	c.Check(record["foo"], Equals, "bar")
}

func (s *InputFileSuite) TestGetLineWithSingleFile(c *C) {
	config := FileConfig{Pattern: "../../tests/resources/plain.a.txt"}

	file := NewFile(&config, new(MockFormat))
	c.Check(file.IsEOF(), Equals, false)

	testReader(c, file, 3)
}

func (s *InputFileSuite) TestGetLineWithPatternGlob(c *C) {
	config := FileConfig{Pattern: "../../tests/resources/plain.*.txt"}

	file := NewFile(&config, new(MockFormat))
	c.Check(file.IsEOF(), Equals, false)

	testReader(c, file, 6)
}

func testReader(c *C, file *File, count int) {
	for i := 0; i <= count; i++ {
		line := file.GetLine()
		if i < count {
			c.Check(line, Not(HasLen), 0)
		}

		if i >= count && len(line) != 0 {
			c.Check(line, HasLen, 0)
		}
	}

	c.Check(file.IsEOF(), Equals, true)
}
