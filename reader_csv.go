package collector

import (
	"fmt"
	"os"
	"strings"
)

import "github.com/gwenn/yacr"

type ReaderCSVConfig struct {
	File   string
	Fields string
}

type ReaderCSV struct {
	file    string
	header  []string
	counter int32
}

func NewReaderCSV(config ReaderCSVConfig) *ReaderCSV {
	reader := new(ReaderCSV)
	reader.SetConfig(config)

	fmt.Print(reader.header)
	return reader
}

func (self *ReaderCSV) SetConfig(config ReaderCSVConfig) {
	for _, field := range strings.Split(config.Fields, ",") {
		self.header = append(self.header, field)
	}

	self.file = config.File
}

func (self *ReaderCSV) ReadIntoChannel(channel chan map[string]string) {
	file, err := os.Open(self.file)
	if err != nil {
		panic(fmt.Sprintf("open %s: %v", self.file, err))
	}

	defer file.Close()

	reader := yacr.DefaultReader(file)

	row := make(map[string]string)
	for reader.Scan() {
		if reader.EmptyLine() { // skip empty line (or line comment)
			continue
		}

		row[self.header[len(row)]] = reader.Text()

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
