package metric

import (
	"runtime"
	"sync"
	"testing"

	. "github.com/mcuadros/harvesterd/src/intf"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MetricsSuite struct{}

var _ = Suite(&MetricsSuite{})

func (s *MetricsSuite) TestProcess(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	metric := NewTerms("foo")

	var wait sync.WaitGroup
	var add = func() {
		for i := 0; i < 10000; i++ {
			metric.Process(Record{"foo": "bar"})
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
	metric.Process(record)

	result := metric.GetValue().(map[string]int)
	c.Assert(result["bar"], Equals, 50000)
	c.Assert(result["qux"], Equals, 1)
}

func (s *MetricsSuite) TestProcessNonString(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	metric := NewTerms("foo")

	var wait sync.WaitGroup
	var add = func() {
		for i := 0; i < 10000; i++ {
			metric.Process(Record{"foo": 1})
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
	metric.Process(record)

	result := metric.GetValue().(map[string]int)
	c.Assert(result["bar"], Equals, 0)
	c.Assert(result["qux"], Equals, 1)
}
