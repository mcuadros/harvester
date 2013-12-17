package collector

import (
	"fmt"
)

type Reader struct {
	counter int32
	format  Format
	input   Input
}

func NewReader(format Format, input Input) *Reader {
	reader := new(Reader)
	reader.SetFormat(format)
	reader.SetInput(input)

	return reader
}

func (self *Reader) SetFormat(format Format) {
	self.format = format
}

func (self *Reader) SetInput(input Input) {
	self.input = input
}

func (self *Reader) ReadIntoChannel(channel chan map[string]string) {
	defer close(channel)

	for !self.input.IsEOF() {
		line := self.input.GetLine()
		row := self.format.Parse(line)

		self.emitRecord(channel, row)
	}
}

func (self *Reader) emitRecord(channel chan map[string]string, row map[string]string) {
	if len(row) > 0 {
		channel <- row
		self.counter++
		if self.counter%1000 == 0 {
			fmt.Println(fmt.Sprintf("%d", self.counter))
		}
	}
}
