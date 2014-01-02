package format

import . "launchpad.net/gocheck"

type FormatRegExpSuite struct{}

var _ = Suite(&FormatRegExpSuite{})

func (s *FormatRegExpSuite) TestGetRecord(c *C) {
	config := RegExpConfig{Pattern: "^(?P<foo>\\S+) (?P<bar>\\S+)$"}

	format := NewRegExp(&config)

	record := format.Parse("qux baz")
	c.Check(record, HasLen, 2)
	c.Check(record["foo"], Equals, "qux")
	c.Check(record["bar"], Equals, "baz")
}

func (s *FormatRegExpSuite) TestGetRecordWithFormat(c *C) {
	config := RegExpConfig{
		Pattern: "^(?P<foo>\\S+) (?P<bar>\\S+)$",
		Format:  "(int)foo",
	}

	format := NewRegExp(&config)

	record := format.Parse("1 baz")
	c.Check(record, HasLen, 2)
	c.Check(record["foo"], Equals, 1)
	c.Check(record["bar"], Equals, "baz")
}
