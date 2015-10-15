package input

import (
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

import . "gopkg.in/check.v1"

type TailFileSuite struct{}

var _ = Suite(&TailFileSuite{})

func (s *TailFileSuite) TestTailFile(c *C) {
	c.Skip("TODO: fix race condition on tests")

	config := TailConfig{File: "../../tests/resources/tail.b.txt"}

	tail := NewTail(&config, new(MockFormat))
	c.Assert(tail.IsEOF(), Equals, false)

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

		time.Sleep(10000 * time.Microsecond)
		tail.Stop()
	}(tail)

	time.Sleep(3000 * time.Microsecond)

	lines := make([]string, 0)
	for !tail.IsEOF() {
		line := tail.GetLine()
		if line != "" {
			lines = append(lines, line)
		}
	}

	tail.Teardown()
	c.Assert(len(lines), Equals, 20)
}

func (s *TailFileSuite) TestTailFileWithPos(c *C) {
	config := TailConfig{File: "../../tests/resources/tail.a.txt"}

	tail := NewTail(&config, new(MockFormat))
	c.Assert(tail.IsEOF(), Equals, false)

	go func(tail *Tail) {
		filename := "../../tests/resources/tail.a.txt"

		time.Sleep(1000 * time.Microsecond)
		file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Microsecond)
			io.WriteString(file, "foo\n")
		}

		time.Sleep(10000 * time.Microsecond)
		tail.Stop()
	}(tail)

	time.Sleep(1000 * time.Microsecond)

	lines := make([]string, 0)
	for !tail.IsEOF() {
		line := tail.GetLine()
		if line != "" {
			lines = append(lines, line)
		}
	}

	tail.Teardown()
	c.Assert(len(lines), Equals, 10)
}

func (s *TailFileSuite) TestTailFileDelete(c *C) {
	filename := "../../tests/resources/tail.c.txt"
	pos := "../../tests/resources/.tail.c.txt.pos"

	CopyFile("../../tests/resources/plain.a.txt", filename)

	config := TailConfig{File: filename}

	tail := NewTail(&config, new(MockFormat))
	c.Assert(tail.IsEOF(), Equals, false)

	go func(tail *Tail) {
		time.Sleep(100 * time.Microsecond)
		file, _ := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

		for i := 0; i <= 10; i++ {
			time.Sleep(1000 * time.Microsecond)
			io.WriteString(file, "foo\n")
		}

		os.Remove(filename)
		time.Sleep(10000 * time.Microsecond)

		tail.Stop()
	}(tail)

	time.Sleep(1000 * time.Microsecond)

	lines := make([]string, 0)
	for !tail.IsEOF() {
		line := tail.GetLine()
		if line != "" {
			lines = append(lines, line)
		}
	}

	time.Sleep(10000 * time.Microsecond)

	positionRaw, _ := ioutil.ReadFile(pos)
	os.Remove(pos)

	position, _ := strconv.ParseInt(string(positionRaw), 10, 0)
	c.Assert(position, Equals, int64(0))

	tail.Teardown()
	c.Assert(len(lines), Equals, 13)
}

func CopyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	// no need to check errors on read only file, we already got everything
	// we need from the filesystem, so nothing can go wrong now.
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return err
	}
	return d.Close()
}
