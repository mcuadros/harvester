package input

import (
	"fmt"
	"harvesterd/intf"
	. "harvesterd/logger"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"
)

import "github.com/ActiveState/tail"

type TailConfig struct {
	Format    string
	File      string // File to be readed
	MustExist bool   // Fail early if the file does not exist
	Poll      bool   // Poll for file changes instead of using inotify
	LimitRate int64  // Maximum read rate (lines per second)ssh
}

type Tail struct {
	tail          *tail.Tail
	format        intf.Format
	file          string
	posFile       string
	eof           bool
	needSavePos   bool
	stopChannel   chan struct{}
	savePosTicker *time.Ticker
}

func NewTail(config *TailConfig, format intf.Format) *Tail {
	input := new(Tail)
	input.SetConfig(config)
	input.SetFormat(format)
	input.Boot()

	return input
}

func (self *Tail) Boot() {
	self.goKeepPosition()
}

func (self *Tail) SetFormat(format intf.Format) {
	self.format = format
}

func (self *Tail) SetConfig(config *TailConfig) {
	self.file = config.File
	Info(self.file)
	self.posFile = fmt.Sprintf("%s/.%s.pos", path.Dir(self.file), path.Base(self.file))

	self.createTailReader(self.translateConfig(config))
}

func (self *Tail) createTailReader(config tail.Config) {
	tail, err := tail.TailFile(self.file, config)
	if err != nil {
		Critical("tail %s: %v", self.file, err)
	}

	self.tail = tail
}

func (self *Tail) translateConfig(original *TailConfig) tail.Config {
	config := tail.Config{Follow: true, ReOpen: true}

	if original.MustExist {
		config.MustExist = true
	}

	if original.Poll {
		config.Poll = true
	}

	if original.LimitRate > 0 {
		config.LimitRate = original.LimitRate
	}

	position := self.readPosition()
	if position > 0 {
		config.Location = &tail.SeekInfo{Offset: position, Whence: 0}
	}

	return config
}

func (self *Tail) readPosition() int64 {
	_, err := os.Stat(self.posFile)
	if os.IsNotExist(err) {
		return 0
	}

	positionRaw, err := ioutil.ReadFile(self.posFile)
	if err != nil {
		Critical("read %s: %v", self.posFile, err)
	}

	position, err := strconv.ParseInt(string(positionRaw), 10, 0)
	if err != nil {
		Critical("malformed content %s: %v", self.posFile, err)
	}

	return position
}

func (self *Tail) GetLine() string {
	line, ok := (<-self.tail.Lines)
	if ok {
		self.needSavePos = true
		return line.Text
	} else {
		self.eof = true
		return ""
	}
}

func (self *Tail) GetRecord() intf.Record {
	line := self.GetLine()
	if line != "" {
		return self.format.Parse(line)
	}

	return nil
}

func (self *Tail) goKeepPosition() {
	self.savePosTicker = time.NewTicker(1 * time.Second)
	self.stopChannel = make(chan struct{})

	go func() {
		for {
			select {
			case <-self.savePosTicker.C:
				if self.needSavePos {
					self.keepPosition()
					self.needSavePos = false
				}
			case <-self.stopChannel:
				self.savePosTicker.Stop()
				return
			}
		}
	}()
}

func (self *Tail) keepPosition() {
	position, _ := self.tail.Tell()
	ioutil.WriteFile(self.posFile, []byte(strconv.FormatInt(position, 10)), 0755)
}

func (self *Tail) IsEOF() bool {
	return self.eof
}

func (self *Tail) Stop() {
	self.keepPosition()
	self.tail.Stop()
}

func (self *Tail) Teardown() {
	close(self.stopChannel)
}
