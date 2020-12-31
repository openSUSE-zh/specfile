package specfile

import (
	"strings"
	"testing"
)

func TestNewTokenizer(t *testing.T) {
	tokenizer := NewTokenizer("Tag", "test")
	if tokenizer.Type != "Tag" || tokenizer.Content != "test" {
		t.Error("[lexer]NewTokenizer test failed")
	}
}

/*func TestTokenizerToString(t *testing.T) {
	tokenizer := NewTokenizer("Tag", "BuildRequires: xz")
	if tokenizer.String() != "BuildRequires: xz" {
		t.Error("[lexer]tokenizer to String test failed")
	}
}*/

func TestConditionalTokenizer(t *testing.T) {
	str := "%if 0%{?suse_version} > 1550\nBuildRequires: xz\n%endif"
	tokenizers, err := NewTokenizers(strings.NewReader(str))
	if err != nil || tokenizers[0].Type != "Conditional" || tokenizers[0].Content != str {
		t.Error("[lexer]conditional tokenizer test failed")
	}
}

func TestMacroTokenizer(t *testing.T) {
	str := "%global suse_version 1550"
	tokenizers, _ := NewTokenizers(strings.NewReader(str))
	if tokenizers[0].Type != "Macro" || tokenizers[0].Content != str {
		t.Error("[lexer]macro tokenizer test failed")
	}
}

func TestTagTokenizer(t *testing.T) {
	str := "BuildRequires: xz"
	tokenizers, _ := NewTokenizers(strings.NewReader(str))
	if tokenizers[0].Type != "Tag" || tokenizers[0].Content != str {
		t.Error("[lexer]tag tokenizer test failed")
	}
}

func TestEmptyTokenizer(t *testing.T) {
	str := "\n"
	tokenizers, _ := NewTokenizers(strings.NewReader(str))
	if tokenizers[0].Type != "Empty" || tokenizers[0].Content != str {
		t.Error("[lexer]empty tokenizer test failed")
	}
}

func TestCommentTokenizer(t *testing.T) {
	str := "# spec file for package gcc10"
	tokenizers, err := NewTokenizers(strings.NewReader(str))
	if err != nil || tokenizers[0].Type != "Comment" || tokenizers[0].Content != str {
		t.Error("[lexer]comment tokenizer test failed")
	}
}

func TestSectionTokenizer(t *testing.T) {
	str := "%description\nCore package for the GNU Compiler Collection, including the C language frontend.\n\n"
	str1 := "%package "
	tokenizers, err := NewTokenizers(strings.NewReader(str + str1))
	if err != nil || tokenizers[0].Type != "Section" || tokenizers[0].Content != str {
		t.Error("[lexer]section tokenizer test failed")
	}
}

func TestSectionTokenizerWithUnclosedIf(t *testing.T) {
	str := "%description\nCore package for the GNU Compiler Collection, including the C language frontend.\n\n"
	str1 := "%if 0%{?suse_version} > 1550\n%package "
	tokenizers, err := NewTokenizers(strings.NewReader(str + str1))
	if err != nil || tokenizers[0].Type != "Section" || tokenizers[0].Content != str {
		t.Error("[lexer]section tokenizer test failed")
	}
}
