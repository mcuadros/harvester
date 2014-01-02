package processor

import (
	. "harvesterd/intf"
	. "harvesterd/logger"
	"regexp"
	"strings"
)

type AggregationConfig struct {
	Flush   int
	Metrics string
}

type Aggregation struct {
	metrics []Metric
	flush   int
}

var configRegExp = regexp.MustCompile("^(\\w+)\\((\\w+)\\)$")

func NewAggregation(config *AggregationConfig) *Aggregation {
	processor := new(Aggregation)
	processor.SetConfig(config)

	return processor
}

func (self *Aggregation) SetConfig(config *AggregationConfig) {
	self.parseMetricsConfig(config.Metrics)
}

func (self *Aggregation) parseMetricsConfig(metricsConfig string) {
	for _, config := range strings.Split(metricsConfig, ",") {
		var metric Metric
		class, field := self.parseMetric(config)

		switch class {
		case "count":
			metric = &CountMetric{field: field, count: make(map[string]int32, 0)}
		}

		self.metrics = append(self.metrics, metric)
	}
}

func (self *Aggregation) parseMetric(metric string) (class string, field string) {
	config := configRegExp.FindStringSubmatch(metric)
	if len(config) != 3 {
		Critical("Malformed metric config \"%s\"", metric)
	}

	return config[1], config[2]
}

func (self *Aggregation) Do(record Record) {
	for _, metric := range self.metrics {
		metric.Process(record)
		record[metric.GetField()] = metric.GetValue()
	}
}

type Metric interface {
	Process(record Record)
	GetValue() interface{}
	GetField() string
	Reset()
}

type CountMetric struct {
	field string
	count map[string]int32
}

func (self *CountMetric) Process(record Record) {
	key := record[self.field].(string)
	if _, ok := self.count[key]; ok {
		self.count[key]++
	} else {
		self.count[key] = 1
	}
}

func (self *CountMetric) GetField() string {
	return self.field
}

func (self *CountMetric) GetValue() interface{} {
	return self.count
}

func (self *CountMetric) Reset() {
	self.count = make(map[string]int32, 0)
}
