package input

import (
	"bufio"
	. "harvesterd/intf"
	. "harvesterd/logger"
	"os"
	"path/filepath"
)

type FileConfig struct {
	Format  string
	Pattern string
}

type File struct {
	files   []*bufio.Scanner
	format  Format
	current int
	empty   bool
	eof     bool
}

func NewFile(config *FileConfig, format Format) *File {
	input := new(File)
	input.SetConfig(config)
	input.SetFormat(format)

	return input
}

func (self *File) SetFormat(format Format) {
	self.format = format
}

func (self *File) SetConfig(config *FileConfig) {
	files, err := filepath.Glob(config.Pattern)
	if err != nil {
		Critical("open %s: %v", config.Pattern, err)
	}

	for _, file := range files {
		self.files = append(self.files, self.createBufioReader(file))
	}

	if len(self.files) == 0 {
		self.empty = true
		self.eof = true
	}
}

func (self *File) createBufioReader(filename string) *bufio.Scanner {
	file, err := os.Open(filename)
	if err != nil {
		Critical("open %s: %v", filename, err)
	}

	return bufio.NewScanner(file)
}

func (self *File) GetLine() string {
	if !self.empty && self.scan() {
		return self.files[self.current].Text()
	}

	return ""
}

func (self *File) GetRecord() Record {
	line := self.GetLine()
	if line != "" {
		return self.format.Parse(line)
	}

	return nil
}

func (self *File) scan() bool {
	if !self.files[self.current].Scan() {
		self.current++

		if self.current >= len(self.files) {
			self.eof = true
			return false
		}

		return self.scan()
	}

	return true
}

func (self *File) IsEOF() bool {
	return self.eof
}

func (self *File) Teardown() {
}
