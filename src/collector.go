package collector

import (
	. "collector/logger"
	"runtime"
	"time"
)

type Collector struct {
	writer  *Writer
	reader  *Reader
	channel chan map[string]string
}

func NewCollector() *Collector {
	collector := new(Collector)

	return collector
}

func (self *Collector) Configure(filename string) {
	GetConfig().LoadFile(filename)
}

func (self *Collector) Boot() {
	self.configureLogger()
	self.configureMaxProcs()
	self.bootWriter()
	self.bootReader()
}

func (self *Collector) configureLogger() {
	Info("Starting ...")
}

func (self *Collector) configureMaxProcs() {
	Info("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (self *Collector) bootWriter() {
	self.writer = GetContainer().GetWriter()
}

func (self *Collector) bootReader() {
	self.reader = GetContainer().GetReader()
}

func (self *Collector) Run() {
	self.channel = self.writer.GoWriteFromChannel()
	self.reader.GoReadIntoChannel(self.channel)
	self.wait()
}

func (self *Collector) wait() {
	print(self.writer.IsAlive())
	for {
		time.Sleep(1 * time.Second)
		Warning("foo")
		print(self.writer.IsAlive())
		//GetLogger().PrintWriterStats(3, self.writer)
	}
}
