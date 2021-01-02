package specfile

import (
	"testing"
)

func TestSectionParse(t *testing.T) {
	token := NewTokenizer("Section", "%post -n fcitx5-configtool -p /sbin/ldconfig\n")
	var s Section
	(&s).Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Belongs != "fcitx5-configtool" || s.Value != "-p /sbin/ldconfig" {
		t.Error("[section]parse single line test failed")
	}
}

func TestSectionParseMultiLine(t *testing.T) {
	token := NewTokenizer("Section", "%post -n fcitx5-configtool\n/sbin/ldconfig\n")
	var s Section
	(&s).Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Belongs != "fcitx5-configtool" || s.Value != "/sbin/ldconfig" {
		t.Error("[section]parse multiple line test failed")
	}
}

func TestSectionParseWithNoBelongs(t *testing.T) {
	token := NewTokenizer("Section", "%post -p /sbin/ldconfig\n")
	var s Section
	(&s).Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Value != "-p /sbin/ldconfig" {
		t.Error("[section]parse single line with no belongs test failed")
	}
}

func TestSectionParseMultiLineWithNoBelongs(t *testing.T) {
	token := NewTokenizer("Section", "%post\n/sbin/ldconfig\n")
	var s Section
	(&s).Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Value != "/sbin/ldconfig" {
		t.Error("[section]parse multiple line test failed")
	}
}
