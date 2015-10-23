package harvester

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	. "github.com/mcuadros/harvester/src/logger"
)

type Harvester struct {
	writerGroup  *WriterGroup
	signsChannel chan os.Signal
	isAlive      bool
}

func NewHarvester() *Harvester {
	harvester := new(Harvester)

	return harvester
}

func (h *Harvester) Configure(filename string) {
	GetConfig().LoadFile(filename)
}

func (h *Harvester) Boot() {
	h.isAlive = true
	h.configureLogger()
	h.configureMaxProcs()
	h.bootSignalWaiter()
	h.bootWriter()
}

func (h *Harvester) bootSignalWaiter() {
	h.signsChannel = make(chan os.Signal, 1)

	signal.Notify(h.signsChannel, syscall.SIGINT, syscall.SIGTERM)
	go h.signalWaiter()
}

func (h *Harvester) signalWaiter() {
	signal := <-h.signsChannel
	Warning("Received signal: %s", signal)
	h.isAlive = false
}

func (h *Harvester) configureLogger() {
	Info("Starting ...")
}

func (h *Harvester) configureMaxProcs() {
	Info("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (h *Harvester) bootWriter() {
	h.writerGroup = GetContainer().GetWriterGroup()
}

func (h *Harvester) Run() {
	h.writerGroup.Setup()
	h.writerGroup.Boot()
	h.wait()
	h.writerGroup.Teardown()
}

func (h *Harvester) wait() {
	for h.writerGroup.IsAlive() && h.isAlive {
		time.Sleep(1 * time.Second)
		h.PrintCounters(1)
	}

	Info("nothing more for read, terminating daemon ...")
}

func (h *Harvester) PrintCounters(elapsedSeconds int) {
	created, failed, _, threads := h.writerGroup.GetCounters()
	h.writerGroup.ResetCounters()

	logFormat := "processed %d document(s), failed %d times(s), %g doc/sec in %d thread(s)"

	rate := float64(created+failed) / float64(elapsedSeconds)
	Info(logFormat, created, failed, rate, threads)
}
