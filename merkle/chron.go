package merkle

// AppendLeaves Append Leaves to a balanced tree
func AppendLeaves(l []*Node, m []*Node) []*Node {
	l = append(l, m...)
	return BuildTree(l)
}

// BuildTree for unbalanced trees with 2n leaves
func BuildTree(l []*Node) []*Node {
	if len(l) == 1 {
		return l
	}
	var k uint
	for n := 0; 1<<uint(n) <= len(l); n++ {
		k = uint(n)
	}
	sub := BuildSubTree(l[:1<<uint(k)])
	if len(l[1<<uint(k):]) > 0 {
		sub = append(sub, BuildTree(l[1<<uint(k):])...)
	}
	return BuildSubTree(sub)
}

// BuildSubTree for balanced trees with 2^n leaves
func BuildSubTree(l []*Node) []*Node {
	if !PowOf2(len(l)) || len(l) == 1 {
		return l
	}
	var parents []*Node
	for n := 0; n < len(l); n += 2 {
		if l[n+1].P == nil {
			p := &Node{L: l[n], R: l[n+1]}
			l[n].P = p
			l[n+1].P = p
		}
		parents = append(parents, l[n].P)
	}
	return BuildSubTree(parents)
}

// PowOf2 comment
func PowOf2(n int) bool {
	return n != 0 && n&(n-1) == 0
}
