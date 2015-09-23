package processor

import (
	"github.com/mcuadros/harvesterd/src/intf"
	"github.com/mcuadros/harvesterd/src/processor/mutate"
	. "gopkg.in/check.v1"
)

type MutateSuite struct{}

var _ = Suite(&MutateSuite{})

func (s *MutateSuite) TestMutateConfigParse(c *C) {
	mc := MutateConfig{
		Cast: []string{
			"fieldA int",
			"fieldB.*.subfieldB1 date 2006-01-02T15:04:05  'May 10th'    2015 '10th May 2015' '10-05-2015' ",
		},
	}
	ops := mc.ParseOperations()

	c.Assert(ops, HasLen, 2)
	c.Assert(ops[0].Id, Equals, mutate.CAST)
	c.Assert(ops[0].Field, HasLen, 1)
	c.Assert(ops[0].Params, HasLen, 1)
	c.Assert(ops[1].Id, Equals, mutate.CAST)
	c.Assert(ops[1].Field, HasLen, 3)
	c.Assert(ops[1].Params, HasLen, 6)
}

func (s *MutateSuite) TestNewMutate(c *C) {
	mc := MutateConfig{
		Cast: []string{
			"fieldA int",
			"fieldB.*.subfieldB1 date 2006-01-02T15:04:05",
		},
	}

	m := NewMutate(&mc)

	c.Assert(m, NotNil)
}

func (s *MutateSuite) TestNewMutateDo(c *C) {
	mc := MutateConfig{
		Cast: []string{
			"fieldA int",
			"fieldB.*.subfieldB1 date 2006-01-02T15:04:05",
		},
	}

	m := NewMutate(&mc)

	var r intf.Record
	var output bool

	r = intf.Record{}
	output = m.Do(r)
	c.Assert(output, Equals, true)

	r = intf.Record{
		"fieldZ": "blah",
	}
	output = m.Do(r)
	c.Assert(output, Equals, true)

	r = intf.Record{
		"fieldA": "not int",
		"fieldB": []interface{}{
			map[string]interface{}{
				"subfieldB1": "not date",
			},
		},
		"fieldZ": "blah",
	}
	output = m.Do(r)
	c.Assert(output, Equals, true)

	r = intf.Record{
		"fieldA": "1234",
		"fieldB": []interface{}{
			map[string]interface{}{
				"subfieldB1": "2015-09-16T09:15:30",
			},
		},
		"fieldZ": "blah",
	}
	output = m.Do(r)
	c.Assert(output, Equals, true)
}
