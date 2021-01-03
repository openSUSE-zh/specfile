package specfile

// Tag normal Tag like "Name: fcitx"
type Tag struct {
	item
}

func (t *Tag) Parse(token *Tokenizer) {
	t.item.Parse(token)
}
