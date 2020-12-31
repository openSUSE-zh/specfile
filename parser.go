package specfile

import (
	"io"
	"strings"
)

// Parser the specfile parser
type Parser struct {
	tokens Tokenizers
	spec   Specfile
}

// NewParser initialize a new specfile parser
func NewParser(rd io.ReaderAt) (Parser, error) {
	tokenizers, err := NewTokenizers(rd)
	if err != nil {
		return Parser{}, err
	}
	return Parser{tokenizers, Specfile{}}, nil
}

// Parse actually parse the tokens to spec
func (f *Parser) Parse() error {
	var last Tokenizer
	systemMacros := initSystemMacros()
	for _, token := range f.tokens {
		switch token.Type {
		case "Conditional":
		case "Macro":
			var macro Macro
			macro.Raw = &token
			err := (&macro).Parse(token.Content)
			if err != nil {
				return err
			}
			if strings.Contains(macro.Value, "%") {
				macro.Value = expandMacro(macro.Value, systemMacros, f.spec.Macros)
			}
			tmp := f.spec.Macros
			tmp = append(tmp, macro)
			f.spec.Macros = tmp
		case "Tag", "Section":
			var item Item
			(&item).Parse(&token)
			if last.Type == "Comment" {
				item.Comment = last.Content
			}
		default:
		}
		last = token
	}
	return nil
}
