package format

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type FileConfig struct {
	File    string
	Pattern string
}

type File struct {
	files   []*bufio.Scanner
	current int
	eof     bool
}

func NewFile(config FileConfig) *File {
	input := new(File)
	input.SetConfig(config)

	return input
}

func (self *File) SetConfig(config FileConfig) {
	if config.File != "" {
		self.files = append(self.files, self.createBufioReader(config.File))
	} else if config.Pattern != "" {
		files, err := filepath.Glob(config.Pattern)
		if err != nil {
			panic(fmt.Sprintf("open %s: %v", config.Pattern, err))
		}

		for _, file := range files {
			self.files = append(self.files, self.createBufioReader(file))
		}
	}
}

func (self *File) createBufioReader(filename string) *bufio.Scanner {
	file, err := os.Open(filename)
	if err != nil {
		panic(fmt.Sprintln("open %s: %v", filename, err))
	}

	return bufio.NewScanner(file)
}

func (self *File) GetLine() string {
	if self.scan() {
		return self.files[self.current].Text()
	}

	return ""
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
