package harvesterd

import (
	. "harvesterd/intf"
)

type WriterGroup struct {
	writers []Writer
}

func NewWriterGroup() *WriterGroup {
	writerGroup := new(WriterGroup)

	return writerGroup
}

func (self *WriterGroup) SetWriters(writers []Writer) {
	self.writers = writers
}

func (self *WriterGroup) Setup() {
	for _, writer := range self.writers {
		writer.Setup()
	}
}

func (self *WriterGroup) Boot() {
	for _, writer := range self.writers {
		writer.Boot()
	}
}

func (self *WriterGroup) IsAlive() bool {
	for _, writer := range self.writers {
		if writer.IsAlive() {
			return true
		}
	}

	return false
}

func (self *WriterGroup) GetCounters() (int32, int32, int32, int32) {
	var created, failed, transferred, threads int32
	for _, writer := range self.writers {
		c, f, t, h := writer.GetCounters()

		created += c
		failed += f
		transferred += t
		threads += h
	}

	return created, failed, transferred, threads
}

func (self *WriterGroup) ResetCounters() {
	for _, writer := range self.writers {
		writer.ResetCounters()
	}
}

func (self *WriterGroup) Teardown() {
	for _, writer := range self.writers {
		writer.Teardown()
	}
}
