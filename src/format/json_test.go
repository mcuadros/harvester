package format

import . "gopkg.in/check.v1"

type FormatJSONSuite struct{}

var _ = Suite(&FormatJSONSuite{})

func (s *FormatJSONSuite) TestGetRecord(c *C) {
	config := JSONConfig{}

	format := NewJSON(&config)

	record := format.Parse("{\"foo\":{\"foo\":\"bar\",\"bar\":\"qux\"}}")

	foo := record["foo"].(map[string]interface{})
	c.Assert(foo["foo"], Equals, "bar")
}
