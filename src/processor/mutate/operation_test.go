package mutate

import . "gopkg.in/check.v1"

type OperationSuite struct{}

var _ = Suite(&OperationSuite{})

func (s *OperationSuite) TestApply(c *C) {
	ops := []Operation{
		Operation{
			Id:     CAST,
			Field:  []string{"fieldA"},
			Params: []string{"int"},
		},
		Operation{
			Id:     CAST,
			Field:  []string{"fieldB", "*", "subfieldB1"},
			Params: []string{"date", "2006-01-02T15:04:05"},
		},
	}

	var r map[string]interface{}
	var err error

	r = map[string]interface{}{}
	err = ops[0].Apply(r)
	c.Assert(err.Error(), Equals, "mutate: key `fieldA` was not found in subtree `map[]`")
	err = ops[1].Apply(r)
	c.Assert(err.Error(), Equals, "mutate: key `fieldB` was not found in subtree `map[]`")

	r = map[string]interface{}{
		"fieldZ": "blah",
	}
	err = ops[0].Apply(r)
	c.Assert(err.Error(), Equals, "mutate: key `fieldA` was not found in subtree `map[fieldZ:blah]`")
	err = ops[1].Apply(r)
	c.Assert(err.Error(), Equals, "mutate: key `fieldB` was not found in subtree `map[fieldZ:blah]`")

	r = map[string]interface{}{
		"fieldA": "not int",
		"fieldB": []interface{}{
			map[string]interface{}{
				"subfieldB1": "not date",
			},
		},
		"fieldZ": "blah",
	}
	err = ops[0].Apply(r)
	c.Assert(err, NotNil)
	c.Assert(r["fieldA"], Equals, "not int")
	err = ops[1].Apply(r)
	c.Assert(err, NotNil)
	c.Assert((r["fieldB"].([]interface{})[0]).(map[string]interface{})["subfieldB1"], Equals, "not date")

	r = map[string]interface{}{
		"fieldA": "1234",
		"fieldB": []interface{}{
			map[string]interface{}{
				"subfieldB1": "2013-09-16T09:15:30",
			},
			map[string]interface{}{
				"subfieldB1": "2015-09-16T09:15:30",
			},
		},
		"fieldZ": "blah",
	}
	err = ops[0].Apply(r)
	c.Assert(err, IsNil)
	c.Assert(r["fieldA"], Not(Equals), "1234")
	err = ops[1].Apply(r)
	c.Assert(err, IsNil)
	c.Assert((r["fieldB"].([]interface{})[0]).(map[string]interface{})["subfieldB1"], Not(Equals), "2013-09-16T09:15:30")
	c.Assert((r["fieldB"].([]interface{})[1]).(map[string]interface{})["subfieldB1"], Not(Equals), "2015-09-16T09:15:30")
}
