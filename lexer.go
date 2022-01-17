package specfile

import (
	"bytes"
	"fmt"
)

// Node a node
type Node struct {
	id        int
	parent    int
	children  []int
	condition []int
	ptr       *Line
	typ       []byte
	val       []byte
	macro     bool
	macrodone bool
}

func (node Node) isNil() bool {
	if node.ptr == nil {
		return true
	}
	return false
}

// Nodes collection of nodes
type Nodes []Node

// String debug nodes
func (nodes Nodes) String() string {
	var str string
	for _, v := range nodes {
		str += fmt.Sprintf("node id %d, parent %d, children %v, condition %d, address %p, type \"%s\", value \"%s\", isMacro %t, isMacroDone %t\n", v.id, v.parent, v.children, v.condition, v.ptr, v.typ, v.val, v.macro, v.macrodone)
	}
	return str
}

func isMacro(b, b1 []byte) (bool, bool) {
	if bytes.HasPrefix(b, []byte("%bcond")) {
		return true, true
	}

	var macro bool

	if string(b) == "%define" || string(b) == "%global" {
		macro = true
	}

	if !macro {
		return false, false
	}

	// %define/%global <macro> <body>
	i := bytes.IndexByte(b1, ' ')
	if i < 0 {
		return true, false
	}

	for j := i; j < len(b1); j++ {
		if b1[j] != ' ' && b1[j] != '\\' && b1[j] != '\n' {
			return true, true
		}
	}

	return true, false
}

func NewLexer(lines Lines) Nodes {
	var nodes Nodes
	var idx int
	var condition []int
	for k := range lines {
		if lines[k].typ < 0 {
			continue
		}

		if lines[k].buf[0] == '%' {
			i := bytes.IndexByte(lines[k].buf, ' ')
			if i >= 0 {
				// the if, elif and else case
				if lines[k].typ > 0 && lines[k].typ < 10 {
					condition = append(condition, idx)
					nodes = append(nodes, Node{idx, -1, []int{}, []int{}, &lines[k], lines[k].buf[:i], bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " "), false, false})
				} else {
					val := bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " ")
					macro, done := isMacro(lines[k].buf[:i], val)
					nodes = append(nodes, Node{idx, -1, []int{}, dup(condition), &lines[k], lines[k].buf[:i],
						val, macro, done})
				}
			} else {
				switch lines[k].typ {
				case 6:
					// the else case
					condition = append(condition, idx)
					nodes = append(nodes, Node{idx, -1, []int{}, []int{}, &lines[k], lines[k].buf[:lines[k].len-1], []byte{}, false, false})
				case 10:
					// the endif case
					// reverse delete
					for i := len(condition) - 1; i >= 0; i-- {
						if nodes[condition[i]].ptr.typ < 6 {
							condition = condition[:i]
							break
						}
					}
					nodes = append(nodes, Node{idx, -1, []int{}, []int{}, &lines[k], lines[k].buf[:lines[k].len-1], []byte{}, false, false})
				default:
					nodes = append(nodes, Node{idx, -1, []int{}, []int{}, &lines[k], []byte{}, lines[k].buf[:lines[k].len-1], false, false})
				}
			}
			idx++
			continue
		}

		if lines[k].buf[0] >= 'A' {
			i := bytes.IndexByte(lines[k].buf, ':')
			if i >= 0 {
				nodes = append(nodes, Node{idx, -1, []int{}, dup(condition), &lines[k], lines[k].buf[:i], bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " "), false, false})
				idx++
				continue
			}
		}
		nodes = append(nodes, Node{idx, -1, []int{}, dup(condition), &lines[k], []byte{}, lines[k].buf[:lines[k].len-1], false, false})
		idx++
	}
	return nodes
}

func dup(b []int) (b1 []int) {
	b1 = make([]int, len(b))
	copy(b1, b)
	return b1
}
