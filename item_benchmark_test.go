package specfile

import (
	"strings"
	"testing"
)

func (i *item) parse1(token *Tokenizer) {
	i.Raw = token

	// Name: xz
	var j int
	name := make([]byte, strings.Index(token.Content, ":"))

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

func (i *item) parse2(token *Tokenizer) {
   i.Raw = token
   arr := strings.Split(token.Content, ":")
   i.Name = arr[0]
   i.Value = strings.TrimSpace(arr[1])
}

var itemtoken = NewTokenizer("Tag", "Name: fcitx5\n")

func BenchmarkItemParse(b *testing.B) {
	b.ResetTimer()
	var it item
	for i := 0; i < b.N; i++ {
		it.Parse(&itemtoken)
	}
}

func BenchmarkItemParse1(b *testing.B) {
	b.ResetTimer()
	var it item
	for i := 0; i < b.N; i++ {
		it.parse1(&itemtoken)
	}
}

func BenchmarkItemParse2(b *testing.B) {
	b.ResetTimer()
	var it item
	for i := 0; i < b.N; i++ {
		it.parse2(&itemtoken)
	}
}
