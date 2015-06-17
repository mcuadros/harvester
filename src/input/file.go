package input

import (
	"bufio"
	"harvesterd/intf"
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
	format  intf.Format
	current int
	empty   bool
	eof     bool
}

func NewFile(config *FileConfig, format intf.Format) *File {
	input := new(File)
	input.SetConfig(config)
	input.SetFormat(format)

	return input
}

func (i *File) SetFormat(format intf.Format) {
	i.format = format
}

func (i *File) SetConfig(config *FileConfig) {
	files, err := filepath.Glob(config.Pattern)
	if err != nil {
		Critical("open %s: %v", config.Pattern, err)
	}

	for _, file := range files {
		i.files = append(i.files, i.createBufioReader(file))
	}

	if len(i.files) == 0 {
		i.empty = true
		i.eof = true
	}
}

func (i *File) createBufioReader(filename string) *bufio.Scanner {
	file, err := os.Open(filename)
	if err != nil {
		Critical("open %s: %v", filename, err)
	}

	return bufio.NewScanner(file)
}

func (i *File) GetLine() string {
	if !i.empty && i.scan() {
		return i.files[i.current].Text()
	}

	return ""
}

func (i *File) GetRecord() intf.Record {
	line := i.GetLine()
	if line != "" {
		return i.format.Parse(line)
	}

	return nil
}

func (i *File) scan() bool {
	if !i.files[i.current].Scan() {
		i.current++

		if i.current >= len(i.files) {
			i.eof = true
			return false
		}

		return i.scan()
	}

	return true
}

func (i *File) IsEOF() bool {
	return i.eof
}

func (i *File) Teardown() {
}
