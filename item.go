package specfile

import "strings"

// item holds universal fields
type item struct {
	Name      string
	Value     string
	Condition string
	Comment   string
	Raw       *Tokenizer
}

// Parse parse the item
func (i *item) Parse(token *Tokenizer) {
	i.Raw = token

	// Name: xz
	var name []rune
	for _, v := range token.Content {
		if v == ':' {
			break
		}
		name = append(name, v)
	}

	i.Name = string(name)
	i.Value = strings.TrimSpace(strings.Replace(token.Content, string(name)+":", "", 1))
}
