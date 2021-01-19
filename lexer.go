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
				last := line.Last
				err1 := readLine(reader, line, &c)
				if err1 != nil && err1 != io.EOF {
					return err1, tmp
				}
				if line.isSection() && line.Last != last {
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

		if line.isDependency() {
			tokenizers = append(tokenizers, NewTokenizer("Dependency", strings.Join(line.Lines, "")))
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
		if line.Len > 0 {
			tokenizers = append(tokenizers, NewTokenizer("Comment", strings.Join(line.Lines, "")))
		}

		return nil, line.Offset
	})

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
