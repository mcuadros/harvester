package output

import (
	"testing"

	"github.com/mcuadros/harvesterd/src/intf"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ElasticsearchSuite struct{}

var _ = Suite(&ElasticsearchSuite{})

func (s *ElasticsearchSuite) TestGetRecordDefault(c *C) {
	config := ElasticsearchConfig{Index: "foo", Type: "qux"}

	output := NewElasticsearch(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/qux", "application/json", "POST", "{\n     \"foo\": \"bar\"\n }")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *ElasticsearchSuite) TestGetRecordDefaultField(c *C) {
	config := ElasticsearchConfig{Index: "foo", Type: "foo", IdField: "foo"}

	output := NewElasticsearch(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/foo", "application/json", "POST", "{\n     \"_id\": \"bar\",\n     \"foo\": \"bar\"\n }")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *ElasticsearchSuite) TestGetRecordConfig(c *C) {
	config := ElasticsearchConfig{Host: "127.0.0.1", Port: 9300, Index: "foo", Type: "bar"}

	output := NewElasticsearch(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9300", "/foo/bar", "application/json", "POST", "{\n     \"foo\": \"bar\"\n }")
	c.Assert(output.PutRecord(record), Equals, true)
}
