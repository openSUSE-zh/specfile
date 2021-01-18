package specfile

// Tag normal Tag like "Name: fcitx"
type Tag struct {
	item
}

func (t *Tag) Parse(token *Tokenizer) {
	t.item.Parse(token)
}

func ParseTag(token, last Tokenizer, spec *Specfile) {
	var t Tag
	if last.Type == "Comment" {
		t.Comment = last.Content
	}
	t.Parse(&token)
	spec.append("Tags", t)
}
