package specfile

import (
	"fmt"
	"io"
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
func (f *Parser) Parse() {
	var last Tokenizer
	for _, token := range f.tokens {
		switch token.Type {
		case "Conditional":
		case "Tag", "Macro", "Section":
			var item Item
			typ := (&item).Parse(&token)
			fmt.Println(typ)
			if last.Type == "Comment" {
				item.Comment = last.Content
			}
		default:
		}
		last = token
	}
}
