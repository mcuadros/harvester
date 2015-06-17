package input

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/mcuadros/harvesterd/src/intf"
	. "github.com/mcuadros/harvesterd/src/logger"

	"github.com/ActiveState/tail"
	"github.com/ActiveState/tail/ratelimiter"
)

type TailConfig struct {
	Format    string `description:"A valid format name"`
	File      string `description:"File to be read"`
	MustExist bool   `description:"Fail early if the file does not exist"`
	Poll      bool   `description:"Poll for file changes instead of using inotify"`
	LimitRate int64  `description:"Maximum read rate (lines per second)"`
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

func (i *Tail) Boot() {
	i.goKeepPosition()
}

func (i *Tail) SetFormat(format intf.Format) {
	i.format = format
}

func (i *Tail) SetConfig(config *TailConfig) {
	i.file = config.File
	Info(i.file)
	i.posFile = fmt.Sprintf("%s/.%s.pos", path.Dir(i.file), path.Base(i.file))

	i.createTailReader(i.translateConfig(config))
}

func (i *Tail) createTailReader(config tail.Config) {
	tail, err := tail.TailFile(i.file, config)
	if err != nil {
		Critical("tail %s: %v", i.file, err)
	}

	i.tail = tail
}

func (i *Tail) translateConfig(original *TailConfig) tail.Config {
	config := tail.Config{Follow: true, ReOpen: true}

	if original.MustExist {
		config.MustExist = true
	}

	if original.Poll {
		config.Poll = true
	}

	if original.LimitRate > 0 {
		config.RateLimiter = ratelimiter.NewLeakyBucket(
			uint16(original.LimitRate),
			time.Second,
		)
	}

	position := i.readPosition()
	if position > 0 {
		config.Location = &tail.SeekInfo{Offset: position, Whence: 0}
	}

	return config
}

func (i *Tail) readPosition() int64 {
	_, err := os.Stat(i.posFile)
	if os.IsNotExist(err) {
		return 0
	}

	positionRaw, err := ioutil.ReadFile(i.posFile)
	if err != nil {
		Critical("read %s: %v", i.posFile, err)
	}

	position, err := strconv.ParseInt(string(positionRaw), 10, 0)
	if err != nil {
		Critical("malformed content %s: %v", i.posFile, err)
	}

	return position
}

func (i *Tail) GetLine() string {
	line, ok := (<-i.tail.Lines)
	if ok {
		i.needSavePos = true
		return line.Text
	} else {
		i.eof = true
		return ""
	}
}

func (i *Tail) GetRecord() intf.Record {
	line := i.GetLine()
	if line != "" {
		return i.format.Parse(line)
	}

	return nil
}

func (i *Tail) goKeepPosition() {
	i.savePosTicker = time.NewTicker(1 * time.Second)
	i.stopChannel = make(chan struct{})

	go func() {
		for {
			select {
			case <-i.savePosTicker.C:
				if i.needSavePos {
					i.keepPosition()
					i.needSavePos = false
				}
			case <-i.stopChannel:
				i.savePosTicker.Stop()
				return
			}
		}
	}()
}

func (i *Tail) keepPosition() {
	position, _ := i.tail.Tell()
	ioutil.WriteFile(i.posFile, []byte(strconv.FormatInt(position, 10)), 0755)
}

func (i *Tail) IsEOF() bool {
	return i.eof
}

func (i *Tail) Stop() {
	i.keepPosition()
	i.tail.Stop()
}

func (i *Tail) Teardown() {
	close(i.stopChannel)
}
