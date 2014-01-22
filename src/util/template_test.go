package util

import (
	"harvesterd/intf"
	"testing"
)

import . "launchpad.net/gocheck"

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type TemplateSuite struct{}

var _ = Suite(&TemplateSuite{})

func (s *TemplateSuite) TestApplyWithString(c *C) {
	template := NewTemplate(`foo %{bar} --- %{qux}`)

	record := intf.Record{"bar": "qux"}
	c.Assert(template.Apply(record), Equals, "foo qux --- ")
}

func (s *TemplateSuite) TestApplyWithMap(c *C) {
	template := NewTemplate(`foo %{bar.foo}`)

	record := intf.Record{"bar": map[string]interface{}{"foo": "qux"}}
	c.Assert(template.Apply(record), Equals, "foo qux")
}

func (s *TemplateSuite) TestApplyWithInt(c *C) {
	template := NewTemplate(`%{int} %{int8} %{int16} %{int32} %{int64}`)

	record := intf.Record{
		"int":   int(1),
		"int8":  int8(8),
		"int16": int16(16),
		"int32": int32(32),
		"int64": int64(64),
	}
	c.Assert(template.Apply(record), Equals, "1 8 16 32 64")
}

func (s *TemplateSuite) TestApplyWithUInt(c *C) {
	template := NewTemplate(`%{uint} %{uint8} %{uint16} %{uint32} %{uint64}`)

	record := intf.Record{
		"uint":   uint(1),
		"uint8":  uint8(8),
		"uint16": uint16(16),
		"uint32": uint32(32),
		"uint64": uint64(64),
	}
	c.Assert(template.Apply(record), Equals, "1 8 16 32 64")
}

func (s *TemplateSuite) TestApplyWithBool(c *C) {
	template := NewTemplate(`%{bool}`)

	record := intf.Record{"bool": true}
	c.Assert(template.Apply(record), Equals, "true")
}

func (s *TemplateSuite) TestApplyWithFloat(c *C) {
	template := NewTemplate(`%{float32} %{float64}`)

	record := intf.Record{
		"float32": float32(32.32),
		"float64": float64(64.64),
	}
	c.Assert(template.Apply(record), Equals, "32.32 64.64")
}
