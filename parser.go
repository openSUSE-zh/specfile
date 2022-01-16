package specfile

import (
  "fmt"
)

var (
	partMap = map[string]struct{}{
		"%package":                {},
		"%prep":                   {},
		"%build":                  {},
		"%install":                {},
		"%check":                  {},
		"%clean":                  {},
		"%files":                  {},
		"%pre":                    {},
		"%post":                   {},
		"%preun":                  {},
		"%postun":                 {},
		"%pretrans":               {},
		"%posttrans":              {},
		"%description":            {},
		"%changelog":              {},
		"%triggerin":              {},
		"%triggerun":              {},
		"%trigger":                {},
		"%verifyscript":           {},
		"%triggerpostun":          {},
		"%triggerprein":           {},
		"%sepolicy":               {},
		"%filetriggerin":          {},
		"%filetrigger":            {},
		"%filetriggerun":          {},
		"%filetriggerpostun":      {},
		"%transfiletriggerin":     {},
		"%transfiletrigger":       {},
		"%transfiletriggerun":     {},
		"%transfiletriggerpostun": {},
		"%patchlist":              {},
		"%sourcelist":             {},
		"%generate_buildrequires": {},
	}
)

type Specfile struct {
	tags         Nodes
	macros       Nodes
	dependencies Nodes
	packages     []Package
	parts        []Nodes
}

type Package struct {
	tags         Nodes
	dependencies Nodes
}

func NewParser(nodes Nodes) Specfile {
	var specfile Specfile
  for _, v := range nodes {
    if v.ptr.typ != 0 {
      continue
    }
    fmt.Printf("node id %d, parent %d, children %v, condition %d, address %p, type \"%s\", value \"%s\"\n", v.id, v.parent, v.children, v.condition, v.ptr, v.typ, v.val)
  }
	return specfile
}
