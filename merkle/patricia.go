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
	if n.Name[(p.Level)/8]&(1<<(7-((p.Level)%8))) != 0 {
		p.RVal = nil
		if p.R == nil {
			p.R = n
			n.Parent = p
		}
		return p.R
	}
	p.LVal = nil
	if p.L == nil {
		p.L = n
		n.Parent = p
	}
	return p.L
}

// ShiftDown comment
func ShiftDown(n *Node) *Node {
	intNode := Node{Level: n.Level, Parent: n.Parent}
	if n.IsR() {
		n.Parent.R = &intNode
		n.Parent.RVal = nil
	} else {
		n.Parent.L = &intNode
		n.Parent.LVal = nil
	}
	n.InsertBelow(&intNode)
	return &intNode
}
