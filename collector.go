package collector

import (
	"fmt"
	//"net/http"
	"runtime"
	"sync"
	"time"
)

import "code.google.com/p/gcfg"

type Config struct {
	Basic struct {
		Threads int
	}
	Logger        LoggerConfig
	CSV           FormatCSVConfig
	Reader        ReaderConfig
	ElasticSearch WriterElasticSearchConfig
}

type Collector struct {
	config Config
	lines  chan map[string]string
	wait   sync.WaitGroup
	writer Writer
	reader Reader
}

func (self *Collector) Configure() {
	err := gcfg.ReadFileInto(&self.config, "config.ini")
	if err != nil {
		panic(fmt.Sprintf("open config: %v", err))
	}
}

func (self *Collector) ReadFile() {
	format := NewFormatCSV(self.config.CSV)

	reader := NewReader(self.config.Reader)
	reader.SetFormat(format)

	go reader.ReadIntoChannel(self.lines)
	go self.PrintStatus()

	self.wait.Wait()

}

func (self *Collector) PrintStatus() {
	for {
		time.Sleep(1 * time.Second)
		GetLogger().PrintWriterStats(3, self.writer)
	}
}

func (self *Collector) Boot() {
	NewLogger(self.config.Logger)
	GetLogger().Info("Starting ...")
	GetLogger().Debug("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())

	self.writer = NewWriterElasticSearch(self.config.ElasticSearch)

	threads := self.config.Basic.Threads
	self.lines = make(chan map[string]string, threads)

	for i := 0; i < threads; i++ {
		self.wait.Add(1)
		go self.writer.WriteFromChannel(self.lines, self.wait)
	}

	GetLogger().Debug("Started %d thread(s)", threads)

}
