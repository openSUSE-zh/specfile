package specfile

import (
	"strings"
)

// Element the element interface, holds Specfile and Item
type Element interface {
	GetName() string
	GetValue() string
	GetCondition() string
	GetComment() string
	Inspect() string
}

// Item represent Tag, Macro, and Section
type Item struct {
	Name      string
	Value     string
	Condition string
	Comment   string
	Raw       *Tokenizer
}

// GetName get the name of the item
func (i Item) GetName() string {
	return i.Name
}

// GetValue get the value of the item
func (i Item) GetValue() string {
	return i.Value
}

// GetCondition get the condition of the item
func (i Item) GetCondition() string {
	return i.Condition
}

// GetComment get the comment of the item
func (i Item) GetComment() string {
	return i.Comment
}

// Inspect print the item as string
func (i Item) Inspect() string {
	//return i.Raw.String()
	return ""
}

// Parse parse the item
func (i *Item) Parse(token *Tokenizer) (typ string) {
	i.Raw = token
	switch token.Type {
	case "Macro":
	case "Section":
		// find name from the first line
		first := strings.Split(token.Content, "\n")[0]
		left := strings.Replace(token.Content, first, "", 1)
		var tmp string
		strs := strings.Split(first, " ")
		i.Name = strs[0]
		if len(strs) > 1 {
			if strs[1] == "-n" {
				typ = strs[2]
				tmp = strings.Join(strs[3:], " ")
			} else {
				typ = strs[1]
			}
		}
		i.Value = tmp + left
	default: //Tag
		m := tagRegex.FindStringSubmatch(token.Content)
		i.Name = m[1]
		i.Value = m[2]
	}
	return typ
}
