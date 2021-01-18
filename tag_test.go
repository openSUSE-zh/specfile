package specfile

import "testing"

func TestParseTag(t *testing.T) {
	var spec Specfile
	token := NewTokenizer("Tag", "BuildRequires: xz")
	last := NewTokenizer("Comment", "This is a comment.")
	ParseTag(token, last, &spec)
	if len(spec.Tags) != 1 || spec.Tags[0].Name != "BuildRequires" || spec.Tags[0].Value != "xz" || spec.Tags[0].Comment != last.Content {
		t.Error("[tag]ParseTag test failed")
	}
}
