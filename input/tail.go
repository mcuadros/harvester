package format

import (
	"fmt"
)

import "github.com/ActiveState/tail"

type TailConfig struct {
	File      string // File to be readed
	MustExist bool   // Fail early if the file does not exist
	Poll      bool   // Poll for file changes instead of using inotify
	LimitRate int64  // Maximum read rate (lines per second)
}

type Tail struct {
	tail *tail.Tail
	file string
	eof  bool
}

func NewTail(config TailConfig) *Tail {
	input := new(Tail)
	input.SetConfig(config)

	return input
}

func (self *Tail) SetConfig(config TailConfig) {
	self.file = config.File

	self.createTailReader(self.translateConfig(config))
}

func (self *Tail) createTailReader(config tail.Config) {
	tail, err := tail.TailFile(self.file, config)
	if err != nil {
		panic(fmt.Sprintf("tail %s: %v", self.file, err))
	}

	self.tail = tail
}

func (self *Tail) translateConfig(original TailConfig) tail.Config {
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

	return config
}

func (self *Tail) GetLine() string {
	line, ok := (<-self.tail.Lines)
	if ok {
		return line.Text
	} else {
		self.eof = true
		return ""
	}
}

func (self *Tail) IsEOF() bool {
	return self.eof
}

func (self *Tail) Stop() {
	self.tail.Stop()
}
