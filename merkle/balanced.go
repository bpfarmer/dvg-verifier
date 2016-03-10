package merkle

// AppendLeaves Append Leaves to a balanced tree
func AppendLeaves(l []*Node, m []*Node) []*Node {
	l = append(l, m...)
	return BuildTree(l)
}

// BuildTree for unbalanced trees with 2n leaves
func BuildTree(l []*Node) []*Node {
	var k uint
	for n := 0; 1<<uint(n) <= len(l); n++ {
		k = uint(n)
	}
	left := BuildSubTree(l[:1<<uint(k)])
	// Calculate the unbalanced righthand side
	if len(l[1<<uint(k):]) > 0 {
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
	if !PowOf2(len(l)) || len(l) == 1 {
		return l
	}
	var parents []*Node
	for n := 0; n < len(l); n += 2 {
		if l[n].Parent == nil {
			p := &Node{L: l[n], R: l[n+1]}
			l[n].Parent = p
			l[n+1].Parent = p
		}
		parents = append(parents, l[n].Parent)
	}
	return BuildSubTree(parents)
}

// PowOf2 comment
func PowOf2(n int) bool {
	return n != 0 && n&(n-1) == 0
}
