package specfile

import (
	"strings"
	"testing"
)

func TestSectionParse(t *testing.T) {
	token := NewTokenizer("Section", "%post -n fcitx5-configtool -p /sbin/ldconfig\n")
	var s Section
	s.Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Belongs != "fcitx5-configtool" || s.Value != "-p /sbin/ldconfig\n" {
		t.Error("[section]parse single line test failed")
	}
}

func TestSectionParseMultiLine(t *testing.T) {
	token := NewTokenizer("Section", "%post -n fcitx5-configtool\n/sbin/ldconfig\n")
	var s Section
	s.Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Belongs != "fcitx5-configtool" || s.Value != "/sbin/ldconfig\n" {
		t.Error("[section]parse multiple line test failed")
	}
}

func TestSectionParseWithNoBelongs(t *testing.T) {
	token := NewTokenizer("Section", "%post -p /sbin/ldconfig\n")
	var s Section
	s.Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Value != "-p /sbin/ldconfig\n" {
		t.Error("[section]parse single line with no belongs test failed")
	}
}

func TestSectionParseMultiLineWithNoBelongs(t *testing.T) {
	token := NewTokenizer("Section", "%post\n/sbin/ldconfig\n")
	var s Section
	s.Parse(&token)
	if s.Raw != &token || s.Name != "%post" || s.Value != "/sbin/ldconfig\n" {
		t.Error("[section]parse multiple line test failed")
	}
}

func TestParseSection(t *testing.T) {
	var spec Specfile
	token := NewTokenizer("Section", "%description\nCore package for the GNU Compiler Collection, including the C language frontend.\nLanguage frontends other than C are split to different sub-packages, namely gcc-ada, gcc-c++, gcc-fortran, gcc-obj, gcc-obj-c++ and gcc-go.\n")
	last := NewTokenizer("Comment", "This is a comment.")
	ParseSection(token, last, &spec)
	if len(spec.Sections) != 1 || spec.Sections[0].Name != "%description" || spec.Sections[0].Value != strings.TrimPrefix(token.Content, "%description\n") || spec.Sections[0].Comment != last.Content {
		t.Error("[section]ParseSection test failed")
	}
}

func TestParseSectionWithSubpackage(t *testing.T) {
	var spec Specfile
	token := NewTokenizer("Section", "%package -n gcc10-testresults\nSummary:        Testsuite results\nLicense:        SUSE-Public-Domain\nGroup:          Development/Languages/C and C++\n")
	ParseSection(token, NewTokenizer("Empty", "\n"), &spec)
	if len(spec.Subpackages) != 1 {
		t.Errorf("[section]ParseSection with subpackage test failed, expected 1 subpackage, got %d\n", len(spec.Subpackages))
	}
	tag, err := spec.Subpackages[0].FindTag("License")
	if err != nil {
		t.Errorf("[section]ParseSection with subpacakge test failed, expected nil error, got %s\n", err)
	}
	if tag.Value != "SUSE-Public-Domain" {
		t.Errorf("[section]ParseSection with subpackage test failed, expected SUSE-Public-Domain, got %s", tag.Value)
	}
	tag, err = spec.Subpackages[0].FindTag("Name")
	if err != nil {
		t.Errorf("[section]ParseSection with subpackage test failed, expected Name to be inserted, but got error %s", err)
	}
	if tag.Value != "gcc10-testresults" {
		t.Errorf("[section]ParseSection with subpackage test failed, expected Name gcc10-testsuite, got %s", tag.Value)
	}
}

func TestParseSectionWithSubpackageSection(t *testing.T) {
	var spec Specfile
	token := NewTokenizer("Section", "%package -n gcc10-testresults\nSummary:        Testsuite results\nLicense:        SUSE-Public-Domain\nGroup:          Development/Languages/C and C++\n")
	ParseSection(token, NewTokenizer("Empty", "\n"), &spec)
	token1 := NewTokenizer("Section", "%description -n gcc10-testresults\nResults from running the gcc and target library testsuites.")
	ParseSection(token1, NewTokenizer("Empty", "\n"), &spec)
	_, err := spec.Subpackages[0].FindSection("%description")
	if err != nil {
		t.Errorf("[section]ParseSection with subpackage section test failed, expected section found, got %s error", err)
	}
}
