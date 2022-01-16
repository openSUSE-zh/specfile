package specfile

import (
	"bytes"
	"fmt"
	"io"
)

var (
	lineTypeMap = map[string]int{
		"%ifarch":   1,
		"%ifnarch":  2,
		"%ifos":     3,
		"%ifnos":    4,
		"%if":       5,
		"%else":     6,
		"%elifarch": 7,
		"%elifos":   8,
		"%elif":     9,
		"%endif":    10,
		"%include":  11,
	}
)

// Line a line
type Line struct {
	num int
	typ int
	len int
	buf []byte
}

// Lines collection of Line
type Lines []Line

func parseLineType(buf []byte) int {
	for k, v := range lineTypeMap {
		if bytes.HasPrefix(buf, []byte(k)) {
			return v
		}
	}

	// comment
	if buf[0] == '#' {
		return -1
	}

	// empty
	if len(buf) == 1 {
		return -2
	}

	// return default line type
	return 0
}

// String debug lines
func (lines Lines) String() string {
	var str string
	for _, v := range lines {
		str += fmt.Sprintf("line %d, type %d, len %d:\n%s", v.num, v.typ, v.len, v.buf)
	}
	return str
}

func write(buf *bytes.Buffer, b []byte) {
	n, err := buf.Write(b)
	if n != len(b) {
		panic(fmt.Sprintf("can not write byte, content %s, length %d, wrote %d\n", b, len(b), n))
	}
	if err == bytes.ErrTooLarge {
		panic(err)
	}
}

// ReadLines read lines of file
func ReadLines(f io.ReaderAt) (lines Lines) {
	b := make([]byte, 32)
	buf := bytes.NewBuffer([]byte{})
	idx := 1
	var off int64

	for {
		n, err := f.ReadAt(b, off)
		if n == 0 && err == io.EOF {
			break
		}

		if i := bytes.IndexByte(b, '\n'); i >= 0 {
			write(buf, b[:i+1])
			b1 := buf.Bytes()
			lines = append(lines, Line{idx, parseLineType(b1), len(b1), b1})
			idx++
			off += int64(i + 1)
			buf = bytes.NewBuffer([]byte{})
			b = make([]byte, 32)
			continue
		}

		write(buf, b)
		off += int64(32)
		b = make([]byte, 32)
	}

	return lines
}
