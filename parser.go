package specfile

type Specfile struct {
	Tags         Nodes
	Macros       Nodes
	Dependencies Nodes
	Packages     []*Specfile
	Parts        Nodes
}

func NewParser(nodes Nodes) Specfile {
	var specfile Specfile
	return specfile
}
