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

	config := MetricsConfig{Metrics: "(terms)foo,(histogram)qux", Flush: 1}
	processor := NewMetrics(&config)

	channel := make(chan Record, 1)
	processor.SetChannel(channel)

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

	processor.Do(Record{"foo": "qux"})
	processor.Teardown()

	record := <-channel
	c.Check(record["foo"].(map[string]int)["bar"], Equals, 50000)
	c.Check(record["qux"].(map[string]interface{})["count"], Equals, int64(50000))
}
