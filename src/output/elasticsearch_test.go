package output

import (
	"harvesterd/intf"
	"io/ioutil"
	"net/http"
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ElasticsearchSuite struct{}

var _ = Suite(&ElasticsearchSuite{})

func (s *ElasticsearchSuite) TestGetRecord(c *C) {
	config := ElasticsearchConfig{Host: "localhost", Port: 9200, Index: "foo", Type: "bar"}

	output := NewElasticsearch(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(":9200", "/foo/bar")
	c.Assert(output.PutRecord(record), Equals, true)
}

func parrotHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)

		w.WriteHeader(http.StatusCreated)
		w.Write(body)
	}
}

func dummyServer(server, path string) {
	http.HandleFunc(path, parrotHandler)
	http.ListenAndServe(server, nil)
}
