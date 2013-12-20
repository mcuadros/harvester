package input

import (
	"io"
	"os"
	"time"
)

import . "launchpad.net/gocheck"

type TailFileSuite struct{}

var _ = Suite(&TailFileSuite{})

func (s *TailFileSuite) TestTailFile(c *C) {
	config := TailConfig{File: "../../tests/resources/tail.b.txt"}

	tail := NewTail(&config, new(MockFormat))
	c.Check(tail.IsEOF(), Equals, false)

	go func(tail *Tail) {
		filename := "../../tests/resources/tail.b.txt"
		pos := "../../tests/resources/.tail.b.txt.pos"

		file, _ := os.Create(filename)
		time.Sleep(100 * time.Microsecond)

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
		os.Remove(pos)

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

	c.Check(lines, HasLen, 20)
}

func (s *TailFileSuite) TestTailFileWithPos(c *C) {

	config := TailConfig{File: "../../tests/resources/tail.a.txt"}

	tail := NewTail(&config, new(MockFormat))
	c.Check(tail.IsEOF(), Equals, false)

	go func(tail *Tail) {
		filename := "../../tests/resources/tail.a.txt"

		time.Sleep(100 * time.Microsecond)
		file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

		for i := 0; i < 10; i++ {
			time.Sleep(1000 * time.Microsecond)
			io.WriteString(file, "foo\n")
		}

		tail.Stop()
	}(tail)

	lines := make([]string, 0)
	for !tail.IsEOF() {
		line := tail.GetLine()
		if line != "" {
			lines = append(lines, line)
		}
	}

	c.Check(len(lines), Equals, 10)
}
