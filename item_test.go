package specfile

import "testing"

func TestItemParse(t *testing.T) {
	token := NewTokenizer("Tag", "Name: fcitx5")
	var i item
	(&i).Parse(&token)
	if i.Raw != &token || i.Name != "Name" || i.Value != "fcitx5" {
		t.Error("[item]parse test failed")
	}
}
