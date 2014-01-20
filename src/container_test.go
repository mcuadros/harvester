package harvesterd

import (
	"harvesterd/format"
	"harvesterd/input"
	"harvesterd/output"
	"harvesterd/processor"
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

	c.Assert(GetContainer().GetInput("foo"), FitsTypeOf, &input.File{})
	c.Assert(GetContainer().GetInput("bar"), FitsTypeOf, &input.Tail{})
}

func (s *ContainerSuite) TestGetFormat(c *C) {
	var raw = string(`
		[input-tail "bar"]
		file = foo
		format = myformat

		[format-csv "csv"]
		fields = foo

		[format-regexp "regexp"]
		pattern = foo

		[format-apache2 "apache2"]
		type = common

		[format-nginx "nginx"]
		type = common

		[format-json "json"]
		empty = foo
	`)

	GetConfig().Load(raw)

	c.Assert(GetContainer().GetFormat("json"), FitsTypeOf, &format.JSON{})
	c.Assert(GetContainer().GetFormat("csv"), FitsTypeOf, &format.CSV{})
	c.Assert(GetContainer().GetFormat("regexp"), FitsTypeOf, &format.RegExp{})
	c.Assert(GetContainer().GetFormat("apache2"), FitsTypeOf, &format.Apache2{})
	c.Assert(GetContainer().GetFormat("nginx"), FitsTypeOf, &format.Nginx{})
}

func (s *ContainerSuite) TestGetOutput(c *C) {
	var raw = string(`
		[output-elasticsearch "es"]
		host = foo

		[output-mongo "mongo"]
		url = mongodb://localhost

		[output-http "http"]
		url = http://localhost

		[output-dummy "dummy"]
		print = true
	`)

	GetConfig().Load(raw)
	c.Assert(GetContainer().GetOutput("dummy"), FitsTypeOf, &output.Dummy{})
	c.Assert(GetContainer().GetOutput("mongo"), FitsTypeOf, &output.Mongo{})
	c.Assert(GetContainer().GetOutput("es"), FitsTypeOf, &output.Elasticsearch{})
	c.Assert(GetContainer().GetOutput("http"), FitsTypeOf, &output.HTTP{})

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

	c.Assert(GetContainer().GetReader(""), FitsTypeOf, &Reader{})
}

func (s *ContainerSuite) TestGetWriter(c *C) {
	var raw = string(`
		[output-elasticsearch "bar"]
		host = foo

		[writer "bar"]
		output = bar
	`)

	GetConfig().Load(raw)

	c.Assert(GetContainer().GetWriter(""), FitsTypeOf, &Writer{})
}

func (s *ContainerSuite) TestGetWriterGroup(c *C) {
	var raw = string(`
		[writer "foo"]
		output = bar

		[writer "bar"]
		output = bar
	`)

	GetConfig().Load(raw)

	c.Assert(GetContainer().GetWriterGroup(), FitsTypeOf, &WriterGroup{})
}

func (s *ContainerSuite) TestGetPostProcessor(c *C) {
	var raw = string(`
		[processor-anonymize "qux"]
		fields = true

		[processor-metrics "bar"]
		metrics = (terms)foo
	`)

	GetConfig().Load(raw)
	c.Assert(GetContainer().GetPostProcessor("qux"), FitsTypeOf, &processor.Anonymize{})
	c.Assert(GetContainer().GetPostProcessor("bar"), FitsTypeOf, &processor.Metrics{})
}
