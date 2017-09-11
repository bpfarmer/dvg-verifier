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
func (n *Node) HashVal() string {
	h := sha256.New()
	if n.IsLeaf() {
		return n.Val
	}
	if n.Deleted {
		h.Write([]byte("EMPTY NODE"))
		return hex.EncodeToString(h.Sum(nil))
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
func (n *Node) InclusionProof(s *Store) []string {
	var p []string
	c := n
	log.Println("Node.InclusionProof():val=" + c.Val)
	for c.PEntry(s) != nil {
		log.Println("Node.InclusionProof():traversing up the tree: val=" + c.Val)
		if c.PEntry(s).LVal == c.HashVal() {
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
	return c.HashVal()
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
