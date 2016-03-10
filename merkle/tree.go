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
