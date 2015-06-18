package input

import (
	"bufio"
	"io"

	"github.com/mcuadros/harvesterd/src/intf"
	. "github.com/mcuadros/harvesterd/src/logger"
)

type ReaderFactory func() io.Reader
type ReaderEOFNotifier func() error

type helper struct {
	factories []ReaderFactory
	readerEOF ReaderEOFNotifier
	format    intf.Format
	current   *bufio.Scanner
	empty     bool
	eof       bool
}

func newHelper(format intf.Format) *helper {
	return &helper{
		format:    format,
		factories: make([]ReaderFactory, 0),
	}
}

func (h *helper) GetLine() string {
	if h.scan() {
		return h.current.Text()
	}

	return ""
}

func (h *helper) GetRecord() intf.Record {
	line := h.GetLine()
	if line != "" {
		return h.format.Parse(line)
	}

	return nil
}

func (h *helper) scan() bool {
	if h.current != nil && h.current.Scan() {
		return true
	}

	if len(h.factories) == 0 {
		h.eof = true
		return false
	}

	if h.readerEOF != nil {
		err := h.readerEOF()
		if err != nil {
			Error("Error finalizing reader: %s", err)
		}
	}

	reader := h.factories[0]()
	h.factories = h.factories[1:]

	if reader == nil {
		return h.scan()
	}

	h.current = bufio.NewScanner(reader)
	return h.scan()
}

func (h *helper) IsEOF() bool {
	return h.eof
}

func (h *helper) Teardown() {
}
