package specfile

import (
	"strings"
	"testing"
)

func trim2(str string) string {
	str = strings.TrimLeftFunc(str, func(r rune) bool {
		return r == '%' || r == '{' || r == '('
	})
	return strings.TrimRightFunc(str, func(r rune) bool {
		return r == '}' || r == ')'
	})
}

func BenchmarkTrim(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim(v)
	}
}

func BenchmarkTrim2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim2(v)
	}
}

func splitConditionalMacro1(str string) (string, string, bool) {
	str = trim(str)
	var neg bool
	var defaultValue string

	// do the ?! and ? judge
	if strings.HasPrefix(str, "!?") {
		neg = true
		str = strings.TrimPrefix(str, "!?")
	}
	if strings.HasPrefix(str, "?") {
		str = strings.TrimPrefix(str, "?")
	}

	if strings.Contains(str, ":") {
		arr := strings.Split(str, ":")
		if arr[0] != str {
			str = arr[0]
			defaultValue = arr[1]
		}
	}

	return str, defaultValue, neg
}

func BenchmarkSplitConditionalMacro(b *testing.B) {
	b.ResetTimer()
	v := "%{?version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro(v)
	}
}

func BenchmarkSplitConditionalMacro_1(b *testing.B) {
	b.ResetTimer()
	v := "%{!?version:5}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro(v)
	}
}

func BenchmarkSplitConditionalMacro_2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro(v)
	}
}

func BenchmarkSplitConditionalMacro1(b *testing.B) {
	b.ResetTimer()
	v := "%{?version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro1(v)
	}
}

func BenchmarkSplitConditionalMacro1_1(b *testing.B) {
	b.ResetTimer()
	v := "%{!?version:5}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro1(v)
	}
}

func BenchmarkSplitConditionalMacro1_2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro1(v)
	}
}
