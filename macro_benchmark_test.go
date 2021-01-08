package specfile

import (
	"regexp"
	"strings"
	"testing"

	"github.com/marguerite/go-stdlib/slice"
)

func trim1(str string) string {
	var tmp []rune
	for i, v := range str {
		if i == 0 && v == '%' {
			continue
		}
		if i == 1 && (v == '{' || v == '(') {
			continue
		}
		if i == len([]rune(str))-1 && (v == '}' || v == ')') {
			break
		}
		tmp = append(tmp, v)
	}
	return string(tmp)
}

func trim2(str string) string {
	str = strings.TrimLeftFunc(str, func(r rune) bool {
		ok, _ := slice.Contains([]rune{'%', '{', '('}, r)
		return ok
	})
	return strings.TrimRightFunc(str, func(r rune) bool {
		ok, _ := slice.Contains([]rune{'}', ')'}, r)
		return ok
	})
}

func trim3(str string) string {
	return strings.TrimFunc(str, func(r rune) bool {
		return r == '%' || r == '{' || r == '}' || r == '(' || r == ')'
	})
}

func BenchmarkTrim(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim(v)
	}
}

func BenchmarkTrim1(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim1(v)
	}
}

func BenchmarkTrim2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim2(v)
	}
}

func BenchmarkTrim3(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		trim3(v)
	}
}

func splitConditionalMacro1(str string) (string, string, int) {
	str = trim(str)
	stat := 0

	var defaultValue string

	if str[0] == '!' {
		stat = -1
	}

	if str[0] == '?' {
		stat = 1
	}

	str = strings.TrimLeft(str, "!?")

	arr := strings.Split(str, ":")
	if arr[0] != str {
		str = arr[0]
		defaultValue = arr[1]
	}

	return str, defaultValue, stat
}

func splitConditionalMacro2(str string) (string, string, int) {
	str = trim(str)
	stat := 0

	var defaultValue string

	if strings.HasPrefix(str, "!?") {
		stat = -1
		str = strings.TrimPrefix(str, "!?")
	}
	if strings.HasPrefix(str, "?") {
		stat = 1
		str = strings.TrimPrefix(str, "?")
	}

	arr := strings.Split(str, ":")
	if arr[0] != str {
		str = arr[0]
		defaultValue = arr[1]
	}

	return str, defaultValue, stat
}

func splitConditionalMacro3(str string) (string, string, int) {
	str = trim(str)
	stat := 0

	var defaultValue string

	if strings.HasPrefix(str, "!?") {
		stat = -1
		str = strings.TrimPrefix(str, "!?")
	}
	if strings.HasPrefix(str, "?") {
		stat = 1
		str = strings.TrimPrefix(str, "?")
	}

	if strings.Contains(str, ":") {
		re := regexp.MustCompile(`^(.*?):(.*)`)
		m := re.FindStringSubmatch(str)
		str = m[1]
		defaultValue = m[2]
	}

	return str, defaultValue, stat
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

func BenchmarkSplitConditionalMacro2(b *testing.B) {
	b.ResetTimer()
	v := "%{?version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro2(v)
	}
}

func BenchmarkSplitConditionalMacro2_1(b *testing.B) {
	b.ResetTimer()
	v := "%{!?version:5}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro2(v)
	}
}

func BenchmarkSplitConditionalMacro2_2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro2(v)
	}
}

func BenchmarkSplitConditionalMacro3(b *testing.B) {
	b.ResetTimer()
	v := "%{?version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro3(v)
	}
}

func BenchmarkSplitConditionalMacro3_1(b *testing.B) {
	b.ResetTimer()
	v := "%{!?version:5}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro3(v)
	}
}

func BenchmarkSplitConditionalMacro3_2(b *testing.B) {
	b.ResetTimer()
	v := "%{version}"
	for i := 0; i < b.N; i++ {
		splitConditionalMacro3(v)
	}
}
