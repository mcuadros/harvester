package processor

import (
	. "harvesterd/intf"
	. "harvesterd/logger"
	. "harvesterd/processor/metric"
	"regexp"
	"strings"
)

type MetricsConfig struct {
	Flush   int
	Metrics string
}

type Metrics struct {
	metrics []Metric
	flush   int
}

type Metric interface {
	Process(record Record)
	GetValue() interface{}
	GetField() string
	Reset()
}

var configRegExp = regexp.MustCompile("^\\((\\w+)\\)(\\w+)$")

func NewMetrics(config *MetricsConfig) *Metrics {
	processor := new(Metrics)
	processor.SetConfig(config)

	return processor
}

func (self *Metrics) SetConfig(config *MetricsConfig) {
	self.parseMetricsConfig(config.Metrics)
}

func (self *Metrics) parseMetricsConfig(metricsConfig string) {
	for _, config := range strings.Split(metricsConfig, ",") {
		var metric Metric
		class, field := self.parseMetric(config)

		switch class {
		case "terms":
			metric = NewTerms(field)
		case "histogram":
			metric = NewHistogram(field)
		}

		self.metrics = append(self.metrics, metric)
	}
}

func (self *Metrics) parseMetric(metric string) (class string, field string) {
	config := configRegExp.FindStringSubmatch(metric)
	if len(config) != 3 {
		Critical("Malformed metric config \"%s\"", metric)
	}

	return config[1], config[2]
}

func (self *Metrics) Do(record Record) {
	for _, metric := range self.metrics {
		metric.Process(record)
		record[metric.GetField()] = metric.GetValue()
	}
}
