package specfile

import (
	"strings"
	"unicode"
)

// Section a section defined in sectionMap
type Section struct {
	Belongs string
	item
}

// Parse actually parse the section
func (s *Section) Parse(token *Tokenizer) {
	s.Raw = token

	lines := strings.Split(token.Content, "-n")

	if len(lines) > 1 {
		// "%post -n fcitx5-configtool -p /sbin/ldconfig"
		s.Name = strings.TrimSpace(lines[0])
		var belongs []rune
		for _, r := range strings.TrimSpace(lines[1]) {
			if unicode.IsSpace(r) {
				break
			}
			belongs = append(belongs, r)
		}
		s.Belongs = string(belongs)
		s.Value = strings.TrimSpace(strings.Replace(lines[1], string(belongs), "", 1))
	} else {
		// "%post -p /sbin/ldconfig"
		var name []rune
		for _, r := range token.Content {
			if unicode.IsSpace(r) {
				break
			}
			name = append(name, r)
		}
		s.Name = string(name)
		s.Value = strings.TrimSpace(strings.Replace(token.Content, s.Name, "", 1))
	}
}
