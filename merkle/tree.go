package merkle

// Tree comment
type Tree struct {
	Root *Node
}

// Store comment
func (t *Tree) Store(s *Store) {
	//s.Save(t)
	walkSave(t.Root, s)
}

func walkSave(n *Node, s *Store) {
	if n != nil {
		//s.Save(n)
		walkSave(n.R, s)
		walkSave(n.L, s)
	}
}

// Leaves comment
func (n *Node) Leaves(l []*Node) []*Node {
	if n.R == nil && n.L == nil {
		return append(l, n)
	} else if n.L != nil {
		l = append(l, n.L.Leaves(l)...)
	} else if n.R != nil {
		l = append(l, n.R.Leaves(l)...)
	}
	return l
}

// CountLeaves comment
func (n *Node) CountLeaves() int {
	if n == nil {
		return 0
	} else if n.IsLeaf() {
		return 1
	} else {
		return n.L.CountLeaves() + n.R.CountLeaves()
	}
}
