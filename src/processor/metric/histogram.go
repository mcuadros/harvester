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

func (m *Histogram) Process(record Record) {
	var value float64

	switch record[m.field].(type) {
	case int:
		value = float64(record[m.field].(int))
	case float64:
		value = record[m.field].(float64)
	default:
		return
	}

	m.histogram.Update(int64(value * float64(m.precision)))

}

func (m *Histogram) GetField() string {
	return m.field
}

func (m *Histogram) GetValue() interface{} {
	result := make(map[string]interface{})
	result["count"] = m.histogram.Count()
	result["min"] = float64(m.histogram.Min()) / float64(m.precision)
	result["max"] = float64(m.histogram.Max()) / float64(m.precision)
	result["mean"] = m.histogram.Mean() / float64(m.precision)
	result["sum"] = float64(m.histogram.Count()) * result["mean"].(float64)
	result["stddev"] = m.histogram.StdDev() / float64(m.precision)

	return result
}

func (m *Histogram) Reset() {
	m.histogram.Clear()
}
