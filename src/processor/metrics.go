package processor

import (
	"harvesterd/intf"
	. "harvesterd/logger"
	. "harvesterd/processor/metric"
	"regexp"
	"strings"
	"sync"
	"time"
)

type MetricsConfig struct {
	Flush   int
	Metrics string
}

type Metrics struct {
	metrics []Metric
	flush   time.Duration
	mutex   sync.Mutex
	channel chan intf.Record
	isAlive bool
}

type Metric interface {
	Process(record intf.Record)
	GetValue() interface{}
	GetField() string
	Reset()
}

var configRegExp = regexp.MustCompile("^\\((\\w+)\\)(\\w+)$")

func NewMetrics(config *MetricsConfig) *Metrics {
	processor := new(Metrics)
	processor.SetConfig(config)
	processor.Setup()

	return processor
}

func (self *Metrics) SetConfig(config *MetricsConfig) {
	self.flush = time.Duration(config.Flush)
	self.parseMetricsConfig(config.Metrics)
}

func (self *Metrics) SetChannel(channel chan intf.Record) {
	self.channel = channel
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
		default:
			Critical("Unknown metric \"%s\", valid: [terms histogram]", class)
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

func (self *Metrics) Do(record intf.Record) bool {
	self.mutex.Lock()

	for _, metric := range self.metrics {
		metric.Process(record)
	}

	self.mutex.Unlock()

	return false
}

func (self *Metrics) Setup() {
	self.isAlive = true
	go self.deliveryRecord()
}

func (self *Metrics) deliveryRecord() {
	Debug("Started metrics routine")
	for {
		time.Sleep(self.flush * time.Second)
		if self.isAlive {
			self.flushMetrics()
		}
	}
}

func (self *Metrics) flushMetrics() {
	var record = make(map[string]interface{})
	for _, metric := range self.metrics {
		record[metric.GetField()] = metric.GetValue()
	}

	self.channel <- record
}

func (self *Metrics) Teardown() {
	self.isAlive = false
	self.flushMetrics()
}
