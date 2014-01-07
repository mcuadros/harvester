package processor

import (
	"harvesterd/intf"
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type AnonymizeSuite struct{}

var _ = Suite(&AnonymizeSuite{})

func (s *AnonymizeSuite) TestDoDefault(c *C) {
	config := AnonymizeConfig{Fields: "foo"}
	processor := NewAnonymize(&config)

	recordA := intf.Record{"foo": "bar"}
	processor.Do(recordA)

	recordB := intf.Record{"foo": "qux"}
	processor.Do(recordB)

	c.Assert(recordA["foo"], Equals, "37b51d194a7513e45b56f6524f2d51f2")
	c.Assert(recordB["foo"], Equals, "d85b1213473c2fd7c2045020a6b9c62b")
}

func (s *AnonymizeSuite) TestDoSHA1(c *C) {
	config := AnonymizeConfig{Fields: "foo", Hash: "sha1"}
	processor := NewAnonymize(&config)

	recordA := intf.Record{"foo": "bar"}
	processor.Do(recordA)

	recordB := intf.Record{"foo": "qux"}
	processor.Do(recordB)

	c.Assert(recordA["foo"], Equals, "62cdb7020ff920e5aa642c3d4066950dd1f01f4d")
	c.Assert(recordB["foo"], Equals, "b54ba7f5621240d403f06815f7246006ef8c7d43")
}

func (s *AnonymizeSuite) TestDoEmailSupport(c *C) {
	config := AnonymizeConfig{Fields: "foo", EmailSupport: true}
	processor := NewAnonymize(&config)

	record := intf.Record{"foo": "bar@qux"}
	processor.Do(record)

	c.Assert(record["foo"], Equals, "37b51d194a7513e45b56f6524f2d51f2@qux")
}
