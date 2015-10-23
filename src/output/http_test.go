package output

import (
	"io/ioutil"
	"net/http"

	"github.com/mcuadros/harvester/src/intf"
	. "gopkg.in/check.v1"
)

type HTTPSuite struct{}

var _ = Suite(&HTTPSuite{})

func (s *HTTPSuite) TestGetRecordDefault(c *C) {
	config := HTTPConfig{
		Url: "http://localhost:9200/foo/a",
	}

	output := NewHTTP(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/a", "application/x-www-form-urlencoded", "POST", "foo=bar")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *HTTPSuite) TestGetRecordTemplatedURL(c *C) {
	config := HTTPConfig{
		Url: `http://localhost:9200/%{foo}/a`,
	}

	output := NewHTTP(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/bar/a", "application/x-www-form-urlencoded", "POST", "foo=bar")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *HTTPSuite) TestGetRecordPOStJson(c *C) {
	config := HTTPConfig{
		Url:         "http://localhost:9200/foo/b",
		Format:      "json",
		Method:      "POST",
		ContentType: "foo/bar",
	}

	output := NewHTTP(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/b", "foo/bar", "POST", "{\n     \"foo\": \"bar\"\n }")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *HTTPSuite) TestGetRecordPOSTForm(c *C) {
	config := HTTPConfig{
		Url:         "http://localhost:9200/foo/c",
		Format:      "form",
		Method:      "POST",
		ContentType: "foo/bar",
	}

	output := NewHTTP(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/c", "foo/bar", "POST", "foo=bar")
	c.Assert(output.PutRecord(record), Equals, true)
}

func (s *HTTPSuite) TestGetRecordGETFrom(c *C) {
	config := HTTPConfig{
		Url:         "http://localhost:9200/foo/d",
		Format:      "form",
		Method:      "POST",
		ContentType: "foo/bar",
	}

	output := NewHTTP(&config)
	record := intf.Record{"foo": "bar"}

	go dummyServer(c, ":9200", "/foo/d", "foo/bar", "POST", "foo=bar")
	c.Assert(output.PutRecord(record), Equals, true)
}

func dummyServer(c *C, server, path, contentType, method, equals string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == method {
			body, _ := ioutil.ReadAll(r.Body)

			c.Assert(string(body), Equals, equals)
			c.Assert(contentType, Equals, r.Header.Get("Content-Type"))

			w.WriteHeader(http.StatusCreated)
			w.Write(body)
		}
	}

	http.HandleFunc(path, handler)
	http.ListenAndServe(server, nil)
}
