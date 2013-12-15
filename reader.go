package collector

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type ReaderConfig struct {
	File    string
	Fields  string
	Pattern string
}

type Reader struct {
	file    string
	pattern string
	header  []string
	counter int32
	format  Format
}

func NewReader(config ReaderConfig) *Reader {
	reader := new(Reader)
	reader.SetConfig(config)

	return reader
}

func (self *Reader) SetConfig(config ReaderConfig) {
	self.file = config.File
	self.pattern = config.Pattern
}

func (self *Reader) SetFormat(format Format) {
	self.format = format
}

func (self *Reader) ReadIntoChannel(channel chan map[string]string) {
	if self.file != "" {
		self.readFileInChannel(self.file, channel)
	} else {
		files, err := filepath.Glob(self.pattern)
		if err != nil {
			panic(fmt.Sprintf("open %s: %v", self.pattern, err))
		}

		fmt.Println(files)
		for _, file := range files {
			self.readFileInChannel(file, channel)
		}
	}
}

func (self *Reader) readFileInChannel(filename string, channel chan map[string]string) {
	GetLogger().Info("Processing '%s'", filename)

	file, err := os.Open(filename)
	if err != nil {
		GetLogger().Error("open %s: %v", self.file, err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	var line string
	for scanner.Scan() {
		line = scanner.Text()

		row := self.format.Parse(line)
		self.emitRecord(filename, channel, row)
	}
}

func (self *Reader) emitRecord(file string, channel chan map[string]string, row map[string]string) {
	if len(row) > 0 {
		channel <- row
		self.counter++
		if self.counter%1000 == 0 {
			fmt.Println(fmt.Sprintf("%d", self.counter))
		}
	}
}
