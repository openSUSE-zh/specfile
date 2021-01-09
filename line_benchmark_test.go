package specfile

import (
	"strings"
	"testing"
	"unicode"
	"unicode/utf8"
)

func (line Line) isTag2() bool {
	if len(strings.TrimSpace(line.Last)) == 0 {
		return false
	}
	r, _ := utf8.DecodeRune([]byte(line.Last))
	return unicode.IsUpper(r)
}

func BenchmarkIsTag(b *testing.B) {
	b.ResetTimer()
	line := NewLine(0, "Name: fcitx5")
	for i := 0; i < b.N; i++ {
		line.isTag()
	}
}

func BenchmarkIsTag2(b *testing.B) {
	b.ResetTimer()
	line := NewLine(0, "Name: fcitx5")
	for i := 0; i < b.N; i++ {
		line.isTag2()
	}
}
