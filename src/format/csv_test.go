package format

import (
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FormatCSVSuite struct{}

var _ = Suite(&FormatCSVSuite{})

func (s *FormatCSVSuite) TestGetRecordDefaultSettings(c *C) {
	config := CSVConfig{Fields: "foo,bar"}

	format := NewCSV(&config)

	record := format.Parse("baz,\"qux  \"")
	c.Assert(record["foo"], Equals, "baz")
	c.Assert(record["bar"], Equals, "qux  ")
}

func (s *FormatCSVSuite) TestGetRecordCustomeSettings(c *C) {
	config := CSVConfig{Fields: "foo,bar", Trim: true, Quote: '\'', Separator: ';'}

	format := NewCSV(&config)

	record := format.Parse("baz;'qux  '")
	c.Assert(record["foo"], Equals, "baz")
	c.Assert(record["bar"], Equals, "qux")
}

func (s *FormatCSVSuite) TestGetRecordCustomeSettingsNotQuoted(c *C) {
	config := CSVConfig{Fields: "foo,bar", NotQuoted: true}

	format := NewCSV(&config)

	record := format.Parse("baz,\"qux  \"")
	c.Assert(record["foo"], Equals, "baz")
	c.Assert(record["bar"], Equals, "\"qux  \"")
}

func (s *FormatCSVSuite) TestGetRecordWithFormat(c *C) {
	config := CSVConfig{
		Fields:    "foo,bar",
		Format:    "(float)foo",
		NotQuoted: true,
	}

	format := NewCSV(&config)

	record := format.Parse("1.2,\"qux  \"")
	c.Assert(record["foo"], Equals, 1.2)
	c.Assert(record["bar"], Equals, "\"qux  \"")
}

func (s *FormatCSVSuite) TestGetRecordWithDiscard(c *C) {
	config := CSVConfig{
		Fields: "_,foo",
	}

	format := NewCSV(&config)

	record := format.Parse("1.2,a")
	c.Assert(record, HasLen, 1)
	c.Assert(record["foo"], Equals, "a")
}
