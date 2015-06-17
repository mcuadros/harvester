package harvesterd

import (
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	. "github.com/mcuadros/harvesterd/src/logger"
)

type Harvesterd struct {
	writerGroup  *WriterGroup
	signsChannel chan os.Signal
	isAlive      bool
}

func NewHarvesterd() *Harvesterd {
	harvesterd := new(Harvesterd)

	return harvesterd
}

func (h *Harvesterd) Configure(filename string) {
	GetConfig().LoadFile(filename)
}

func (h *Harvesterd) Boot() {
	h.isAlive = true
	h.configureLogger()
	h.configureMaxProcs()
	h.bootSignalWaiter()
	h.bootWriter()
}

func (h *Harvesterd) bootSignalWaiter() {
	h.signsChannel = make(chan os.Signal, 1)

	signal.Notify(h.signsChannel, syscall.SIGINT, syscall.SIGTERM)
	go h.signalWaiter()
}

func (h *Harvesterd) signalWaiter() {
	signal := <-h.signsChannel
	Warning("Received signal: %s", signal)
	h.isAlive = false
}

func (h *Harvesterd) configureLogger() {
	Info("Starting ...")
}

func (h *Harvesterd) configureMaxProcs() {
	Info("Number of max. process %d", runtime.NumCPU())
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func (h *Harvesterd) bootWriter() {
	h.writerGroup = GetContainer().GetWriterGroup()
}

func (h *Harvesterd) Run() {
	h.writerGroup.Setup()
	h.writerGroup.Boot()
	h.wait()
	h.writerGroup.Teardown()
}

func (h *Harvesterd) wait() {
	for h.writerGroup.IsAlive() && h.isAlive {
		time.Sleep(1 * time.Second)
		h.PrintCounters(1)
	}

	Info("nothing more for read, terminating daemon ...")
}

func (h *Harvesterd) PrintCounters(elapsedSeconds int) {
	created, failed, _, threads := h.writerGroup.GetCounters()
	h.writerGroup.ResetCounters()

	logFormat := "processed %d document(s), failed %d times(s), %g doc/sec in %d thread(s)"

	rate := float64(created+failed) / float64(elapsedSeconds)
	Info(logFormat, created, failed, rate, threads)
}
