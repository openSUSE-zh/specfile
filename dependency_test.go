package specfile

import "testing"

func TestParseDependency(t *testing.T) {
	var spec Specfile
	token := NewTokenizer("Dependency", "BuildRequires: xz")
	last := NewTokenizer("Comment", "This is a comment.")
	ParseDependency(token, last, &spec)
	if len(spec.Dependencies) != 1 || spec.Dependencies[0].Name != "BuildRequires" || spec.Dependencies[0].Value != "xz" || spec.Dependencies[0].Comment != last.Content {
		t.Error("[dependency]ParseDependency test failed")
	}
}
