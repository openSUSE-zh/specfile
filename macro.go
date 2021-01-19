package specfile

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	dirutils "github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
)

var (
	macroDirs   = []string{"/usr/lib/rpm/macros.d", "/etc/rpm"}
	macroFiles  = []string{"/usr/lib/rpm/macros", "/usr/lib/rpm/suse/macros"}
	buildConfig = "/usr/lib/build/configs/default.conf"
)

func getFunctionName(str string) string {
	var tmp []byte
	for _, b := range []byte(str) {
		if b == '(' {
			break
		}
		tmp = append(tmp, b)
	}
	return string(tmp)
}

// Macros rpm macros
type Macros []Macro

// Find find a specific macro through macros
func (macros Macros) Find(m Macro) int {
	for i, v := range macros {
		if v.Condition == m.Condition {
			if (v.Type == "variable" && v.Name == m.Name) || (v.Type == "function" && getFunctionName(v.Name) == m.Name) {
				return i
			}
		}
	}
	return -1
}

// Concat concat two macro slice
func (macros *Macros) Concat(macros1 Macros) {
	for _, v := range macros1 {
		if i := macros.Find(v); i >= 0 {
			(*macros)[i].Update(v.Value)
		} else {
			*macros = append(*macros, v)
		}
	}
}

// Macro represent a rpm macro
type Macro struct {
	Indicator string // %global or %define
	Type      string // function or variable
	item
}

// Parse actually parse the macro
func (m *Macro) Parse(str string) error {
	var indicator, name string
	var tmp []byte

	for i, v := range []byte(str) {
		if i == 0 && v != '%' {
			return fmt.Errorf("not a macro")
		}
		if v == '\\' {
			name = string(tmp)
			break
		}
		r, _ := utf8.DecodeRune([]byte{v})
		if unicode.IsSpace(r) {
			if string(tmp) == "%global" || string(tmp) == "%define" {
				indicator = string(tmp)
				tmp = []byte{}
				continue
			} else {
				name = string(tmp)
				break
			}
		}
		tmp = append(tmp, v)
	}

	m.Indicator = indicator
	m.Name = name
	if strings.Contains(name, "(") {
		m.Type = "function"
	} else {
		m.Type = "variable"
	}
	tmp1 := name
	if len(indicator) > 0 {
		tmp1 = indicator + " " + tmp1
	}
	str = strings.Replace(str, tmp1, "", 1)
	str = strings.TrimLeft(str, "\\")
	m.Value = strings.TrimSpace(str)
	m.Name = strings.Replace(m.Name, "%", "", 1)
	return nil
}

// Update update macro definition
func (m *Macro) Update(val string) {
	m.Value = val
}

// ParseMacro parse a macro token into spec
func ParseMacro(token, last Tokenizer, macros Macros, spec *Specfile) {
	var m Macro
	if last.Type == "Comment" {
		m.Comment = last.Content
	}
	m.Raw = &token
	err := m.Parse(token.Content)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	if strings.Contains(m.Value, "%") {
		m.Value = expandMacro(m, macros, spec.Macros, spec.Tags)
	}
	spec.append("Macros", m)
}

// initSystemMacros load system defined rpm macros
func initSystemMacros() Macros {
	var macros Macros
	var files []string
	for _, v := range macroDirs {
		dirs, err1 := dirutils.Glob(filepath.Join(v, "macros.*"))
		if err1 != nil {
			fmt.Printf("can not find macros in %s\n", v)
			continue
		}
		slice.Concat(&files, dirs)
	}
	slice.Concat(&macroFiles, files)
	slice.Concat(&macroFiles, buildConfig)

	for _, v := range macroFiles {
		f, err1 := os.Open(v)
		if err1 != nil {
			fmt.Printf("can not open %s\n", v)
			continue
		}

		var tmp Macros
		if strings.HasSuffix(v, ".conf") {
			tmp, err1 = parseBuildConfig(f)
		} else {
			tmp, err1 = parseMacroFile(f)
		}

		if err1 != nil {
			fmt.Printf("failed to parse %s, %v\n", v, err1)
			break
		}
		macros.Concat(tmp)
		f.Close()
	}

	return macros
}

// parseMacroFile parse .macros file
func parseMacroFile(f io.ReaderAt) (Macros, error) {
	var macros Macros
	err := walkFile(f, false, func(rd io.ReaderAt, line *Line) (error, int64) {
		// skip comment and empty line
		if !strings.HasPrefix(line.Last, "#") && line.Len != 0 {
			var m Macro
			err1 := (&m).Parse(strings.Join(line.Lines, ""))
			// not a macro
			if err1 != nil {
				return err1, line.Offset
			}
			macros = append(macros, m)
		}
		return nil, line.Offset
	}, "Parentheses")

	return macros, err
}

// parseBuildConfig parse /usr/lib/build/configs/default.conf openSUSE specific place for some macros
func parseBuildConfig(f io.ReaderAt) (Macros, error) {
	var macros Macros
	isMacro := false
	err := walkFile(f, false, func(rd io.ReaderAt, line *Line) (error, int64) {
		var macro Macro
		str := strings.Join(line.Lines, "")
		if strings.HasPrefix(str, "%define") {
			err1 := (&macro).Parse(str)
			if err1 != nil {
				return err1, line.Offset
			}
			macros = append(macros, macro)
		}
		if strings.HasPrefix(line.Last, "Macros:") {
			isMacro = true
			return nil, line.Offset
		}
		if strings.HasPrefix(line.Last, ":Macros") {
			isMacro = false
		}

		if isMacro {
			// skip comment and empty line
			if !strings.HasPrefix(line.Last, "#") && line.Lines[0] != "\n" && len(line.Lines[0]) != 0 {
				err1 := (&macro).Parse(str)
				if err1 != nil {
					return err1, line.Offset
				}
				macros = append(macros, macro)
			}
		}

		return nil, line.Offset
	}, "Parentheses")
	return macros, err
}

func expandMacro(macro Macro, system, local Macros, tags []Tag) string {
	// no macro at all
	str := macro.Value
	if !strings.Contains(str, "%") || macro.Type == "function" {
		return str
	}
	if strings.Contains(str, "expand:") {
		str = expand(str)
	}

	var start, useCounter bool
	var idx int
	var c Counter

	var tmp []byte
	var records []string

	for i, v := range []byte(str) {
		if v == '%' {
			start = true
			idx = i
		}
		if start {
			// don't allow nested macro, find the most inner macro first
			if v == '%' {
				tmp = []byte{'%'}
				idx = i
				useCounter = false
				continue
			}
			tmp = append(tmp, v)
			// the next is '{' or '(', we should find the corresponding '}' or ')' to close
			if i == idx+1 && (v == '{' || v == '(') {
				useCounter = true
			}
			// eg '%ix86 x86_64 %arm' stop at whitespace or end of str
			r, _ := utf8.DecodeRune([]byte{v})

			if !useCounter && (unicode.IsSpace(r) || i == len(str)-1) {
				// the space was appended to tmp
				records = append(records, strings.TrimSpace(string(tmp)))
				tmp = []byte{}
				start = false
			}
			if useCounter {
				c.Count(tmp)
				if c.Valid() {
					records = append(records, string(tmp))
					tmp = []byte{}
					useCounter = false
					start = false
				}
				c.Reset()
			}
		}
	}

	for _, v := range records {
		str = strings.Replace(str, v, fillupMacroWithValue(v, system, local, tags), 1)
	}

	// the outer
	if strings.Contains(trim(str), "%") {
		newMacro := macro
		newMacro.Value = str
		newMacro.Type = "variable"
		str = expandMacro(newMacro, system, local, tags)
	}

	// shell commands
	if len(str) > 1 && str[1] == '(' {
		str = callShell(trim(str))
	}
	// macro function
	if len(str) > 1 && str[1] == '{' {
		str = execMacroFunction(str, system, local)
		newMacro := macro
		newMacro.Value = str
		newMacro.Type = "variable"
		str = expandMacro(newMacro, system, local, tags)
	}
	return str
}

func execMacroFunction(s string, system, local Macros) string {
	str := trim(s)
	arr := strings.Split(str, " ")

	if arr[0] == str {
		// not a macro function
		return s
	}

	name := arr[0]
	num := len(arr) - 1
	if i := local.Find(Macro{"", "", item{name, "", "", "", nil}}); i >= 0 {
		val := local[i].Value
		for j := 1; j <= num; j++ {
			if strings.Contains(val, "%{"+strconv.Itoa(j)+"}") {
				val = strings.Replace(val, "%{"+strconv.Itoa(j)+"}", arr[j], -1)
			}
		}
		return val
	}
	return ""
}

// callShell implementation of rpm %()
func callShell(str string) string {
	out, err := exec.Command("/bin/sh", "-c", str).Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}

// expand implementation of rpm %{expand: }
func expand(str string) string {
	idx := strings.LastIndex(str, "expand:")

	if idx < 0 {
		return str
	}

	var c Counter
	replacer := strings.NewReplacer("expand:", "", "%%", "%")

	arr := []byte{'%', '{'}

	for i := idx; i < len(str); i++ {
		arr = append(arr, str[i])
		c.Count(arr)
		if c.Valid() {
			break
		}
		c.Reset()
	}

	s := string(arr)
	tmp := trim(replacer.Replace(s))
	str = strings.Replace(str, s, tmp, 1)

	if strings.Contains(str, "expand:") {
		str = expand(str)
	}

	return str
}

// trim trim the surrounding "%{}"
func trim(str string) string {
	length := len(str)

	if length < 3 {
		return str
	}

	start := 0
	stop := length

	if str[length-1] == '}' || str[length-1] == ')' {
		stop = length - 1
	}

	if str[0] == '%' {
		start = 1
	}

	if str[1] == '{' || str[1] == '(' {
		start = 2
	}

	return str[start:stop]
}

// splitConditionalMacro split conditional macro like "%{!?version:5}" or "%{?version}"
// to the macro "version", default value "5", and a negative symbol
func splitConditionalMacro(str string) (string, string, bool) {
	str = trim(str)
	var neg bool
	var start int
	var defaultValue string

	if str[0] == '!' {
		neg = true
		start = 2
	}

	if str[0] == '?' {
		start = 1
	}

	str = str[start:]

	if strings.Contains(str, ":") {
		arr := strings.Split(str, ":")
		if arr[0] != str {
			str = arr[0]
			defaultValue = arr[1]
		}
	}

	return str, defaultValue, neg
}

func fillupMacroWithValue(str string, system, local Macros, tags []Tag) string {
	str, defaultValue, stat := splitConditionalMacro(str)

	if i := local.Find(Macro{"", "", item{str, "", "", "", nil}}); i >= 0 {
		if stat {
			return ""
		}
		return local[i].Value
	}
	if i := system.Find(Macro{"", "", item{str, "", "", "", nil}}); i >= 0 {
		if stat {
			return ""
		}
		return system[i].Value
	}
	// things like %{name} or %name
	for _, t := range tags {
		if str == strings.ToLower(t.Name) {
			if stat {
				return ""
			}
			return t.Value
		}
	}
	if stat {
		return defaultValue
	}
	return ""
}
