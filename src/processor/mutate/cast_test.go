package mutate

import (
	"time"

	. "gopkg.in/check.v1"
)

type CastSuite struct{}

var _ = Suite(&CastSuite{})

func (s *CastSuite) TestCastINT(c *C) {
	var value interface{}
	var err error

	value = "1234"
	value, err = cast(value, INT)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, 1234)

	value = 1234
	value, err = cast(value, INT)
	c.Assert(err, IsNil)
	c.Assert(value, Equals, 1234)

	value = "int 1234 that will be stripped"
	value, err = cast(value, INT, "strip")
	c.Assert(err, IsNil)
	c.Assert(value, Equals, 1234)

	value = "not int"
	value, err = cast(value, INT)
	c.Assert(err, NotNil)
	c.Assert(value, Equals, "not int")

	value = struct{}{}
	value, err = cast(value, INT)
	c.Assert(err, NotNil)
	c.Assert(value, Equals, struct{}{})
}

func (s *CastSuite) TestCastDATE(c *C) {
	var value interface{}
	var err error
	var t0 time.Time
	var t time.Time

	value = "2015-09-16T09:15:30"
	value, err = cast(value, DATE, "2006-01-02T15:04:05")
	c.Assert(err, IsNil)
	t = value.(time.Time)
	c.Assert(t.Unix(), Equals, int64(1442394930))

	value = "Sep 15"
	value, err = cast(value, DATE, "2006-01-02T15:04:05", "Jan 06")
	c.Assert(err, IsNil)
	t = value.(time.Time)
	c.Assert(t.Unix(), Equals, int64(1441065600))

	value = 1442394930
	value, err = cast(value, DATE)
	c.Assert(err, IsNil)
	t = value.(time.Time)
	c.Assert(t.Unix(), Equals, int64(1442394930))

	value = "Present"
	t0 = time.Now()
	value, err = cast(value, DATE, "2006-01-02T15:04:05", "present")
	c.Assert(err, IsNil)
	t = value.(time.Time)
	c.Assert(t.Unix() >= t0.Unix(), Equals, true)
	c.Assert(t.Unix() <= time.Now().Unix(), Equals, true)

	value = "bad nullable date"
	value, err = cast(value, DATE, "2006-01-02T15:04:05", "null")
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	value = "bad date"
	value, err = cast(value, DATE, "2006-01-02T15:04:05")
	c.Assert(err, NotNil)
	c.Assert(value, Equals, "bad date")

	value = struct{}{}
	value, err = cast(value, DATE, "2006-01-02T15:04:05")
	c.Assert(err, NotNil)
	c.Assert(value, Equals, struct{}{})
}

func (s *CastSuite) TestCastUnknown(c *C) {
	var value interface{}
	var err error

	value = "irrelevant test value"
	value, err = cast(value, "unsupported-function")
	c.Assert(err, NotNil)
	c.Assert(value, Equals, "irrelevant test value")
}
