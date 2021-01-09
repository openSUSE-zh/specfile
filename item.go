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

	arr := strings.Split(token.Content, ":")

	i.Name = arr[0]
	i.Value = strings.TrimSpace(arr[1])
}
