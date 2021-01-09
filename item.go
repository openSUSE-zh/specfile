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
	var j int
	name := make([]byte, 0, 80)

	for _, v := range []byte(token.Content) {
		if v == ':' {
			break
		}
		name[j] = v
		j++
	}

	i.Name = string(name)
	i.Value = strings.TrimSpace(strings.Replace(token.Content, string(name)+":", "", 1))
}
