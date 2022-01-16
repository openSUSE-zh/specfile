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
}

// Nodes collection of nodes
type Nodes []Node

// String debug nodes
func (nodes Nodes) String() string {
	var str string
	for _, v := range nodes {
		str += fmt.Sprintf("node id %d, parent %d, children %v, condition %d, address %p, type \"%s\", value \"%s\"\n", v.id, v.parent, v.children, v.condition, v.ptr, v.typ, v.val)
	}
	return str
}

func (nodes Nodes) FindByID(id int) Node {
	for _, v := range nodes {
		if v.id == id {
			return v
		}
	}
	return Node{}
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
					nodes = append(nodes, Node{idx, 0, []int{}, []int{}, &lines[k], lines[k].buf[:i], bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " ")})
				} else {
					c := make([]int, len(condition))
					copy(c, condition)
					nodes = append(nodes, Node{idx, 0, []int{}, c, &lines[k], lines[k].buf[:i], bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " ")})
				}
			} else {
				switch lines[k].typ {
				case 6:
					// the else case
					condition = append(condition, idx)
				case 10:
					// the endif case
					// reverse delete
					for i := len(condition) - 1; i >= 0; i-- {
						if v := nodes.FindByID(condition[i]); v.ptr.typ < 6 {
							condition = condition[:i]
							break
						}
					}
				}
				nodes = append(nodes, Node{idx, 0, []int{}, []int{}, &lines[k], []byte{}, lines[k].buf[:lines[k].len-1]})
			}
			idx++
			continue
		}

		if lines[k].buf[0] >= 'A' {
			i := bytes.IndexByte(lines[k].buf, ':')
			if i >= 0 {
				c := make([]int, len(condition))
				copy(c, condition)
				nodes = append(nodes, Node{idx, 0, []int{}, c, &lines[k], lines[k].buf[:i], bytes.TrimLeft(lines[k].buf[i+1:lines[k].len-1], " ")})
				idx++
				continue
			}
		}
		c := make([]int, len(condition))
		copy(c, condition)
		nodes = append(nodes, Node{idx, 0, []int{}, c, &lines[k], []byte{}, lines[k].buf[:lines[k].len-1]})
		idx++
	}
	return nodes
}
