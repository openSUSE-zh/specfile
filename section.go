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
		// fix me: whitespace has many variant, may result error
		s.Value = strings.TrimPrefix(token.Content, strings.Join(arr[:3], " "))[1:]
	} else {
		// "%post -p /sbin/ldconfig\n"
		s.Name = arr[0]
		s.Value = strings.TrimPrefix(token.Content, arr[0])[1:]
	}
}

// ParseSection parse section
func ParseSection(token, last Tokenizer, spec *Specfile) error {
	var s Section
	if last.Type == "Comment" {
		s.Comment = last.Content
	}
	s.Parse(&token)

	if len(s.Belongs) > 0 {
		if s.Name == "%package" {
			// subpackage
			parser, err := NewParser(strings.NewReader(s.Value))
			if err != nil {
				return err
			}
			parser.Parse()

			// add subpackage name
			var name Tag
			token := NewTokenizer("Tag", "Name: "+s.Belongs+"\n")
			name.Parse(&token)
			parser.Spec.append("Tags", name)
			spec.append("Subpackages", parser.Spec)
		} else {
			for i := 0; i < len(spec.Subpackages); i++ {
				if t, err := spec.Subpackages[i].FindTag("Name"); err == nil {
					if t.Value == s.Belongs {
						spec.Subpackages[i].append("Sections", s)
					}
				}
			}
		}
	} else {
		spec.append("Sections", s)
	}
	return nil
}
