package specfile

import (
	"fmt"
	"io"
	"strings"
)

// readConditionalLine continously read from reader until the current Line has no unclosed conditions
func readConditionalLine(reader io.ReaderAt, line *Line, c *Counter, numIf int) error {
	if strings.HasPrefix(line.Last, "%if") {
		numIf++
	}
	if strings.HasPrefix(line.Last, "%end") {
		numIf--
	}

	if numIf > 0 {
		old := line.Lines
		line.Lines = []string{}
		line.Len = 0
		err := readLine(reader, line, c)
		line.Concat(true, old...)
		// err for the second loop
		if err != nil {
			return err
		}
		err = readConditionalLine(reader, line, c, numIf)
		// actually break the second loop here
		if err != nil {
			return err
		}
	}
	return nil
}

// walkFile read the whole file and let you do things Line by Line
func walkFile(reader io.ReaderAt, brk bool, fn func(rd io.ReaderAt, line *Line) (error, int64), readLineOptions ...string) error {
	var err error
	var offset int64
	var c Counter

	for {
		line := NewLine(offset)
		err1 := readLine(reader, &line, &c, readLineOptions...)
		// don't break EOF because we need to handle the read line first
		if err1 != nil && err1 != io.EOF {
			err = err1
			break
		}

		err2, offset1 := fn(reader, &line)
		offset = offset1

		if err2 != nil {
			if err2 == io.EOF || err1 == io.EOF {
				break
			}
			if brk {
				err = err2
				break
			} else {
				continue
			}
		}
		// break the EOF here
		if err1 == io.EOF {
			break
		}
	}

	return err
}

// readLine read syntactically valid line from io.ReaderAt with no unclosed brackets
func readLine(reader io.ReaderAt, line *Line, c *Counter, options ...string) error {
	b, err := read(reader, line)
	c.Count(b)
	line.Concat(false, string(b))

	if err != nil {
		if !c.Valid(options...) {
			c.Reset()
			return fmt.Errorf("The line is incomplete in syntax: %s", strings.Join(line.Lines, ""))
		}

		c.Reset()
		if err == io.EOF {
			return err
		}
		return fmt.Errorf("The read operation is partial successful: %s", strings.Join(line.Lines, ""))
	}

	if !c.Valid(options...) {
		if c.NextLineConcats != 0 {
			c.NextLineConcats--
		}
		err = readLine(reader, line, c, options...)
		if err != nil {
			return err
		}
	}

	c.Reset()
	return nil
}

// read a valid line from io.ReaderAt
func read(reader io.ReaderAt, line *Line) (bytes []byte, err error) {
	for {
		// sometimes the next-line concat symbol "\" doesn't follow immediately and is very very far away,
		// I think a 255 buf is enough to reach it
		buf := make([]byte, 255)
		n, err1 := reader.ReadAt(buf, line.Offset)

		// save every byte we read
		found := false
		for i := 0; i < n; i++ {
			bytes = append(bytes, buf[i])
			if buf[i] == '\n' {
				line.Offset += int64(i + 1)
				found = true
				break
			}
		}
		if found {
			break
		}

		line.Offset += int64(n)

		// when n == len(buf), the err can be EOF or nil
		if n == 255 && err1 == nil {
			break
		}

		if err1 != nil {
			err = err1
			break
		}
	}

	// bytes may be empty
	return bytes, err
}
