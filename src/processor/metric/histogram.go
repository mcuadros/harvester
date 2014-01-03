package metric

import (
	. "harvesterd/intf"
)

import "github.com/rcrowley/go-metrics"

type Histogram struct {
	field     string
	histogram metrics.Histogram
	precision int
}

func NewHistogram(field string) *Histogram {
	counter := &Histogram{
		field:     field,
		histogram: metrics.NewHistogram(metrics.NewExpDecaySample(1028, 0.015)),
		precision: 1000000,
	}

	return counter
}

func (self *Histogram) Process(record Record) {
	var value float64

	switch record[self.field].(type) {
	case int:
		value = float64(record[self.field].(int))
	case float64:
		value = record[self.field].(float64)
	default:
		return
	}

	self.histogram.Update(int64(value * float64(self.precision)))

}

func (self *Histogram) GetField() string {
	return self.field
}

func (self *Histogram) GetValue() interface{} {
	result := make(map[string]interface{})
	result["count"] = self.histogram.Count()
	result["min"] = float64(self.histogram.Min()) / float64(self.precision)
	result["max"] = float64(self.histogram.Max()) / float64(self.precision)
	result["mean"] = self.histogram.Mean() / float64(self.precision)
	result["sum"] = float64(self.histogram.Count()) * result["mean"].(float64)
	result["stddev"] = self.histogram.StdDev() / float64(self.precision)

	return result
}

func (self *Histogram) Reset() {
	self.histogram.Clear()
}
