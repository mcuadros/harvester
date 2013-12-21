package collector

import (
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type CofigSuite struct{}

var _ = Suite(&CofigSuite{})

func (s *CofigSuite) TestFormat(c *C) {
	var raw = string(`
		[reader]
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
	`)

	GetConfig().Load(raw)

	c.Check(len(GetConfig().Reader.Input), Equals, 2)
	c.Check(GetConfig().Format_CSV["foo"].Fields, Equals, "foo")
}
