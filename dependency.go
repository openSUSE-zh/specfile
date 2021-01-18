package specfile

// Dependency the dependencies
type Dependency struct {
	item
}

func (d *Dependency) Parse(token *Tokenizer) {
	d.item.Parse(token)
}

func ParseDependency(token, last Tokenizer, spec *Specfile) {
	var dependency Dependency
	if last.Type == "Comment" {
		dependency.Comment = last.Content
	}
	dependency.Parse(&token)
	spec.append("Dependencies", dependency)
}
