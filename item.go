package specfile

import (
	"strings"
)

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
        idx := strings.Index(token.Content, ":")

	i.Name = token.Content[:idx]
	i.Value = strings.TrimSpace(token.Content[idx+1:])
}
