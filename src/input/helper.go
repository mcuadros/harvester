package input

import (
	"bufio"
	"io"

	"github.com/mcuadros/harvester/src/intf"
	. "github.com/mcuadros/harvester/src/logger"
)

type ReaderFactory func() io.Reader
type ReaderEOFNotifier func() error

type helper struct {
	factories []ReaderFactory
	readerEOF ReaderEOFNotifier
	format    intf.Format
	current   *bufio.Reader
	empty     bool
	eof       bool
}

func newHelper(format intf.Format) *helper {
	return &helper{
		format:    format,
		factories: make([]ReaderFactory, 0),
	}
}

func (h *helper) GetRecord() intf.Record {
	line := h.getLine()
	if line != "" {
		return h.format.Parse(line)
	}

	return nil
}

func (h *helper) getLine() string {
	if h.current == nil && !h.next() {
		return ""
	}

	line, err := h.current.ReadString('\n')
	if err == io.EOF {
		h.current = nil
		if line != "" {
			return line
		}

		return h.getLine()
	} else if err != nil {
		Error("Error readering: %s", err)
		return ""
	}

	return line[:len(line)-1]
}

func (h *helper) next() bool {
	if len(h.factories) == 0 {
		h.eof = true
		return false
	}

	if h.readerEOF != nil {
		if err := h.readerEOF(); err != nil {
			Error("Error finalizing reader: %s", err)
		}
	}

	reader := h.factories[0]()
	h.factories = h.factories[1:]

	if reader == nil {
		return h.next()
	}

	h.current = bufio.NewReader(reader)
	return true
}

func (h *helper) IsEOF() bool {
	return h.eof
}

func (h *helper) Teardown() {
}
