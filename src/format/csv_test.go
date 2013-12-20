package format

import (
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type FormatCSVSuite struct{}

var _ = Suite(&FormatCSVSuite{})

func (s *FormatCSVSuite) TestGetRecordDefaultSettings(c *C) {
	config := CSVConfig{Fields: "foo,bar"}

	format := NewCSV(&config)

	record := format.Parse("baz,\"qux  \"")
	c.Check(record["foo"], Equals, "baz")
	c.Check(record["bar"], Equals, "qux  ")
}

func (s *FormatCSVSuite) TestGetRecordCustomeSettings(c *C) {
	config := CSVConfig{Fields: "foo,bar", Trim: true, Quote: '\'', Separator: ';'}

	format := NewCSV(&config)

	record := format.Parse("baz;'qux  '")
	c.Check(record["foo"], Equals, "baz")
	c.Check(record["bar"], Equals, "qux")
}

func (s *FormatCSVSuite) TestGetRecordCustomeSettingsNotQuoted(c *C) {
	config := CSVConfig{Fields: "foo,bar", NotQuoted: true}

	format := NewCSV(&config)

	record := format.Parse("baz,\"qux  \"")
	c.Check(record["foo"], Equals, "baz")
	c.Check(record["bar"], Equals, "\"qux  \"")
}
