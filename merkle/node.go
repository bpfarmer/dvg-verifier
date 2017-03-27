package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

// Node comment
type Node struct {
	P, L, R         *Node
	LVal, RVal, Val string
	Epoch           uint
	// For DB purposes, probably unnecessary to include
	ID int
}

// HashVal comment
func (n *Node) HashVal() string {
	h := sha256.New()
	if n.IsLeaf() {
		return n.Val
	}
	if len(n.LVal) == 0 && n.L != nil {
		n.LVal = n.L.HashVal()
	}
	if len(n.RVal) == 0 && n.R != nil {
		n.RVal = n.R.HashVal()
	}
	io.WriteString(h, hashEmpty(n.LVal))
	io.WriteString(h, hashEmpty(n.RVal))
	return hex.EncodeToString(h.Sum(nil))
}

// Reset comment
func (n *Node) Reset() {
	n.LVal = ""
	n.RVal = ""
}

// subHash comment
func hashEmpty(subHash string) string {
	if len(subHash) > 0 {
		return subHash
	}
	h := sha256.New()
	h.Write([]byte("EMPTY NODE"))
	return hex.EncodeToString(h.Sum(nil))
}

// InclusionProof comment
func (n *Node) InclusionProof() []string {
	var p []string
	c := n
	fmt.Println("--INCLUSION PROOF--")
	fmt.Println(c.Val)
	for c.P != nil {
		if c.P.LVal == c.HashVal() {
			p = append(p, c.P.RVal+"_R")
		} else {
			p = append(p, c.P.LVal+"_L")
		}
		c = c.P
	}
	return p
}

// RootHash comment
func (n *Node) RootHash() string {
	c := n
	for c.P != nil {
		c = c.P
	}
	return c.HashVal()
}

// IsLeaf comment
func (n *Node) IsLeaf() bool {
	return n.L == nil && n.R == nil
}

var leaves []*Node

/*
// Leaves comment
func (n *Node) Leaves() []*Node {
	leaves = []*Node{}
	return n.leaves()
}

// Leaves comment
func (n *Node) leaves() []*Node {
	if n.IsLeaf() {
		leaves = append(leaves, n)
	} else {
		if n.L != nil {
			n.L.leaves()
		}
		if n.R != nil {
			n.R.leaves()
		}
	}
	return leaves
}*/

// Path comment
func (n *Node) Path(s *Store) []*Node {
	var path []*Node
	return n.calculatePath(path, s)
}

func (n *Node) calculatePath(nodes []*Node, s *Store) []*Node {
	nodes = append(nodes, n)
	if n.PEntry(s) != nil {
		nodes = n.PEntry(s).calculatePath(nodes, s)
	}
	return nodes
}

/*
// Addr comment
func (n *Node) Addr() string {
	if n.Parent != nil {
		return n.Parent.Addr() + "/" + strconv.Itoa(n.Parent.ID)
	}
	return ""
}

/*
// FindPath comment
func (n *Node) FindPath(s *Store) []*Node {
	ids := strings.Split(n.Path, "/")
	ids = ids[1:]
	q := "select * from nodes where id = any($1::integer[])"
	rows, err := s.DB.Query(q, fmt.Sprintf("{%s}", strings.Join(ids, ", ")))
	if err != nil {
		log.Fatal(err)
	}
	path := MapToNodes(rows)
	path = append(path, n)
	AssocNodes(path)
	OrderPath(path)
	return path
}

// OrderPath comment
func OrderPath(n []*Node) []*Node {
	var leaf int
	for p, v := range n {
		if len(v.Val) > 0 {
			leaf = p
		}
	}
	var o []*Node
	c := n[leaf]
	for c != nil {
		o = append(o, c)
		c = c.Parent
	}
	return o
}
*/
