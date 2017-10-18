package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"
)

// Node comment
type Node struct {
	P, L, R         *Node
	LVal, RVal, Val string
	Epoch           uint
	Deleted         bool
	// For DB purposes, probably unnecessary to include
	ID int
}

// HashVal comment
func (n *Node) HashVal(s *Store) string {
	h := sha256.New()
	if n.IsLeaf() {
		return n.Val
	}
	if len(n.LVal) == 0 && n.L != nil {
		n.LVal = n.LEntry(s).HashVal(s)
	}
	if len(n.RVal) == 0 && n.R != nil {
		n.RVal = n.REntry(s).HashVal(s)
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
func (n *Node) InclusionProof(s *Store) []string {
	var p []string
	c := n
	log.Println("Node.InclusionProof():val=" + c.Val)
	for c.PEntry(s) != nil {
		log.Println("Node.InclusionProof():traversing up the tree: val=" + c.Val)
		if c.PEntry(s).LVal == c.HashVal(s) {
			p = append(p, c.PEntry(s).RVal+"_R")
		} else {
			p = append(p, c.PEntry(s).LVal+"_L")
		}
		c = c.PEntry(s)
	}
	return p
}

// RootHash comment
func (n *Node) RootHash(s *Store) string {
	c := n
	for c.PEntry(s) != nil {
		c = c.PEntry(s)
	}
	return c.HashVal(s)
}

// IsLeaf comment
func (n *Node) IsLeaf() bool {
	return n.L == nil && n.R == nil
}

var leaves []*Node

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

// RMostEntry returns the right-most descendant
func (n *Node) RMostEntry(s *Store) *Node {
	for n.REntry(s) != nil {
		n = n.REntry(s)
	}
	return n
}

// LEntry comment
func (n *Node) LEntry(s *Store) *Node {
	if n.L != nil {
		return n.L
	}
	n.L = FindNode(s, n.LVal)
	if n.L == nil {
		return nil
	}
	n.L.P = n
	return n.L
}

// REntry comment
func (n *Node) REntry(s *Store) *Node {
	if n.R != nil {
		return n.R
	}
	n.R = FindNode(s, n.RVal)
	if n.R == nil {
		return nil
	}
	n.R.P = n
	return n.R
}

// PEntry comment
func (n *Node) PEntry(s *Store) *Node {
	if n.P != nil {
		return n.P
	}
	n.P = FindParent(s, n.Val)
	if n.P == nil {
		return nil
	}
	if n.P.LVal == n.Val {
		n.P.L = n
		return n.P
	}
	n.P.R = n
	return n.P
}

// FindParent comment
func FindParent(s *Store, val string) *Node {
	q := "select * from nodes where l_val = $1 or r_val = $1"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	n := MapToNodes(rows)
	if len(n) > 0 {
		return n[0]
	}
	return nil
}

// FindNode comment
func FindNode(s *Store, val string) *Node {
	q := "select * from nodes where val = $1"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	n := MapToNodes(rows)
	if len(n) > 0 {
		return n[0]
	}
	return nil
}

// FindNodes comment
func FindNodes(s *Store, val string) []*Node {
	q := "select * from nodes where val = $1"
	rows, err := s.DB.Query(q, val)
	if err != nil {
		log.Fatal(err)
	}
	return MapToNodes(rows)
}
