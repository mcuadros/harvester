package processor

import (
	. "harvesterd/intf"
	"runtime"
	"sync"
)

import . "launchpad.net/gocheck"

type MetricsSuite struct{}

var _ = Suite(&MetricsSuite{})

func (s *MetricsSuite) TestDoCount(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config := MetricsConfig{Metrics: "(terms)foo,(histogram)qux"}
	processor := NewMetrics(&config)

	var wait sync.WaitGroup
	var add = func() {
		for i := 0; i < 10000; i++ {
			processor.Do(Record{"foo": "bar", "qux": 1})
		}

		wait.Done()
	}

	count := 5
	for i := 0; i < count; i++ {
		go add()
	}

	wait.Add(count)
	wait.Wait()

	record := Record{"foo": "qux"}
	processor.Do(record)

	c.Check(record["foo"].(map[string]int)["bar"], Equals, 50000)
	c.Check(record["qux"].(map[string]interface{})["count"], Equals, int64(50000))
}
