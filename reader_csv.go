package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

import "github.com/gwenn/yacr"

type ReaderCSVConfig struct {
	File    string
	Fields  string
	Pattern string
}

type ReaderCSV struct {
	file    string
	pattern string
	header  []string
	counter int32
}

func NewReaderCSV(config ReaderCSVConfig) *ReaderCSV {
	reader := new(ReaderCSV)
	reader.SetConfig(config)

	return reader
}

func (self *ReaderCSV) SetConfig(config ReaderCSVConfig) {
	for _, field := range strings.Split(config.Fields, ",") {
		self.header = append(self.header, field)
	}

	self.file = config.File
	self.pattern = config.Pattern
}

func (self *ReaderCSV) ReadIntoChannel(channel chan map[string]string) {
	if self.file != "" {
		self.readFileInChannel(self.file, channel)
	} else {
		files, err := filepath.Glob(self.pattern)
		if err != nil {
			panic(fmt.Sprintf("open %s: %v", self.pattern, err))
		}

		for _, file := range files {
			self.readFileInChannel(file, channel)
		}

		fmt.Println(files)
	}
}

func (self *ReaderCSV) readFileInChannel(filename string, channel chan map[string]string) {
	GetLogger().Info("Processing '%s'", filename)

	file, err := os.Open(filename)
	if err != nil {
		GetLogger().Error("open %s: %v", self.file, err)
	}

	defer file.Close()

	reader := yacr.DefaultReader(file)
	headers := len(self.header)
	row := make(map[string]string)
	for reader.Scan() {
		if reader.EmptyLine() { // skip empty line (or line comment)
			continue
		}

		if len(row) < headers {
			row[self.header[len(row)]] = reader.Text()
		}

		if reader.EndOfRecord() {
			self.emitRecord(channel, row)
			row = make(map[string]string)
		}
	}

	self.emitRecord(channel, row)
}

func (self *ReaderCSV) emitRecord(channel chan map[string]string, row map[string]string) {
	if len(row) > 0 {
		channel <- row
		self.counter++
	}
}
