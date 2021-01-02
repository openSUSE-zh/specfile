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
			if strings.Contains(macro.Value, "%") {
				macro.Value = expandMacro(macro, systemMacros, f.Spec.Macros, f.Spec.Tags)
			}
			f.Spec.append("Macros", macro)
		case "Dependency":
			var i Dependency
			if last.Type == "Comment" {
				i.Comment = last.Content
			}
			(&i).Parse(&token)
			f.Spec.append("Dependencies", i)
		case "Section":
			var section Section
			if last.Type == "Comment" {
				section.Comment = last.Content
			}
			(&section).Parse(&token)
			f.Spec.append("Sections", section)
		case "Tag":
			var i Tag
			if last.Type == "Comment" {
				i.Comment = last.Content
			}
			(&i).Parse(&token)
			f.Spec.append("Tags", i)
		}
		last = token
	}
	return nil
}
