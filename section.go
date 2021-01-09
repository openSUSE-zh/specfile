package specfile

import (
	"strings"
)

// Section a section defined in sectionMap
type Section struct {
	Belongs string
	item
}

// Parse actually parse the section
func (s *Section) Parse(token *Tokenizer) {
	s.Raw = token

	arr := strings.Fields(token.Content)

	if arr[1] == "-n" {
		// "%post -n fcitx5-configtool -p /sbin/ldconfig\n"
		s.Name = arr[0]
		s.Belongs = arr[2]
		s.Value = strings.Join(arr[3:], " ")
	} else {
		// "%post -p /sbin/ldconfig\n"
		s.Name = arr[0]
		s.Value = strings.Join(arr[1:], " ")
	}
}
