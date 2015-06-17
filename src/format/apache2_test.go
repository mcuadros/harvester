package format

import (
	"time"
)

import . "gopkg.in/check.v1"

type Apache2Suite struct{}

var _ = Suite(&Apache2Suite{})
var apache2CommonExample = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326"
var apache2CombinedExample = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326 \"http://www.example.com/start.html\" \"Mozilla/4.08 [en] (Win98; I ;Nav)\""
var apache2ErrorExample = "[Wed Jan 01 20:22:10 2014] [error] [mod_pagespeed 0.10.22.4-1633 @31413] /var/www/mod_pagespeed/cache/5ooIhZKDe5hPOR1Lv9I3.outputlock:0: failed to stat (code=2 No such file or directory)"

func (s *Apache2Suite) TestGetRecordCommon(c *C) {
	config := Apache2Config{Type: "common"}

	format := NewApache2(&config)

	record := format.Parse(apache2CommonExample)
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

func (s *Apache2Suite) TestGetRecordCombined(c *C) {
	config := Apache2Config{Type: "combined"}

	format := NewApache2(&config)

	record := format.Parse(apache2CombinedExample)
	c.Assert(record["host"], Equals, "127.0.0.1")
	c.Assert(record["identd"], Equals, "-")
	c.Assert(record["user"], Equals, "frank")
	c.Assert(record["time"].(time.Time).String(), Equals, "2000-10-10 13:55:36 -0700 -0700")
	c.Assert(record["method"], Equals, "GET")
	c.Assert(record["path"], Equals, "/apache_pb.gif")
	c.Assert(record["version"], Equals, "HTTP/1.0")
	c.Assert(record["status"], Equals, 200)
	c.Assert(record["size"], Equals, 2326)
	c.Assert(record["referer"], Equals, "http://www.example.com/start.html")
	c.Assert(record["agent"], Equals, "Mozilla/4.08 [en] (Win98; I ;Nav)")

}

func (s *Apache2Suite) TestGetRecordError(c *C) {
	config := Apache2Config{Type: "error"}

	format := NewApache2(&config)

	record := format.Parse(apache2ErrorExample)
	c.Assert(record["time"].(time.Time).String(), Equals, "2014-01-01 20:22:10 +0000 UTC")
	c.Assert(record["severity"], Equals, "error")
	c.Assert(record["identifier"], Equals, "mod_pagespeed 0.10.22.4-1633 @31413")
	c.Assert(record["message"], Equals, "/var/www/mod_pagespeed/cache/5ooIhZKDe5hPOR1Lv9I3.outputlock:0: failed to stat (code=2 No such file or directory)")

}
