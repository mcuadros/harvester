package harvesterd

import (
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CofigSuite struct{}

var _ = Suite(&CofigSuite{})

func (s *CofigSuite) TestLoad(c *C) {
	var raw = string(`
		[reader "foo"]
		input = bar
		input = foo

		[format-csv "foo"]
		fields = foo

		[input-tail "bar"]
		file = foo
		format = myformat

		[input-file "foo"]
		pattern = foo
		format = myformat

		[input-json "foo"]
	`)

	GetConfig().Load(raw)

	c.Assert(len(GetConfig().Reader["foo"].Input), Equals, 2)
	c.Assert(GetConfig().Format_CSV["foo"].Fields, Equals, "foo")
	c.Assert(GetConfig().Format_JSON, HasLen, 0)

}

func (s *CofigSuite) TestGetDescription(c *C) {
	definition := GetConfig().GetDescription()

	c.Assert(definition, HasLen, 15)

	c.Assert(definition[0].Name, Equals, "Logger")
	c.Assert(definition[0].AllowMultiple, Equals, false)
	c.Assert(definition[0].Fields[0].Name, Equals, "Level")
	c.Assert(definition[0].Fields[0].Default, Equals, "info")

	c.Assert(definition[1].Name, Equals, "Writer")
	c.Assert(definition[1].AllowMultiple, Equals, true)
}
