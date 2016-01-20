package harvester

import (
	"testing"
)

import . "gopkg.in/check.v1"

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

	c.Assert(definition, HasLen, 19)
	c.Assert(definition[0].Name, Equals, "logger")
	c.Assert(definition[0].AllowMultiple, Equals, false)
	c.Assert(definition[0].Fields[0].Name, Equals, "level")
	c.Assert(definition[0].Fields[0].Default, Equals, "")

	c.Assert(definition[1].Name, Equals, "writer")
	c.Assert(definition[1].AllowMultiple, Equals, true)
}
