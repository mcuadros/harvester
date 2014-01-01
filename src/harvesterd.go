package harvesterd

import (
	. "harvesterd/intf"
	. "harvesterd/logger"
	"runtime"
	"time"
)

type Harvesterd struct {
	writer  *Writer
	reader  *Reader
	channel chan Record
}

func NewHarvesterd() *Harvesterd {
	harvesterd := new(Harvesterd)

	return harvesterd
}

func (self *Harvesterd) Configure(filename string) {
	GetConfig().LoadFile(filename)
}

func (self *Harvesterd) Boot() {
	self.configureLogger()
	self.configureMaxProcs()
	self.bootWriter()
	self.bootReader()
}

func (self *Harvesterd) configureLogger() {
	Info("Starting ...")
}

func (self *Harvesterd) configureMaxProcs() {
	Info("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (self *Harvesterd) bootWriter() {
	self.writer = GetContainer().GetWriter()
}

func (self *Harvesterd) bootReader() {
	self.reader = GetContainer().GetReader()
}

func (self *Harvesterd) Run() {
	self.channel = self.writer.GoWriteFromChannel()
	self.reader.GoReadIntoChannel(self.channel)
	self.wait()
	self.reader.Finish()
}

func (self *Harvesterd) wait() {
	for self.writer.IsAlive() {
		time.Sleep(1 * time.Second)
		self.writer.PrintCounters(1)
	}

	Info("nothing more for read, terminating daemon ...")
}
