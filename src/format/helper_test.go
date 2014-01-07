package format

import (
	"time"
)
import . "launchpad.net/gocheck"

type FormatHelperSuite struct{}

var _ = Suite(&FormatHelperSuite{})

func (s *FormatHelperSuite) TestFormatInt(c *C) {
	helper := NewFormatHelper("(int)bar")

	c.Assert(helper.Format("bar", "1"), Equals, 1)
	c.Assert(helper.Format("bar", "foo"), Equals, 0)
	c.Assert(helper.Format("bar", " 1"), Equals, 1)
}

func (s *FormatHelperSuite) TestFormatFloat(c *C) {
	helper := NewFormatHelper("(float)bar")

	c.Assert(helper.Format("bar", "1"), Equals, 1.0)
	c.Assert(helper.Format("bar", "foo"), Equals, 0.0)
	c.Assert(helper.Format("bar", " 1"), Equals, 1.0)
	c.Assert(helper.Format("bar", "1,0"), Equals, 1.0)
	c.Assert(helper.Format("bar", "1.0"), Equals, 1.0)
}

func (s *FormatHelperSuite) TestFormatBool(c *C) {
	helper := NewFormatHelper("(bool)bar")

	c.Assert(helper.Format("bar", "foo"), Equals, nil)
	c.Assert(helper.Format("bar", " 1"), Equals, true)
	c.Assert(helper.Format("bar", "1"), Equals, true)
	c.Assert(helper.Format("bar", "T"), Equals, true)
	c.Assert(helper.Format("bar", "true"), Equals, true)
	c.Assert(helper.Format("bar", "0"), Equals, false)
	c.Assert(helper.Format("bar", "F"), Equals, false)
	c.Assert(helper.Format("bar", "false"), Equals, false)
}

func (s *FormatHelperSuite) TestFormatString(c *C) {
	helper := NewFormatHelper("(string)bar")

	c.Assert(helper.Format("bar", " 1 "), Equals, "1")
}

func (s *FormatHelperSuite) TestFormatTime(c *C) {
	helper := NewFormatHelper("(time:\"Jan 2, 2006 at 3:04pm (MST)\")bar,(int)foo")

	result := helper.Format("bar", "Jul 9, 2012 at 5:02am (CEST)")
	c.Assert(result.(time.Time).Unix(), Equals, int64(1341802920))
}
