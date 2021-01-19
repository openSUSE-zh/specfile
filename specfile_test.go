package specfile

import "testing"

func TestFindTag(t *testing.T) {
	var spec Specfile
	var tag Tag
	token := NewTokenizer("Tag", "Name: gcc10\n")
	tag.Parse(&token)
	spec.append("Tags", tag)
	if _, err := spec.FindTag("Name"); err != nil {
		t.Error("[specfile]FindTag test failed")
	}
}

func TestFindSection(t *testing.T) {
	var spec Specfile
	var section Section
	token := NewTokenizer("Section", "%description -n gcc10-testresults\nResults from running the gcc and target library testsuites.")
	section.Parse(&token)
	spec.append("Sections", section)
	if _, err := spec.FindSection("%description"); err != nil {
		t.Error("[specfile]FindSection test failed")
	}
}

func TestFindSubpackage(t *testing.T) {
	var spec, spec1 Specfile
	var tag Tag
	token := NewTokenizer("Tag", "Name: gcc10\n")
	tag.Parse(&token)
	spec1.append("Tags", tag)
	spec.append("Subpackages", spec1)
	if _, err := spec.FindSubpackage("gcc10"); err != nil {
		t.Error("[specfile]FindSubpackage test failed")
	}
}
