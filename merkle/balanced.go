package merkle

// AppendLeaves Append Leaves to a balanced tree
func (t *Tree) AppendLeaves(m []*Node) {
	l := t.Root.Leaves(nil)
	l = append(l, m...)
	BuildTree(l)
}

// BuildTree comment
func BuildTree(l []*Node) []*Node {
	// Odd number of nodes case
	if len(l)%2 != 0 {
		l = append(l, &Node{})
	}
	// Calculate the largest balanced subtree
	var k uint
	for n := 0; 1<<uint(n) <= len(l); n++ {
		k = uint(n)
	}
	left := BuildSubTree(l[:1<<uint(k)])
	// Calculate the unbalanced righthand side
	if len(left) < len(l) {
		right := l[1<<uint(k):]
		for len(right) > 1 {
			right = BuildTree(right)
		}
		left = append(left, right...)
	}
	return BuildSubTree(left)
}

// BuildSubTree for balanced trees with 2^n leaves
func BuildSubTree(l []*Node) []*Node {
	if !PowOf2(len(l)) {
		return nil
	}
	var parents []*Node
	for n := 0; n < len(l); n += 2 {
		if l[n+1].Parent == nil {
			p := &Node{L: l[n], R: l[n+1]}
			l[n].Parent = p
			l[n+1].Parent = p
		}
		parents = append(parents, l[n+1].Parent)
	}
	if len(parents) > 1 {
		return BuildSubTree(parents)
	}
	return parents
}

// PowOf2 comment
func PowOf2(n int) bool {
	return n != 0 && n&(n-1) == 0
}
