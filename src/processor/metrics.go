package processor

import (
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
	. "github.com/mcuadros/harvester/src/processor/metric"
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

func (p *Metrics) SetConfig(config *MetricsConfig) {
	p.flush = time.Duration(config.Flush)
	p.parseMetricsConfig(config.Metrics)
}

func (p *Metrics) SetChannel(channel chan intf.Record) {
	p.channel = channel
}

func (p *Metrics) parseMetricsConfig(metricsConfig string) {
	for _, config := range strings.Split(metricsConfig, ",") {
		var metric Metric
		class, field := p.parseMetric(config)

		switch class {
		case "terms":
			metric = NewTerms(field)
		case "histogram":
			metric = NewHistogram(field)
		default:
			Critical("Unknown metric \"%s\", valid: [terms histogram]", class)
		}

		p.metrics = append(p.metrics, metric)
	}
}

func (p *Metrics) parseMetric(metric string) (class string, field string) {
	config := configRegExp.FindStringSubmatch(metric)
	if len(config) != 3 {
		Critical("Malformed metric config \"%s\"", metric)
	}

	return config[1], config[2]
}

func (p *Metrics) Do(record intf.Record) bool {
	p.mutex.Lock()

	for _, metric := range p.metrics {
		metric.Process(record)
	}

	p.mutex.Unlock()

	return false
}

func (p *Metrics) Setup() {
	p.isAlive = true
	go p.deliveryRecord()
}

func (p *Metrics) deliveryRecord() {
	Debug("Started metrics routine")
	for {
		time.Sleep(p.flush * time.Second)
		if p.isAlive {
			p.flushMetrics()
		}
	}
}

func (p *Metrics) flushMetrics() {
	var record = make(map[string]interface{})
	for _, metric := range p.metrics {
		record[metric.GetField()] = metric.GetValue()
	}

	p.channel <- record
}

func (p *Metrics) Teardown() {
	p.isAlive = false
	p.flushMetrics()
}
