package harvesterd

import (
	. "harvesterd/logger"

	"runtime"
	"time"
)

type Harvesterd struct {
	writerGroup *WriterGroup
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
}

func (self *Harvesterd) configureLogger() {
	Info("Starting ...")
}

func (self *Harvesterd) configureMaxProcs() {
	Info("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (self *Harvesterd) bootWriter() {
	self.writerGroup = GetContainer().GetWriterGroup()
}

func (self *Harvesterd) Run() {
	self.writerGroup.Setup()
	self.writerGroup.Boot()
	self.wait()
	self.writerGroup.Teardown()
}

func (self *Harvesterd) wait() {
	for self.writerGroup.IsAlive() {
		time.Sleep(1 * time.Second)
		self.PrintCounters(1)
	}

	Info("nothing more for read, terminating daemon ...")
}

func (self *Harvesterd) PrintCounters(elapsedSeconds int) {
	created, failed, _, threads := self.writerGroup.GetCounters()
	self.writerGroup.ResetCounters()

	logFormat := "Created %d document(s), failed %d times(s), %g doc/sec in %d thread(s)"

	rate := float64(created+failed) / float64(elapsedSeconds)
	Info(logFormat, created, failed, rate, threads)
}
