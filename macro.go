package specfile

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	dirutils "github.com/marguerite/util/dir"
	"github.com/marguerite/util/slice"
)

var (
	macroDirs   = []string{"/usr/lib/rpm/macros.d", "/etc/rpm"}
	macroFiles  = []string{"/usr/lib/rpm/macros", "/usr/lib/rpm/suse/macros"}
	buildConfig = "/usr/lib/build/configs/default.conf"
)

// Macros rpm macros
type Macros []Macro

// Find find a specific macro through macros
func (macros Macros) Find(m Macro) int {
	for i, v := range macros {
		if v.Indicator == m.Indicator &&
			v.Type == m.Type &&
			v.Name == m.Name &&
			v.Conditional == m.Conditional {
			return i
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
	Indicator   string // %global or %define
	Type        string // function or variable
	Name        string
	Value       string
	Conditional string
	Raw         *Tokenizer
}

// Parse actually parse the macro
func (m *Macro) Parse(str string) error {
	var indicator, name string
	var tmp []byte

	bytes := []byte(str)

	for i := 0; i < len(bytes); i++ {
		if i == 0 && bytes[i] != '%' {
			return fmt.Errorf("not a macro")
		}
		if bytes[i] == '\\' {
			name = string(tmp)
			break
		}
		if bytes[i] == ' ' || bytes[i] == '\t' {
			if string(tmp) == "%global" || string(tmp) == "%define" {
				indicator = string(tmp)
				tmp = []byte{}
				continue
			} else {
				name = string(tmp)
				break
			}
		}
		tmp = append(tmp, bytes[i])
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
	return nil
}

// Update update macro definition
func (m *Macro) Update(val string) {
	m.Value = val
}

// InitSystemMacros load system defined rpm macros
func InitSystemMacros() Macros {
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
		if !strings.HasPrefix(line.Last, "#") && line.Lines[0] != "\n" && len(line.Lines[0]) != 0 {
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
