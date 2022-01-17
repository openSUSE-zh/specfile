package specfile

import (
	"bytes"
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

func isPart(node Node) bool {
	if _, ok := partMap[string(node.typ)]; ok {
		return true
	}
	for k := range partMap {
		if bytes.HasPrefix(node.val, []byte(k)) {
			return true
		}
	}
	return false
}

type Specfile struct {
	tags   []int
	macros []int
	parts  []int
	src    Nodes
}

func NewParser(nodes Nodes) Specfile {
	var specfile Specfile
	var parent int
	var children []int
	var macroundone bool

	specfile.src = nodes

	for _, v := range specfile.src {
		if v.ptr.typ != 0 {
			continue
		}

		if v.macro {
			if v.macrodone {
				if macroundone {
					children = append(children, v.id)
					specfile.src[v.id].parent = parent
				}
			} else {
				parent = v.id
				macroundone = true
			}
			specfile.macros = append(specfile.macros, v.id)
			continue
		}

		if macroundone {
			if bytes.HasSuffix(bytes.TrimRight(v.val, "}"), []byte("nil")) {
				macroundone = false
				children = append(children, v.id)
				specfile.src[parent].children = dup(children)
				specfile.src[v.id].parent = parent
				children = []int{}
				parent = -1
				continue
			}
			children = append(children, v.id)
			specfile.src[v.id].parent = parent
			continue
		}
	}

	for _, k := range specfile.macros {
		v := specfile.src[k]
		fmt.Printf("node id %d, parent %d, children %v, condition %d, address %p, type \"%s\", value \"%s\", isMacro %t, isMacroDone %t\n", v.id, v.parent, v.children, v.condition, v.ptr, v.typ, v.val, v.macro, v.macrodone)
	}
	return specfile
}
