package specfile

// Dependency the dependencies
type Dependency struct {
	item
}

func (d *Dependency) Parse(token *Tokenizer) {
	d.item.Parse(token)
}
