package specfile

import (
	"testing"
)

func TestValidCounter(t *testing.T) {
	var c Counter
	c.NextLineConcats = 1
	c.CurlyBrackets = 1
	c.SquareBrackets = 1
	c.Parentheses = 1
	if c.Valid() == true {
		t.Error("[counter]Valid test failed")
	}
}

func TestValidCounterWithOptions(t *testing.T) {
	var c Counter
	c.Parentheses = 1
	if c.Valid("Parentheses") == false {
		t.Error("[counter]Valid with options test failed")
	}
}

func TestResetCounter(t *testing.T) {
	var c Counter
	c.NextLineConcats = 1
	(&c).Reset()
	if c.NextLineConcats != 0 {
		t.Error("[counter]Reset test failed")
	}
}

func TestCountNextLineConcats(t *testing.T) {
	var c Counter
	str := RandStringBytes(10)
	str += "\\"
	(&c).Count([]byte(str))
	if c.NextLineConcats != 1 {
		t.Error("[counter]Count NextLine Concats failed")
	}
}

func TestCountNextLineConcatsInTheMiddle(t *testing.T) {
	var c Counter
	str := RandStringBytes(10)
	str += "\\"
	str += RandStringBytes(10)
	(&c).Count([]byte(str))
	if c.NextLineConcats != 0 {
		t.Error("[counter]Count NextLine Concats failed")
	}
}

func TestCountCurlyBrackets(t *testing.T) {
	var c Counter
	str := "{"
	str += RandStringBytes(10)
	str += "}"
	(&c).Count([]byte(str))
	if c.CurlyBrackets != 0 {
		t.Error("[counter]Count Curly Brackets failed")
	}
}

func TestCountSquareBrackets(t *testing.T) {
	var c Counter
	str := "["
	str += RandStringBytes(10)
	str += "]"
	(&c).Count([]byte(str))
	if c.SquareBrackets != 0 {
		t.Error("[counter]Count Square Brackets failed")
	}
}

func TestCountParentheses(t *testing.T) {
	var c Counter
	str := "("
	str += RandStringBytes(10)
	str += ")"
	(&c).Count([]byte(str))
	if c.Parentheses != 0 {
		t.Error("[counter]Count Parentheses failed")
	}
}
