package specfile

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseMacroWithNotMacro(t *testing.T) {
	str := "test"
	var m Macro
	err := (&m).Parse(str)
	if err == nil {
		t.Error("[macro]parse a non-macro but the err is nil")
	}
}

func TestParseMacroWithSingleLineVariable(t *testing.T) {
	str := "%define fcitx5_version 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version, got '%s'", m.Name)
	}
	if m.Value != "5.0.1" {
		t.Errorf("[macro]parsed value is wrong, expected 5.0.1, got '%s'", m.Value)
	}
}

func TestParseMacroWithMultiLineVariable(t *testing.T) {
	str := "%define fcitx5_version\\\n 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version, got '%s'", m.Name)
	}
	if m.Value != "5.0.1" {
		t.Errorf("[macro]parsed value is wrong, expected 5.0.1, got '%s'", m.Value)
	}
}

func TestParseMacroWithSingleLineFunction(t *testing.T) {
	str := "%define fcitx5_version() %(fcitx5 -v)"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version()" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version(), got '%s'", m.Name)
	}
	if m.Value != "%(fcitx5 -v)" {
		t.Errorf("[macro]parsed value is wrong, expected %%(fcitx5 -v), got '%s'", m.Value)
	}
}

func TestParseMacroWithMultiLineFunction(t *testing.T) {
	str := "%define fcitx5_version()\\\n%(fcitx5 -v)"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version()" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version(), got '%s'", m.Name)
	}
	if m.Value != "%(fcitx5 -v)" {
		t.Errorf("[macro]parsed value is wrong, expected %%(fcitx5 -v), got '%s'", m.Value)
	}
}

func TestParseMacroWithSpaceSeparator(t *testing.T) {
	str := "%define fcitx5_version 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version, got '%s'", m.Name)
	}
}

func TestParseMacroWithTabSeparator(t *testing.T) {
	str := "%define\tfcitx5_version 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version, got '%s'", m.Name)
	}
}

func TestParseMacroWithDefinition(t *testing.T) {
	str := "%define fcitx5_version 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "%define" {
		t.Errorf("[macro]parsed indicator is wrong, expected %%define, got '%s'", m.Indicator)
	}
}

func TestParseMacroWithNoDefinition(t *testing.T) {
	str := "%fcitx5_version 5.0.1"
	var m Macro
	err := (&m).Parse(str)
	if err != nil {
		t.Error("[macro]parse a macro but the err is non-nil")
	}
	if m.Indicator != "" {
		t.Errorf("[macro]parsed indicator is wrong, expected empty, got '%s'", m.Indicator)
	}
	if m.Name != "fcitx5_version" {
		t.Errorf("[macro]parsed name is wrong, expected fcitx5_version, got '%s'", m.Name)
	}
	if m.Value != "5.0.1" {
		t.Errorf("[macro]parsed value is wrong, expected 5.0.1, got '%s'", m.Value)
	}
}

func TestParseMacroFile(t *testing.T) {
	macros, err := parseMacroFile(strings.NewReader("%define fcitx5_version() %(fcitx5 -v)\n%fcitx5_name fcitx5"))
	result := Macros{Macro{"%define", "function", item{"fcitx5_version()", "%(fcitx5 -v)", "", "", nil}}, Macro{"", "variable", item{"fcitx5_name", "fcitx5", "", "", nil}}}
	if !reflect.DeepEqual(macros, result) || err != nil {
		t.Error("[macro]parseMacroFile test failed")
	}
}

func TestUpdateMacro(t *testing.T) {
	m := Macro{"", "variable", item{"fcitx_name", "fcitx5", "", "", nil}}
	m.Update("fcitx6")
	if m.Value != "fcitx6" {
		t.Error("[macro]update macro test failed")
	}
}

func TestFindMacro(t *testing.T) {
	m := Macro{"", "variable", item{"%fcitx_name", "fcitx5", "", "", nil}}
	macros := Macros{m}
	if macros.Find(m) < 0 {
		t.Error("[macro]find macro test failed")
	}
}

func TestConcatMacros(t *testing.T) {
	macros := Macros{Macro{"", "variable", item{"%fcitx5_name", "fcitx5", "", "", nil}}}
	macros1 := Macros{Macro{"", "variable", item{"%fcitx5_name", "fcitx6", "", "", nil}}}
	macros.Concat(macros1)
	if !reflect.DeepEqual(macros, macros1) {
		t.Error("[macro]concat macros test failed")
	}
}

func TestParseBuildConfig(t *testing.T) {
	str := "%define gcc_version 5\nConflict: kiwi:systemd-mini\nMacros:\n%rubySTOP() %nil\n:Macros"
	macros, err := parseBuildConfig(strings.NewReader(str))
	result := Macros{Macro{"%define", "variable", item{"gcc_version", "5", "", "", nil}}, Macro{"", "function", item{"rubySTOP()", "%nil", "", "", nil}}}
	if err != nil || !reflect.DeepEqual(macros, result) {
		t.Error("[macro]parseBuildConfig test failed")
	}
}

func TestTrim(t *testing.T) {
	str := "%{%{version}}"
	if trim(str) != "%{version}" {
		t.Error("[macro]trim test failed")
	}
}

func TestSplitConditionalMacro(t *testing.T) {
	str := "%{?version}"
	str, dft, ok := splitConditionalMacro(str)
	if str != "version" || len(dft) != 0 || ok {
		t.Error("[macro]splitConditionalMacro test failed")
	}
}

func TestSplitConditionalMacroWithNoSymbol(t *testing.T) {
	str := "%{version}"
	str, dft, ok := splitConditionalMacro(str)
	if str != "version" || len(dft) != 0 || ok {
		t.Error("[macro]splitConditionalMacro no '!?' or '?' test failed")
	}
}

func TestSplitConditionalMacroWithNonExistence(t *testing.T) {
	str := "%{!?version}"
	str, dft, ok := splitConditionalMacro(str)
	if str != "version" || len(dft) != 0 || !ok {
		t.Error("[macro]splitConditionalMacro with !? test failed")
	}
}

func TestSplitConditionalMacroWithNonExistenceAndDefaultValue(t *testing.T) {
	str := "%{!?version:5}"
	str, dft, ok := splitConditionalMacro(str)
	if str != "version" || dft != "5" || !ok {
		t.Error("[macro]splitConditionalMacro with !? and default value test failed")
	}
}

func TestCallShell(t *testing.T) {
	if callShell("echo true") != "true" {
		t.Error("[macro]callShell test failed")
	}
}

func TestExpand(t *testing.T) {
	str := "%{expand:%%{expand:%%{version}}}"
	if expand(str) != "%{version}" {
		t.Error("[macro]expand macro test failed")
	}
}

func TestExpandWithNoExpand(t *testing.T) {
	str := "%{version}"
	if expand(str) != "%{version}" {
		t.Error("[macro]expand macro with no expand test failed")
	}
}

func TestExpandWithSingleExpand(t *testing.T) {
	str := "%{expand:%%{version}}"
	if expand(str) != "%{version}" {
		t.Error("[macro]expand macro with single expand test failed")
	}
}
