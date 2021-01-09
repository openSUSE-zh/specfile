package specfile

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func (s *Section) parse1(token *Tokenizer) {
	s.Raw = token

	lines := strings.Split(token.Content, "-n")

	if len(lines) > 1 {
		// "%post -n fcitx5-configtool -p /sbin/ldconfig"
		s.Name = strings.TrimSpace(lines[0])
		var belongs []byte
		for _, v := range []byte(strings.TrimSpace(lines[1])) {
			r, _ := utf8.DecodeRune([]byte{v})
			if unicode.IsSpace(r) {
				break
			}
			belongs = append(belongs, v)
		}
		s.Belongs = string(belongs)
		s.Value = strings.TrimSpace(strings.Replace(lines[1], string(belongs), "", 1))
	} else {
		// "%post -p /sbin/ldconfig"
		var name []byte
		for _, v := range []byte(token.Content) {
			r, _ := utf8.DecodeRune([]byte{v})
			if unicode.IsSpace(r) {
				break
			}
			name = append(name, v)
		}
		s.Name = string(name)
		s.Value = strings.TrimSpace(strings.Replace(token.Content, s.Name, "", 1))
	}
}

var sectiontoken = NewTokenizer("Section", "%post -n fcitx5-configtool -p /sbin/ldconfig")

func BenchmarkSectionParse(b *testing.B) {
	b.ResetTimer()
	var s Section
	for i := 0; i < b.N; i++ {
		s.Parse(&sectiontoken)
	}
}

func BenchmarkSectionParse1(b *testing.B) {
	b.ResetTimer()
	var s Section
	for i := 0; i < b.N; i++ {
		s.parse1(&sectiontoken)
	}
}
