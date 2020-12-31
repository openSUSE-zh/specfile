package specfile

import (
	"io"
	"strings"
)

// Tokenizers flat list of Tokenizer
type Tokenizers []Tokenizer

// NewTokenizers Dump the specfile plain text to Tokenizers
func NewTokenizers(rd io.ReaderAt) (tokenizers Tokenizers, err error) {
	err = walkFile(rd, true, func(reader io.ReaderAt, line *Line) (error, int64) {
		var c Counter

		if line.isConditional() {
			tmp := line.Offset
			err1 := readConditionalLine(reader, line, &c, 0)
			if err1 != nil && err1 != io.EOF {
				return err1, tmp
			}

			tokenizers = append(tokenizers, NewTokenizer("Conditional", strings.Join(line.Lines, "")))
			if err1 == io.EOF {
				return err1, line.Offset
			}
			return nil, line.Offset
		}

		if line.isSection() {
			for {
				tmp := line.Offset
				err1 := readLine(reader, line, &c)
				if err1 != nil && err1 != io.EOF {
					return err1, tmp
				}
				if line.isSection() {
					line.Lines = line.Lines[:line.Len-1]
					line.Len--
					line.Last = line.Lines[line.Len-1]
					line.Offset = tmp
					break
				}
				if err1 == io.EOF {
					return err1, line.Offset
				}
			}

			// the unclosed if here
			if strings.HasPrefix(line.Last, "%if") || strings.HasPrefix(line.Last, "%else") {
				n := int64(len(line.Last))
				line.Lines = line.Lines[:line.Len-1]
				line.Len--
				line.Last = line.Lines[line.Len-1]
				line.Offset -= n
			}

			tokenizers = append(tokenizers, NewTokenizer("Section", strings.Join(line.Lines, "")))
			return nil, line.Offset
		}

		if line.isMacro() {
			tokenizers = append(tokenizers, NewTokenizer("Macro", strings.Join(line.Lines, "")))
			return nil, line.Offset
		}

		if line.isTag() {
			tokenizers = append(tokenizers, NewTokenizer("Tag", strings.Join(line.Lines, "")))
			return nil, line.Offset
		}

		// empty line
		if strings.Join(line.Lines, "\n") == "\n" {
			tokenizers = append(tokenizers, NewTokenizer("Empty", "\n"))
			return nil, line.Offset
		}
		// comment
		tokenizers = append(tokenizers, NewTokenizer("Comment", strings.Join(line.Lines, "")))

		return nil, line.Offset
	})

	/*var c Counter
	var offset int64

	for {
		line := NewLine(offset)
		err1 := readLine(rd, &line, &c)

		// don't break EOF because we need to handle the read line first
		if err1 != nil && err1 != io.EOF {
			err = err1
			break
		}

		if line.isConditional() {
			err2 := readConditionalLine(rd, &line, &c, 0)
			if err2 != nil && err2 != io.EOF {
				err = err2
				break
			}

			offset = line.Offset
			tokenizers = append(tokenizers, NewTokenizer("Conditional", strings.Join(line.Lines, "")))
			if err2 == io.EOF {
				break
			}
			continue
		}

		if line.isSection() {
			for {
				offset1 := line.Offset
				err3 := readLine(rd, &line, &c)
				if err3 != nil && err3 != io.EOF {
					err = err3
					break
				}
				if line.isSection() {
					line.Lines = line.Lines[:line.Len-1]
					line.Len--
					line.Last = line.Lines[line.Len-1]
					line.Offset = offset1
					break
				}
				if err3 == io.EOF {
					break
				}
			}

			// the unclosed if here
			if strings.HasPrefix(line.Last, "%if") || strings.HasPrefix(line.Last, "%else") {
				n := int64(len(line.Last))
				line.Lines = line.Lines[:line.Len-1]
				line.Len--
				line.Last = line.Lines[line.Len-1]
				line.Offset -= n
			}

			offset = line.Offset
			tokenizers = append(tokenizers, NewTokenizer("Section", strings.Join(line.Lines, "")))
			continue
		}

		if line.isMacro() {
			offset = line.Offset
			tokenizers = append(tokenizers, NewTokenizer("Macro", strings.Join(line.Lines, "")))
			continue
		}

		if line.isTag() {
			offset = line.Offset
			tokenizers = append(tokenizers, NewTokenizer("Tag", strings.Join(line.Lines, "")))
			continue
		}

		offset = line.Offset
		// empty line
		if strings.Join(line.Lines, "\n") == "\n" {
			tokenizers = append(tokenizers, NewTokenizer("Empty", "\n"))
			continue
		}
		// comment
		tokenizers = append(tokenizers, NewTokenizer("Comment", strings.Join(line.Lines, "")))

		// break the EOF here
		if err1 == io.EOF {
			break
		}
	}*/

	return tokenizers, err
}

// Tokenizer like Tokenizer{"Macro", "%define fcitx5_version 5.0.1\n"}
type Tokenizer struct {
	Type    string
	Content string
}

// NewTokenizer return a new tokenizer
func NewTokenizer(typ, content string) Tokenizer {
	return Tokenizer{typ, content}
}

// String return the raw content
/*func (token Tokenizer) String() string {
	return token.Content
}*/
