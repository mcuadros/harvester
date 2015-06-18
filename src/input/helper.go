package input

import (
	"bufio"

	"github.com/mcuadros/harvesterd/src/intf"
)

type helper struct {
	files   []*bufio.Scanner
	format  intf.Format
	current int
	empty   bool
	eof     bool
}

func newHelper(format intf.Format) *helper {
	return &helper{
		format: format,
		files:  make([]*bufio.Scanner, 0),
	}
}

func (h *helper) GetLine() string {
	if !h.empty && h.scan() {
		return h.files[h.current].Text()
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
	if !h.files[h.current].Scan() {
		h.current++

		if h.current >= len(h.files) {
			h.eof = true
			return false
		}

		return h.scan()
	}

	return true
}

func (h *helper) IsEOF() bool {
	return h.eof
}

func (h *helper) Teardown() {
}
