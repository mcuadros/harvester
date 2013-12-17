package input

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

import "github.com/ActiveState/tail"

type TailConfig struct {
	File      string // File to be readed
	MustExist bool   // Fail early if the file does not exist
	Poll      bool   // Poll for file changes instead of using inotify
	LimitRate int64  // Maximum read rate (lines per second)
}

type Tail struct {
	tail    *tail.Tail
	file    string
	posFile string
	counter int
	eof     bool
}

func NewTail(config TailConfig) *Tail {
	input := new(Tail)
	input.SetConfig(config)

	return input
}

func (self *Tail) SetConfig(config TailConfig) {
	self.file = config.File
	self.posFile = fmt.Sprintf("%s/.%s.pos", path.Dir(self.file), path.Base(self.file))

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
		panic(fmt.Sprintf("read %s: %v", self.posFile, err))
	}

	position, err := strconv.ParseInt(string(positionRaw), 10, 0)
	if err != nil {
		panic(fmt.Sprintf("malformed content %s: %v", self.posFile, err))
	}

	return position
}

func (self *Tail) GetLine() string {
	line, ok := (<-self.tail.Lines)
	if ok {
		self.counter++
		go self.keepPosition()
		return line.Text
	} else {
		self.eof = true
		return ""
	}
}

func (self *Tail) keepPosition() {
	if self.counter > 1 {
		position, _ := self.tail.Tell()
		ioutil.WriteFile(self.posFile, []byte(strconv.FormatInt(position, 10)), 0755)
	}
}

func (self *Tail) IsEOF() bool {
	return self.eof
}

func (self *Tail) Stop() {
	self.tail.Stop()
}
