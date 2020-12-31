package specfile

import (
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func TestReadFull(t *testing.T) {
	str := RandStringBytes(80)
	var line Line
	b, err := read(strings.NewReader(str), &line)
	if !reflect.DeepEqual(b, []byte(str)) || err != nil || line.Offset != 80 {
		t.Error("[readline]read full buf test failed")
	}
}

func TestReadPartial(t *testing.T) {
	str := RandStringBytes(70)
	var line Line
	b, err := read(strings.NewReader(str), &line)
	if !reflect.DeepEqual(b, []byte(str)) || err != io.EOF || line.Offset != 70 {
		t.Error("[readline]read partial buf test failed")
	}
}

func TestReadFullWithBreak(t *testing.T) {
	str := RandStringBytes(79)
	str += "\n"
	var line Line
	b, err := read(strings.NewReader(str), &line)
	if !reflect.DeepEqual(b, []byte(str)) || err != nil || line.Offset != 80 {
		t.Error("[readline]read full buf test failed")
	}
}

func TestReadLine(t *testing.T) {
	str := RandStringBytes(10)
	var line Line
	var c Counter
	err := readLine(strings.NewReader(str), &line, &c)
	if err != io.EOF || line.Last != str || line.Len != 1 || !c.Valid() {
		t.Error("[readline]readLine with single line test failed")
	}
}

func TestReadLineWithMultipleLines(t *testing.T) {
	str := RandStringBytes(10)
	str += "\\\n"
	str1 := RandStringBytes(10)
	var line Line
	var c Counter
	err := readLine(strings.NewReader(str+str1), &line, &c)
	if err != io.EOF || !reflect.DeepEqual(line.Lines, []string{str, str1}) || line.Last != str1 || line.Len != 2 || !c.Valid() {
		t.Error("[readline]readLine with multple lines test failed")
	}
}

func TestReadConditionalLine(t *testing.T) {
	str := "%if 0%{?suse_version} > 1550\n" // 29
	str += RandStringBytes(10)
	str += "\n%endif\n"
	line := NewLine(29, "%if 0%{?suse_version} > 1550\n")
	var c Counter
	err := readConditionalLine(strings.NewReader(str), &line, &c, 0)
	if err != nil || line.Last != "%endif\n" || line.Len != 3 {
		t.Error("[readline]readConditionalLine test failed")
	}
}

func TestReadConditionalLineWithNoBreak(t *testing.T) {
	str := "%if 0%{?suse_version} > 1550\n" // 29
	str += RandStringBytes(10)
	str += "\n%endif"
	line := NewLine(29, "%if 0%{?suse_version} > 1550\n")
	var c Counter
	err := readConditionalLine(strings.NewReader(str), &line, &c, 0)
	if err != io.EOF || line.Last != "%endif" || line.Len != 3 {
		t.Error("[readline]readConditionalLine test failed")
	}
}

func TestWalkFileWithContinue(t *testing.T) {
	var str string
	test := "this is the first line\nthis is the second line\nthis is the last line"
	err := walkFile(strings.NewReader(test), false,
		func(rd io.ReaderAt, line *Line) (error, int64) {
			str1 := strings.Join(line.Lines, "")
			str += str1
			if strings.Contains(str1, "last") {
				return fmt.Errorf("this is a test error"), line.Offset
			}
			return nil, line.Offset
		})
	if err != nil {
		t.Errorf("[readline]walkFile should return nil but the error is %s", err)
	}
	if str != test {
		t.Errorf("[readline]walkFile failed, expected read content %s, got %s", test, str)
	}
}

func TestWalkFileWithBreak(t *testing.T) {
	var str string
	test := "this is the first line\nthis is the second line\nthis is the last line"
	result := "this is the first line\nthis is the second line\n"
	err := walkFile(strings.NewReader(test), true,
		func(rd io.ReaderAt, line *Line) (error, int64) {
			str1 := strings.Join(line.Lines, "")
			if strings.Contains(str1, "last") {
				return fmt.Errorf("this is a test error"), line.Offset
			}
			str += str1
			return nil, line.Offset
		})
	if err != nil {
		t.Errorf("[readline]walkFile should return nil but the error is %s", err)
	}
	if str != result {
		t.Errorf("[readline]walkFile failed, expected read content '%s', got '%s'", result, str)
	}
}
