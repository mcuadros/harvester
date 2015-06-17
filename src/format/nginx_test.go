package format

import (
	"time"

	. "gopkg.in/check.v1"
)

type NginxSuite struct{}

var _ = Suite(&NginxSuite{})
var nginxCombinedExample = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326 \"http://www.example.com/start.html\" \"Mozilla/4.08 [en] (Win98; I ;Nav)\""
var nginxErrorExample = "2011/06/10 13:30:10 [error] 23263#0: *1 directory index of \"/var/www/ssl/\" is forbidden, client: 127.0.0.1, server: example.com, request: \"GET /foo.gif HTTP/1.1\", host: \"front-1.example.com\""

func (s *NginxSuite) TestGetRecordError(c *C) {
	config := NginxConfig{Type: "error"}

	format := NewNginx(&config)

	record := format.Parse(nginxErrorExample)
	c.Assert(record["time"].(time.Time).String(), Equals, "2011-06-10 13:30:10 +0000 UTC")
	c.Assert(record["severity"], Equals, "error")
	c.Assert(record["message"], Equals, "*1 directory index of \"/var/www/ssl/\" is forbidden")
	c.Assert(record["client"], Equals, "127.0.0.1")
	c.Assert(record["server"], Equals, "example.com")
	c.Assert(record["method"], Equals, "GET")
	c.Assert(record["path"], Equals, "/foo.gif")
	c.Assert(record["version"], Equals, "HTTP/1.1")
	c.Assert(record["host"], Equals, "front-1.example.com")
}

func (s *NginxSuite) TestGetRecordCombined(c *C) {
	config := NginxConfig{Type: "combined"}

	format := NewNginx(&config)

	record := format.Parse(nginxCombinedExample)
	c.Assert(record["host"], Equals, "127.0.0.1")
	c.Assert(record["identd"], Equals, "-")
	c.Assert(record["user"], Equals, "frank")
	c.Assert(record["time"].(time.Time).String(), Equals, "2000-10-10 13:55:36 -0700 -0700")
	c.Assert(record["method"], Equals, "GET")
	c.Assert(record["path"], Equals, "/apache_pb.gif")
	c.Assert(record["version"], Equals, "HTTP/1.0")
	c.Assert(record["status"], Equals, 200)
	c.Assert(record["size"], Equals, 2326)
}
