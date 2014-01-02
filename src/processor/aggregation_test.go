package processor

import (
	"fmt"
	. "harvesterd/intf"
)

import . "launchpad.net/gocheck"

type AggregationSuite struct{}

var _ = Suite(&AggregationSuite{})

func (s *AggregationSuite) TestDoCount(c *C) {
	config := AggregationConfig{Metrics: "count(foo)"}
	processor := NewAggregation(&config)

	processor.Do(Record{"foo": "bar"})
	processor.Do(Record{"foo": "bar"})

	record := Record{"foo": "qux"}
	processor.Do(record)

	fmt.Println(record)
}
