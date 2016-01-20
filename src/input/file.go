package input

import (
	"io"
	"os"
	"path/filepath"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
)

type FileConfig struct {
	Format  string
	Pattern string
}

type File struct {
	*helper
}

func NewFile(config *FileConfig, format intf.Format) *File {
	input := &File{newHelper(format)}
	input.SetConfig(config)

	return input
}

func (i *File) SetConfig(config *FileConfig) {
	files, err := filepath.Glob(config.Pattern)
	if err != nil {
		Critical("open %s: %v", config.Pattern, err)
	}

	for _, file := range files {
		i.factories = append(i.factories, i.createReaderFactory(file))
	}
}

func (i *File) createReaderFactory(filename string) ReaderFactory {
	return func() io.Reader {
		file, err := os.Open(filename)
		if err != nil {
			Critical("open %s: %v", filename, err)
		}

		return file
	}
}
