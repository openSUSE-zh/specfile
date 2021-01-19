package specfile

import (
	"io"
)

// Parser the specfile parser
type Parser struct {
	Tokens Tokenizers
	Spec   Specfile
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
	for _, token := range f.Tokens {
		switch token.Type {
		case "Conditional":
		case "Macro":
			ParseMacro(token, last, systemMacros, &(f.Spec))
		case "Dependency":
			ParseDependency(token, last, &(f.Spec))
		case "Section":
			err := ParseSection(token, last, &(f.Spec))
			if err != nil {
				return err
			}
		case "Tag":
			ParseTag(token, last, &(f.Spec))
		}
		last = token
	}
	return nil
}
