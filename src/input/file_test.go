package format

import (
	"testing"
)

func TestSingleFile(t *testing.T) {
	config := FileConfig{File: "../../tests/resources/plain.a.txt"}

	file := NewFile(config)
	if file.IsEOF() {
		t.Errorf("FAIL: Wrong IsEOF behavior")
	}

	testReader(t, file, 3)
}

func TestPatternGlob(t *testing.T) {
	config := FileConfig{Pattern: "../../tests/resources/plain.*.txt"}

	file := NewFile(config)
	if file.IsEOF() {
		t.Errorf("FAIL: Wrong IsEOF behavior")
	}

	testReader(t, file, 6)
}

func testReader(t *testing.T, file *File, count int) {
	for i := 0; i <= count; i++ {
		line := file.GetLine()
		if i < count && len(line) == 0 {
			t.Errorf("FAIL: Unable to read a line at position %d", i)
		}

		if i >= count && len(line) != 0 {
			t.Errorf("FAIL: Readed out of scope at position %d. %s", i, line)
		}
	}

	if !file.IsEOF() {
		t.Errorf("FAIL: Wrong IsEOF behavior")
	}
}
