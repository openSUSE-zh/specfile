package specfile

import (
	"reflect"
	"testing"
)

func TestNewLineWithSingleLine(t *testing.T) {
	line := NewLine(0, "test")
	if line.Last != "test" || line.Offset != 0 || line.Len != 1 {
		t.Error("[line]NewLine with single line test failed")
	}
}

func TestNewLineWithMultipleLines(t *testing.T) {
	line := NewLine(0, "test", "test1")
	if !reflect.DeepEqual(line.Lines, []string{"test", "test1"}) || line.Last != "test1" || line.Offset != 0 || line.Len != 2 {
		t.Error("[line]NewLine with multiple lines test failed")
	}
}

func TestConditionalLine(t *testing.T) {
	for _, c := range conditionalIndicators {
		line := NewLine(0, c)
		if !line.isConditional() {
			t.Error("[line]Conditional test failed")
		}
	}
}

func TestSectionLine(t *testing.T) {
	line := NewLine(0, "%description\n")
	if !line.isSection() {
		t.Error("[line]Section test failed")
	}
}

func TestMacroLine(t *testing.T) {
	line := NewLine(0, "%define fcitx5_version 5.0.1\n")
	if !line.isMacro() {
		t.Error("[line]Macro test failed")
	}
}

func TestTagLine(t *testing.T) {
	line := NewLine(0, "BuildRequires: xz\n")
	if !line.isTag() {
		t.Error("[line]Tag test failed")
	}
}

func TestAppendLineWithSingleLine(t *testing.T) {
	line := NewLine(0, "test")
	line.Concat(false, "test1")
	if !reflect.DeepEqual(line.Lines, []string{"test", "test1"}) || line.Last != "test1" || line.Len != 2 {
		t.Error("[line]Line.Concat with single line append test failed")
	}
}

func TestAppendLineWithMultipleLines(t *testing.T) {
	line := NewLine(0, "test")
	line.Concat(false, "test1", "test2")
	if !reflect.DeepEqual(line.Lines, []string{"test", "test1", "test2"}) || line.Last != "test2" || line.Len != 3 {
		t.Error("[line]Line.Concat with multiple lines append test failed")
	}
}

func TestPrependLineWithSingleLine(t *testing.T) {
	line := NewLine(0, "test")
	line.Concat(true, "test1")
	if !reflect.DeepEqual(line.Lines, []string{"test1", "test"}) || line.Last != "test" || line.Len != 2 {
		t.Error("[line]Line.Concat with single line prepend test failed")
	}
}

func TestPrependLineWithMultipleLines(t *testing.T) {
	line := NewLine(0, "test")
	line.Concat(true, "test1", "test2")
	if !reflect.DeepEqual(line.Lines, []string{"test1", "test2", "test"}) || line.Last != "test" || line.Len != 3 {
		t.Error("[line]Line.Concat with multiple lines prepend test failed")
	}
}
