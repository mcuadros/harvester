package format

import (
	"io"
	"os"
	"testing"
	"time"
)

func TestTailFile(t *testing.T) {
	config := TailConfig{File: "../tests/resources/tail.b.txt"}

	tail := NewTail(config)
	if tail.IsEOF() {
		t.Errorf("FAIL: Wrong IsEOF behavior")
	}

	go func(tail *Tail) {
		filename := "../tests/resources/tail.b.txt"

		time.Sleep(100 * time.Microsecond)
		file, _ := os.Create(filename)

		for i := 0; i < 10; i++ {
			time.Sleep(1000 * time.Microsecond)
			io.WriteString(file, "foo\n")
		}

		file.Close()

		file, _ = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

		for i := 0; i < 10; i++ {
			time.Sleep(1000 * time.Microsecond)
			io.WriteString(file, "foo\n")
		}

		time.Sleep(100 * time.Microsecond)
		os.Remove(filename)
		time.Sleep(100 * time.Microsecond)
		tail.Stop()
	}(tail)

	lines := make([]string, 0)
	for !tail.IsEOF() {
		line := tail.GetLine()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) != 20 {
		t.Errorf("Incorrect number of lines: %d should be %d", len(lines), 20)
	}
}
