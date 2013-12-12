package collector

import (
	"fmt"
	//"net/http"
	"runtime"
	"sync"
	//"time"
)

import "code.google.com/p/gcfg"

type Config struct {
	CSV           ReaderCSVConfig
	ElasticSearch WriterElasticSearchConfig
}

type Collector struct {
	config Config
	lines  chan map[string]string
	wait   sync.WaitGroup
}

func (self *Collector) Configure() {
	err := gcfg.ReadFileInto(&self.config, "config.ini")
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}
}

func (self *Collector) ReadFile() {
	reader := NewReaderCSV(self.config.CSV)
	reader.ReadIntoChannel(self.lines)
}

func (self *Collector) Boot() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := NewWriterElasticSearch(self.config.ElasticSearch)

	self.lines = make(chan map[string]string, 24)

	for i := 0; i < 24; i++ {
		self.wait.Add(1)
		go writer.WriteFromChannel(self.lines, self.wait)
	}

}
