package specfile

import (
	"bytes"
	"fmt"
	"io"
)

// Line a syntax-complete line
type Line struct {
	Num int
	ptr *bytes.Buffer
}

// Lines collection of Line
type Lines []Line

// NewLine initialize a new Line
func NewLine(num int, ptr *bytes.Buffer) Line {
	return Line{num, ptr}
}

// Debug debug lines
func (lines Lines) Debug() string {
	var str string
	for _, v := range lines {
		str += fmt.Sprintf("%d\n%s", v.Num, v.ptr.String())
	}
	return str
}

// FindByNum find a Line's content and position
func (lines Lines) FindByNum(i int) ([]byte, int) {
	for j, v := range lines {
		if v.Num == i {
			return v.ptr.Bytes(), j
		}
	}
	return []byte{}, 0
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
			// "\" means to concat a line
			if j := bytes.IndexByte(b[:i], '\\'); j < 0 {
				n1, err1 := buf.Write(b[:i+1])
				if n1 != i+1 {
					panic(fmt.Sprintf("can not write byte, content %s, length %d, wrote %d\n", b, i+1, n1))
				}
				if err1 == bytes.ErrTooLarge {
					panic(err1)
				}

				x := bytes.Count(buf.Bytes(), []byte("%if"))
				y := bytes.Count(buf.Bytes(), []byte("%endif"))

				if x == 0 || x == y {
					lines = append(lines, NewLine(idx, buf))
					idx++
					buf = bytes.NewBuffer([]byte{})
				}
				b = make([]byte, 32)
				off += int64(i + 1)
				continue
			}
		}

		n1, err1 := buf.Write(b)
		if n1 != 32 {
			panic(fmt.Sprintf("can not write byte, content %s, length %d, wrote %d\n", b, 32, n1))
		}
		if err1 == bytes.ErrTooLarge {
			panic(err1)
		}
		off += int64(32)
		b = make([]byte, 32)
	}

	return lines
}
