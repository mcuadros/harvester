package harvesterd

import "github.com/mcuadros/harvesterd/src/intf"

type WriterGroup struct {
	writers []intf.Writer
}

func NewWriterGroup() *WriterGroup {
	writerGroup := new(WriterGroup)

	return writerGroup
}

func (wg *WriterGroup) SetWriters(writers []intf.Writer) {
	wg.writers = writers
}

func (wg *WriterGroup) Setup() {
	for _, writer := range wg.writers {
		writer.Setup()
	}
}

func (wg *WriterGroup) Boot() {
	for _, writer := range wg.writers {
		writer.Boot()
	}
}

func (wg *WriterGroup) IsAlive() bool {
	for _, writer := range wg.writers {
		if writer.IsAlive() {
			return true
		}
	}

	return false
}

func (wg *WriterGroup) GetCounters() (int32, int32, int32, int32) {
	var created, failed, transferred, threads int32
	for _, writer := range wg.writers {
		c, f, t, h := writer.GetCounters()

		created += c
		failed += f
		transferred += t
		threads += h
	}

	return created, failed, transferred, threads
}

func (wg *WriterGroup) ResetCounters() {
	for _, writer := range wg.writers {
		writer.ResetCounters()
	}
}

func (wg *WriterGroup) Teardown() {
	for _, writer := range wg.writers {
		writer.Teardown()
	}
}
