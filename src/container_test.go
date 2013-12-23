package collector

import (
	"collector/format"
	"collector/input"
	"collector/output"
)

import . "launchpad.net/gocheck"

type ContainerSuite struct{}

var _ = Suite(&ContainerSuite{})

func (s *ContainerSuite) TestGetInput(c *C) {
	var raw = string(`
		[format-csv "myformat"]
		fields = foo

		[input-tail "bar"]
		file = foo
		format = myformat

		[input-file "foo"]
		pattern = foo
		format = myformat
	`)

	GetConfig().Load(raw)

	c.Check(GetContainer().GetInput("foo"), FitsTypeOf, &input.File{})
	c.Check(GetContainer().GetInput("bar"), FitsTypeOf, &input.Tail{})
}

func (s *ContainerSuite) TestGetFormat(c *C) {
	var raw = string(`
		[input-tail "bar"]
		file = foo
		format = myformat

		[format-csv "foo"]
		fields = foo

		[format-regexp "bar"]
		pattern = foo
	`)

	GetConfig().Load(raw)

	c.Check(GetContainer().GetFormat("foo"), FitsTypeOf, &format.CSV{})
	c.Check(GetContainer().GetFormat("bar"), FitsTypeOf, &format.RegExp{})
}

func (s *ContainerSuite) TestGetOutput(c *C) {
	var raw = string(`
		[output-elasticsearch "bar"]
		host = foo

		[output-mongo "foo"]
		url = mongodb://localhost

		[output-dummy "qux"]
		print = true
	`)

	GetConfig().Load(raw)
	c.Check(GetContainer().GetOutput("qux"), FitsTypeOf, &output.Dummy{})
	c.Check(GetContainer().GetOutput("foo"), FitsTypeOf, &output.Mongo{})
	c.Check(GetContainer().GetOutput("bar"), FitsTypeOf, &output.Elasticsearch{})
}

func (s *ContainerSuite) TestGetReader(c *C) {
	var raw = string(`
		[input-tail "bar"]
		file = foo
		format = myformat

		[reader]
		input = bar

	`)

	GetConfig().Load(raw)

	c.Check(GetContainer().GetReader(), FitsTypeOf, &Reader{})
}

func (s *ContainerSuite) TestGetWriter(c *C) {
	var raw = string(`
		[output-elasticsearch "bar"]
		host = foo

		[writer]
		output = bar

	`)

	GetConfig().Load(raw)

	c.Check(GetContainer().GetWriter(), FitsTypeOf, &Writer{})
}
