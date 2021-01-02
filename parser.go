package specfile

import (
	"io"
	"strings"
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
			var macro Macro
			if last.Type == "Comment" {
				macro.Comment = last.Content
			}
			macro.Raw = &token
			err := (&macro).Parse(token.Content)
			if err != nil {
				return err
			}
			if strings.Contains(macro.Value, "%") && macro.Type != "function" {
				macro.Value = expandMacro(macro.Value, systemMacros, f.Spec.Macros)
			}
			f.Spec.append("Macros", macro)
		case "Dependency":
			var item item
			if last.Type == "Comment" {
				item.Comment = last.Content
			}
			(&item).Parse(&token)
			f.Spec.append("Dependencies", item)
		case "Section":
			var section Section
			if last.Type == "Comment" {
				section.Comment = last.Content
			}
			(&section).Parse(&token)
			f.Spec.append("Sections", section)
		case "Tag":
			var item item
			if last.Type == "Comment" {
				item.Comment = last.Content
			}
			(&item).Parse(&token)
			f.Spec.append("Tags", item)
		}
		last = token
	}
	return nil
}
