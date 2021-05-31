package specfile

import (
	"bufio"
	"strings"
)

// Conditional the conditionals like "%if 0%{?suse_version} >= 1550\nBuildRequires: xz\n%endif\n"
type Conditional struct {
	item
}

// Parse parse the conditional
func (conditional *Conditional) Parse(token *Tokenizer) {

}

// ParseConditional parse conditional from token
func ParseConditional(token, last Tokenizer, macros Macros, spec *Specfile) {
	scanner := bufio.NewScanner(strings.NewReader(token.Content))
	var conditions []string
	var lastIf, secondlastIf, idx int

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "%if") || strings.HasPrefix(line, "%else") {
			conditions = append(conditions, line)
			if strings.HasPrefix(line, "%if") {
				secondlastIf = lastIf
				lastIf = idx
			}
			idx++
			continue
		}
		if strings.HasPrefix(line, "%endif") {
			// ["%if 0%{?suse_version} >= 1310", "%if 0%{?suse_version} > 1330"]
			conditions = conditions[:lastIf]
			lastIf = secondlastIf
			secondlastIf = 0
			idx -= len(conditions) - lastIf
			continue
		}
		parser, _ := NewParser(strings.NewReader(line))
		parser.Parse()
	}
}
