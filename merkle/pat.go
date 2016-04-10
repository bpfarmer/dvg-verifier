package merkle

// InsertLeaf Patricia Leaf
func (t *Tree) InsertLeaf(leaf *Node) {
	curNode := t.Root
	var offset uint
	for curNode != leaf {
		if curNode.IsLeaf() && curNode.Parent != nil {
			curNode = ShiftDown(curNode)
		} else {
			curNode = leaf.InsertBelow(curNode)
		}
		offset++
	}
}

// InsertBelow comment
func (n *Node) InsertBelow(p *Node) *Node {
	n.Level = p.Level + 1
	p.Reset()
	if n.Name[(p.Level)/8]&(1<<(7-((p.Level)%8))) != 0 {
		if p.R == nil {
			p.R = n
			n.SetParent(p)
		}
		return p.R
	}
	if p.L == nil {
		p.L = n
		n.SetParent(p)
	}
	return p.L
}

// ShiftDown comment
func ShiftDown(n *Node) *Node {
	intNode := Node{Level: n.Level, Parent: n.Parent}
	n.Parent.Reset()
	if n.IsR() {
		n.Parent.R = &intNode
	} else {
		n.Parent.L = &intNode
	}
	n.InsertBelow(&intNode)
	return &intNode
}
