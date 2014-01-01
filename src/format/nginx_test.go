package format

import . "launchpad.net/gocheck"

type NginxSuite struct{}

var _ = Suite(&NginxSuite{})
var nginxCombinedExample = "127.0.0.1 - frank [10/Oct/2000:13:55:36 -0700] \"GET /apache_pb.gif HTTP/1.0\" 200 2326 \"http://www.example.com/start.html\" \"Mozilla/4.08 [en] (Win98; I ;Nav)\""
var nginxErrorExample = "2011/06/10 13:30:10 [error] 23263#0: *1 directory index of \"/var/www/ssl/\" is forbidden, client: 127.0.0.1, server: example.com, request: \"GET /foo.gif HTTP/1.1\", host: \"front-1.example.com\""

func (s *NginxSuite) TestGetRecordError(c *C) {
	config := NginxConfig{Type: "error"}

	format := NewNginx(&config)

	record := format.Parse(nginxErrorExample)
	c.Check(record["time"], Equals, "2011/06/10 13:30:10")
	c.Check(record["severity"], Equals, "error")
	c.Check(record["message"], Equals, "*1 directory index of \"/var/www/ssl/\" is forbidden")
	c.Check(record["client"], Equals, "127.0.0.1")
	c.Check(record["server"], Equals, "example.com")
	c.Check(record["method"], Equals, "GET")
	c.Check(record["path"], Equals, "/foo.gif")
	c.Check(record["version"], Equals, "HTTP/1.1")
	c.Check(record["host"], Equals, "front-1.example.com")
}

func (s *NginxSuite) TestGetRecordCombined(c *C) {
	config := Apache2Config{Type: "combined"}

	format := NewApache2(&config)

	record := format.Parse(nginxCombinedExample)
	c.Check(record["host"], Equals, "127.0.0.1")
	c.Check(record["identd"], Equals, "-")
	c.Check(record["user"], Equals, "frank")
	c.Check(record["time"], Equals, "10/Oct/2000:13:55:36 -0700")
	c.Check(record["method"], Equals, "GET")
	c.Check(record["path"], Equals, "/apache_pb.gif")
	c.Check(record["version"], Equals, "HTTP/1.0")
	c.Check(record["status"], Equals, "200")
	c.Check(record["size"], Equals, "2326")
}
